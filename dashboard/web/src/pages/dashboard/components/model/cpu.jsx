export class CPU {
  constructor(
    totalCPUCores,
    requestedCPU,
    availableCPU,
    CPUUsage,
    totalCPUArray,
    systemReservedCPUArray,
    requestedCPUArray,
    availableCPUArray) {
    this.totalCPUCores = totalCPUCores;
    this.requestedCPU = requestedCPU;
    this.availableCPU = availableCPU;
    this.CPUUsage = CPUUsage;
    this.totalCPUArray = totalCPUArray;
    this.systemReservedCPUArray = systemReservedCPUArray;
    this.systemReservedCPUArray = systemReservedCPUArray;
    this.requestedCPUArray = requestedCPUArray;
    this.availableCPUArray = availableCPUArray;
  }
}
