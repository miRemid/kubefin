// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React, { useState, useEffect } from "react";
import {
  AppLayout,
  ContentLayout,
  Grid,
  TopNavigation,
  SideNavigation,
  Box,
  Popover,
} from "@cloudscape-design/components";
import { useParams } from "react-router";
import "@cloudscape-design/global-styles/dark-mode-utils.css";

import "../../../../styles/top-navigation.scss";

import "../../../../styles/dashboard.scss";
import "../../../../styles/density-switch-images.scss";

import { appLayoutLabels } from "../../../../common/labels";

import { SingleClusterDashboardHeader } from "./single-cluster-header";
import { Notifications } from "../../../commons/common-components";
import ClusterInfo from "./cluster-info";
import apiClient from "../../../../common/network/http-common";
import { Cluster } from "../model/cluster";
import { Pod } from "../model/pod";
import { Node } from "../model/node";
import { Memory } from "../model/memory";
import { CPU } from "../model/cpu";
import ClusterNodes from "./cluster-nodes";
import ClusterPods from "./cluster-pods";
import ClusterCPU from "./cluster-cpu";
import ClusterMemory from "./cluster-memory";
import ClusterCPUUtilization from "./cluster-cpu-utilization";
import ClusterMemoryUtilization from "./cluster-memory-utilization";
import { keepTwoDecimal } from "../components-common";

function Content(props) {
  return (
    <Grid
      gridDefinition={[
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 6, m: 6, default: 12 } },
        { colspan: { l: 6, m: 6, default: 12 } },
        { colspan: { l: 6, m: 6, default: 12 } },
        { colspan: { l: 6, m: 6, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
      ]}
    >
      <ClusterInfo cluster={props.cluster} />
      <ClusterNodes cluster={props.cluster} />
      <ClusterPods cluster={props.cluster} />
      <ClusterCPU cluster={props.cluster} />
      <ClusterMemory cluster={props.cluster} />
      <ClusterCPUUtilization cluster={props.cluster} />
      <ClusterMemoryUtilization cluster={props.cluster} />
    </Grid>
  );
}

export function SingleClusterDashboard() {
  const { clusterId } = useParams();
  const [toolsOpen, setToolsOpen] = useState(false);
  const [toolsContent, setToolsContent] = useState();
  const [navigationOpen, setNavigationOpen] = useState(true);
  const [cluster, setCluster] = useState();
  const [activeHref, setActiveHref] = useState("/dashboard/" + clusterId);

  useEffect(() => {
    const fetchClusterData = async () => {
      await Promise.all([
        apiClient.get("/costs/clusters/" + clusterId + "/summary"),
        apiClient.get("/metrics/clusters/" + clusterId + "/summary"),
        apiClient.get("/metrics/clusters/" + clusterId + "/cpu"),
        apiClient.get("/metrics/clusters/" + clusterId + "/memory"),
      ])
        .then((data) => {
          console.log("data=",data)
          const costSummary = data[0].data;
          const metricsSummary = data[1].data;
          const cpuInfo = data[2].data;
          const memoryInfo = data[3].data;
          const requestedCPUArray = cpuInfo.resourceRequestValues;
          const allocatableCPUArray = cpuInfo.resourceAllocatableValues;
          const totalCPUArray = cpuInfo.resourceTotalValues;
          const requestedMemoryArray = memoryInfo.resourceRequestValues;
          const allocatableMemoryArray = memoryInfo.resourceAllocatableValues;
          const totalMemoryArray = memoryInfo.resourceTotalValues;

          const pod = new Pod(
            metricsSummary.podTotalCurrent,
            metricsSummary.podScheduledCurrent,
            metricsSummary.podUnscheduledCurrent,
            metricsSummary.podScheduledCurrent / metricsSummary.podTotalCurrent
          );
          const node = new Node(
            metricsSummary.nodeNumbersCurrent,
            metricsSummary.onDemandBillingNodeNumbersCurrent,
            metricsSummary.spotBillingNodeNumbersCurrent,
            metricsSummary.fallbackBillingNodeNumbersCurrent
          );
          const memory = new Memory(
            keepTwoDecimal(metricsSummary.ramGBTotal),
            keepTwoDecimal(metricsSummary.ramGBRequest),
            keepTwoDecimal(metricsSummary.ramGBCapacity) - keepTwoDecimal(metricsSummary.ramGBRequest),
            keepTwoDecimal(metricsSummary.ramGBUsage),
            requestedMemoryArray,
            allocatableMemoryArray,
            totalMemoryArray
          );
          const cpu = new CPU(
            keepTwoDecimal(metricsSummary.cpuCoreTotal),
            keepTwoDecimal(metricsSummary.cpuCoreRequest),
            keepTwoDecimal(
              metricsSummary.cpuCoreCapacity - metricsSummary.cpuCoreRequest
            ),
            keepTwoDecimal(metricsSummary.cpuCoreUsage),
            requestedCPUArray,
            allocatableCPUArray,
            totalCPUArray
          );

          const cluster = new Cluster(
            costSummary.clusterId,
            costSummary.clusterName,
            costSummary.clusterRegion,
            costSummary.cloudProvider,
            pod,
            node,
            memory,
            cpu,
            costSummary.clusterConnectionSate === "running"
              ? "Active"
              : "Deactivated"
          );
          setCluster(cluster);
        })
        .catch((err) => {
          setCluster(null);
        });
    };
    fetchClusterData();
  }, [clusterId]);

  const clusterName =
    cluster === undefined || cluster === null ? "Loading" : cluster.clusterName;
  const KFNavHeader = { text: clusterName, href: "#/" };
  const navItems = [
    {
      type: "section",
      text: "Dashboard",
      items: [
        {
          type: "link",
          text: "all clusters",
          href: "/dashboard",
        },
        {
          type: "link",
          text: clusterName,
          href: "/dashboard/" + clusterId,
          info: (
            <Box color="text-status-info" display="inline">
              <Popover
                header="cluster dashboard"
                size="medium"
                triggerType="text"
                content={<>Show the resources information in your cluster</>}
                renderWithPortal={true}
                dismissAriaLabel="Close"
              >
                <Box
                  color="text-status-info"
                  fontSize="body-s"
                  fontWeight="bold"
                  data-testid="new-feature-announcement-trigger"
                >
                  Info
                </Box>
              </Popover>
            </Box>
          ),
        },
      ],
    },
    {
      type: "section",
      text: "Cost allocation",
      items: [
        {
          type: "link",
          text: "cluster cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/cluster",
        },
        {
          type: "link",
          text: "workload cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/workload",
        },
        {
          type: "link",
          text: "namespace cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/namespace",
        },
      ],
    },
    // {
    //   type: "section",
    //   text: "Cost Optimization",
    //   items: [{ type: "link", text: "Available Savings", href: "#/database" }],
    // },
  ];

  const i18nStrings = {
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
          <ContentLayout
            header={
              <SingleClusterDashboardHeader
                updateTools={updateTools}
                cluster={cluster}
              />
            }
          >
            <Content updateTools={updateTools} cluster={cluster} />
          </ContentLayout>
        }
        headerSelector="#header"
        tools={toolsContent}
        toolsOpen={toolsOpen}
        navigationOpen={navigationOpen}
        navigation={
          <SideNavigation
            activeHref={activeHref}
            header={KFNavHeader}
            items={navItems}
          />
        }
        onNavigationChange={({ detail }) => setNavigationOpen(detail.open)}
        onToolsChange={({ detail }) => setToolsOpen(detail.open)}
        ariaLabels={appLayoutLabels}
        notifications={<Notifications />}
      />
    </>
  );
}
