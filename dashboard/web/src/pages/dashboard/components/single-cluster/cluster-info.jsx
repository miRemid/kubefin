// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from 'react';
import { ColumnLayout, Container, Header, Box, StatusIndicator } from '@cloudscape-design/components';
import CopyText from "../../../commons/copy-text";
import { Cluster } from '../model/cluster';

export default function ClusterInfo(props) {
  let cluster = props.cluster;
  if(cluster === undefined || cluster === null){
    cluster = new Cluster(
      "-",
      "-",
      "-",
      "-"
    );
  }
  return (
    <Container header={<Header variant="h2">Cluster Info</Header>}>
    <ColumnLayout columns={4} variant="text-grid">
      <div>
        <Box variant="awsui-key-label">Cloud provider</Box>
        <div>{cluster.cloudProvider}</div>
      </div>
      <div>
        <Box variant="awsui-key-label">Region</Box>
        <div>{cluster.clusterRegion}</div>
      </div>
      <div>
        <Box variant="awsui-key-label">Cluster Id</Box>
        <CopyText
          copyText={cluster.clusterId}
          copyButtonLabel="Copy cluster id"
          successText="cluster id copy successful"
          errorText="cluster id copy failed"
        />
      </div>
      <div>
        <Box variant="awsui-key-label">State</Box>
        <StatusIndicator type={cluster.state === 'Deactivated' ? 'error' : 'success'}>{cluster.state}</StatusIndicator>
      </div>
    </ColumnLayout>
  </Container>
  );
}
