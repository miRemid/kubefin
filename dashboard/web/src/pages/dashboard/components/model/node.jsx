export class Node {
  constructor(nodeNum, ondemandNodes, spotNodes, fallbackNodes) {
    this.nodeNum = nodeNum;
    this.ondemandNodes = ondemandNodes ? ondemandNodes : 0;
    this.spotNodes = spotNodes ? spotNodes : 0;
    this.fallbackNodes = fallbackNodes ? fallbackNodes : 0;
  }
}
