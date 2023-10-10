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

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/cmd/kubefin-cost-analyzer/app/options"
	"github.com/kubefin/kubefin/pkg/query"
	pkgrouter "github.com/kubefin/kubefin/pkg/router"
)

// NewAnalyzerCommand creates a *cobra.Command object with parameters
func NewAnalyzerCommand(ctx context.Context) *cobra.Command {
	opts := options.NewAnalyzerOptions()

	cmd := &cobra.Command{
		Use:  "kubefin-cost-analyzer",
		Long: `kubefin-cost-analyzer used to do cost analyzer with prometheus sql`,
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

	analyzerFlagSet := fss.FlagSet("analyzer")
	opts.AddFlags(analyzerFlagSet)

	logFlagSet := fss.FlagSet("log")
	klog.InitFlags(flag.CommandLine)
	logFlagSet.AddGoFlagSet(flag.CommandLine)

	cmd.Flags().AddFlagSet(analyzerFlagSet)
	cmd.Flags().AddFlagSet(logFlagSet)

	return cmd
}

func Run(ctx context.Context, opts *options.AnalyzerOptions) error {
	klog.Infof("Start kubefin-cost-analyzer...")
	stopCh := ctx.Done()

	query.InitPromQueryClient(opts.QueryBackendEndpoint)

	router := pkgrouter.NewServerRouter()
	if err := router.Run(":8080"); err != nil {
		klog.Errorf("run server failed:%v", err)
	}

	<-stopCh
	return nil
}
