package status

import (
	"context"
	v1 "lfm-operator/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Option interface {
	ApplyOption(vg *v1.VolumeGroup)
	GetResult() (reconcile.Result, error)
}

type OptionBuilder interface {
	GetOptions() []Option
}

// Update takes the options provided by the given option builder, applies them all and then updates the resource
func Update(statusWriter client.StatusWriter, vg *v1.VolumeGroup, optionBuilder OptionBuilder) (reconcile.Result, error) {
	options := optionBuilder.GetOptions()
	for _, opt := range options {
		opt.ApplyOption(vg)
	}

	if err := statusWriter.Update(context.TODO(), vg); err != nil {
		return reconcile.Result{}, err
	}

	return determineReconciliationResult(options)
}

func determineReconciliationResult(options []Option) (reconcile.Result, error) {
	for _, opt := range options {
		res, err := opt.GetResult()
		if err != nil {
			return res, err
		}
	}
	for _, opt := range options {
		res, _ := opt.GetResult()
		if res.Requeue || res.RequeueAfter > 0 {
			return res, nil
		}
	}
	return reconcile.Result{}, nil
}
