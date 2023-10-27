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

package query

import "github.com/kubefin/kubefin/pkg/values"

var (
	QlSumNodesResourceTotalFromCluster       = "sum(" + values.NodeResourceTotalMetricsName + "{cluster_id='%s',resource='%s'})"
	QlSumNodesResourceAvailableFromCluster   = "sum(" + values.NodeResourceAvailableMetricsName + "{cluster_id='%s',resource='%s'})"
	QlSumNodesResourceUsageFromCluster       = "sum(" + values.NodeResourceUsageMetricsName + "{cluster_id='%s',resource='%s'})"
	QlSumNodesResourceSystemTakenFromCluster = "sum(" + values.NodeResourceSystemTakenName + "{cluster_id='%s',resource='%s'})"
	QlSumPodResourceRequestFromCluster       = "sum(" + values.PodResourceRequestMetricsName + "{cluster_id='%s',resource='%s'})"

	QlNodesNumber         = "count(" + values.NodeTotalHourlyCostMetricsName + ") by (cluster_id,billing_mode)"
	QlPodsNumber          = "count(" + values.PodResoueceCostMetricsName + ") by (cluster_id,scheduled)"
	QlResourceTotal       = "sum(" + values.NodeResourceTotalMetricsName + ") by (cluster_id,resource)"
	QlResourceUsage       = "sum(" + values.NodeResourceUsageMetricsName + ") by (cluster_id,resource)"
	QlResourceRequest     = "sum(" + values.PodResourceRequestMetricsName + ") by (cluster_id,resource)"
	QlResourceAvailable   = "sum(" + values.NodeResourceAvailableMetricsName + ") by (cluster_id,resource)"
	QlResoruceSystemTaken = "sum(" + values.NodeResourceSystemTakenName + ") by (cluster_id,resource)"

	QlNodesNumberFromCluster         = "count(" + values.NodeTotalHourlyCostMetricsName + "{cluster_id='%s'}) by (billing_mode)"
	QlPodsNumberFromCluster          = "count(" + values.PodResoueceCostMetricsName + "{cluster_id='%s'}) by (scheduled)"
	QlResourceTotalFromCluster       = "sum(" + values.NodeResourceTotalMetricsName + "{cluster_id='%s'}) by (resource)"
	QlResourceUsageFromCluster       = "sum(" + values.NodeResourceUsageMetricsName + "{cluster_id='%s'}) by (resource)"
	QlResourceRequestFromCluster     = "sum(" + values.PodResourceRequestMetricsName + "{cluster_id='%s'}) by (resource)"
	QlResourceAvailableFromCluster   = "sum(" + values.NodeResourceAvailableMetricsName + "{cluster_id='%s'}) by (resource)"
	QlResourceSystemTakenFromCluster = "sum(" + values.NodeResourceSystemTakenName + "{cluster_id='%s'}) by (resource)"

	QlTotalPodsNumberFromCluster            = "count(count(" + values.PodResourceRequestMetricsName + "{cluster_id='%s'}) by (pod))"
	QlPodsNumberByScheduleStatusFromCluster = "count(count(" + values.PodResourceRequestMetricsName + "{cluster_id='%s',scheduled='%s'}) by (pod))"

	QlClusterActiveTimeWithTimeRange = "count_over_time(" + values.ClusterActiveMetricsName + "{cluster_id='%s'}[%ds])*15"

	QlNodesTotalHourlyCostFromClusterWithTimeRange            = "sum(sum_over_time(" + values.NodeTotalHourlyCostMetricsName + "{cluster_id='%s'}[%ds]))/240"
	QlNodesTotalHourlyBillingModeCostFromClusterWithTimeRange = "sum(sum_over_time(" + values.NodeTotalHourlyCostMetricsName + "{cluster_id='%s'}[%ds])/240) by (billing_mode)"

	// TODO: NodeCPUHourlyCostMetricsName/NodeRAMHourlyCostMetricsName could be merged as one
	QlNodeCPUTotalCostFromClusterWithTimeRange       = "sum(sum_over_time(" + values.NodeResourceHourlyCostMetricsName + "{cluster_id='%s',resource='cpu'}[%ds]))/240"
	QlNodeResourceTotalCostFromClusterWithTimeRange  = "sum(sum_over_time(" + values.NodeResourceHourlyCostMetricsName + "{cluster_id='%s'}[%ds])/240) by (resource)"
	QlNodeCPUTotalCostWithTimeRange                  = "sum(sum_over_time(" + values.NodeResourceHourlyCostMetricsName + "{resource='cpu'}[%ds])/240) by (cluster_id)"
	QlNodeResourceTotalCountFromClusterWithTimeRange = "sum(sum_over_time(" + values.NodeResourceTotalMetricsName + "{cluster_id='%s',resource='%s'}[%ds]))/240"
	QlNodeCPUTotalCountWithTimeRange                 = "sum(sum_over_time(" + values.NodeResourceTotalMetricsName + "{resource='%s'}[%ds])/240) by (cluster_id)"
	QlNodeResourceUsageCountFromClusterWithTimeRange = "sum(sum_over_time(" + values.NodeResourceUsageMetricsName + "{cluster_id='%s',resource='%s'}[%ds]))/240"

	QlPodTotalCostFromClusterWithTimeRange            = "sum(sum_over_time(" + values.PodResoueceCostMetricsName + "{cluster_id='%s'}[%ds])/240) by (pod,namespace)"
	QlPodResourceRequestFromClusterWithTimeRange      = "sum(sum_over_time(" + values.PodResourceRequestMetricsName + "{cluster_id='%s'}[%ds])/240) by (pod,namespace,resource)"
	QlPodResourceUsageFromClusterWithTimeRange        = "sum(sum_over_time(" + values.PodResourceUsageMetricsName + "{cluster_id='%s'}[%ds])/240) by (pod,namespace,resource)"
	QlWorkloadTotalCostFromClusterWithTimeRange       = "sum(sum_over_time(" + values.WorkloadResourceCostMetricsName + "{cluster_id='%s',workload_type=~'%s'}[%ds])/240) by (namespace,workload_name,workload_type)"
	QlWorkloadPodFromClusterWithTimeRange             = "sum(sum_over_time(" + values.WorkloadPodCountMetricsName + "{cluster_id='%s',workload_type=~'%s'}[%ds])) by (namespace,workload_name,workload_type)"
	QlWorkloadResourceRequestFromClusterWithTimeRange = "sum(sum_over_time(" + values.WorkloadResourceRequestMetricsName + "{cluster_id='%s',workload_type=~'%s'}[%ds])/240) by (namespace,workload_name,workload_type,resource)"
	QlWorkloadResourceUsageFromClusterWithTimeRange   = "sum(sum_over_time(" + values.WorkloadResourceUsageMetricsName + "{cluster_id='%s',workload_type=~'%s'}[%ds])/240) by (namespace,workload_name,workload_type,resource)"
	QlNSTotalCostFromClusterWithTimeRange             = "sum(sum_over_time(" + values.PodResoueceCostMetricsName + "{cluster_id='%s'}[%ds])/240) by (namespace)"
	// TODO: Check the scheduled label has effect on this
	QlNSPodFromClusterWithTimeRange             = "sum(count_over_time(" + values.PodResoueceCostMetricsName + "{cluster_id='%s'}[%ds])) by (namespace)"
	QlNSResourceRequestFromClusterWithTimeRange = "sum(sum_over_time(" + values.PodResourceRequestMetricsName + "{cluster_id='%s'}[%ds])/240) by (namespace,resource)"
	QlNSResourceUsageFromClusterWithTimeRange   = "sum(sum_over_time(" + values.PodResourceUsageMetricsName + "{cluster_id='%s'}[%ds])/240) by (namespace,resource)"

	// QlNodesTotalCostsFromClusterWithTimeRange get all nodes cost with time range, we sample metrics
	// every 15 seconds, so 240 is used to transform it to one hour
	QlNodesTotalCostsFromClusterWithTimeRange = "sum(sum_over_time(" + values.NodeTotalHourlyCostMetricsName + "{cluster_id='%s'}[%ds]))/240"
	QlNodesTotalCostsWithTimeRange            = "sum(sum_over_time(" + values.NodeTotalHourlyCostMetricsName + "[%ds])/240) by (cluster_id)"

	QlAllClustersActivity   = "kubefin_cluster_active"
	QlClusterActivity       = "kubefin_cluster_active{cluster_id='%s'}"
	QlClusterActiveTime     = "count_over_time(" + values.ClusterActiveMetricsName + "{cluster_id='%s'}[%ds])*15"
	QlAllClustersActiveTime = "count_over_time(" + values.ClusterActiveMetricsName + "[%ds])*15"
)
