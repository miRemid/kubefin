// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React, { useState, useEffect } from "react";
import { Header, DateRangePicker } from "@cloudscape-design/components";
import { useParams } from "react-router";
import apiClient from "../../../../common/network/http-common";
import moment from "moment";
import { NamespaceCostInfo } from "../model/namespace/namespace-cost";

export function NamespaceCostHeader(props) {
  const { clusterId } = useParams();
  const [onChangeDetails, setOnChangeDetails] = useState({
    type: "relative",
    amount: 7,
    unit: "hour",
  });

  const fetchClusterData = async (startTime, endTime, stepSeconds) => {
    const params = {
      startTime: startTime,
      endTime: endTime,
      stepSeconds: stepSeconds,
    };
    await Promise.all([
      apiClient.get("/costs/clusters/" + clusterId + "/namespace", { params }),
    ])
      .then((data) => {
        const nsCostResponse = data[0].data.items;
        let nsCostMap = new Map();

        nsCostResponse.map((element) => {
          const namespace = element.namespace;
          let nsCostInfo = nsCostMap.get(namespace);
          nsCostInfo = new NamespaceCostInfo();

          let podCount = 0;
          let podCountIndex = 0;
          let cpuRequest = 0;
          let cpuRequestIndex = 0;
          let ramGBRequest = 0;
          let ramGBRequestIndex = 0;
          let totalCost = 0;

          element.costList.map((cost) => {
            podCount += cost.podCount === undefined ? 0 : cost.podCount;
            cpuRequest += cost.cpuRequest === undefined ? 0 : cost.cpuRequest;
            ramGBRequest +=
              cost.ramGBRequest === undefined ? 0 : cost.ramGBRequest;
            totalCost += cost.totalCost === undefined ? 0 : cost.totalCost;

            podCountIndex += 1.0;
            cpuRequestIndex += 1.0;
            ramGBRequestIndex += 1.0;
          });

          nsCostInfo.namespace = namespace
          nsCostInfo.podCount = podCount / podCountIndex;
          nsCostInfo.cpuRequest = cpuRequest / cpuRequestIndex;
          nsCostInfo.ramGBRequest = ramGBRequest / ramGBRequestIndex;
          nsCostInfo.totalCost = totalCost;

          nsCostMap.set(namespace, nsCostInfo);
          return nsCostMap;
        });

        console.log("namespaceCostMap#header=", nsCostMap);
        props.onChangeNamespaceCostMap(nsCostMap);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  useEffect(() => {
    let startDate = 0;
    let endDate = 0;

    if (onChangeDetails?.type === "absolute") {
      startDate = moment(onChangeDetails.startDate).valueOf() / 1e3;
      endDate = moment(onChangeDetails.endDate).valueOf() / 1e3;
    } else if (onChangeDetails?.type === "relative") {
      const currentTimestamp = moment().valueOf() / 1e3;
      const twelveDaysAgoTimestamp =
        moment()
          .subtract(onChangeDetails.amount, onChangeDetails.unit)
          .valueOf() / 1e3;
      endDate = currentTimestamp;
      startDate = twelveDaysAgoTimestamp;
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
      Namespace Cost Analysis
    </Header>
  );
}
