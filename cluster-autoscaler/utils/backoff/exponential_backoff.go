package backoff

import (
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type exponentialBackoff struct {
	maxBackoffDuration	time.Duration
	initialBackoffDuration	time.Duration
	backoffResetTimeout	time.Duration
	backoffInfo		map[string]exponentialBackoffInfo
	nodeGroupKey		func(nodeGroup cloudprovider.NodeGroup) string
}
type exponentialBackoffInfo struct {
	duration		time.Duration
	backoffUntil		time.Time
	lastFailedExecution	time.Time
}

func NewExponentialBackoff(initialBackoffDuration time.Duration, maxBackoffDuration time.Duration, backoffResetTimeout time.Duration, nodeGroupKey func(nodeGroup cloudprovider.NodeGroup) string) Backoff {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &exponentialBackoff{maxBackoffDuration: maxBackoffDuration, initialBackoffDuration: initialBackoffDuration, backoffResetTimeout: backoffResetTimeout, backoffInfo: make(map[string]exponentialBackoffInfo), nodeGroupKey: nodeGroupKey}
}
func NewIdBasedExponentialBackoff(initialBackoffDuration time.Duration, maxBackoffDuration time.Duration, backoffResetTimeout time.Duration) Backoff {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewExponentialBackoff(initialBackoffDuration, maxBackoffDuration, backoffResetTimeout, func(nodeGroup cloudprovider.NodeGroup) string {
		return nodeGroup.Id()
	})
}
func (b *exponentialBackoff) Backoff(nodeGroup cloudprovider.NodeGroup, currentTime time.Time) time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	duration := b.initialBackoffDuration
	key := b.nodeGroupKey(nodeGroup)
	if backoffInfo, found := b.backoffInfo[key]; found {
		if backoffInfo.backoffUntil.Before(currentTime) {
			duration = 2 * backoffInfo.duration
			if duration > b.maxBackoffDuration {
				duration = b.maxBackoffDuration
			}
		}
	}
	backoffUntil := currentTime.Add(duration)
	b.backoffInfo[key] = exponentialBackoffInfo{duration: duration, backoffUntil: backoffUntil, lastFailedExecution: currentTime}
	return backoffUntil
}
func (b *exponentialBackoff) IsBackedOff(nodeGroup cloudprovider.NodeGroup, currentTime time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoffInfo, found := b.backoffInfo[b.nodeGroupKey(nodeGroup)]
	return found && backoffInfo.backoffUntil.After(currentTime)
}
func (b *exponentialBackoff) RemoveBackoff(nodeGroup cloudprovider.NodeGroup) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	delete(b.backoffInfo, b.nodeGroupKey(nodeGroup))
}
func (b *exponentialBackoff) RemoveStaleBackoffData(currentTime time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for key, backoffInfo := range b.backoffInfo {
		if backoffInfo.lastFailedExecution.Add(b.backoffResetTimeout).Before(currentTime) {
			delete(b.backoffInfo, key)
		}
	}
}
