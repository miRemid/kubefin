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

package api

import (
	"github.com/prometheus/common/model"
)

const (
	KubeFinAPIVersion = "v1"
	KubeFinStatusKind = "Status"
	KubeFinListKind   = "List"

	CloudProviderAck     = "ack"
	CloudProviderDefault = "default"
)

var (
	QueryFailedStatus = "QueryFailed"
	QueryFailedReason = "Query backend with promql error"

	QueryNotFoundStatus = "QueryNotFound"
	QueryNotFoundReason = "Query parameter not found"

	QueryParaErrorStatus = "BadParameters"
	QueryParaErrorReason = "Query parameters are wrong"
)

const (
	AggregateByAll         = "all"
	AggregateByPod         = "pod"
	AggregateByDeployment  = "deployment"
	AggregateByStatefulSet = "statefulset"
	AggregateByDaemonSet   = "daemonset"
	AggregateByLabel       = "label"

	QueryStartTimePara   = "startTime"
	QueryEndTimePara     = "endTime"
	QueryStepSecondsPara = "stepSeconds"
	QueryAggregateBy     = "aggregateBy"
)

type StatusError struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Reason     string `json:"reason"`
	Code       int    `json:"code"`
}

type StandardInfoList struct {
	Items []interface{} `json:"items"`
}

type ClusterMetricsSummaryList struct {
	Items []ClusterMetricsSummary `json:"items"`
}

type ClusterCostsSummaryList struct {
	Items []*ClusterCostsSummary `json:"items"`
}

type ClusterCostsSummary struct {
	ClusterBasicProperty
	// ClusterMonthCostCurrent means the total cost in current month
	ClusterMonthCostCurrent float64 `json:"clusterMonthCostCurrent,omitempty"`
	// ClusterMonthEstimateCost means the estimating cost with previous 7 days costs
	ClusterMonthEstimateCost float64 `json:"clusterMonthEstimateCost,omitempty"`
	// ClusterAvgDailyCost means average daily costs in current month
	ClusterAvgDailyCost float64 `json:"clusterAvgDailyCost,omitempty"`
	// ClusterAvgDailyCost means average core/hour costs in current month
	ClusterAvgHourlyCoreCost float64 `json:"ClusterAvgHourlyCoreCost,omitempty"`
}

type ClusterResourceCostList struct {
	ClusterId string                 `json:"clusterId"`
	Items     []*ClusterResourceCost `json:"items"`
}

type ClusterResourceCost struct {
	// Timestamp is in unix timestamp format, you can transform it to any you want
	Timestamp               int64   `json:"timestamp"`
	TotalCost               float64 `json:"totalCost,omitempty"`
	CostOnDemandBillingMode float64 `json:"costOnDemandBillingMode,omitempty"`
	CostSpotBillingMode     float64 `json:"costSpotBillingMode,omitempty"`
	CostPeriodBillingMode   float64 `json:"costPeriodBillingMode,omitempty"`
	// TODO: implement this type
	CostFallbackBillingMode float64 `json:"costFallbackBillingMode,omitempty"`

	// CPUCoreCount means the average core hour count in this period
	CPUCoreCount float64 `json:"cpuCoreCount,omitempty"`
	// CPUCoreUsage means the average core hour usage in this period
	CPUCoreUsage float64 `json:"cpuCoreUsage,omitempty"`
	CPUCost      float64 `json:"cpuCost,omitempty"`

	// RAMGBCount means the average ram hour count in this period
	RAMGBCount float64 `json:"ramGBCount,omitempty"`
	// RAMGBUsage means the average ram hour usage in this period
	RAMGBUsage float64 `json:"ramGBUsage,omitempty"`
	RAMCost    float64 `json:"ramCost,omitempty"`
}

type ClusterWorkloadCostList struct {
	ClusterId string                 `json:"clusterId"`
	Items     []*ClusterWorkloadCost `json:"items"`
}

type ClusterWorkloadCost struct {
	Namespace    string `json:"namespace"`
	WorkloadName string `json:"workloadName"`
	// WorkloadType could be pod/daemonset/statefulset/deployment
	WorkloadType string                       `json:"workloadType"`
	CostList     []*ClusterWorkloadCostDetail `json:"costList"`
}

type ClusterWorkloadCostDetail struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	// PodCount means the average pod count in this period
	PodCount       float64 `json:"podCount,omitempty"`
	CPUCoreRequest float64 `json:"cpuCoreRequest,omitempty"`
	CPUCoreUsage   float64 `json:"cpuCoreUsage,omitempty"`
	RAMGBRequest   float64 `json:"ramGBRequest,omitempty"`
	RAMGBUsage     float64 `json:"ramGBUsage,omitempty"`
	TotalCost      float64 `json:"totalCost,omitempty"`
}

type ClusterNamespaceCostList struct {
	ClusterId string                  `json:"clusterId,omitempty"`
	Items     []*ClusterNamespaceCost `json:"items,omitempty"`
}

type ClusterNamespaceCost struct {
	Namespace string                        `json:"namespace,omitempty"`
	CostList  []*ClusterNamespaceCostDetail `json:"costList,omitempty"`
}

type ClusterNamespaceCostDetail struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	// PodCount means the average pod count in this period
	PodCount       float64 `json:"podCount,omitempty"`
	CPUCoreRequest float64 `json:"cpuRequest,omitempty"`
	CPUCoreUsage   float64 `json:"cpuCoreUsage,omitempty"`
	RAMGBRequest   float64 `json:"ramGBRequest,omitempty"`
	RAMGBUsage     float64 `json:"ramGBUsage,omitempty"`
	TotalCost      float64 `json:"totalCost,omitempty"`
}

type ClusterMetricsSummary struct {
	ClusterBasicProperty
	NodeNumbersCurrent                int64 `json:"nodeNumbersCurrent"`
	OnDemandBillingNodeNumbersCurrent int64 `json:"onDemandBillingNodeNumbersCurrent,omitempty"`
	PeriodBillingNodeNumbersCurrent   int64 `json:"periodBillingNodeNumbersCurrent,omitempty"`
	SpotBillingNodeNumbersCurrent     int64 `json:"spotBillingNodeNumbersCurrent,omitempty"`
	FallbackBillingNodeNumbersCurrent int64 `json:"fallbackBillingNodeNumbersCurrent,omitempty"`
	PodTotalCurrent                   int64 `json:"podTotalCurrent,omitempty"`
	PodScheduledCurrent               int64 `json:"podScheduledCurrent,omitempty"`
	PodUnscheduledCurrent             int64 `json:"podUnscheduledCurrent,omitempty"`
	// CPUCoreTotal means all nodes' cpu core
	CPUCoreTotal float64 `json:"cpuCoreTotal,omitempty"`
	// CPUCoreRequest means all pods' cpu core request
	CPUCoreRequest float64 `json:"cpuCoreRequest,omitempty"`
	// CPUCoreCapacity means all nodes' cpu core capacity
	CPUCoreCapacity float64 `json:"cpuCoreCapacity,omitempty"`
	// CPUCoreUsage means all pods' cpu core usage
	CPUCoreUsage float64 `json:"cpuCoreUsage,omitempty"`
	// RAMGBTotal means all nodes' ram gb
	RAMGBTotal float64 `json:"ramGBTotal,omitempty"`
	// RAMGBRequest all pods' ram gb request
	RAMGBRequest float64 `json:"ramGBRequest,omitempty"`
	// RAMGBCapacity means all nodes' ram gb capacity
	RAMGBCapacity float64 `json:"ramGBCapacity,omitempty"`
	// RAMGBUsage means all pods' ram gb usage
	RAMGBUsage float64 `json:"ramGBUsage,omitempty"`
}

type ClusterResourceMetrics struct {
	ClusterId                 string             `json:"clusterId"`
	ResourceType              string             `json:"resourceType"`
	Unit                      string             `json:"unit"`
	ResourceRequestValues     []model.SamplePair `json:"resourceRequestValues"`
	ResourceAllocatableValues []model.SamplePair `json:"resourceAllocatableValues"`
	ResourceUsageValues       []model.SamplePair `json:"resourceUsageValues"`
	ResourceCapacityValues    []model.SamplePair `json:"resourceCapacityValues"`
	ResourceTotalValues       []model.SamplePair `json:"resourceTotalValues"`
}

type ClusterBasicProperty struct {
	ClusterName    string `json:"clusterName"`
	ClusterId      string `json:"clusterId"`
	CloudProvider  string `json:"cloudProvider"`
	ClusterRegion  string `json:"clusterRegion"`
	LastActiveTime int64  `json:"lastActiveTime"`
	// ClusterConnectionSate can be running/connect_failed
	ClusterConnectionSate string `json:"clusterConnectionSate"`
	// ClusterActiveTime shows the cluster active time in seconds
	ClusterActiveTime float64 `json:"clusterActiveTime"`
	// ConnectionTime shows the time the cluster connected
	ConnectionTime int64 `json:"connectionTime,omitempty"`
}
