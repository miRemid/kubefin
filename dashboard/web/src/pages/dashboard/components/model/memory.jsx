export class Memory {
  constructor(
    totalMEM,
    requestedMEM,
    allocatableMEM,
    MEMusage,
    requestedMemoryArray,
    allocatableMemoryArray,
    totalMemoryArray
  ) {
    this.totalMEM = totalMEM;
    this.requestedMEM = requestedMEM;
    this.allocatableMEM = allocatableMEM;
    this.MEMusage = MEMusage;
    this.requestedMemoryArray = requestedMemoryArray;
    this.allocatableMemoryArray = allocatableMemoryArray;
    this.totalMemoryArray = totalMemoryArray;
  }
}
