package gce

import (
	"net/http"
	"testing"
	"time"
	test_util "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gce_api "google.golang.org/api/compute/v1"
)

func newTestAutoscalingGceClient(t *testing.T, projectId, url string) *autoscalingGceClientV1 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := &http.Client{}
	gceClient, err := NewAutoscalingGceClientV1(client, projectId)
	if !assert.NoError(t, err) {
		t.Fatalf("fatal error: %v", err)
	}
	gceClient.gceService.BasePath = url
	return gceClient
}

const operationRunningResponse = `{
  "name": "operation-1505728466148-d16f5197",
  "zone": "us-central1-a",
  "operationType": "CREATE_NODE_POOL",
  "status": "RUNNING",
  "selfLink": "https://container.googleapis.com/v1/projects/601024681890/locations/us-central1-a/operations/operation-1505728466148-d16f5197",
  "targetLink": "https://container.googleapis.com/v1/projects/601024681890/locations/us-central1-a/clusters/cluster-1/nodePools/nodeautoprovisioning-323233232",
  "startTime": "2017-09-18T09:54:26.148507311Z",
  "endTime": "2017-09-18T09:54:35.124878859Z"
}`
const operationDoneResponse = `{
  "name": "operation-1505728466148-d16f5197",
  "zone": "us-central1-a",
  "operationType": "CREATE_NODE_POOL",
  "status": "DONE",
  "selfLink": "https://container.googleapis.com/v1/projects/601024681890/locations/us-central1-a/operations/operation-1505728466148-d16f5197",
  "targetLink": "https://container.googleapis.com/v1/projects/601024681890/locations/us-central1-a/clusters/cluster-1/nodePools/nodeautoprovisioning-323233232",
  "startTime": "2017-09-18T09:54:26.148507311Z",
  "endTime": "2017-09-18T09:54:35.124878859Z"
}`

func TestWaitForOp(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	server := test_util.NewHttpServerMock()
	defer server.Close()
	g := newTestAutoscalingGceClient(t, "project1", server.URL)
	g.operationPollInterval = 1 * time.Millisecond
	g.operationWaitTimeout = 500 * time.Millisecond
	server.On("handle", "/project1/zones/us-central1-b/operations/operation-1505728466148-d16f5197").Return(operationRunningResponse).Times(3)
	server.On("handle", "/project1/zones/us-central1-b/operations/operation-1505728466148-d16f5197").Return(operationDoneResponse).Once()
	operation := &gce_api.Operation{Name: "operation-1505728466148-d16f5197"}
	err := g.waitForOp(operation, projectId, zoneB)
	assert.NoError(t, err)
	mock.AssertExpectationsForObjects(t, server)
}
func TestWaitForOpTimeout(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	server := test_util.NewHttpServerMock()
	defer server.Close()
	g := newTestAutoscalingGceClient(t, "project1", server.URL)
	g.operationPollInterval = 10 * time.Millisecond
	g.operationWaitTimeout = 49 * time.Millisecond
	server.On("handle", "/project1/zones/us-central1-b/operations/operation-1505728466148-d16f5197").Return(operationRunningResponse).Times(4)
	operation := &gce_api.Operation{Name: "operation-1505728466148-d16f5197"}
	err := g.waitForOp(operation, projectId, zoneB)
	assert.Error(t, err)
}
