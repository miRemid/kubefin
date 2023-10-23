// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import { useColumnWidths } from "../../../commons/use-column-widths";
import {
  NS_COST_COLUMN_DEFINITIONS,
  NS_COST_FILTERING_PROPERTIES,
} from "./table-property-filter-config";
import { NSCostListItem } from "./ns-cost-list-item";
import "../../../../styles/base.scss";
import { NSPropertyFilterTable } from "./ns-property-filter-table";

export default function NSCostTableList(props) {
  const [columnDefinitions, saveWidths] = useColumnWidths(
    "React-TableServerSide-Widths",
    NS_COST_COLUMN_DEFINITIONS
  );
  const namespaceCostMap = props.namespaceCostMap;
  if (!(namespaceCostMap instanceof Map) || namespaceCostMap.size === 0) {
    return <div>The namespace Cost is empty or not a valid Map.</div>;
  }
  // console.log("namespaceCostMap@ClusterNsCostTableList = ", namespaceCostMap);

  const nsCostList = Array.from(namespaceCostMap).map(
    ([key, value]) => {
      return new NSCostListItem(
        key,
        value.podCount.toFixed(2),
        value.cpuRequest.toFixed(2),
        value.ramGBRequest.toFixed(2),
        value.totalCost.toFixed(2)
      );
    }
  );

  return (
    <NSPropertyFilterTable
      data={nsCostList}
      columnDefinitions={columnDefinitions}
      saveWidths={saveWidths}
      filteringProperties={NS_COST_FILTERING_PROPERTIES}
    />
  );
}
