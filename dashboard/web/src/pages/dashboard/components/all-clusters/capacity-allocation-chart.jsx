// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import {
  Container,
  Header,
  PieChart,
  Button,
  Box,
  ColumnLayout,
  HelpPanel,
  Icon,
} from "@cloudscape-design/components";
import { keepTwoDecimal, percentageFormatter } from "../components-common";
import { ExternalLinkItem, InfoLink } from "../../../commons/common-components";

function CapacityAllocationInfo() {
  return (
    <HelpPanel
      header={<h2>CPU and memory allocation</h2>}
      footer={
        <>
          <h3>
            Learn More{" "}
            <span role="img" aria-label="Icon external Link">
              <Icon name="external" />
            </span>
          </h3>
          <ul>
            <li>
              <ExternalLinkItem
                href="https://kubernetes.io/docs/concepts/architecture/nodes/#capacity"
                text="Kubernetes node capacity and allocable"
              />
            </li>
          </ul>
        </>
      }
    >
      <p>
        This section shows the CPU and Memory allocation in all your K8s clusters.
      </p>
      <ul>
        <li>System reserved: kubelet reserved.</li>
        <li>Workload reserved: all pods requested</li>
      </ul>
    </HelpPanel>
  );
}
const i18nStrings = {
  detailsValue: "Value",
  detailsPercentage: "Percentage",
  filterLabel: "Filter displayed data",
  filterPlaceholder: "Filter data",
  filterSelectedAriaLabel: "selected",
  detailPopoverDismissAriaLabel: "Dismiss",
  legendAriaLabel: "Legend",
  chartAriaRoleDescription: "pie chart",
  segmentAriaRoleDescription: "segment",
};
export default function CapacityAllocation(props) {
  const totalCPUCores = props.totalCPUCores;
  const systemReservedCPU = props.systemReservedCPU;
  const workloadReservedCPU = props.workloadReservedCPU;
  const allocatableCPU = props.allocatableCPU;
  const totalMemory = props.totalMemory;
  const systemReservedMemory = props.systemReservedMemory;
  const workloadReservedMemory = props.workloadReservedMemory;
  const allocatableMemory = props.allocatableMemory;

  return (
    <Container
      className="custom-dashboard-container"
      header={
        <Header
          variant="h2"
          info={
            <InfoLink
              onFollow={() => props.updateTools(<CapacityAllocationInfo />)}
              ariaLabel={"Kubernetes node capacity and allocable"}
            />
          }
        >
          CPU and memory allocation
        </Header>
      }
    >
      <ColumnLayout columns="2">
        <PieChart
          size="large"
          data={[
            { title: "Workload reserved", value: keepTwoDecimal(workloadReservedCPU) },
            { title: "System reserved", value: keepTwoDecimal(systemReservedCPU) },
            { title: "Available", value: keepTwoDecimal(allocatableCPU) },
          ]}
          ariaLabel="CPU allocation chart"
          ariaDescription="CPU allocation chart."
          hideFilter={true}
          segmentDescription={(datum, sum) =>
            `${datum.value} Core, ${percentageFormatter(datum.value / sum)}`
          }
          i18nStrings={i18nStrings}
          errorText="Error loading data."
          hideLegend={true}
          innerMetricDescription="Core"
          innerMetricValue={keepTwoDecimal(totalCPUCores)}
          loadingText="Loading chart"
          recoveryText="Retry"
          variant="donut"
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
        <PieChart
          size="large"
          data={[
            { title: "Workload reserved", value: keepTwoDecimal(workloadReservedMemory) },
            { title: "System reserved", value: keepTwoDecimal(systemReservedMemory) },
            { title: "Available", value: keepTwoDecimal(allocatableMemory) },
          ]}
          ariaLabel="Memory allocation chart"
          ariaDescription="Memory allocation chart."
          hideFilter={true}
          segmentDescription={(datum, sum) =>
            `${datum.value} GiB, ${percentageFormatter(datum.value / sum)}`
          }
          i18nStrings={i18nStrings}
          errorText="Error loading data."
          hideLegend={true}
          innerMetricDescription="GiB"
          innerMetricValue={keepTwoDecimal(totalMemory)}
          loadingText="Loading chart"
          recoveryText="Retry"
          variant="donut"
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
      </ColumnLayout>
    </Container>
  );
}
