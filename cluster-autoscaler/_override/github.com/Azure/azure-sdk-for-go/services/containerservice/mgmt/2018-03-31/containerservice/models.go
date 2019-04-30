package containerservice

import (
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

type NetworkPlugin string

const (
	Azure	NetworkPlugin	= "azure"
	Kubenet	NetworkPlugin	= "kubenet"
)

func PossibleNetworkPluginValues() []NetworkPlugin {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []NetworkPlugin{Azure, Kubenet}
}

type NetworkPolicy string

const (
	Calico NetworkPolicy = "calico"
)

func PossibleNetworkPolicyValues() []NetworkPolicy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []NetworkPolicy{Calico}
}

type OrchestratorTypes string

const (
	Custom		OrchestratorTypes	= "Custom"
	DCOS		OrchestratorTypes	= "DCOS"
	DockerCE	OrchestratorTypes	= "DockerCE"
	Kubernetes	OrchestratorTypes	= "Kubernetes"
	Swarm		OrchestratorTypes	= "Swarm"
)

func PossibleOrchestratorTypesValues() []OrchestratorTypes {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []OrchestratorTypes{Custom, DCOS, DockerCE, Kubernetes, Swarm}
}

type OSType string

const (
	Linux	OSType	= "Linux"
	Windows	OSType	= "Windows"
)

func PossibleOSTypeValues() []OSType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []OSType{Linux, Windows}
}

type StorageProfileTypes string

const (
	ManagedDisks	StorageProfileTypes	= "ManagedDisks"
	StorageAccount	StorageProfileTypes	= "StorageAccount"
)

func PossibleStorageProfileTypesValues() []StorageProfileTypes {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []StorageProfileTypes{ManagedDisks, StorageAccount}
}

type VMSizeTypes string

const (
	StandardA1		VMSizeTypes	= "Standard_A1"
	StandardA10		VMSizeTypes	= "Standard_A10"
	StandardA11		VMSizeTypes	= "Standard_A11"
	StandardA1V2		VMSizeTypes	= "Standard_A1_v2"
	StandardA2		VMSizeTypes	= "Standard_A2"
	StandardA2mV2		VMSizeTypes	= "Standard_A2m_v2"
	StandardA2V2		VMSizeTypes	= "Standard_A2_v2"
	StandardA3		VMSizeTypes	= "Standard_A3"
	StandardA4		VMSizeTypes	= "Standard_A4"
	StandardA4mV2		VMSizeTypes	= "Standard_A4m_v2"
	StandardA4V2		VMSizeTypes	= "Standard_A4_v2"
	StandardA5		VMSizeTypes	= "Standard_A5"
	StandardA6		VMSizeTypes	= "Standard_A6"
	StandardA7		VMSizeTypes	= "Standard_A7"
	StandardA8		VMSizeTypes	= "Standard_A8"
	StandardA8mV2		VMSizeTypes	= "Standard_A8m_v2"
	StandardA8V2		VMSizeTypes	= "Standard_A8_v2"
	StandardA9		VMSizeTypes	= "Standard_A9"
	StandardB2ms		VMSizeTypes	= "Standard_B2ms"
	StandardB2s		VMSizeTypes	= "Standard_B2s"
	StandardB4ms		VMSizeTypes	= "Standard_B4ms"
	StandardB8ms		VMSizeTypes	= "Standard_B8ms"
	StandardD1		VMSizeTypes	= "Standard_D1"
	StandardD11		VMSizeTypes	= "Standard_D11"
	StandardD11V2		VMSizeTypes	= "Standard_D11_v2"
	StandardD11V2Promo	VMSizeTypes	= "Standard_D11_v2_Promo"
	StandardD12		VMSizeTypes	= "Standard_D12"
	StandardD12V2		VMSizeTypes	= "Standard_D12_v2"
	StandardD12V2Promo	VMSizeTypes	= "Standard_D12_v2_Promo"
	StandardD13		VMSizeTypes	= "Standard_D13"
	StandardD13V2		VMSizeTypes	= "Standard_D13_v2"
	StandardD13V2Promo	VMSizeTypes	= "Standard_D13_v2_Promo"
	StandardD14		VMSizeTypes	= "Standard_D14"
	StandardD14V2		VMSizeTypes	= "Standard_D14_v2"
	StandardD14V2Promo	VMSizeTypes	= "Standard_D14_v2_Promo"
	StandardD15V2		VMSizeTypes	= "Standard_D15_v2"
	StandardD16sV3		VMSizeTypes	= "Standard_D16s_v3"
	StandardD16V3		VMSizeTypes	= "Standard_D16_v3"
	StandardD1V2		VMSizeTypes	= "Standard_D1_v2"
	StandardD2		VMSizeTypes	= "Standard_D2"
	StandardD2sV3		VMSizeTypes	= "Standard_D2s_v3"
	StandardD2V2		VMSizeTypes	= "Standard_D2_v2"
	StandardD2V2Promo	VMSizeTypes	= "Standard_D2_v2_Promo"
	StandardD2V3		VMSizeTypes	= "Standard_D2_v3"
	StandardD3		VMSizeTypes	= "Standard_D3"
	StandardD32sV3		VMSizeTypes	= "Standard_D32s_v3"
	StandardD32V3		VMSizeTypes	= "Standard_D32_v3"
	StandardD3V2		VMSizeTypes	= "Standard_D3_v2"
	StandardD3V2Promo	VMSizeTypes	= "Standard_D3_v2_Promo"
	StandardD4		VMSizeTypes	= "Standard_D4"
	StandardD4sV3		VMSizeTypes	= "Standard_D4s_v3"
	StandardD4V2		VMSizeTypes	= "Standard_D4_v2"
	StandardD4V2Promo	VMSizeTypes	= "Standard_D4_v2_Promo"
	StandardD4V3		VMSizeTypes	= "Standard_D4_v3"
	StandardD5V2		VMSizeTypes	= "Standard_D5_v2"
	StandardD5V2Promo	VMSizeTypes	= "Standard_D5_v2_Promo"
	StandardD64sV3		VMSizeTypes	= "Standard_D64s_v3"
	StandardD64V3		VMSizeTypes	= "Standard_D64_v3"
	StandardD8sV3		VMSizeTypes	= "Standard_D8s_v3"
	StandardD8V3		VMSizeTypes	= "Standard_D8_v3"
	StandardDS1		VMSizeTypes	= "Standard_DS1"
	StandardDS11		VMSizeTypes	= "Standard_DS11"
	StandardDS11V2		VMSizeTypes	= "Standard_DS11_v2"
	StandardDS11V2Promo	VMSizeTypes	= "Standard_DS11_v2_Promo"
	StandardDS12		VMSizeTypes	= "Standard_DS12"
	StandardDS12V2		VMSizeTypes	= "Standard_DS12_v2"
	StandardDS12V2Promo	VMSizeTypes	= "Standard_DS12_v2_Promo"
	StandardDS13		VMSizeTypes	= "Standard_DS13"
	StandardDS132V2		VMSizeTypes	= "Standard_DS13-2_v2"
	StandardDS134V2		VMSizeTypes	= "Standard_DS13-4_v2"
	StandardDS13V2		VMSizeTypes	= "Standard_DS13_v2"
	StandardDS13V2Promo	VMSizeTypes	= "Standard_DS13_v2_Promo"
	StandardDS14		VMSizeTypes	= "Standard_DS14"
	StandardDS144V2		VMSizeTypes	= "Standard_DS14-4_v2"
	StandardDS148V2		VMSizeTypes	= "Standard_DS14-8_v2"
	StandardDS14V2		VMSizeTypes	= "Standard_DS14_v2"
	StandardDS14V2Promo	VMSizeTypes	= "Standard_DS14_v2_Promo"
	StandardDS15V2		VMSizeTypes	= "Standard_DS15_v2"
	StandardDS1V2		VMSizeTypes	= "Standard_DS1_v2"
	StandardDS2		VMSizeTypes	= "Standard_DS2"
	StandardDS2V2		VMSizeTypes	= "Standard_DS2_v2"
	StandardDS2V2Promo	VMSizeTypes	= "Standard_DS2_v2_Promo"
	StandardDS3		VMSizeTypes	= "Standard_DS3"
	StandardDS3V2		VMSizeTypes	= "Standard_DS3_v2"
	StandardDS3V2Promo	VMSizeTypes	= "Standard_DS3_v2_Promo"
	StandardDS4		VMSizeTypes	= "Standard_DS4"
	StandardDS4V2		VMSizeTypes	= "Standard_DS4_v2"
	StandardDS4V2Promo	VMSizeTypes	= "Standard_DS4_v2_Promo"
	StandardDS5V2		VMSizeTypes	= "Standard_DS5_v2"
	StandardDS5V2Promo	VMSizeTypes	= "Standard_DS5_v2_Promo"
	StandardE16sV3		VMSizeTypes	= "Standard_E16s_v3"
	StandardE16V3		VMSizeTypes	= "Standard_E16_v3"
	StandardE2sV3		VMSizeTypes	= "Standard_E2s_v3"
	StandardE2V3		VMSizeTypes	= "Standard_E2_v3"
	StandardE3216sV3	VMSizeTypes	= "Standard_E32-16s_v3"
	StandardE328sV3		VMSizeTypes	= "Standard_E32-8s_v3"
	StandardE32sV3		VMSizeTypes	= "Standard_E32s_v3"
	StandardE32V3		VMSizeTypes	= "Standard_E32_v3"
	StandardE4sV3		VMSizeTypes	= "Standard_E4s_v3"
	StandardE4V3		VMSizeTypes	= "Standard_E4_v3"
	StandardE6416sV3	VMSizeTypes	= "Standard_E64-16s_v3"
	StandardE6432sV3	VMSizeTypes	= "Standard_E64-32s_v3"
	StandardE64sV3		VMSizeTypes	= "Standard_E64s_v3"
	StandardE64V3		VMSizeTypes	= "Standard_E64_v3"
	StandardE8sV3		VMSizeTypes	= "Standard_E8s_v3"
	StandardE8V3		VMSizeTypes	= "Standard_E8_v3"
	StandardF1		VMSizeTypes	= "Standard_F1"
	StandardF16		VMSizeTypes	= "Standard_F16"
	StandardF16s		VMSizeTypes	= "Standard_F16s"
	StandardF16sV2		VMSizeTypes	= "Standard_F16s_v2"
	StandardF1s		VMSizeTypes	= "Standard_F1s"
	StandardF2		VMSizeTypes	= "Standard_F2"
	StandardF2s		VMSizeTypes	= "Standard_F2s"
	StandardF2sV2		VMSizeTypes	= "Standard_F2s_v2"
	StandardF32sV2		VMSizeTypes	= "Standard_F32s_v2"
	StandardF4		VMSizeTypes	= "Standard_F4"
	StandardF4s		VMSizeTypes	= "Standard_F4s"
	StandardF4sV2		VMSizeTypes	= "Standard_F4s_v2"
	StandardF64sV2		VMSizeTypes	= "Standard_F64s_v2"
	StandardF72sV2		VMSizeTypes	= "Standard_F72s_v2"
	StandardF8		VMSizeTypes	= "Standard_F8"
	StandardF8s		VMSizeTypes	= "Standard_F8s"
	StandardF8sV2		VMSizeTypes	= "Standard_F8s_v2"
	StandardG1		VMSizeTypes	= "Standard_G1"
	StandardG2		VMSizeTypes	= "Standard_G2"
	StandardG3		VMSizeTypes	= "Standard_G3"
	StandardG4		VMSizeTypes	= "Standard_G4"
	StandardG5		VMSizeTypes	= "Standard_G5"
	StandardGS1		VMSizeTypes	= "Standard_GS1"
	StandardGS2		VMSizeTypes	= "Standard_GS2"
	StandardGS3		VMSizeTypes	= "Standard_GS3"
	StandardGS4		VMSizeTypes	= "Standard_GS4"
	StandardGS44		VMSizeTypes	= "Standard_GS4-4"
	StandardGS48		VMSizeTypes	= "Standard_GS4-8"
	StandardGS5		VMSizeTypes	= "Standard_GS5"
	StandardGS516		VMSizeTypes	= "Standard_GS5-16"
	StandardGS58		VMSizeTypes	= "Standard_GS5-8"
	StandardH16		VMSizeTypes	= "Standard_H16"
	StandardH16m		VMSizeTypes	= "Standard_H16m"
	StandardH16mr		VMSizeTypes	= "Standard_H16mr"
	StandardH16r		VMSizeTypes	= "Standard_H16r"
	StandardH8		VMSizeTypes	= "Standard_H8"
	StandardH8m		VMSizeTypes	= "Standard_H8m"
	StandardL16s		VMSizeTypes	= "Standard_L16s"
	StandardL32s		VMSizeTypes	= "Standard_L32s"
	StandardL4s		VMSizeTypes	= "Standard_L4s"
	StandardL8s		VMSizeTypes	= "Standard_L8s"
	StandardM12832ms	VMSizeTypes	= "Standard_M128-32ms"
	StandardM12864ms	VMSizeTypes	= "Standard_M128-64ms"
	StandardM128ms		VMSizeTypes	= "Standard_M128ms"
	StandardM128s		VMSizeTypes	= "Standard_M128s"
	StandardM6416ms		VMSizeTypes	= "Standard_M64-16ms"
	StandardM6432ms		VMSizeTypes	= "Standard_M64-32ms"
	StandardM64ms		VMSizeTypes	= "Standard_M64ms"
	StandardM64s		VMSizeTypes	= "Standard_M64s"
	StandardNC12		VMSizeTypes	= "Standard_NC12"
	StandardNC12sV2		VMSizeTypes	= "Standard_NC12s_v2"
	StandardNC12sV3		VMSizeTypes	= "Standard_NC12s_v3"
	StandardNC24		VMSizeTypes	= "Standard_NC24"
	StandardNC24r		VMSizeTypes	= "Standard_NC24r"
	StandardNC24rsV2	VMSizeTypes	= "Standard_NC24rs_v2"
	StandardNC24rsV3	VMSizeTypes	= "Standard_NC24rs_v3"
	StandardNC24sV2		VMSizeTypes	= "Standard_NC24s_v2"
	StandardNC24sV3		VMSizeTypes	= "Standard_NC24s_v3"
	StandardNC6		VMSizeTypes	= "Standard_NC6"
	StandardNC6sV2		VMSizeTypes	= "Standard_NC6s_v2"
	StandardNC6sV3		VMSizeTypes	= "Standard_NC6s_v3"
	StandardND12s		VMSizeTypes	= "Standard_ND12s"
	StandardND24rs		VMSizeTypes	= "Standard_ND24rs"
	StandardND24s		VMSizeTypes	= "Standard_ND24s"
	StandardND6s		VMSizeTypes	= "Standard_ND6s"
	StandardNV12		VMSizeTypes	= "Standard_NV12"
	StandardNV24		VMSizeTypes	= "Standard_NV24"
	StandardNV6		VMSizeTypes	= "Standard_NV6"
)

func PossibleVMSizeTypesValues() []VMSizeTypes {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []VMSizeTypes{StandardA1, StandardA10, StandardA11, StandardA1V2, StandardA2, StandardA2mV2, StandardA2V2, StandardA3, StandardA4, StandardA4mV2, StandardA4V2, StandardA5, StandardA6, StandardA7, StandardA8, StandardA8mV2, StandardA8V2, StandardA9, StandardB2ms, StandardB2s, StandardB4ms, StandardB8ms, StandardD1, StandardD11, StandardD11V2, StandardD11V2Promo, StandardD12, StandardD12V2, StandardD12V2Promo, StandardD13, StandardD13V2, StandardD13V2Promo, StandardD14, StandardD14V2, StandardD14V2Promo, StandardD15V2, StandardD16sV3, StandardD16V3, StandardD1V2, StandardD2, StandardD2sV3, StandardD2V2, StandardD2V2Promo, StandardD2V3, StandardD3, StandardD32sV3, StandardD32V3, StandardD3V2, StandardD3V2Promo, StandardD4, StandardD4sV3, StandardD4V2, StandardD4V2Promo, StandardD4V3, StandardD5V2, StandardD5V2Promo, StandardD64sV3, StandardD64V3, StandardD8sV3, StandardD8V3, StandardDS1, StandardDS11, StandardDS11V2, StandardDS11V2Promo, StandardDS12, StandardDS12V2, StandardDS12V2Promo, StandardDS13, StandardDS132V2, StandardDS134V2, StandardDS13V2, StandardDS13V2Promo, StandardDS14, StandardDS144V2, StandardDS148V2, StandardDS14V2, StandardDS14V2Promo, StandardDS15V2, StandardDS1V2, StandardDS2, StandardDS2V2, StandardDS2V2Promo, StandardDS3, StandardDS3V2, StandardDS3V2Promo, StandardDS4, StandardDS4V2, StandardDS4V2Promo, StandardDS5V2, StandardDS5V2Promo, StandardE16sV3, StandardE16V3, StandardE2sV3, StandardE2V3, StandardE3216sV3, StandardE328sV3, StandardE32sV3, StandardE32V3, StandardE4sV3, StandardE4V3, StandardE6416sV3, StandardE6432sV3, StandardE64sV3, StandardE64V3, StandardE8sV3, StandardE8V3, StandardF1, StandardF16, StandardF16s, StandardF16sV2, StandardF1s, StandardF2, StandardF2s, StandardF2sV2, StandardF32sV2, StandardF4, StandardF4s, StandardF4sV2, StandardF64sV2, StandardF72sV2, StandardF8, StandardF8s, StandardF8sV2, StandardG1, StandardG2, StandardG3, StandardG4, StandardG5, StandardGS1, StandardGS2, StandardGS3, StandardGS4, StandardGS44, StandardGS48, StandardGS5, StandardGS516, StandardGS58, StandardH16, StandardH16m, StandardH16mr, StandardH16r, StandardH8, StandardH8m, StandardL16s, StandardL32s, StandardL4s, StandardL8s, StandardM12832ms, StandardM12864ms, StandardM128ms, StandardM128s, StandardM6416ms, StandardM6432ms, StandardM64ms, StandardM64s, StandardNC12, StandardNC12sV2, StandardNC12sV3, StandardNC24, StandardNC24r, StandardNC24rsV2, StandardNC24rsV3, StandardNC24sV2, StandardNC24sV3, StandardNC6, StandardNC6sV2, StandardNC6sV3, StandardND12s, StandardND24rs, StandardND24s, StandardND6s, StandardNV12, StandardNV24, StandardNV6}
}

type AccessProfile struct {
	KubeConfig *[]byte `json:"kubeConfig,omitempty"`
}
type AgentPoolProfile struct {
	Name		*string			`json:"name,omitempty"`
	Count		*int32			`json:"count,omitempty"`
	VMSize		VMSizeTypes		`json:"vmSize,omitempty"`
	OsDiskSizeGB	*int32			`json:"osDiskSizeGB,omitempty"`
	DNSPrefix	*string			`json:"dnsPrefix,omitempty"`
	Fqdn		*string			`json:"fqdn,omitempty"`
	Ports		*[]int32		`json:"ports,omitempty"`
	StorageProfile	StorageProfileTypes	`json:"storageProfile,omitempty"`
	VnetSubnetID	*string			`json:"vnetSubnetID,omitempty"`
	OsType		OSType			`json:"osType,omitempty"`
}
type ContainerService struct {
	autorest.Response	`json:"-"`
	*Properties		`json:"properties,omitempty"`
	ID			*string			`json:"id,omitempty"`
	Name			*string			`json:"name,omitempty"`
	Type			*string			`json:"type,omitempty"`
	Location		*string			`json:"location,omitempty"`
	Tags			map[string]*string	`json:"tags"`
}

func (cs ContainerService) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if cs.Properties != nil {
		objectMap["properties"] = cs.Properties
	}
	if cs.ID != nil {
		objectMap["id"] = cs.ID
	}
	if cs.Name != nil {
		objectMap["name"] = cs.Name
	}
	if cs.Type != nil {
		objectMap["type"] = cs.Type
	}
	if cs.Location != nil {
		objectMap["location"] = cs.Location
	}
	if cs.Tags != nil {
		objectMap["tags"] = cs.Tags
	}
	return json.Marshal(objectMap)
}
func (cs *ContainerService) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var properties Properties
				err = json.Unmarshal(*v, &properties)
				if err != nil {
					return err
				}
				cs.Properties = &properties
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				cs.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				cs.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				cs.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				cs.Location = &location
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				cs.Tags = tags
			}
		}
	}
	return nil
}

type ContainerServicesCreateOrUpdateFutureType struct{ azure.Future }

func (future *ContainerServicesCreateOrUpdateFutureType) Result(client ContainerServicesClient) (cs ContainerService, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesCreateOrUpdateFutureType", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("containerservice.ContainerServicesCreateOrUpdateFutureType")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if cs.Response.Response, err = future.GetResult(sender); err == nil && cs.Response.Response.StatusCode != http.StatusNoContent {
		cs, err = client.CreateOrUpdateResponder(cs.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesCreateOrUpdateFutureType", "Result", cs.Response.Response, "Failure responding to request")
		}
	}
	return
}

type ContainerServicesDeleteFutureType struct{ azure.Future }

func (future *ContainerServicesDeleteFutureType) Result(client ContainerServicesClient) (ar autorest.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesDeleteFutureType", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("containerservice.ContainerServicesDeleteFutureType")
		return
	}
	ar.Response = future.Response()
	return
}

type CustomProfile struct {
	Orchestrator *string `json:"orchestrator,omitempty"`
}
type DiagnosticsProfile struct {
	VMDiagnostics *VMDiagnostics `json:"vmDiagnostics,omitempty"`
}
type KeyVaultSecretRef struct {
	VaultID		*string	`json:"vaultID,omitempty"`
	SecretName	*string	`json:"secretName,omitempty"`
	Version		*string	`json:"version,omitempty"`
}
type LinuxProfile struct {
	AdminUsername	*string			`json:"adminUsername,omitempty"`
	SSH		*SSHConfiguration	`json:"ssh,omitempty"`
}
type ListResult struct {
	autorest.Response	`json:"-"`
	Value			*[]ContainerService	`json:"value,omitempty"`
	NextLink		*string			`json:"nextLink,omitempty"`
}
type ListResultIterator struct {
	i	int
	page	ListResultPage
}

func (iter *ListResultIterator) Next() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err := iter.page.Next()
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}
func (iter ListResultIterator) NotDone() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}
func (iter ListResultIterator) Response() ListResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return iter.page.Response()
}
func (iter ListResultIterator) Value() ContainerService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !iter.page.NotDone() {
		return ContainerService{}
	}
	return iter.page.Values()[iter.i]
}
func (lr ListResult) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return lr.Value == nil || len(*lr.Value) == 0
}
func (lr ListResult) listResultPreparer() (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if lr.NextLink == nil || len(to.String(lr.NextLink)) < 1 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{}, autorest.AsJSON(), autorest.AsGet(), autorest.WithBaseURL(to.String(lr.NextLink)))
}

type ListResultPage struct {
	fn	func(ListResult) (ListResult, error)
	lr	ListResult
}

func (page *ListResultPage) Next() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	next, err := page.fn(page.lr)
	if err != nil {
		return err
	}
	page.lr = next
	return nil
}
func (page ListResultPage) NotDone() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !page.lr.IsEmpty()
}
func (page ListResultPage) Response() ListResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return page.lr
}
func (page ListResultPage) Values() []ContainerService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if page.lr.IsEmpty() {
		return nil
	}
	return *page.lr.Value
}

type ManagedCluster struct {
	autorest.Response		`json:"-"`
	*ManagedClusterProperties	`json:"properties,omitempty"`
	ID				*string			`json:"id,omitempty"`
	Name				*string			`json:"name,omitempty"`
	Type				*string			`json:"type,omitempty"`
	Location			*string			`json:"location,omitempty"`
	Tags				map[string]*string	`json:"tags"`
}

func (mc ManagedCluster) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if mc.ManagedClusterProperties != nil {
		objectMap["properties"] = mc.ManagedClusterProperties
	}
	if mc.ID != nil {
		objectMap["id"] = mc.ID
	}
	if mc.Name != nil {
		objectMap["name"] = mc.Name
	}
	if mc.Type != nil {
		objectMap["type"] = mc.Type
	}
	if mc.Location != nil {
		objectMap["location"] = mc.Location
	}
	if mc.Tags != nil {
		objectMap["tags"] = mc.Tags
	}
	return json.Marshal(objectMap)
}
func (mc *ManagedCluster) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var managedClusterProperties ManagedClusterProperties
				err = json.Unmarshal(*v, &managedClusterProperties)
				if err != nil {
					return err
				}
				mc.ManagedClusterProperties = &managedClusterProperties
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				mc.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				mc.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				mc.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				mc.Location = &location
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				mc.Tags = tags
			}
		}
	}
	return nil
}

type ManagedClusterAADProfile struct {
	ClientAppID	*string	`json:"clientAppID,omitempty"`
	ServerAppID	*string	`json:"serverAppID,omitempty"`
	ServerAppSecret	*string	`json:"serverAppSecret,omitempty"`
	TenantID	*string	`json:"tenantID,omitempty"`
}
type ManagedClusterAccessProfile struct {
	autorest.Response	`json:"-"`
	*AccessProfile		`json:"properties,omitempty"`
	ID			*string			`json:"id,omitempty"`
	Name			*string			`json:"name,omitempty"`
	Type			*string			`json:"type,omitempty"`
	Location		*string			`json:"location,omitempty"`
	Tags			map[string]*string	`json:"tags"`
}

func (mcap ManagedClusterAccessProfile) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if mcap.AccessProfile != nil {
		objectMap["properties"] = mcap.AccessProfile
	}
	if mcap.ID != nil {
		objectMap["id"] = mcap.ID
	}
	if mcap.Name != nil {
		objectMap["name"] = mcap.Name
	}
	if mcap.Type != nil {
		objectMap["type"] = mcap.Type
	}
	if mcap.Location != nil {
		objectMap["location"] = mcap.Location
	}
	if mcap.Tags != nil {
		objectMap["tags"] = mcap.Tags
	}
	return json.Marshal(objectMap)
}
func (mcap *ManagedClusterAccessProfile) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var accessProfile AccessProfile
				err = json.Unmarshal(*v, &accessProfile)
				if err != nil {
					return err
				}
				mcap.AccessProfile = &accessProfile
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				mcap.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				mcap.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				mcap.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				mcap.Location = &location
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				mcap.Tags = tags
			}
		}
	}
	return nil
}

type ManagedClusterAddonProfile struct {
	Enabled	*bool			`json:"enabled,omitempty"`
	Config	map[string]*string	`json:"config"`
}

func (mcap ManagedClusterAddonProfile) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if mcap.Enabled != nil {
		objectMap["enabled"] = mcap.Enabled
	}
	if mcap.Config != nil {
		objectMap["config"] = mcap.Config
	}
	return json.Marshal(objectMap)
}

type ManagedClusterAgentPoolProfile struct {
	Name		*string			`json:"name,omitempty"`
	Count		*int32			`json:"count,omitempty"`
	VMSize		VMSizeTypes		`json:"vmSize,omitempty"`
	OsDiskSizeGB	*int32			`json:"osDiskSizeGB,omitempty"`
	DNSPrefix	*string			`json:"dnsPrefix,omitempty"`
	Fqdn		*string			`json:"fqdn,omitempty"`
	Ports		*[]int32		`json:"ports,omitempty"`
	StorageProfile	StorageProfileTypes	`json:"storageProfile,omitempty"`
	VnetSubnetID	*string			`json:"vnetSubnetID,omitempty"`
	MaxPods		*int32			`json:"maxPods,omitempty"`
	OsType		OSType			`json:"osType,omitempty"`
}
type ManagedClusterListResult struct {
	autorest.Response	`json:"-"`
	Value			*[]ManagedCluster	`json:"value,omitempty"`
	NextLink		*string			`json:"nextLink,omitempty"`
}
type ManagedClusterListResultIterator struct {
	i	int
	page	ManagedClusterListResultPage
}

func (iter *ManagedClusterListResultIterator) Next() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err := iter.page.Next()
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}
func (iter ManagedClusterListResultIterator) NotDone() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}
func (iter ManagedClusterListResultIterator) Response() ManagedClusterListResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return iter.page.Response()
}
func (iter ManagedClusterListResultIterator) Value() ManagedCluster {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !iter.page.NotDone() {
		return ManagedCluster{}
	}
	return iter.page.Values()[iter.i]
}
func (mclr ManagedClusterListResult) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mclr.Value == nil || len(*mclr.Value) == 0
}
func (mclr ManagedClusterListResult) managedClusterListResultPreparer() (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mclr.NextLink == nil || len(to.String(mclr.NextLink)) < 1 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{}, autorest.AsJSON(), autorest.AsGet(), autorest.WithBaseURL(to.String(mclr.NextLink)))
}

type ManagedClusterListResultPage struct {
	fn	func(ManagedClusterListResult) (ManagedClusterListResult, error)
	mclr	ManagedClusterListResult
}

func (page *ManagedClusterListResultPage) Next() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	next, err := page.fn(page.mclr)
	if err != nil {
		return err
	}
	page.mclr = next
	return nil
}
func (page ManagedClusterListResultPage) NotDone() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !page.mclr.IsEmpty()
}
func (page ManagedClusterListResultPage) Response() ManagedClusterListResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return page.mclr
}
func (page ManagedClusterListResultPage) Values() []ManagedCluster {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if page.mclr.IsEmpty() {
		return nil
	}
	return *page.mclr.Value
}

type ManagedClusterPoolUpgradeProfile struct {
	KubernetesVersion	*string		`json:"kubernetesVersion,omitempty"`
	Name			*string		`json:"name,omitempty"`
	OsType			OSType		`json:"osType,omitempty"`
	Upgrades		*[]string	`json:"upgrades,omitempty"`
}
type ManagedClusterProperties struct {
	ProvisioningState	*string					`json:"provisioningState,omitempty"`
	KubernetesVersion	*string					`json:"kubernetesVersion,omitempty"`
	DNSPrefix		*string					`json:"dnsPrefix,omitempty"`
	Fqdn			*string					`json:"fqdn,omitempty"`
	AgentPoolProfiles	*[]ManagedClusterAgentPoolProfile	`json:"agentPoolProfiles,omitempty"`
	LinuxProfile		*LinuxProfile				`json:"linuxProfile,omitempty"`
	ServicePrincipalProfile	*ServicePrincipalProfile		`json:"servicePrincipalProfile,omitempty"`
	AddonProfiles		map[string]*ManagedClusterAddonProfile	`json:"addonProfiles"`
	NodeResourceGroup	*string					`json:"nodeResourceGroup,omitempty"`
	EnableRBAC		*bool					`json:"enableRBAC,omitempty"`
	NetworkProfile		*NetworkProfile				`json:"networkProfile,omitempty"`
	AadProfile		*ManagedClusterAADProfile		`json:"aadProfile,omitempty"`
}

func (mcp ManagedClusterProperties) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if mcp.ProvisioningState != nil {
		objectMap["provisioningState"] = mcp.ProvisioningState
	}
	if mcp.KubernetesVersion != nil {
		objectMap["kubernetesVersion"] = mcp.KubernetesVersion
	}
	if mcp.DNSPrefix != nil {
		objectMap["dnsPrefix"] = mcp.DNSPrefix
	}
	if mcp.Fqdn != nil {
		objectMap["fqdn"] = mcp.Fqdn
	}
	if mcp.AgentPoolProfiles != nil {
		objectMap["agentPoolProfiles"] = mcp.AgentPoolProfiles
	}
	if mcp.LinuxProfile != nil {
		objectMap["linuxProfile"] = mcp.LinuxProfile
	}
	if mcp.ServicePrincipalProfile != nil {
		objectMap["servicePrincipalProfile"] = mcp.ServicePrincipalProfile
	}
	if mcp.AddonProfiles != nil {
		objectMap["addonProfiles"] = mcp.AddonProfiles
	}
	if mcp.NodeResourceGroup != nil {
		objectMap["nodeResourceGroup"] = mcp.NodeResourceGroup
	}
	if mcp.EnableRBAC != nil {
		objectMap["enableRBAC"] = mcp.EnableRBAC
	}
	if mcp.NetworkProfile != nil {
		objectMap["networkProfile"] = mcp.NetworkProfile
	}
	if mcp.AadProfile != nil {
		objectMap["aadProfile"] = mcp.AadProfile
	}
	return json.Marshal(objectMap)
}

type ManagedClustersCreateOrUpdateFuture struct{ azure.Future }

func (future *ManagedClustersCreateOrUpdateFuture) Result(client ManagedClustersClient) (mc ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersCreateOrUpdateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("containerservice.ManagedClustersCreateOrUpdateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if mc.Response.Response, err = future.GetResult(sender); err == nil && mc.Response.Response.StatusCode != http.StatusNoContent {
		mc, err = client.CreateOrUpdateResponder(mc.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersCreateOrUpdateFuture", "Result", mc.Response.Response, "Failure responding to request")
		}
	}
	return
}

type ManagedClustersDeleteFuture struct{ azure.Future }

func (future *ManagedClustersDeleteFuture) Result(client ManagedClustersClient) (ar autorest.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersDeleteFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("containerservice.ManagedClustersDeleteFuture")
		return
	}
	ar.Response = future.Response()
	return
}

type ManagedClustersUpdateTagsFuture struct{ azure.Future }

func (future *ManagedClustersUpdateTagsFuture) Result(client ManagedClustersClient) (mc ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersUpdateTagsFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("containerservice.ManagedClustersUpdateTagsFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if mc.Response.Response, err = future.GetResult(sender); err == nil && mc.Response.Response.StatusCode != http.StatusNoContent {
		mc, err = client.UpdateTagsResponder(mc.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersUpdateTagsFuture", "Result", mc.Response.Response, "Failure responding to request")
		}
	}
	return
}

type ManagedClusterUpgradeProfile struct {
	autorest.Response			`json:"-"`
	ID					*string	`json:"id,omitempty"`
	Name					*string	`json:"name,omitempty"`
	Type					*string	`json:"type,omitempty"`
	*ManagedClusterUpgradeProfileProperties	`json:"properties,omitempty"`
}

func (mcup ManagedClusterUpgradeProfile) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if mcup.ID != nil {
		objectMap["id"] = mcup.ID
	}
	if mcup.Name != nil {
		objectMap["name"] = mcup.Name
	}
	if mcup.Type != nil {
		objectMap["type"] = mcup.Type
	}
	if mcup.ManagedClusterUpgradeProfileProperties != nil {
		objectMap["properties"] = mcup.ManagedClusterUpgradeProfileProperties
	}
	return json.Marshal(objectMap)
}
func (mcup *ManagedClusterUpgradeProfile) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				mcup.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				mcup.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				mcup.Type = &typeVar
			}
		case "properties":
			if v != nil {
				var managedClusterUpgradeProfileProperties ManagedClusterUpgradeProfileProperties
				err = json.Unmarshal(*v, &managedClusterUpgradeProfileProperties)
				if err != nil {
					return err
				}
				mcup.ManagedClusterUpgradeProfileProperties = &managedClusterUpgradeProfileProperties
			}
		}
	}
	return nil
}

type ManagedClusterUpgradeProfileProperties struct {
	ControlPlaneProfile	*ManagedClusterPoolUpgradeProfile	`json:"controlPlaneProfile,omitempty"`
	AgentPoolProfiles	*[]ManagedClusterPoolUpgradeProfile	`json:"agentPoolProfiles,omitempty"`
}
type MasterProfile struct {
	Count				*int32			`json:"count,omitempty"`
	DNSPrefix			*string			`json:"dnsPrefix,omitempty"`
	VMSize				VMSizeTypes		`json:"vmSize,omitempty"`
	OsDiskSizeGB			*int32			`json:"osDiskSizeGB,omitempty"`
	VnetSubnetID			*string			`json:"vnetSubnetID,omitempty"`
	FirstConsecutiveStaticIP	*string			`json:"firstConsecutiveStaticIP,omitempty"`
	StorageProfile			StorageProfileTypes	`json:"storageProfile,omitempty"`
	Fqdn				*string			`json:"fqdn,omitempty"`
}
type NetworkProfile struct {
	NetworkPlugin		NetworkPlugin	`json:"networkPlugin,omitempty"`
	NetworkPolicy		NetworkPolicy	`json:"networkPolicy,omitempty"`
	PodCidr			*string		`json:"podCidr,omitempty"`
	ServiceCidr		*string		`json:"serviceCidr,omitempty"`
	DNSServiceIP		*string		`json:"dnsServiceIP,omitempty"`
	DockerBridgeCidr	*string		`json:"dockerBridgeCidr,omitempty"`
}
type OperationListResult struct {
	autorest.Response	`json:"-"`
	Value			*[]OperationValue	`json:"value,omitempty"`
}
type OperationValue struct {
	Origin			*string	`json:"origin,omitempty"`
	Name			*string	`json:"name,omitempty"`
	*OperationValueDisplay	`json:"display,omitempty"`
}

func (ov OperationValue) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if ov.Origin != nil {
		objectMap["origin"] = ov.Origin
	}
	if ov.Name != nil {
		objectMap["name"] = ov.Name
	}
	if ov.OperationValueDisplay != nil {
		objectMap["display"] = ov.OperationValueDisplay
	}
	return json.Marshal(objectMap)
}
func (ov *OperationValue) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "origin":
			if v != nil {
				var origin string
				err = json.Unmarshal(*v, &origin)
				if err != nil {
					return err
				}
				ov.Origin = &origin
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				ov.Name = &name
			}
		case "display":
			if v != nil {
				var operationValueDisplay OperationValueDisplay
				err = json.Unmarshal(*v, &operationValueDisplay)
				if err != nil {
					return err
				}
				ov.OperationValueDisplay = &operationValueDisplay
			}
		}
	}
	return nil
}

type OperationValueDisplay struct {
	Operation	*string	`json:"operation,omitempty"`
	Resource	*string	`json:"resource,omitempty"`
	Description	*string	`json:"description,omitempty"`
	Provider	*string	`json:"provider,omitempty"`
}
type OrchestratorProfile struct {
	OrchestratorType	*string	`json:"orchestratorType,omitempty"`
	OrchestratorVersion	*string	`json:"orchestratorVersion,omitempty"`
}
type OrchestratorProfileType struct {
	OrchestratorType	OrchestratorTypes	`json:"orchestratorType,omitempty"`
	OrchestratorVersion	*string			`json:"orchestratorVersion,omitempty"`
}
type OrchestratorVersionProfile struct {
	OrchestratorType	*string			`json:"orchestratorType,omitempty"`
	OrchestratorVersion	*string			`json:"orchestratorVersion,omitempty"`
	Default			*bool			`json:"default,omitempty"`
	Upgrades		*[]OrchestratorProfile	`json:"upgrades,omitempty"`
}
type OrchestratorVersionProfileListResult struct {
	autorest.Response			`json:"-"`
	ID					*string	`json:"id,omitempty"`
	Name					*string	`json:"name,omitempty"`
	Type					*string	`json:"type,omitempty"`
	*OrchestratorVersionProfileProperties	`json:"properties,omitempty"`
}

func (ovplr OrchestratorVersionProfileListResult) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if ovplr.ID != nil {
		objectMap["id"] = ovplr.ID
	}
	if ovplr.Name != nil {
		objectMap["name"] = ovplr.Name
	}
	if ovplr.Type != nil {
		objectMap["type"] = ovplr.Type
	}
	if ovplr.OrchestratorVersionProfileProperties != nil {
		objectMap["properties"] = ovplr.OrchestratorVersionProfileProperties
	}
	return json.Marshal(objectMap)
}
func (ovplr *OrchestratorVersionProfileListResult) UnmarshalJSON(body []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				ovplr.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				ovplr.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				ovplr.Type = &typeVar
			}
		case "properties":
			if v != nil {
				var orchestratorVersionProfileProperties OrchestratorVersionProfileProperties
				err = json.Unmarshal(*v, &orchestratorVersionProfileProperties)
				if err != nil {
					return err
				}
				ovplr.OrchestratorVersionProfileProperties = &orchestratorVersionProfileProperties
			}
		}
	}
	return nil
}

type OrchestratorVersionProfileProperties struct {
	Orchestrators *[]OrchestratorVersionProfile `json:"orchestrators,omitempty"`
}
type Properties struct {
	ProvisioningState	*string				`json:"provisioningState,omitempty"`
	OrchestratorProfile	*OrchestratorProfileType	`json:"orchestratorProfile,omitempty"`
	CustomProfile		*CustomProfile			`json:"customProfile,omitempty"`
	ServicePrincipalProfile	*ServicePrincipalProfile	`json:"servicePrincipalProfile,omitempty"`
	MasterProfile		*MasterProfile			`json:"masterProfile,omitempty"`
	AgentPoolProfiles	*[]AgentPoolProfile		`json:"agentPoolProfiles,omitempty"`
	WindowsProfile		*WindowsProfile			`json:"windowsProfile,omitempty"`
	LinuxProfile		*LinuxProfile			`json:"linuxProfile,omitempty"`
	DiagnosticsProfile	*DiagnosticsProfile		`json:"diagnosticsProfile,omitempty"`
}
type Resource struct {
	ID		*string			`json:"id,omitempty"`
	Name		*string			`json:"name,omitempty"`
	Type		*string			`json:"type,omitempty"`
	Location	*string			`json:"location,omitempty"`
	Tags		map[string]*string	`json:"tags"`
}

func (r Resource) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if r.ID != nil {
		objectMap["id"] = r.ID
	}
	if r.Name != nil {
		objectMap["name"] = r.Name
	}
	if r.Type != nil {
		objectMap["type"] = r.Type
	}
	if r.Location != nil {
		objectMap["location"] = r.Location
	}
	if r.Tags != nil {
		objectMap["tags"] = r.Tags
	}
	return json.Marshal(objectMap)
}

type ServicePrincipalProfile struct {
	ClientID		*string			`json:"clientId,omitempty"`
	Secret			*string			`json:"secret,omitempty"`
	KeyVaultSecretRef	*KeyVaultSecretRef	`json:"keyVaultSecretRef,omitempty"`
}
type SSHConfiguration struct {
	PublicKeys *[]SSHPublicKey `json:"publicKeys,omitempty"`
}
type SSHPublicKey struct {
	KeyData *string `json:"keyData,omitempty"`
}
type TagsObject struct {
	Tags map[string]*string `json:"tags"`
}

func (toVar TagsObject) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	objectMap := make(map[string]interface{})
	if toVar.Tags != nil {
		objectMap["tags"] = toVar.Tags
	}
	return json.Marshal(objectMap)
}

type VMDiagnostics struct {
	Enabled		*bool	`json:"enabled,omitempty"`
	StorageURI	*string	`json:"storageUri,omitempty"`
}
type WindowsProfile struct {
	AdminUsername	*string	`json:"adminUsername,omitempty"`
	AdminPassword	*string	`json:"adminPassword,omitempty"`
}
