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

package values

var (
	// cluster level metrics name
	ClusterActiveMetricsName = "kubefin_cluster_active"

	// Node level metrics name
	NodeCPUCoreHourlyCostMetricsName   = "kubefin_node_cpu_core_hourly_cost"
	NodeCPUHourlyCostMetricsName       = "kubefin_node_cpu_hourly_cost"
	NodeRAMGBHourlyCostMetricsName     = "kubefin_node_ram_gb_hourly_cost"
	NodeResourceHourlyCostMetricsName  = "kubefin_node_resource_hourly_cost"
	NodeRAMHourlyCostMetricsName       = "kubefin_node_ram_hourly_cost"
	NodeTotalHourlyCostMetricsName     = "kubefin_node_total_hourly_cost"
	NodeResourceUsageMetricsName       = "kubefin_node_resource_usage"
	NodeResourceAllocatableMetricsName = "kubefin_node_resource_allocatable"
	NodeResourceCapacityMetricsName    = "kubefin_node_resource_capacity"
	// NodeResourceTotalMetricsName means total resource of the node, contains the resource used by os&kubelet
	NodeResourceTotalMetricsName = "kubefin_node_resource_total"

	// workload level metric name
	WorkloadResourceCostMetricsName    = "kubefin_workload_resource_cost"
	WorkloadResourceRequestMetricsName = "kubefin_workload_resource_request"
	WorkloadResourceUsageMetricsName   = "kubefin_workload_resource_usage"
	WorkloadPodCountMetricsName        = "kubefin_workload_pod_count"

	// Pod level metrics name
	PodResourceRequestMetricsName = "kubefin_pod_resource_request"
	PodResourceUsageMetricsName   = "kubefin_pod_resource_usage"
	PodResoueceCostMetricsName    = "kubefin_pod_resource_cost"

	// metrics labels
	ClusterNameLabelKey       = "cluster_name"
	ClusterIdLabelKey         = "cluster_id"
	NamespaceLabelKey         = "namespace"
	LabelsLabelKey            = "labels"
	ResourceTypeLabelKey      = "resource"
	BillingModeLabelKey       = "billing_mode"
	NodeNameLabelKey          = "node"
	ContainerNameLabelKey     = "container"
	WorkloadTypeLabelKey      = "workload_type"
	WorkloadNameLabelKey      = "workload_name"
	NodeInstanceTypeLabelKey  = "instance_type"
	NodeBillingPeriodLabelKey = "billing_period"
	RegionLabelKey            = "region"
	CloudProviderLabelKey     = "cloud_provider"
	PodNameLabelKey           = "pod"
	PodScheduledKey           = "scheduled"
)
