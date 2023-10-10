//go:build e2e
// +build e2e

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

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/test/e2e/utils"
)

func TestSpecificClusterAllWorkloadCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster all workload costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster all workload costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters all workload costs response error:%s", string(body))
	}
}

func TestSpecificClusterAllWorkloadCostsWithTime(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	timeBefore := timeNow.AddDate(0, 0, -1)
	timeNowStr := strconv.FormatInt(timeNow.Unix(), 10)
	timeBeforeStr := strconv.FormatInt(timeBefore.Unix(), 10)

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	path += "?startTime=" + timeBeforeStr + "&endTime=" + timeNowStr + "&stepSeconds=3600"
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster all workload costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster all workload costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters all workload costs response error:%s", string(body))
	}
}

func TestSpecificClusterPodsCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=pod", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster pods costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster workload costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters pods costs response error:%s", string(body))
	}
}

func TestSpecificClusterPodsCostsWithTime(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	timeBefore := timeNow.AddDate(0, 0, -1)
	timeNowStr := strconv.FormatInt(timeNow.Unix(), 10)
	timeBeforeStr := strconv.FormatInt(timeBefore.Unix(), 10)

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=pod", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	path += "&startTime=" + timeBeforeStr + "&endTime=" + timeNowStr + "&stepSeconds=3600"
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster pods costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster pods costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters pods costs response error:%s", string(body))
	}
}

func TestSpecificClusterDeploymentsCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=deployment", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster deployments costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster deployments costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters deployments costs response error:%s", string(body))
	}
}

func TestSpecificClusterDeploymentsCostsWithTime(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	timeBefore := timeNow.AddDate(0, 0, -1)
	timeNowStr := strconv.FormatInt(timeNow.Unix(), 10)
	timeBeforeStr := strconv.FormatInt(timeBefore.Unix(), 10)

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=deployment", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	path += "&startTime=" + timeBeforeStr + "&endTime=" + timeNowStr + "&stepSeconds=3600"
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster deployments costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster deployments costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters deployments costs response error:%s", string(body))
	}
}

func TestSpecificClusterStatefulSetsCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=statefulset", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster statefulsets costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster statefulsets costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters statefulsets costs response error:%s", string(body))
	}
}

func TestSpecificClusterStatefulSetsWithTime(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	timeBefore := timeNow.AddDate(0, 0, -1)
	timeNowStr := strconv.FormatInt(timeNow.Unix(), 10)
	timeBeforeStr := strconv.FormatInt(timeBefore.Unix(), 10)

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=statefulset", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	path += "&startTime=" + timeBeforeStr + "&endTime=" + timeNowStr + "&stepSeconds=3600"
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster statefulsets costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster statefulsets costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters statefulsets costs response error:%s", string(body))
	}
}

func TestSpecificClusterDaemonsetsCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=daemonset", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster daemonsets costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster daemonsets costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters daemonsets costs response error:%s", string(body))
	}
}

func TestSpecificClusterDaemonsetWithTime(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	timeBefore := timeNow.AddDate(0, 0, -1)
	timeNowStr := strconv.FormatInt(timeNow.Unix(), 10)
	timeBeforeStr := strconv.FormatInt(timeBefore.Unix(), 10)

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath+"?aggregateBy=daemonset", "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	path += "&startTime=" + timeBeforeStr + "&endTime=" + timeNowStr + "&stepSeconds=3600"
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster daemonsets costs error:%v, %d", err, code)
	}
	clusterWorkloadCosts := api.ClusterWorkloadCostList{}
	err = json.Unmarshal(body, &clusterWorkloadCosts)
	if err != nil {
		t.Fatalf("Marshal cluster daemonsets costs error:%v", err)
	}
	if !utils.ValidateSpecificClusterWorkloadCosts(&clusterWorkloadCosts) {
		t.Fatalf("Validate clusters daemonsets costs response error:%s", string(body))
	}
}

func TestNoneClusterWorkloadCosts(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterWorkloadCostsPath, "not-exists")
	_, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil {
		t.Fatalf("Get specific cluster compute costs error:%v", err)
	}
	if code != http.StatusNotFound {
		t.Fatalf("Get none cluster compute costs error:%d", code)
	}
}
