// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React, { useState } from "react";
import {
  AppLayout,
  ContentLayout,
  Grid,
  TopNavigation,
  SideNavigation,
  Box,
  Popover
} from "@cloudscape-design/components";
import { useParams } from "react-router";
import "@cloudscape-design/global-styles/dark-mode-utils.css";

import "../../../../styles/top-navigation.scss";

import "../../../../styles/dashboard.scss";

import { appLayoutLabels } from "../../../../common/labels";

import { Notifications } from "../../../commons/common-components";
import { ClusterCostHeader } from "./cluster-cost-header";
import { ClusterComputeCost } from "../model/cluster-compute-cost";

import ClusterCostOverview from "./cluster-cost-overview";
import ClusterComputeCostChart from "./compute-cost-chart";
import ClusterComputeCostList from "../table/cluster-compute-cost-list";

function Content(props) {
  return (
    <Grid
      gridDefinition={[
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
      ]}
    >
      <ClusterCostOverview
        clusterCostInfo={props.clusterCostInfo}
        updateTools={props.updateTools}
      />
      <ClusterComputeCostChart clusterCostInfo={props.clusterCostInfo} />
      <ClusterComputeCostList clusterCostInfo={props.clusterCostInfo} />
    </Grid>
  );
}
export function ClusterCostDashboard() {
  const { clusterId } = useParams();
  const [toolsOpen, setToolsOpen] = useState(false);
  const [toolsContent, setToolsContent] = useState();
  const [navigationOpen, setNavigationOpen] = useState(true);
  const [clusterCostInfo, setClusterCostInfo] = useState({
    clusterName: "-",
    clusterCostCurrent: 0,
    clusterMonthEstimateCost: 0,
    clusterAvgDailyCost: 0,
    clusterAvgHourlyCoreCost: 0,
    clusterComputeCostArray: [new ClusterComputeCost()],
  });

  function handleChangeClusterCostInfo(clusterCostInfo) {
    setClusterCostInfo(clusterCostInfo);
  }

  const clusterName =
    clusterCostInfo === undefined ? "Loading" : clusterCostInfo.clusterName;
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
              <ClusterCostHeader
                updateTools={updateTools}
                onChangeClusterCostInfo={(clusterCostInfo) =>
                  handleChangeClusterCostInfo(clusterCostInfo)
                }
              />
            }
          >
            <Content
              updateTools={updateTools}
              clusterCostInfo={clusterCostInfo}
            />
          </ContentLayout>
        }
        headerSelector="#header"
        tools={toolsContent}
        toolsOpen={toolsOpen}
        navigationOpen={navigationOpen}
        navigation={
          <SideNavigation
            activeHref={"/cost/" + clusterId}
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
