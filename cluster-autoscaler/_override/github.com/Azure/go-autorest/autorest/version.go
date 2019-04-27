package autorest

import "github.com/Azure/go-autorest/version"

func Version() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return version.Number
}
