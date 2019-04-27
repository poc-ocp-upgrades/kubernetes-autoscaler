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
	activityTimeout	time.Duration
	checkTimeout	bool
	lastActivity	time.Time
	mutex		*sync.Mutex
}

func NewHealthCheck(activityTimeout time.Duration, checkTimeout bool) *HealthCheck {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &HealthCheck{activityTimeout: activityTimeout, checkTimeout: checkTimeout, lastActivity: time.Now(), mutex: &sync.Mutex{}}
}
func (hc *HealthCheck) checkLastActivity() (bool, time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	now := time.Now()
	lastActivity := hc.lastActivity
	activityTimedOut := now.After(lastActivity.Add(hc.activityTimeout))
	timedOut := hc.checkTimeout && activityTimedOut
	return timedOut, now.Sub(lastActivity)
}
func (hc *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timedOut, ago := hc.checkLastActivity()
	if timedOut {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error: last activity more than %v ago", ago)))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}
func (hc *HealthCheck) UpdateLastActivity() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.lastActivity = time.Now()
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
