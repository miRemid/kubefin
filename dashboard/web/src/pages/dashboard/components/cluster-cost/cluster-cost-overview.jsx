/* eslint-disable react-hooks/exhaustive-deps */
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import {
  Box,
  Container,
  Header,
  ColumnLayout,
  HelpPanel,
  Icon,
} from "@cloudscape-design/components";
import {
  CounterLink,
  ExternalLinkItem,
  InfoLink,
} from "../../../commons/common-components";
import "../../../../styles/cluster-cost.scss";
import { ClusterCostInfo } from "../model/cluster-cost";
import { keepTwoDecimal } from "../components-common";

function CostInfo() {
  return (
    <HelpPanel
      header={<h2>Cluster cost overview</h2>}
      footer={
        <>
          <h3>
            Learn more{" "}
            <span role="img" aria-label="Icon external Link">
              <Icon name="external" />
            </span>
          </h3>
          <ul>
            <li>
              <ExternalLinkItem href="#" text="Service health dashboard" />
            </li>
            <li>
              <ExternalLinkItem href="#" text="Personal health dashboard" />
            </li>
          </ul>
        </>
      }
    >
      <p>There're 4 metrics to measure cluster cost:</p>
      <ul>
        <li>
          Current month spend: Current month factual compute spend that includes
          the end of the current day projection
        </li>
        <li>
          Monthly forecast: Takes into account factual workload uptime and
          forecast for the rest of the month. Forecast is calculated using 7 day
          moving average formula.
        </li>
        <li>Avg. daily cost: Average daily compute cost</li>
        <li>
          Avg. daily cost per resource: Average daily cost per provisioned
          resource
        </li>
      </ul>
    </HelpPanel>
  );
}

export default function ClusterCostOverview(props) {
  let clusterCostInfo = props.clusterCostInfo;
  if (clusterCostInfo === undefined || clusterCostInfo === null) {
    clusterCostInfo = new ClusterCostInfo();
  }
  return (
    <Container
      className="custom-dashboard-container"
      header={
        <Header
          variant="h2"
          info={<InfoLink onFollow={() => props.updateTools(<CostInfo />)} />}
        >
          Cluster cost overview
        </Header>
      }
    >
      <ColumnLayout columns="4" variant="text-grid">
        <div id="compute-spend-container">
          <Box variant="awsui-key-label">Current month spend</Box>
          <CounterLink>
            {clusterCostInfo.clusterCostCurrent.toFixed(2)}
          </CounterLink>
        </div>

        <div id="avg-monthly-cost-container">
          <Box variant="awsui-key-label">Monthly forecast</Box>
          <CounterLink>
            {keepTwoDecimal(clusterCostInfo.clusterMonthEstimateCost)}
          </CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Avg. daily cost</Box>
          <CounterLink>
            {keepTwoDecimal(clusterCostInfo.clusterAvgDailyCost)}
          </CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Avg. daily cost per resource</Box>
          <CounterLink>
            {keepTwoDecimal(clusterCostInfo.clusterAvgHourlyCoreCost)}
          </CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
