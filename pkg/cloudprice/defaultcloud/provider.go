/*
Copyright 2022 The KubeFin Authors

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

package defaultcloud

import (
	"context"
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/values"
)

const (
	defaultCpuCoreHourlyPrice = 0.08
	defaultRamGBHourlyPrice   = 0.02
)

type DefaultCloudProvider struct {
	client             kubernetes.Interface
	CpuCoreHourlyPrice float64
	RamGBHourlyPrice   float64
	CPUCoreReserved    string
	RAMGBReserved      string
}

func NewDefaultCloudProvider(client kubernetes.Interface, agentOptions *options.AgentOptions) (*DefaultCloudProvider, error) {
	var err error

	cpuCoreHourlyPrice, ramGBHourlyPrice := defaultCpuCoreHourlyPrice, defaultRamGBHourlyPrice
	if agentOptions.CustomCPUCoreHourPrice != "" {
		cpuCoreHourlyPrice, err = strconv.ParseFloat(agentOptions.CustomCPUCoreHourPrice, 64)
		if err != nil {
			return nil, err
		}
	} else {
		cpuCoreHourlyPrice = defaultCpuCoreHourlyPrice
	}

	if agentOptions.CustomRAMGBHourPrice != "" {
		ramGBHourlyPrice, err = strconv.ParseFloat(agentOptions.CustomRAMGBHourPrice, 64)
		if err != nil {
			return nil, err
		}
	} else {
		ramGBHourlyPrice = defaultRamGBHourlyPrice
	}

	defaultCloud := DefaultCloudProvider{
		client:             client,
		CpuCoreHourlyPrice: cpuCoreHourlyPrice,
		RamGBHourlyPrice:   ramGBHourlyPrice,
		CPUCoreReserved:    agentOptions.CPUCoreReserved,
		RAMGBReserved:      agentOptions.RAMGBReserved,
	}

	return &defaultCloud, nil
}

func (c *DefaultCloudProvider) ParseClusterInfo(agentOptions *options.AgentOptions) error {
	if agentOptions.ClusterName == "" {
		return fmt.Errorf("please set the cluster name via env CLUSTER_NAME in agent manifest")
	}

	if agentOptions.ClusterId != "" {
		return nil
	}

	systemNS, err := c.client.CoreV1().Namespaces().Get(context.Background(), metav1.NamespaceSystem, metav1.GetOptions{})
	if err != nil {
		return err
	}
	agentOptions.ClusterId = string(systemNS.UID)

	return nil
}

func (c *DefaultCloudProvider) GetNodeHourlyPrice(node *v1.Node) (*api.InstancePriceInfo, error) {
	cpuCoresQuantity := node.Status.Capacity[v1.ResourceCPU]
	ramBytesQuantity := node.Status.Capacity[v1.ResourceMemory]
	cpuCores := cpuCoresQuantity.AsApproximateFloat64()
	ramBytes := ramBytesQuantity.AsApproximateFloat64()

	// we can't get real cpu/ram from node.status, take env and add it
	var err error
	CPUCoreReserved := 0.0
	if c.CPUCoreReserved != "" {
		CPUCoreReserved, err = strconv.ParseFloat(c.CPUCoreReserved, 64)
		if err != nil {
			klog.Errorf("Parse %s error:%v", c.CPUCoreReserved, err)
			return nil, err
		}
	}

	RAMGBReserved := 0.0
	if c.RAMGBReserved != "" {
		RAMGBReserved, err = strconv.ParseFloat(c.RAMGBReserved, 64)
		if err != nil {
			klog.Errorf("Parse %s error:%v", c.RAMGBReserved, err)
			return nil, err
		}
	}

	return &api.InstancePriceInfo{
		NodeTotalHourlyPrice: c.CpuCoreHourlyPrice*cpuCores + c.RamGBHourlyPrice*(ramBytes/values.GBInBytes),
		CPUCore:              cpuCores + CPUCoreReserved,
		CPUCoreHourlyPrice:   c.CpuCoreHourlyPrice,
		RamGB:                (ramBytes / values.GBInBytes) + RAMGBReserved,
		RAMGBHourlyPrice:     c.RamGBHourlyPrice,
		InstanceType:         "default_instance_type",
		BillingMode:          values.BillingModeOnDemand,
		BillingPeriod:        0,
		Region:               "default_region",
		CloudProvider:        api.CloudProviderDefault,
	}, nil
}
