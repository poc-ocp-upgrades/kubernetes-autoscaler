package alicloud

import (
	"github.com/stretchr/testify/assert"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	"testing"
)

func TestBuildGenericLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	template := &sgTemplate{InstanceType: &instanceType{instanceTypeID: "gn5-4c-8g", vcpu: 4, memoryInBytes: 8 * 1024 * 1024 * 1024, gpu: 1}, Region: "cn-hangzhou", Zone: "cn-hangzhou-a"}
	nodeName := "virtual-node"
	labels := buildGenericLabels(template, nodeName)
	assert.Equal(t, labels[kubeletapis.LabelInstanceType], template.InstanceType.instanceTypeID)
}
