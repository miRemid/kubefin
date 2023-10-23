// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import { HelpPanel, Icon, Header } from "@cloudscape-design/components";
import { ExternalLinkItem, InfoLink } from "../../commons/common-components";

export function AllClustersDashboardInfo() {
  return (
    <HelpPanel
      header={<h2>All clusters dashboard</h2>}
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
              <ExternalLinkItem
                href="https://kubefin.dev/docs/"
                text="Get start with KubeFin"
              />
            </li>
          </ul>
        </>
      }
    >
      <p>
          This dashboard contains the information of node and spend in all clusters you have connected to KubeFin.
      </p>
    </HelpPanel>
  );
}

export function DashboardHeader(props) {
  return (
    <Header
      variant="awsui-h1-sticky"
      info={
        <InfoLink
          onFollow={() => props.updateTools(<AllClustersDashboardInfo />)}
          ariaLabel={"Information about service dashboard."}
        />
      }
    >
      All clusters dashboard
    </Header>
  );
}
