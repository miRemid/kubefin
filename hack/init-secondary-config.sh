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

set -o errexit
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/../../
EXE_ROOT=$(dirname "${BASH_SOURCE[0]}")
source "${EXE_ROOT}"/utils.sh

function usage() {
  echo "This script will deploy KubeFin components on primary cluster"
  echo "Usage: hack/install-kubefin-primary.sh <CLUSTER_PROVIDER> <CLUSTER_NAME> <METRICS_WRITE_ADDRESS>"
  echo "Example: hack/install-kubefin-primary.sh default kubefin-server"
}

if [[ $# -ne 3 ]]; then
  usage
  exit 1
fi

# shellcheck disable=SC2001
metrics_push_address=$(echo "$3" | sed 's/\//\\\//g')
init_secondary_config "$1" "$2" "${metrics_push_address}"
