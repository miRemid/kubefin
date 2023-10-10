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
	"testing"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/test/e2e/utils"
)

func TestClustersMetricsSummary(t *testing.T) {
	t.Parallel()

	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, utils.AllClusterMetricsSummaryPath)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get all clusters summary error:%v, %d", err, code)
	}
	allClustersSummary := api.ClusterMetricsSummaryList{}
	err = json.Unmarshal(body, &allClustersSummary)
	if err != nil {
		t.Fatalf("Marshal clusters summary error:%v", err)
	}
	if !utils.ValidateAllClustersMetricsSummary(&allClustersSummary) {
		t.Fatalf("Validate clusters summary response error:%s", string(body))
	}
}

func TestSpecificClusterMetricsSummary(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterMetricsSummaryPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster metrics summary error:%v, %d", err, code)
	}
	clusterSummary := api.ClusterMetricsSummary{}
	err = json.Unmarshal(body, &clusterSummary)
	if err != nil {
		t.Fatalf("Marshal cluster metrics summary error:%v", err)
	}
	if !utils.ValidateSpecificClusterMetricsSummary(&clusterSummary) {
		t.Fatalf("Validate cluster summary response error:%s", string(body))
	}
}

func TestNoneClusterMetricsSummary(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterMetricsSummaryPath, "not-exists")
	_, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil {
		t.Fatalf("Get specific cluster metrics summary error:%v", err)
	}
	if code != http.StatusNotFound {
		t.Fatalf("Get none cluster metrics summary error:%d", code)
	}
}

func TestClustersCostsSummary(t *testing.T) {
	t.Parallel()

	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, utils.AllClusterCostsSummaryPath)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get all clusters costs summary error:%v, %d", err, code)
	}
	allClustersCostsSummary := api.ClusterCostsSummaryList{}
	err = json.Unmarshal(body, &allClustersCostsSummary)
	if err != nil {
		t.Fatalf("Marshal clusters costs summary error:%v", err)
	}
	if !utils.ValidateAllClustersCostsSummary(&allClustersCostsSummary) {
		t.Fatalf("Validate clusters costs summary response error:%s", string(body))
	}
}

func TestSpecificClusterCostsSummary(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterCostsSummaryPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster costs summary error:%v, %d", err, code)
	}
	clusterSummary := api.ClusterCostsSummary{}
	err = json.Unmarshal(body, &clusterSummary)
	if err != nil {
		t.Fatalf("Marshal cluster costs summary error:%v", err)
	}
	if !utils.ValidateSpecificClusterCostsSummary(&clusterSummary) {
		t.Fatalf("Validate cluster summary response error:%s", string(body))
	}
}

func TestNoneClusterCostsSummary(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterCostsSummaryPath, "not-exists")
	_, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil {
		t.Fatalf("Get specific cluster costs summary error:%v, %d", err, code)
	}
	if code != http.StatusNotFound {
		t.Fatalf("Get none cluster costs summary error:%d", code)
	}
}
