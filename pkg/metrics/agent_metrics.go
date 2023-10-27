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

package metrics

import (
	"context"
	"time"

	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/metrics/core"
)

type AgentMetricsCollector struct {
	agentOptions                  *options.AgentOptions
	ctx                           context.Context
	interval                      time.Duration
	clusterMetricsCollector       *core.ClusterLevelMetricsCollector
	nodeLevelMetricsCollector     *core.NodeLevelMetricsCollector
	podLevelMetricsCollector      *core.PodLevelMetricsCollector
	workloadLevelMetricsCollector *core.WorkloadLevelMetricsCollector
}

func NewAgentMetricsCollector(ctx context.Context,
	options *options.AgentOptions,
	coreResourceInformerLister *api.CoreResourceInformerLister,
	provider cloudprice.CloudProviderInterface,
	metricsClientSet *versioned.Clientset) *AgentMetricsCollector {
	return &AgentMetricsCollector{
		ctx:                       ctx,
		agentOptions:              options,
		interval:                  options.ScrapMetricsInterval,
		clusterMetricsCollector:   core.NewClusterLevelMetricsCollector(provider, coreResourceInformerLister.NodeLister),
		nodeLevelMetricsCollector: core.NewNodeLevelMetricsCollector(metricsClientSet, provider, coreResourceInformerLister),
		podLevelMetricsCollector: core.NewPodLevelMetricsCollector(
			metricsClientSet, provider,
			coreResourceInformerLister.PodLister,
			coreResourceInformerLister.NodeLister),
		workloadLevelMetricsCollector: core.NewWorkloadLevelMetricsCollector(
			metricsClientSet, provider,
			coreResourceInformerLister.PodLister,
			coreResourceInformerLister.NodeLister,
			coreResourceInformerLister.DaemonSetLister,
			coreResourceInformerLister.DeploymentLister,
			coreResourceInformerLister.StatefulSetLister),
	}
}

func (a *AgentMetricsCollector) StartAgentMetricsCollector() {
	go a.nodeLevelMetricsCollector.StartCollectNodeLevelMetrics(a.ctx, a.interval, a.agentOptions)
	go a.podLevelMetricsCollector.StartCollectPodLevelMetrics(a.ctx, a.interval, a.agentOptions)
	go a.workloadLevelMetricsCollector.StartCollectWorkloadLevelMetrics(a.ctx, a.interval, a.agentOptions)
	go a.clusterMetricsCollector.StartCollectClusterLevelMetrics(a.ctx, a.interval, a.agentOptions)
}
