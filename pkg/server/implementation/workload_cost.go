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

package implementation

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/prometheus/common/model"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/query"
	"github.com/kubefin/kubefin/pkg/values"
)

func generatePodNamespaceNameKey(labels model.Metric) string {
	namespace := labels[model.LabelName(values.NamespaceLabelKey)]
	name := labels[model.LabelName(values.PodNameLabelKey)]

	return fmt.Sprintf("pod/%s/%s", namespace, name)
}

func generateWorkloadNamespaceNameKey(labels model.Metric) string {
	workloadType := labels[model.LabelName(values.WorkloadTypeLabelKey)]
	namespace := labels[model.LabelName(values.NamespaceLabelKey)]
	name := labels[model.LabelName(values.WorkloadNameLabelKey)]

	return fmt.Sprintf("%s/%s/%s", workloadType, namespace, name)
}

func parseWorkloadNamespaceNameType(key string) (namespace, name, workloadType string) {
	parts := strings.Split(key, "/")
	workloadType = parts[0]
	namespace = parts[1]
	name = parts[2]

	return workloadType, namespace, name
}

func QueryWorkloadCostsWithTimeRange(tenantId, clusterId string,
	start, end, stepSeconds int64, aggregateBy string) (*api.ClusterWorkloadCostList, error) {
	var errs []error
	var err error
	var wg sync.WaitGroup

	var podCosts []*api.ClusterWorkloadCost
	if aggregateBy == api.AggregateByPod || aggregateBy == api.AggregateByAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			podCosts, err = queryPodCostsWithTimeRange(tenantId, clusterId, start, end, stepSeconds)
			if err != nil {
				errs = append(errs, err)
			}
		}()
	}

	var workloadCosts []*api.ClusterWorkloadCost
	if aggregateBy != api.AggregateByPod {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workloadCosts, err = queryHighLevelWorkloadCostsWithTimeRange(tenantId, clusterId, start, end, stepSeconds, aggregateBy)
			if err != nil {
				errs = append(errs, err)
			}
		}()
	}

	wg.Wait()
	if len(errs) > 0 {
		return nil, errors.NewAggregate(errs)
	}

	ret := &api.ClusterWorkloadCostList{ClusterId: clusterId, Items: []*api.ClusterWorkloadCost{}}
	ret.Items = append(ret.Items, podCosts...)
	ret.Items = append(ret.Items, workloadCosts...)
	return ret, nil
}

func queryPodCostsWithTimeRange(tenantId, clusterId string, start, end, stepSeconds int64) ([]*api.ClusterWorkloadCost, error) {
	var totalCosts map[string]map[int64]float64
	var cpuRequest map[string]map[int64]float64
	var ramRequest map[string]map[int64]float64
	var cpuUsage map[string]map[int64]float64
	var ramUsage map[string]map[int64]float64

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(3)
	go func() {
		defer wg.Done()
		totalCosts, err = queryPodTotalCostsWithTimeRange(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuRequest, ramRequest, err = queryPodResourceRequest(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuUsage, ramUsage, err = queryPodResourceUsage(tenantId, clusterId, start, end, stepSeconds)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}
	podWorkloadCost := make(map[string]map[int64]*api.ClusterWorkloadCostDetail)
	parsePodTotalCost(podWorkloadCost, totalCosts)
	parsePodResourceRequest(podWorkloadCost, cpuRequest, ramRequest, stepSeconds)
	parsePodResourceUsage(podWorkloadCost, cpuUsage, ramUsage, stepSeconds)

	return convertClusterWorkloadCostToList(podWorkloadCost), nil
}

func parsePodTotalCost(podWorkloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail, totalCosts map[string]map[int64]float64) {
	for pod, details := range totalCosts {
		podWorkloadCost[pod] = make(map[int64]*api.ClusterWorkloadCostDetail)
		for timeStamp, v := range details {
			podWorkloadCost[pod][timeStamp] = &api.ClusterWorkloadCostDetail{
				Timestamp: timeStamp,
				TotalCost: v,
				// PodCount is always 1
				PodCount: 1,
			}
		}
	}
}

func parsePodResourceRequest(podWorkloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	cpuRequest map[string]map[int64]float64, ramRequest map[string]map[int64]float64, stepSeconds int64) {
	for pod, details := range cpuRequest {
		item, ok := podWorkloadCost[pod]
		if !ok {
			podWorkloadCost[pod] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
					// PodCount is always 1
					PodCount: 1,
				}
				item[timeStamp].CPUCoreRequest = v / float64(stepSeconds) * values.HourInSeconds
			}
		}
	}

	for pod, details := range ramRequest {
		item, ok := podWorkloadCost[pod]
		if !ok {
			podWorkloadCost[pod] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
					// PodCount is always 1
					PodCount: 1,
				}
			}
			item[timeStamp].RAMGiBRequest = v / float64(stepSeconds) * values.HourInSeconds
		}
	}
}

func parsePodResourceUsage(podWorkloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	cpuUsage map[string]map[int64]float64, ramUsage map[string]map[int64]float64, stepSeconds int64) {
	for pod, details := range cpuUsage {
		item, ok := podWorkloadCost[pod]
		if !ok {
			podWorkloadCost[pod] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
					// PodCount is always 1
					PodCount: 1,
				}
			}
			item[timeStamp].CPUCoreUsage = v / float64(stepSeconds) * values.HourInSeconds
		}
	}

	for pod, details := range ramUsage {
		item, ok := podWorkloadCost[pod]
		if !ok {
			podWorkloadCost[pod] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
					// PodCount is always 1
					PodCount: 1,
				}
			}
			item[timeStamp].RAMGiBUsage = v / float64(stepSeconds) * values.HourInSeconds
		}
	}
}

func queryPodTotalCostsWithTimeRange(tenantId, clusterId string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	totalCosts := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlPodTotalCostFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) pod costs error:%v", clusterId, err)
		return nil, err
	}
	for _, pod := range ret {
		key := generatePodNamespaceNameKey(pod.Metric)
		totalCosts[key] = make(map[int64]float64)
		for _, v := range pod.Values {
			totalCosts[key][v.Timestamp.Unix()] = float64(v.Value)
		}
	}

	return totalCosts, nil
}

func queryPodResourceRequest(tenantId, clusterId string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuRequest := make(map[string]map[int64]float64)
	ramRequest := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlPodResourceRequestFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) pod resource request error:%v", err)
		return nil, nil, err
	}
	for _, pod := range ret {
		key := generatePodNamespaceNameKey(pod.Metric)
		resourceType := pod.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuRequest[key] = make(map[int64]float64)
			for _, v := range pod.Values {
				cpuRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramRequest[key] = make(map[int64]float64)
			for _, v := range pod.Values {
				ramRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}
	return cpuRequest, ramRequest, nil
}

func queryPodResourceUsage(tenantId, clusterId string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuUsage := make(map[string]map[int64]float64)
	ramUsage := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlPodResourceUsageFromClusterWithTimeRange, clusterId, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) pod resoruce usage error:%v", err)
		return nil, nil, err
	}
	for _, pod := range ret {
		key := generatePodNamespaceNameKey(pod.Metric)
		resourceType := pod.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuUsage[key] = make(map[int64]float64)
			for _, v := range pod.Values {
				cpuUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramUsage[key] = make(map[int64]float64)
			for _, v := range pod.Values {
				ramUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}

	return cpuUsage, ramUsage, nil
}

func queryHighLevelWorkloadCostsWithTimeRange(tenantId, clusterId string, start, end, stepSeconds int64, aggregateBy string) ([]*api.ClusterWorkloadCost, error) {
	queryRe := aggregateBy
	if aggregateBy == api.AggregateByAll {
		queryRe = "deployment|statefulset|daemonset"
	}

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
		totalCosts, err = queryHighLevelWorkloadTotalCost(tenantId, clusterId, queryRe, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		podCount, err = queryHighLevelWorkloadPodCount(tenantId, clusterId, queryRe, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuRequest, ramRequest, err = queryHighLevelWorkloadResourceRequest(tenantId, clusterId, queryRe, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuUsage, ramUsage, err = queryHighLevelWorkloadResourceUsage(tenantId, clusterId, queryRe, start, end, stepSeconds)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}
	workloadCost := make(map[string]map[int64]*api.ClusterWorkloadCostDetail)
	parseHighLevelWorkloadTotalCost(workloadCost, totalCosts)
	parseHighLevelWorkloadPodCount(workloadCost, podCount, stepSeconds)
	parseHighLevelWorkloadResourceRequest(workloadCost, cpuRequest, ramRequest, stepSeconds)
	parseHighLevelWorkloadResourceUsage(workloadCost, cpuUsage, ramUsage, stepSeconds)

	return convertClusterWorkloadCostToList(workloadCost), nil
}

func parseHighLevelWorkloadTotalCost(workloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	totalCosts map[string]map[int64]float64) {
	for workload, details := range totalCosts {
		workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		for timeStamp, v := range details {
			workloadCost[workload][timeStamp] = &api.ClusterWorkloadCostDetail{
				Timestamp: timeStamp,
				TotalCost: v,
			}
		}
	}
}

func parseHighLevelWorkloadPodCount(workloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	podCount map[string]map[int64]float64, stepSeconds int64) {
	for workload, details := range podCount {
		item, ok := workloadCost[workload]
		if !ok {
			workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].PodCount = v / (float64(stepSeconds) / values.MetricsPeriodInSeconds)
		}
	}
}

func parseHighLevelWorkloadResourceRequest(workloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	cpuRequest map[string]map[int64]float64,
	ramRequest map[string]map[int64]float64, stepSeconds int64) {
	for workload, details := range cpuRequest {
		item, ok := workloadCost[workload]
		if !ok {
			workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].CPUCoreRequest = v / (float64(stepSeconds) / values.HourInSeconds)
		}
	}

	for workload, details := range ramRequest {
		item, ok := workloadCost[workload]
		if !ok {
			workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].RAMGiBRequest = v / (float64(stepSeconds) / values.HourInSeconds)
		}
	}
}

func parseHighLevelWorkloadResourceUsage(workloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail,
	cpuUsage map[string]map[int64]float64,
	ramUsage map[string]map[int64]float64, stepSeconds int64) {
	for workload, details := range cpuUsage {
		item, ok := workloadCost[workload]
		if !ok {
			workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].CPUCoreUsage = v / (float64(stepSeconds) / values.HourInSeconds)
		}
	}

	for workload, details := range ramUsage {
		item, ok := workloadCost[workload]
		if !ok {
			workloadCost[workload] = make(map[int64]*api.ClusterWorkloadCostDetail)
		}
		for timeStamp, v := range details {
			if _, ok := item[timeStamp]; !ok {
				item[timeStamp] = &api.ClusterWorkloadCostDetail{
					Timestamp: timeStamp,
				}
			}
			item[timeStamp].RAMGiBUsage = v / (float64(stepSeconds) / values.HourInSeconds)
		}
	}
}

func queryHighLevelWorkloadTotalCost(tenantId, clusterId, queryRe string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	totalCosts := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlWorkloadTotalCostFromClusterWithTimeRange, clusterId, queryRe, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) total workload costs error:%v", err)
		return nil, err
	}
	for _, workload := range ret {
		key := generateWorkloadNamespaceNameKey(workload.Metric)
		totalCosts[key] = make(map[int64]float64)
		for _, v := range workload.Values {
			totalCosts[key][v.Timestamp.Unix()] = float64(v.Value)
		}
	}
	return totalCosts, nil
}

func queryHighLevelWorkloadPodCount(tenantId, clusterId, queryRe string, start, end, stepSeconds int64) (map[string]map[int64]float64, error) {
	podCount := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlWorkloadPodFromClusterWithTimeRange, clusterId, queryRe, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) pod count error:%v", err)
		return nil, err
	}
	for _, workload := range ret {
		key := generateWorkloadNamespaceNameKey(workload.Metric)
		podCount[key] = make(map[int64]float64)
		for _, v := range workload.Values {
			podCount[key][v.Timestamp.Unix()] = float64(v.Value)
		}
	}
	return podCount, nil
}

func queryHighLevelWorkloadResourceRequest(tenantId, clusterId, queryRe string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuRequest := make(map[string]map[int64]float64)
	ramRequest := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlWorkloadResourceRequestFromClusterWithTimeRange, clusterId, queryRe, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource request error:%v", err)
		return nil, nil, err
	}
	for _, workload := range ret {
		key := generateWorkloadNamespaceNameKey(workload.Metric)
		resourceType := workload.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuRequest[key] = make(map[int64]float64)
			for _, v := range workload.Values {
				cpuRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramRequest[key] = make(map[int64]float64)
			for _, v := range workload.Values {
				ramRequest[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}

	return cpuRequest, ramRequest, nil
}

func queryHighLevelWorkloadResourceUsage(tenantId, clusterId, queryRe string,
	start, end, stepSeconds int64) (map[string]map[int64]float64, map[string]map[int64]float64, error) {
	cpuUsage := make(map[string]map[int64]float64)
	ramUsage := make(map[string]map[int64]float64)
	promql := fmt.Sprintf(query.QlWorkloadResourceUsageFromClusterWithTimeRange, clusterId, queryRe, stepSeconds)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource usage error:%v", err)
		return nil, nil, err
	}
	for _, workload := range ret {
		key := generateWorkloadNamespaceNameKey(workload.Metric)
		resourceType := workload.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		switch resourceType {
		case "cpu":
			cpuUsage[key] = make(map[int64]float64)
			for _, v := range workload.Values {
				cpuUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		case "memory":
			ramUsage[key] = make(map[int64]float64)
			for _, v := range workload.Values {
				ramUsage[key][v.Timestamp.Unix()] = float64(v.Value)
			}
		}
	}

	return cpuUsage, ramUsage, nil
}

func convertClusterWorkloadCostToList(workloadCost map[string]map[int64]*api.ClusterWorkloadCostDetail) []*api.ClusterWorkloadCost {
	ret := []*api.ClusterWorkloadCost{}
	for workloadKey, details := range workloadCost {
		workloadType, namesapce, name := parseWorkloadNamespaceNameType(workloadKey)
		cost := &api.ClusterWorkloadCost{
			Namespace:    namesapce,
			WorkloadName: name,
			WorkloadType: workloadType,
			CostList:     []*api.ClusterWorkloadCostDetail{},
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
