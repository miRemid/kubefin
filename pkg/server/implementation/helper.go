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

package implementation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/api"
	"github.com/kubefin/kubefin/pkg/utils"
)

func GetStartEndStepsTimeFromCtx(ctx *gin.Context, stepSecondsIfNone int64) (int64, int64, int64, error) {
	// the time format should be unix format
	startTime, endTime, err := GetStartEndTimeFromCtx(ctx)
	if err != nil {
		return 0, 0, 0, err
	}
	stepSecondsPara := ctx.Query(api.QueryStepSecondsPara)
	if stepSecondsPara == "" {
		return startTime, endTime, stepSecondsIfNone, nil
	}
	stepSeconds, err := strconv.ParseInt(stepSecondsPara, 0, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return startTime, endTime, stepSeconds, nil
}

func GetStartEndTimeFromCtx(ctx *gin.Context) (int64, int64, error) {
	// the time format should be unix format
	startTimeStr := ctx.Query(api.QueryStartTimePara)
	endTimeStr := ctx.Query(api.QueryEndTimePara)
	if startTimeStr == "" || endTimeStr == "" {
		startTimePara, endTimePara, err := utils.GetCurrentMonthFirstLastDay()
		if err != nil {
			klog.Errorf("Query current time error:%v", err)
			return 0, 0, err
		}

		return startTimePara, endTimePara, nil
	}

	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return startTime, endTime, nil
}
