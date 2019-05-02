package errors

import (
 "encoding/json"
 "fmt"
 "github.com/jmespath/go-jmespath"
)

var wrapperList = []ServerErrorWrapper{&SignatureDostNotMatchWrapper{}}

type ServerError struct {
 httpStatus int
 requestId  string
 hostId     string
 errorCode  string
 recommend  string
 message    string
 comment    string
}
type ServerErrorWrapper interface {
 tryWrap(error *ServerError, wrapInfo map[string]string) (bool, *ServerError)
}

func (err *ServerError) Error() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("SDK.ServerError\nErrorCode: %s\nRecommend: %s\nRequestId: %s\nMessage: %s", err.errorCode, err.comment+err.recommend, err.requestId, err.message)
}
func NewServerError(httpStatus int, responseContent, comment string) Error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := &ServerError{httpStatus: httpStatus, message: responseContent, comment: comment}
 var data interface{}
 err := json.Unmarshal([]byte(responseContent), &data)
 if err == nil {
  requestId, _ := jmespath.Search("RequestId", data)
  hostId, _ := jmespath.Search("HostId", data)
  errorCode, _ := jmespath.Search("Code", data)
  recommend, _ := jmespath.Search("Recommend", data)
  message, _ := jmespath.Search("Message", data)
  if requestId != nil {
   result.requestId = requestId.(string)
  }
  if hostId != nil {
   result.hostId = hostId.(string)
  }
  if errorCode != nil {
   result.errorCode = errorCode.(string)
  }
  if recommend != nil {
   result.recommend = recommend.(string)
  }
  if message != nil {
   result.message = message.(string)
  }
 }
 return result
}
func WrapServerError(originError *ServerError, wrapInfo map[string]string) *ServerError {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, wrapper := range wrapperList {
  ok, newError := wrapper.tryWrap(originError, wrapInfo)
  if ok {
   return newError
  }
 }
 return originError
}
func (err *ServerError) HttpStatus() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.httpStatus
}
func (err *ServerError) ErrorCode() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.errorCode
}
func (err *ServerError) Message() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.message
}
func (err *ServerError) OriginError() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil
}
func (err *ServerError) HostId() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.hostId
}
func (err *ServerError) RequestId() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.requestId
}
func (err *ServerError) Recommend() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.recommend
}
func (err *ServerError) Comment() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return err.comment
}
