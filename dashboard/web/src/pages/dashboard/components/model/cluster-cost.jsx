export class ClusterCostInfo {
    constructor(
      clusterName,
      dataStep,
      clusterCostCurrent,
      clusterMonthEstimateCost,
      clusterAvgDailyCost,
      clusterAvgHourlyCoreCost,
      clusterComputeCostArray,
    ) {
      this.clusterName = clusterName ? clusterName : "-";
      this.dataStep = dataStep ? dataStep : 3600;
      this.clusterCostCurrent = clusterCostCurrent ? clusterCostCurrent : 0;
      this.clusterMonthEstimateCost = clusterMonthEstimateCost ? clusterMonthEstimateCost : 0;
      this.clusterAvgDailyCost = clusterAvgDailyCost ? clusterAvgDailyCost : 0;
      this.clusterAvgHourlyCoreCost = clusterAvgHourlyCoreCost ? clusterAvgHourlyCoreCost : 0;
      this.clusterComputeCostArray = clusterComputeCostArray ? clusterComputeCostArray : [];
    }
  }
