package nanny

import (
 "fmt"
 "time"
 api "k8s.io/kubernetes/pkg/api"
 apiv1 "k8s.io/kubernetes/pkg/api/v1"
 cache "k8s.io/kubernetes/pkg/client/cache"
 client "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
 runtime "k8s.io/kubernetes/pkg/runtime"
 wait "k8s.io/kubernetes/pkg/util/wait"
 watch "k8s.io/kubernetes/pkg/watch"
)

type kubernetesClient struct {
 namespace  string
 deployment string
 pod        string
 container  string
 clientset  *client.Clientset
 nodeStore  cache.Store
 reflector  *cache.Reflector
}

func (k *kubernetesClient) CountNodes() (uint64, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 err := wait.PollImmediate(time.Second, time.Minute, func() (bool, error) {
  if k.reflector.LastSyncResourceVersion() == "" {
   return false, nil
  }
  return true, nil
 })
 if err != nil {
  return 0, err
 }
 return uint64(len(k.nodeStore.List())), nil
}
func (k *kubernetesClient) ContainerResources() (*apiv1.ResourceRequirements, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, err := k.clientset.CoreClient.Pods(k.namespace).Get(k.pod)
 if err != nil {
  return nil, err
 }
 for _, container := range pod.Spec.Containers {
  if container.Name == k.container {
   return &container.Resources, nil
  }
 }
 return nil, fmt.Errorf("Container %s was not found in deployment %s in namespace %s.", k.container, k.deployment, k.namespace)
}
func (k *kubernetesClient) UpdateDeployment(resources *apiv1.ResourceRequirements) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 dep, err := k.clientset.Extensions().Deployments(k.namespace).Get(k.deployment)
 if err != nil {
  return err
 }
 for i, container := range dep.Spec.Template.Spec.Containers {
  if container.Name == k.container {
   dep.Spec.Template.Spec.Containers[i].Resources = *resources
   _, err = k.clientset.ExtensionsClient.Deployments(k.namespace).Update(dep)
   return err
  }
 }
 return fmt.Errorf("Container %s was not found in the deployment %s in namespace %s.", k.container, k.deployment, k.namespace)
}
func NewKubernetesClient(namespace, deployment, pod, container string, clientset *client.Clientset) KubernetesClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := &kubernetesClient{namespace: namespace, deployment: deployment, pod: pod, container: container, clientset: clientset, nodeStore: cache.NewStore(cache.MetaNamespaceKeyFunc)}
 nodeListWatch := &cache.ListWatch{ListFunc: func(options api.ListOptions) (runtime.Object, error) {
  return clientset.Core().Nodes().List(options)
 }, WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
  return clientset.Core().Nodes().Watch(options)
 }}
 result.reflector = cache.NewReflector(nodeListWatch, &apiv1.Node{}, result.nodeStore, 0)
 result.reflector.Run()
 return result
}
