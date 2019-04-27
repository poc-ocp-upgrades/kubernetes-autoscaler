package auth

import (
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/signers"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/errors"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
	"reflect"
)

type Signer interface {
	GetName() string
	GetType() string
	GetVersion() string
	GetAccessKeyId() (string, error)
	GetExtraParam() map[string]string
	Sign(stringToSign, secretSuffix string) string
	Shutdown()
}

func NewSignerWithCredential(credential Credential, commonApi func(request *requests.CommonRequest, signer interface{}) (response *responses.CommonResponse, err error)) (signer Signer, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch instance := credential.(type) {
	case *credentials.AccessKeyCredential:
		{
			signer, err = signers.NewAccessKeySigner(instance)
		}
	case *credentials.StsTokenCredential:
		{
			signer, err = signers.NewStsTokenSigner(instance)
		}
	case *credentials.RamRoleArnCredential:
		{
			signer, err = signers.NewRamRoleArnSigner(instance, commonApi)
		}
	case *credentials.RsaKeyPairCredential:
		{
			signer, err = signers.NewSignerKeyPair(instance, commonApi)
		}
	case *credentials.EcsRamRoleCredential:
		{
			signer, err = signers.NewEcsRamRoleSigner(instance, commonApi)
		}
	case *credentials.BaseCredential:
		{
			signer, err = signers.NewAccessKeySigner(instance.ToAccessKeyCredential())
		}
	case *credentials.StsRoleArnCredential:
		{
			signer, err = signers.NewRamRoleArnSigner(instance.ToRamRoleArnCredential(), commonApi)
		}
	case *credentials.StsRoleNameOnEcsCredential:
		{
			signer, err = signers.NewEcsRamRoleSigner(instance.ToEcsRamRoleCredential(), commonApi)
		}
	default:
		message := fmt.Sprintf(errors.UnsupportedCredentialErrorMessage, reflect.TypeOf(credential))
		err = errors.NewClientError(errors.UnsupportedCredentialErrorCode, message, nil)
	}
	return
}
func Sign(request requests.AcsRequest, signer Signer, regionId string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch request.GetStyle() {
	case requests.ROA:
		{
			signRoaRequest(request, signer, regionId)
		}
	case requests.RPC:
		{
			err = signRpcRequest(request, signer, regionId)
		}
	default:
		message := fmt.Sprintf(errors.UnknownRequestTypeErrorMessage, reflect.TypeOf(request))
		err = errors.NewClientError(errors.UnknownRequestTypeErrorCode, message, nil)
	}
	return
}
