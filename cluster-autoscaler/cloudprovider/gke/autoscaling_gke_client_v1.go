package gke

import (
	"fmt"
	"net/http"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	gke_api "google.golang.org/api/container/v1"
)

type autoscalingGkeClientV1 struct {
	gkeService	*gke_api.Service
	clusterPath	string
	operationPath	string
}

func NewAutoscalingGkeClientV1(client *http.Client, projectId, location, clusterName string) (*autoscalingGkeClientV1, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	autoscalingGkeClient := &autoscalingGkeClientV1{clusterPath: fmt.Sprintf(clusterPathPrefix, projectId, location, clusterName), operationPath: fmt.Sprintf(operationPathPrefix, projectId, location)}
	gkeService, err := gke_api.New(client)
	if err != nil {
		return nil, err
	}
	if *GkeAPIEndpoint != "" {
		gkeService.BasePath = *GkeAPIEndpoint
	}
	autoscalingGkeClient.gkeService = gkeService
	return autoscalingGkeClient, nil
}
func (m *autoscalingGkeClientV1) GetCluster() (Cluster, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("clusters", "get")
	clusterResponse, err := m.gkeService.Projects.Locations.Clusters.Get(m.clusterPath).Do()
	if err != nil {
		return Cluster{}, err
	}
	nodePools := []NodePool{}
	for _, pool := range clusterResponse.NodePools {
		if pool.Autoscaling != nil && pool.Autoscaling.Enabled {
			nodePools = append(nodePools, NodePool{Name: pool.Name, InstanceGroupUrls: pool.InstanceGroupUrls, Autoscaled: pool.Autoscaling.Enabled, MinNodeCount: pool.Autoscaling.MinNodeCount, MaxNodeCount: pool.Autoscaling.MaxNodeCount})
		}
	}
	return Cluster{Locations: clusterResponse.Locations, NodePools: nodePools}, nil
}
func (m *autoscalingGkeClientV1) DeleteNodePool(toBeRemoved string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
func (m *autoscalingGkeClientV1) CreateNodePool(mig *GkeMig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
