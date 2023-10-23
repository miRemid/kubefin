// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import {
  Container,
  Header,
  PieChart,
  Button,
  Box,
} from "@cloudscape-design/components";
import { keepThreeDecimal, percentageFormatter } from "../components-common";

const i18nStrings = {
  detailsValue: "Value",
  detailsPercentage: "Percentage",
  filterLabel: "Filter displayed data",
  filterPlaceholder: "Filter data",
  filterSelectedAriaLabel: "selected",
  detailPopoverDismissAriaLabel: "Dismiss",
  legendAriaLabel: "Legend",
  chartAriaRoleDescription: "pie chart",
  segmentAriaRoleDescription: "segment",
};
export default function NamespaceCostAllocation(props) {
  const namespaceCostMap = props.namespaceCostMap;
  console.log("namespaceCostMap#header",JSON.stringify(namespaceCostMap));
  if (!(namespaceCostMap instanceof Map) || namespaceCostMap.size === 0) {
    return <div>The namespace Cost is empty or not a valid Map.</div>;
  }

  const totalCostMap = new Map();
  let totalCost = 0;
  namespaceCostMap.forEach((data, key) => {
    const workloadCost =  data.totalCost;
    totalCost += workloadCost;
    totalCostMap.set(key, workloadCost);
  });
  totalCostMap.forEach((value, key) => {
    totalCostMap.set(key, keepThreeDecimal(value));
  });
  console.log("totalCostMap = ", totalCostMap);
  const totalCostData = Array.from(totalCostMap).map(([title, value]) => ({
    title,
    value,
  }));

  return (
    <Container
      className="custom-dashboard-container"
      header={<Header variant="h2">Namespace Cost Allocation</Header>}
    >
      <PieChart
        size="large"
        data={totalCostData}
        ariaLabel="Namespace allocation chart"
        ariaDescription="Namespace allocation chart."
        hideFilter={true}
        segmentDescription={(datum, sum) =>
          `$${datum.value}, ${percentageFormatter(datum.value / sum)}`
        }
        i18nStrings={i18nStrings}
        errorText="Error loading data."
        hideLegend={false}
        innerMetricDescription="Total Cost"
        innerMetricValue={"$" + keepThreeDecimal(totalCost)}
        loadingText="Loading chart"
        recoveryText="Retry"
        variant="donut"
        empty={
          <Box textAlign="center" color="inherit">
            <b>No data available</b>
            <Box variant="p" color="inherit">
              There is no data available
            </Box>
          </Box>
        }
        noMatch={
          <Box textAlign="center" color="inherit">
            <b>No matching data</b>
            <Box variant="p" color="inherit">
              There is no matching data to display
            </Box>
            <Button>Clear filter</Button>
          </Box>
        }
      />
    </Container>
  );
}
