import React, { useState } from "react";
import {
  AppLayout,
  Table,
  Pagination,
  TopNavigation,
  TextFilter,
  SideNavigation,
  Box,
  Popover,
} from "@cloudscape-design/components";
import { useParams } from "react-router";
import { useCollection } from "@cloudscape-design/collection-hooks";
import { appLayoutLabels } from "../../../../common/labels";
import { getFilterCounterText } from "../../../../common/tableCounterStrings";
import {
  TableEmptyState,
  TableNoMatchState,
} from "../../../commons/common-components";
import { paginationLabels } from "../../../../common/labels";
import { WorkloadCostToolsContent } from "../table/common-components";
import { SELECTION_LABELS, COLUMN_DEFINITIONS_WORKLOAD } from "./table-config";
import { useSplitPanel } from "./utils";
import { WorkloadCostHeader } from "./workload-cost-header";

// It's necessary to import a scss file or the build will fail
import "../../../../styles/base.scss";

export function WorkloadCostDashboard(props) {
  const { clusterId } = useParams();
  const { clusterName } = useParams();
  const [workloadCostInfo, setWorkloadCostInfo] = useState({});
  function handleChangeWorkloadCostInfo(workloadCostInfo) {
    console.log("workloadCostInfo=", workloadCostInfo);
    setWorkloadCostInfo(workloadCostInfo);
  }

  const workloadCostJsonArray = Array.from(workloadCostInfo).map((item) => ({
    namespace: item.namespace,
    workloadType: item.workloadType,
    workloadName: item.workloadName,
    totalCost: item.totalCost,
    podCount: item.podCount,
    cpuCoreRequest: item.cpuCoreRequest,
    cpuCoreUsage: item.cpuCoreUsage,
    ramGBRequest: item.ramGBRequest,
    ramGBUsage: item.ramGBUsage,
  }));

  const navHeader = { text: "Workload Cost", href: "#/" };
  const navItems = [
    {
      type: "section",
      text: "Dashboard",
      items: [
        {
          type: "link",
          text: "all clusters",
          href: "/dashboard",
        },
        {
          type: "link",
          text: clusterName,
          href: "/dashboard/" + clusterId,
          info: (
            <Box color="text-status-info" display="inline">
              <Popover
                header="cluster dashboard"
                size="medium"
                triggerType="text"
                content={<>Show the resources information in your cluster</>}
                renderWithPortal={true}
                dismissAriaLabel="Close"
              >
                <Box
                  color="text-status-info"
                  fontSize="body-s"
                  fontWeight="bold"
                  data-testid="new-feature-announcement-trigger"
                >
                  Info
                </Box>
              </Popover>
            </Box>
          ),
        },
      ],
    },
    {
      type: "section",
      text: "Cost allocation",
      items: [
        {
          type: "link",
          text: "cluster cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/cluster",
        },
        {
          type: "link",
          text: "workload cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/workload",
        },
        {
          type: "link",
          text: "namespace cost",
          href: "/cost/" + clusterId + "/" + clusterName + "/namespace",
        },
      ],
    },
  ];
  const {
    items,
    actions,
    filteredItemsCount,
    collectionProps,
    filterProps,
    paginationProps,
  } = useCollection(workloadCostJsonArray, {
    filtering: {
      empty: <TableEmptyState resourceName="Workload" />,
      noMatch: (
        <TableNoMatchState onClearFilter={() => actions.setFiltering("")} />
      ),
    },
    selection: {},
  });

  const [toolsOpen, setToolsOpen] = useState(false);
  const [navigationOpen, setNavigationOpen] = useState(true);

  const {
    splitPanelOpen,
    onSplitPanelToggle,
    splitPanelSize,
    onSplitPanelResize,
  } = useSplitPanel(collectionProps.selectedItems);
  const i18nStrings = {
    errorIconAriaLabel: "Error",
    navigationAriaLabel: "Steps",
    cancelButton: "Cancel",
    previousButton: "Previous",
    nextButton: "Next",
    submitButton: "I run the script",
    optional: "optional",
    searchIconAriaLabel: "Search",
    searchDismissIconAriaLabel: "Close search",
    overflowMenuTriggerText: "More",
    overflowMenuTitleText: "All",
    overflowMenuBackIconAriaLabel: "Back",
    overflowMenuDismissIconAriaLabel: "Close menu",
  };

  return (
    <>
      <TopNavigation
        i18nStrings={i18nStrings}
        identity={{
          href: "/",
          title: "KubeFin",
        }}
        utilities={[
          {
            type: "button",
            text: "Documentation",
            href: "https://kubefin.dev",
            external: true,
            externalIconAriaLabel: " (opens in a new tab)",
          },
        ]}
      />
      <AppLayout
        contentType="table"
        headerSelector="#header"
        navigationOpen={navigationOpen}
        navigation={<SideNavigation header={navHeader} items={navItems} />}
        onNavigationChange={({ detail }) => setNavigationOpen(detail.open)}
        tools={<WorkloadCostToolsContent />}
        toolsOpen={toolsOpen}
        onToolsChange={({ detail }) => setToolsOpen(detail.open)}
        splitPanelOpen={splitPanelOpen}
        onSplitPanelToggle={onSplitPanelToggle}
        splitPanelSize={splitPanelSize}
        onSplitPanelResize={onSplitPanelResize}
        content={
          <Table
            {...collectionProps}
            header={
              <WorkloadCostHeader
                updateTools={() => setToolsOpen(true)}
                onChangeWorkloadCostInfo={(workloadCostInfo) =>
                  handleChangeWorkloadCostInfo(workloadCostInfo)
                }
              />
            }
            variant="full-page"
            stickyHeader={true}
            workloadCostInfo={workloadCostInfo}
            columnDefinitions={COLUMN_DEFINITIONS_WORKLOAD}
            items={items}
            ariaLabels={SELECTION_LABELS}
            filter={
              <TextFilter
                {...filterProps}
                filteringAriaLabel="Filter workload"
                filteringPlaceholder="Search workload"
                countText={getFilterCounterText(filteredItemsCount)}
              />
            }
            pagination={
              <Pagination {...paginationProps} ariaLabels={paginationLabels} />
            }
          />
        }
        ariaLabels={appLayoutLabels}
      />
    </>
  );
}
