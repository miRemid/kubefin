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

package eks

import (
	"strconv"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	cloudpriceapis "github.com/kubefin/kubefin/pkg/cloudprice/apis"
)

type EksCloudProvider struct {
	client kubernetes.Interface

	cpuMemoryCostRatio float64

	// nodePriceMap maps [region name][node type]price
	nodePriceMap  map[string]map[string]float64
	nodePriceLock sync.Mutex

	// nodeSpecMap maps [node type]NodeSpec
	nodeSpecMap  map[string]cloudpriceapis.NodeSpec
	nodeSpecLock sync.Mutex
}

func NewEksCloudProvider(client kubernetes.Interface, agentOptions *options.AgentOptions) (*EksCloudProvider, error) {
	var err error

	cpuMemoryCostRatio := cloudpriceapis.DefaultCPUMemoryCostRatio
	if agentOptions.CPUMemoryCostRatio != "" {
		cpuMemoryCostRatio, err = strconv.ParseFloat(agentOptions.CPUMemoryCostRatio, 64)
		if err != nil {
			return nil, err
		}
	}
	eksCloud := EksCloudProvider{
		client:             client,
		cpuMemoryCostRatio: cpuMemoryCostRatio,
		nodePriceMap:       map[string]map[string]float64{},
		nodeSpecMap:        map[string]cloudpriceapis.NodeSpec{},
	}

	return &eksCloud, nil
}

func (e *EksCloudProvider) ParseClusterInfo(agentOptions *options.AgentOptions) error {
	return nil
}

func (e *EksCloudProvider) GetNodeHourlyPrice(node *v1.Node) (*api.InstancePriceInfo, error) {
	return nil, nil
}
