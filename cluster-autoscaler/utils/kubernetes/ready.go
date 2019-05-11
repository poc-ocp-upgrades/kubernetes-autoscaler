package kubernetes

import (
	"fmt"
	"time"
	apiv1 "k8s.io/api/core/v1"
)

func IsNodeReadyAndSchedulable(node *apiv1.Node) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ready, _, _ := GetReadinessState(node)
	if !ready {
		return false
	}
	if node.Spec.Unschedulable {
		return false
	}
	return true
}
func GetReadinessState(node *apiv1.Node) (isNodeReady bool, lastTransitionTime time.Time, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	canNodeBeReady, readyFound := true, false
	lastTransitionTime = time.Time{}
	for _, cond := range node.Status.Conditions {
		switch cond.Type {
		case apiv1.NodeReady:
			readyFound = true
			if cond.Status == apiv1.ConditionFalse || cond.Status == apiv1.ConditionUnknown {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		case apiv1.NodeOutOfDisk:
			if cond.Status == apiv1.ConditionTrue {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		case apiv1.NodeNetworkUnavailable:
			if cond.Status == apiv1.ConditionTrue {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		}
	}
	if !readyFound {
		return false, time.Time{}, fmt.Errorf("readiness information not found")
	}
	return canNodeBeReady, lastTransitionTime, nil
}
