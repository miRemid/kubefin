export class ClusterComputeCost {
    constructor(
      timestamp,
      totalCost,
      costFallbackBillingMode,
      costOnDemandBillingMode,
      costSpotBillingMode,
      cpuCoreCount,
      cpuCoreUsage,
      cpuCost,
      ramGBCount,
      ramUsage,
      ramCost,
    ) {
      this.timestamp = timestamp;
      this.totalCost = totalCost;
      this.costFallbackBillingMode = costFallbackBillingMode;
      this.costOnDemandBillingMode = costOnDemandBillingMode;
      this.costSpotBillingMode = costSpotBillingMode;
      this.cpuCoreCount = cpuCoreCount;
      this.cpuCoreUsage = cpuCoreUsage;
      this.cpuCost = cpuCost;
      this.ramGBCount = ramGBCount;
      this.ramUsage = ramUsage;
      this.ramCost = ramCost;
    }
  }
