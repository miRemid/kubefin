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
)

// ClustersCostsSummaryHandler   godoc
//
//	@Summary		Get all clusters costs summary
//	@Description	Get all clusters costs summary in current two month
//	@Tags			Costs
//	@Produce		json
//	@Success		200	{object}	api.ClusterCostsSummaryList
//	@Failure		500	{object}	api.StatusError
//	@Router			/costs/summary   [get]
func ClustersCostsSummaryHandler(ctx *gin.Context) {
	klog.V(6).Info("Start to query clusters costs summary")
	tenantId := utils.ParserTenantIdFromCtx(ctx)
	// If data not comes up in two-month period, we will ignore it
	start, end := utils.GetCurrentTwoMonthStartEndTime()
	allClustersProperty, err := implementation.QueryAllClustersBasicProperty(tenantId, start, end)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	allClustersSummary, err := implementation.QueryAllClustersCurrentMonthCost(tenantId)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	summaries := implementation.ConvertToMultiClustersCostsList(allClustersSummary, allClustersProperty)
	bodyBytes, err := json.Marshal(summaries)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	ctx.Data(http.StatusOK, "application/json", bodyBytes)
}

// ClusterCostsSummaryHandler    godoc
//
//	@Summary		Get specific cluster costs summary
//	@Description	Get specific cluster costs summary in current two month
//	@Tags			Costs
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster Id"
//	@Success		200			{object}	api.ClusterCostsSummary
//	@Failure		500			{object}	api.StatusError
//	@Router			/costs/clusters/{cluster_id}/summary [get]
func ClusterCostsSummaryHandler(ctx *gin.Context) {
	klog.V(4).Info("Start to query specific cluster costs summary")
	tenantId := utils.ParserTenantIdFromCtx(ctx)
	clusterId := utils.ParseClusterFromCtx(ctx)
	// If data not comes up in two-month period, we will ignore it
	start, end := utils.GetCurrentTwoMonthStartEndTime()
	clusterProperty, err := implementation.QueryClusterBasicProperty(tenantId, clusterId, start, end)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusNotFound,
			api.QueryNotFoundStatus, api.QueryNotFoundReason, "no clusters found")
		return
	}
	summary, err := implementation.QueryClusterCurrentMonthCost(tenantId, clusterId)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	summary.ClusterBasicProperty = *clusterProperty
	klog.Infof("%v", summary)
	bodyBytes, err := json.Marshal(summary)
	if err != nil {
		utils.ForwardStatusError(ctx, http.StatusInternalServerError,
			api.QueryFailedStatus, api.QueryFailedReason, err.Error())
		return
	}
	ctx.Data(http.StatusOK, "application/json", bodyBytes)
}
