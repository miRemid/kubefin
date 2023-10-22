// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import {
  Box,
  BarChart,
  Container,
  Header,
  ColumnLayout,
  SpaceBetween,
} from "@cloudscape-design/components";
import { ClusterCostInfo } from "../model/cluster-cost";
import { CounterLink } from "../../../commons/common-components";
import {
  commonChartProps,
  barChartInstructions,
} from "./common";
import { dateDayFormatter, dateHourFormatter, keepTwoDecimal } from "../components-common";

export default function ClusterComputeCostChart(props) {
  let clusterCostInfo = props.clusterCostInfo;
  if (clusterCostInfo === undefined) {
    clusterCostInfo = new ClusterCostInfo();
  }
  const clusterComputeCostArrayByTimestamp =
    clusterCostInfo.clusterComputeCostArray;

  let totalFallbackSpend = 0;
  let totalOndemandSpend = 0;
  let totalSpotSpend = 0;
  let nodesData = [];
  let dayMaxCost = 1;

  // eslint-disable-next-line array-callback-return
  clusterComputeCostArrayByTimestamp.map((element) => {
    let fallbackCost = element.costFallbackBillingMode ? element.costFallbackBillingMode : 0;
    let ondemandCost = element.costOnDemandBillingMode ? element.costOnDemandBillingMode : 0;
    let spotCost = element.costSpotBillingMode ? element.costSpotBillingMode : 0;

    totalFallbackSpend += fallbackCost;
    totalOndemandSpend += ondemandCost;
    totalSpotSpend += spotCost;

    if (fallbackCost + ondemandCost + spotCost > dayMaxCost) {
      dayMaxCost = fallbackCost + ondemandCost + spotCost;
    }

    nodesData.push({
      // The returned time is in seconds, we should transform it to millsecond
      date: new Date(element.timestamp * 1000),
      ondemand: keepTwoDecimal(ondemandCost),
      fallback: keepTwoDecimal(fallbackCost),
      spot: keepTwoDecimal(spotCost),
    });
  });

  const nodesDomain = nodesData.map(({ date }) => date);

  const nodesSeries = [
    {
      title: "ondemand",
      type: "bar",
      data: nodesData.map((datum) => ({ x: datum.date, y: datum["ondemand"] })),
    },
    {
      title: "spot",
      type: "bar",
      data: nodesData.map((datum) => ({ x: datum.date, y: datum["spot"] })),
    },
  ];

  let xDateFormatter = dateHourFormatter;
  if (clusterCostInfo.dataStep > 3600) {
    xDateFormatter = dateDayFormatter;
  }
  return (
    <Container
      className="custom-dashboard-container"
      header={<Header variant="h2">Compute spend</Header>}
    >
      <SpaceBetween size="xxl">
        <ColumnLayout columns="3" variant="text-grid">
          <div>
            <Box variant="awsui-key-label">Total spend</Box>
            <CounterLink>
              {keepTwoDecimal(totalOndemandSpend + totalFallbackSpend + totalSpotSpend)}
            </CounterLink>
          </div>
          <div>
            <Box variant="awsui-key-label">Ondemand spend</Box>
            <CounterLink>{keepTwoDecimal(totalOndemandSpend)}</CounterLink>
          </div>
          <div>
            <Box variant="awsui-key-label">Spot spend</Box>
            <CounterLink>{keepTwoDecimal(totalSpotSpend)}</CounterLink>
          </div>
        </ColumnLayout>

        <BarChart
          {...commonChartProps}
          height={220}
          hideFilter={true}
          // The y-axis needs to be handled flexibly in order to achieve better display.
          yDomain={[0, dayMaxCost * 1.2]}
          xDomain={nodesDomain}
          xScaleType="categorical"
          stackedBars={true}
          series={nodesSeries}
          xTitle="Date"
          yTitle="Spend"
          ariaDescription={`Bar chart showing total instance hours per instance type over the last 15 days. ${barChartInstructions}`}
          i18nStrings={{
            ...commonChartProps.i18nStrings,
            filterLabel: "Filter displayed instance types",
            filterPlaceholder: "Filter instance types",
            xTickFormatter: xDateFormatter,
          }}
        />
      </SpaceBetween>
    </Container>
  );
}
