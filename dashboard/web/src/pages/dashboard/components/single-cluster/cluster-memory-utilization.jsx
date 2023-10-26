// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import {
  Button,
  Box,
  AreaChart,
  Container,
  Header,
} from "@cloudscape-design/components";
import { Cluster } from "../model/cluster";
import { Pod } from "../model/pod";
import { Node } from "../model/node";
import { Memory } from "../model/memory";
import { CPU } from "../model/cpu";
import { keepTwoDecimal } from "../components-common";
import { ParseArrayIntoTimeSeriesData } from "../../../commons/common-components";

export default function ClusterMemoryUtilization(props) {
  let cluster = props.cluster;
  if (cluster === undefined || cluster === null) {
    cluster = new Cluster(
      "-",
      "-",
      "-",
      "-",
      new Pod(0, 0, 0, 0),
      new Node(0, 0, 0, 0),
      new Memory(0, 0, 0, 0, [0, 0], [0, 0], [0, 0], [0, 0]),
      new CPU(0, 0, 0, 0, [0, 0], [0, 0], [0, 0], [0, 0])
    );
  }

  const requestedMemoryArray = cluster.Memory.requestedMemoryArray;
  const availableMemoryArray = cluster.Memory.availableMemoryArray;
  const systemReservedMemoryArray = cluster.Memory.systemReservedMemoryArray;
  const memoryTotal = cluster.Memory.totalMemoryArray;

  console.log(systemReservedMemoryArray);
  var memoryRequestedData = ParseArrayIntoTimeSeriesData(requestedMemoryArray);
  var memoryAvaiableData = ParseArrayIntoTimeSeriesData(availableMemoryArray);
  var memorySystemTakenData = ParseArrayIntoTimeSeriesData(systemReservedMemoryArray);

  var yMax = 1;
  memoryTotal.map(function (item) {
    if (Number(item[1]) > yMax) {
      yMax = Number(item[1]);
    }
    return yMax;
  });

  var series = []
  if (memoryRequestedData.length !== 0) {
    series.push({
      title: "System reserved Memory",
      type: "area",
      data: memorySystemTakenData,
      color: "#FFFF00",
      valueFormatter: function o(e) {
        return keepTwoDecimal(e) + " C";
      },
    })
  }
  if (memoryRequestedData.length !== 0) {
    series.push({
      title: "Requested Memory",
      type: "area",
      data: memoryRequestedData,
      color: "#FF0000",
      valueFormatter: function o(e) {
        return keepTwoDecimal(e) + " C";
      },
    })
  }
  if (memoryAvaiableData.length !== 0) {
    series.push({
      title: "Available Memory",
      type: "area",
      color: "#0000CD",
      data: memoryAvaiableData,
      valueFormatter: function o(e) {
        return keepTwoDecimal(e) + " C";
      },
    })
  };

  return (
    <Container
      className="custom-dashboard-container"
      header={<Header variant="h2">Memory utilization</Header>}
    >
      <AreaChart
        series={series}
        xDomain={[memoryAvaiableData[0].x, memoryAvaiableData[memoryAvaiableData.length - 1].x]}
        yDomain={[0, yMax * 1.2]}
        i18nStrings={{
          filterLabel: "Filter displayed data",
          filterPlaceholder: "Filter data",
          filterSelectedAriaLabel: "selected",
          detailPopoverDismissAriaLabel: "Dismiss",
          legendAriaLabel: "Legend",
          chartAriaRoleDescription: "line chart",
          detailTotalLabel: "Total",
          xTickFormatter: (e) =>
            e
              .toLocaleDateString("en-US", {
                month: "short",
                day: "numeric",
                hour: "numeric",
                minute: "numeric",
                hour12: !1,
              })
              .split(",")
              .join("\n"),
          yTickFormatter: function o(e) {
            return keepTwoDecimal(e) + " G"
          },
        }}
        ariaLabel="Stacked area chart"
        errorText="Error loading data."
        height={300}
        hideFilter
        loadingText="Loading chart"
        recoveryText="Retry"
        xScaleType="time"
        xTitle="Time (UTC)"
        yTitle="Memory (GiB)"
        empty={
          <Box textAlign="center" color="inherit">
            <b>No data available</b>
            <Box variant="p" color="inherit">
              There is no data available
            </Box>
          </Box>
        }
        noMatch={
          <Box textAlign="center" color="inherit">
            <b>No matching data</b>
            <Box variant="p" color="inherit">
              There is no matching data to display
            </Box>
            <Button>Clear filter</Button>
          </Box>
        }
      />
    </Container>
  );
}
