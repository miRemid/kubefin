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

package values

import "time"

const (
	KubeFinNamespace = "kubefin"
	KubeFinAgentName = "kubefin-agent"

	// DefaultLeaseDuration is the defaultcloud LeaseDuration for leader election.
	DefaultLeaseDuration = 15 * time.Second
	// DefaultRenewDeadline is the defaultcloud RenewDeadline for leader election.
	DefaultRenewDeadline = 10 * time.Second
	// DefaultRetryPeriod is the defaultcloud RetryPeriod for leader election.
	DefaultRetryPeriod = 5 * time.Second

	LostConnectionTimeoutThreshold = time.Minute * 3 / time.Second

	GBInBytes              = 1024.0 * 1024.0 * 1024.0
	CoreInMCore            = 1000.0
	HourInSeconds          = 3600.0
	MetricsPeriodInSeconds = 15.0

	BillingModeOnDemand = "ondemand"
	BillingModeMonthly  = "monthly"
	BillingModeYearly   = "yearly"
	BillingModeSpot     = "spot"
	BillingModeFallback = "fallback"

	ClusterStateRunning        = "running"
	ClusterStateLostConnection = "connect_failed"

	CloudProviderEnv          = "CLOUD_PROVIDER"
	ClusterNameEnv            = "CLUSTER_NAME"
	ClusterIdEnv              = "CLUSTER_ID"
	LeaderElectionIDEnv       = "LEADER_ELECTION_ID"
	QueryBackendEndpointEnv   = "QUERY_BACKEND_ENDPOINT"
	NodeCPUDeviationEnv       = "NODE_CPU_DEVIATION"
	NodeRAMDeviationEnv       = "NODE_RAM_DEVIATION"
	CPUMemoryCostRatioEnv     = "CPUCORE_RAMGB_PRICE_RATIO"
	CustomCPUCoreHourPriceEnv = "CUSTOM_CPU_CORE_HOUR_PRICE"
	CustomRAMGBHourPriceEnv   = "CUSTOM_RAM_GB_HOUR_PRICE"

	MultiTenantHeader       = "X-Scope-OrgID"
	ClusterIdQueryParameter = "cluster_id"

	DefaultStepSeconds = 3600
	// DefaultDetailStepSeconds is used to show the fine-grained line chart of cpu/memory data
	DefaultDetailStepSeconds = 600
)
