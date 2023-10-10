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

package cloudprice

import (
	"github.com/kubefin/kubefin/cmd/kubefin-agent/app/options"
	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/cloudprice/ack"
	"github.com/kubefin/kubefin/pkg/cloudprice/defaultcloud"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type CloudProviderInterface interface {
	GetNodeHourlyPrice(node *v1.Node) (*api.InstancePriceInfo, error)
	ParseClusterInfo(agentOptions *options.AgentOptions) error
}

func NewCloudProvider(client kubernetes.Interface, agentOptions *options.AgentOptions) (CloudProviderInterface, error) {
	switch agentOptions.CloudProvider {
	case api.CloudProviderAck:
		return ack.NewAckCloudProvider(client, agentOptions)
	default:
		return defaultcloud.NewDefaultCloudProvider(client, agentOptions)
	}
}
