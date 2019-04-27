package status

import (
	"fmt"
	"strings"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/context"
)

type EventingScaleUpStatusProcessor struct{}

func (p *EventingScaleUpStatusProcessor) Process(context *context.AutoscalingContext, status *ScaleUpStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, noScaleUpInfo := range status.PodsRemainUnschedulable {
		context.Recorder.Event(noScaleUpInfo.Pod, apiv1.EventTypeNormal, "NotTriggerScaleUp", fmt.Sprintf("pod didn't trigger scale-up (it wouldn't fit if a new node is added): %s", ReasonsMessage(noScaleUpInfo)))
	}
	if len(status.ScaleUpInfos) > 0 {
		for _, pod := range status.PodsTriggeredScaleUp {
			context.Recorder.Eventf(pod, apiv1.EventTypeNormal, "TriggeredScaleUp", "pod triggered scale-up: %v", status.ScaleUpInfos)
		}
	}
}
func (p *EventingScaleUpStatusProcessor) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func ReasonsMessage(noScaleUpInfo NoScaleUpInfo) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	messages := []string{}
	aggregated := map[string]int{}
	for _, reasons := range noScaleUpInfo.RejectedNodeGroups {
		for _, reason := range reasons.Reasons() {
			aggregated[reason]++
		}
	}
	for _, reasons := range noScaleUpInfo.SkippedNodeGroups {
		for _, reason := range reasons.Reasons() {
			aggregated[reason]++
		}
	}
	for msg, count := range aggregated {
		messages = append(messages, fmt.Sprintf("%d %s", count, msg))
	}
	return strings.Join(messages, ", ")
}
