package glogx

import (
	"k8s.io/klog"
)

type quota struct {
	limit	int
	left	int
}

func NewLoggingQuota(n int) *quota {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &quota{n, n}
}
func (q *quota) Left() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return q.left
}
func (q *quota) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	q.left = q.limit
}
func UpTo(quota *quota) klog.Verbose {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	quota.left--
	return quota.left >= 0
}
func Over(quota *quota) klog.Verbose {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return quota.left < 0
}
func V(n klog.Level) Verbose {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Verbose(klog.V(n))
}

type Verbose klog.Verbose

func (v Verbose) UpTo(quota *quota) klog.Verbose {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if v {
		return UpTo(quota)
	}
	return klog.Verbose(false)
}
func (v Verbose) Over(quota *quota) klog.Verbose {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if v {
		return Over(quota)
	}
	return klog.Verbose(false)
}
