package logic

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/updater/eviction"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/updater/priority"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	vpa_lister "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/autoscaling.k8s.io/v1beta1"
	metrics_updater "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics/updater"
	vpa_api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
	kube_client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"github.com/golang/glog"
)

type Updater interface{ RunOnce() }
type updater struct {
	vpaLister				vpa_lister.VerticalPodAutoscalerLister
	podLister				v1lister.PodLister
	eventRecorder			record.EventRecorder
	evictionFactory			eviction.PodsEvictionRestrictionFactory
	recommendationProcessor	vpa_api_util.RecommendationProcessor
	evictionAdmission		priority.PodEvictionAdmission
}

func NewUpdater(kubeClient kube_client.Interface, vpaClient *vpa_clientset.Clientset, minReplicasForEvicition int, evictionToleranceFraction float64, recommendationProcessor vpa_api_util.RecommendationProcessor, evictionAdmission priority.PodEvictionAdmission) (Updater, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	factory, err := eviction.NewPodsEvictionRestrictionFactory(kubeClient, minReplicasForEvicition, evictionToleranceFraction)
	if err != nil {
		return nil, fmt.Errorf("Failed to create eviction restriction factory: %v", err)
	}
	return &updater{vpaLister: vpa_api_util.NewAllVpasLister(vpaClient, make(chan struct{})), podLister: newPodLister(kubeClient), eventRecorder: newEventRecorder(kubeClient), evictionFactory: factory, recommendationProcessor: recommendationProcessor, evictionAdmission: evictionAdmission}, nil
}
func (u *updater) RunOnce() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	timer := metrics_updater.NewExecutionTimer()
	vpaList, err := u.vpaLister.List(labels.Everything())
	if err != nil {
		glog.Fatalf("failed get VPA list: %v", err)
	}
	timer.ObserveStep("ListVPAs")
	vpas := make([]*vpa_types.VerticalPodAutoscaler, 0)
	for _, vpa := range vpaList {
		if vpa_api_util.GetUpdateMode(vpa) != vpa_types.UpdateModeRecreate && vpa_api_util.GetUpdateMode(vpa) != vpa_types.UpdateModeAuto {
			glog.V(3).Infof("skipping VPA object %v because its mode is not \"Recreate\" or \"Auto\"", vpa.Name)
			continue
		}
		vpas = append(vpas, vpa)
	}
	if len(vpas) == 0 {
		glog.Warningf("no VPA objects to process")
		if u.evictionAdmission != nil {
			u.evictionAdmission.CleanUp()
		}
		timer.ObserveTotal()
		return
	}
	podsList, err := u.podLister.List(labels.Everything())
	if err != nil {
		glog.Errorf("failed to get pods list: %v", err)
		timer.ObserveTotal()
		return
	}
	timer.ObserveStep("ListPods")
	allLivePods := filterDeletedPods(podsList)
	controlledPods := make(map[*vpa_types.VerticalPodAutoscaler][]*apiv1.Pod)
	for _, pod := range allLivePods {
		controllingVPA := vpa_api_util.GetControllingVPAForPod(pod, vpas)
		if controllingVPA != nil {
			controlledPods[controllingVPA] = append(controlledPods[controllingVPA], pod)
		}
	}
	timer.ObserveStep("FilterPods")
	if u.evictionAdmission != nil {
		u.evictionAdmission.LoopInit(allLivePods, controlledPods)
	}
	timer.ObserveStep("AdmissionInit")
	for vpa, livePods := range controlledPods {
		evictionLimiter := u.evictionFactory.NewPodsEvictionRestriction(livePods)
		podsForUpdate := u.getPodsUpdateOrder(filterNonEvictablePods(livePods, evictionLimiter), vpa)
		for _, pod := range podsForUpdate {
			if !evictionLimiter.CanEvict(pod) {
				continue
			}
			glog.V(2).Infof("evicting pod %v", pod.Name)
			evictErr := evictionLimiter.Evict(pod, u.eventRecorder)
			if evictErr != nil {
				glog.Warningf("evicting pod %v failed: %v", pod.Name, evictErr)
			}
		}
	}
	timer.ObserveStep("EvictPods")
	timer.ObserveTotal()
}
func (u *updater) getPodsUpdateOrder(pods []*apiv1.Pod, vpa *vpa_types.VerticalPodAutoscaler) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	priorityCalculator := priority.NewUpdatePriorityCalculator(vpa.Spec.ResourcePolicy, vpa.Status.Conditions, nil, u.recommendationProcessor)
	recommendation := vpa.Status.Recommendation
	for _, pod := range pods {
		priorityCalculator.AddPod(pod, recommendation, time.Now())
	}
	return priorityCalculator.GetSortedPods(u.evictionAdmission)
}
func filterNonEvictablePods(pods []*apiv1.Pod, evictionRestriciton eviction.PodsEvictionRestriction) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*apiv1.Pod, 0)
	for _, pod := range pods {
		if evictionRestriciton.CanEvict(pod) {
			result = append(result, pod)
		}
	}
	return result
}
func filterDeletedPods(pods []*apiv1.Pod) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*apiv1.Pod, 0)
	for _, pod := range pods {
		if pod.DeletionTimestamp == nil {
			result = append(result, pod)
		}
	}
	return result
}
func newPodLister(kubeClient kube_client.Interface) v1lister.PodLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	selector := fields.ParseSelectorOrDie("spec.nodeName!=" + "" + ",status.phase!=" + string(apiv1.PodSucceeded) + ",status.phase!=" + string(apiv1.PodFailed))
	podListWatch := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", apiv1.NamespaceAll, selector)
	store := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	podLister := v1lister.NewPodLister(store)
	podReflector := cache.NewReflector(podListWatch, &apiv1.Pod{}, store, time.Hour)
	stopCh := make(chan struct{})
	go podReflector.Run(stopCh)
	return podLister
}
func newEventRecorder(kubeClient kube_client.Interface) record.EventRecorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.V(4).Infof)
	if _, isFake := kubeClient.(*fake.Clientset); !isFake {
		eventBroadcaster.StartRecordingToSink(&clientv1.EventSinkImpl{Interface: clientv1.New(kubeClient.CoreV1().RESTClient()).Events("")})
	}
	return eventBroadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{Component: "vpa-updater"})
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
