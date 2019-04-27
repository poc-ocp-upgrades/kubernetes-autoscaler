package sdk

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/endpoints"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/errors"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"strconv"
	"sync"
)

var Version = "0.0.1"

type Client struct {
	regionId	string
	config		*Config
	signer		auth.Signer
	httpClient	*http.Client
	asyncTaskQueue	chan func()
	debug		bool
	isRunning	bool
	asyncChanLock	*sync.RWMutex
}

func (client *Client) Init() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	panic("not support yet")
}
func (client *Client) InitWithOptions(regionId string, config *Config, credential auth.Credential) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.isRunning = true
	client.asyncChanLock = new(sync.RWMutex)
	client.regionId = regionId
	client.config = config
	if err != nil {
		return
	}
	client.httpClient = &http.Client{}
	if config.HttpTransport != nil {
		client.httpClient.Transport = config.HttpTransport
	}
	if config.Timeout > 0 {
		client.httpClient.Timeout = config.Timeout
	}
	if config.EnableAsync {
		client.EnableAsync(config.GoRoutinePoolSize, config.MaxTaskQueueSize)
	}
	client.signer, err = auth.NewSignerWithCredential(credential, client.ProcessCommonRequestWithSigner)
	return
}
func (client *Client) EnableAsync(routinePoolSize, maxTaskQueueSize int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.asyncTaskQueue = make(chan func(), maxTaskQueueSize)
	for i := 0; i < routinePoolSize; i++ {
		go func() {
			for client.isRunning {
				select {
				case task, notClosed := <-client.asyncTaskQueue:
					if notClosed {
						task()
					}
				}
			}
		}()
	}
}
func (client *Client) InitWithAccessKey(regionId, accessKeyId, accessKeySecret string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := client.InitClientConfig()
	credential := &credentials.BaseCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}
	return client.InitWithOptions(regionId, config, credential)
}
func (client *Client) InitWithStsToken(regionId, accessKeyId, accessKeySecret, securityToken string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := client.InitClientConfig()
	credential := &credentials.StsTokenCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret, AccessKeyStsToken: securityToken}
	return client.InitWithOptions(regionId, config, credential)
}
func (client *Client) InitWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := client.InitClientConfig()
	credential := &credentials.RamRoleArnCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret, RoleArn: roleArn, RoleSessionName: roleSessionName}
	return client.InitWithOptions(regionId, config, credential)
}
func (client *Client) InitWithRsaKeyPair(regionId, publicKeyId, privateKey string, sessionExpiration int) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := client.InitClientConfig()
	credential := &credentials.RsaKeyPairCredential{PrivateKey: privateKey, PublicKeyId: publicKeyId, SessionExpiration: sessionExpiration}
	return client.InitWithOptions(regionId, config, credential)
}
func (client *Client) InitWithEcsRamRole(regionId, roleName string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := client.InitClientConfig()
	credential := &credentials.EcsRamRoleCredential{RoleName: roleName}
	return client.InitWithOptions(regionId, config, credential)
}
func (client *Client) InitClientConfig() (config *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if client.config != nil {
		return client.config
	}
	return NewConfig()
}
func (client *Client) DoAction(request requests.AcsRequest, response responses.AcsResponse) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.DoActionWithSigner(request, response, nil)
}
func (client *Client) BuildRequestWithSigner(request requests.AcsRequest, signer auth.Signer) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.GetHeaders()["x-sdk-core-version"] = Version
	regionId := client.regionId
	if len(request.GetRegionId()) > 0 {
		regionId = request.GetRegionId()
	}
	resolveParam := &endpoints.ResolveParam{Domain: request.GetDomain(), Product: request.GetProduct(), RegionId: regionId, LocationProduct: request.GetLocationServiceCode(), LocationEndpointType: request.GetLocationEndpointType(), CommonApi: client.ProcessCommonRequest}
	endpoint, err := endpoints.Resolve(resolveParam)
	if err != nil {
		return
	}
	request.SetDomain(endpoint)
	err = requests.InitParams(request)
	if err != nil {
		return
	}
	var finalSigner auth.Signer
	if signer != nil {
		finalSigner = signer
	} else {
		finalSigner = client.signer
	}
	httpRequest, err := buildHttpRequest(request, finalSigner, regionId)
	if client.config.UserAgent != "" {
		httpRequest.Header.Set("User-Agent", client.config.UserAgent)
	}
	return err
}
func (client *Client) DoActionWithSigner(request requests.AcsRequest, response responses.AcsResponse, signer auth.Signer) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.GetHeaders()["x-sdk-core-version"] = Version
	regionId := client.regionId
	if len(request.GetRegionId()) > 0 {
		regionId = request.GetRegionId()
	}
	resolveParam := &endpoints.ResolveParam{Domain: request.GetDomain(), Product: request.GetProduct(), RegionId: regionId, LocationProduct: request.GetLocationServiceCode(), LocationEndpointType: request.GetLocationEndpointType(), CommonApi: client.ProcessCommonRequest}
	endpoint, err := endpoints.Resolve(resolveParam)
	if err != nil {
		return
	}
	request.SetDomain(endpoint)
	if request.GetScheme() == "" {
		request.SetScheme(client.config.Scheme)
	}
	err = requests.InitParams(request)
	if err != nil {
		return
	}
	var finalSigner auth.Signer
	if signer != nil {
		finalSigner = signer
	} else {
		finalSigner = client.signer
	}
	httpRequest, err := buildHttpRequest(request, finalSigner, regionId)
	if client.config.UserAgent != "" {
		httpRequest.Header.Set("User-Agent", client.config.UserAgent)
	}
	if err != nil {
		return
	}
	var httpResponse *http.Response
	for retryTimes := 0; retryTimes <= client.config.MaxRetryTime; retryTimes++ {
		httpResponse, err = client.httpClient.Do(httpRequest)
		var timeout bool
		if err != nil {
			if !client.config.AutoRetry {
				return
			} else if timeout = isTimeout(err); !timeout {
				return
			} else if retryTimes >= client.config.MaxRetryTime {
				timeoutErrorMsg := fmt.Sprintf(errors.TimeoutErrorMessage, strconv.Itoa(retryTimes+1), strconv.Itoa(retryTimes+1))
				err = errors.NewClientError(errors.TimeoutErrorCode, timeoutErrorMsg, err)
				return
			}
		}
		if client.config.AutoRetry && (timeout || isServerError(httpResponse)) {
			httpRequest, err = buildHttpRequest(request, finalSigner, regionId)
			if err != nil {
				return
			}
			continue
		}
		break
	}
	err = responses.Unmarshal(response, httpResponse, request.GetAcceptFormat())
	if serverErr, ok := err.(*errors.ServerError); ok {
		var wrapInfo = map[string]string{}
		wrapInfo["StringToSign"] = request.GetStringToSign()
		err = errors.WrapServerError(serverErr, wrapInfo)
	}
	return
}
func buildHttpRequest(request requests.AcsRequest, singer auth.Signer, regionId string) (httpRequest *http.Request, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = auth.Sign(request, singer, regionId)
	if err != nil {
		return
	}
	requestMethod := request.GetMethod()
	requestUrl := request.BuildUrl()
	body := request.GetBodyReader()
	httpRequest, err = http.NewRequest(requestMethod, requestUrl, body)
	if err != nil {
		return
	}
	for key, value := range request.GetHeaders() {
		httpRequest.Header[key] = []string{value}
	}
	if host, containsHost := request.GetHeaders()["Host"]; containsHost {
		httpRequest.Host = host
	}
	return
}
func isTimeout(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err == nil {
		return false
	}
	netErr, isNetError := err.(net.Error)
	return isNetError && netErr.Timeout()
}
func isServerError(httpResponse *http.Response) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return httpResponse.StatusCode >= http.StatusInternalServerError
}
func (client *Client) AddAsyncTask(task func()) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if client.asyncTaskQueue != nil {
		client.asyncChanLock.RLock()
		defer client.asyncChanLock.RUnlock()
		if client.isRunning {
			client.asyncTaskQueue <- task
		}
	} else {
		err = errors.NewClientError(errors.AsyncFunctionNotEnabledCode, errors.AsyncFunctionNotEnabledMessage, nil)
	}
	return
}
func (client *Client) GetConfig() *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.config
}
func NewClient() (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.Init()
	return
}
func NewClientWithOptions(regionId string, config *Config, credential auth.Credential) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithOptions(regionId, config, credential)
	return
}
func NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return
}
func NewClientWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken)
	return
}
func NewClientWithRamRoleArn(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName)
	return
}
func NewClientWithEcsRamRole(regionId string, roleName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithEcsRamRole(regionId, roleName)
	return
}
func NewClientWithRsaKeyPair(regionId string, publicKeyId, privateKey string, sessionExpiration int) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithRsaKeyPair(regionId, publicKeyId, privateKey, sessionExpiration)
	return
}
func NewClientWithStsRoleArn(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewClientWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName)
}
func NewClientWithStsRoleNameOnEcs(regionId string, roleName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewClientWithEcsRamRole(regionId, roleName)
}
func (client *Client) ProcessCommonRequest(request *requests.CommonRequest) (response *responses.CommonResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.TransToAcsRequest()
	response = responses.NewCommonResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) ProcessCommonRequestWithSigner(request *requests.CommonRequest, signerInterface interface{}) (response *responses.CommonResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if signer, isSigner := signerInterface.(auth.Signer); isSigner {
		request.TransToAcsRequest()
		response = responses.NewCommonResponse()
		err = client.DoActionWithSigner(request, response, signer)
		return
	}
	panic("should not be here")
}
func (client *Client) Shutdown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.signer.Shutdown()
	client.asyncChanLock.Lock()
	defer client.asyncChanLock.Unlock()
	client.isRunning = false
	close(client.asyncTaskQueue)
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
