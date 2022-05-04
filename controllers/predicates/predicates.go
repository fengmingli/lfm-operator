package predicates

import (
	v1 "lfm-operator/api/v1"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func OnlyOnSpecChange() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldResource := e.ObjectOld.(*v1.VolumeGroup)
			newResource := e.ObjectNew.(*v1.VolumeGroup)
			specChanged := !reflect.DeepEqual(oldResource.Spec, newResource.Spec)
			return specChanged
		},
	}
}
