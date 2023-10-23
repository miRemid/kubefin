export class Cluster {
    constructor(
      clusterId,
      clusterName,
      clusterRegion,
      cloudProvider,
      Pod,
      Node,
      Memory,
      CPU,
      state
    ) {
      this.clusterId = clusterId;
      this.clusterName = clusterName;
      this.clusterRegion = clusterRegion;
      this.cloudProvider = cloudProvider;
      this.Pod = Pod;
      this.Node = Node;
      this.Memory = Memory;
      this.CPU = CPU;
      this.state = state;
    }
  }
