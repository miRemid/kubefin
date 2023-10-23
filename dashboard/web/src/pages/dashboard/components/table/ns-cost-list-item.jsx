export class NSCostListItem {
    constructor(
      namespace,
      pods,
      cpuRequested,
      ramRequested,
      totalCost
    ) {
      this.namespace = namespace;
      this.pods = pods;
      this.cpuRequested = cpuRequested;
      this.ramRequested = ramRequested;
      this.totalCost = totalCost;
    }
  }
