package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"lfm-operator/pkg/kube/client"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PersistentVolumeReconciler reconciles a PersistentVolumeClaim object
type PersistentVolumeReconciler struct {
	client client.Client
}

func NewPersistentVolumeReconciler(mgr manager.Manager) *PersistentVolumeReconciler {
	return &PersistentVolumeReconciler{
		client: client.NewClient(mgr.GetClient()),
	}
}

//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;delete

// Reconcile finalize PVC
func (r *PersistentVolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof("=======pv==========%v", req)
	pvc := &corev1.PersistentVolume{}
	err := r.client.Get(ctx, req.NamespacedName, pvc)
	switch {
	case err == nil:
	case apierrors.IsNotFound(err):
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, err
	}
	//delete event
	if pvc.DeletionTimestamp == nil {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controllers with the Manager.
func (r *PersistentVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.Funcs{
		CreateFunc:  func(event.CreateEvent) bool { return false },
		DeleteFunc:  func(event.DeleteEvent) bool { return true },
		UpdateFunc:  func(event.UpdateEvent) bool { return false },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(pred).
		For(&corev1.PersistentVolume{}).
		Complete(r)
}
