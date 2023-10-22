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
  Popover,
} from "@cloudscape-design/components";
import { useParams } from "react-router";
import "@cloudscape-design/global-styles/dark-mode-utils.css";

import "../../../../styles/top-navigation.scss";

import "../../../../styles/dashboard.scss";

import { appLayoutLabels } from "../../../../common/labels";

import { Notifications } from "../../../commons/common-components";
import { NamespaceCostHeader } from "./namespace-cost-header";

import NamespaceCostAllocationPieChart from "./namespace-cost-allocation-pie-chart";
import NSCostTableList from "../table/ns-cost-table-list";

function Content(props) {
  return (
    <Grid
      gridDefinition={[
        { colspan: { l: 12, m: 12, default: 12 } },
        { colspan: { l: 12, m: 12, default: 12 } },
      ]}
    >
      <NamespaceCostAllocationPieChart
        updateTools={props.updateTools}
        namespaceCostMap={props.namespaceCostMap}
      />
      <NSCostTableList namespaceCostMap={props.namespaceCostMap} />
    </Grid>
  );
}
export function NamespaceCostDashboard() {
  const { clusterId } = useParams();
  const { clusterName } = useParams();

  const [navigationOpen, setNavigationOpen] = useState(true);
  const [namespaceCostMap, setNamespaceCostMap] = useState({});

  function handleChangeNamespaceCostMap(namespaceCostMap) {
    setNamespaceCostMap(namespaceCostMap);
  }
  const KFNavHeader = { text: "Namespace cost", href: "#/" };
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
  ];
  const i18nStrings = {
    errorIconAriaLabel: "Error",
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
              <NamespaceCostHeader
                onChangeNamespaceCostMap={(namespaceCostMap) =>
                  handleChangeNamespaceCostMap(namespaceCostMap)
                }
              />
            }
          >
            <Content
              namespaceCostMap={namespaceCostMap}
            />
          </ContentLayout>
        }
        headerSelector="#header"
        navigationOpen={navigationOpen}
        navigation={
          <SideNavigation
            activeHref={"/cost/" + clusterId}
            header={KFNavHeader}
            items={navItems}
          />
        }
        onNavigationChange={({ detail }) => setNavigationOpen(detail.open)}
        ariaLabels={appLayoutLabels}
        notifications={<Notifications />}
      />
    </>
  );
}
