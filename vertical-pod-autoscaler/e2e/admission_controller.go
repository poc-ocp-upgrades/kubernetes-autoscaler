package autoscaling

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/kubernetes/test/e2e/framework"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = AdmissionControllerE2eDescribe("Admission-controller", func() {
	f := framework.NewDefaultFramework("vertical-pod-autoscaling")
	ginkgo.It("starts pods with new recommended request", func() {
		d := NewHamsterDeploymentWithResources(f, ParseQuantityOrDie("100m"), ParseQuantityOrDie("100Mi"))
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		vpaCRD.Status.Recommendation = &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: "hamster", Target: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("250m"), apiv1.ResourceMemory: ParseQuantityOrDie("200Mi")}}}}
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]).To(gomega.Equal(ParseQuantityOrDie("250m")))
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceMemory]).To(gomega.Equal(ParseQuantityOrDie("200Mi")))
		}
	})
	ginkgo.It("caps request to limit set by the user", func() {
		d := NewHamsterDeploymentWithResources(f, ParseQuantityOrDie("100m"), ParseQuantityOrDie("100Mi"))
		d.Spec.Template.Spec.Containers[0].Resources.Limits = apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("222m"), apiv1.ResourceMemory: ParseQuantityOrDie("123Mi")}
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		vpaCRD.Status.Recommendation = &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: "hamster", Target: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("250m"), apiv1.ResourceMemory: ParseQuantityOrDie("200Mi")}}}}
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]).To(gomega.Equal(ParseQuantityOrDie("222m")))
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceMemory]).To(gomega.Equal(ParseQuantityOrDie("123Mi")))
		}
	})
	ginkgo.It("caps request to max set in VPA", func() {
		d := NewHamsterDeploymentWithResources(f, ParseQuantityOrDie("100m"), ParseQuantityOrDie("100Mi"))
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		vpaCRD.Status.Recommendation = &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: "hamster", Target: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("250m"), apiv1.ResourceMemory: ParseQuantityOrDie("200Mi")}}}}
		vpaCRD.Spec.ResourcePolicy = &vpa_types.PodResourcePolicy{ContainerPolicies: []vpa_types.ContainerResourcePolicy{{ContainerName: "hamster", MaxAllowed: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("233m"), apiv1.ResourceMemory: ParseQuantityOrDie("150Mi")}}}}
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]).To(gomega.Equal(ParseQuantityOrDie("233m")))
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceMemory]).To(gomega.Equal(ParseQuantityOrDie("150Mi")))
		}
	})
	ginkgo.It("raises request to min set in VPA", func() {
		d := NewHamsterDeploymentWithResources(f, ParseQuantityOrDie("100m"), ParseQuantityOrDie("100Mi"))
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		vpaCRD.Status.Recommendation = &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: "hamster", Target: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("50m"), apiv1.ResourceMemory: ParseQuantityOrDie("60Mi")}}}}
		vpaCRD.Spec.ResourcePolicy = &vpa_types.PodResourcePolicy{ContainerPolicies: []vpa_types.ContainerResourcePolicy{{ContainerName: "hamster", MinAllowed: apiv1.ResourceList{apiv1.ResourceCPU: ParseQuantityOrDie("90m"), apiv1.ResourceMemory: ParseQuantityOrDie("80Mi")}}}}
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]).To(gomega.Equal(ParseQuantityOrDie("90m")))
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceMemory]).To(gomega.Equal(ParseQuantityOrDie("80Mi")))
		}
	})
	ginkgo.It("leaves users request when no recommendation", func() {
		d := NewHamsterDeploymentWithResources(f, ParseQuantityOrDie("100m"), ParseQuantityOrDie("100Mi"))
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]).To(gomega.Equal(ParseQuantityOrDie("100m")))
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceMemory]).To(gomega.Equal(ParseQuantityOrDie("100Mi")))
		}
	})
	ginkgo.It("passes empty request when no recommendation and no user-specified request", func() {
		d := NewHamsterDeployment(f)
		ginkgo.By("Setting up a VPA CRD")
		vpaCRD := NewVPA(f, "hamster-vpa", &metav1.LabelSelector{MatchLabels: d.Spec.Template.Labels})
		InstallVPA(f, vpaCRD)
		ginkgo.By("Setting up a hamster deployment")
		podList := startDeploymentPods(f, d)
		for _, pod := range podList.Items {
			gomega.Expect(pod.Spec.Containers[0].Resources.Requests).To(gomega.BeEmpty())
		}
	})
})

func startDeploymentPods(f *framework.Framework, deployment *appsv1.Deployment) *apiv1.PodList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, ns := f.ClientSet, f.Namespace.Name
	deployment, err := c.AppsV1().Deployments(ns).Create(deployment)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = framework.WaitForDeploymentComplete(c, deployment)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	podList, err := framework.GetPodsForDeployment(c, deployment)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return podList
}
