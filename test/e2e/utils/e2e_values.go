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

const (
	E2ETestEndpoint = "http://127.0.0.1:8080"

	AllClusterMetricsSummaryPath      = "/api/v1/metrics/summary"
	SpecificClusterMetricsSummaryPath = "/api/v1/metrics/clusters/%s/summary"

	SpecificClusterCPUMetricsPath    = "/api/v1/metrics/clusters/%s/cpu"
	SpecificClusterMemoryMetricsPath = "/api/v1/metrics/clusters/%s/memory"

	AllClusterCostsSummaryPath      = "/api/v1/costs/summary"
	SpecificClusterCostsSummaryPath = "/api/v1/costs/clusters/%s/summary"

	SpecificClusterComputeCostsPath  = "/api/v1/costs/clusters/%s/compute"
	SpecificClusterWorkloadCostsPath = "/api/v1/costs/clusters/%s/workload"
)
