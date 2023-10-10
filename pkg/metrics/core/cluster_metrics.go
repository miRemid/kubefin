/*
Copyright 2023 The KubeFin Authors

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

package core

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/values"
)

type ClusterLevelMetricsCollector struct {
	clusterActiveGV *prometheus.GaugeVec
	provider        cloudprice.CloudProviderInterface
	nodeLister      v1.NodeLister
}

func NewClusterLevelMetricsCollector(provider cloudprice.CloudProviderInterface, lister v1.NodeLister) *ClusterLevelMetricsCollector {
	metricsLabelKey := []string{
		values.RegionLabelKey,
		values.CloudProviderLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
	}
	clusterActiveGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.ClusterActiveMetricsName,
		Help: "The cluster alive metrics"}, metricsLabelKey)

	prometheus.MustRegister(clusterActiveGV)
	return &ClusterLevelMetricsCollector{
		clusterActiveGV: clusterActiveGV,
		nodeLister:      lister,
		provider:        provider,
	}
}

func (c *ClusterLevelMetricsCollector) StartCollectClusterLevelMetrics(ctx context.Context,
	interval time.Duration, agentOptions *options.AgentOptions) {
	ticker := time.NewTicker(interval)

	klog.Infof("Start collecting cluster level metrics")
	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			klog.Infof("Stop cluster level metrics collect")
			return
		case <-ticker.C:
			c.collectClusterMetrics(agentOptions)
		}
	}
}

func (c *ClusterLevelMetricsCollector) collectClusterMetrics(agentOptions *options.AgentOptions) {
	nodes, err := c.nodeLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all nodes error:%v", err)
		return
	}

	nodeCostInfo, err := c.provider.GetNodeHourlyPrice(nodes[0])
	if err != nil {
		klog.Errorf("Get node hourly price error:%v", err)
		return
	}
	labels := prometheus.Labels{
		values.RegionLabelKey:        nodeCostInfo.Region,
		values.CloudProviderLabelKey: agentOptions.CloudProvider,
		values.ClusterNameLabelKey:   agentOptions.ClusterName,
		values.ClusterIdLabelKey:     agentOptions.ClusterId,
	}
	c.clusterActiveGV.With(labels).Set(1)
}
