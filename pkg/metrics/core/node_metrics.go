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

package core

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/utils"
	"github.com/kubefin/kubefin/pkg/values"
)

type nodeResourceInfo struct {
	allocatableResource corev1.ResourceList
	requestedResource   corev1.ResourceList
}

type NodeLevelMetricsCollector struct {
	metricsClient *versioned.Clientset
	provider      cloudprice.CloudProviderInterface
	nodeLister    v1.NodeLister

	mutex        sync.Mutex
	nodeResouece map[string]nodeResourceInfo

	nodeCPUCoreCostGV             *prometheus.GaugeVec
	nodeRAMGBCostGV               *prometheus.GaugeVec
	nodeResourceHourlyTotalCostGV *prometheus.GaugeVec
	nodeTotalCostGV               *prometheus.GaugeVec

	nodeResourceTotalGV       *prometheus.GaugeVec
	nodeResourceSystemTakenGV *prometheus.GaugeVec
	nodeResourceAvailableGV   *prometheus.GaugeVec
	nodeResourceUsageGV       *prometheus.GaugeVec
}

func NewNodeLevelMetricsCollector(client *versioned.Clientset, provider cloudprice.CloudProviderInterface,
	coreResourceInformerLister *api.CoreResourceInformerLister) *NodeLevelMetricsCollector {
	metricsCostLabelKey := []string{
		values.NodeNameLabelKey,
		values.NodeInstanceTypeLabelKey,
		values.BillingModeLabelKey,
		values.NodeBillingPeriodLabelKey,
		values.RegionLabelKey,
		values.CloudProviderLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
	}
	nodeCPUCoreHourlyCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeCPUCoreHourlyCostMetricsName,
		Help: "The node hourly cpu-core cost for the node"}, metricsCostLabelKey)
	nodeRAMGBHourlyCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeRAMGBHourlyCostMetricsName,
		Help: "The node hourly ram-gb cost for the node"}, metricsCostLabelKey)
	nodeTotalCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeTotalHourlyCostMetricsName,
		Help: "The node total hourly cost for the node"}, metricsCostLabelKey)

	metricsCostUnifiedLabelKey := []string{
		values.NodeNameLabelKey,
		values.NodeInstanceTypeLabelKey,
		values.BillingModeLabelKey,
		values.NodeBillingPeriodLabelKey,
		values.RegionLabelKey,
		values.CloudProviderLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.ResourceTypeLabelKey,
	}
	nodeResourceHourlyCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeResourceHourlyCostMetricsName,
		Help: "The node hourly cpu/ram(total cores) cost for the node"}, metricsCostUnifiedLabelKey)

	resourceMetricsLabelKey := []string{
		values.NodeNameLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.ResourceTypeLabelKey,
		values.BillingModeLabelKey,
	}

	nodeResourceTotalGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeResourceTotalMetricsName,
		Help: "The total node resource for the node"}, resourceMetricsLabelKey)
	nodeResourceSystemTakenGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeResourceSystemTakenName,
		Help: "The total node resoruce taken by system"}, resourceMetricsLabelKey)
	nodeResourceAvailableGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeResourceAvailableMetricsName,
		Help: "The node resource allocatable for the node"}, resourceMetricsLabelKey)
	nodeResourceUsageGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.NodeResourceUsageMetricsName,
		Help: "The node resource usage for the node"}, resourceMetricsLabelKey)

	prometheus.MustRegister(nodeCPUCoreHourlyCostGV,
		nodeRAMGBHourlyCostGV, nodeResourceHourlyCostGV, nodeResourceUsageGV,
		nodeTotalCostGV, nodeResourceTotalGV, nodeResourceSystemTakenGV, nodeResourceAvailableGV)

	nodeMetricsCollector := &NodeLevelMetricsCollector{
		metricsClient:                 client,
		provider:                      provider,
		nodeLister:                    coreResourceInformerLister.NodeLister,
		nodeResouece:                  make(map[string]nodeResourceInfo),
		nodeCPUCoreCostGV:             nodeCPUCoreHourlyCostGV,
		nodeRAMGBCostGV:               nodeRAMGBHourlyCostGV,
		nodeResourceHourlyTotalCostGV: nodeResourceHourlyCostGV,
		nodeTotalCostGV:               nodeTotalCostGV,
		nodeResourceTotalGV:           nodeResourceTotalGV,
		nodeResourceSystemTakenGV:     nodeResourceSystemTakenGV,
		nodeResourceAvailableGV:       nodeResourceAvailableGV,
		nodeResourceUsageGV:           nodeResourceUsageGV,
	}
	nodeMetricsCollector.registerNodeResourceEventHandler(coreResourceInformerLister)

	return nodeMetricsCollector
}

func (n *NodeLevelMetricsCollector) registerNodeResourceEventHandler(coreResourceInformerLister *api.CoreResourceInformerLister) {
	coreResourceInformerLister.NodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if !ok {
				return
			}
			n.handleNodeAddition(node)
		},
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if !ok {
				return
			}
			n.handleNodeDeletion(node)
		},
	})

	coreResourceInformerLister.PodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			n.handlePodAddition(pod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}
			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}
			n.handlePodUpdate(oldPod, newPod)
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			n.handlePodDeletion(pod)
		},
	})
}

func (n *NodeLevelMetricsCollector) handleNodeAddition(node *corev1.Node) {
	if _, ok := n.nodeResouece[node.Name]; !ok {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		n.nodeResouece[node.Name] = nodeResourceInfo{
			allocatableResource: corev1.ResourceList{},
			requestedResource:   corev1.ResourceList{},
		}

		for resourceName, resourceValue := range node.Status.Allocatable {
			n.nodeResouece[node.Name].allocatableResource[resourceName] = resourceValue
		}
	}
}

func (n *NodeLevelMetricsCollector) handleNodeDeletion(node *corev1.Node) {
	if _, ok := n.nodeResouece[node.Name]; !ok {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		delete(n.nodeResouece, node.Name)
	}
}

func (n *NodeLevelMetricsCollector) addPodResourceRequested(pod *corev1.Pod) {
	for _, container := range pod.Spec.Containers {
		for resourceName, resourceValue := range container.Resources.Requests {
			requested, ok := n.nodeResouece[pod.Spec.NodeName].requestedResource[resourceName]
			if !ok {
				requested = resource.Quantity{}
			}
			requested.Add(resourceValue)
			n.nodeResouece[pod.Spec.NodeName].requestedResource[resourceName] = requested
		}
	}
}

func (n *NodeLevelMetricsCollector) handlePodAddition(pod *corev1.Pod) {
	if pod.Spec.NodeName != "" {
		n.mutex.Lock()
		defer n.mutex.Unlock()

		if _, ok := n.nodeResouece[pod.Spec.NodeName]; !ok {
			klog.Warningf("Node %s not found in cluster", pod.Spec.NodeName)
		}
		n.addPodResourceRequested(pod)
	}
}

func (n *NodeLevelMetricsCollector) handlePodUpdate(oldPod, newPod *corev1.Pod) {
	if oldPod.Spec.NodeName == "" && newPod.Spec.NodeName != "" {
		n.mutex.Lock()
		defer n.mutex.Unlock()

		if _, ok := n.nodeResouece[newPod.Spec.NodeName]; !ok {
			klog.Warningf("Node %s not found in cluster", newPod.Spec.NodeName)
		}
		n.addPodResourceRequested(newPod)
	}
}

func (n *NodeLevelMetricsCollector) deletePodResourceRequested(pod *corev1.Pod) {
	for _, container := range pod.Spec.Containers {
		for resourceName, resourceValue := range container.Resources.Requests {
			requested, ok := n.nodeResouece[pod.Spec.NodeName].requestedResource[resourceName]
			if !ok {
				continue
			}
			requested.Sub(resourceValue)
			n.nodeResouece[pod.Spec.NodeName].requestedResource[resourceName] = requested
		}
	}
}

func (n *NodeLevelMetricsCollector) handlePodDeletion(pod *corev1.Pod) {
	if pod.Spec.NodeName != "" {
		n.mutex.Lock()
		defer n.mutex.Unlock()

		n.deletePodResourceRequested(pod)
	}
}

func (n *NodeLevelMetricsCollector) StartCollectNodeLevelMetrics(ctx context.Context,
	interval time.Duration, agentOptions *options.AgentOptions) {
	ticker := time.NewTicker(interval)

	klog.Infof("Start collecting Node level metrics")
	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			klog.Infof("Stop node level metrics collect")
			return
		case <-ticker.C:
			n.collectNodeCost(agentOptions)
			n.collectNodeResourceUsage(ctx, agentOptions)
			n.collectNodeResourceMetrics(agentOptions)
		}
	}
}

func (n *NodeLevelMetricsCollector) collectNodeCost(agentOptions *options.AgentOptions) {
	nodes, err := n.nodeLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all nodes error:%v", err)
		return
	}

	for _, node := range nodes {
		nodeCostInfo, err := n.provider.GetNodeHourlyPrice(node)
		if err != nil {
			klog.Errorf("Get node price from cloud provider error:%v", err)
			continue
		}
		metricsLabelValues := prometheus.Labels{
			values.NodeNameLabelKey:          node.Name,
			values.NodeInstanceTypeLabelKey:  nodeCostInfo.InstanceType,
			values.BillingModeLabelKey:       nodeCostInfo.BillingMode,
			values.NodeBillingPeriodLabelKey: strconv.Itoa(nodeCostInfo.BillingPeriod),
			values.RegionLabelKey:            nodeCostInfo.Region,
			values.CloudProviderLabelKey:     nodeCostInfo.CloudProvider,
			values.ClusterNameLabelKey:       agentOptions.ClusterName,
			values.ClusterIdLabelKey:         agentOptions.ClusterId,
		}
		n.nodeCPUCoreCostGV.With(metricsLabelValues).Set(nodeCostInfo.CPUCoreHourlyPrice)
		n.nodeRAMGBCostGV.With(metricsLabelValues).Set(nodeCostInfo.RAMGBHourlyPrice)
		n.nodeTotalCostGV.With(metricsLabelValues).Set(nodeCostInfo.NodeTotalHourlyPrice)

		metricsLabelValues[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
		n.nodeResourceHourlyTotalCostGV.With(metricsLabelValues).Set(nodeCostInfo.CPUCoreHourlyPrice * nodeCostInfo.CPUCore)
		metricsLabelValues[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
		n.nodeResourceHourlyTotalCostGV.With(metricsLabelValues).Set(nodeCostInfo.RAMGBHourlyPrice * nodeCostInfo.RamGiB)
	}
}

func (n *NodeLevelMetricsCollector) collectNodeResourceUsage(ctx context.Context, agentOptions *options.AgentOptions) {
	nodes, err := n.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List all node metrics error:%v", err)
		return
	}

	for _, node := range nodes.Items {
		nodeCostInfo, err := n.getNodeCostInfo(node.Name)
		if err != nil {
			continue
		}

		metricsLabels := prometheus.Labels{
			values.NodeNameLabelKey:    node.Name,
			values.ClusterNameLabelKey: agentOptions.ClusterName,
			values.ClusterIdLabelKey:   agentOptions.ClusterId,
			values.BillingModeLabelKey: nodeCostInfo.BillingMode,
		}
		cpu, memory := utils.ParseNodeResourceUsage(node)
		metricsLabels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
		n.nodeResourceUsageGV.With(metricsLabels).Set(cpu)
		metricsLabels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
		n.nodeResourceUsageGV.With(metricsLabels).Set(memory)
	}
}

func (n *NodeLevelMetricsCollector) collectNodeResourceMetrics(agentOptions *options.AgentOptions) {
	nodes, err := n.nodeLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all nodes error:%v", err)
		return
	}

	for _, node := range nodes {
		nodeCostInfo, err := n.getNodeCostInfo(node.Name)
		if err != nil {
			continue
		}

		metricsLabels := prometheus.Labels{
			values.NodeNameLabelKey:    node.Name,
			values.ClusterNameLabelKey: agentOptions.ClusterName,
			values.ClusterIdLabelKey:   agentOptions.ClusterId,
			values.BillingModeLabelKey: nodeCostInfo.BillingMode,
		}

		metricsLabels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
		n.nodeResourceTotalGV.With(metricsLabels).Set(nodeCostInfo.CPUCore)

		n.mutex.Lock()
		if _, ok := n.nodeResouece[node.Name]; ok {
			allocatable := n.nodeResouece[node.Name].allocatableResource[corev1.ResourceCPU]
			requested := n.nodeResouece[node.Name].requestedResource[corev1.ResourceCPU]

			resourceSystemTaken := nodeCostInfo.CPUCore - utils.ConvertQualityToCore(allocatable)
			allocatable.Sub(requested)
			resoruceAvailable := utils.ConvertQualityToCore(allocatable)

			n.nodeResourceSystemTakenGV.With(metricsLabels).Set(resourceSystemTaken)
			n.nodeResourceAvailableGV.With(metricsLabels).Set(resoruceAvailable)
		}
		n.mutex.Unlock()

		metricsLabels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
		n.nodeResourceTotalGV.With(metricsLabels).Set(nodeCostInfo.RamGiB)

		n.mutex.Lock()
		if _, ok := n.nodeResouece[node.Name]; ok {
			allocatable := n.nodeResouece[node.Name].allocatableResource[corev1.ResourceMemory]
			requested := n.nodeResouece[node.Name].requestedResource[corev1.ResourceMemory]

			resourceSystemTaken := nodeCostInfo.RamGiB - utils.ConvertQualityToGiB(allocatable)
			allocatable.Sub(requested)
			resoruceAvailable := utils.ConvertQualityToGiB(allocatable)

			n.nodeResourceSystemTakenGV.With(metricsLabels).Set(resourceSystemTaken)
			n.nodeResourceAvailableGV.With(metricsLabels).Set(resoruceAvailable)
		}
		n.mutex.Unlock()
	}
}

func (n *NodeLevelMetricsCollector) getNodeCostInfo(nodeName string) (*api.InstancePriceInfo, error) {
	node, err := n.nodeLister.Get(nodeName)
	if err != nil {
		klog.Errorf("Get node from lister error:%v", err)
		return nil, err
	}

	nodeCostInfo, err := n.provider.GetNodeHourlyPrice(node)
	if err != nil {
		klog.Errorf("Get node price from cloud provider error:%v", err)
		return nil, err
	}
	return nodeCostInfo, nil
}
