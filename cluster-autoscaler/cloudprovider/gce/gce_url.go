package gce

import (
	"fmt"
	"strings"
)

const (
	gceUrlSchema		= "https"
	gceDomainSuffix		= "googleapis.com/compute/v1/projects/"
	gcePrefix		= gceUrlSchema + "://content." + gceDomainSuffix
	instanceUrlTemplate	= gcePrefix + "%s/zones/%s/instances/%s"
	migUrlTemplate		= gcePrefix + "%s/zones/%s/instanceGroups/%s"
)

func ParseMigUrl(url string) (project string, zone string, name string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return parseGceUrl(url, "instanceGroups")
}
func ParseIgmUrl(url string) (project string, zone string, name string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return parseGceUrl(url, "instanceGroupManagers")
}
func ParseInstanceUrl(url string) (project string, zone string, name string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return parseGceUrl(url, "instances")
}
func ParseInstanceUrlRef(url string) (GceRef, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	project, zone, name, err := parseGceUrl(url, "instances")
	if err != nil {
		return GceRef{}, err
	}
	return GceRef{Project: project, Zone: zone, Name: name}, nil
}
func GenerateInstanceUrl(ref GceRef) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf(instanceUrlTemplate, ref.Project, ref.Zone, ref.Name)
}
func GenerateMigUrl(ref GceRef) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf(migUrlTemplate, ref.Project, ref.Zone, ref.Name)
}
func parseGceUrl(url, expectedResource string) (project string, zone string, name string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	errMsg := fmt.Errorf("Wrong url: expected format https://content.googleapis.com/compute/v1/projects/<project-id>/zones/<zone>/%s/<name>, got %s", expectedResource, url)
	if !strings.Contains(url, gceDomainSuffix) {
		return "", "", "", errMsg
	}
	if !strings.HasPrefix(url, gceUrlSchema) {
		return "", "", "", errMsg
	}
	splitted := strings.Split(strings.Split(url, gceDomainSuffix)[1], "/")
	if len(splitted) != 5 || splitted[1] != "zones" {
		return "", "", "", errMsg
	}
	if splitted[3] != expectedResource {
		return "", "", "", fmt.Errorf("Wrong resource in url: expected %s, got %s", expectedResource, splitted[3])
	}
	project = splitted[0]
	zone = splitted[2]
	name = splitted[4]
	return project, zone, name, nil
}
