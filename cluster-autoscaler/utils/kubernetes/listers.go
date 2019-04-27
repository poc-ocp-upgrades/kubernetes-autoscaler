package kubernetes

import (
	"time"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1 "k8s.io/api/extensions/v1beta1"
	policyv1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	client "k8s.io/client-go/kubernetes"
	v1lister "k8s.io/client-go/listers/core/v1"
	v1extensionslister "k8s.io/client-go/listers/extensions/v1beta1"
	v1policylister "k8s.io/client-go/listers/policy/v1beta1"
	"k8s.io/client-go/tools/cache"
	podv1 "k8s.io/kubernetes/pkg/api/v1/pod"
)

type ListerRegistry interface {
	AllNodeLister() NodeLister
	ReadyNodeLister() NodeLister
	ScheduledPodLister() PodLister
	UnschedulablePodLister() PodLister
	PodDisruptionBudgetLister() PodDisruptionBudgetLister
	DaemonSetLister() DaemonSetLister
}
type listerRegistryImpl struct {
	allNodeLister			NodeLister
	readyNodeLister			NodeLister
	scheduledPodLister		PodLister
	unschedulablePodLister		PodLister
	podDisruptionBudgetLister	PodDisruptionBudgetLister
	daemonSetLister			DaemonSetLister
}

func NewListerRegistry(allNode NodeLister, readyNode NodeLister, scheduledPod PodLister, unschedulablePod PodLister, podDisruptionBudgetLister PodDisruptionBudgetLister, daemonSetLister DaemonSetLister) ListerRegistry {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return listerRegistryImpl{allNodeLister: allNode, readyNodeLister: readyNode, scheduledPodLister: scheduledPod, unschedulablePodLister: unschedulablePod, podDisruptionBudgetLister: podDisruptionBudgetLister, daemonSetLister: daemonSetLister}
}
func NewListerRegistryWithDefaultListers(kubeClient client.Interface, stopChannel <-chan struct{}) ListerRegistry {
	_logClusterCodePath()
	defer _logClusterCodePath()
	unschedulablePodLister := NewUnschedulablePodLister(kubeClient, stopChannel)
	scheduledPodLister := NewScheduledPodLister(kubeClient, stopChannel)
	readyNodeLister := NewReadyNodeLister(kubeClient, stopChannel)
	allNodeLister := NewAllNodeLister(kubeClient, stopChannel)
	podDisruptionBudgetLister := NewPodDisruptionBudgetLister(kubeClient, stopChannel)
	daemonSetLister := NewDaemonSetLister(kubeClient, stopChannel)
	return NewListerRegistry(allNodeLister, readyNodeLister, scheduledPodLister, unschedulablePodLister, podDisruptionBudgetLister, daemonSetLister)
}
func (r listerRegistryImpl) AllNodeLister() NodeLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.allNodeLister
}
func (r listerRegistryImpl) ReadyNodeLister() NodeLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.readyNodeLister
}
func (r listerRegistryImpl) ScheduledPodLister() PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.scheduledPodLister
}
func (r listerRegistryImpl) UnschedulablePodLister() PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.unschedulablePodLister
}
func (r listerRegistryImpl) PodDisruptionBudgetLister() PodDisruptionBudgetLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.podDisruptionBudgetLister
}
func (r listerRegistryImpl) DaemonSetLister() DaemonSetLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.daemonSetLister
}

type PodLister interface{ List() ([]*apiv1.Pod, error) }
type UnschedulablePodLister struct{ podLister v1lister.PodLister }

func (unschedulablePodLister *UnschedulablePodLister) List() ([]*apiv1.Pod, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var unschedulablePods []*apiv1.Pod
	allPods, err := unschedulablePodLister.podLister.List(labels.Everything())
	if err != nil {
		return unschedulablePods, err
	}
	for _, pod := range allPods {
		_, condition := podv1.GetPodCondition(&pod.Status, apiv1.PodScheduled)
		if condition != nil && condition.Status == apiv1.ConditionFalse && condition.Reason == apiv1.PodReasonUnschedulable {
			unschedulablePods = append(unschedulablePods, pod)
		}
	}
	return unschedulablePods, nil
}
func NewUnschedulablePodLister(kubeClient client.Interface, stopchannel <-chan struct{}) PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewUnschedulablePodInNamespaceLister(kubeClient, apiv1.NamespaceAll, stopchannel)
}
func NewUnschedulablePodInNamespaceLister(kubeClient client.Interface, namespace string, stopchannel <-chan struct{}) PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	selector := fields.ParseSelectorOrDie("spec.nodeName==" + "" + ",status.phase!=" + string(apiv1.PodSucceeded) + ",status.phase!=" + string(apiv1.PodFailed))
	podListWatch := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", namespace, selector)
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	podLister := v1lister.NewPodLister(store)
	podReflector := cache.NewReflector(podListWatch, &apiv1.Pod{}, store, time.Hour)
	go podReflector.Run(stopchannel)
	return &UnschedulablePodLister{podLister: podLister}
}

type ScheduledPodLister struct{ podLister v1lister.PodLister }

func (lister *ScheduledPodLister) List() ([]*apiv1.Pod, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return lister.podLister.List(labels.Everything())
}
func NewScheduledPodLister(kubeClient client.Interface, stopchannel <-chan struct{}) PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	selector := fields.ParseSelectorOrDie("spec.nodeName!=" + "" + ",status.phase!=" + string(apiv1.PodSucceeded) + ",status.phase!=" + string(apiv1.PodFailed))
	podListWatch := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", apiv1.NamespaceAll, selector)
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	podLister := v1lister.NewPodLister(store)
	podReflector := cache.NewReflector(podListWatch, &apiv1.Pod{}, store, time.Hour)
	go podReflector.Run(stopchannel)
	return &ScheduledPodLister{podLister: podLister}
}

type NodeLister interface{ List() ([]*apiv1.Node, error) }
type ReadyNodeLister struct{ nodeLister v1lister.NodeLister }

func (readyNodeLister *ReadyNodeLister) List() ([]*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodes, err := readyNodeLister.nodeLister.List(labels.Everything())
	if err != nil {
		return []*apiv1.Node{}, err
	}
	readyNodes := make([]*apiv1.Node, 0, len(nodes))
	for _, node := range nodes {
		if IsNodeReadyAndSchedulable(node) {
			readyNodes = append(readyNodes, node)
		}
	}
	return readyNodes, nil
}
func NewReadyNodeLister(kubeClient client.Interface, stopChannel <-chan struct{}) NodeLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listWatcher := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "nodes", apiv1.NamespaceAll, fields.Everything())
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	nodeLister := v1lister.NewNodeLister(store)
	reflector := cache.NewReflector(listWatcher, &apiv1.Node{}, store, time.Hour)
	go reflector.Run(stopChannel)
	return &ReadyNodeLister{nodeLister: nodeLister}
}

type AllNodeLister struct{ nodeLister v1lister.NodeLister }

func (allNodeLister *AllNodeLister) List() ([]*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodes, err := allNodeLister.nodeLister.List(labels.Everything())
	if err != nil {
		return []*apiv1.Node{}, err
	}
	allNodes := append(make([]*apiv1.Node, 0, len(nodes)), nodes...)
	return allNodes, nil
}
func NewAllNodeLister(kubeClient client.Interface, stopchannel <-chan struct{}) NodeLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listWatcher := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "nodes", apiv1.NamespaceAll, fields.Everything())
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	nodeLister := v1lister.NewNodeLister(store)
	reflector := cache.NewReflector(listWatcher, &apiv1.Node{}, store, time.Hour)
	go reflector.Run(stopchannel)
	return &AllNodeLister{nodeLister: nodeLister}
}

type PodDisruptionBudgetLister interface {
	List() ([]*policyv1.PodDisruptionBudget, error)
}
type PodDisruptionBudgetListerImpl struct {
	pdbLister v1policylister.PodDisruptionBudgetLister
}

func (lister *PodDisruptionBudgetListerImpl) List() ([]*policyv1.PodDisruptionBudget, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return lister.pdbLister.List(labels.Everything())
}
func NewPodDisruptionBudgetLister(kubeClient client.Interface, stopchannel <-chan struct{}) PodDisruptionBudgetLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listWatcher := cache.NewListWatchFromClient(kubeClient.Policy().RESTClient(), "poddisruptionbudgets", apiv1.NamespaceAll, fields.Everything())
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	pdbLister := v1policylister.NewPodDisruptionBudgetLister(store)
	reflector := cache.NewReflector(listWatcher, &policyv1.PodDisruptionBudget{}, store, time.Hour)
	go reflector.Run(stopchannel)
	return &PodDisruptionBudgetListerImpl{pdbLister: pdbLister}
}

type DaemonSetLister interface {
	List() ([]*extensionsv1.DaemonSet, error)
}
type DaemonSetListerImpl struct {
	daemonSetLister v1extensionslister.DaemonSetLister
}

func (lister *DaemonSetListerImpl) List() ([]*extensionsv1.DaemonSet, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return lister.daemonSetLister.List(labels.Everything())
}
func NewDaemonSetLister(kubeClient client.Interface, stopchannel <-chan struct{}) DaemonSetLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listWatcher := cache.NewListWatchFromClient(kubeClient.Extensions().RESTClient(), "daemonsets", apiv1.NamespaceAll, fields.Everything())
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	lister := v1extensionslister.NewDaemonSetLister(store)
	reflector := cache.NewReflector(listWatcher, &extensionsv1.DaemonSet{}, store, time.Hour)
	go reflector.Run(stopchannel)
	return &DaemonSetListerImpl{daemonSetLister: lister}
}
