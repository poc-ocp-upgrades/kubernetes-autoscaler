package metrics

import (
	"net/http/httptest"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func getTestResponse(start time.Time, activityTimeout, successTimeout time.Duration, checkMonitoring bool) *httptest.ResponseRecorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := httptest.NewRequest("GET", "/health-check", nil)
	w := httptest.NewRecorder()
	healthCheck := NewHealthCheck(activityTimeout, successTimeout)
	if checkMonitoring {
		healthCheck.StartMonitoring()
	}
	healthCheck.lastActivity = start
	healthCheck.lastSuccessfulRun = start
	healthCheck.ServeHTTP(w, req)
	return w
}
func TestOkServeHTTP(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	w := getTestResponse(time.Now(), time.Second, time.Second, true)
	assert.Equal(t, 200, w.Code)
}
func TestFailTimeoutServeHTTP(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	w := getTestResponse(time.Now().Add(time.Second*-2), time.Second, time.Second, true)
	assert.Equal(t, 500, w.Code)
}
func TestMonitoringOffAfterTimeout(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	w := getTestResponse(time.Now().Add(time.Second*-2), time.Second, time.Second, false)
	assert.Equal(t, 200, w.Code)
}
func TestMonitoringOffBeforeTimeout(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	w := getTestResponse(time.Now().Add(time.Second*2), time.Second, time.Second, false)
	assert.Equal(t, 200, w.Code)
}
func TestUpdateLastActivity(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timeout := time.Second
	start := time.Now().Add(timeout * -2)
	lastSuccess := time.Now().Add(timeout * 10)
	req := httptest.NewRequest("GET", "/health-check", nil)
	healthCheck := NewHealthCheck(timeout, timeout)
	healthCheck.StartMonitoring()
	healthCheck.lastActivity = start
	healthCheck.lastSuccessfulRun = lastSuccess
	w := httptest.NewRecorder()
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	w = httptest.NewRecorder()
	healthCheck.UpdateLastActivity(time.Now())
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
func TestUpdateActivityAtUpdateLastSuccessfulRun(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timeout := time.Second
	start := time.Now().Add(timeout * -2)
	lastSuccess := time.Now().Add(timeout * 10)
	req := httptest.NewRequest("GET", "/health-check", nil)
	healthCheck := NewHealthCheck(timeout, timeout)
	healthCheck.StartMonitoring()
	healthCheck.lastActivity = start
	healthCheck.lastSuccessfulRun = lastSuccess
	w := httptest.NewRecorder()
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	w = httptest.NewRecorder()
	healthCheck.UpdateLastSuccessfulRun(time.Now())
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, healthCheck.lastSuccessfulRun.After(healthCheck.lastActivity))
}
func TestUpdateLastSuccessfulRun(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timeout := time.Second
	start := time.Now().Add(timeout * -2)
	lastActivity := time.Now().Add(timeout * 10)
	req := httptest.NewRequest("GET", "/health-check", nil)
	healthCheck := NewHealthCheck(timeout, timeout)
	healthCheck.StartMonitoring()
	healthCheck.lastActivity = lastActivity
	healthCheck.lastSuccessfulRun = start
	w := httptest.NewRecorder()
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	w = httptest.NewRecorder()
	healthCheck.UpdateLastSuccessfulRun(time.Now())
	healthCheck.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, healthCheck.lastActivity.After(healthCheck.lastSuccessfulRun))
}
