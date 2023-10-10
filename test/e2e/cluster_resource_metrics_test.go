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

func TestSpecificClusterCPUMetricsRaw(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterCPUMetricsPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster cpu metrics error:%v, %d", err, code)
	}
	clusterResourceMetrics := api.ClusterResourceMetrics{}
	err = json.Unmarshal(body, &clusterResourceMetrics)
	if err != nil {
		t.Fatalf("Marshal cluster cpu metrics error:%v", err)
	}
	if !utils.ValidateSpecificClusterResourceMetrics(&clusterResourceMetrics) {
		t.Fatalf("Validate clusters cpu metrics response error:%s", string(body))
	}
}

func TestNoneClusterCPUMetricsRaw(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterCPUMetricsPath, "not-exists")
	_, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil {
		t.Fatalf("Get none cluster cpu metrics error:%v", err)
	}
	if code != http.StatusNotFound {
		t.Fatalf("Get none cluster cpu metrics error:%d", code)
	}
}

func TestSpecificClusterMemoryMetricsRaw(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterMemoryMetricsPath, "a4b52fe0-afd0-4050-9ecb-93edcadef48e")
	body, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil || code != http.StatusOK {
		t.Fatalf("Get specific cluster memory metrics error:%v, %d", err, code)
	}
	clusterResourceMetrics := api.ClusterResourceMetrics{}
	err = json.Unmarshal(body, &clusterResourceMetrics)
	if err != nil {
		t.Fatalf("Marshal cluster memory metrics error:%v", err)
	}
	if !utils.ValidateSpecificClusterResourceMetrics(&clusterResourceMetrics) {
		t.Fatalf("Validate clusters memory metrics response error:%s", string(body))
	}
}

func TestNoneClusterMemoryMetricsRaw(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(utils.SpecificClusterCPUMetricsPath, "not-exists")
	_, code, err := utils.DoGetRequest(utils.E2ETestEndpoint, path)
	if err != nil {
		t.Fatalf("Get specific cluster memory metrics error:%v", err)
	}
	if code != http.StatusNotFound {
		t.Fatalf("Get none cluster memory metrics error:%d", code)
	}
}
