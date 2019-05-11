package autoscaling

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"
	ginkgo "github.com/onsi/ginkgo"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

const (
	dynamicConsumptionTimeInSeconds	= 30
	staticConsumptionTimeInSeconds	= 3600
	dynamicRequestSizeInMillicores	= 20
	dynamicRequestSizeInMegabytes	= 100
	dynamicRequestSizeCustomMetric	= 10
	port							= 80
	targetPort						= 8080
	timeoutRC						= 120 * time.Second
	startServiceTimeout				= time.Minute
	startServiceInterval			= 5 * time.Second
	rcIsNil							= "ERROR: replicationController = nil"
	deploymentIsNil					= "ERROR: deployment = nil"
	rsIsNil							= "ERROR: replicaset = nil"
	invalidKind						= "ERROR: invalid workload kind for resource consumer"
	customMetricName				= "QPS"
	serviceInitializationTimeout	= 2 * time.Minute
	serviceInitializationInterval	= 15 * time.Second
)

var (
	resourceConsumerImage			= imageutils.GetE2EImage(imageutils.ResourceConsumer)
	resourceConsumerControllerImage	= imageutils.GetE2EImage(imageutils.ResourceController)
)
var (
	KindRC			= schema.GroupVersionKind{Version: "v1", Kind: "ReplicationController"}
	KindDeployment	= schema.GroupVersionKind{Group: "apps", Version: "v1beta2", Kind: "Deployment"}
	KindReplicaSet	= schema.GroupVersionKind{Group: "apps", Version: "v1beta2", Kind: "ReplicaSet"}
	subresource		= "scale"
)

type ResourceConsumer struct {
	name						string
	controllerName				string
	kind						schema.GroupVersionKind
	nsName						string
	clientSet					clientset.Interface
	internalClientset			*internalclientset.Clientset
	cpu							chan int
	mem							chan int
	customMetric				chan int
	stopCPU						chan int
	stopMem						chan int
	stopCustomMetric			chan int
	stopWaitGroup				sync.WaitGroup
	consumptionTimeInSeconds	int
	sleepTime					time.Duration
	requestSizeInMillicores		int
	requestSizeInMegabytes		int
	requestSizeCustomMetric		int
}

func GetResourceConsumerImage() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return resourceConsumerImage
}
func NewDynamicResourceConsumer(name, nsName string, kind schema.GroupVersionKind, replicas, initCPUTotal, initMemoryTotal, initCustomMetric int, cpuRequest, memRequest resource.Quantity, clientset clientset.Interface, internalClientset *internalclientset.Clientset) *ResourceConsumer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newResourceConsumer(name, nsName, kind, replicas, initCPUTotal, initMemoryTotal, initCustomMetric, dynamicConsumptionTimeInSeconds, dynamicRequestSizeInMillicores, dynamicRequestSizeInMegabytes, dynamicRequestSizeCustomMetric, cpuRequest, memRequest, clientset, internalClientset)
}
func NewStaticResourceConsumer(name, nsName string, replicas, initCPUTotal, initMemoryTotal, initCustomMetric int, cpuRequest, memRequest resource.Quantity, clientset clientset.Interface, internalClientset *internalclientset.Clientset) *ResourceConsumer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newResourceConsumer(name, nsName, KindRC, replicas, initCPUTotal, initMemoryTotal, initCustomMetric, staticConsumptionTimeInSeconds, initCPUTotal/replicas, initMemoryTotal/replicas, initCustomMetric/replicas, cpuRequest, memRequest, clientset, internalClientset)
}
func newResourceConsumer(name, nsName string, kind schema.GroupVersionKind, replicas, initCPUTotal, initMemoryTotal, initCustomMetric, consumptionTimeInSeconds, requestSizeInMillicores, requestSizeInMegabytes int, requestSizeCustomMetric int, cpuRequest, memRequest resource.Quantity, clientset clientset.Interface, internalClientset *internalclientset.Clientset) *ResourceConsumer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	runServiceAndWorkloadForResourceConsumer(clientset, internalClientset, nsName, name, kind, replicas, cpuRequest, memRequest)
	rc := &ResourceConsumer{name: name, controllerName: name + "-ctrl", kind: kind, nsName: nsName, clientSet: clientset, internalClientset: internalClientset, cpu: make(chan int), mem: make(chan int), customMetric: make(chan int), stopCPU: make(chan int), stopMem: make(chan int), stopCustomMetric: make(chan int), consumptionTimeInSeconds: consumptionTimeInSeconds, sleepTime: time.Duration(consumptionTimeInSeconds) * time.Second, requestSizeInMillicores: requestSizeInMillicores, requestSizeInMegabytes: requestSizeInMegabytes, requestSizeCustomMetric: requestSizeCustomMetric}
	go rc.makeConsumeCPURequests()
	rc.ConsumeCPU(initCPUTotal)
	go rc.makeConsumeMemRequests()
	rc.ConsumeMem(initMemoryTotal)
	go rc.makeConsumeCustomMetric()
	rc.ConsumeCustomMetric(initCustomMetric)
	return rc
}
func (rc *ResourceConsumer) ConsumeCPU(millicores int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	framework.Logf("RC %s: consume %v millicores in total", rc.name, millicores)
	rc.cpu <- millicores
}
func (rc *ResourceConsumer) ConsumeMem(megabytes int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	framework.Logf("RC %s: consume %v MB in total", rc.name, megabytes)
	rc.mem <- megabytes
}
func (rc *ResourceConsumer) ConsumeCustomMetric(amount int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	framework.Logf("RC %s: consume custom metric %v in total", rc.name, amount)
	rc.customMetric <- amount
}
func (rc *ResourceConsumer) makeConsumeCPURequests() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer ginkgo.GinkgoRecover()
	rc.stopWaitGroup.Add(1)
	defer rc.stopWaitGroup.Done()
	sleepTime := time.Duration(0)
	millicores := 0
	for {
		select {
		case millicores = <-rc.cpu:
			framework.Logf("RC %s: setting consumption to %v millicores in total", rc.name, millicores)
		case <-time.After(sleepTime):
			framework.Logf("RC %s: sending request to consume %d millicores", rc.name, millicores)
			rc.sendConsumeCPURequest(millicores)
			sleepTime = rc.sleepTime
		case <-rc.stopCPU:
			framework.Logf("RC %s: stopping CPU consumer", rc.name)
			return
		}
	}
}
func (rc *ResourceConsumer) makeConsumeMemRequests() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer ginkgo.GinkgoRecover()
	rc.stopWaitGroup.Add(1)
	defer rc.stopWaitGroup.Done()
	sleepTime := time.Duration(0)
	megabytes := 0
	for {
		select {
		case megabytes = <-rc.mem:
			framework.Logf("RC %s: setting consumption to %v MB in total", rc.name, megabytes)
		case <-time.After(sleepTime):
			framework.Logf("RC %s: sending request to consume %d MB", rc.name, megabytes)
			rc.sendConsumeMemRequest(megabytes)
			sleepTime = rc.sleepTime
		case <-rc.stopMem:
			framework.Logf("RC %s: stopping mem consumer", rc.name)
			return
		}
	}
}
func (rc *ResourceConsumer) makeConsumeCustomMetric() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer ginkgo.GinkgoRecover()
	rc.stopWaitGroup.Add(1)
	defer rc.stopWaitGroup.Done()
	sleepTime := time.Duration(0)
	delta := 0
	for {
		select {
		case delta := <-rc.customMetric:
			framework.Logf("RC %s: setting bump of metric %s to %d in total", rc.name, customMetricName, delta)
		case <-time.After(sleepTime):
			framework.Logf("RC %s: sending request to consume %d of custom metric %s", rc.name, delta, customMetricName)
			rc.sendConsumeCustomMetric(delta)
			sleepTime = rc.sleepTime
		case <-rc.stopCustomMetric:
			framework.Logf("RC %s: stopping metric consumer", rc.name)
			return
		}
	}
}
func (rc *ResourceConsumer) sendConsumeCPURequest(millicores int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := context.WithTimeout(context.Background(), framework.SingleCallTimeout)
	defer cancel()
	err := wait.PollImmediate(serviceInitializationInterval, serviceInitializationTimeout, func() (bool, error) {
		proxyRequest, err := framework.GetServicesProxyRequest(rc.clientSet, rc.clientSet.CoreV1().RESTClient().Post())
		framework.ExpectNoError(err)
		req := proxyRequest.Namespace(rc.nsName).Context(ctx).Name(rc.controllerName).Suffix("ConsumeCPU").Param("millicores", strconv.Itoa(millicores)).Param("durationSec", strconv.Itoa(rc.consumptionTimeInSeconds)).Param("requestSizeMillicores", strconv.Itoa(rc.requestSizeInMillicores))
		framework.Logf("ConsumeCPU URL: %v", *req.URL())
		_, err = req.DoRaw()
		if err != nil {
			framework.Logf("ConsumeCPU failure: %v", err)
			return false, nil
		}
		return true, nil
	})
	framework.ExpectNoError(err)
}
func (rc *ResourceConsumer) sendConsumeMemRequest(megabytes int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := context.WithTimeout(context.Background(), framework.SingleCallTimeout)
	defer cancel()
	err := wait.PollImmediate(serviceInitializationInterval, serviceInitializationTimeout, func() (bool, error) {
		proxyRequest, err := framework.GetServicesProxyRequest(rc.clientSet, rc.clientSet.CoreV1().RESTClient().Post())
		framework.ExpectNoError(err)
		req := proxyRequest.Namespace(rc.nsName).Context(ctx).Name(rc.controllerName).Suffix("ConsumeMem").Param("megabytes", strconv.Itoa(megabytes)).Param("durationSec", strconv.Itoa(rc.consumptionTimeInSeconds)).Param("requestSizeMegabytes", strconv.Itoa(rc.requestSizeInMegabytes))
		framework.Logf("ConsumeMem URL: %v", *req.URL())
		_, err = req.DoRaw()
		if err != nil {
			framework.Logf("ConsumeMem failure: %v", err)
			return false, nil
		}
		return true, nil
	})
	framework.ExpectNoError(err)
}
func (rc *ResourceConsumer) sendConsumeCustomMetric(delta int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := context.WithTimeout(context.Background(), framework.SingleCallTimeout)
	defer cancel()
	err := wait.PollImmediate(serviceInitializationInterval, serviceInitializationTimeout, func() (bool, error) {
		proxyRequest, err := framework.GetServicesProxyRequest(rc.clientSet, rc.clientSet.CoreV1().RESTClient().Post())
		framework.ExpectNoError(err)
		req := proxyRequest.Namespace(rc.nsName).Context(ctx).Name(rc.controllerName).Suffix("BumpMetric").Param("metric", customMetricName).Param("delta", strconv.Itoa(delta)).Param("durationSec", strconv.Itoa(rc.consumptionTimeInSeconds)).Param("requestSizeMetrics", strconv.Itoa(rc.requestSizeCustomMetric))
		framework.Logf("ConsumeCustomMetric URL: %v", *req.URL())
		_, err = req.DoRaw()
		if err != nil {
			framework.Logf("ConsumeCustomMetric failure: %v", err)
			return false, nil
		}
		return true, nil
	})
	framework.ExpectNoError(err)
}
func (rc *ResourceConsumer) GetReplicas() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch rc.kind {
	case KindRC:
		replicationController, err := rc.clientSet.CoreV1().ReplicationControllers(rc.nsName).Get(rc.name, metav1.GetOptions{})
		framework.ExpectNoError(err)
		if replicationController == nil {
			framework.Failf(rcIsNil)
		}
		return int(replicationController.Status.ReadyReplicas)
	case KindDeployment:
		deployment, err := rc.clientSet.ExtensionsV1beta1().Deployments(rc.nsName).Get(rc.name, metav1.GetOptions{})
		framework.ExpectNoError(err)
		if deployment == nil {
			framework.Failf(deploymentIsNil)
		}
		return int(deployment.Status.ReadyReplicas)
	case KindReplicaSet:
		rs, err := rc.clientSet.ExtensionsV1beta1().ReplicaSets(rc.nsName).Get(rc.name, metav1.GetOptions{})
		framework.ExpectNoError(err)
		if rs == nil {
			framework.Failf(rsIsNil)
		}
		return int(rs.Status.ReadyReplicas)
	default:
		framework.Failf(invalidKind)
	}
	return 0
}
func (rc *ResourceConsumer) WaitForReplicas(desiredReplicas int, duration time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interval := 20 * time.Second
	err := wait.PollImmediate(interval, duration, func() (bool, error) {
		replicas := rc.GetReplicas()
		framework.Logf("waiting for %d replicas (current: %d)", desiredReplicas, replicas)
		return replicas == desiredReplicas, nil
	})
	framework.ExpectNoErrorWithOffset(1, err, "timeout waiting %v for %d replicas", duration, desiredReplicas)
}
func (rc *ResourceConsumer) EnsureDesiredReplicas(desiredReplicas int, duration time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interval := 10 * time.Second
	err := wait.PollImmediate(interval, duration, func() (bool, error) {
		replicas := rc.GetReplicas()
		framework.Logf("expecting there to be %d replicas (are: %d)", desiredReplicas, replicas)
		if replicas != desiredReplicas {
			return false, fmt.Errorf("number of replicas changed unexpectedly")
		}
		return false, nil
	})
	if err == wait.ErrWaitTimeout {
		framework.Logf("Number of replicas was stable over %v", duration)
		return
	}
	framework.ExpectNoErrorWithOffset(1, err)
}
func (rc *ResourceConsumer) Pause() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ginkgo.By(fmt.Sprintf("HPA pausing RC %s", rc.name))
	rc.stopCPU <- 0
	rc.stopMem <- 0
	rc.stopCustomMetric <- 0
	rc.stopWaitGroup.Wait()
}
func (rc *ResourceConsumer) Resume() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ginkgo.By(fmt.Sprintf("HPA resuming RC %s", rc.name))
	go rc.makeConsumeCPURequests()
	go rc.makeConsumeMemRequests()
	go rc.makeConsumeCustomMetric()
}
func (rc *ResourceConsumer) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ginkgo.By(fmt.Sprintf("Removing consuming RC %s", rc.name))
	close(rc.stopCPU)
	close(rc.stopMem)
	close(rc.stopCustomMetric)
	rc.stopWaitGroup.Wait()
	time.Sleep(10 * time.Second)
	kind := rc.kind.GroupKind()
	framework.ExpectNoError(framework.DeleteResourceAndWaitForGC(rc.clientSet, kind, rc.nsName, rc.name))
	framework.ExpectNoError(rc.clientSet.CoreV1().Services(rc.nsName).Delete(rc.name, nil))
	framework.ExpectNoError(framework.DeleteResourceAndWaitForGC(rc.clientSet, api.Kind("ReplicationController"), rc.nsName, rc.controllerName))
	framework.ExpectNoError(rc.clientSet.CoreV1().Services(rc.nsName).Delete(rc.controllerName, nil))
}
func runServiceAndWorkloadForResourceConsumer(c clientset.Interface, internalClient internalclientset.Interface, ns, name string, kind schema.GroupVersionKind, replicas int, cpuLimit, memLimit resource.Quantity) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ginkgo.By(fmt.Sprintf("Running consuming RC %s via %s with %v replicas", name, kind, replicas))
	_, err := c.CoreV1().Services(ns).Create(&v1.Service{ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Port: port, TargetPort: intstr.FromInt(targetPort)}}, Selector: map[string]string{"name": name}}})
	framework.ExpectNoError(err)
	rcConfig := testutils.RCConfig{Client: c, InternalClient: internalClient, Image: resourceConsumerImage, Name: name, Namespace: ns, Timeout: timeoutRC, Replicas: replicas, CpuRequest: cpuLimit.MilliValue(), MemRequest: memLimit.Value()}
	switch kind {
	case KindRC:
		framework.ExpectNoError(framework.RunRC(rcConfig))
		break
	case KindDeployment:
		dpConfig := testutils.DeploymentConfig{RCConfig: rcConfig}
		framework.ExpectNoError(framework.RunDeployment(dpConfig))
		break
	case KindReplicaSet:
		rsConfig := testutils.ReplicaSetConfig{RCConfig: rcConfig}
		ginkgo.By(fmt.Sprintf("creating replicaset %s in namespace %s", rsConfig.Name, rsConfig.Namespace))
		framework.ExpectNoError(framework.RunReplicaSet(rsConfig))
		break
	default:
		framework.Failf(invalidKind)
	}
	ginkgo.By(fmt.Sprintf("Running controller"))
	controllerName := name + "-ctrl"
	_, err = c.CoreV1().Services(ns).Create(&v1.Service{ObjectMeta: metav1.ObjectMeta{Name: controllerName}, Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Port: port, TargetPort: intstr.FromInt(targetPort)}}, Selector: map[string]string{"name": controllerName}}})
	framework.ExpectNoError(err)
	dnsClusterFirst := v1.DNSClusterFirst
	controllerRcConfig := testutils.RCConfig{Client: c, Image: resourceConsumerControllerImage, Name: controllerName, Namespace: ns, Timeout: timeoutRC, Replicas: 1, Command: []string{"/controller", "--consumer-service-name=" + name, "--consumer-service-namespace=" + ns, "--consumer-port=80"}, DNSPolicy: &dnsClusterFirst}
	framework.ExpectNoError(framework.RunRC(controllerRcConfig))
	framework.ExpectNoError(framework.WaitForServiceEndpointsNum(c, ns, controllerName, 1, startServiceInterval, startServiceTimeout))
}
func CreateCPUHorizontalPodAutoscaler(rc *ResourceConsumer, cpu, minReplicas, maxRepl int32) *autoscalingv1.HorizontalPodAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hpa := &autoscalingv1.HorizontalPodAutoscaler{ObjectMeta: metav1.ObjectMeta{Name: rc.name, Namespace: rc.nsName}, Spec: autoscalingv1.HorizontalPodAutoscalerSpec{ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{APIVersion: rc.kind.GroupVersion().String(), Kind: rc.kind.Kind, Name: rc.name}, MinReplicas: &minReplicas, MaxReplicas: maxRepl, TargetCPUUtilizationPercentage: &cpu}}
	hpa, errHPA := rc.clientSet.AutoscalingV1().HorizontalPodAutoscalers(rc.nsName).Create(hpa)
	framework.ExpectNoError(errHPA)
	return hpa
}
func DeleteHorizontalPodAutoscaler(rc *ResourceConsumer, autoscalerName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rc.clientSet.AutoscalingV1().HorizontalPodAutoscalers(rc.nsName).Delete(autoscalerName, nil)
}
