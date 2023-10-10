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

package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/utils"
	"github.com/kubefin/kubefin/pkg/values"
)

type PodLevelMetricsCollector struct {
	metricsClient *versioned.Clientset
	provider      cloudprice.CloudProviderInterface

	podLister  v1.PodLister
	nodeLister v1.NodeLister

	// podResourceCostGV will be the node price * pod request resource
	podResourceCostGV *prometheus.GaugeVec

	podResourceRequestGV *prometheus.GaugeVec
	podResourceUsageGV   *prometheus.GaugeVec
}

func NewPodLevelMetricsCollector(client *versioned.Clientset, provider cloudprice.CloudProviderInterface,
	podLister v1.PodLister, nodeLister v1.NodeLister) *PodLevelMetricsCollector {
	containerNoneCareLabelKey := []string{
		values.NamespaceLabelKey,
		values.PodNameLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.ResourceTypeLabelKey,
		values.PodScheduledKey,
		values.LabelsLabelKey,
	}
	podResourceCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.PodResoueceCostMetricsName,
		Help: "The pod level resource cost"}, containerNoneCareLabelKey)

	containerCareLabelKey := []string{
		values.NamespaceLabelKey,
		values.PodNameLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.ResourceTypeLabelKey,
		values.LabelsLabelKey,
		// For multiple container pod, this metrics is needed for cpu/memory size recommendation
		values.ContainerNameLabelKey,
	}
	podResourceRequestGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.PodResourceRequestMetricsName,
		Help: "The pod container level resource requested"}, containerCareLabelKey)
	podResourceUsageGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.PodResourceUsageMetricsName,
		Help: "The pod container level resource usage"}, containerCareLabelKey)

	prometheus.MustRegister(podResourceRequestGV, podResourceUsageGV, podResourceCostGV)
	return &PodLevelMetricsCollector{
		metricsClient:        client,
		podLister:            podLister,
		nodeLister:           nodeLister,
		provider:             provider,
		podResourceRequestGV: podResourceRequestGV,
		podResourceUsageGV:   podResourceUsageGV,
		podResourceCostGV:    podResourceCostGV,
	}
}

func (p *PodLevelMetricsCollector) StartCollectPodLevelMetrics(ctx context.Context,
	interval time.Duration, agentOptions *options.AgentOptions) {
	ticker := time.NewTicker(interval)

	klog.Infof("Start collecting Pod level metrics")
	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			klog.Infof("Stop collecting Pod level metrics")
			return
		case <-ticker.C:
			p.collectPodResourceCost(agentOptions)
			p.collectPodResourceRequest(agentOptions)
			p.collectPodResourceUsage(ctx, agentOptions)
		}
	}
}

func (p *PodLevelMetricsCollector) collectPodResourceCost(agentOptions *options.AgentOptions) {
	pods, err := p.podLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}
	for _, pod := range pods {
		podLabels, err := json.Marshal(pod.Labels)
		if err != nil {
			klog.Errorf("Marshal pod labels error:%v", err)
			return
		}
		cost := 0.0
		scheduled := "false"
		if pod.Spec.NodeName != "" {
			cost = utils.ParsePodResourceCost(pod, p.provider, p.nodeLister)
			scheduled = "true"
		}
		labels := prometheus.Labels{
			values.NamespaceLabelKey:    pod.Namespace,
			values.PodNameLabelKey:      pod.Name,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(podLabels),
			values.PodScheduledKey:      scheduled,
			values.ResourceTypeLabelKey: "cost",
		}
		p.podResourceCostGV.With(labels).Set(cost)
	}
}

func (p *PodLevelMetricsCollector) collectPodResourceRequest(agentOptions *options.AgentOptions) {
	pods, err := p.podLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}
	for _, pod := range pods {
		podLabels, err := json.Marshal(pod.Labels)
		if err != nil {
			klog.Errorf("Marshal pod labels error:%v", err)
			return
		}
		labels := prometheus.Labels{
			values.NamespaceLabelKey:   pod.Namespace,
			values.PodNameLabelKey:     pod.Name,
			values.ClusterNameLabelKey: agentOptions.ClusterName,
			values.ClusterIdLabelKey:   agentOptions.ClusterId,
			values.LabelsLabelKey:      string(podLabels),
		}
		cpuRequest, memoryRequest := utils.ParsePodResourceRequest(pod.Spec.Containers)

		labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
		for containerName, cpu := range cpuRequest {
			labels[values.ContainerNameLabelKey] = containerName
			p.podResourceRequestGV.With(labels).Set(cpu)
		}
		labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
		for containerName, memory := range memoryRequest {
			labels[values.ContainerNameLabelKey] = containerName
			p.podResourceRequestGV.With(labels).Set(memory)
		}
	}
}

func (p *PodLevelMetricsCollector) collectPodResourceUsage(ctx context.Context, agentOptions *options.AgentOptions) {
	pods, err := p.metricsClient.MetricsV1beta1().PodMetricses(corev1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List all pod metrics error:%v, kubernetes metrics server may not be installed", err)
		return
	}

	for _, pod := range pods.Items {
		podStandard, err := p.podLister.Pods(pod.Namespace).Get(pod.Name)
		if err != nil {
			klog.Errorf("Get pod error:%v", err)
			return
		}
		podLabels, err := json.Marshal(podStandard.Labels)
		if err != nil {
			klog.Errorf("Marshal pod labels error:%v", err)
			return
		}
		labels := prometheus.Labels{
			values.NamespaceLabelKey:   pod.Namespace,
			values.PodNameLabelKey:     pod.Name,
			values.ClusterNameLabelKey: agentOptions.ClusterName,
			values.ClusterIdLabelKey:   agentOptions.ClusterId,
			values.LabelsLabelKey:      string(podLabels),
		}
		cpuUsage, memoryUsage := utils.ParsePodResourceUsage(pod.Containers)

		labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
		for containerName, cpu := range cpuUsage {
			labels[values.ContainerNameLabelKey] = containerName
			p.podResourceUsageGV.With(labels).Set(cpu)
		}
		labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
		for containerName, memory := range memoryUsage {
			labels[values.ContainerNameLabelKey] = containerName
			p.podResourceUsageGV.With(labels).Set(memory)
		}
	}
}
