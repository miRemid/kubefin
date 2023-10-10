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

	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/query"
	"github.com/kubefin/kubefin/pkg/values"
	"github.com/prometheus/common/model"
)

func QueryNamespaceCostsWithTimeRange(tenantId, clusterId string,
	start, end, stepSeconds int64) (*api.ClusterNamespaceCostList, error) {
	var totalCosts map[string]map[int64]float64
	var podCount map[string]map[int64]float64
	var cpuRequest map[string]map[int64]float64
	var ramRequest map[string]map[int64]float64
	var cpuUsage map[string]map[int64]float64
	var ramUsage map[string]map[int64]float64

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(4)
	go func() {
		defer wg.Done()
		totalCosts, err = queryNamespaceTotalCost(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		podCount, err = queryNamespacePodCount(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuRequest, ramRequest, err = queryNamespaceResourceRequest(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuUsage, ramUsage, err = queryNamespaceResourceUsage(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	nsCost := make(map[string]map[int64]*api.ClusterNamespaceCostDetail)
	parseNamepsaceTotalCosts(nsCost, totalCosts)
	parseNamespacePodCount(nsCost, podCount, stepSeconds)
	parseNamespaceResourceRequest(nsCost, cpuRequest, ramRequest, stepSeconds)
	parseNamespaceResourceUsage(nsCost, cpuUsage, ramUsage, stepSeconds)

	nsCosts := convertClusterNSCostToList(nsCost)
	ret := &api.ClusterNamespaceCostList{ClusterId: clusterId, Items: []*api.ClusterNamespaceCost{}}
	ret.Items = append(ret.Items, nsCosts...)
	return ret, nil
}

func parseNamepsaceTotalCosts(nsCost map[string]map[int64]*api.ClusterNamespaceCostDetail, totalCosts map[string]map[int64]float64) {
	for ns, details := range totalCosts {
		nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		for timeStamp, v := range details {
			nsCost[ns][timeStamp] = &api.ClusterNamespaceCostDetail{
				Timestamp: timeStamp,
				TotalCost: v,
			}
		}
	}
}

func parseNamespacePodCount(nsCost map[string]map[int64]*api.ClusterNamespaceCostDetail,
	podCount map[string]map[int64]float64, stepSeconds int64) {
	for ns, details := range podCount {
		item, ok := nsCost[ns]
		if !ok {
			nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterNamespaceCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].PodCount = v / (float64(stepSeconds) / values.MetricsPeriodInSeconds)
		}
	}
}

func parseNamespaceResourceRequest(nsCost map[string]map[int64]*api.ClusterNamespaceCostDetail,
	cpuRequest, ramRequest map[string]map[int64]float64, stepSeconds int64) {
	for ns, details := range cpuRequest {
		item, ok := nsCost[ns]
		if !ok {
			nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterNamespaceCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].CPUCoreRequest = v / float64(stepSeconds) * values.HourInSeconds
		}
	}

	for ns, details := range ramRequest {
		item, ok := nsCost[ns]
		if !ok {
			nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterNamespaceCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].RAMGBRequest = v / float64(stepSeconds) * values.HourInSeconds
		}
	}
}

func parseNamespaceResourceUsage(nsCost map[string]map[int64]*api.ClusterNamespaceCostDetail,
	cpuUsage, ramUsage map[string]map[int64]float64, stepSeconds int64) {
	for ns, details := range cpuUsage {
		item, ok := nsCost[ns]
		if !ok {
			nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterNamespaceCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].CPUCoreUsage = v / float64(stepSeconds) * values.HourInSeconds
		}
	}

	for ns, details := range ramUsage {
		item, ok := nsCost[ns]
		if !ok {
			nsCost[ns] = make(map[int64]*api.ClusterNamespaceCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; ok {
				item[timeStamp].RAMGBUsage = v
				continue
			}
			item[timeStamp] = &api.ClusterNamespaceCostDetail{
				Timestamp:  timeStamp,
				RAMGBUsage: v,
			}
		}
	}
}

func queryNamespaceTotalCost(tenantId, clusterId string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	totalCosts := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlNSTotalCostFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) namespace cost error:%v", clusterId, err)
		return nil, err
	}
	for _, ns := range ret {
		key := ns.Metric[model.LabelName(values.NamespaceLabelKey)]
		totalCosts[string(key)] = make(map[int64]float64)
		for _, v := range ns.Values {
			totalCosts[string(key)][v.Timestamp.Unix()] = float64(v.Value)
		}
	}

	return totalCosts, nil
}

func queryNamespacePodCount(tenantId, clusterId string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	podCount := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlNSPodFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) namespace pod count error:%v", clusterId, err)
		return nil, err
	}
	for _, ns := range ret {
		key := ns.Metric[model.LabelName(values.NamespaceLabelKey)]
		podCount[string(key)] = make(map[int64]float64)
		for _, v := range ns.Values {
			podCount[string(key)][v.Timestamp.Unix()] = float64(v.Value)
		}
	}

	return podCount, nil
}

func queryNamespaceResourceRequest(tenantId, clusterId string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuRequest := make(map[string]map[int64]float64)
	ramRequest := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlNSResourceRequestFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) namespace resource request error:%v", clusterId, err)
		return nil, nil, err
	}
	for _, ns := range ret {
		key := string(ns.Metric[model.LabelName(values.NamespaceLabelKey)])
		resourceType := ns.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuRequest[key] = make(map[int64]float64)
			for _, v := range ns.Values {
				cpuRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramRequest[key] = make(map[int64]float64)
			for _, v := range ns.Values {
				ramRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}

	return cpuRequest, ramRequest, nil
}

func queryNamespaceResourceUsage(tenantId, clusterId string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuUsage := make(map[string]map[int64]float64)
	ramUsage := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlNSResourceUsageFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) namespace resource usage error:%v", clusterId, err)
		return nil, nil, err
	}
	for _, ns := range ret {
		key := string(ns.Metric[model.LabelName(values.NamespaceLabelKey)])
		resourceType := ns.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuUsage[key] = make(map[int64]float64)
			for _, v := range ns.Values {
				cpuUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramUsage[key] = make(map[int64]float64)
			for _, v := range ns.Values {
				ramUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}

	return cpuUsage, ramUsage, nil
}

func convertClusterNSCostToList(nsCost map[string]map[int64]*api.ClusterNamespaceCostDetail) []*api.ClusterNamespaceCost {
	ret := []*api.ClusterNamespaceCost{}
	for nsKey, details := range nsCost {
		cost := &api.ClusterNamespaceCost{
			Namespace: nsKey,
			CostList:  []*api.ClusterNamespaceCostDetail{},
		}
		for _, v := range details {
			cost.CostList = append(cost.CostList, v)
		}
		sort.Slice(cost.CostList, func(i, j int) bool {
			return cost.CostList[i].Timestamp < cost.CostList[j].Timestamp
		})
		ret = append(ret, cost)
	}

	return ret
}
