package gke

import (
	"net/http"
	"testing"
	"time"
	test_util "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gke_api_beta "google.golang.org/api/container/v1beta1"
)

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

func newTestAutoscalingGkeClientV1beta1(t *testing.T, project, location, clusterName, url string) *autoscalingGkeClientV1beta1 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*GkeAPIEndpoint = url
	client := &http.Client{}
	gkeClient, err := NewAutoscalingGkeClientV1beta1(client, project, location, clusterName)
	if !assert.NoError(t, err) {
		t.Fatalf("fatal error: %v", err)
	}
	return gkeClient
}
func TestWaitForGkeOp(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	server := test_util.NewHttpServerMock()
	defer server.Close()
	g := newTestAutoscalingGkeClientV1beta1(t, "project1", "us-central1-b", "cluster-1", server.URL)
	g.operationPollInterval = 1 * time.Millisecond
	g.operationWaitTimeout = 500 * time.Millisecond
	server.On("handle", "/v1beta1/projects/project1/locations/us-central1-b/operations/operation-1505728466148-d16f5197").Return(operationRunningResponse).Once()
	server.On("handle", "/v1beta1/projects/project1/locations/us-central1-b/operations/operation-1505728466148-d16f5197").Return(operationDoneResponse).Once()
	operation := &gke_api_beta.Operation{Name: "operation-1505728466148-d16f5197"}
	err := g.waitForGkeOp(operation)
	assert.NoError(t, err)
	mock.AssertExpectationsForObjects(t, server)
}
