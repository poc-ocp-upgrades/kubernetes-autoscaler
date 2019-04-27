package gce

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	proj, zone, name, err := parseGceUrl("https://www.googleapis.com/compute/v1/projects/mwielgus-proj/zones/us-central1-b/instanceGroups/kubernetes-minion-group", "instanceGroups")
	assert.Nil(t, err)
	assert.Equal(t, "mwielgus-proj", proj)
	assert.Equal(t, "us-central1-b", zone)
	assert.Equal(t, "kubernetes-minion-group", name)
	proj, zone, name, err = parseGceUrl("https://content.googleapis.com/compute/v1/projects/mwielgus-proj/zones/us-central1-b/instanceGroups/kubernetes-minion-group", "instanceGroups")
	assert.Nil(t, err)
	assert.Equal(t, "mwielgus-proj", proj)
	assert.Equal(t, "us-central1-b", zone)
	assert.Equal(t, "kubernetes-minion-group", name)
	proj, zone, name, err = parseGceUrl("www.onet.pl", "instanceGroups")
	assert.NotNil(t, err)
	proj, zone, name, err = parseGceUrl("https://content.googleapis.com/compute/vabc/projects/mwielgus-proj/zones/us-central1-b/instanceGroups/kubernetes-minion-group", "instanceGroups")
	assert.NotNil(t, err)
}
