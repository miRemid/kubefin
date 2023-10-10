/*
Copyright 2022 The KubeFin Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicabl e law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"os"
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	baseconfig "k8s.io/component-base/config"

	"github.com/kubefin/kubefin/pkg/values"
)

type AgentOptions struct {
	LeaderElection       baseconfig.LeaderElectionConfiguration
	ScrapMetricsInterval time.Duration
	LeaderElectionID     string

	CloudProvider string
	ClusterName   string
	ClusterId     string

	CPUCoreReserved        string
	RAMGBReserved          string
	CPUMemoryCostRatio     string
	CustomCPUCoreHourPrice string
	CustomRAMGBHourPrice   string
}

// NewAgentOptions builds an empty options.
func NewAgentOptions() *AgentOptions {
	return &AgentOptions{
		LeaderElection: baseconfig.LeaderElectionConfiguration{
			ResourceLock:      resourcelock.LeasesResourceLock,
			ResourceNamespace: values.KubeFinNamespace,
			ResourceName:      values.KubeFinAgentName,
			LeaseDuration:     metav1.Duration{Duration: values.DefaultLeaseDuration},
			RenewDeadline:     metav1.Duration{Duration: values.DefaultRenewDeadline},
			RetryPeriod:       metav1.Duration{Duration: values.DefaultRetryPeriod},
		},
		ScrapMetricsInterval:   time.Second,
		LeaderElectionID:       os.Getenv(values.LeaderElectionIDEnv),
		CloudProvider:          os.Getenv(values.CloudProviderEnv),
		ClusterName:            os.Getenv(values.ClusterNameEnv),
		ClusterId:              os.Getenv(values.ClusterIdEnv),
		CPUMemoryCostRatio:     os.Getenv(values.CPUMemoryCostRatioEnv),
		CustomCPUCoreHourPrice: os.Getenv(values.CustomCPUCoreHourPriceEnv),
		CustomRAMGBHourPrice:   os.Getenv(values.CustomRAMGBHourPriceEnv),
		CPUCoreReserved:        os.Getenv(values.CPUCorReservedEnv),
		RAMGBReserved:          os.Getenv(values.RAMGBReservedEnv),
	}
}

func (o *AgentOptions) Complete() error {
	return nil
}

func (o *AgentOptions) Validate() error {
	return nil
}

func (o *AgentOptions) ApplyTo() {
}

func (o *AgentOptions) AddFlags(flags *pflag.FlagSet) {
}
