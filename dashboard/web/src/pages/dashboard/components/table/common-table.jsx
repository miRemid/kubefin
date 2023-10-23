// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from 'react';
import { useCollection } from '@cloudscape-design/collection-hooks';
import {
  Table,
  Pagination,
} from "@cloudscape-design/components";

import { ClusterComputeSpendTableHeader } from './common-components';
import { paginationLabels, distributionSelectionLabels } from '../../../../common/labels';

import '../../../../styles/base.scss';

export function CommonTable({
  data,
  columnDefinitions,
  saveWidths,
}) {
  const { items, collectionProps, paginationProps } = useCollection(
    data,
    {
      pagination: { pageSize: 20 },
      sorting: { defaultState: { sortingColumn: columnDefinitions[0] } },
      selection: {},
    }
  );

  return (
    <Table
      {...collectionProps}
      items={items}
      columnDefinitions={columnDefinitions}
      ariaLabels={distributionSelectionLabels}
      stickyHeader={true}
      resizableColumns={true}
      onColumnWidthsChange={saveWidths}
      header={
        <ClusterComputeSpendTableHeader
          serverSide={false}
        />
      }
      loadingText="Loading"
      pagination={<Pagination {...paginationProps} ariaLabels={paginationLabels} />}
    />
  );
}
