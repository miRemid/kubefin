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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	listersappv1 "k8s.io/client-go/listers/apps/v1"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/utils"
	"github.com/kubefin/kubefin/pkg/values"
)

// WorkloadLevelMetricsCollector collects metrics about deployment/statefulset/daemonset
type WorkloadLevelMetricsCollector struct {
	metricsClient *versioned.Clientset
	provider      cloudprice.CloudProviderInterface

	podLister         listercorev1.PodLister
	nodeLister        listercorev1.NodeLister
	daemonSetLister   listersappv1.DaemonSetLister
	deploymentLister  listersappv1.DeploymentLister
	statefulSetLister listersappv1.StatefulSetLister

	workloadResourceCostGV    *prometheus.GaugeVec
	workloadPodCountGV        *prometheus.GaugeVec
	workloadResourceRequestGV *prometheus.GaugeVec
	workloadResourceUsageGV   *prometheus.GaugeVec
}

func NewWorkloadLevelMetricsCollector(client *versioned.Clientset, provider cloudprice.CloudProviderInterface,
	podLister listercorev1.PodLister, nodeLister listercorev1.NodeLister, daemonSetLister listersappv1.DaemonSetLister,
	deploymentLister listersappv1.DeploymentLister, statefulSetLister listersappv1.StatefulSetLister) *WorkloadLevelMetricsCollector {
	containerNoneCareLabelKey := []string{
		values.WorkloadTypeLabelKey,
		values.WorkloadNameLabelKey,
		values.NamespaceLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.LabelsLabelKey,
		values.ResourceTypeLabelKey,
	}
	workloadResourceCostGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.WorkloadResourceCostMetricsName,
		Help: "The workload resource cost"}, containerNoneCareLabelKey)
	workloadPodCountGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.WorkloadPodCountMetricsName,
		Help: "The workload pod count"}, containerNoneCareLabelKey)

	containerCareLabelKey := []string{
		values.WorkloadTypeLabelKey,
		values.WorkloadNameLabelKey,
		values.NamespaceLabelKey,
		values.ClusterNameLabelKey,
		values.ClusterIdLabelKey,
		values.LabelsLabelKey,
		values.ResourceTypeLabelKey,
		// For multiple container workload, this metrics is needed for cpu/memory size recommendation
		values.ContainerNameLabelKey,
	}
	workloadResourceRequestGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.WorkloadResourceRequestMetricsName,
		Help: "The workload resource request",
	}, containerCareLabelKey)
	workloadResourceUsageGV := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: values.WorkloadResourceUsageMetricsName,
		Help: "The workload resource usage",
	}, containerCareLabelKey)

	prometheus.MustRegister(workloadResourceCostGV, workloadResourceRequestGV, workloadResourceUsageGV, workloadPodCountGV)
	return &WorkloadLevelMetricsCollector{
		metricsClient:             client,
		provider:                  provider,
		podLister:                 podLister,
		nodeLister:                nodeLister,
		daemonSetLister:           daemonSetLister,
		deploymentLister:          deploymentLister,
		statefulSetLister:         statefulSetLister,
		workloadResourceCostGV:    workloadResourceCostGV,
		workloadResourceRequestGV: workloadResourceRequestGV,
		workloadResourceUsageGV:   workloadResourceUsageGV,
		workloadPodCountGV:        workloadPodCountGV,
	}
}

func (w *WorkloadLevelMetricsCollector) StartCollectWorkloadLevelMetrics(ctx context.Context,
	interval time.Duration, agentOptions *options.AgentOptions) {
	ticker := time.NewTicker(interval)

	klog.Infof("Start collecting workload level metrics")
	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			klog.Infof("Stop collecting DaemonSet level metrics")
			return
		case <-ticker.C:
			w.collectDaemonSetResourceMetrics(agentOptions)
			w.collectStatefulSetResourceMetrics(agentOptions)
			w.collectDeploymentResourceMetrics(agentOptions)
		}
	}
}

func (w *WorkloadLevelMetricsCollector) collectDaemonSetResourceMetrics(agentOptions *options.AgentOptions) {
	daemonSets, err := w.daemonSetLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all daemonSets error:%v", err)
		return
	}
	for _, daemonSet := range daemonSets {
		daemonSetLabels, err := json.Marshal(daemonSet.Labels)
		if err != nil {
			klog.Errorf("Marshal daemonSet labels error:%v", err)
			continue
		}

		matchLabels := daemonSet.Spec.Selector.MatchLabels
		selector := labels.Set(matchLabels).AsSelector()

		containerNoneCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "daemonset",
			values.WorkloadNameLabelKey: daemonSet.Name,
			values.NamespaceLabelKey:    daemonSet.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(daemonSetLabels),
		}
		w.collectDaemonSetPodMetrics(daemonSet, agentOptions, selector, containerNoneCareLabelValues)
		w.collectDaemonSetCostMetrics(daemonSet, agentOptions, selector, containerNoneCareLabelValues)

		containerCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "daemonset",
			values.WorkloadNameLabelKey: daemonSet.Name,
			values.NamespaceLabelKey:    daemonSet.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(daemonSetLabels),
		}
		w.collectDaemonSetRequestMetrics(daemonSet, agentOptions, selector, containerCareLabelValues)
		w.collectDaemonSetUsageMetrics(daemonSet, agentOptions, selector, containerCareLabelValues)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDaemonSetPodMetrics(ds *appsv1.DaemonSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(ds.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}
	podCount := float64(len(pods))
	labels[values.ResourceTypeLabelKey] = "pod"
	w.workloadPodCountGV.With(labels).Set(podCount)
}

func (w *WorkloadLevelMetricsCollector) collectDaemonSetRequestMetrics(ds *appsv1.DaemonSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(ds.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	cpuTotalRequest, ramTotalRequest := map[string]float64{}, map[string]float64{}
	for _, pod := range pods {
		cpu, ram := utils.ParsePodResourceRequest(pod.Spec.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalRequest[containerName]; !ok {
				cpuTotalRequest[containerName] = 0
			}
			cpuTotalRequest[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := ramTotalRequest[containerName]; !ok {
				ramTotalRequest[containerName] = 0
			}
			ramTotalRequest[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(cpu)
	}
	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range ramTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDaemonSetUsageMetrics(ds *appsv1.DaemonSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.metricsClient.MetricsV1beta1().PodMetricses(ds.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List all pod metrics error:%v, kubernetes metrics server may not be installed", err)
		return
	}
	cpuTotalUsage, memoryTotalUsage := map[string]float64{}, map[string]float64{}
	for _, pod := range pods.Items {
		cpu, ram := utils.ParsePodResourceUsage(pod.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalUsage[containerName]; !ok {
				cpuTotalUsage[containerName] = 0
			}
			cpuTotalUsage[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := memoryTotalUsage[containerName]; !ok {
				memoryTotalUsage[containerName] = 0
			}
			memoryTotalUsage[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(cpu)
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range memoryTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDaemonSetCostMetrics(ds *appsv1.DaemonSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(ds.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	totalCost := 0.0
	for _, pod := range pods {
		if pod.Spec.NodeName == "" {
			continue
		}
		totalCost += utils.ParsePodResourceCost(pod, w.provider, w.nodeLister)
	}

	labels[values.ResourceTypeLabelKey] = "cost"
	w.workloadResourceCostGV.With(labels).Set(totalCost)
}

func (w *WorkloadLevelMetricsCollector) collectStatefulSetResourceMetrics(agentOptions *options.AgentOptions) {
	statefulSets, err := w.statefulSetLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all statefulSet error:%v", err)
		return
	}
	for _, statefulSet := range statefulSets {
		statefulSetLabels, err := json.Marshal(statefulSet.Labels)
		if err != nil {
			klog.Errorf("Marshal daemonSet labels error:%v", err)
			continue
		}

		matchLabels := statefulSet.Spec.Selector.MatchLabels
		selector := labels.Set(matchLabels).AsSelector()

		containerNoneCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "statefulset",
			values.WorkloadNameLabelKey: statefulSet.Name,
			values.NamespaceLabelKey:    statefulSet.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(statefulSetLabels),
		}
		w.collectStatefulSetPodMetrics(statefulSet, agentOptions, selector, containerNoneCareLabelValues)
		w.collectStatefulSetCostMetrics(statefulSet, agentOptions, selector, containerNoneCareLabelValues)

		containerCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "statefulset",
			values.WorkloadNameLabelKey: statefulSet.Name,
			values.NamespaceLabelKey:    statefulSet.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(statefulSetLabels),
		}
		w.collectStatefulSetRequestMetrics(statefulSet, agentOptions, selector, containerCareLabelValues)
		w.collectStatefulSetUsageMetrics(statefulSet, agentOptions, selector, containerCareLabelValues)
	}
}

func (w *WorkloadLevelMetricsCollector) collectStatefulSetPodMetrics(sf *appsv1.StatefulSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(sf.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}
	podCount := float64(len(pods))
	labels[values.ResourceTypeLabelKey] = "pod"
	w.workloadPodCountGV.With(labels).Set(podCount)
}

func (w *WorkloadLevelMetricsCollector) collectStatefulSetRequestMetrics(sf *appsv1.StatefulSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(sf.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	cpuTotalRequest, ramTotalRequest := map[string]float64{}, map[string]float64{}
	for _, pod := range pods {
		cpu, ram := utils.ParsePodResourceRequest(pod.Spec.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalRequest[containerName]; !ok {
				cpuTotalRequest[containerName] = 0
			}
			cpuTotalRequest[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := ramTotalRequest[containerName]; !ok {
				ramTotalRequest[containerName] = 0
			}
			ramTotalRequest[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(cpu)
	}
	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range ramTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectStatefulSetUsageMetrics(sf *appsv1.StatefulSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.metricsClient.MetricsV1beta1().PodMetricses(sf.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List all pod metrics error:%v, kubernetes metrics server may not be installed", err)
		return
	}
	cpuTotalUsage, memoryTotalUsage := map[string]float64{}, map[string]float64{}
	for _, pod := range pods.Items {
		cpu, ram := utils.ParsePodResourceUsage(pod.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalUsage[containerName]; !ok {
				cpuTotalUsage[containerName] = 0
			}
			cpuTotalUsage[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := memoryTotalUsage[containerName]; !ok {
				memoryTotalUsage[containerName] = 0
			}
			memoryTotalUsage[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(cpu)
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range memoryTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectStatefulSetCostMetrics(sf *appsv1.StatefulSet, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(sf.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	totalCost := 0.0
	for _, pod := range pods {
		totalCost += utils.ParsePodResourceCost(pod, w.provider, w.nodeLister)
	}

	labels[values.ResourceTypeLabelKey] = "cost"
	w.workloadResourceCostGV.With(labels).Set(totalCost)
}

func (w *WorkloadLevelMetricsCollector) collectDeploymentResourceMetrics(agentOptions *options.AgentOptions) {
	deployments, err := w.deploymentLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("List all deployment error:%v", err)
		return
	}
	for _, deployment := range deployments {
		deploymentLabels, err := json.Marshal(deployment.Labels)
		if err != nil {
			klog.Errorf("Marshal daemonSet labels error:%v", err)
			continue
		}
		matchLabels := deployment.Spec.Selector.MatchLabels
		selector := labels.Set(matchLabels).AsSelector()

		containerNoneCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "deployment",
			values.WorkloadNameLabelKey: deployment.Name,
			values.NamespaceLabelKey:    deployment.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(deploymentLabels),
		}
		w.collectDeploymentPodMetrics(deployment, agentOptions, selector, containerNoneCareLabelValues)
		w.collectDeploymentCostMetrics(deployment, agentOptions, selector, containerNoneCareLabelValues)

		containerCareLabelValues := prometheus.Labels{
			values.WorkloadTypeLabelKey: "deployment",
			values.WorkloadNameLabelKey: deployment.Name,
			values.NamespaceLabelKey:    deployment.Namespace,
			values.ClusterNameLabelKey:  agentOptions.ClusterName,
			values.ClusterIdLabelKey:    agentOptions.ClusterId,
			values.LabelsLabelKey:       string(deploymentLabels),
		}
		w.collectDeploymentRequestMetrics(deployment, agentOptions, selector, containerCareLabelValues)
		w.collectDeploymentUsageMetrics(deployment, agentOptions, selector, containerCareLabelValues)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDeploymentPodMetrics(dm *appsv1.Deployment, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(dm.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}
	podCount := float64(len(pods))
	labels[values.ResourceTypeLabelKey] = "pod"
	w.workloadPodCountGV.With(labels).Set(podCount)
}

func (w *WorkloadLevelMetricsCollector) collectDeploymentRequestMetrics(dm *appsv1.Deployment, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(dm.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	cpuTotalRequest, ramTotalRequest := map[string]float64{}, map[string]float64{}
	for _, pod := range pods {
		cpu, ram := utils.ParsePodResourceRequest(pod.Spec.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalRequest[containerName]; !ok {
				cpuTotalRequest[containerName] = 0
			}
			cpuTotalRequest[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := ramTotalRequest[containerName]; !ok {
				ramTotalRequest[containerName] = 0
			}
			ramTotalRequest[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(cpu)
	}
	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range ramTotalRequest {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceRequestGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDeploymentUsageMetrics(dm *appsv1.Deployment, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.metricsClient.MetricsV1beta1().PodMetricses(dm.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("List all pod metrics error:%v, kubernetes metrics server may not be installed", err)
		return
	}
	cpuTotalUsage, memoryTotalUsage := map[string]float64{}, map[string]float64{}
	for _, pod := range pods.Items {
		cpu, ram := utils.ParsePodResourceUsage(pod.Containers)
		for containerName, value := range cpu {
			if _, ok := cpuTotalUsage[containerName]; !ok {
				cpuTotalUsage[containerName] = 0
			}
			cpuTotalUsage[containerName] += value
		}
		for containerName, value := range ram {
			if _, ok := memoryTotalUsage[containerName]; !ok {
				memoryTotalUsage[containerName] = 0
			}
			memoryTotalUsage[containerName] += value
		}
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceCPU)
	for containerName, cpu := range cpuTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(cpu)
	}

	labels[values.ResourceTypeLabelKey] = string(corev1.ResourceMemory)
	for containerName, ram := range memoryTotalUsage {
		labels[values.ContainerNameLabelKey] = containerName
		w.workloadResourceUsageGV.With(labels).Set(ram)
	}
}

func (w *WorkloadLevelMetricsCollector) collectDeploymentCostMetrics(dm *appsv1.Deployment, agentOptions *options.AgentOptions, selector labels.Selector, labels prometheus.Labels) {
	pods, err := w.podLister.Pods(dm.Namespace).List(selector)
	if err != nil {
		klog.Errorf("List all pods error:%v", err)
		return
	}

	totalCost := 0.0
	for _, pod := range pods {
		totalCost += utils.ParsePodResourceCost(pod, w.provider, w.nodeLister)
	}

	labels[values.ResourceTypeLabelKey] = "cost"
	w.workloadResourceCostGV.With(labels).Set(totalCost)
}
