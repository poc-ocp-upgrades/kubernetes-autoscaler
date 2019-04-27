package containerservice

import "github.com/Azure/azure-sdk-for-go/version"

func UserAgent() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "Azure-SDK-For-Go/" + version.Number + " containerservice/2018-03-31"
}
func Version() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return version.Number
}
