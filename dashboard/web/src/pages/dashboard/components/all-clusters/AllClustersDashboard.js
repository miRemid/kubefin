// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React, { useState, useEffect } from "react";
import {
  AppLayout,
  ContentLayout,
  Grid,
  TopNavigation
} from "@cloudscape-design/components";

import "@cloudscape-design/global-styles/dark-mode-utils.css";

import "../../../../styles/top-navigation.scss";

import "../../../../styles/dashboard.scss";
import "../../../../styles/density-switch-images.scss";

import { appLayoutLabels } from "../../../../common/labels";

import { DashboardHeader, AllClustersDashboardInfo } from "../header";
import { Notifications } from "../../../commons/common-components";
import ClustersCostOverview from "./clusters-cost-overview";
import AllNodesOverview from "./all-nodes-overview";
import ClustersTableList from "./clusters-table-list";
import CapacityAllocation from "./capacity-allocation-chart";
import apiClient from "../../../../common/network/http-common";

function Content(props) {
  return (
    <Grid
      gridDefinition={[
        { colspan: { l: 4, m: 4, default: 12 } },
        { colspan: { l: 8, m: 8, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
      ]}
    >
      <ClustersCostOverview
        updateTools={props.updateTools}
        totalCost={props.totalCost}
      />
      <AllNodesOverview
        totalNodes={props.totalNodes}
        totalOndemandNodes={props.totalOndemandNodes}
        totalSpotNodes={props.totalSpotNodes}
        totalFallbackNodes={props.totalFallbackNodes}
      />
      <CapacityAllocation
        updateTools={props.updateTools}
        totalCPUCores={props.totalCPUCores}
        systemReservedCPU={props.systemReservedCPU}
        workloadReservedCPU={props.workloadReservedCPU}
        allocatableCPU={props.allocatableCPU}
        totalMemory={props.totalMemory}
        systemReservedMemory={props.systemReservedMemory}
        workloadReservedMemory={props.workloadReservedMemory}
        allocatableMemory={props.allocatableMemory}
      />
      <ClustersTableList
        costSummary = {props.costSummary}
        metricsSummary = {props.metricsSummary}
      />
    </Grid>
  );
}

export function AllClustersDashboard() {
  const [toolsOpen, setToolsOpen] = useState(false);
  const [toolsContent, setToolsContent] = useState(<AllClustersDashboardInfo />);
  const [totalCost, setTotalCost] = useState(0);
  const [totalNodes, setTotalNodes] = useState(0);
  const [totalOndemandNodes, setTotalOndemandNodes] = useState(0);
  const [totalSpotNodes, setTotalSpotNodes] = useState(0);
  const [totalFallbackNodes, setTotalFallbackNodes] = useState(0);
  const [totalCPUCores, setTotalCPUCores] = useState(0);
  const [systemReservedCPU, setSystemReservedCPU] = useState(0);
  const [workloadReservedCPU, setWorkloadReservedCPU] = useState(0);
  const [allocatableCPU, setAllocatableCPU] = useState(0);
  const [totalMemory, setTotalMemory] = useState(0);
  const [systemReservedMemory, setSystemReservedMemory] = useState(0);
  const [workloadReservedMemory, setWorkloadReservedMemory] = useState(0);
  const [allocatableMemory, setAllocatableMemory] = useState(0);
  const [costSummary, setCostSummary] = useState([]);
  const [metricsSummary, setMetricsSummary] = useState([]);

  const fetchClustersData = async () => {
    await Promise.all([
      apiClient.get("/costs/summary"),
      apiClient.get("/metrics/summary"),
    ])
      .then((data) => {
        const costSummary = data[0].data.items;
        const metricsSummary = data[1].data.items;
        //cost overview
        const costItems = costSummary.map(
          (item) => item.clusterMonthCostCurrent
        );
        //nodes overview
        const nodesItems = metricsSummary.map(
          (item) => item.nodeNumbersCurrent
        );
        const ondemandNodesItems = metricsSummary.map(
          (item) => item.onDemandBillingNodeNumbersCurrent ? item.onDemandBillingNodeNumbersCurrent : 0
        );
        const spotNodesItems = metricsSummary.map(
          (item) => item.spotBillingNodeNumbersCurrent ? item.spotBillingNodeNumbersCurrent : 0
        );
        const fallbackNodes = metricsSummary.map(
          (item) => item.fallbackBillingNodeNumbersCurrent ? item.fallbackBillingNodeNumbersCurrent : 0
        );
        // cpu and memory overview
        const CPUCoresItems = metricsSummary.map((item) => item.cpuCoreTotal);
        const requestedCPUItems = metricsSummary.map(
          (item) => item.cpuCoreRequest
        );
        const availableCPUItems = metricsSummary.map(
          (item) => item.cpuCoreCapacity
        );
        const MemoryItems = metricsSummary.map((item) => item.ramGBTotal);
        const requestedMemoryItems = metricsSummary.map(
          (item) => item.ramGBRequest
        );
        const availableMemoryItems = metricsSummary.map(
          (item) => item.ramGBCapacity
        );
        const totalCost = costItems.reduce(function (result, item) {
          return result + item;
        }, 0);
        const totalNodes = nodesItems.reduce(function (result, item) {
          return result + item;
        }, 0);
        const totalOndemandNodes = ondemandNodesItems.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        const totalSpotNodes = spotNodesItems.reduce(function (result, item) {
          return result + item;
        }, 0);
        const totalFallbackNodes = fallbackNodes.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        const totalCPUCores = CPUCoresItems.reduce(function (result, item) {
          return result + item;
        }, 0);
        const totalRequestedCPU = requestedCPUItems.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        const totalAvailableCPU = availableCPUItems.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        const totalMemory = MemoryItems.reduce(function (result, item) {
          return result + item;
        }, 0);
        const totalRequestedMemory = requestedMemoryItems.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        const totalAvailableMemory = availableMemoryItems.reduce(function (
          result,
          item
        ) {
          return result + item;
        },
        0);
        setTotalCost(totalCost);
        setTotalNodes(totalNodes);
        setTotalOndemandNodes(totalOndemandNodes);
        setTotalSpotNodes(totalSpotNodes);
        setTotalFallbackNodes(totalFallbackNodes);
        setTotalCPUCores(totalCPUCores);
        setSystemReservedCPU(totalCPUCores - totalAvailableCPU);
        setWorkloadReservedCPU(totalRequestedCPU);
        setAllocatableCPU(totalAvailableCPU - totalRequestedCPU);
        setTotalMemory(totalMemory);
        setSystemReservedMemory(totalMemory - totalAvailableMemory);
        setWorkloadReservedMemory(totalRequestedMemory);
        setAllocatableMemory(totalAvailableMemory - totalRequestedMemory);
        setCostSummary(costSummary);
        setMetricsSummary(metricsSummary);
      })
      .catch((err) => {
        setTotalCost(0);
        setTotalNodes(0);
        setTotalOndemandNodes(0);
        setTotalSpotNodes(0);
        setTotalFallbackNodes(0);
        setTotalCPUCores(0);
        setSystemReservedCPU(0);
        setWorkloadReservedCPU(0);
        setAllocatableCPU(0);
        setTotalMemory(0);
        setSystemReservedMemory(0);
        setWorkloadReservedMemory(0);
        setAllocatableMemory(0);
        setCostSummary([]);
        setMetricsSummary([]);
      });
  };

  useEffect(() => {
    fetchClustersData();
  }, []);
  const i18nStrings = {
    stepNumberLabel: (stepNumber) => `Step ${stepNumber}`,
    collapsedStepsLabel: (stepNumber, stepsCount) =>
      `Step ${stepNumber} of ${stepsCount}`,
    errorIconAriaLabel: "Error",
    navigationAriaLabel: "Steps",
    cancelButton: "Cancel",
    previousButton: "Previous",
    nextButton: "Next",
    submitButton: "I run the script",
    optional: "optional",
    searchIconAriaLabel: "Search",
    searchDismissIconAriaLabel: "Close search",
    overflowMenuTriggerText: "More",
    overflowMenuTitleText: "All",
    overflowMenuBackIconAriaLabel: "Back",
    overflowMenuDismissIconAriaLabel: "Close menu",
  };
  const updateTools = (toolsContent) => {
    setToolsOpen(true);
    setToolsContent(toolsContent);
  };

  return (
    <>
      <TopNavigation
        i18nStrings={i18nStrings}
        identity={{
          href: "/",
          title: "KubeFin",
        }}
        utilities={[
          {
            type: "button",
            text: "Documentation",
            href: "https://kubefin.dev",
            external: true,
            externalIconAriaLabel: " (opens in a new tab)",
          },
        ]}
      />
      <AppLayout
        content={
          <ContentLayout header={<DashboardHeader updateTools={updateTools} />}>
            <Content
              updateTools={updateTools}
              totalCost={totalCost}
              totalNodes={totalNodes}
              totalOndemandNodes={totalOndemandNodes}
              totalSpotNodes={totalSpotNodes}
              totalFallbackNodes={totalFallbackNodes}
              totalCPUCores={totalCPUCores}
              systemReservedCPU={systemReservedCPU}
              workloadReservedCPU={workloadReservedCPU}
              allocatableCPU={allocatableCPU}
              totalMemory={totalMemory}
              systemReservedMemory={systemReservedMemory}
              workloadReservedMemory={workloadReservedMemory}
              allocatableMemory={allocatableMemory}
              costSummary = {costSummary}
              metricsSummary = {metricsSummary}
            />
          </ContentLayout>
        }
        headerSelector="#header"
        tools={toolsContent}
        toolsOpen={toolsOpen}
        navigationHide={true}
        onToolsChange={({ detail }) => setToolsOpen(detail.open)}
        ariaLabels={appLayoutLabels}
        notifications={<Notifications />}
      />
    </>
  );
}
