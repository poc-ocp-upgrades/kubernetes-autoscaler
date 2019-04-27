package autoscaling

import (
	"testing"
	"k8s.io/kubernetes/test/e2e/framework"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	framework.ViperizeFlags()
}
func TestVpaE2E(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	RunE2ETests(t)
}
