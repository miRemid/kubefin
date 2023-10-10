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

package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/cloudprice"
	"github.com/kubefin/kubefin/pkg/metrics"
)

// NewAgentCommand creates a *cobra.Command object with defaultcloud parameters
func NewAgentCommand(ctx context.Context) *cobra.Command {
	opts := options.NewAgentOptions()

	cmd := &cobra.Command{
		Use:  "kubefin-agent",
		Long: `kubefin-agent used to scrap metrics to storage store such as thanos`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Complete(); err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}

			if err := Run(ctx, opts); err != nil {
				return err
			}

			return nil
		},
	}

	fss := cliflag.NamedFlagSets{}

	agentFlagSet := fss.FlagSet("agent")
	opts.AddFlags(agentFlagSet)

	logFlagSet := fss.FlagSet("log")
	klog.InitFlags(flag.CommandLine)
	logFlagSet.AddGoFlagSet(flag.CommandLine)

	cmd.Flags().AddFlagSet(agentFlagSet)
	cmd.Flags().AddFlagSet(logFlagSet)

	return cmd
}

func Run(ctx context.Context, opts *options.AgentOptions) error {
	klog.Infof("Start kubefin-agent...")
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("find values for connect kube-apiserver error:%v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create client to connect kube-apiserver error:%v", err)
	}

	metricsClientSet, err := metricsv.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create metrics client to connect kube-apiserver error:%v", err)
	}

	factory := informers.NewSharedInformerFactory(clientSet, 0)
	coreResourceInformerLister := getAllCoreResourceLister(factory)

	stopCh := ctx.Done()
	factory.Start(stopCh)

	klog.Infof("Wait node cache sync...")
	if ok := cache.WaitForCacheSync(stopCh,
		coreResourceInformerLister.NamespaceInformer.HasSynced,
		coreResourceInformerLister.NodeInformer.HasSynced,
		coreResourceInformerLister.PodInformer.HasSynced); !ok {
		return fmt.Errorf("wait core resource cache sync failed")
	}

	provider, err := cloudprice.NewCloudProvider(clientSet, opts)
	if err != nil {
		return fmt.Errorf("create cloud provider error:%v", err)
	}
	if err := provider.ParseClusterInfo(opts); err != nil {
		return err
	}

	metricsCollector := metrics.NewAgentMetricsCollector(ctx, opts, coreResourceInformerLister, provider, metricsClientSet)
	runFunc := func(runCtx context.Context) {
		metricsCollector.StartAgentMetricsCollector()
	}

	klog.Infof("Start metrics http server")
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			klog.Fatalf("Start http server error:%v", err)
		}
	}()

	if err := runLeaderElection(ctx, clientSet, opts, runFunc); err != nil {
		return fmt.Errorf("run leader election error:%v", err)
	}

	<-stopCh
	return nil
}

func runLeaderElection(ctx context.Context, clientset kubernetes.Interface,
	opts *options.AgentOptions, runFunc func(ctx context.Context)) error {
	rl, err := resourcelock.New(opts.LeaderElection.ResourceLock,
		opts.LeaderElection.ResourceNamespace,
		opts.LeaderElection.ResourceName,
		clientset.CoreV1(),
		clientset.CoordinationV1(),
		resourcelock.ResourceLockConfig{Identity: opts.LeaderElectionID})
	if err != nil {
		return fmt.Errorf("couldn't create resource lock: %v", err)
	}
	leaderElectionCfg := leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: opts.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: opts.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   opts.LeaderElection.RetryPeriod.Duration,
		WatchDog:      leaderelection.NewLeaderHealthzAdaptor(time.Second * 20),
		Name:          opts.LeaderElection.ResourceName,
	}
	leaderElectionCfg.Callbacks = leaderelection.LeaderCallbacks{
		OnStartedLeading: runFunc,
		OnStoppedLeading: func() {
			klog.Infof("Leader election lost")
		},
	}
	leaderElector, err := leaderelection.NewLeaderElector(leaderElectionCfg)
	if err != nil {
		return fmt.Errorf("couldn't create leader elector:%v", err)
	}
	leaderElector.Run(ctx)
	return nil
}

func getAllCoreResourceLister(factory informers.SharedInformerFactory) *api.CoreResourceInformerLister {
	coreResource := factory.Core().V1()
	appsResource := factory.Apps().V1()
	return &api.CoreResourceInformerLister{
		NodeInformer:        coreResource.Nodes().Informer(),
		NamespaceInformer:   coreResource.Namespaces().Informer(),
		PodInformer:         coreResource.Pods().Informer(),
		DeploymentInformer:  appsResource.Deployments().Informer(),
		StatefulSetInformer: appsResource.StatefulSets().Informer(),
		DaemonSetInformer:   appsResource.DaemonSets().Informer(),
		NodeLister:          coreResource.Nodes().Lister(),
		PodLister:           coreResource.Pods().Lister(),
		DeploymentLister:    appsResource.Deployments().Lister(),
		StatefulSetLister:   appsResource.StatefulSets().Lister(),
		DaemonSetLister:     appsResource.DaemonSets().Lister(),
	}
}
