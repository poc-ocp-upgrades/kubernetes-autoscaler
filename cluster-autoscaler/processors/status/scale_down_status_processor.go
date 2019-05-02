package status

import (
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/context"
 "k8s.io/autoscaler/cluster-autoscaler/simulator"
)

type ScaleDownStatus struct {
 Result            ScaleDownResult
 ScaledDownNodes   []*ScaleDownNode
 NodeDeleteResults map[string]error
}
type ScaleDownNode struct {
 Node        *apiv1.Node
 NodeGroup   cloudprovider.NodeGroup
 EvictedPods []*apiv1.Pod
 UtilInfo    simulator.UtilizationInfo
}
type ScaleDownResult int

const (
 ScaleDownError ScaleDownResult = iota
 ScaleDownNoUnneeded
 ScaleDownNoNodeDeleted
 ScaleDownNodeDeleted
 ScaleDownNodeDeleteStarted
 ScaleDownNotTried
 ScaleDownInCooldown
 ScaleDownInProgress
)

type ScaleDownStatusProcessor interface {
 Process(context *context.AutoscalingContext, status *ScaleDownStatus)
 CleanUp()
}

func NewDefaultScaleDownStatusProcessor() ScaleDownStatusProcessor {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &NoOpScaleDownStatusProcessor{}
}

type NoOpScaleDownStatusProcessor struct{}

func (p *NoOpScaleDownStatusProcessor) Process(context *context.AutoscalingContext, status *ScaleDownStatus) {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
func (p *NoOpScaleDownStatusProcessor) CleanUp() {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
