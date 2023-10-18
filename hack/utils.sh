#!/usr/bin/env bash

# Copyright 2023 The KubeFin Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

function get_mimir_server_ip() {
  kubeconfig=$1
  server_node_ip=$(kubectl get node --kubeconfig="${kubeconfig}" -ojson | jq .items[0].status.addresses[0].address | sed 's/"//g')
  echo "${server_node_ip}"
}

function get_mimir_server_port() {
  kubeconfig=$1
  server_node_port=$(kubectl get svc -nkubefin mimir --kubeconfig="${kubeconfig}" -ojson | jq .spec.ports[0].nodePort)
  echo "${server_node_port}"
}

function init_primary_config() {
  cloud_provider=$1
  cluster_name=$2
  multi_cluster_enable=$3
  dashboard_tag=$4

  mimir_service_type="ClusterIP"
  if [[ "${multi_cluster_enable}" == "true" ]]; then
    mimir_service_type="LoadBalancer"
  fi

  rm -rf config_primary
  cp -r config_template config_primary

  sed -i'' -e "s/{REMOTE_WRITE_ADDRESS}/http:\/\/mimir.kubefin.svc.cluster.local:9009\/api\/v1\/push/g" config_primary/core/configmap/otel.yaml
  sed -i'' -e "s/{KUBEFIN_AGENT_IAMGE}/ko:\/\/github.com\/kubefin\/kubefin\/cmd\/kubefin-agent/g" config_primary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLOUD_PROVIDER}/${cloud_provider}/g" config_primary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLUSTER_ID}//g" config_primary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLUSTER_NAME}/${cluster_name}/g" config_primary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{MIMIR_SERVICE_TYPE}/${mimir_service_type}/g" config_primary/third_party/mimir.yaml
  sed -i'' -e "s/{KUBEFIN_DASHBOARD_VERSION}/${dashboard_tag}/" config_primary/core/deployments/kubefin-cost-analyzer.yaml
}

function init_secondary_config() {
  cloud_provider=$1
  cluster_name=$2
  metrics_push_addr=$3

  rm -rf config_secondary
  cp -r config_template config_secondary

  rm -rf config_secondary/third_party
  rm -rf config_secondary/core/deployments/kubefin-cost-analyzer.yaml

  sed -i'' -e "s/{REMOTE_WRITE_ADDRESS}/${metrics_push_addr}/g" config_secondary/core/configmap/otel.yaml
  sed -i'' -e "s/{KUBEFIN_AGENT_IAMGE}/ko:\/\/github.com\/kubefin\/kubefin\/cmd\/kubefin-agent/g" config_secondary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLOUD_PROVIDER}/${cloud_provider}/g" config_secondary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLUSTER_ID}//g" config_secondary/core/deployments/kubefin-agent.yaml
  sed -i'' -e "s/{CLUSTER_NAME}/${cluster_name}/g" config_secondary/core/deployments/kubefin-agent.yaml
}

function get_build_arch() {
  platform=$(uname -m)

  if [[ "$platform" == "aarch64" || "$platform" == "arm64" ]]; then
    echo "arm64"
  else
    echo "amd64"
  fi
}

function echo_info() {
  echo -e "\033[34m[INFO] $1\033[0m"
}

function echo_note() {
  echo -e "\033[33m$1\033[0m"
}
