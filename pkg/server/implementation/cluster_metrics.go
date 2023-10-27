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
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/query"
	"github.com/kubefin/kubefin/pkg/values"
	"github.com/prometheus/common/model"
)

func QueryClusterMetricsSummaryWithTimeRange(tenantId, clusterId string,
	resourceType v1.ResourceName, start, end, stepSeconds int64) (*api.ClusterResourceMetrics, error) {
	var usage []model.SamplePair
	var total []model.SamplePair
	var available []model.SamplePair
	var request []model.SamplePair
	var systemTaken []model.SamplePair

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(5)
	go func() {
		defer wg.Done()
		total, err = queryClusterResourceTotalWithTimeRange(tenantId, clusterId, resourceType, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		available, err = queryClusterResourceAvailableWithTimeRange(tenantId, clusterId, resourceType, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		systemTaken, err = queryClusterResourceSystemTakenWithTimeRange(tenantId, clusterId, resourceType, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		request, err = queryClusterResourceRequestWithTimeRange(tenantId, clusterId, resourceType, start, end, stepSeconds)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		usage, err = queryClusterResourceUsageWithTimeRange(tenantId, clusterId, resourceType, start, end, stepSeconds)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	unit := "core"
	if resourceType == v1.ResourceMemory {
		unit = "gb"
	}
	return &api.ClusterResourceMetrics{
		ClusterId:                 clusterId,
		ResourceType:              string(resourceType),
		Unit:                      unit,
		ResourceTotalValues:       total,
		ResourceSystemTakenValues: systemTaken,
		ResourceAvailableValues:   available,
		ResourceRequestValues:     request,
		ResourceUsageValues:       usage,
	}, nil
}

func queryClusterResourceTotalWithTimeRange(tenantId, clusterId string,
	resourceType v1.ResourceName, start, end, stepSeconds int64) ([]model.SamplePair, error) {
	var total []model.SamplePair
	promql := fmt.Sprintf(query.QlSumNodesResourceTotalFromCluster, clusterId, resourceType)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource total error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) != 1 {
		err := fmt.Errorf("query cluster(%s) resource total error, the data not correct", clusterId)
		klog.Errorf("%v", err)
		return nil, err
	}
	total = ret[0].Values

	return total, nil
}

func queryClusterResourceAvailableWithTimeRange(tenantId, clusterId string, resourceType v1.ResourceName, start, end, stepSeconds int64) ([]model.SamplePair, error) {
	var capacity []model.SamplePair
	promql := fmt.Sprintf(query.QlSumNodesResourceAvailableFromCluster, clusterId, resourceType)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource available error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) != 1 {
		err := fmt.Errorf("query cluster(%s) resource available error, the data not correct", clusterId)
		klog.Errorf("%v", err)
		return nil, err
	}
	capacity = ret[0].Values
	return capacity, nil
}

func queryClusterResourceSystemTakenWithTimeRange(tenantId, clusterId string, resourceType v1.ResourceName, start, end, stepSeconds int64) ([]model.SamplePair, error) {
	var systemTaken []model.SamplePair
	promql := fmt.Sprintf(query.QlSumNodesResourceSystemTakenFromCluster, clusterId, resourceType)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource system takne error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) != 1 {
		err := fmt.Errorf("query cluster(%s) resource system takne error, the data not correct", clusterId)
		klog.Errorf("%v", err)
		return nil, err
	}
	systemTaken = ret[0].Values
	return systemTaken, nil
}

func queryClusterResourceRequestWithTimeRange(tenantId, clusterId string, resourceType v1.ResourceName, start, end, stepSeconds int64) ([]model.SamplePair, error) {
	var request []model.SamplePair
	promql := fmt.Sprintf(query.QlSumPodResourceRequestFromCluster, clusterId, resourceType)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource request error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) != 1 {
		err := fmt.Errorf("query cluster(%s) resource request error, the data not correct", clusterId)
		klog.Errorf("%v", err)
		return nil, err
	}
	request = ret[0].Values

	return request, nil
}

func queryClusterResourceUsageWithTimeRange(tenantId, clusterId string, resourceType v1.ResourceName, start, end, stepSeconds int64) ([]model.SamplePair, error) {
	var usage []model.SamplePair
	promql := fmt.Sprintf(query.QlSumNodesResourceUsageFromCluster, clusterId, resourceType)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryRangeWithStep(promql, start, end, stepSeconds)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource usage error:%v", clusterId, err)
		return nil, err
	}
	if len(ret) != 1 {
		err := fmt.Errorf("query cluster(%s) resource usage error, the data not correct", clusterId)
		klog.Errorf("%v", err)
		return nil, err
	}
	usage = ret[0].Values

	return usage, nil
}

func QueryAllClustersBasicProperty(tenantId string, start, end int64) (map[string]*api.ClusterBasicProperty, error) {
	var err error
	var allClustersActiveTime []*model.Sample
	queryAllClustersActiveTimeFunc := func() error {
		promql := fmt.Sprintf(query.QlAllClustersActiveTime, end-start)
		allClustersActiveTime, err = query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
		if err != nil {
			klog.Errorf("Query cluster activity data error:%v", err)
			return err
		}

		return nil
	}

	var allClustersLastActiveInfo []*model.Sample
	queryAllClustersLastActiveFunc := func() error {
		promql := fmt.Sprintf(query.QlAllClustersActivity)
		allClustersLastActiveInfo, err = query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
		if err != nil {
			klog.Errorf("Query cluster last active time data error:%v", err)
			return err
		}
		return nil
	}

	var wg sync.WaitGroup
	var errs []error

	wg.Add(2)
	go func() {
		defer wg.Done()
		errs = append(errs, queryAllClustersActiveTimeFunc())
	}()
	go func() {
		defer wg.Done()
		errs = append(errs, queryAllClustersLastActiveFunc())
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	return ParseMultiClustersBasicProperty(allClustersActiveTime, allClustersLastActiveInfo), nil
}

func QueryClusterBasicProperty(tenantId, clusterId string, start, end int64) (*api.ClusterBasicProperty, error) {
	var err error
	var clusterActiveTime []*model.Sample
	queryClusterActiveTimeFunc := func() error {
		promql := fmt.Sprintf(query.QlClusterActiveTime, clusterId, end-start)
		clusterActiveTime, err = query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
		if err != nil {
			klog.Errorf("Query cluster activity data error:%v", err)
			return err
		}
		if len(clusterActiveTime) == 0 {
			err := fmt.Errorf("no such cluster:%s", clusterId)
			klog.Errorf("%v", err)
			return err
		}

		return nil
	}

	var clusterLastActiveInfo []*model.Sample
	queryClusterLastActiveFunc := func() error {
		promql := fmt.Sprintf(query.QlClusterActivity, clusterId)
		clusterLastActiveInfo, err = query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
		if err != nil {
			klog.Errorf("Query cluster activity data error:%v", err)
			return err
		}
		if len(clusterLastActiveInfo) == 0 {
			err := fmt.Errorf("no such cluster:%s", clusterId)
			klog.Errorf("%v", err)
			return err
		}

		return nil
	}

	var wg sync.WaitGroup
	var errs []error

	wg.Add(2)
	go func() {
		defer wg.Done()
		errs = append(errs, queryClusterActiveTimeFunc())
	}()
	go func() {
		defer wg.Done()
		errs = append(errs, queryClusterLastActiveFunc())
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	return ParseClusterBasicProperty(clusterActiveTime[0], clusterLastActiveInfo[0]), nil
}

func ParseClusterBasicProperty(clusterActiveTime *model.Sample, clusterLasterActive *model.Sample) *api.ClusterBasicProperty {
	clusterId := string(clusterActiveTime.Metric[model.LabelName(values.ClusterIdLabelKey)])
	clusterName := string(clusterActiveTime.Metric[model.LabelName(values.ClusterNameLabelKey)])
	provider := string(clusterActiveTime.Metric[model.LabelName(values.CloudProviderLabelKey)])
	regin := string(clusterActiveTime.Metric[model.LabelName(values.RegionLabelKey)])

	timeNow := time.Now().Unix()
	lastActiveTime := clusterLasterActive.Timestamp.Unix()

	clusterState := "running"
	if timeNow-lastActiveTime > int64(values.LostConnectionTimeoutThreshold) {
		clusterState = "connect_failed"
	}

	return &api.ClusterBasicProperty{
		ClusterId:             clusterId,
		ClusterName:           clusterName,
		CloudProvider:         provider,
		ClusterRegion:         regin,
		ClusterActiveTime:     float64(clusterActiveTime.Value),
		LastActiveTime:        lastActiveTime,
		ClusterConnectionSate: clusterState,
	}
}

func ParseMultiClustersBasicProperty(clusterActiveTime []*model.Sample, clusterLasterActive []*model.Sample) map[string]*api.ClusterBasicProperty {
	clustersInfo := make(map[string]*api.ClusterBasicProperty)
	for _, sampleActiveTime := range clusterActiveTime {
		clusterIdx := string(sampleActiveTime.Metric[model.LabelName(values.ClusterIdLabelKey)])
		hasMatchData := false
		for _, sameLastActive := range clusterLasterActive {
			clusterIdy := string(sameLastActive.Metric[model.LabelName(values.ClusterIdLabelKey)])
			if clusterIdx != clusterIdy {
				continue
			}
			hasMatchData = true
			clustersInfo[clusterIdx] = ParseClusterBasicProperty(sampleActiveTime, sameLastActive)
		}

		// This should happen generally
		if !hasMatchData {
			klog.Warningf("Cluster information is not correct:%s, ignore it")
			continue
		}
	}

	return clustersInfo
}

func ConvertToMultiClustersCostsList(clustersSummary map[string]*api.ClusterCostsSummary, clustersProperty map[string]*api.ClusterBasicProperty) *api.ClusterCostsSummaryList {
	for clusterId := range clustersSummary {
		clustersSummary[clusterId].ClusterBasicProperty = *clustersProperty[clusterId]
	}

	retList := &api.ClusterCostsSummaryList{}
	for _, cluster := range clustersSummary {
		retList.Items = append(retList.Items, cluster)
	}

	return retList
}

func QueryAllClustersCurrentMetrics(tenantId string) (map[string]*api.ClusterMetricsSummary, error) {
	var nodesNumber map[string]map[string]int64
	var podsNumber map[string]map[string]int64
	var resourceTotal map[string]map[string]float64
	var resourceUsage map[string]map[string]float64
	var resourceRequest map[string]map[string]float64
	var resourceAvailable map[string]map[string]float64
	var resourceSystemTaken map[string]map[string]float64

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(7)
	go func() {
		defer wg.Done()
		nodesNumber, err = queryAllClustersNodesNumer(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		podsNumber, err = queryAllClustersPodsNumber(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceTotal, err = queryAllClustersResourceTotal(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceUsage, err = queryAllClustersResourceUsage(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceRequest, err = queryAllClustersResourceRequest(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceAvailable, err = queryAllClustersResourceAvailable(tenantId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceSystemTaken, err = queryAllClustersResourceSystemTaken(tenantId)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}

	ret := make(map[string]*api.ClusterMetricsSummary)
	parseAllClustersNodesNumber(ret, nodesNumber)
	parseAllClustersPodsNumber(ret, podsNumber)
	parseAllClustersResourceTotal(ret, resourceTotal)
	parseAllClustersResourceUsage(ret, resourceUsage)
	parseAllClustersResourceRequest(ret, resourceRequest)
	parseAllClustersResourceAvailable(ret, resourceAvailable)
	parseAllClustersResourceSystemTaken(ret, resourceSystemTaken)

	return ret, nil
}

func parseAllClustersNodesNumber(data map[string]*api.ClusterMetricsSummary,
	nodesNumber map[string]map[string]int64) {
	for clusterId, v := range nodesNumber {
		totalNodes := int64(0)
		ondemandNodes := int64(0)
		spotNodes := int64(0)
		for billingMode, num := range v {
			totalNodes += int64(num)
			switch billingMode {
			case values.BillingModeOnDemand:
				ondemandNodes += int64(num)
			case values.BillingModeSpot:
				spotNodes += int64(num)
			}
		}
		data[clusterId] = &api.ClusterMetricsSummary{
			NodeNumbersCurrent:                totalNodes,
			OnDemandBillingNodeNumbersCurrent: ondemandNodes,
			SpotBillingNodeNumbersCurrent:     spotNodes,
		}
	}
}

func parseAllClustersPodsNumber(data map[string]*api.ClusterMetricsSummary,
	podsNumber map[string]map[string]int64) {
	for clusterId, v := range podsNumber {
		podTotal := int64(0)
		podScheduled := int64(0)
		podUnscheduled := int64(0)
		for scheduled, num := range v {
			podTotal += int64(num)
			switch scheduled {
			case "true":
				podScheduled += int64(num)
			case "false":
				podUnscheduled += int64(num)
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].PodTotalCurrent = podTotal
		data[clusterId].PodScheduledCurrent = podScheduled
		data[clusterId].PodUnscheduledCurrent = podUnscheduled
	}
}

func parseAllClustersResourceTotal(data map[string]*api.ClusterMetricsSummary,
	resourceTotal map[string]map[string]float64) {
	for clusterId, v := range resourceTotal {
		cpuTotal := float64(0.0)
		ramTotal := float64(0.0)
		for resourceType, num := range v {
			switch resourceType {
			case "cpu":
				cpuTotal += num
			case "memory":
				ramTotal += num
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].CPUCoreTotal = cpuTotal
		data[clusterId].RAMGiBTotal = ramTotal
	}
}

func parseAllClustersResourceUsage(data map[string]*api.ClusterMetricsSummary,
	resourceUsage map[string]map[string]float64) {
	for clusterId, v := range resourceUsage {
		cpuUsage := float64(0.0)
		ramUsage := float64(0.0)
		for resourceType, num := range v {
			switch resourceType {
			case "cpu":
				cpuUsage += num
			case "memory":
				ramUsage += num
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].CPUCoreUsage = cpuUsage
		data[clusterId].RAMGiBUsage = ramUsage
	}
}

func parseAllClustersResourceRequest(data map[string]*api.ClusterMetricsSummary,
	resourceRequest map[string]map[string]float64) {
	for clusterId, v := range resourceRequest {
		cpuRequest := float64(0.0)
		ramRequest := float64(0.0)
		for resourceType, num := range v {
			switch resourceType {
			case "cpu":
				cpuRequest += num
			case "memory":
				ramRequest += num
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].CPUCoreRequest = cpuRequest
		data[clusterId].RAMGiBRequest = ramRequest
	}
}

func parseAllClustersResourceAvailable(data map[string]*api.ClusterMetricsSummary,
	resourceAvailable map[string]map[string]float64) {
	for clusterId, v := range resourceAvailable {
		cpuAvailable := float64(0.0)
		ramAvailable := float64(0.0)
		for resourceType, num := range v {
			switch resourceType {
			case "cpu":
				cpuAvailable += num
			case "memory":
				ramAvailable += num
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].CPUCoreAvailable = cpuAvailable
		data[clusterId].RAMGiBAvailable = ramAvailable
	}
}

func parseAllClustersResourceSystemTaken(data map[string]*api.ClusterMetricsSummary,
	resourceSystemTaken map[string]map[string]float64) {
	for clusterId, v := range resourceSystemTaken {
		cpuSystemTaken := float64(0.0)
		ramSystemTaken := float64(0.0)
		for resourceType, num := range v {
			switch resourceType {
			case "cpu":
				cpuSystemTaken += num
			case "memory":
				ramSystemTaken += num
			}
		}
		if _, ok := data[clusterId]; !ok {
			data[clusterId] = &api.ClusterMetricsSummary{}
		}
		data[clusterId].CPUCoreSystemTaken = cpuSystemTaken
		data[clusterId].RAMGiBSystemTaken = ramSystemTaken
	}
}

func queryAllClustersNodesNumer(tenantId string) (map[string]map[string]int64, error) {
	// maps [cluster id][billing mode]count
	nodesNumber := make(map[string]map[string]int64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlNodesNumber)
	if err != nil {
		klog.Errorf("Query all clusters nodes error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := nodesNumber[string(clusterId)]; !ok {
			nodesNumber[string(clusterId)] = map[string]int64{}
		}
		billingMode := t.Metric[model.LabelName(values.BillingModeLabelKey)]
		nodesNumber[string(clusterId)][string(billingMode)] = int64(t.Value)
	}
	return nodesNumber, nil
}

func queryAllClustersPodsNumber(tenantId string) (map[string]map[string]int64, error) {
	// maps [cluster id][schedule status]count
	podsNumber := make(map[string]map[string]int64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlPodsNumber)
	if err != nil {
		klog.Errorf("Query all clusters pods error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		scheduleStatus := t.Metric[model.LabelName(values.PodScheduledKey)]
		podsNumber[string(clusterId)] = map[string]int64{
			string(scheduleStatus): int64(t.Value),
		}
	}
	return podsNumber, nil
}

func queryAllClustersResourceTotal(tenantId string) (map[string]map[string]float64, error) {
	// maps [cluster id][cpu/memory]float64
	resourceTotal := make(map[string]map[string]float64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlResourceTotal)
	if err != nil {
		klog.Errorf("Query all clusters resource total error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := resourceTotal[string(clusterId)]; !ok {
			resourceTotal[string(clusterId)] = map[string]float64{}
		}
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceTotal[string(clusterId)][string(resourceType)] = float64(t.Value)
	}
	return resourceTotal, nil
}

func queryAllClustersResourceUsage(tenantId string) (map[string]map[string]float64, error) {
	// maps [cluster id][cpu/memory]float64
	resourceUsage := make(map[string]map[string]float64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlResourceUsage)
	if err != nil {
		klog.Errorf("Query all clusters resource usage error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := resourceUsage[string(clusterId)]; !ok {
			resourceUsage[string(clusterId)] = map[string]float64{}
		}
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceUsage[string(clusterId)][string(resourceType)] = float64(t.Value)
	}
	return resourceUsage, nil
}

func queryAllClustersResourceRequest(tenantId string) (map[string]map[string]float64, error) {
	// maps [cluster id][cpu/memory]float64
	resourceRequest := make(map[string]map[string]float64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlResourceRequest)
	if err != nil {
		klog.Errorf("Query all clusters resource request error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := resourceRequest[string(clusterId)]; !ok {
			resourceRequest[string(clusterId)] = map[string]float64{}
		}
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceRequest[string(clusterId)][string(resourceType)] = float64(t.Value)
	}
	return resourceRequest, nil
}

func queryAllClustersResourceAvailable(tenantId string) (map[string]map[string]float64, error) {
	resourceAvailable := make(map[string]map[string]float64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlResourceAvailable)
	if err != nil {
		klog.Errorf("Query all clusters resource available error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := resourceAvailable[string(clusterId)]; !ok {
			resourceAvailable[string(clusterId)] = map[string]float64{}
		}
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceAvailable[string(clusterId)][string(resourceType)] = float64(t.Value)
	}
	return resourceAvailable, nil
}

func queryAllClustersResourceSystemTaken(tenantId string) (map[string]map[string]float64, error) {
	resourceSystemTaken := make(map[string]map[string]float64)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(query.QlResoruceSystemTaken)
	if err != nil {
		klog.Errorf("Query all clusters resource system taken error:%v", err)
		return nil, err
	}
	for _, t := range ret {
		clusterId := t.Metric[model.LabelName(values.ClusterIdLabelKey)]
		if _, ok := resourceSystemTaken[string(clusterId)]; !ok {
			resourceSystemTaken[string(clusterId)] = map[string]float64{}
		}
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceSystemTaken[string(clusterId)][string(resourceType)] = float64(t.Value)
	}
	return resourceSystemTaken, nil
}

func ConvertToMultiClustersMetricsList(clusterMetrics map[string]*api.ClusterMetricsSummary, clustersProperty map[string]*api.ClusterBasicProperty) *api.ClusterMetricsSummaryList {
	for cluterId := range clusterMetrics {
		clusterMetrics[cluterId].ClusterBasicProperty = *clustersProperty[cluterId]
	}

	retList := &api.ClusterMetricsSummaryList{}
	for _, v := range clusterMetrics {
		retList.Items = append(retList.Items, *v)
	}

	return retList
}

func QueryClusterCurrentMetrics(tenantId, clusterId string) (*api.ClusterMetricsSummary, error) {
	var nodesNumber map[string]int64
	var podsNumber map[string]int64
	var resourceTotal map[string]float64
	var resourceUsage map[string]float64
	var resourceRequest map[string]float64
	var resourceAvailable map[string]float64
	var resourceSystemTaken map[string]float64

	var wg sync.WaitGroup
	var err error
	var errs []error

	wg.Add(7)
	go func() {
		defer wg.Done()
		nodesNumber, err = queryClusterNodesNumber(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		podsNumber, err = queryClusterPodsNumber(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceTotal, err = queryClusterResourceTotal(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceUsage, err = queryClusterResoruceUsage(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceRequest, err = queryClusterResoruceRequest(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceAvailable, err = queryClusterResourceAvailable(tenantId, clusterId)
		errs = append(errs, err)
	}()
	go func() {
		defer wg.Done()
		resourceSystemTaken, err = queryClusterResourceSystemTaken(tenantId, clusterId)
		errs = append(errs, err)
	}()

	wg.Wait()
	if errors.NewAggregate(errs) != nil {
		return nil, errors.NewAggregate(errs)
	}
	ret := &api.ClusterMetricsSummary{}
	parseClusterNodesNumber(ret, nodesNumber)
	parseClusterPodsNumber(ret, podsNumber)
	parseClusterResourceTotal(ret, resourceTotal)
	parseClusterResourceUsage(ret, resourceUsage)
	parseClusterResourceRequest(ret, resourceRequest)
	parseClusterResourceAvailable(ret, resourceAvailable)
	parseClusterResourceSystemTaken(ret, resourceSystemTaken)

	return ret, nil
}

func parseClusterNodesNumber(data *api.ClusterMetricsSummary, nodesNumber map[string]int64) {
	for billingMode, num := range nodesNumber {
		data.NodeNumbersCurrent += int64(num)
		switch billingMode {
		case values.BillingModeOnDemand:
			data.OnDemandBillingNodeNumbersCurrent += int64(num)
		case values.BillingModeSpot:
			data.SpotBillingNodeNumbersCurrent += int64(num)
		}
	}
}

func parseClusterPodsNumber(data *api.ClusterMetricsSummary, podsNumber map[string]int64) {
	for scheduled, num := range podsNumber {
		data.PodTotalCurrent += int64(num)
		switch scheduled {
		case "true":
			data.PodScheduledCurrent += int64(num)
		case "false":
			data.PodUnscheduledCurrent += int64(num)
		}
	}
}

func parseClusterResourceTotal(data *api.ClusterMetricsSummary, resourceTotal map[string]float64) {
	for resourceType, num := range resourceTotal {
		switch resourceType {
		case "cpu":
			data.CPUCoreTotal += num
		case "memory":
			data.RAMGiBTotal += num
		}
	}
}

func parseClusterResourceUsage(data *api.ClusterMetricsSummary, resourceUsage map[string]float64) {
	for resourceType, num := range resourceUsage {
		switch resourceType {
		case "cpu":
			data.CPUCoreUsage += num
		case "memory":
			data.RAMGiBUsage += num
		}
	}
}

func parseClusterResourceRequest(data *api.ClusterMetricsSummary, resourceRequest map[string]float64) {
	for resourceType, num := range resourceRequest {
		switch resourceType {
		case "cpu":
			data.CPUCoreRequest += num
		case "memory":
			data.RAMGiBRequest += num
		}
	}
}

func parseClusterResourceAvailable(data *api.ClusterMetricsSummary, resourceAvailable map[string]float64) {
	for resourceType, num := range resourceAvailable {
		switch resourceType {
		case "cpu":
			data.CPUCoreAvailable += num
		case "memory":
			data.RAMGiBAvailable += num
		}
	}
}

func parseClusterResourceSystemTaken(data *api.ClusterMetricsSummary, resourceSystemTaken map[string]float64) {
	for resourceType, num := range resourceSystemTaken {
		switch resourceType {
		case "cpu":
			data.CPUCoreSystemTaken += num
		case "memory":
			data.RAMGiBSystemTaken += num
		}
	}
}

func queryClusterNodesNumber(tenantId, clusterId string) (map[string]int64, error) {
	// maps [billing mode]count
	nodesNumber := make(map[string]int64)
	promql := fmt.Sprintf(query.QlNodesNumberFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) nodes number error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		billiingMode := t.Metric[model.LabelName(values.BillingModeLabelKey)]
		nodesNumber[string(billiingMode)] = int64(t.Value)
	}
	return nodesNumber, nil
}

func queryClusterPodsNumber(tenantId, clusterId string) (map[string]int64, error) {
	// maps [schedule status]count
	podsNumber := make(map[string]int64)
	promql := fmt.Sprintf(query.QlPodsNumberFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) pods number error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		billingMode := t.Metric[model.LabelName(values.PodScheduledKey)]
		podsNumber[string(billingMode)] = int64(t.Value)
	}
	return podsNumber, nil
}

func queryClusterResourceTotal(tenantId, clusterId string) (map[string]float64, error) {
	// maps [cpu/memory]float64
	resourceTotal := make(map[string]float64)
	promql := fmt.Sprintf(query.QlResourceTotalFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource total error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceTotal[string(resourceType)] = float64(t.Value)
	}
	return resourceTotal, nil
}

func queryClusterResoruceUsage(tenantId, clusterId string) (map[string]float64, error) {
	// maps [cpu/memory]float64
	resourceUsage := make(map[string]float64)
	promql := fmt.Sprintf(query.QlResourceUsageFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource usage error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceUsage[string(resourceType)] = float64(t.Value)
	}
	return resourceUsage, nil
}

func queryClusterResoruceRequest(tenantId, clusterId string) (map[string]float64, error) {
	// maps [cpu/memory]float64
	resourceRequest := make(map[string]float64)
	promql := fmt.Sprintf(query.QlResourceRequestFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource request error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceRequest[string(resourceType)] = float64(t.Value)
	}
	return resourceRequest, nil
}

func queryClusterResourceAvailable(tenantId, clusterId string) (map[string]float64, error) {
	// maps [cpu/memory]float64
	resourceAvailale := make(map[string]float64)
	promql := fmt.Sprintf(query.QlResourceAvailableFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource capacity error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceAvailale[string(resourceType)] = float64(t.Value)
	}
	return resourceAvailale, nil
}

func queryClusterResourceSystemTaken(tenantId, clusterId string) (map[string]float64, error) {
	// maps [cpu/memory]float64
	resourceSystemTaken := make(map[string]float64)
	promql := fmt.Sprintf(query.QlResourceSystemTakenFromCluster, clusterId)
	ret, err := query.GetPromQueryClient().WithTenantId(tenantId).QueryInstant(promql)
	if err != nil {
		klog.Errorf("Query cluster(%s) resource capacity error:%v", clusterId, err)
		return nil, err
	}
	for _, t := range ret {
		resourceType := t.Metric[model.LabelName(values.ResourceTypeLabelKey)]
		resourceSystemTaken[string(resourceType)] = float64(t.Value)
	}
	return resourceSystemTaken, nil
}
