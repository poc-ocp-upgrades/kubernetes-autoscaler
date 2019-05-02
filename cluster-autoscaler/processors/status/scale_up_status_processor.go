package status

import (
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/autoscaler/cluster-autoscaler/context"
 "k8s.io/autoscaler/cluster-autoscaler/processors/nodegroupset"
)

type ScaleUpStatus struct {
 Result                  ScaleUpResult
 ScaleUpInfos            []nodegroupset.ScaleUpInfo
 PodsTriggeredScaleUp    []*apiv1.Pod
 PodsRemainUnschedulable []NoScaleUpInfo
 PodsAwaitEvaluation     []*apiv1.Pod
}
type NoScaleUpInfo struct {
 Pod                *apiv1.Pod
 RejectedNodeGroups map[string]Reasons
 SkippedNodeGroups  map[string]Reasons
}
type ScaleUpResult int

const (
 ScaleUpSuccessful ScaleUpResult = iota
 ScaleUpError
 ScaleUpNoOptionsAvailable
 ScaleUpNotNeeded
 ScaleUpNotTried
 ScaleUpInCooldown
)

func (s *ScaleUpStatus) WasSuccessful() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return s.Result == ScaleUpSuccessful
}

type Reasons interface{ Reasons() []string }
type ScaleUpStatusProcessor interface {
 Process(context *context.AutoscalingContext, status *ScaleUpStatus)
 CleanUp()
}

func NewDefaultScaleUpStatusProcessor() ScaleUpStatusProcessor {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &EventingScaleUpStatusProcessor{}
}

type NoOpScaleUpStatusProcessor struct{}

func (p *NoOpScaleUpStatusProcessor) Process(context *context.AutoscalingContext, status *ScaleUpStatus) {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
func (p *NoOpScaleUpStatusProcessor) CleanUp() {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
