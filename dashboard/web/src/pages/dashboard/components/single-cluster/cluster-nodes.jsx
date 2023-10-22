// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from 'react';
import { Box, Container, Header, ColumnLayout } from '@cloudscape-design/components';
import { CounterLink } from '../../../commons/common-components';
import { Cluster } from '../model/cluster';
import { Pod } from '../model/pod';
import { Node } from '../model/node';
import { Memory } from '../model/memory';
import { CPU } from '../model/cpu';

export default function ClusterNodes(props) {
  let cluster = props.cluster;
  if(cluster === undefined || cluster === null){
    cluster = new Cluster(
      "-",
      "-",
      "-",
      "-",
      new Pod(0,0,0,0),
      new Node(0,0,0,0),
      new Memory(0,0,0,0),
      new CPU(0,0,0,0)
    );
  }

  return (
    <Container
      className="custom-dashboard-container"
      header={
        <Header variant="h2">
          Nodes
        </Header>
      }
    >
      <ColumnLayout columns="3" variant="text-grid">
        <div>
          <Box variant="awsui-key-label">All nodes</Box>
          <CounterLink>{cluster.Node.nodeNum}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Ondemand nodes</Box>
          <CounterLink>{cluster.Node.ondemandNodes}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Spot nodes</Box>
          <CounterLink>{cluster.Node.spotNodes}</CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
