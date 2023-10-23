// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from 'react';
import { Box } from '@cloudscape-design/components';

export const percentageFormatter = value => `${(value * 100).toFixed(1)}%`;

export const numberTickFormatter = value => {
  if (Math.abs(value) < 1000) {
    return value.toString();
  }
  return (value / 1000).toFixed() + 'k';
};

export const dateHourFormatter = date =>
  date
    .toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      hour12: true,
    })
    .split(',')
    .join('\n');

export const dateDayFormatter = date =>
  date
    .toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour12: false,
    })
    .split(' ')
    .join('\n');

export const commonChartProps = {
  loadingText: 'Loading chart',
  errorText: 'Error loading data.',
  recoveryText: 'Retry',
  empty: (
    <Box textAlign="center" color="inherit">
      <b>No data available</b>
      <Box variant="p" color="inherit">
        There is no data available
      </Box>
    </Box>
  ),
  noMatch: (
    <Box textAlign="center" color="inherit">
      <b>No matching data</b>
      <Box variant="p" color="inherit">
        There is no matching data to display
      </Box>
    </Box>
  ),
  i18nStrings: {
    filterLabel: 'Filter displayed data',
    filterPlaceholder: 'Filter data',
    filterSelectedAriaLabel: 'selected',
    legendAriaLabel: 'Legend',
    chartAriaRoleDescription: 'line chart',
    xAxisAriaRoleDescription: 'x axis',
    yAxisAriaRoleDescription: 'y axis',
    yTickFormatter: numberTickFormatter,
  },
};

export const lineChartInstructions =
  'Use up/down arrow keys to navigate between series, and left/right arrow keys to navigate within a series.';

export const barChartInstructions = 'Use left/right arrow keys to navigate between data groups.';

export function keepTwoDecimal(n) {
  if (!n) {
    return 0
  }
  return Math.floor(n * 100) / 100;
}

export function keepThreeDecimal(n) {
  if (!n) {
    return 0
  }
  return Math.floor(n * 1000) / 1000;
}
