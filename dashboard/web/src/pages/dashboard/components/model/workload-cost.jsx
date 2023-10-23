export class WorkloadCost {
    constructor(
      namespace,
      workloadType,
      workloadName,
      totalCost,
      podCount,
      cpuCoreRequest,
      cpuCoreUsage,
      ramGBRequest,
      ramGBUsage,
    ) {
      this.namespace = namespace ? namespace : "-";
      this.workloadType = workloadType ? workloadType : "-";
      this.workloadName = workloadName ? workloadName : "-";
      this.totalCost = totalCost ? totalCost : 0;
      this.podCount = podCount ? podCount : 0;
      this.cpuCoreRequest = cpuCoreRequest ? cpuCoreRequest : 0;
      this.cpuCoreUsage = cpuCoreUsage ? cpuCoreUsage : 0;
      this.ramGBRequest = ramGBRequest ? ramGBRequest : 0;
      this.ramGBUsage = ramGBUsage ? ramGBUsage : 0;
    }
  }
