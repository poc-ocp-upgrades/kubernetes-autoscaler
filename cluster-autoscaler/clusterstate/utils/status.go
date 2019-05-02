package utils

import (
 "errors"
 "fmt"
 "time"
 apiv1 "k8s.io/api/core/v1"
 kube_errors "k8s.io/apimachinery/pkg/api/errors"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/apimachinery/pkg/runtime"
 kube_client "k8s.io/client-go/kubernetes"
 "k8s.io/client-go/tools/record"
 "k8s.io/klog"
)

const (
 StatusConfigMapName       = "cluster-autoscaler-status"
 ConfigMapLastUpdatedKey   = "cluster-autoscaler.kubernetes.io/last-updated"
 ConfigMapLastUpdateFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

type LogEventRecorder struct {
 recorder     record.EventRecorder
 statusObject runtime.Object
 active       bool
}

func (ler *LogEventRecorder) Event(eventtype, reason, message string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if ler.active && ler.statusObject != nil {
  ler.recorder.Event(ler.statusObject, eventtype, reason, message)
 }
}
func (ler *LogEventRecorder) Eventf(eventtype, reason, message string, args ...interface{}) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if ler.active && ler.statusObject != nil {
  ler.recorder.Eventf(ler.statusObject, eventtype, reason, message, args...)
 }
}
func NewStatusMapRecorder(kubeClient kube_client.Interface, namespace string, recorder record.EventRecorder, active bool) (*LogEventRecorder, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var mapObj runtime.Object
 var err error
 if active {
  mapObj, err = WriteStatusConfigMap(kubeClient, namespace, "Initializing", nil)
  if err != nil {
   return nil, errors.New("Failed to init status ConfigMap")
  }
 }
 return &LogEventRecorder{recorder: recorder, statusObject: mapObj, active: active}, nil
}
func WriteStatusConfigMap(kubeClient kube_client.Interface, namespace string, msg string, logRecorder *LogEventRecorder) (*apiv1.ConfigMap, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 statusUpdateTime := time.Now().Format(ConfigMapLastUpdateFormat)
 statusMsg := fmt.Sprintf("Cluster-autoscaler status at %s:\n%v", statusUpdateTime, msg)
 var configMap *apiv1.ConfigMap
 var getStatusError, writeStatusError error
 var errMsg string
 maps := kubeClient.CoreV1().ConfigMaps(namespace)
 configMap, getStatusError = maps.Get(StatusConfigMapName, metav1.GetOptions{})
 if getStatusError == nil {
  configMap.Data["status"] = statusMsg
  if configMap.ObjectMeta.Annotations == nil {
   configMap.ObjectMeta.Annotations = make(map[string]string)
  }
  configMap.ObjectMeta.Annotations[ConfigMapLastUpdatedKey] = statusUpdateTime
  configMap, writeStatusError = maps.Update(configMap)
 } else if kube_errors.IsNotFound(getStatusError) {
  configMap = &apiv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: StatusConfigMapName, Annotations: map[string]string{ConfigMapLastUpdatedKey: statusUpdateTime}}, Data: map[string]string{"status": statusMsg}}
  configMap, writeStatusError = maps.Create(configMap)
 } else {
  errMsg = fmt.Sprintf("Failed to retrieve status configmap for update: %v", getStatusError)
 }
 if writeStatusError != nil {
  errMsg = fmt.Sprintf("Failed to write status configmap: %v", writeStatusError)
 }
 if errMsg != "" {
  klog.Error(errMsg)
  return nil, errors.New(errMsg)
 }
 klog.V(8).Infof("Successfully wrote status configmap with body \"%v\"", statusMsg)
 if logRecorder != nil {
  logRecorder.statusObject = configMap
 }
 return configMap, nil
}
func DeleteStatusConfigMap(kubeClient kube_client.Interface, namespace string) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 maps := kubeClient.CoreV1().ConfigMaps(namespace)
 err := maps.Delete(StatusConfigMapName, &metav1.DeleteOptions{})
 if err != nil {
  klog.Error("Failed to delete status configmap")
 }
 return err
}
