// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import { useColumnWidths } from "../../../commons/use-column-widths";
import { CLUSTER_COMPUTE_COST_COLUMN_DEFINITIONS } from "./table-property-filter-config";
import { CommonTable } from "./common-table";
import { ClusterComputeCostListItem } from "./cluster-compute-cost-list-item";
import { ClusterCostInfo } from "../model/cluster-cost";
import { keepTwoDecimal } from "../components-common";

import "../../../../styles/base.scss";

export default function ClusterComputeCostList(props) {
  const [columnDefinitions, saveWidths] = useColumnWidths(
    "React-TableServerSide-Widths",
    CLUSTER_COMPUTE_COST_COLUMN_DEFINITIONS
  );
  let clusterCostInfo = props.clusterCostInfo;
  if (clusterCostInfo === undefined) {
    clusterCostInfo = new ClusterCostInfo();
  }
  const clusterComputeCostArrayByTimestamp =
    clusterCostInfo.clusterComputeCostArray;
  const costMap = buildCostMap(clusterComputeCostArrayByTimestamp);
  console.log("costMap = ", costMap);

  const costArray = Array.from(costMap);
  console.log("costArray = ", costArray);

  const data = costArray.map((item) => {
    const clusterComputeCostListItem = new ClusterComputeCostListItem(
      item[0],
      keepTwoDecimal(item[1].totalCost),
      keepTwoDecimal(item[1].costFallbackBillingMode),
      keepTwoDecimal(item[1].costOnDemandBillingMode),
      keepTwoDecimal(item[1].costSpotBillingMode),
      keepTwoDecimal(item[1].cpuCoreCount),
      item[1].cpuCoreCountIndex,
      keepTwoDecimal(item[1].cpuCoreUsage),
      keepTwoDecimal(item[1].cpuCost),
      keepTwoDecimal(item[1].ramGBCount),
      item[1].ramGBCountIndex,
      keepTwoDecimal(item[1].ramGBUsage),
      keepTwoDecimal(item[1].ramCost),
    );
    return clusterComputeCostListItem;
  });

  function buildCostMap(data) {
    const costMap = new Map();

    data.forEach((item) => {
      // Calculate the cost sum for each item
      const clusterComputeCostListItem = new ClusterComputeCostListItem(
        item.timestamp,
        keepTwoDecimal(item.totalCost),
        keepTwoDecimal(item.costFallbackBillingMode),
        keepTwoDecimal(item.costOnDemandBillingMode),
        keepTwoDecimal(item.costSpotBillingMode),
        keepTwoDecimal(item.cpuCoreCount),
        1.0,
        keepTwoDecimal(item.cpuCoreUsage),
        keepTwoDecimal(item.cpuCost),
        keepTwoDecimal(item.ramGBCount),
        1.0,
        keepTwoDecimal(item.ramUsage),
        keepTwoDecimal(item.ramCost),
      );

      // Check if there's an entry in the map for the item's timestamp day
      const itemTimestamp = clusterComputeCostListItem.timestamp;
      const itemDate = new Date(itemTimestamp * 1000);
      console.log("itemTimestamp=",itemTimestamp);
      console.log("itemDate=",itemDate);

      const itemDay = new Date(
        itemDate.getFullYear(),
        itemDate.getMonth(),
        itemDate.getDate()
      ).getTime();

      if (costMap.has(itemDay)) {
        // If there's an entry, add the cost sum to the existing value
        var currentClusterComputeCostListItem = costMap.get(itemDay);
        currentClusterComputeCostListItem.totalCost += clusterComputeCostListItem.totalCost;
        currentClusterComputeCostListItem.costFallbackBillingMode +=
          clusterComputeCostListItem.costFallbackBillingMode;
        currentClusterComputeCostListItem.costOnDemandBillingMode +=
          clusterComputeCostListItem.costOnDemandBillingMode;
        currentClusterComputeCostListItem.costSpotBillingMode +=
          clusterComputeCostListItem.costSpotBillingMode;

        currentClusterComputeCostListItem.cpuCost += clusterComputeCostListItem.cpuCost;
        currentClusterComputeCostListItem.cpuCoreCount += clusterComputeCostListItem.cpuCoreCount;
        currentClusterComputeCostListItem.cpuCoreCountIndex += 1.0;
        currentClusterComputeCostListItem.ramCost += clusterComputeCostListItem.ramCost;
        currentClusterComputeCostListItem.ramGBCount += clusterComputeCostListItem.ramGBCount;
        currentClusterComputeCostListItem.ramGBCountIndex += 1.0;
        costMap.set(itemDay, currentClusterComputeCostListItem);
      } else {
        // If no entry exists, create a new entry with the cost sum
        costMap.set(itemDay, clusterComputeCostListItem);
      }
    });

    return costMap;
  }

  return (
    <CommonTable
      data={data}
      columnDefinitions={columnDefinitions}
      saveWidths={saveWidths}
    />
  );
}
