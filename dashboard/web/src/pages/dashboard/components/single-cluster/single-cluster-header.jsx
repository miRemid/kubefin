// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import { Header } from "@cloudscape-design/components";

export function SingleClusterDashboardHeader(props) {
  const clusterName =
    props.cluster === undefined || props.cluster === null ? "Connecting" : props.cluster.clusterName;
    const now = new Date();
    const options = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    };
    const formattedTime = now.toLocaleString('zh-CN', options);
  return (
    <Header variant="awsui-h1-sticky" description={formattedTime}>
      {clusterName}
    </Header>
  );
}
