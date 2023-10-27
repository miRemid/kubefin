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

package utils

import (
	"github.com/kubefin/kubefin/pkg/api"
)

func ValidateAllClustersMetricsSummary(clustersSummary *api.ClusterMetricsSummaryList) bool {
	if len(clustersSummary.Items) != 0 &&
		clustersSummary.Items[0].ClusterConnectionSate == "running" &&
		clustersSummary.Items[0].NodeNumbersCurrent != 0 &&
		clustersSummary.Items[0].PodTotalCurrent != 0 &&
		clustersSummary.Items[0].PodScheduledCurrent != 0 &&
		clustersSummary.Items[0].CPUCoreTotal != 0 &&
		clustersSummary.Items[0].RAMGiBTotal != 0 {
		return true
	}

	return false
}

func ValidateSpecificClusterMetricsSummary(clusterSummary *api.ClusterMetricsSummary) bool {
	if clusterSummary.ClusterConnectionSate == "running" &&
		clusterSummary.NodeNumbersCurrent != 0 &&
		clusterSummary.PodTotalCurrent != 0 &&
		clusterSummary.PodScheduledCurrent != 0 &&
		clusterSummary.CPUCoreTotal != 0 &&
		clusterSummary.RAMGiBTotal != 0 {
		return true
	}

	return false
}

func ValidateSpecificClusterResourceMetrics(metrics *api.ClusterResourceMetrics) bool {
	if metrics.ClusterId != "" &&
		len(metrics.ResourceTotalValues) != 0 &&
		len(metrics.ResourceUsageValues) != 0 &&
		len(metrics.ResourceRequestValues) != 0 &&
		len(metrics.ResourceSystemTakenValues) != 0 {
		return true
	}
	return false
}

func ValidateAllClustersCostsSummary(clustersCostsSummary *api.ClusterCostsSummaryList) bool {
	if len(clustersCostsSummary.Items) != 0 &&
		clustersCostsSummary.Items[0].ClusterConnectionSate == "running" &&
		clustersCostsSummary.Items[0].ClusterMonthCostCurrent != 0 &&
		clustersCostsSummary.Items[0].ClusterMonthEstimateCost != 0 &&
		clustersCostsSummary.Items[0].ClusterAvgDailyCost != 0 &&
		clustersCostsSummary.Items[0].ClusterAvgHourlyCoreCost != 0 {
		return true
	}

	return false
}

func ValidateSpecificClusterCostsSummary(clusterCostsSummary *api.ClusterCostsSummary) bool {
	if clusterCostsSummary.ClusterConnectionSate == "running" &&
		clusterCostsSummary.ClusterMonthCostCurrent != 0 &&
		clusterCostsSummary.ClusterMonthEstimateCost != 0 &&
		clusterCostsSummary.ClusterAvgDailyCost != 0 &&
		clusterCostsSummary.ClusterAvgHourlyCoreCost != 0 {
		return true
	}

	return false
}

func ValidateSpecificClusterComputeCosts(computeCosts *api.ClusterResourceCostList) bool {
	// if computeCosts.ClusterId != "" &&
	// 	len(computeCosts.Items) != 0 &&
	// 	computeCosts.Items[0].CostOnDemandBillingMode != 0 &&
	// 	computeCosts.Items[0].CPUCostOnDemandBillingMode != 0 &&
	// 	computeCosts.Items[0].CPUCoreCountOnDemandBillingMode != 0 &&
	// 	computeCosts.Items[0].RAMCostOnDemandBillingMode != 0 &&
	// 	computeCosts.Items[0].RAMGBCountOnDemandBillingMode != 0 {
	// 	return true
	// }

	return false
}

func ValidateSpecificClusterWorkloadCosts(workloadCosts *api.ClusterWorkloadCostList) bool {
	// if workloadCosts.ClusterId == "" ||
	// 	len(workloadCosts.Items) == 0 {
	// 	return false
	// }

	// for _, ele := range workloadCosts.Items {
	// 	// In kind cluster, these two workloads don't have resource requirements by default
	// 	if strings.Contains(ele.WorkloadName, "kube-proxy") ||
	// 		strings.Contains(ele.WorkloadName, "local-path") {
	// 		continue
	// 	}
	// 	if ele.CostList[0].CostOnDemandBillingMode == 0 {
	// 		return false
	// 	}
	// }

	return true
}
