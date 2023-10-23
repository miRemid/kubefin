// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

import { keepThreeDecimal, keepTwoDecimal } from "../components-common";

// SPDX-License-Identifier: MIT-0
export const COLUMN_DEFINITIONS_WORKLOAD = [
  {
    id: "workloadName",
    header: "Workload name",
    cell: (item) => item.workloadName,
  },
  {
    id: "workloadType",
    header: "Type",
    cell: (item) => item.workloadType,
  },
  {
    id: "namespace",
    header: "Namespace",
    cell: (item) => item.namespace,
  },
  {
    id: "pods",
    header: "Pods",
    cell: (item) => keepTwoDecimal(item.podCount),
    sortingField: "pods",
  },
  {
    id: "cpu",
    header: "CPU requested",
    cell: (item) => keepThreeDecimal(item.cpuCoreRequest),
    sortingField: "cpu",
  },
  {
    id: "memory",
    header: "Memory requested",
    cell: (item) => keepThreeDecimal(item.ramGBRequest),
    sortingField: "memory",
  },
  {
    id: "totalCost",
    header: "Total cost",
    cell: (item) => keepThreeDecimal(item.totalCost),
    sortingField: "totalCost",
  }
];

export const COLUMN_DEFINITIONS_PANEL_CONTENT_SINGLE = [
  {
    id: "type",
    header: "Type",
    cell: (item) => item.type,
  },
  {
    id: "protocol",
    header: "Protocol",
    cell: (item) => item.protocol,
  },
  {
    id: "portRange",
    header: "Port range",
    cell: (item) => item.portRange,
  },
  {
    id: "source",
    header: "Source",
    cell: (item) => item.source,
  },
  {
    id: "description",
    header: "Description",
    cell: (item) => item.description,
  },
];

export const SELECTION_LABELS = {
  itemSelectionLabel: (data, row) => `select ${row.id}`,
  allItemsSelectionLabel: () => "select all",
  selectionGroupLabel: "Instance selection",
};

