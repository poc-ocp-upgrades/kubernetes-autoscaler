package errors

import "strings"

const (
 SignatureDostNotMatchErrorCode = "SignatureDoesNotMatch"
 MessagePrefix                  = "Specified signature is not matched with our calculation. server string to sign is:"
)

type SignatureDostNotMatchWrapper struct{}

func (*SignatureDostNotMatchWrapper) tryWrap(error *ServerError, wrapInfo map[string]string) (bool, *ServerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 clientStringToSign := wrapInfo["StringToSign"]
 if error.errorCode == SignatureDostNotMatchErrorCode && clientStringToSign != "" {
  message := error.message
  if strings.HasPrefix(message, MessagePrefix) {
   serverStringToSign := message[len(MessagePrefix):]
   if clientStringToSign == serverStringToSign {
    error.recommend = "Please check you AccessKeySecret"
   } else {
    error.recommend = "This may be a bug with the SDK and we hope you can submit this question in the " + "github issue(https://github.com/aliyun/alibaba-cloud-sdk-go/issues), thanks very much"
   }
  }
  return true, error
 }
 return false, nil
}
