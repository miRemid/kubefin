// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React, { useState, useEffect } from "react";
import { Header, DateRangePicker } from "@cloudscape-design/components";
import { useParams } from "react-router";
import { ClusterCostInfo } from "../model/cluster-cost";
import { ClusterComputeCost } from "../model/cluster-compute-cost";
import apiClient from "../../../../common/network/http-common";
import moment from "moment";

export function ClusterCostHeader(props) {
  const { clusterId } = useParams();
  const [onChangeDetails, setOnChangeDetails] = useState({
    type: "relative",
    amount: 7,
    unit: "hour",
  });
  const [clusterCostInfo, setClusterCostInfo] = useState({
    clusterName: "-",
    clusterCostCurrent: 0,
    clusterMonthEstimateCost: 0,
    clusterAvgDailyCost: 0,
    clusterAvgHourlyCoreCost: 0,
    clusterComputeCostArray: [new ClusterComputeCost()],
  });

  const clusterName =
    clusterCostInfo === null ? "Loading" : clusterCostInfo.clusterName;

  const fetchClusterData = async (startTime, endTime, stepSeconds) => {
    const params = {
      startTime: startTime,
      endTime: endTime,
      stepSeconds: stepSeconds,
    };

    // TODO: Do fetch separately
    await Promise.all([
      apiClient.get("/costs/clusters/" + clusterId + "/summary", { params }),
      apiClient.get("/costs/clusters/" + clusterId + "/resource", { params }),
    ])
      .then((data) => {
        const costSummary = data[0].data;
        const computeCostResponse = data[1].data.items;
        const computeCostArray = [];

        computeCostResponse.map((element) =>
          computeCostArray.push(
            new ClusterComputeCost(
              element.timestamp,
              element.totalCost,
              element.costFallbackBillingMode,
              element.costOnDemandBillingMode,
              element.costSpotBillingMode,
              element.cpuCoreCount,
              element.cpuCoreUsage,
              element.cpuCost,
              element.ramGBCount,
              element.ramGBUsage,
              element.ramCost,
            )
          )
        );

        const clusterCostInfo = new ClusterCostInfo(
          costSummary.clusterName,
          stepSeconds,
          costSummary.clusterMonthCostCurrent,
          costSummary.clusterMonthEstimateCost,
          costSummary.clusterAvgDailyCost,
          costSummary.ClusterAvgHourlyCoreCost,
          computeCostArray
        );

        setClusterCostInfo(clusterCostInfo);
        console.log("clusterCostInfo=", clusterCostInfo);
        props.onChangeClusterCostInfo(clusterCostInfo);
      })
      .catch((err) => {
        console.error(err)
        setClusterCostInfo(null);
      });
  };

  useEffect(() => {
    let startDate = 0;
    let endDate = 0;

    if (onChangeDetails?.type === "absolute") {
      // valueOf return in milliseconds, backend need in seconds
      startDate = moment(onChangeDetails.startDate).valueOf() / 1e3;
      endDate = moment(onChangeDetails.endDate).valueOf() / 1e3;
    } else if (onChangeDetails?.type === "relative") {
      const currentTimestamp = moment().valueOf() / 1e3;
      const twelveDaysAgoTimestamp = moment()
        .subtract(onChangeDetails.amount, onChangeDetails.unit)
        .valueOf() / 1e3;
      startDate = twelveDaysAgoTimestamp;
      endDate = currentTimestamp;
    }

    let startDateInt = Math.floor(startDate);
    let endDateInt = Math.floor(endDate);
    let stepSeconds = 3600;
    // if the start time is less 24 hours before the end time, the step time shoud be days
    if (endDateInt - startDateInt >= 86400) {
      stepSeconds = 86400;
    }
    fetchClusterData(startDateInt, endDateInt, stepSeconds);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [onChangeDetails]);

  const handleChange = ({ detail }) => {
    setOnChangeDetails(detail.value);
  };

  return (
    <Header
      variant="awsui-h1-sticky"
      actions={
        <DateRangePicker
          value={onChangeDetails}
          onChange={handleChange}
          relativeOptions={[
            {
              key: "previous-5-minutes",
              amount: 5,
              unit: "minute",
              type: "relative",
            },
            {
              key: "previous-30-minutes",
              amount: 30,
              unit: "minute",
              type: "relative",
            },
            {
              key: "previous-1-hour",
              amount: 1,
              unit: "hour",
              type: "relative",
            },
            {
              key: "previous-6-hours",
              amount: 6,
              unit: "hour",
              type: "relative",
            },
          ]}
          isValidRange={(range) => {
            if (range.type === "absolute") {
              const [startDateWithoutTime] = range.startDate.split("T");
              const [endDateWithoutTime] = range.endDate.split("T");
              if (!startDateWithoutTime || !endDateWithoutTime) {
                return {
                  valid: false,
                  errorMessage:
                    "The selected date range is incomplete. Select a start and end date for the date range.",
                };
              }
              if (new Date(range.startDate) - new Date(range.endDate) > 0) {
                return {
                  valid: false,
                  errorMessage:
                    "The selected date range is invalid. The start date must be before the end date.",
                };
              }
            }
            return { valid: true };
          }}
          i18nStrings={{
            todayAriaLabel: "Today",
            nextMonthAriaLabel: "Next month",
            previousMonthAriaLabel: "Previous month",
            customRelativeRangeDurationLabel: "Duration",
            customRelativeRangeDurationPlaceholder: "Enter duration",
            customRelativeRangeOptionLabel: "Custom range",
            customRelativeRangeOptionDescription:
              "Set a custom range in the past",
            customRelativeRangeUnitLabel: "Unit of time",
            formatRelativeRange: (e) => {
              const n = 1 === e.amount ? e.unit : `${e.unit}s`;
              return `Last ${e.amount} ${n}`;
            },
            formatUnit: (e, n) => (1 === n ? e : `${e}s`),
            dateTimeConstraintText:
              "Range is 6 to 30 days. For date, use YYYY/MM/DD. For time, use 24 hr format.",
            relativeModeTitle: "Relative range",
            absoluteModeTitle: "Absolute range",
            relativeRangeSelectionHeading: "Choose a range",
            startDateLabel: "Start date",
            endDateLabel: "End date",
            startTimeLabel: "Start time",
            endTimeLabel: "End time",
            clearButtonLabel: "Clear and dismiss",
            cancelButtonLabel: "Cancel",
            applyButtonLabel: "Apply",
          }}
          placeholder="Filter by a date and time range"
        />
      }
    >
      {clusterName}
    </Header>
  );
}
