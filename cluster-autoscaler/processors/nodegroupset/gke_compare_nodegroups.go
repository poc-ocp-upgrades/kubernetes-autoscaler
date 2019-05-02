package nodegroupset

import (
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const GkeNodepoolLabel = "cloud.google.com/gke-nodepool"

func nodesFromSameGkeNodePool(n1, n2 *schedulercache.NodeInfo) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 n1GkeNodePool := n1.Node().Labels[GkeNodepoolLabel]
 n2GkeNodePool := n2.Node().Labels[GkeNodepoolLabel]
 return n1GkeNodePool != "" && n1GkeNodePool == n2GkeNodePool
}
func IsGkeNodeInfoSimilar(n1, n2 *schedulercache.NodeInfo) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if nodesFromSameGkeNodePool(n1, n2) {
  return true
 }
 return IsNodeInfoSimilar(n1, n2)
}
