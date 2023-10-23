// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import { useColumnWidths } from "../../../commons/use-column-widths";
import {
  CLUSTERS_COLUMN_DEFINITIONS,
  FILTERING_PROPERTIES,
} from "../table/table-property-filter-config";
import { PropertyFilterTable } from "../table/property-filter-table";
import { ClusterListItem } from "../table/cluster-list-item";
import "../../../../styles/base.scss";

export default function ClustersTableList(props) {
  const [columnDefinitions, saveWidths] = useColumnWidths(
    "React-TableServerSide-Widths",
    CLUSTERS_COLUMN_DEFINITIONS
  );
  const costSummary = props.costSummary;
  const metricsSummary = props.metricsSummary;
  const data = costSummary.map((item) => {
    const clusterItem = new ClusterListItem(
      item.clusterId,
      item.clusterName,
      item.clusterRegion,
      item.clusterMonthCostCurrent,
      item.clusterConnectionSate === "running" ? "Active" : "Deactivated"
    );
    return clusterItem;
  });
  data.forEach((costItem) => {
    metricsSummary.forEach((metricsItem) => {
      if (metricsItem.clusterId === costItem.clusterId) {
        costItem.nodes = metricsItem.nodeNumbersCurrent;
        costItem.cpu = metricsItem.cpuCoreTotal;
        costItem.memory = metricsItem.ramGBTotal;
      }
    });
  });

  return (
    <PropertyFilterTable
      data={data}
      columnDefinitions={columnDefinitions}
      saveWidths={saveWidths}
      filteringProperties={FILTERING_PROPERTIES}
    />
  );
}
