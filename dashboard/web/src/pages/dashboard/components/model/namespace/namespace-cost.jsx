export class NamespaceCostInfo {
    constructor(
      namespace,
      podCount,
      cpuRequest,
      ramGBRequest,
      totalCost
    ) {
      this.namespace = namespace;
      this.podCount = podCount;
      this.cpuRequest = cpuRequest;
      this.ramGBRequest = ramGBRequest;
      this.totalCost=totalCost;
  }
}
