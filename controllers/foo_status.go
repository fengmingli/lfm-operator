package controllers

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "lfm-operator/api/v1"
	"lfm-operator/pkg/util/result"
	"lfm-operator/pkg/util/status"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// severity indicates the severity level
// at which the message should be logged

// optionBuilder is in charge of constructing a slice of options that
// will be applied on top of the MongoDB resource that has been provided
type optionBuilder struct {
	options []status.Option
}

// GetOptions implements the OptionBuilder interface
func (o *optionBuilder) GetOptions() []status.Option {
	return o.options
}

// statusOptions returns an initialized optionBuilder
func statusOptions() *optionBuilder {
	return &optionBuilder{
		options: []status.Option{},
	}
}

type LvCreateStatus int32

const (
	MkdirDirectory  LvCreateStatus = 1
	LvCreate        LvCreateStatus = 2
	Ext42FileSystem LvCreateStatus = 3
	MountPath       LvCreateStatus = 4
)

type statusOption struct {
	lvCreateStatus LvCreateStatus
}

func (p statusOption) ApplyOption(mdb *v1.VolumeGroup) {
	//mdb.Status.PvcStatus = p.PvcStatus
}

func (p statusOption) GetResult() (reconcile.Result, error) {
	return result.OK()
}

func (o *optionBuilder) withPvcItem(pvName string, size resource.Quantity) *optionBuilder {
	o.options = append(o.options,
		statusOption{})
	return o
}
