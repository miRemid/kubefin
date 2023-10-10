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

package main

import (
	"os"

	pkgserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/kubefin/kubefin/cmd/kubefin-cost-analyzer/app"
)

//	@title			KubeFin API
//	@version		0.1
//	@description	KubeFin Cost Analyzer API.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/api/v1
func main() {
	ctx := pkgserver.SetupSignalContext()
	cmd := app.NewAnalyzerCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
