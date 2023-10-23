export class Pod {
  constructor(podNum, scheduledPodNum, unscheduledPodNum, scheduledRatio) {
    this.podNum = podNum ? podNum : 0;
    this.scheduledPodNum = scheduledPodNum ? scheduledPodNum : 0;
    this.unscheduledPodNum = unscheduledPodNum ? unscheduledPodNum : 0;
    this.scheduledRatio = scheduledRatio ? scheduledRatio * 100 : 0;
  }
}
