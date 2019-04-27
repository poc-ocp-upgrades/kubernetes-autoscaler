package api

import (
	"fmt"
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
)

func prepareConditions() (health, scaleUp ClusterAutoscalerCondition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	healthCondition := ClusterAutoscalerCondition{Type: ClusterAutoscalerHealth, Status: ClusterAutoscalerHealthy, Message: "HEALTH_MESSAGE"}
	scaleUpCondition := ClusterAutoscalerCondition{Type: ClusterAutoscalerScaleUp, Status: ClusterAutoscalerNotNeeded, Message: "SCALE_UP_MESSAGE"}
	return healthCondition, scaleUpCondition
}
func TestGetStringForEmptyStatus(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var empty ClusterAutoscalerStatus
	assert.Regexp(t, regexp.MustCompile("\\s*Health:\\s*<unknown>"), empty.GetReadableString())
}
func TestGetStringNothingGoingOn(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var status ClusterAutoscalerStatus
	healthCondition, scaleUpCondition := prepareConditions()
	status.ClusterwideConditions = append(status.ClusterwideConditions, healthCondition)
	status.ClusterwideConditions = append(status.ClusterwideConditions, scaleUpCondition)
	result := status.GetReadableString()
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v:\\s*%v", ClusterAutoscalerHealth, ClusterAutoscalerHealthy)), result)
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v.*HEALTH_MESSAGE", ClusterAutoscalerHealth)), result)
	assert.NotRegexp(t, regexp.MustCompile(fmt.Sprintf("%v.*SCALE_UP_MESSAGE", ClusterAutoscalerHealth)), result)
	assert.NotRegexp(t, regexp.MustCompile("NodeGroups"), result)
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v:\\s*%v", ClusterAutoscalerScaleUp, ClusterAutoscalerNotNeeded)), result)
	var reorderedStatus ClusterAutoscalerStatus
	reorderedStatus.ClusterwideConditions = append(status.ClusterwideConditions, scaleUpCondition)
	reorderedStatus.ClusterwideConditions = append(status.ClusterwideConditions, healthCondition)
	reorderedResult := reorderedStatus.GetReadableString()
	assert.Equal(t, result, reorderedResult)
}
func TestGetStringScalingUp(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var status ClusterAutoscalerStatus
	healthCondition, scaleUpCondition := prepareConditions()
	scaleUpCondition.Status = ClusterAutoscalerInProgress
	status.ClusterwideConditions = append(status.ClusterwideConditions, healthCondition)
	status.ClusterwideConditions = append(status.ClusterwideConditions, scaleUpCondition)
	result := status.GetReadableString()
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v:\\s*%v.*SCALE_UP_MESSAGE", ClusterAutoscalerScaleUp, ClusterAutoscalerInProgress)), result)
}
func TestGetStringNodeGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var status ClusterAutoscalerStatus
	healthCondition, scaleUpCondition := prepareConditions()
	status.ClusterwideConditions = append(status.ClusterwideConditions, healthCondition)
	status.ClusterwideConditions = append(status.ClusterwideConditions, scaleUpCondition)
	var ng1, ng2 NodeGroupStatus
	ng1.ProviderID = "ng1"
	ng1.Conditions = status.ClusterwideConditions
	ng2.ProviderID = "ng2"
	ng2.Conditions = status.ClusterwideConditions
	status.NodeGroupStatuses = append(status.NodeGroupStatuses, ng1)
	status.NodeGroupStatuses = append(status.NodeGroupStatuses, ng2)
	result := status.GetReadableString()
	assert.Regexp(t, regexp.MustCompile("(?ms)NodeGroups:.*Name:\\s*ng1"), result)
	assert.Regexp(t, regexp.MustCompile("(?ms)NodeGroups:.*Name:\\s*ng2"), result)
}
