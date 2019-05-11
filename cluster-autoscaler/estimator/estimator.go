package estimator

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/klog"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
	BasicEstimatorName		= "basic"
	BinpackingEstimatorName	= "binpacking"
)

func deprecated(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (DEPRECATED)", name)
}

var AvailableEstimators = []string{BinpackingEstimatorName, deprecated(BasicEstimatorName)}

type Estimator interface {
	Estimate([]*apiv1.Pod, *schedulercache.NodeInfo, []*schedulercache.NodeInfo) int
}
type EstimatorBuilder func(*simulator.PredicateChecker) Estimator

func NewEstimatorBuilder(name string) (EstimatorBuilder, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch name {
	case BinpackingEstimatorName:
		return func(predicateChecker *simulator.PredicateChecker) Estimator {
			return NewBinpackingNodeEstimator(predicateChecker)
		}, nil
	case BasicEstimatorName:
		klog.Warning(basicEstimatorDeprecationMessage)
		return func(_ *simulator.PredicateChecker) Estimator {
			return NewBasicNodeEstimator()
		}, nil
	}
	return nil, fmt.Errorf("Unknown estimator: %s", name)
}
