export class Memory {
  constructor(
    totalMem,
    requestedMem,
    availableMem,
    MemUsage,
    totalMemoryArray,
    systemReservedMemoryArray,
    requestedMemoryArray,
    availableMemoryArray,
  ) {
    this.totalMem = totalMem;
    this.requestedMem = requestedMem;
    this.availableMem = availableMem;
    this.MemUsage = MemUsage;
    this.totalMemoryArray = totalMemoryArray;
    this.systemReservedMemoryArray = systemReservedMemoryArray;
    this.requestedMemoryArray = requestedMemoryArray;
    this.availableMemoryArray = availableMemoryArray;
  }
}
