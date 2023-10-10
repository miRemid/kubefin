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
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/query"
	"github.com/kubefin/kubefin/pkg/utils"
	"github.com/kubefin/kubefin/pkg/values"
	"github.com/prometheus/common/model"
)

func QueryAllClustersCurrentMonthCost(tenantId string) (map[string]*api.ClusterCostsSummary, error) {
	start, end, err := utils.GetCurrentMonthFirstLastDay()
	if err != nil {
		klog.Errorf("Query current time error:%v", err)
		return nil, err
	}

	var allClustersActiveTime map[string]float64
	var monthCostCurrent map[string]float64
	var cpuTotalCost map[string]float64
	var cpuTotalCount map[string]float64

	var errs []error
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		allClustersActiveTime, err = queryAllClustersActiveTime(tenantId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		monthCostCurrent, err = queryAllClustersCurrentMonthCost(tenantId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuTotalCost, err = queryAllClustersCPUTotalCost(tenantId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuTotalCount, err = queryAllClustersCPUTotalCount(tenantId, start, end)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	clusterCostSummary := make(map[string]*api.ClusterCostsSummary)
	for clusterId, activeTime := range allClustersActiveTime {
		costCurrent, cpuCost, cpuCount := monthCostCurrent[clusterId], cpuTotalCost[clusterId], cpuTotalCount[clusterId]
		clusterCostSummary[clusterId] = &api.ClusterCostsSummary{
			ClusterMonthCostCurrent:  costCurrent,
			ClusterMonthEstimateCost: 730 * costCurrent / (activeTime / values.HourInSeconds),
			ClusterAvgDailyCost:      24 * costCurrent / (activeTime / values.HourInSeconds),
			ClusterAvgHourlyCoreCost: cpuCost / cpuCount,
		}
	}

	return clusterCostSummary, nil
}

func queryAllClustersActiveTime(tenantId string, start, end int64) (map[string]float64, error) {
	allClustersActiveTime := make(map[string]float64)
	allClustersProperty, err := QueryAllClustersBasicProperty(tenantId, start, end)
	if err != nil {
		return nil, err
	}
	for clusterId, sample := range allClustersProperty {
		allClustersActiveTime[clusterId] = sample.ClusterActiveTime
	}
	return allClustersActiveTime, nil
}

func queryAllClustersCurrentMonthCost(tenantId string, start, end int64) (map[string]float64, error) {
	monthCostCurrent := make(map[string]float64)
	promql := fmt.Sprintf(query.QlNodesTotalCostsWithTimeRange, end-start)
	allCotalCost, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query clusters current month cost error:%v,%s", err, promql)
		return nil, err
	}

	for _, t := range allCotalCost {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		monthCostCurrent[string(clusterId)] = float64(t.Value)
	}
	return monthCostCurrent, nil
}

func queryAllClustersCPUTotalCost(tenantId string, start, end int64) (map[string]float64, error) {
	cpuTotalCost := make(map[string]float64)
	promql := fmt.Sprintf(query.QlNodeCPUTotalCostWithTimeRange, end-start)
	allCpuCost, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query clusters current month cpu cost error:%v", err)
		return nil, err
	}
	for _, c := range allCpuCost {
		clusterId := c.Metric[model.LabelName(values.ClusterIdLabelKey)]
		cpuTotalCost[string(clusterId)] = float64(c.Value)
	}
	return cpuTotalCost, nil
}

func queryAllClustersCPUTotalCount(tenantId string, start, end int64) (map[string]float64, error) {
	cpuTotalCount := make(map[string]float64)
	promql := fmt.Sprintf(query.QlNodeCPUTotalCountWithTimeRange, corev1.ResourceCPU, end-start)
	allCpuCount, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query clusters current month cpu count error:%v,%s", err, promql)
		return nil, err
	}
	for _, c := range allCpuCount {
		clusterId := c.Metric[model.LabelName(values.ClusterIdLabelKey)]
		cpuTotalCount[string(clusterId)] = float64(c.Value)
	}
	return cpuTotalCount, nil
}

func QueryClusterCurrentMonthCost(tenantId, clusterId string) (*api.ClusterCostsSummary, error) {
	start, end, err := utils.GetCurrentMonthFirstLastDay()
	if err != nil {
		klog.Errorf("Query current time error:%v", err)
		return nil, err
	}

	var clusterActiveTime float64
	var monthCostCurrent float64
	var cpuTotalCost float64
	var cpuTotalCount float64

	var errs []error
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		clusterActiveTime, err = queryClusterActiveTime(tenantId, clusterId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		monthCostCurrent, err = queryClusterCurrentMonthCost(tenantId, clusterId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuTotalCost, err = queryClusterCPUTotalCost(tenantId, clusterId, start, end)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		cpuTotalCount, err = queryClusterCPUTotalCount(tenantId, clusterId, start, end)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	return &api.ClusterCostsSummary{
		ClusterMonthCostCurrent:  monthCostCurrent,
		ClusterMonthEstimateCost: 730 * monthCostCurrent / (clusterActiveTime / values.HourInSeconds),
		ClusterAvgDailyCost:      24 * monthCostCurrent / (clusterActiveTime / values.HourInSeconds),
		ClusterAvgHourlyCoreCost: cpuTotalCost / cpuTotalCount,
	}, nil
}

func queryClusterActiveTime(tenantId, clusterId string, start, end int64) (float64, error) {
	var clusterActiveTime float64
	clusterProperty, err := QueryClusterBasicProperty(tenantId, clusterId, start, end)
	if err != nil {
		return 0, err
	}
	clusterActiveTime = float64(clusterProperty.ClusterActiveTime)
	return clusterActiveTime, nil
}

func queryClusterCurrentMonthCost(tenantId, clusterId string, start, end int64) (float64, error) {
	monthCostCurrent := float64(0)
	promql := fmt.Sprintf(query.QlNodesTotalCostsFromClusterWithTimeRange, clusterId, end-start)
	totalCost, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) current month cost error:%v", clusterId, err)
		return 0, err
	}
	if len(totalCost) == 0 {
		return 0, nil
	}
	monthCostCurrent = float64(totalCost[0].Value)
	return monthCostCurrent, nil
}

func queryClusterCPUTotalCost(tenantId, clusterId string, start, end int64) (float64, error) {
	cpuTotalCost := float64(0)
	promql := fmt.Sprintf(query.QlNodeCPUTotalCostFromClusterWithTimeRange, clusterId, end-start)
	cpuCost, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) current month cpu cost error:%v", clusterId, err)
		return 0, err
	}
	if len(cpuCost) == 0 {
		return 0, nil
	}
	cpuTotalCost = float64(cpuCost[0].Value)
	return cpuTotalCost, nil
}

func queryClusterCPUTotalCount(tenantId, clusterId string, start, end int64) (float64, error) {
	cpuTotalCount := float64(0)
	promql := fmt.Sprintf(query.QlNodeResourceTotalCountFromClusterWithTimeRange, clusterId, corev1.ResourceCPU, end-start)
	cpuCount, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) current month cpu count error:%v", clusterId, err)
		return 0, err
	}
	if len(cpuCount) == 0 {
		return 0, nil
	}
	cpuTotalCount = float64(cpuCount[0].Value)
	return cpuTotalCount, nil
}
