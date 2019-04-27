package nanny

import (
	"time"
	log "github.com/golang/glog"
	api "k8s.io/kubernetes/pkg/api/v1"
)

func checkResource(estimatorResult *EstimatorResult, actual api.ResourceList, res api.ResourceName) *api.ResourceList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	val, ok := actual[res]
	expMinVal, expMinOk := estimatorResult.AcceptableRange.lower[res]
	expMaxVal, expMaxOk := estimatorResult.AcceptableRange.upper[res]
	if ok != expMinOk || ok != expMaxOk {
		return &estimatorResult.RecommendedRange.lower
	}
	if !ok && !expMinOk && !expMaxOk {
		return nil
	}
	if val.Cmp(expMinVal) == -1 {
		return &estimatorResult.RecommendedRange.lower
	}
	if val.Cmp(expMaxVal) == 1 {
		return &estimatorResult.RecommendedRange.upper
	}
	return nil
}
func shouldOverwriteResources(estimatorResult *EstimatorResult, limits, reqs api.ResourceList) *api.ResourceRequirements {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, list := range []api.ResourceList{limits, reqs} {
		for _, resourceType := range []api.ResourceName{api.ResourceCPU, api.ResourceMemory, api.ResourceStorage} {
			newReqs := checkResource(estimatorResult, list, resourceType)
			if newReqs != nil {
				log.V(4).Infof("Resource %s is out of bounds.", resourceType)
				return &api.ResourceRequirements{Limits: *newReqs, Requests: *newReqs}
			}
		}
	}
	return nil
}

type KubernetesClient interface {
	CountNodes() (uint64, error)
	ContainerResources() (*api.ResourceRequirements, error)
	UpdateDeployment(resources *api.ResourceRequirements) error
}
type ResourceEstimator interface {
	scaleWithNodes(numNodes uint64) *EstimatorResult
}

func PollAPIServer(k8s KubernetesClient, est ResourceEstimator, pollPeriod time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; true; i++ {
		if i != 0 {
			time.Sleep(pollPeriod)
		}
		num, err := k8s.CountNodes()
		if err != nil {
			log.Error(err)
			continue
		}
		log.V(4).Infof("The number of nodes is %d", num)
		resources, err := k8s.ContainerResources()
		if err != nil {
			log.Errorf("Error while querying apiserver for resources: %v", err)
			continue
		}
		estimation := est.scaleWithNodes(num)
		overwrite := shouldOverwriteResources(estimation, resources.Limits, resources.Requests)
		if overwrite == nil {
			log.V(4).Infof("Resources are within the expected limits. Actual: %+v, accepted range: %+v", *resources, estimation.AcceptableRange)
			continue
		}
		log.Infof("Resources are not within the expected limits, updating the deployment. Actual: %+v New: %+v", *resources, *overwrite)
		if err := k8s.UpdateDeployment(overwrite); err != nil {
			log.Error(err)
			continue
		}
	}
}
