package azure

import (
 "encoding/json"
 "fmt"
 "io/ioutil"
 "net/http"
 "regexp"
 "strconv"
 "strings"
 "github.com/Azure/go-autorest/autorest"
)

const (
 HeaderClientID       = "x-ms-client-request-id"
 HeaderReturnClientID = "x-ms-return-client-request-id"
 HeaderRequestID      = "x-ms-request-id"
)

type ServiceError struct {
 Code           string                   `json:"code"`
 Message        string                   `json:"message"`
 Target         *string                  `json:"target"`
 Details        []map[string]interface{} `json:"details"`
 InnerError     map[string]interface{}   `json:"innererror"`
 AdditionalInfo []map[string]interface{} `json:"additionalInfo"`
}

func (se ServiceError) Error() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := fmt.Sprintf("Code=%q Message=%q", se.Code, se.Message)
 if se.Target != nil {
  result += fmt.Sprintf(" Target=%q", *se.Target)
 }
 if se.Details != nil {
  d, err := json.Marshal(se.Details)
  if err != nil {
   result += fmt.Sprintf(" Details=%v", se.Details)
  }
  result += fmt.Sprintf(" Details=%v", string(d))
 }
 if se.InnerError != nil {
  d, err := json.Marshal(se.InnerError)
  if err != nil {
   result += fmt.Sprintf(" InnerError=%v", se.InnerError)
  }
  result += fmt.Sprintf(" InnerError=%v", string(d))
 }
 if se.AdditionalInfo != nil {
  d, err := json.Marshal(se.AdditionalInfo)
  if err != nil {
   result += fmt.Sprintf(" AdditionalInfo=%v", se.AdditionalInfo)
  }
  result += fmt.Sprintf(" AdditionalInfo=%v", string(d))
 }
 return result
}
func (se *ServiceError) UnmarshalJSON(b []byte) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 type serviceError1 struct {
  Code           string                   `json:"code"`
  Message        string                   `json:"message"`
  Target         *string                  `json:"target"`
  Details        []map[string]interface{} `json:"details"`
  InnerError     map[string]interface{}   `json:"innererror"`
  AdditionalInfo []map[string]interface{} `json:"additionalInfo"`
 }
 type serviceError2 struct {
  Code           string                   `json:"code"`
  Message        string                   `json:"message"`
  Target         *string                  `json:"target"`
  Details        map[string]interface{}   `json:"details"`
  InnerError     map[string]interface{}   `json:"innererror"`
  AdditionalInfo []map[string]interface{} `json:"additionalInfo"`
 }
 se1 := serviceError1{}
 err := json.Unmarshal(b, &se1)
 if err == nil {
  se.populate(se1.Code, se1.Message, se1.Target, se1.Details, se1.InnerError, se1.AdditionalInfo)
  return nil
 }
 se2 := serviceError2{}
 err = json.Unmarshal(b, &se2)
 if err == nil {
  se.populate(se2.Code, se2.Message, se2.Target, nil, se2.InnerError, se2.AdditionalInfo)
  se.Details = append(se.Details, se2.Details)
  return nil
 }
 return err
}
func (se *ServiceError) populate(code, message string, target *string, details []map[string]interface{}, inner map[string]interface{}, additional []map[string]interface{}) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 se.Code = code
 se.Message = message
 se.Target = target
 se.Details = details
 se.InnerError = inner
 se.AdditionalInfo = additional
}

type RequestError struct {
 autorest.DetailedError
 ServiceError *ServiceError `json:"error"`
 RequestID    string
}

func (e RequestError) Error() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("autorest/azure: Service returned an error. Status=%v %v", e.StatusCode, e.ServiceError)
}
func IsAzureError(e error) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 _, ok := e.(*RequestError)
 return ok
}

type Resource struct {
 SubscriptionID string
 ResourceGroup  string
 Provider       string
 ResourceType   string
 ResourceName   string
}

func ParseResourceID(resourceID string) (Resource, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 const resourceIDPatternText = `(?i)subscriptions/(.+)/resourceGroups/(.+)/providers/(.+?)/(.+?)/(.+)`
 resourceIDPattern := regexp.MustCompile(resourceIDPatternText)
 match := resourceIDPattern.FindStringSubmatch(resourceID)
 if len(match) == 0 {
  return Resource{}, fmt.Errorf("parsing failed for %s. Invalid resource Id format", resourceID)
 }
 v := strings.Split(match[5], "/")
 resourceName := v[len(v)-1]
 result := Resource{SubscriptionID: match[1], ResourceGroup: match[2], Provider: match[3], ResourceType: match[4], ResourceName: resourceName}
 return result, nil
}
func NewErrorWithError(original error, packageType string, method string, resp *http.Response, message string, args ...interface{}) RequestError {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if v, ok := original.(*RequestError); ok {
  return *v
 }
 statusCode := autorest.UndefinedStatusCode
 if resp != nil {
  statusCode = resp.StatusCode
 }
 return RequestError{DetailedError: autorest.DetailedError{Original: original, PackageType: packageType, Method: method, StatusCode: statusCode, Message: fmt.Sprintf(message, args...)}}
}
func WithReturningClientID(uuid string) autorest.PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 preparer := autorest.CreatePreparer(WithClientID(uuid), WithReturnClientID(true))
 return func(p autorest.Preparer) autorest.Preparer {
  return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err != nil {
    return r, err
   }
   return preparer.Prepare(r)
  })
 }
}
func WithClientID(uuid string) autorest.PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return autorest.WithHeader(HeaderClientID, uuid)
}
func WithReturnClientID(b bool) autorest.PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return autorest.WithHeader(HeaderReturnClientID, strconv.FormatBool(b))
}
func ExtractClientID(resp *http.Response) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return autorest.ExtractHeaderValue(HeaderClientID, resp)
}
func ExtractRequestID(resp *http.Response) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return autorest.ExtractHeaderValue(HeaderRequestID, resp)
}
func WithErrorUnlessStatusCode(codes ...int) autorest.RespondDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(r autorest.Responder) autorest.Responder {
  return autorest.ResponderFunc(func(resp *http.Response) error {
   err := r.Respond(resp)
   if err == nil && !autorest.ResponseHasStatusCode(resp, codes...) {
    var e RequestError
    defer resp.Body.Close()
    b, decodeErr := autorest.CopyAndDecode(autorest.EncodedAsJSON, resp.Body, &e)
    resp.Body = ioutil.NopCloser(&b)
    if decodeErr != nil {
     return fmt.Errorf("autorest/azure: error response cannot be parsed: %q error: %v", b.String(), decodeErr)
    }
    if e.ServiceError == nil {
     if err := json.Unmarshal(b.Bytes(), &e.ServiceError); err != nil {
      return err
     }
    }
    if e.ServiceError.Message == "" {
     rawBody := map[string]interface{}{}
     if err := json.Unmarshal(b.Bytes(), &rawBody); err != nil {
      return err
     }
     e.ServiceError = &ServiceError{Code: "Unknown", Message: "Unknown service error"}
     if len(rawBody) > 0 {
      e.ServiceError.Details = []map[string]interface{}{rawBody}
     }
    }
    e.Response = resp
    e.RequestID = ExtractRequestID(resp)
    if e.StatusCode == nil {
     e.StatusCode = resp.StatusCode
    }
    err = &e
   }
   return err
  })
 }
}
