package signers

import (
 "encoding/json"
 "fmt"
 "github.com/jmespath/go-jmespath"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/errors"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
 "net/http"
 "strconv"
)

type SignerKeyPair struct {
 *credentialUpdater
 sessionCredential *SessionCredential
 credential        *credentials.RsaKeyPairCredential
 commonApi         func(request *requests.CommonRequest, signer interface{}) (response *responses.CommonResponse, err error)
}

func NewSignerKeyPair(credential *credentials.RsaKeyPairCredential, commonApi func(*requests.CommonRequest, interface{}) (response *responses.CommonResponse, err error)) (signer *SignerKeyPair, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 signer = &SignerKeyPair{credential: credential, commonApi: commonApi}
 signer.credentialUpdater = &credentialUpdater{credentialExpiration: credential.SessionExpiration, buildRequestMethod: signer.buildCommonRequest, responseCallBack: signer.refreshCredential, refreshApi: signer.refreshApi}
 if credential.SessionExpiration > 0 {
  if credential.SessionExpiration >= 900 && credential.SessionExpiration <= 3600 {
   signer.credentialExpiration = credential.SessionExpiration
  } else {
   err = errors.NewClientError(errors.InvalidParamErrorCode, "Key Pair session duration should be in the range of 15min - 1Hr", nil)
  }
 } else {
  signer.credentialExpiration = defaultDurationSeconds
 }
 return
}
func (*SignerKeyPair) GetName() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return "HMAC-SHA1"
}
func (*SignerKeyPair) GetType() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ""
}
func (*SignerKeyPair) GetVersion() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return "1.0"
}
func (signer *SignerKeyPair) GetAccessKeyId() (accessKeyId string, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if signer.sessionCredential == nil || signer.needUpdateCredential() {
  err = signer.updateCredential()
 }
 if err != nil && (signer.sessionCredential == nil || len(signer.sessionCredential.AccessKeyId) <= 0) {
  return "", err
 }
 return signer.sessionCredential.AccessKeyId, err
}
func (signer *SignerKeyPair) GetExtraParam() map[string]string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if signer.sessionCredential == nil || signer.needUpdateCredential() {
  signer.updateCredential()
 }
 if signer.sessionCredential == nil || len(signer.sessionCredential.AccessKeyId) <= 0 {
  return make(map[string]string)
 }
 return make(map[string]string)
}
func (signer *SignerKeyPair) Sign(stringToSign, secretSuffix string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 secret := signer.sessionCredential.AccessKeyId + secretSuffix
 return ShaHmac1(stringToSign, secret)
}
func (signer *SignerKeyPair) buildCommonRequest() (request *requests.CommonRequest, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 request = requests.NewCommonRequest()
 request.Product = "Sts"
 request.Version = "2015-04-01"
 request.ApiName = "GenerateSessionAccessKey"
 request.Scheme = requests.HTTPS
 request.QueryParams["PublicKeyId"] = signer.credential.PublicKeyId
 request.QueryParams["DurationSeconds"] = strconv.Itoa(signer.credentialExpiration)
 return
}
func (signer *SignerKeyPair) refreshApi(request *requests.CommonRequest) (response *responses.CommonResponse, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 signerV2, err := NewSignerV2(signer.credential)
 return signer.commonApi(request, signerV2)
}
func (signer *SignerKeyPair) refreshCredential(response *responses.CommonResponse) (err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if response.GetHttpStatus() != http.StatusOK {
  message := "refresh session AccessKey failed"
  err = errors.NewServerError(response.GetHttpStatus(), response.GetHttpContentString(), message)
  return
 }
 var data interface{}
 err = json.Unmarshal(response.GetHttpContentBytes(), &data)
 if err != nil {
  fmt.Println("refresh KeyPair err, json.Unmarshal fail", err)
  return
 }
 accessKeyId, err := jmespath.Search("SessionAccessKey.SessionAccessKeyId", data)
 if err != nil {
  fmt.Println("refresh KeyPair err, fail to get SessionAccessKeyId", err)
  return
 }
 accessKeySecret, err := jmespath.Search("SessionAccessKey.SessionAccessKeySecret", data)
 if err != nil {
  fmt.Println("refresh KeyPair err, fail to get SessionAccessKeySecret", err)
  return
 }
 if accessKeyId == nil || accessKeySecret == nil {
  return
 }
 signer.sessionCredential = &SessionCredential{AccessKeyId: accessKeyId.(string), AccessKeySecret: accessKeySecret.(string)}
 return
}
func (signer *SignerKeyPair) Shutdown() {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
