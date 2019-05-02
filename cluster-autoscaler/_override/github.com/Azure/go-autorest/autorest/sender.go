package autorest

import (
 "fmt"
 "log"
 "math"
 "net/http"
 "strconv"
 "time"
)

type Sender interface {
 Do(*http.Request) (*http.Response, error)
}
type SenderFunc func(*http.Request) (*http.Response, error)

func (sf SenderFunc) Do(r *http.Request) (*http.Response, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return sf(r)
}

type SendDecorator func(Sender) Sender

func CreateSender(decorators ...SendDecorator) Sender {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return DecorateSender(&http.Client{}, decorators...)
}
func DecorateSender(s Sender, decorators ...SendDecorator) Sender {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, decorate := range decorators {
  s = decorate(s)
 }
 return s
}
func Send(r *http.Request, decorators ...SendDecorator) (*http.Response, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return SendWithSender(&http.Client{}, r, decorators...)
}
func SendWithSender(s Sender, r *http.Request, decorators ...SendDecorator) (*http.Response, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return DecorateSender(s, decorators...).Do(r)
}
func AfterDelay(d time.Duration) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   if !DelayForBackoff(d, 0, r.Context().Done()) {
    return nil, fmt.Errorf("autorest: AfterDelay canceled before full delay")
   }
   return s.Do(r)
  })
 }
}
func AsIs() SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   return s.Do(r)
  })
 }
}
func DoCloseIfError() SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   resp, err := s.Do(r)
   if err != nil {
    Respond(resp, ByDiscardingBody(), ByClosing())
   }
   return resp, err
  })
 }
}
func DoErrorIfStatusCode(codes ...int) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   resp, err := s.Do(r)
   if err == nil && ResponseHasStatusCode(resp, codes...) {
    err = NewErrorWithResponse("autorest", "DoErrorIfStatusCode", resp, "%v %v failed with %s", resp.Request.Method, resp.Request.URL, resp.Status)
   }
   return resp, err
  })
 }
}
func DoErrorUnlessStatusCode(codes ...int) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   resp, err := s.Do(r)
   if err == nil && !ResponseHasStatusCode(resp, codes...) {
    err = NewErrorWithResponse("autorest", "DoErrorUnlessStatusCode", resp, "%v %v failed with %s", resp.Request.Method, resp.Request.URL, resp.Status)
   }
   return resp, err
  })
 }
}
func DoPollForStatusCodes(duration time.Duration, delay time.Duration, codes ...int) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
   resp, err = s.Do(r)
   if err == nil && ResponseHasStatusCode(resp, codes...) {
    r, err = NewPollingRequestWithContext(r.Context(), resp)
    for err == nil && ResponseHasStatusCode(resp, codes...) {
     Respond(resp, ByDiscardingBody(), ByClosing())
     resp, err = SendWithSender(s, r, AfterDelay(GetRetryAfter(resp, delay)))
    }
   }
   return resp, err
  })
 }
}
func DoRetryForAttempts(attempts int, backoff time.Duration) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
   rr := NewRetriableRequest(r)
   for attempt := 0; attempt < attempts; attempt++ {
    err = rr.Prepare()
    if err != nil {
     return resp, err
    }
    resp, err = s.Do(rr.Request())
    if err == nil {
     return resp, err
    }
    if !DelayForBackoff(backoff, attempt, r.Context().Done()) {
     return nil, r.Context().Err()
    }
   }
   return resp, err
  })
 }
}
func DoRetryForStatusCodes(attempts int, backoff time.Duration, codes ...int) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
   rr := NewRetriableRequest(r)
   attempts++
   for attempt := 0; attempt < attempts; {
    err = rr.Prepare()
    if err != nil {
     return resp, err
    }
    resp, err = s.Do(rr.Request())
    if err != nil && !IsTemporaryNetworkError(err) {
     return nil, err
    }
    if err == nil && !ResponseHasStatusCode(resp, codes...) || IsTokenRefreshError(err) {
     return resp, err
    }
    delayed := DelayWithRetryAfter(resp, r.Context().Done())
    if !delayed && !DelayForBackoff(backoff, attempt, r.Context().Done()) {
     return resp, r.Context().Err()
    }
    if resp == nil || resp.StatusCode != http.StatusTooManyRequests {
     attempt++
    }
   }
   return resp, err
  })
 }
}
func DelayWithRetryAfter(resp *http.Response, cancel <-chan struct{}) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if resp == nil {
  return false
 }
 retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
 if resp.StatusCode == http.StatusTooManyRequests && retryAfter > 0 {
  select {
  case <-time.After(time.Duration(retryAfter) * time.Second):
   return true
  case <-cancel:
   return false
  }
 }
 return false
}
func DoRetryForDuration(d time.Duration, backoff time.Duration) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
   rr := NewRetriableRequest(r)
   end := time.Now().Add(d)
   for attempt := 0; time.Now().Before(end); attempt++ {
    err = rr.Prepare()
    if err != nil {
     return resp, err
    }
    resp, err = s.Do(rr.Request())
    if err == nil {
     return resp, err
    }
    if !DelayForBackoff(backoff, attempt, r.Context().Done()) {
     return nil, r.Context().Err()
    }
   }
   return resp, err
  })
 }
}
func WithLogging(logger *log.Logger) SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s Sender) Sender {
  return SenderFunc(func(r *http.Request) (*http.Response, error) {
   logger.Printf("Sending %s %s", r.Method, r.URL)
   resp, err := s.Do(r)
   if err != nil {
    logger.Printf("%s %s received error '%v'", r.Method, r.URL, err)
   } else {
    logger.Printf("%s %s received %s", r.Method, r.URL, resp.Status)
   }
   return resp, err
  })
 }
}
func DelayForBackoff(backoff time.Duration, attempt int, cancel <-chan struct{}) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 select {
 case <-time.After(time.Duration(backoff.Seconds()*math.Pow(2, float64(attempt))) * time.Second):
  return true
 case <-cancel:
  return false
 }
}
