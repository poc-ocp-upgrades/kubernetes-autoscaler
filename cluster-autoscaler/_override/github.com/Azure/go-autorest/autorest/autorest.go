package autorest

import (
 "context"
 "net/http"
 "time"
)

const (
 HeaderLocation   = "Location"
 HeaderRetryAfter = "Retry-After"
)

func ResponseHasStatusCode(resp *http.Response, codes ...int) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if resp == nil {
  return false
 }
 return containsInt(codes, resp.StatusCode)
}
func GetLocation(resp *http.Response) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return resp.Header.Get(HeaderLocation)
}
func GetRetryAfter(resp *http.Response, defaultDelay time.Duration) time.Duration {
 _logClusterCodePath()
 defer _logClusterCodePath()
 retry := resp.Header.Get(HeaderRetryAfter)
 if retry == "" {
  return defaultDelay
 }
 d, err := time.ParseDuration(retry + "s")
 if err != nil {
  return defaultDelay
 }
 return d
}
func NewPollingRequest(resp *http.Response, cancel <-chan struct{}) (*http.Request, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 location := GetLocation(resp)
 if location == "" {
  return nil, NewErrorWithResponse("autorest", "NewPollingRequest", resp, "Location header missing from response that requires polling")
 }
 req, err := Prepare(&http.Request{Cancel: cancel}, AsGet(), WithBaseURL(location))
 if err != nil {
  return nil, NewErrorWithError(err, "autorest", "NewPollingRequest", nil, "Failure creating poll request to %s", location)
 }
 return req, nil
}
func NewPollingRequestWithContext(ctx context.Context, resp *http.Response) (*http.Request, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 location := GetLocation(resp)
 if location == "" {
  return nil, NewErrorWithResponse("autorest", "NewPollingRequestWithContext", resp, "Location header missing from response that requires polling")
 }
 req, err := Prepare((&http.Request{}).WithContext(ctx), AsGet(), WithBaseURL(location))
 if err != nil {
  return nil, NewErrorWithError(err, "autorest", "NewPollingRequestWithContext", nil, "Failure creating poll request to %s", location)
 }
 return req, nil
}
