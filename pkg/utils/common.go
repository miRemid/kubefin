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

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/values"
)

func ForwardStatusError(ctx *gin.Context, httpCode int, status, reason, message string) {
	apiError := newStatusError(status, reason, message, httpCode)
	forwardRaw, err := json.Marshal(apiError)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}
	ctx.Data(httpCode, "application/json", forwardRaw)
}

func newStatusError(status, reason, message string, httpCode int) *api.StatusError {
	return &api.StatusError{
		APIVersion: api.KubeFinAPIVersion,
		Kind:       api.KubeFinStatusKind,
		Status:     status,
		Message:    message,
		Reason:     reason,
		Code:       httpCode,
	}
}

func GetPodNamespaceNameKey(labels model.Metric) string {
	return fmt.Sprintf("%s/%s",
		labels[model.LabelName(values.NamespaceLabelKey)],
		labels[model.LabelName(values.LabelsLabelKey)])
}

func GetWorkloadNamespaceNameKey(labels model.Metric) string {
	return fmt.Sprintf("%s/%s",
		labels[model.LabelName(values.NamespaceLabelKey)],
		labels[model.LabelName(values.WorkloadNameLabelKey)])
}

func GetNamespace(labels model.Metric) string {
	return string(labels[model.LabelName(values.NamespaceLabelKey)])
}

func GetPodName(labels model.Metric) string {
	return string(labels[model.LabelName(values.LabelsLabelKey)])
}

func GetWorkloadName(labels model.Metric) string {
	return string(labels[model.LabelName(values.WorkloadNameLabelKey)])
}

func GetNodeName(labels model.Metric) string {
	return string(labels[model.LabelName(values.NodeNameLabelKey)])
}

func ParserTenantIdFromCtx(ctx *gin.Context) string {
	// Logic based on: https://grafana.com/docs/mimir/latest/operators-guide/secure/authentication-and-authorization/
	tenantId := ctx.GetHeader(values.MultiTenantHeader)
	return tenantId
}

func ParseClusterFromCtx(ctx *gin.Context) string {
	return ctx.Param(values.ClusterIdQueryParameter)
}

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func GetCurrentMonthFirstLastDay() (int64, int64, error) {
	nowUnixTime := GetCurrentTime()
	now := time.Unix(nowUnixTime, 0)
	firstDay := now.AddDate(0, 0, -now.Day()+1)
	lastDay := now.AddDate(0, 0, 1)
	startTime := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 0, 0, 0, 0, now.Location())
	endTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 0, 0, 0, 0, now.Location())
	return startTime.Unix(), endTime.Unix(), nil
}

func GetCurrentTwoMonthStartEndTime() (start, end int64) {
	return time.Now().AddDate(0, -2, 0).Unix(), time.Now().Unix()
}
