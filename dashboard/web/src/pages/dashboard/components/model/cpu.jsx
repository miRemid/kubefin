export class CPU {
  constructor(totalCPUCores, requestedCPU,
    allocatableCPU, CPUusage, requestedCPUArray, allocatableCPUArray, totalCPUArray) {
    this.totalCPUCores = totalCPUCores;
    this.requestedCPU = requestedCPU;
    this.allocatableCPU = allocatableCPU;
    this.CPUusage = CPUusage;
    this.requestedCPUArray = requestedCPUArray;
    this.allocatableCPUArray = allocatableCPUArray;
    this.totalCPUArray = totalCPUArray;
  }
}
