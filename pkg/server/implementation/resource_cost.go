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

package implementation

import (
	"fmt"
	"sort"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/query"
	"github.com/kubefin/kubefin/pkg/values"
	"github.com/prometheus/common/model"
)

// QueryClusterResourceCost queries cluster resource cots
func QueryClusterResourceCost(tenantId, clusterId string,
	start, end, stepSeconds int64) (*api.ClusterResourceCostList, error) {
	var totalCosts map[int64]float64
	var billingModeCosts map[string]map[int64]float64
	var resourceTotalCost map[string]map[int64]float64
	var cpuTotalHourCount map[int64]float64
	var cpuUsageHourCount map[int64]float64
	var ramTotalHourCount map[int64]float64
	var ramUsageHourCount map[int64]float64

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(7)
	go func() {
		defer wg.Done()
		totalCosts, err = queryNodeTotalCost(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		billingModeCosts, err = queryNodeBillingModeCost(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceTotalCost, err = queryNodeResourceTotalCost(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuTotalHourCount, err = queryNodeCPUTotalHour(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuUsageHourCount, err = queryNodeCPUUsageHour(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		ramTotalHourCount, err = queryNodeRAMTotalHour(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		ramUsageHourCount, err = queryNodeRAMUsageHour(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	clusterResourceCost := make(map[int64]*api.ClusterResourceCost)
	parseResourceTotalCost(clusterResourceCost, totalCosts)
	parseResourceBillingModeCost(clusterResourceCost, billingModeCosts)
	parseNodeResourceTotalCost(clusterResourceCost, resourceTotalCost)
	parseResourceCPUTotalHour(clusterResourceCost, cpuTotalHourCount, stepSeconds)
	parseResourceRAMTotalHour(clusterResourceCost, ramTotalHourCount, stepSeconds)
	parseResourceCPUUsageHour(clusterResourceCost, cpuUsageHourCount, stepSeconds)
	parseResourceRAMUsageHour(clusterResourceCost, ramUsageHourCount, stepSeconds)

	return convertClusterResourceCostToList(clusterId, clusterResourceCost), nil
}

func parseResourceTotalCost(clusterResourceCost map[int64]*api.ClusterResourceCost, totalCosts map[int64]float64) {
	for timeStamp, v := range totalCosts {
		clusterResourceCost[timeStamp] = &api.ClusterResourceCost{
			Timestamp: timeStamp,
			TotalCost: v,
		}
	}
}

func parseResourceBillingModeCost(clusterResourceCost map[int64]*api.ClusterResourceCost, billingModeCosts map[string]map[int64]float64) {
	for billingMode, cost := range billingModeCosts {
		for timeStamp, v := range cost {
			cost, ok := clusterResourceCost[timeStamp]
			if !ok {
				cost = &api.ClusterResourceCost{
					Timestamp: timeStamp,
				}
			}
			switch billingMode {
			case values.BillingModeOnDemand:
				cost.CostOnDemandBillingMode = v
			case values.BillingModeSpot:
				cost.CostSpotBillingMode = v
			case values.BillingModeFallback:
				cost.CostFallbackBillingMode = v
			default:
				klog.Errorf("Billing mode %s not supported yet", billingMode)
			}
			clusterResourceCost[timeStamp] = cost
		}
	}
}

func parseNodeResourceTotalCost(clusterResourceCost map[int64]*api.ClusterResourceCost, resourceTotalCost map[string]map[int64]float64) {
	for resourceType, v := range resourceTotalCost {
		switch resourceType {
		case "cpu":
			for timeStamp, costVale := range v {
				cost, ok := clusterResourceCost[timeStamp]
				if !ok {
					cost = &api.ClusterResourceCost{
						Timestamp: timeStamp,
					}
				}
				cost.CPUCost = costVale
				clusterResourceCost[timeStamp] = cost
			}
		case "memory":
			for timeStamp, costVale := range v {
				cost, ok := clusterResourceCost[timeStamp]
				if !ok {
					cost = &api.ClusterResourceCost{
						Timestamp: timeStamp,
					}
				}
				cost.RAMCost = costVale
				clusterResourceCost[timeStamp] = cost
			}
		}
	}
}

func parseResourceCPUTotalHour(clusterResourceCost map[int64]*api.ClusterResourceCost, cpuTotalHourCount map[int64]float64, stepSeconds int64) {
	for timeStamp, v := range cpuTotalHourCount {
		cost, ok := clusterResourceCost[timeStamp]
		if !ok {
			cost = &api.ClusterResourceCost{
				Timestamp: timeStamp,
			}
		}
		cost.CPUCoreCount = v / float64(stepSeconds) * values.HourInSeconds
		clusterResourceCost[timeStamp] = cost
	}
}

func parseResourceRAMTotalHour(clusterResourceCost map[int64]*api.ClusterResourceCost, ramTotalHourCount map[int64]float64, stepSeconds int64) {
	for timeStamp, v := range ramTotalHourCount {
		cost, ok := clusterResourceCost[timeStamp]
		if !ok {
			cost = &api.ClusterResourceCost{
				Timestamp: timeStamp,
			}
		}
		cost.RAMGBCount = v / float64(stepSeconds) * values.HourInSeconds
		clusterResourceCost[timeStamp] = cost
	}
}

func parseResourceCPUUsageHour(clusterResourceCost map[int64]*api.ClusterResourceCost, cpuUsageHourCount map[int64]float64, stepSeconds int64) {
	for timeStamp, v := range cpuUsageHourCount {
		cost, ok := clusterResourceCost[timeStamp]
		if !ok {
			cost = &api.ClusterResourceCost{
				Timestamp: timeStamp,
			}
		}
		cost.CPUCoreUsage = v / float64(stepSeconds) * values.HourInSeconds
		clusterResourceCost[timeStamp] = cost
	}
}

func parseResourceRAMUsageHour(clusterResourceCost map[int64]*api.ClusterResourceCost, ramUsageHourCount map[int64]float64, stepSeconds int64) {
	for timeStamp, v := range ramUsageHourCount {
		cost, ok := clusterResourceCost[timeStamp]
		if !ok {
			cost = &api.ClusterResourceCost{
				Timestamp: timeStamp,
			}
		}
		cost.RAMGBUsage = v / float64(stepSeconds) * values.HourInSeconds
		clusterResourceCost[timeStamp] = cost
	}
}

func queryNodeTotalCost(tenantId, clusterId string, start, end, stepSeconds int64) (map[int64]float64, error) {
	totalCosts := make(map[int64]float64)
	promql := fmt.Sprintf(query.QlNodesTotalHourlyCostFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) total node cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret[0].Values {
		totalCosts[v.Timestamp.Unix()] = float64(v.Value)
	}
	return totalCosts, nil
}

func queryNodeBillingModeCost(tenantId, clusterId string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	billingModeCosts := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlNodesTotalHourlyBillingModeCostFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) billing mode cost error:%v", clusterId, err)
		return nil, err
	}
	for _, v := range ret {
		billingMode := string(v.Metric[model.LabelName(values.BillingModeLabelKey)])
		billingModeCosts[billingMode] = make(map[int64]float64)
		for _, vv := range v.Values {
			billingModeCosts[billingMode][vv.Timestamp.Unix()] = float64(vv.Value)
		}
	}
	return billingModeCosts, nil
}

func queryNodeResourceTotalCost(tenantId, clusterId string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	// maps [cpu/memory][timestamp]cost
	resourceTotalCost := make(map[string]map[int64]float64)
	promal := fmt.Sprintf(query.QlNodeResourceTotalCostFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promal, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) cpu total cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret {
		resourceType := string(v.Metric[model.LabelName(values.ResourceTypeLabelKey)])
		switch string(resourceType) {
		case "cpu":
			resourceTotalCost["cpu"] = map[int64]float64{}
			for _, vv := range v.Values {
				resourceTotalCost["cpu"][vv.Timestamp.Unix()] = float64(vv.Value)
			}
		case "memory":
			resourceTotalCost["memory"] = map[int64]float64{}
			for _, vv := range v.Values {
				resourceTotalCost["memory"][vv.Timestamp.Unix()] = float64(vv.Value)
			}
		}
	}
	return resourceTotalCost, nil
}

func queryNodeCPUTotalHour(tenantId, clusterId string, start, end, stepSeconds int64) (map[int64]float64, error) {
	cpuTotalHourCount := make(map[int64]float64)
	promql := fmt.Sprintf(query.QlNodeResourceTotalCountFromClusterWithTimeRange, clusterId, corev1.ResourceCPU, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) cpu core hour cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret[0].Values {
		cpuTotalHourCount[v.Timestamp.Unix()] = float64(v.Value)
	}
	return cpuTotalHourCount, nil
}

func queryNodeCPUUsageHour(tenantId, clusterId string, start, end, stepSeconds int64) (map[int64]float64, error) {
	cpuUsageHourCount := make(map[int64]float64)
	promql := fmt.Sprintf(query.QlNodeResourceUsageCountFromClusterWithTimeRange, clusterId, corev1.ResourceCPU, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) cpu usage hour cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret[0].Values {
		cpuUsageHourCount[v.Timestamp.Unix()] = float64(v.Value)
	}
	return cpuUsageHourCount, nil
}

func queryNodeRAMTotalHour(tenantId, clusterId string, start, end, stepSeconds int64) (map[int64]float64, error) {
	ramTotalHourCount := make(map[int64]float64)
	promql := fmt.Sprintf(query.QlNodeResourceTotalCountFromClusterWithTimeRange, clusterId, corev1.ResourceMemory, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) ram GB hour cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret[0].Values {
		ramTotalHourCount[v.Timestamp.Unix()] = float64(v.Value)
	}
	return ramTotalHourCount, nil
}

func queryNodeRAMUsageHour(tenantId, clusterId string, start, end, stepSeconds int64) (map[int64]float64, error) {
	ramUsageHourCount := make(map[int64]float64)
	promql := fmt.Sprintf(query.QlNodeResourceUsageCountFromClusterWithTimeRange, clusterId, corev1.ResourceMemory, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) ram usage hour cost error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	for _, v := range ret[0].Values {
		ramUsageHourCount[v.Timestamp.Unix()] = float64(v.Value)
	}
	return ramUsageHourCount, nil
}

func convertClusterResourceCostToList(clusterId string, nodeCost map[int64]*api.ClusterResourceCost) *api.ClusterResourceCostList {
	ret := &api.ClusterResourceCostList{
		ClusterId: clusterId,
		Items:     []*api.ClusterResourceCost{},
	}
	for _, v := range nodeCost {
		ret.Items = append(ret.Items, v)
	}
	sort.Slice(ret.Items, func(i, j int) bool {
		return ret.Items[i].Timestamp < ret.Items[j].Timestamp
	})

	return ret
}
