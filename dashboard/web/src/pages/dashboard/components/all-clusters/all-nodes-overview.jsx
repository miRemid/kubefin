// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from 'react';
import { Box, Container, Header, ColumnLayout } from '@cloudscape-design/components';
import { CounterLink } from '../../../commons/common-components';

export default function AllNodesOverview(props) {
  return (
    <Container
      className="custom-dashboard-container"
      header={
        <Header variant="h2">
          Nodes overview
        </Header>
      }
    >
      <ColumnLayout columns="3" variant="text-grid">
        <div>
          <Box variant="awsui-key-label">Total nodes</Box>
          <CounterLink>{props.totalNodes}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">On-demand nodes</Box>
          <CounterLink>{props.totalOndemandNodes}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Spot nodes</Box>
          <CounterLink>{props.totalSpotNodes}</CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
