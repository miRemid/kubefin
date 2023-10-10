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

	"github.com/spf13/pflag"

	"github.com/kubefin/kubefin/pkg/values"
)

type AnalyzerOptions struct {
	QueryBackendEndpoint string
}

// NewAnalyzerOptions builds an empty options.
func NewAnalyzerOptions() *AnalyzerOptions {
	return &AnalyzerOptions{
		QueryBackendEndpoint: os.Getenv(values.QueryBackendEndpointEnv),
	}
}

func (o *AnalyzerOptions) Complete() error {
	return nil
}

func (o *AnalyzerOptions) Validate() error {
	return nil
}

func (o *AnalyzerOptions) ApplyTo() {
}

func (o *AnalyzerOptions) AddFlags(flags *pflag.FlagSet) {
}
