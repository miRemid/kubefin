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

package costs_handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/server/implementation"
	"github.com/kubefin/kubefin/pkg/utils"
	"github.com/kubefin/kubefin/pkg/values"
)

// ClusterWorkloadsCostsHandler  godoc
//
//	@Summary		Get specific cluster workloads costs
//	@Description	Get specific cluster workloads costs with time range
//	@Tags			Costs
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster Id"
//	@Param			startTime	query		uint64	false	"The start time to query"
//	@Param			endTime		query		uint64	false	"The end time to query"
//	@Param			stepSeconds	query		uint64	false	"The step seconds of the data to return"
//	@Param			aggregateBy	query		string	false	"The aggregated way to show workload costs"
//	@Success		200			{object}	api.ClusterWorkloadCostList
//	@Failure		500			{object}	api.StatusError
//	@Router			/costs/clusters/{cluster_id}/workload [get]
func ClusterWorkloadsCostsHandler(ctx *gin.Context) {
	klog.V(6).Info("Start to query clusters workload cost")
	tenantId := utils.ParserTenantIdFromCtx(ctx)
	clusterId := utils.ParseClusterFromCtx(ctx)
	if clusterId == "" {
		utils.ForwardStatusError(ctx, http.StatusBadRequest,
			api.QueryParaErrorStatus, api.QueryParaErrorReason, "")
		return
	}
	aggregateBy := ctx.Query(api.QueryAggregateBy)
	if aggregateBy == "" {
		aggregateBy = api.AggregateByAll
	}
	startTime, endTime, stepSeconds, err := implementation.GetStartEndStepsTimeFromCtx(ctx, values.DefaultStepSeconds)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusBadRequest,
			api.QueryParaErrorStatus, api.QueryParaErrorReason, err.Error())
		return
	}

	workloadCost, err := implementation.QueryWorkloadCostsWithTimeRange(tenantId, clusterId, startTime, endTime, stepSeconds, aggregateBy)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	bodyBytes, err := json.Marshal(workloadCost)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	ctx.Data(http.StatusOK, "application/json", bodyBytes)
}
