/*
Copyright 2023 The KubeFin Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/values"
)

const (
	ackClusterIdLabelKey  = "ack.aliyun.com"
	ackNodeTypeLabelKey   = "node.kubernetes.io/instance-type"
	ackNodeRegionLabelKey = "topology.kubernetes.io/region"

	nodePriceQueryUrl = "https://buy-api.aliyun.com/price/getLightWeightPrice2.json?tenant=TenantCalculator"
	nodeSpecQueryUrl  = "https://query.aliyun.com/rest/sell.ecs.allInstanceTypes?domain=aliyun&saleStrategy=PostPaid"

	defaultCPUMemoryCostRatio = 3.0
)

type AckCloudProvider struct {
	client kubernetes.Interface

	cpuMemoryCostRatio float64

	// nodePriceMap maps [region name][node type]price
	nodePriceMap  map[string]map[string]float64
	nodePriceLock sync.Mutex

	// nodeSpecMap maps [node type]NodeSpec
	nodeSpecMap  map[string]NodeSpec
	nodeSpecLock sync.Mutex
}

func NewAckCloudProvider(client kubernetes.Interface, agentOptions *options.AgentOptions) (*AckCloudProvider, error) {
	var err error

	cpuMemoryCostRatio := defaultCPUMemoryCostRatio
	if agentOptions.CPUMemoryCostRatio != "" {
		cpuMemoryCostRatio, err = strconv.ParseFloat(agentOptions.CPUMemoryCostRatio, 64)
		if err != nil {
			return nil, err
		}
	}
	ackCloud := AckCloudProvider{
		client:             client,
		cpuMemoryCostRatio: cpuMemoryCostRatio,
		nodePriceMap:       map[string]map[string]float64{},
		nodeSpecMap:        map[string]NodeSpec{},
	}

	return &ackCloud, nil
}

func (c *AckCloudProvider) ParseClusterInfo(agentOptions *options.AgentOptions) error {
	// We have no way to get the cluster name currently
	if agentOptions.ClusterName == "" {
		return fmt.Errorf("please set the cluster name via env CLUSTER_NAME in agent manifest")
	}

	if agentOptions.ClusterId != "" {
		return nil
	}

	nodes, err := c.client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodes.Items {
		labels := node.Labels
		if clusterId, ok := labels[ackClusterIdLabelKey]; ok {
			agentOptions.ClusterId = clusterId
			break
		}
	}

	return nil
}

func (c *AckCloudProvider) GetNodeHourlyPrice(node *v1.Node) (*api.InstancePriceInfo, error) {
	if node.Labels == nil {
		return nil, fmt.Errorf("node(%s) has no labels", node.Name)
	}

	nodeRegion, ok := node.Labels[ackNodeRegionLabelKey]
	if !ok {
		return nil, fmt.Errorf("node(%s) has no label %s", node.Name, ackNodeRegionLabelKey)
	}
	nodeType, ok := node.Labels[ackNodeTypeLabelKey]
	if !ok {
		return nil, fmt.Errorf("node(%s) has no label %s", node.Name, ackNodeTypeLabelKey)
	}

	var err error

	c.nodeSpecLock.Lock()
	defer c.nodeSpecLock.Unlock()
	nodeSpec, ok := c.nodeSpecMap[nodeType]
	if !ok {
		err = queryNodeSpecFromCloud(c.nodeSpecMap)
		if err != nil {
			klog.Errorf("Query AliCloud ecs spec error:%v", err)
			return nil, err
		}
		if nodeSpec, ok = c.nodeSpecMap[nodeType]; !ok {
			klog.Errorf("Could not find node type:%s", nodeType)
			return nil, fmt.Errorf("could not find node type:%s", nodeType)
		}
	}

	c.nodePriceLock.Lock()
	defer c.nodePriceLock.Unlock()
	_, ok = c.nodePriceMap[nodeRegion]
	if !ok {
		c.nodePriceMap[nodeRegion] = map[string]float64{}
	}
	nodePrice, ok := c.nodePriceMap[nodeRegion][nodeType]
	if !ok {
		nodePrice, err = queryNodePriceFromCloud(nodeRegion, nodeType)
		if err != nil {
			klog.Errorf("Query AliCloud ecs price error:%v", err)
			return nil, err
		}
		c.nodePriceMap[nodeRegion][nodeType] = nodePrice
	}
	price := c.nodePriceMap[nodeRegion][nodeType]
	return &api.InstancePriceInfo{
		NodeTotalHourlyPrice: price,
		CPUCore:              nodeSpec.CPUCount,
		CPUCoreHourlyPrice:   nodePrice * c.cpuMemoryCostRatio / (c.cpuMemoryCostRatio + 1),
		RamGiB:               nodeSpec.RAMGBCount,
		RAMGBHourlyPrice:     nodePrice / (c.cpuMemoryCostRatio + 1),
		InstanceType:         nodeType,
		BillingMode:          values.BillingModeOnDemand,
		BillingPeriod:        0,
		Region:               nodeRegion,
		CloudProvider:        api.CloudProviderAck,
	}, nil
}

func queryNodePriceFromCloud(nodeRegion, nodeType string) (float64, error) {
	queryPara := newNodePriceQueryPara(nodeRegion, nodeType)
	jsonData, err := json.Marshal(queryPara)
	if err != nil {
		klog.Errorf("Marshal TenantCalculator error:%v", err)
		return 0, err
	}

	resp, err := http.Post(nodePriceQueryUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		klog.Errorf("Query AliCloud ecs price error:%v", err)
		return 0, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Read response body error:%v", err)
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		klog.Errorf("Query AliCloud ecs price error:%s", string(data))
		return 0, fmt.Errorf("query AliCloud ecs price error:%s", string(data))
	}

	priceResult := &TenantCalculatorResult{}
	err = json.Unmarshal(data, priceResult)
	if err != nil {
		klog.Errorf("Unmarshal response body error:%v", err)
		return 0, err
	}

	return priceResult.Data.Order.TradeAmount, nil
}

func queryNodeSpecFromCloud(nodeSpec map[string]NodeSpec) error {
	resp, err := http.Get(nodeSpecQueryUrl)
	if err != nil {
		klog.Errorf("Query AliCloud ecs spec error:%v", err)
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Read response body error:%v", err)
		return err
	}

	queryResult := &NodeSpecQueryResult{}
	err = json.Unmarshal(data, queryResult)
	if err != nil {
		klog.Errorf("Unmarshal response body error:%v", err)
		return err
	}

	for _, instance := range queryResult.Data.Components.InstanceType.InstanceType {
		nodeType := instance.InstanceTypeId
		if _, ok := nodeSpec[nodeType]; !ok {
			cpuCount, err := strconv.ParseFloat(instance.CPUCoreCount, 64)
			if err != nil {
				klog.Errorf("Can not parse cpu count(%s):%v", instance.CPUCoreCount, err)
			}
			memoryCount, err := strconv.ParseFloat(instance.MemorySize, 64)
			if err != nil {
				klog.Errorf("Can not parse memory count(%s):%v", instance.MemorySize, err)
				return err
			}
			nodeSpec[nodeType] = NodeSpec{
				CPUCount:   cpuCount,
				RAMGBCount: memoryCount,
			}
		}
	}
	return nil
}

func newNodePriceQueryPara(instanceRegion, instanceType string) *TenantCalculator {
	return &TenantCalculator{
		Tenant: "TenantCalculator",
		Configurations: []TenantCalculatorConfiguration{
			{
				CommodityCode:   "ecs",
				SpecCode:        "ecs",
				ChargeType:      "POSTPAY",
				OrderType:       "BUY",
				Quantity:        1,
				Duration:        1,
				PricingCycle:    "Hour",
				UseTimeUnit:     "Hour",
				UseTimeQuantity: 1,
				Components: []TenantCalculatorComponent{
					{
						ComponentCode: "vm_region_no",
						InstanceProperty: []TenantCalculatorInstanceProperty{
							{Code: "vm_region_no", Value: instanceRegion},
						},
					},
					{
						ComponentCode: "instance_type",
						InstanceProperty: []TenantCalculatorInstanceProperty{
							{Code: "instance_type", Value: instanceType},
						},
					},
				},
			},
		},
	}
}
