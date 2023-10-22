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
import { keepTwoDecimal } from '../components-common';

export default function ClusterMemory(props) {
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
          Memory (GiB)
        </Header>
      }
    >
      <ColumnLayout columns="4" variant="text-grid">
        <div>
          <Box variant="awsui-key-label">Total memory</Box>
          <CounterLink>{keepTwoDecimal(cluster.Memory.totalMEM)}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Memory requested</Box>
          <CounterLink>{keepTwoDecimal(cluster.Memory.requestedMEM)}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Memory allocatable</Box>
          <CounterLink>{keepTwoDecimal(cluster.Memory.allocatableMEM)}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Memory utilization</Box>
          <CounterLink>{keepTwoDecimal(cluster.Memory.MEMusage/cluster.Memory.totalMEM*100)}%</CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
