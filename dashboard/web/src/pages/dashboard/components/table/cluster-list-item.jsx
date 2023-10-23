export class ClusterListItem {
    constructor(
      clusterId,
      clusterName,
      clusterRegion,
      computeCost,
      state,
      nodes,
      cpu,
      memory
    ) {
      this.clusterId = clusterId;
      this.clusterName = clusterName;
      this.clusterRegion = clusterRegion;
      this.computeCost = computeCost;
      this.state = state;
      this.nodes = nodes;
      this.cpu = cpu;
      this.memory = memory;
      this.href = '/dashboard/' + clusterId;
    }
  }
