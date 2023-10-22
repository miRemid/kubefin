// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import { StatusIndicator, Link } from "@cloudscape-design/components";
import { addColumnSortLabels } from "../../../../common/labels";
import { keepTwoDecimal } from "../components-common";

export const CLUSTERS_COLUMN_DEFINITIONS = addColumnSortLabels([
  {
    id: "clusterId",
    sortingField: "clusterId",
    header: "cluster id",
    cell: (item) => (
      <div>
        <Link href={item.href}>{item.clusterId}</Link>
      </div>
    ),
    minWidth: 180,
  },
  {
    id: "clusterName",
    sortingField: "clusterName",
    cell: (item) => item.clusterName,
    header: "cluster name",
    minWidth: 160,
  },
  {
    id: "clusterRegion",
    sortingField: "clusterRegion",
    header: "region",
    cell: (item) => item.clusterRegion,
    minWidth: 100,
  },
  {
    id: "nodes",
    sortingField: "nodes",
    header: "total nodes",
    cell: (item) => keepTwoDecimal(item.nodes),
    minWidth: 100,
  },
  {
    id: "cpu",
    sortingField: "cpu",
    header: "CPU(vCPU)",
    cell: (item) => keepTwoDecimal(item.cpu),
    minWidth: 100,
  },
  {
    id: "memory",
    sortingField: "memory",
    header: "Memory(GiB)",
    cell: (item) => keepTwoDecimal(item.memory),
    minWidth: 100,
  },
  {
    id: "computeCost",
    sortingField: "computeCost",
    header: "total cost",
    cell: (item) => keepTwoDecimal(item.computeCost),
    minWidth: 100,
  },
  {
    id: "state",
    sortingField: "state",
    header: "state",
    cell: (item) => (
      <StatusIndicator
        type={item.state === "Deactivated" ? "error" : "success"}
      >
        {item.state}
      </StatusIndicator>
    ),
    minWidth: 100,
  },
]);
export const CLUSTER_COMPUTE_COST_COLUMN_DEFINITIONS = addColumnSortLabels([
  {
    id: "date",
    sortingField: "date",
    header: "date",
    cell: (item) =>
      new Intl.DateTimeFormat("en-US", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
      }).format(item.timestamp),
    minWidth: 180,
  },
  {
    id: "provisionedCPU",
    sortingField: "provisionedCPU",
    cell: (item) =>
      keepTwoDecimal(item.cpuCoreCount / item.cpuCoreCountIndex),
    header: "Provisioned CPU",
    minWidth: 160,
  },
  {
    id: "provisionedMEM",
    sortingField: "provisionedMEM",
    header: "Provisioned Memory",
    cell: (item) =>
      keepTwoDecimal(item.ramGBCount / item.ramGBCountIndex),
    minWidth: 100,
  },
  {
    id: "costCPU",
    sortingField: "costCPU",
    header: "CPU cost",
    cell: (item) =>
      keepTwoDecimal(item.cpuCost),
    minWidth: 100,
  },
  {
    id: "costMEM",
    sortingField: "costMEM",
    header: "Memory cost",
    cell: (item) =>
      keepTwoDecimal(item.ramCost),
    minWidth: 100,
  },
  {
    id: "totalCost",
    sortingField: "totalCost",
    header: "total cost",
    cell: (item) =>
      item.totalCost,
    minWidth: 100,
  },
]);
export const NS_COST_COLUMN_DEFINITIONS = addColumnSortLabels([
  {
    id: "namespace",
    sortingField: "namespace",
    header: "namespace",
    cell: (item) => item.namespace,
    minWidth: 180,
  },
  {
    id: "pods",
    sortingField: "pods",
    cell: (item) => item.pods,
    header: "pods",
    minWidth: 160,
  },
  {
    id: "cpuRequested",
    sortingField: "cpuRequested",
    header: "CPU Requested",
    cell: (item) => item.cpuRequested,
    minWidth: 100,
  },
  {
    id: "ramRequested",
    sortingField: "ramRequested",
    header: "MEM Requested",
    cell: (item) => item.ramRequested,
    minWidth: 100,
  },
  {
    id: "totalCost",
    sortingField: "totalCost",
    header: "total cost",
    cell: (item) => item.totalCost,
    minWidth: 100,
  },
]);
export const PAGE_SIZE_OPTIONS = [
  { value: 10, label: "10 Distributions" },
  { value: 30, label: "30 Distributions" },
  { value: 50, label: "50 Distributions" },
];

export const FILTERING_PROPERTIES = [
  {
    propertyLabel: "Cluster ID",
    key: "clusterId",
    groupValuesLabel: "Cluster ID values",
    operators: [":", "!:", "=", "!="],
  },
  {
    propertyLabel: "Cluster name",
    key: "clusterName",
    groupValuesLabel: "Cluster name values",
    operators: [":", "!:", "=", "!="],
  },
  {
    propertyLabel: "clusterRegion",
    key: "clusterRegion",
    groupValuesLabel: "Region values",
    operators: [":", "!:", "=", "!="],
  },
  {
    propertyLabel: "State",
    key: "state",
    groupValuesLabel: "State values",
    operators: [":", "!:", "=", "!="],
  },
];
export const NS_COST_FILTERING_PROPERTIES = [
  {
    propertyLabel: "namespace",
    key: "namespace",
    groupValuesLabel: "namespace values",
    operators: [":", "!:", "=", "!="],
  },
];
export const PROPERTY_FILTERING_I18N_CONSTANTS = {
  filteringAriaLabel: "your choice",
  dismissAriaLabel: "Dismiss",

  filteringPlaceholder: "Search",
  groupValuesText: "Values",
  groupPropertiesText: "Properties",
  operatorsText: "Operators",

  operationAndText: "and",
  operationOrText: "or",

  operatorLessText: "Less than",
  operatorLessOrEqualText: "Less than or equal",
  operatorGreaterText: "Greater than",
  operatorGreaterOrEqualText: "Greater than or equal",
  operatorContainsText: "Contains",
  operatorDoesNotContainText: "Does not contain",
  operatorEqualsText: "Equals",
  operatorDoesNotEqualText: "Does not equal",

  editTokenHeader: "Edit filter",
  propertyText: "Property",
  operatorText: "Operator",
  valueText: "Value",
  cancelActionText: "Cancel",
  applyActionText: "Apply",
  allPropertiesLabel: "All properties",

  tokenLimitShowMore: "Show more",
  tokenLimitShowFewer: "Show fewer",
  clearFiltersText: "Clear filters",
  removeTokenButtonAriaLabel: () => "Remove token",
  enteredTextLabel: (text) => `Use: "${text}"`,
};
