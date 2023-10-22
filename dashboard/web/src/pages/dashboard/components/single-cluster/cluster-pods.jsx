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

export default function ClusterPods(props) {
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
          Pods
        </Header>
      }
    >
      <ColumnLayout columns="4" variant="text-grid">
        <div>
          <Box variant="awsui-key-label">All pod</Box>
          <CounterLink>{cluster.Pod.podNum}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Scheduled pods</Box>
          <CounterLink>{cluster.Pod.scheduledPodNum}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">unscheduled pods</Box>
          <CounterLink>{cluster.Pod.unscheduledPodNum}</CounterLink>
        </div>
        <div>
          <Box variant="awsui-key-label">Scheduled pods percentage</Box>
          <CounterLink>{parseFloat(cluster.Pod.scheduledRatio.toFixed(2))}%</CounterLink>
        </div>
      </ColumnLayout>
    </Container>
  );
}
