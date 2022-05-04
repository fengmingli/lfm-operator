package controllers

import (
	"context"
	"crypto/sha256"
	"fmt"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	apiv1 "lfm-operator/api/v1"
	"lfm-operator/controllers/predicates"
	"lfm-operator/controllers/watch"
	"lfm-operator/pkg/kube/client"
	"lfm-operator/pkg/kube/configmap"
	"lfm-operator/pkg/kube/secret"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

/**
 * @Author: LFM
 * @Date: 2022/5/3 4:59 下午
 * @Since: 1.0.0
 * @Desc: TODO
 */

func NewVgReconciler(mgr manager.Manager) *VgReconciler {
	secretWatcher := watch.New()
	return &VgReconciler{
		client:        client.NewClient(mgr.GetClient()),
		scheme:        mgr.GetScheme(),
		secretWatcher: &secretWatcher,
	}
}

type VgReconciler struct {
	client        client.Client
	scheme        *runtime.Scheme
	secretWatcher *watch.ResourceWatcher
}

func (r *VgReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	klog.Infof("=======vg==========%v", req)
	vg := apiv1.VolumeGroup{}

	if err := r.client.Get(ctx, req.NamespacedName, &vg); err != nil {
		if !apierrs.IsNotFound(err) {
			klog.Warningf("=======>%v<<<<<<<<<", vg)
			return ctrl.Result{}, err
		}
	}
	if vg.Name == "vg-1" {
	}

	if vg.Name == "vg-2" {
		ensureCASecret(r.client, r.client, r.client, vg)
	}
	klog.V(0).Infof(">>>>>>>>>>%v<<<<<<<<<", vg)
	return reconcile.Result{}, nil
}

func (r *VgReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		For(&apiv1.VolumeGroup{}, builder.WithPredicates(predicates.OnlyOnSpecChange())).
		Complete(r)
}

// ensureCASecret will create or update the operator managed Secret containing
// the CA certficate from the user provided Secret or ConfigMap.
func ensureCASecret(cmGetter configmap.Getter, secretGetter secret.Getter, getUpdateCreator secret.GetUpdateCreator, mdb apiv1.VolumeGroup) error {

	caFileName := tlsOperatorSecretFileName("lfm")

	operatorSecret := secret.Builder().
		SetName(mdb.Name).
		SetNamespace(mdb.Namespace).
		SetField(caFileName, "lfm").
		SetOwnerReferences(mdb.GetOwnerReferences()).
		Build()

	return secret.CreateOrUpdate(getUpdateCreator, operatorSecret)
}

func tlsOperatorSecretFileName(certKey string) string {
	hash := sha256.Sum256([]byte(certKey))
	return fmt.Sprintf("%x.pem", hash)
}
