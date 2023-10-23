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

export default function ClusterCPU(props) {
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
      new CPU(0,0,0,0,[0,0],[0,0])
    );
  }

  return (
    <Container
      className="custom-dashboard-container"
      header={
        <Header variant="h2">
          CPU (vCPU)
        </Header>
      }
    >
      <ColumnLayout columns="4" variant="text-grid">
        <div>
          <Box variant="awsui-key-label">All CPU cores</Box>
          <CounterLink>{cluster.CPU.totalCPUCores}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">CPU requested</Box>
          <CounterLink>{cluster.CPU.requestedCPU}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">CPU allocatable</Box>
          <CounterLink>{cluster.CPU.allocatableCPU}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">CPU utilization</Box>
          <CounterLink>{keepTwoDecimal(cluster.CPU.CPUusage/cluster.CPU.totalCPUCores*100)}%</CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
