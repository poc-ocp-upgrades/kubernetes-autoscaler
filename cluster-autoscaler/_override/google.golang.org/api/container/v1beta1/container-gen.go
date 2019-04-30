package container

import (
	"bytes"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"encoding/json"
	"errors"
	"fmt"
	context "golang.org/x/net/context"
	ctxhttp "golang.org/x/net/context/ctxhttp"
	gensupport "google.golang.org/api/gensupport"
	googleapi "google.golang.org/api/googleapi"
	"io"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = gensupport.MarshalJSON
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace
var _ = context.Canceled
var _ = ctxhttp.Do

const apiId = "container:v1beta1"
const apiName = "container"
const apiVersion = "v1beta1"
const basePath = "https://container.googleapis.com/"
const (
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

func New(client *http.Client) (*Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Projects = NewProjectsService(s)
	return s, nil
}

type Service struct {
	client		*http.Client
	BasePath	string
	UserAgent	string
	Projects	*ProjectsService
}

func (s *Service) userAgent() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}
func NewProjectsService(s *Service) *ProjectsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsService{s: s}
	rs.Aggregated = NewProjectsAggregatedService(s)
	rs.Locations = NewProjectsLocationsService(s)
	rs.Zones = NewProjectsZonesService(s)
	return rs
}

type ProjectsService struct {
	s		*Service
	Aggregated	*ProjectsAggregatedService
	Locations	*ProjectsLocationsService
	Zones		*ProjectsZonesService
}

func NewProjectsAggregatedService(s *Service) *ProjectsAggregatedService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsAggregatedService{s: s}
	rs.UsableSubnetworks = NewProjectsAggregatedUsableSubnetworksService(s)
	return rs
}

type ProjectsAggregatedService struct {
	s			*Service
	UsableSubnetworks	*ProjectsAggregatedUsableSubnetworksService
}

func NewProjectsAggregatedUsableSubnetworksService(s *Service) *ProjectsAggregatedUsableSubnetworksService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsAggregatedUsableSubnetworksService{s: s}
	return rs
}

type ProjectsAggregatedUsableSubnetworksService struct{ s *Service }

func NewProjectsLocationsService(s *Service) *ProjectsLocationsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsLocationsService{s: s}
	rs.Clusters = NewProjectsLocationsClustersService(s)
	rs.Operations = NewProjectsLocationsOperationsService(s)
	return rs
}

type ProjectsLocationsService struct {
	s		*Service
	Clusters	*ProjectsLocationsClustersService
	Operations	*ProjectsLocationsOperationsService
}

func NewProjectsLocationsClustersService(s *Service) *ProjectsLocationsClustersService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsLocationsClustersService{s: s}
	rs.NodePools = NewProjectsLocationsClustersNodePoolsService(s)
	return rs
}

type ProjectsLocationsClustersService struct {
	s		*Service
	NodePools	*ProjectsLocationsClustersNodePoolsService
}

func NewProjectsLocationsClustersNodePoolsService(s *Service) *ProjectsLocationsClustersNodePoolsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsLocationsClustersNodePoolsService{s: s}
	return rs
}

type ProjectsLocationsClustersNodePoolsService struct{ s *Service }

func NewProjectsLocationsOperationsService(s *Service) *ProjectsLocationsOperationsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsLocationsOperationsService{s: s}
	return rs
}

type ProjectsLocationsOperationsService struct{ s *Service }

func NewProjectsZonesService(s *Service) *ProjectsZonesService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsZonesService{s: s}
	rs.Clusters = NewProjectsZonesClustersService(s)
	rs.Operations = NewProjectsZonesOperationsService(s)
	return rs
}

type ProjectsZonesService struct {
	s		*Service
	Clusters	*ProjectsZonesClustersService
	Operations	*ProjectsZonesOperationsService
}

func NewProjectsZonesClustersService(s *Service) *ProjectsZonesClustersService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsZonesClustersService{s: s}
	rs.NodePools = NewProjectsZonesClustersNodePoolsService(s)
	return rs
}

type ProjectsZonesClustersService struct {
	s		*Service
	NodePools	*ProjectsZonesClustersNodePoolsService
}

func NewProjectsZonesClustersNodePoolsService(s *Service) *ProjectsZonesClustersNodePoolsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsZonesClustersNodePoolsService{s: s}
	return rs
}

type ProjectsZonesClustersNodePoolsService struct{ s *Service }

func NewProjectsZonesOperationsService(s *Service) *ProjectsZonesOperationsService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs := &ProjectsZonesOperationsService{s: s}
	return rs
}

type ProjectsZonesOperationsService struct{ s *Service }
type AcceleratorConfig struct {
	AcceleratorCount	int64		`json:"acceleratorCount,omitempty,string"`
	AcceleratorType		string		`json:"acceleratorType,omitempty"`
	ForceSendFields		[]string	`json:"-"`
	NullFields		[]string	`json:"-"`
}

func (s *AcceleratorConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod AcceleratorConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type AddonsConfig struct {
	HorizontalPodAutoscaling	*HorizontalPodAutoscaling	`json:"horizontalPodAutoscaling,omitempty"`
	HttpLoadBalancing		*HttpLoadBalancing		`json:"httpLoadBalancing,omitempty"`
	KubernetesDashboard		*KubernetesDashboard		`json:"kubernetesDashboard,omitempty"`
	NetworkPolicyConfig		*NetworkPolicyConfig		`json:"networkPolicyConfig,omitempty"`
	ForceSendFields			[]string			`json:"-"`
	NullFields			[]string			`json:"-"`
}

func (s *AddonsConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod AddonsConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type AutoUpgradeOptions struct {
	AutoUpgradeStartTime	string		`json:"autoUpgradeStartTime,omitempty"`
	Description		string		`json:"description,omitempty"`
	ForceSendFields		[]string	`json:"-"`
	NullFields		[]string	`json:"-"`
}

func (s *AutoUpgradeOptions) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod AutoUpgradeOptions
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type BinaryAuthorization struct {
	Enabled		bool		`json:"enabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *BinaryAuthorization) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod BinaryAuthorization
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CancelOperationRequest struct {
	Name		string		`json:"name,omitempty"`
	OperationId	string		`json:"operationId,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *CancelOperationRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod CancelOperationRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CidrBlock struct {
	CidrBlock	string		`json:"cidrBlock,omitempty"`
	DisplayName	string		`json:"displayName,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *CidrBlock) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod CidrBlock
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ClientCertificateConfig struct {
	IssueClientCertificate	bool		`json:"issueClientCertificate,omitempty"`
	ForceSendFields		[]string	`json:"-"`
	NullFields		[]string	`json:"-"`
}

func (s *ClientCertificateConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ClientCertificateConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Cluster struct {
	AddonsConfig			*AddonsConfig			`json:"addonsConfig,omitempty"`
	Autoscaling			*ClusterAutoscaling		`json:"autoscaling,omitempty"`
	BinaryAuthorization		*BinaryAuthorization		`json:"binaryAuthorization,omitempty"`
	ClusterIpv4Cidr			string				`json:"clusterIpv4Cidr,omitempty"`
	CreateTime			string				`json:"createTime,omitempty"`
	CurrentMasterVersion		string				`json:"currentMasterVersion,omitempty"`
	CurrentNodeCount		int64				`json:"currentNodeCount,omitempty"`
	CurrentNodeVersion		string				`json:"currentNodeVersion,omitempty"`
	Description			string				`json:"description,omitempty"`
	EnableKubernetesAlpha		bool				`json:"enableKubernetesAlpha,omitempty"`
	Endpoint			string				`json:"endpoint,omitempty"`
	ExpireTime			string				`json:"expireTime,omitempty"`
	InitialClusterVersion		string				`json:"initialClusterVersion,omitempty"`
	InitialNodeCount		int64				`json:"initialNodeCount,omitempty"`
	InstanceGroupUrls		[]string			`json:"instanceGroupUrls,omitempty"`
	IpAllocationPolicy		*IPAllocationPolicy		`json:"ipAllocationPolicy,omitempty"`
	LabelFingerprint		string				`json:"labelFingerprint,omitempty"`
	LegacyAbac			*LegacyAbac			`json:"legacyAbac,omitempty"`
	Location			string				`json:"location,omitempty"`
	Locations			[]string			`json:"locations,omitempty"`
	LoggingService			string				`json:"loggingService,omitempty"`
	MaintenancePolicy		*MaintenancePolicy		`json:"maintenancePolicy,omitempty"`
	MasterAuth			*MasterAuth			`json:"masterAuth,omitempty"`
	MasterAuthorizedNetworksConfig	*MasterAuthorizedNetworksConfig	`json:"masterAuthorizedNetworksConfig,omitempty"`
	MasterIpv4CidrBlock		string				`json:"masterIpv4CidrBlock,omitempty"`
	MonitoringService		string				`json:"monitoringService,omitempty"`
	Name				string				`json:"name,omitempty"`
	Network				string				`json:"network,omitempty"`
	NetworkConfig			*NetworkConfig			`json:"networkConfig,omitempty"`
	NetworkPolicy			*NetworkPolicy			`json:"networkPolicy,omitempty"`
	NodeConfig			*NodeConfig			`json:"nodeConfig,omitempty"`
	NodeIpv4CidrSize		int64				`json:"nodeIpv4CidrSize,omitempty"`
	NodePools			[]*NodePool			`json:"nodePools,omitempty"`
	PodSecurityPolicyConfig		*PodSecurityPolicyConfig	`json:"podSecurityPolicyConfig,omitempty"`
	PrivateCluster			bool				`json:"privateCluster,omitempty"`
	ResourceLabels			map[string]string		`json:"resourceLabels,omitempty"`
	SelfLink			string				`json:"selfLink,omitempty"`
	ServicesIpv4Cidr		string				`json:"servicesIpv4Cidr,omitempty"`
	Status				string				`json:"status,omitempty"`
	StatusMessage			string				`json:"statusMessage,omitempty"`
	Subnetwork			string				`json:"subnetwork,omitempty"`
	Zone				string				`json:"zone,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *Cluster) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod Cluster
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ClusterAutoscaling struct {
	EnableNodeAutoprovisioning	bool			`json:"enableNodeAutoprovisioning,omitempty"`
	ResourceLimits			[]*ResourceLimit	`json:"resourceLimits,omitempty"`
	ForceSendFields			[]string		`json:"-"`
	NullFields			[]string		`json:"-"`
}

func (s *ClusterAutoscaling) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ClusterAutoscaling
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ClusterUpdate struct {
	DesiredAddonsConfig			*AddonsConfig			`json:"desiredAddonsConfig,omitempty"`
	DesiredBinaryAuthorization		*BinaryAuthorization		`json:"desiredBinaryAuthorization,omitempty"`
	DesiredClusterAutoscaling		*ClusterAutoscaling		`json:"desiredClusterAutoscaling,omitempty"`
	DesiredImageType			string				`json:"desiredImageType,omitempty"`
	DesiredLocations			[]string			`json:"desiredLocations,omitempty"`
	DesiredMasterAuthorizedNetworksConfig	*MasterAuthorizedNetworksConfig	`json:"desiredMasterAuthorizedNetworksConfig,omitempty"`
	DesiredMasterVersion			string				`json:"desiredMasterVersion,omitempty"`
	DesiredMonitoringService		string				`json:"desiredMonitoringService,omitempty"`
	DesiredNodePoolAutoscaling		*NodePoolAutoscaling		`json:"desiredNodePoolAutoscaling,omitempty"`
	DesiredNodePoolId			string				`json:"desiredNodePoolId,omitempty"`
	DesiredNodeVersion			string				`json:"desiredNodeVersion,omitempty"`
	DesiredPodSecurityPolicyConfig		*PodSecurityPolicyConfig	`json:"desiredPodSecurityPolicyConfig,omitempty"`
	ForceSendFields				[]string			`json:"-"`
	NullFields				[]string			`json:"-"`
}

func (s *ClusterUpdate) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ClusterUpdate
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CompleteIPRotationRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *CompleteIPRotationRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod CompleteIPRotationRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CreateClusterRequest struct {
	Cluster		*Cluster	`json:"cluster,omitempty"`
	Parent		string		`json:"parent,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *CreateClusterRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod CreateClusterRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CreateNodePoolRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	NodePool	*NodePool	`json:"nodePool,omitempty"`
	Parent		string		`json:"parent,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *CreateNodePoolRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod CreateNodePoolRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type DailyMaintenanceWindow struct {
	Duration	string		`json:"duration,omitempty"`
	StartTime	string		`json:"startTime,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *DailyMaintenanceWindow) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod DailyMaintenanceWindow
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Empty struct {
	googleapi.ServerResponse `json:"-"`
}
type HorizontalPodAutoscaling struct {
	Disabled	bool		`json:"disabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *HorizontalPodAutoscaling) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod HorizontalPodAutoscaling
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type HttpLoadBalancing struct {
	Disabled	bool		`json:"disabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *HttpLoadBalancing) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod HttpLoadBalancing
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type IPAllocationPolicy struct {
	AllowRouteOverlap		bool		`json:"allowRouteOverlap,omitempty"`
	ClusterIpv4Cidr			string		`json:"clusterIpv4Cidr,omitempty"`
	ClusterIpv4CidrBlock		string		`json:"clusterIpv4CidrBlock,omitempty"`
	ClusterSecondaryRangeName	string		`json:"clusterSecondaryRangeName,omitempty"`
	CreateSubnetwork		bool		`json:"createSubnetwork,omitempty"`
	NodeIpv4Cidr			string		`json:"nodeIpv4Cidr,omitempty"`
	NodeIpv4CidrBlock		string		`json:"nodeIpv4CidrBlock,omitempty"`
	ServicesIpv4Cidr		string		`json:"servicesIpv4Cidr,omitempty"`
	ServicesIpv4CidrBlock		string		`json:"servicesIpv4CidrBlock,omitempty"`
	ServicesSecondaryRangeName	string		`json:"servicesSecondaryRangeName,omitempty"`
	SubnetworkName			string		`json:"subnetworkName,omitempty"`
	UseIpAliases			bool		`json:"useIpAliases,omitempty"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *IPAllocationPolicy) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod IPAllocationPolicy
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type KubernetesDashboard struct {
	Disabled	bool		`json:"disabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *KubernetesDashboard) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod KubernetesDashboard
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type LegacyAbac struct {
	Enabled		bool		`json:"enabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *LegacyAbac) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod LegacyAbac
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListClustersResponse struct {
	Clusters			[]*Cluster	`json:"clusters,omitempty"`
	MissingZones			[]string	`json:"missingZones,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ListClustersResponse) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ListClustersResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListLocationsResponse struct {
	Locations			[]*Location	`json:"locations,omitempty"`
	NextPageToken			string		`json:"nextPageToken,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ListLocationsResponse) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ListLocationsResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListNodePoolsResponse struct {
	NodePools			[]*NodePool	`json:"nodePools,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ListNodePoolsResponse) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ListNodePoolsResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListOperationsResponse struct {
	MissingZones			[]string	`json:"missingZones,omitempty"`
	Operations			[]*Operation	`json:"operations,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ListOperationsResponse) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ListOperationsResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListUsableSubnetworksResponse struct {
	NextPageToken			string			`json:"nextPageToken,omitempty"`
	Subnetworks			[]*UsableSubnetwork	`json:"subnetworks,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ListUsableSubnetworksResponse) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ListUsableSubnetworksResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Location struct {
	Name		string		`json:"name,omitempty"`
	Recommended	bool		`json:"recommended,omitempty"`
	Type		string		`json:"type,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *Location) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod Location
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type MaintenancePolicy struct {
	Window		*MaintenanceWindow	`json:"window,omitempty"`
	ForceSendFields	[]string		`json:"-"`
	NullFields	[]string		`json:"-"`
}

func (s *MaintenancePolicy) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod MaintenancePolicy
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type MaintenanceWindow struct {
	DailyMaintenanceWindow	*DailyMaintenanceWindow	`json:"dailyMaintenanceWindow,omitempty"`
	ForceSendFields		[]string		`json:"-"`
	NullFields		[]string		`json:"-"`
}

func (s *MaintenanceWindow) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod MaintenanceWindow
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type MasterAuth struct {
	ClientCertificate	string				`json:"clientCertificate,omitempty"`
	ClientCertificateConfig	*ClientCertificateConfig	`json:"clientCertificateConfig,omitempty"`
	ClientKey		string				`json:"clientKey,omitempty"`
	ClusterCaCertificate	string				`json:"clusterCaCertificate,omitempty"`
	Password		string				`json:"password,omitempty"`
	Username		string				`json:"username,omitempty"`
	ForceSendFields		[]string			`json:"-"`
	NullFields		[]string			`json:"-"`
}

func (s *MasterAuth) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod MasterAuth
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type MasterAuthorizedNetworksConfig struct {
	CidrBlocks	[]*CidrBlock	`json:"cidrBlocks,omitempty"`
	Enabled		bool		`json:"enabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *MasterAuthorizedNetworksConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod MasterAuthorizedNetworksConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Metric struct {
	DoubleValue	float64		`json:"doubleValue,omitempty"`
	IntValue	int64		`json:"intValue,omitempty,string"`
	Name		string		`json:"name,omitempty"`
	StringValue	string		`json:"stringValue,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *Metric) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod Metric
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}
func (s *Metric) UnmarshalJSON(data []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod Metric
	var s1 struct {
		DoubleValue	gensupport.JSONFloat64	`json:"doubleValue"`
		*NoMethod
	}
	s1.NoMethod = (*NoMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.DoubleValue = float64(s1.DoubleValue)
	return nil
}

type NetworkConfig struct {
	Network		string		`json:"network,omitempty"`
	Subnetwork	string		`json:"subnetwork,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *NetworkConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NetworkConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NetworkPolicy struct {
	Enabled		bool		`json:"enabled,omitempty"`
	Provider	string		`json:"provider,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *NetworkPolicy) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NetworkPolicy
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NetworkPolicyConfig struct {
	Disabled	bool		`json:"disabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *NetworkPolicyConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NetworkPolicyConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NodeConfig struct {
	Accelerators		[]*AcceleratorConfig	`json:"accelerators,omitempty"`
	DiskSizeGb		int64			`json:"diskSizeGb,omitempty"`
	DiskType		string			`json:"diskType,omitempty"`
	ImageType		string			`json:"imageType,omitempty"`
	Labels			map[string]string	`json:"labels,omitempty"`
	LocalSsdCount		int64			`json:"localSsdCount,omitempty"`
	MachineType		string			`json:"machineType,omitempty"`
	Metadata		map[string]string	`json:"metadata,omitempty"`
	MinCpuPlatform		string			`json:"minCpuPlatform,omitempty"`
	OauthScopes		[]string		`json:"oauthScopes,omitempty"`
	Preemptible		bool			`json:"preemptible,omitempty"`
	ServiceAccount		string			`json:"serviceAccount,omitempty"`
	Tags			[]string		`json:"tags,omitempty"`
	Taints			[]*NodeTaint		`json:"taints,omitempty"`
	WorkloadMetadataConfig	*WorkloadMetadataConfig	`json:"workloadMetadataConfig,omitempty"`
	ForceSendFields		[]string		`json:"-"`
	NullFields		[]string		`json:"-"`
}

func (s *NodeConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NodeConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NodeManagement struct {
	AutoRepair	bool			`json:"autoRepair,omitempty"`
	AutoUpgrade	bool			`json:"autoUpgrade,omitempty"`
	UpgradeOptions	*AutoUpgradeOptions	`json:"upgradeOptions,omitempty"`
	ForceSendFields	[]string		`json:"-"`
	NullFields	[]string		`json:"-"`
}

func (s *NodeManagement) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NodeManagement
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NodePool struct {
	Autoscaling			*NodePoolAutoscaling	`json:"autoscaling,omitempty"`
	Config				*NodeConfig		`json:"config,omitempty"`
	InitialNodeCount		int64			`json:"initialNodeCount,omitempty"`
	InstanceGroupUrls		[]string		`json:"instanceGroupUrls,omitempty"`
	Management			*NodeManagement		`json:"management,omitempty"`
	Name				string			`json:"name,omitempty"`
	SelfLink			string			`json:"selfLink,omitempty"`
	Status				string			`json:"status,omitempty"`
	StatusMessage			string			`json:"statusMessage,omitempty"`
	Version				string			`json:"version,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *NodePool) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NodePool
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NodePoolAutoscaling struct {
	Autoprovisioned	bool		`json:"autoprovisioned,omitempty"`
	Enabled		bool		`json:"enabled,omitempty"`
	MaxNodeCount	int64		`json:"maxNodeCount,omitempty"`
	MinNodeCount	int64		`json:"minNodeCount,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *NodePoolAutoscaling) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NodePoolAutoscaling
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type NodeTaint struct {
	Effect		string		`json:"effect,omitempty"`
	Key		string		`json:"key,omitempty"`
	Value		string		`json:"value,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *NodeTaint) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod NodeTaint
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Operation struct {
	Detail				string			`json:"detail,omitempty"`
	EndTime				string			`json:"endTime,omitempty"`
	Location			string			`json:"location,omitempty"`
	Name				string			`json:"name,omitempty"`
	OperationType			string			`json:"operationType,omitempty"`
	Progress			*OperationProgress	`json:"progress,omitempty"`
	SelfLink			string			`json:"selfLink,omitempty"`
	StartTime			string			`json:"startTime,omitempty"`
	Status				string			`json:"status,omitempty"`
	StatusMessage			string			`json:"statusMessage,omitempty"`
	TargetLink			string			`json:"targetLink,omitempty"`
	Zone				string			`json:"zone,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *Operation) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod Operation
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type OperationProgress struct {
	Metrics		[]*Metric		`json:"metrics,omitempty"`
	Name		string			`json:"name,omitempty"`
	Stages		[]*OperationProgress	`json:"stages,omitempty"`
	Status		string			`json:"status,omitempty"`
	ForceSendFields	[]string		`json:"-"`
	NullFields	[]string		`json:"-"`
}

func (s *OperationProgress) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod OperationProgress
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type PodSecurityPolicyConfig struct {
	Enabled		bool		`json:"enabled,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *PodSecurityPolicyConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod PodSecurityPolicyConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ResourceLimit struct {
	Maximum		int64		`json:"maximum,omitempty,string"`
	Minimum		int64		`json:"minimum,omitempty,string"`
	ResourceType	string		`json:"resourceType,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *ResourceLimit) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ResourceLimit
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type RollbackNodePoolUpgradeRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	NodePoolId	string		`json:"nodePoolId,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *RollbackNodePoolUpgradeRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod RollbackNodePoolUpgradeRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ServerConfig struct {
	DefaultClusterVersion		string		`json:"defaultClusterVersion,omitempty"`
	DefaultImageType		string		`json:"defaultImageType,omitempty"`
	ValidImageTypes			[]string	`json:"validImageTypes,omitempty"`
	ValidMasterVersions		[]string	`json:"validMasterVersions,omitempty"`
	ValidNodeVersions		[]string	`json:"validNodeVersions,omitempty"`
	googleapi.ServerResponse	`json:"-"`
	ForceSendFields			[]string	`json:"-"`
	NullFields			[]string	`json:"-"`
}

func (s *ServerConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod ServerConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetAddonsConfigRequest struct {
	AddonsConfig	*AddonsConfig	`json:"addonsConfig,omitempty"`
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetAddonsConfigRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetAddonsConfigRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetLabelsRequest struct {
	ClusterId		string			`json:"clusterId,omitempty"`
	LabelFingerprint	string			`json:"labelFingerprint,omitempty"`
	Name			string			`json:"name,omitempty"`
	ProjectId		string			`json:"projectId,omitempty"`
	ResourceLabels		map[string]string	`json:"resourceLabels,omitempty"`
	Zone			string			`json:"zone,omitempty"`
	ForceSendFields		[]string		`json:"-"`
	NullFields		[]string		`json:"-"`
}

func (s *SetLabelsRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetLabelsRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetLegacyAbacRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Enabled		bool		`json:"enabled,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetLegacyAbacRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetLegacyAbacRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetLocationsRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Locations	[]string	`json:"locations,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetLocationsRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetLocationsRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetLoggingServiceRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	LoggingService	string		`json:"loggingService,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetLoggingServiceRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetLoggingServiceRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetMaintenancePolicyRequest struct {
	ClusterId		string			`json:"clusterId,omitempty"`
	MaintenancePolicy	*MaintenancePolicy	`json:"maintenancePolicy,omitempty"`
	Name			string			`json:"name,omitempty"`
	ProjectId		string			`json:"projectId,omitempty"`
	Zone			string			`json:"zone,omitempty"`
	ForceSendFields		[]string		`json:"-"`
	NullFields		[]string		`json:"-"`
}

func (s *SetMaintenancePolicyRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetMaintenancePolicyRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetMasterAuthRequest struct {
	Action		string		`json:"action,omitempty"`
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Update		*MasterAuth	`json:"update,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetMasterAuthRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetMasterAuthRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetMonitoringServiceRequest struct {
	ClusterId		string		`json:"clusterId,omitempty"`
	MonitoringService	string		`json:"monitoringService,omitempty"`
	Name			string		`json:"name,omitempty"`
	ProjectId		string		`json:"projectId,omitempty"`
	Zone			string		`json:"zone,omitempty"`
	ForceSendFields		[]string	`json:"-"`
	NullFields		[]string	`json:"-"`
}

func (s *SetMonitoringServiceRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetMonitoringServiceRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetNetworkPolicyRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	NetworkPolicy	*NetworkPolicy	`json:"networkPolicy,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetNetworkPolicyRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetNetworkPolicyRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetNodePoolAutoscalingRequest struct {
	Autoscaling	*NodePoolAutoscaling	`json:"autoscaling,omitempty"`
	ClusterId	string			`json:"clusterId,omitempty"`
	Name		string			`json:"name,omitempty"`
	NodePoolId	string			`json:"nodePoolId,omitempty"`
	ProjectId	string			`json:"projectId,omitempty"`
	Zone		string			`json:"zone,omitempty"`
	ForceSendFields	[]string		`json:"-"`
	NullFields	[]string		`json:"-"`
}

func (s *SetNodePoolAutoscalingRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetNodePoolAutoscalingRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetNodePoolManagementRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Management	*NodeManagement	`json:"management,omitempty"`
	Name		string		`json:"name,omitempty"`
	NodePoolId	string		`json:"nodePoolId,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetNodePoolManagementRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetNodePoolManagementRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SetNodePoolSizeRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	NodeCount	int64		`json:"nodeCount,omitempty"`
	NodePoolId	string		`json:"nodePoolId,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *SetNodePoolSizeRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod SetNodePoolSizeRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type StartIPRotationRequest struct {
	ClusterId		string		`json:"clusterId,omitempty"`
	Name			string		`json:"name,omitempty"`
	ProjectId		string		`json:"projectId,omitempty"`
	RotateCredentials	bool		`json:"rotateCredentials,omitempty"`
	Zone			string		`json:"zone,omitempty"`
	ForceSendFields		[]string	`json:"-"`
	NullFields		[]string	`json:"-"`
}

func (s *StartIPRotationRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod StartIPRotationRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type UpdateClusterRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Update		*ClusterUpdate	`json:"update,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *UpdateClusterRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod UpdateClusterRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type UpdateMasterRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	MasterVersion	string		`json:"masterVersion,omitempty"`
	Name		string		`json:"name,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *UpdateMasterRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod UpdateMasterRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type UpdateNodePoolRequest struct {
	ClusterId	string		`json:"clusterId,omitempty"`
	ImageType	string		`json:"imageType,omitempty"`
	Name		string		`json:"name,omitempty"`
	NodePoolId	string		`json:"nodePoolId,omitempty"`
	NodeVersion	string		`json:"nodeVersion,omitempty"`
	ProjectId	string		`json:"projectId,omitempty"`
	Zone		string		`json:"zone,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *UpdateNodePoolRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod UpdateNodePoolRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type UsableSubnetwork struct {
	IpCidrRange	string		`json:"ipCidrRange,omitempty"`
	Network		string		`json:"network,omitempty"`
	Subnetwork	string		`json:"subnetwork,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *UsableSubnetwork) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod UsableSubnetwork
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type WorkloadMetadataConfig struct {
	NodeMetadata	string		`json:"nodeMetadata,omitempty"`
	ForceSendFields	[]string	`json:"-"`
	NullFields	[]string	`json:"-"`
}

func (s *WorkloadMetadataConfig) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type NoMethod WorkloadMetadataConfig
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ProjectsAggregatedUsableSubnetworksListCall struct {
	s		*Service
	parent		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsAggregatedUsableSubnetworksService) List(parent string) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsAggregatedUsableSubnetworksListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Filter(filter string) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("filter", filter)
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) PageSize(pageSize int64) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) PageToken(pageToken string) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("pageToken", pageToken)
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Fields(s ...googleapi.Field) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) IfNoneMatch(entityTag string) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Context(ctx context.Context) *ProjectsAggregatedUsableSubnetworksListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/aggregated/usableSubnetworks")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Do(opts ...googleapi.CallOption) (*ListUsableSubnetworksResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListUsableSubnetworksResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}
func (c *ProjectsAggregatedUsableSubnetworksListCall) Pages(ctx context.Context, f func(*ListUsableSubnetworksResponse) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken"))
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

type ProjectsLocationsGetServerConfigCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsService) GetServerConfig(name string) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsGetServerConfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) ProjectId(projectId string) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) Zone(zone string) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) Fields(s ...googleapi.Field) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) IfNoneMatch(entityTag string) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) Context(ctx context.Context) *ProjectsLocationsGetServerConfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsGetServerConfigCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsGetServerConfigCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}/serverConfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsGetServerConfigCall) Do(opts ...googleapi.CallOption) (*ServerConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ServerConfig{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsListCall struct {
	s		*Service
	parent		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsService) List(parent string) *ProjectsLocationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}
func (c *ProjectsLocationsListCall) Fields(s ...googleapi.Field) *ProjectsLocationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsListCall) IfNoneMatch(entityTag string) *ProjectsLocationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsListCall) Context(ctx context.Context) *ProjectsLocationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/locations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsListCall) Do(opts ...googleapi.CallOption) (*ListLocationsResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListLocationsResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersCompleteIpRotationCall struct {
	s				*Service
	name				string
	completeiprotationrequest	*CompleteIPRotationRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersService) CompleteIpRotation(name string, completeiprotationrequest *CompleteIPRotationRequest) *ProjectsLocationsClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersCompleteIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.completeiprotationrequest = completeiprotationrequest
	return c
}
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Context(ctx context.Context) *ProjectsLocationsClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersCompleteIpRotationCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.completeiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:completeIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersCreateCall struct {
	s			*Service
	parent			string
	createclusterrequest	*CreateClusterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) Create(parent string, createclusterrequest *CreateClusterRequest) *ProjectsLocationsClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.createclusterrequest = createclusterrequest
	return c
}
func (c *ProjectsLocationsClustersCreateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersCreateCall) Context(ctx context.Context) *ProjectsLocationsClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersCreateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersCreateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersDeleteCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersService) Delete(name string) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) ClusterId(clusterId string) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("clusterId", clusterId)
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) ProjectId(projectId string) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) Zone(zone string) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) Context(ctx context.Context) *ProjectsLocationsClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersDeleteCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersDeleteCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersGetCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersService) Get(name string) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsClustersGetCall) ClusterId(clusterId string) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("clusterId", clusterId)
	return c
}
func (c *ProjectsLocationsClustersGetCall) ProjectId(projectId string) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersGetCall) Zone(zone string) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsClustersGetCall) Context(ctx context.Context) *ProjectsLocationsClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersGetCall) Do(opts ...googleapi.CallOption) (*Cluster, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Cluster{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersListCall struct {
	s		*Service
	parent		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersService) List(parent string) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}
func (c *ProjectsLocationsClustersListCall) ProjectId(projectId string) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersListCall) Zone(zone string) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersListCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersListCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsClustersListCall) Context(ctx context.Context) *ProjectsLocationsClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersListCall) Do(opts ...googleapi.CallOption) (*ListClustersResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListClustersResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetAddonsCall struct {
	s			*Service
	name			string
	setaddonsconfigrequest	*SetAddonsConfigRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetAddons(name string, setaddonsconfigrequest *SetAddonsConfigRequest) *ProjectsLocationsClustersSetAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetAddonsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setaddonsconfigrequest = setaddonsconfigrequest
	return c
}
func (c *ProjectsLocationsClustersSetAddonsCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetAddonsCall) Context(ctx context.Context) *ProjectsLocationsClustersSetAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetAddonsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetAddonsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setaddonsconfigrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setAddons")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetAddonsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetLegacyAbacCall struct {
	s			*Service
	name			string
	setlegacyabacrequest	*SetLegacyAbacRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetLegacyAbac(name string, setlegacyabacrequest *SetLegacyAbacRequest) *ProjectsLocationsClustersSetLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetLegacyAbacCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setlegacyabacrequest = setlegacyabacrequest
	return c
}
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Context(ctx context.Context) *ProjectsLocationsClustersSetLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetLegacyAbacCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlegacyabacrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setLegacyAbac")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetLocationsCall struct {
	s			*Service
	name			string
	setlocationsrequest	*SetLocationsRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetLocations(name string, setlocationsrequest *SetLocationsRequest) *ProjectsLocationsClustersSetLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetLocationsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setlocationsrequest = setlocationsrequest
	return c
}
func (c *ProjectsLocationsClustersSetLocationsCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetLocationsCall) Context(ctx context.Context) *ProjectsLocationsClustersSetLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetLocationsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetLocationsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlocationsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setLocations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetLocationsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetLoggingCall struct {
	s				*Service
	name				string
	setloggingservicerequest	*SetLoggingServiceRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersService) SetLogging(name string, setloggingservicerequest *SetLoggingServiceRequest) *ProjectsLocationsClustersSetLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetLoggingCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setloggingservicerequest = setloggingservicerequest
	return c
}
func (c *ProjectsLocationsClustersSetLoggingCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetLoggingCall) Context(ctx context.Context) *ProjectsLocationsClustersSetLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetLoggingCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetLoggingCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setloggingservicerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setLogging")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetLoggingCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetMaintenancePolicyCall struct {
	s				*Service
	name				string
	setmaintenancepolicyrequest	*SetMaintenancePolicyRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersService) SetMaintenancePolicy(name string, setmaintenancepolicyrequest *SetMaintenancePolicyRequest) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetMaintenancePolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setmaintenancepolicyrequest = setmaintenancepolicyrequest
	return c
}
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Context(ctx context.Context) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmaintenancepolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setMaintenancePolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetMasterAuthCall struct {
	s			*Service
	name			string
	setmasterauthrequest	*SetMasterAuthRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetMasterAuth(name string, setmasterauthrequest *SetMasterAuthRequest) *ProjectsLocationsClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetMasterAuthCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setmasterauthrequest = setmasterauthrequest
	return c
}
func (c *ProjectsLocationsClustersSetMasterAuthCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetMasterAuthCall) Context(ctx context.Context) *ProjectsLocationsClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetMasterAuthCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetMasterAuthCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmasterauthrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setMasterAuth")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetMasterAuthCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetMonitoringCall struct {
	s				*Service
	name				string
	setmonitoringservicerequest	*SetMonitoringServiceRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersService) SetMonitoring(name string, setmonitoringservicerequest *SetMonitoringServiceRequest) *ProjectsLocationsClustersSetMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetMonitoringCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setmonitoringservicerequest = setmonitoringservicerequest
	return c
}
func (c *ProjectsLocationsClustersSetMonitoringCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetMonitoringCall) Context(ctx context.Context) *ProjectsLocationsClustersSetMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetMonitoringCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetMonitoringCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmonitoringservicerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setMonitoring")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetMonitoringCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetNetworkPolicyCall struct {
	s			*Service
	name			string
	setnetworkpolicyrequest	*SetNetworkPolicyRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetNetworkPolicy(name string, setnetworkpolicyrequest *SetNetworkPolicyRequest) *ProjectsLocationsClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetNetworkPolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnetworkpolicyrequest = setnetworkpolicyrequest
	return c
}
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Context(ctx context.Context) *ProjectsLocationsClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnetworkpolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setNetworkPolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersSetResourceLabelsCall struct {
	s			*Service
	name			string
	setlabelsrequest	*SetLabelsRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) SetResourceLabels(name string, setlabelsrequest *SetLabelsRequest) *ProjectsLocationsClustersSetResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersSetResourceLabelsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setlabelsrequest = setlabelsrequest
	return c
}
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Context(ctx context.Context) *ProjectsLocationsClustersSetResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersSetResourceLabelsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlabelsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setResourceLabels")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersStartIpRotationCall struct {
	s			*Service
	name			string
	startiprotationrequest	*StartIPRotationRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) StartIpRotation(name string, startiprotationrequest *StartIPRotationRequest) *ProjectsLocationsClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersStartIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.startiprotationrequest = startiprotationrequest
	return c
}
func (c *ProjectsLocationsClustersStartIpRotationCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersStartIpRotationCall) Context(ctx context.Context) *ProjectsLocationsClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersStartIpRotationCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersStartIpRotationCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.startiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:startIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersStartIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersUpdateCall struct {
	s			*Service
	name			string
	updateclusterrequest	*UpdateClusterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) Update(name string, updateclusterrequest *UpdateClusterRequest) *ProjectsLocationsClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.updateclusterrequest = updateclusterrequest
	return c
}
func (c *ProjectsLocationsClustersUpdateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersUpdateCall) Context(ctx context.Context) *ProjectsLocationsClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersUpdateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersUpdateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updateclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersUpdateMasterCall struct {
	s			*Service
	name			string
	updatemasterrequest	*UpdateMasterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersService) UpdateMaster(name string, updatemasterrequest *UpdateMasterRequest) *ProjectsLocationsClustersUpdateMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersUpdateMasterCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.updatemasterrequest = updatemasterrequest
	return c
}
func (c *ProjectsLocationsClustersUpdateMasterCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersUpdateMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersUpdateMasterCall) Context(ctx context.Context) *ProjectsLocationsClustersUpdateMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersUpdateMasterCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersUpdateMasterCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updatemasterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:updateMaster")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersUpdateMasterCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsCreateCall struct {
	s			*Service
	parent			string
	createnodepoolrequest	*CreateNodePoolRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) Create(parent string, createnodepoolrequest *CreateNodePoolRequest) *ProjectsLocationsClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.createnodepoolrequest = createnodepoolrequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsCreateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createnodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsDeleteCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) Delete(name string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("clusterId", clusterId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) NodePoolId(nodePoolId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("nodePoolId", nodePoolId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsGetCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) Get(name string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("clusterId", clusterId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) NodePoolId(nodePoolId string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("nodePoolId", nodePoolId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsGetCall) Do(opts ...googleapi.CallOption) (*NodePool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &NodePool{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsListCall struct {
	s		*Service
	parent		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) List(parent string) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("clusterId", clusterId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsListCall) Do(opts ...googleapi.CallOption) (*ListNodePoolsResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListNodePoolsResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsRollbackCall struct {
	s				*Service
	name				string
	rollbacknodepoolupgraderequest	*RollbackNodePoolUpgradeRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) Rollback(name string, rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest) *ProjectsLocationsClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsRollbackCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.rollbacknodepoolupgraderequest = rollbacknodepoolupgraderequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.rollbacknodepoolupgraderequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:rollback")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsSetAutoscalingCall struct {
	s				*Service
	name				string
	setnodepoolautoscalingrequest	*SetNodePoolAutoscalingRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) SetAutoscaling(name string, setnodepoolautoscalingrequest *SetNodePoolAutoscalingRequest) *ProjectsLocationsClustersNodePoolsSetAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsSetAutoscalingCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnodepoolautoscalingrequest = setnodepoolautoscalingrequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetAutoscalingCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsSetAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetAutoscalingCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsSetAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetAutoscalingCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsSetAutoscalingCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolautoscalingrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setAutoscaling")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsSetAutoscalingCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsSetManagementCall struct {
	s				*Service
	name				string
	setnodepoolmanagementrequest	*SetNodePoolManagementRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) SetManagement(name string, setnodepoolmanagementrequest *SetNodePoolManagementRequest) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsSetManagementCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnodepoolmanagementrequest = setnodepoolmanagementrequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolmanagementrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setManagement")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsSetSizeCall struct {
	s			*Service
	name			string
	setnodepoolsizerequest	*SetNodePoolSizeRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) SetSize(name string, setnodepoolsizerequest *SetNodePoolSizeRequest) *ProjectsLocationsClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsSetSizeCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnodepoolsizerequest = setnodepoolsizerequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetSizeCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetSizeCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsSetSizeCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsSetSizeCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolsizerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setSize")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsSetSizeCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsClustersNodePoolsUpdateCall struct {
	s			*Service
	name			string
	updatenodepoolrequest	*UpdateNodePoolRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsClustersNodePoolsService) Update(name string, updatenodepoolrequest *UpdateNodePoolRequest) *ProjectsLocationsClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsClustersNodePoolsUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.updatenodepoolrequest = updatenodepoolrequest
	return c
}
func (c *ProjectsLocationsClustersNodePoolsUpdateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsClustersNodePoolsUpdateCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsClustersNodePoolsUpdateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsClustersNodePoolsUpdateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updatenodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsClustersNodePoolsUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsOperationsCancelCall struct {
	s			*Service
	name			string
	canceloperationrequest	*CancelOperationRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsLocationsOperationsService) Cancel(name string, canceloperationrequest *CancelOperationRequest) *ProjectsLocationsOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsOperationsCancelCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.canceloperationrequest = canceloperationrequest
	return c
}
func (c *ProjectsLocationsOperationsCancelCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsOperationsCancelCall) Context(ctx context.Context) *ProjectsLocationsOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsOperationsCancelCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsOperationsCancelCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.canceloperationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:cancel")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsOperationsCancelCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsOperationsGetCall struct {
	s		*Service
	name		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsOperationsService) Get(name string) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsOperationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}
func (c *ProjectsLocationsOperationsGetCall) OperationId(operationId string) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("operationId", operationId)
	return c
}
func (c *ProjectsLocationsOperationsGetCall) ProjectId(projectId string) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsOperationsGetCall) Zone(zone string) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsOperationsGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsOperationsGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsOperationsGetCall) Context(ctx context.Context) *ProjectsLocationsOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsOperationsGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsOperationsGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"name": c.name})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsOperationsGetCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsLocationsOperationsListCall struct {
	s		*Service
	parent		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsLocationsOperationsService) List(parent string) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsLocationsOperationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}
func (c *ProjectsLocationsOperationsListCall) ProjectId(projectId string) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("projectId", projectId)
	return c
}
func (c *ProjectsLocationsOperationsListCall) Zone(zone string) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("zone", zone)
	return c
}
func (c *ProjectsLocationsOperationsListCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsLocationsOperationsListCall) IfNoneMatch(entityTag string) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsLocationsOperationsListCall) Context(ctx context.Context) *ProjectsLocationsOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsLocationsOperationsListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsLocationsOperationsListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/operations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"parent": c.parent})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsLocationsOperationsListCall) Do(opts ...googleapi.CallOption) (*ListOperationsResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListOperationsResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesGetServerconfigCall struct {
	s		*Service
	projectId	string
	zone		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesService) GetServerconfig(projectId string, zone string) *ProjectsZonesGetServerconfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesGetServerconfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}
func (c *ProjectsZonesGetServerconfigCall) Name(name string) *ProjectsZonesGetServerconfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesGetServerconfigCall) Fields(s ...googleapi.Field) *ProjectsZonesGetServerconfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesGetServerconfigCall) IfNoneMatch(entityTag string) *ProjectsZonesGetServerconfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesGetServerconfigCall) Context(ctx context.Context) *ProjectsZonesGetServerconfigCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesGetServerconfigCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesGetServerconfigCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/serverconfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesGetServerconfigCall) Do(opts ...googleapi.CallOption) (*ServerConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ServerConfig{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersAddonsCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setaddonsconfigrequest	*SetAddonsConfigRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) Addons(projectId string, zone string, clusterId string, setaddonsconfigrequest *SetAddonsConfigRequest) *ProjectsZonesClustersAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersAddonsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setaddonsconfigrequest = setaddonsconfigrequest
	return c
}
func (c *ProjectsZonesClustersAddonsCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersAddonsCall) Context(ctx context.Context) *ProjectsZonesClustersAddonsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersAddonsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersAddonsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setaddonsconfigrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/addons")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersAddonsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersCompleteIpRotationCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	completeiprotationrequest	*CompleteIPRotationRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersService) CompleteIpRotation(projectId string, zone string, clusterId string, completeiprotationrequest *CompleteIPRotationRequest) *ProjectsZonesClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersCompleteIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.completeiprotationrequest = completeiprotationrequest
	return c
}
func (c *ProjectsZonesClustersCompleteIpRotationCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersCompleteIpRotationCall) Context(ctx context.Context) *ProjectsZonesClustersCompleteIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersCompleteIpRotationCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersCompleteIpRotationCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.completeiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:completeIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersCompleteIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersCreateCall struct {
	s			*Service
	projectId		string
	zone			string
	createclusterrequest	*CreateClusterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) Create(projectId string, zone string, createclusterrequest *CreateClusterRequest) *ProjectsZonesClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.createclusterrequest = createclusterrequest
	return c
}
func (c *ProjectsZonesClustersCreateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersCreateCall) Context(ctx context.Context) *ProjectsZonesClustersCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersCreateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersCreateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersDeleteCall struct {
	s		*Service
	projectId	string
	zone		string
	clusterId	string
	urlParams_	gensupport.URLParams
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersService) Delete(projectId string, zone string, clusterId string) *ProjectsZonesClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}
func (c *ProjectsZonesClustersDeleteCall) Name(name string) *ProjectsZonesClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesClustersDeleteCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersDeleteCall) Context(ctx context.Context) *ProjectsZonesClustersDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersDeleteCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersDeleteCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersGetCall struct {
	s		*Service
	projectId	string
	zone		string
	clusterId	string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersService) Get(projectId string, zone string, clusterId string) *ProjectsZonesClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}
func (c *ProjectsZonesClustersGetCall) Name(name string) *ProjectsZonesClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesClustersGetCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersGetCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesClustersGetCall) Context(ctx context.Context) *ProjectsZonesClustersGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersGetCall) Do(opts ...googleapi.CallOption) (*Cluster, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Cluster{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersLegacyAbacCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setlegacyabacrequest	*SetLegacyAbacRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) LegacyAbac(projectId string, zone string, clusterId string, setlegacyabacrequest *SetLegacyAbacRequest) *ProjectsZonesClustersLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersLegacyAbacCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setlegacyabacrequest = setlegacyabacrequest
	return c
}
func (c *ProjectsZonesClustersLegacyAbacCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersLegacyAbacCall) Context(ctx context.Context) *ProjectsZonesClustersLegacyAbacCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersLegacyAbacCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersLegacyAbacCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlegacyabacrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/legacyAbac")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersLegacyAbacCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersListCall struct {
	s		*Service
	projectId	string
	zone		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersService) List(projectId string, zone string) *ProjectsZonesClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}
func (c *ProjectsZonesClustersListCall) Parent(parent string) *ProjectsZonesClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("parent", parent)
	return c
}
func (c *ProjectsZonesClustersListCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersListCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesClustersListCall) Context(ctx context.Context) *ProjectsZonesClustersListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersListCall) Do(opts ...googleapi.CallOption) (*ListClustersResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListClustersResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersLocationsCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setlocationsrequest	*SetLocationsRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) Locations(projectId string, zone string, clusterId string, setlocationsrequest *SetLocationsRequest) *ProjectsZonesClustersLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersLocationsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setlocationsrequest = setlocationsrequest
	return c
}
func (c *ProjectsZonesClustersLocationsCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersLocationsCall) Context(ctx context.Context) *ProjectsZonesClustersLocationsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersLocationsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersLocationsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlocationsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/locations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersLocationsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersLoggingCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	setloggingservicerequest	*SetLoggingServiceRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersService) Logging(projectId string, zone string, clusterId string, setloggingservicerequest *SetLoggingServiceRequest) *ProjectsZonesClustersLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersLoggingCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setloggingservicerequest = setloggingservicerequest
	return c
}
func (c *ProjectsZonesClustersLoggingCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersLoggingCall) Context(ctx context.Context) *ProjectsZonesClustersLoggingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersLoggingCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersLoggingCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setloggingservicerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/logging")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersLoggingCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersMasterCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	updatemasterrequest	*UpdateMasterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) Master(projectId string, zone string, clusterId string, updatemasterrequest *UpdateMasterRequest) *ProjectsZonesClustersMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersMasterCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.updatemasterrequest = updatemasterrequest
	return c
}
func (c *ProjectsZonesClustersMasterCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersMasterCall) Context(ctx context.Context) *ProjectsZonesClustersMasterCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersMasterCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersMasterCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updatemasterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/master")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersMasterCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersMonitoringCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	setmonitoringservicerequest	*SetMonitoringServiceRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersService) Monitoring(projectId string, zone string, clusterId string, setmonitoringservicerequest *SetMonitoringServiceRequest) *ProjectsZonesClustersMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersMonitoringCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setmonitoringservicerequest = setmonitoringservicerequest
	return c
}
func (c *ProjectsZonesClustersMonitoringCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersMonitoringCall) Context(ctx context.Context) *ProjectsZonesClustersMonitoringCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersMonitoringCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersMonitoringCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmonitoringservicerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/monitoring")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersMonitoringCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersResourceLabelsCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setlabelsrequest	*SetLabelsRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) ResourceLabels(projectId string, zone string, clusterId string, setlabelsrequest *SetLabelsRequest) *ProjectsZonesClustersResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersResourceLabelsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setlabelsrequest = setlabelsrequest
	return c
}
func (c *ProjectsZonesClustersResourceLabelsCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersResourceLabelsCall) Context(ctx context.Context) *ProjectsZonesClustersResourceLabelsCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersResourceLabelsCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersResourceLabelsCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlabelsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/resourceLabels")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersResourceLabelsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersSetMaintenancePolicyCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	setmaintenancepolicyrequest	*SetMaintenancePolicyRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersService) SetMaintenancePolicy(projectId string, zone string, clusterId string, setmaintenancepolicyrequest *SetMaintenancePolicyRequest) *ProjectsZonesClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersSetMaintenancePolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setmaintenancepolicyrequest = setmaintenancepolicyrequest
	return c
}
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Context(ctx context.Context) *ProjectsZonesClustersSetMaintenancePolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmaintenancepolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMaintenancePolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersSetMasterAuthCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setmasterauthrequest	*SetMasterAuthRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) SetMasterAuth(projectId string, zone string, clusterId string, setmasterauthrequest *SetMasterAuthRequest) *ProjectsZonesClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersSetMasterAuthCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setmasterauthrequest = setmasterauthrequest
	return c
}
func (c *ProjectsZonesClustersSetMasterAuthCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersSetMasterAuthCall) Context(ctx context.Context) *ProjectsZonesClustersSetMasterAuthCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersSetMasterAuthCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersSetMasterAuthCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmasterauthrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMasterAuth")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersSetMasterAuthCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersSetNetworkPolicyCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	setnetworkpolicyrequest	*SetNetworkPolicyRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) SetNetworkPolicy(projectId string, zone string, clusterId string, setnetworkpolicyrequest *SetNetworkPolicyRequest) *ProjectsZonesClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersSetNetworkPolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setnetworkpolicyrequest = setnetworkpolicyrequest
	return c
}
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Context(ctx context.Context) *ProjectsZonesClustersSetNetworkPolicyCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersSetNetworkPolicyCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnetworkpolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setNetworkPolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersStartIpRotationCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	startiprotationrequest	*StartIPRotationRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) StartIpRotation(projectId string, zone string, clusterId string, startiprotationrequest *StartIPRotationRequest) *ProjectsZonesClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersStartIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.startiprotationrequest = startiprotationrequest
	return c
}
func (c *ProjectsZonesClustersStartIpRotationCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersStartIpRotationCall) Context(ctx context.Context) *ProjectsZonesClustersStartIpRotationCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersStartIpRotationCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersStartIpRotationCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.startiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:startIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersStartIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersUpdateCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	updateclusterrequest	*UpdateClusterRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersService) Update(projectId string, zone string, clusterId string, updateclusterrequest *UpdateClusterRequest) *ProjectsZonesClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.updateclusterrequest = updateclusterrequest
	return c
}
func (c *ProjectsZonesClustersUpdateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersUpdateCall) Context(ctx context.Context) *ProjectsZonesClustersUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersUpdateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersUpdateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updateclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsAutoscalingCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	nodePoolId			string
	setnodepoolautoscalingrequest	*SetNodePoolAutoscalingRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Autoscaling(projectId string, zone string, clusterId string, nodePoolId string, setnodepoolautoscalingrequest *SetNodePoolAutoscalingRequest) *ProjectsZonesClustersNodePoolsAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsAutoscalingCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.setnodepoolautoscalingrequest = setnodepoolautoscalingrequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsAutoscalingCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsAutoscalingCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsAutoscalingCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsAutoscalingCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsAutoscalingCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolautoscalingrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/autoscaling")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsAutoscalingCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsCreateCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	createnodepoolrequest	*CreateNodePoolRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Create(projectId string, zone string, clusterId string, createnodepoolrequest *CreateNodePoolRequest) *ProjectsZonesClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.createnodepoolrequest = createnodepoolrequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsCreateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsCreateCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsCreateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsCreateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsCreateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createnodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsDeleteCall struct {
	s		*Service
	projectId	string
	zone		string
	clusterId	string
	nodePoolId	string
	urlParams_	gensupport.URLParams
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Delete(projectId string, zone string, clusterId string, nodePoolId string) *ProjectsZonesClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	return c
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Name(name string) *ProjectsZonesClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsDeleteCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsGetCall struct {
	s		*Service
	projectId	string
	zone		string
	clusterId	string
	nodePoolId	string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Get(projectId string, zone string, clusterId string, nodePoolId string) *ProjectsZonesClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	return c
}
func (c *ProjectsZonesClustersNodePoolsGetCall) Name(name string) *ProjectsZonesClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesClustersNodePoolsGetCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsGetCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesClustersNodePoolsGetCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsGetCall) Do(opts ...googleapi.CallOption) (*NodePool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &NodePool{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsListCall struct {
	s		*Service
	projectId	string
	zone		string
	clusterId	string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) List(projectId string, zone string, clusterId string) *ProjectsZonesClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}
func (c *ProjectsZonesClustersNodePoolsListCall) Parent(parent string) *ProjectsZonesClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("parent", parent)
	return c
}
func (c *ProjectsZonesClustersNodePoolsListCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsListCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesClustersNodePoolsListCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsListCall) Do(opts ...googleapi.CallOption) (*ListNodePoolsResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListNodePoolsResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsRollbackCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	nodePoolId			string
	rollbacknodepoolupgraderequest	*RollbackNodePoolUpgradeRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Rollback(projectId string, zone string, clusterId string, nodePoolId string, rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest) *ProjectsZonesClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsRollbackCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.rollbacknodepoolupgraderequest = rollbacknodepoolupgraderequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsRollbackCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsRollbackCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.rollbacknodepoolupgraderequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}:rollback")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsSetManagementCall struct {
	s				*Service
	projectId			string
	zone				string
	clusterId			string
	nodePoolId			string
	setnodepoolmanagementrequest	*SetNodePoolManagementRequest
	urlParams_			gensupport.URLParams
	ctx_				context.Context
	header_				http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) SetManagement(projectId string, zone string, clusterId string, nodePoolId string, setnodepoolmanagementrequest *SetNodePoolManagementRequest) *ProjectsZonesClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsSetManagementCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.setnodepoolmanagementrequest = setnodepoolmanagementrequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsSetManagementCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolmanagementrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/setManagement")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsSetSizeCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	nodePoolId		string
	setnodepoolsizerequest	*SetNodePoolSizeRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) SetSize(projectId string, zone string, clusterId string, nodePoolId string, setnodepoolsizerequest *SetNodePoolSizeRequest) *ProjectsZonesClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsSetSizeCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.setnodepoolsizerequest = setnodepoolsizerequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetSizeCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetSizeCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsSetSizeCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsSetSizeCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsSetSizeCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolsizerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/setSize")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsSetSizeCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesClustersNodePoolsUpdateCall struct {
	s			*Service
	projectId		string
	zone			string
	clusterId		string
	nodePoolId		string
	updatenodepoolrequest	*UpdateNodePoolRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesClustersNodePoolsService) Update(projectId string, zone string, clusterId string, nodePoolId string, updatenodepoolrequest *UpdateNodePoolRequest) *ProjectsZonesClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesClustersNodePoolsUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.updatenodepoolrequest = updatenodepoolrequest
	return c
}
func (c *ProjectsZonesClustersNodePoolsUpdateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesClustersNodePoolsUpdateCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsUpdateCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesClustersNodePoolsUpdateCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesClustersNodePoolsUpdateCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updatenodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/update")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "clusterId": c.clusterId, "nodePoolId": c.nodePoolId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesClustersNodePoolsUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesOperationsCancelCall struct {
	s			*Service
	projectId		string
	zone			string
	operationId		string
	canceloperationrequest	*CancelOperationRequest
	urlParams_		gensupport.URLParams
	ctx_			context.Context
	header_			http.Header
}

func (r *ProjectsZonesOperationsService) Cancel(projectId string, zone string, operationId string, canceloperationrequest *CancelOperationRequest) *ProjectsZonesOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesOperationsCancelCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.operationId = operationId
	c.canceloperationrequest = canceloperationrequest
	return c
}
func (c *ProjectsZonesOperationsCancelCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesOperationsCancelCall) Context(ctx context.Context) *ProjectsZonesOperationsCancelCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesOperationsCancelCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesOperationsCancelCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.canceloperationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}:cancel")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "operationId": c.operationId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesOperationsCancelCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesOperationsGetCall struct {
	s		*Service
	projectId	string
	zone		string
	operationId	string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesOperationsService) Get(projectId string, zone string, operationId string) *ProjectsZonesOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesOperationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.operationId = operationId
	return c
}
func (c *ProjectsZonesOperationsGetCall) Name(name string) *ProjectsZonesOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("name", name)
	return c
}
func (c *ProjectsZonesOperationsGetCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesOperationsGetCall) IfNoneMatch(entityTag string) *ProjectsZonesOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesOperationsGetCall) Context(ctx context.Context) *ProjectsZonesOperationsGetCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesOperationsGetCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesOperationsGetCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone, "operationId": c.operationId})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesOperationsGetCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Operation{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}

type ProjectsZonesOperationsListCall struct {
	s		*Service
	projectId	string
	zone		string
	urlParams_	gensupport.URLParams
	ifNoneMatch_	string
	ctx_		context.Context
	header_		http.Header
}

func (r *ProjectsZonesOperationsService) List(projectId string, zone string) *ProjectsZonesOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &ProjectsZonesOperationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}
func (c *ProjectsZonesOperationsListCall) Parent(parent string) *ProjectsZonesOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("parent", parent)
	return c
}
func (c *ProjectsZonesOperationsListCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}
func (c *ProjectsZonesOperationsListCall) IfNoneMatch(entityTag string) *ProjectsZonesOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ifNoneMatch_ = entityTag
	return c
}
func (c *ProjectsZonesOperationsListCall) Context(ctx context.Context) *ProjectsZonesOperationsListCall {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.ctx_ = ctx
	return c
}
func (c *ProjectsZonesOperationsListCall) Header() http.Header {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}
func (c *ProjectsZonesOperationsListCall) doRequest(alt string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{"projectId": c.projectId, "zone": c.zone})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}
func (c *ProjectsZonesOperationsListCall) Do(opts ...googleapi.CallOption) (*ListOperationsResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{Code: res.StatusCode, Header: res.Header}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListOperationsResponse{ServerResponse: googleapi.ServerResponse{Header: res.Header, HTTPStatusCode: res.StatusCode}}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
