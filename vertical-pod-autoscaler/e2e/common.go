package autoscaling

import (
	"fmt"
	"time"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/kubernetes/test/e2e/framework"
)

const (
	recommenderComponent			= "recommender"
	updateComponent					= "updater"
	admissionControllerComponent	= "admission-controller"
	fullVpaSuite					= "full-vpa"
	actuationSuite					= "actuation"
	pollInterval					= 10 * time.Second
	pollTimeout						= 15 * time.Minute
	VpaEvictionTimeout				= 3 * time.Minute
	defaultHamsterReplicas			= int32(3)
)

var hamsterLabels = map[string]string{"app": "hamster"}

func SIGDescribe(text string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ginkgo.Describe(fmt.Sprintf("[sig-autoscaling] %v", text), body)
}
func E2eDescribe(scenario, name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SIGDescribe(fmt.Sprintf("[VPA] [%s] %s", scenario, name), body)
}
func RecommenderE2eDescribe(name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return E2eDescribe(recommenderComponent, name, body)
}
func UpdaterE2eDescribe(name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return E2eDescribe(updateComponent, name, body)
}
func AdmissionControllerE2eDescribe(name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return E2eDescribe(admissionControllerComponent, name, body)
}
func FullVpaE2eDescribe(name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return E2eDescribe(fullVpaSuite, name, body)
}
func ActuationSuiteE2eDescribe(name string, body func()) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return E2eDescribe(actuationSuite, name, body)
}
func SetupHamsterDeployment(f *framework.Framework, cpu, memory string, replicas int32) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuQuantity := ParseQuantityOrDie(cpu)
	memoryQuantity := ParseQuantityOrDie(memory)
	d := NewHamsterDeploymentWithResources(f, cpuQuantity, memoryQuantity)
	d.Spec.Replicas = &replicas
	d, err := f.ClientSet.AppsV1().Deployments(f.Namespace.Name).Create(d)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = framework.WaitForDeploymentComplete(f.ClientSet, d)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return d
}
func NewHamsterDeployment(f *framework.Framework) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := framework.NewDeployment("hamster-deployment", defaultHamsterReplicas, hamsterLabels, "hamster", "k8s.gcr.io/ubuntu-slim:0.1", appsv1.RollingUpdateDeploymentStrategyType)
	d.ObjectMeta.Namespace = f.Namespace.Name
	d.Spec.Template.Spec.Containers[0].Command = []string{"/bin/sh"}
	d.Spec.Template.Spec.Containers[0].Args = []string{"-c", "/usr/bin/yes >/dev/null"}
	return d
}
func NewHamsterDeploymentWithResources(f *framework.Framework, cpuQuantity, memoryQuantity resource.Quantity) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := NewHamsterDeployment(f)
	d.Spec.Template.Spec.Containers[0].Resources.Requests = apiv1.ResourceList{apiv1.ResourceCPU: cpuQuantity, apiv1.ResourceMemory: memoryQuantity}
	return d
}
func GetHamsterPods(f *framework.Framework) (*apiv1.PodList, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	label := labels.SelectorFromSet(labels.Set(hamsterLabels))
	selector := fields.ParseSelectorOrDie("status.phase!=" + string(apiv1.PodSucceeded) + ",status.phase!=" + string(apiv1.PodFailed))
	options := metav1.ListOptions{LabelSelector: label.String(), FieldSelector: selector.String()}
	return f.ClientSet.CoreV1().Pods(f.Namespace.Name).List(options)
}
func SetupVPA(f *framework.Framework, cpu string, mode vpa_types.UpdateMode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: hamsterLabels})
	vpaCRD.Spec.UpdatePolicy.UpdateMode = &mode
	cpuQuantity := ParseQuantityOrDie(cpu)
	resourceList := apiv1.ResourceList{apiv1.ResourceCPU: cpuQuantity}
	vpaCRD.Status.Recommendation = &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: "hamster", Target: resourceList, LowerBound: resourceList, UpperBound: resourceList}}}
	InstallVPA(f, vpaCRD)
}
func NewVPA(f *framework.Framework, name string, selector *metav1.LabelSelector) *vpa_types.VerticalPodAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	updateMode := vpa_types.UpdateModeAuto
	vpa := vpa_types.VerticalPodAutoscaler{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: f.Namespace.Name}, Spec: vpa_types.VerticalPodAutoscalerSpec{Selector: selector, UpdatePolicy: &vpa_types.PodUpdatePolicy{UpdateMode: &updateMode}, ResourcePolicy: &vpa_types.PodResourcePolicy{ContainerPolicies: []vpa_types.ContainerResourcePolicy{}}}}
	return &vpa
}
func InstallVPA(f *framework.Framework, vpa *vpa_types.VerticalPodAutoscaler) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ns := f.Namespace.Name
	config, err := framework.LoadConfig()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	vpaClientSet := vpa_clientset.NewForConfigOrDie(config)
	vpaClient := vpaClientSet.AutoscalingV1beta1()
	_, err = vpaClient.VerticalPodAutoscalers(ns).Create(vpa)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
func ParseQuantityOrDie(text string) resource.Quantity {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quantity, err := resource.ParseQuantity(text)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return quantity
}

type PodSet map[string]types.UID

func MakePodSet(pods *apiv1.PodList) PodSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(PodSet)
	if pods == nil {
		return result
	}
	for _, p := range pods.Items {
		result[p.Name] = p.UID
	}
	return result
}
func WaitForPodsRestarted(f *framework.Framework, podList *apiv1.PodList) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	initialPodSet := MakePodSet(podList)
	err := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
		currentPodList, err := GetHamsterPods(f)
		if err != nil {
			return false, err
		}
		currentPodSet := MakePodSet(currentPodList)
		return WerePodsSuccessfullyRestarted(currentPodSet, initialPodSet), nil
	})
	if err != nil {
		return fmt.Errorf("Waiting for set of pods changed: %v", err)
	}
	return nil
}
func WaitForPodsEvicted(f *framework.Framework, podList *apiv1.PodList) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	initialPodSet := MakePodSet(podList)
	err := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
		currentPodList, err := GetHamsterPods(f)
		if err != nil {
			return false, err
		}
		currentPodSet := MakePodSet(currentPodList)
		return GetEvictedPodsCount(currentPodSet, initialPodSet) > 0, nil
	})
	if err != nil {
		return fmt.Errorf("Waiting for set of pods changed: %v", err)
	}
	return nil
}
func WerePodsSuccessfullyRestarted(currentPodSet PodSet, initialPodSet PodSet) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(currentPodSet) < len(initialPodSet) {
		framework.Logf("Restart in progress")
		return false
	}
	evictedCount := GetEvictedPodsCount(currentPodSet, initialPodSet)
	framework.Logf("%v of initial pods were already evicted", evictedCount)
	return evictedCount > 0
}
func GetEvictedPodsCount(currentPodSet PodSet, initialPodSet PodSet) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	diffs := 0
	for name, initialUID := range initialPodSet {
		currentUID, inCurrent := currentPodSet[name]
		if !inCurrent {
			diffs += 1
		} else if initialUID != currentUID {
			diffs += 1
		}
	}
	return diffs
}
func CheckNoPodsEvicted(f *framework.Framework, initialPodSet PodSet) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	time.Sleep(VpaEvictionTimeout)
	currentPodList, err := GetHamsterPods(f)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	restarted := GetEvictedPodsCount(MakePodSet(currentPodList), initialPodSet)
	gomega.Expect(restarted).To(gomega.Equal(0))
}
