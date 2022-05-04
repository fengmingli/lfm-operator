package controllers

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"lfm-operator/pkg/kube/client"
	"lfm-operator/pkg/util/result"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PersistentVolumeClaimReconciler reconciles a PersistentVolumeClaim object
type PersistentVolumeClaimReconciler struct {
	client client.Client
}

func NewPersistentVolumeClaim(mgr manager.Manager) *PersistentVolumeClaimReconciler {
	return &PersistentVolumeClaimReconciler{
		client: client.NewClient(mgr.GetClient()),
	}
}

//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;delete

// Reconcile finalize PVC
func (r *PersistentVolumeClaimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof("=======pvc==========%v", req)
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(ctx, req.NamespacedName, pvc)
	switch {
	case err == nil:
	case apierrors.IsNotFound(err):
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, err
	}

	if pvc.DeletionTimestamp == nil {
		return ctrl.Result{}, nil
	}

	//确保是需要进行创建的pvc资源
	isNeedPvcResource := r.ensureNeedPvcResource(pvc)
	if !isNeedPvcResource {
		return result.OK()
	}

	//校验sc是否需要lvm进行提供
	isLvmProvisioner, err := r.validateStorageClassConfig(pvc)
	if err != nil {
		klog.Error("Get sc fail errMsg:%s", err)
		return result.OK()
	}
	if !isLvmProvisioner {
		return result.OK()
	}

	r.ensureLvmResource(pvc)

	r.ensurePvResource(pvc)

	return ctrl.Result{}, nil
}

//validateLvmConfig 校验lvm相关配置
func (r *PersistentVolumeClaimReconciler) validateStorageClassConfig(pvc *corev1.PersistentVolumeClaim) (bool, error) {

	className := pvc.Spec.StorageClassName
	if *className == "" {
		return false, fmt.Errorf("namesapce:%s ,pvc:%s storageClass empty", pvc.Namespace, pvc.Name)
	}
	//获取sc，判断sc是的provision是否是lvm-simple

	//判断params的vgType参数是否为空
	return true, nil
}

//validateLvmConfig 校验lvm相关配置
func (r *PersistentVolumeClaimReconciler) ensureNeedPvcResource(pvc *corev1.PersistentVolumeClaim) bool {
	//只处理pvc为pending状态的资源
	if pvc.Status.Phase != corev1.ClaimPending {
		return false
	}

	//确保annotation存在调度器已经选择可调度节点
	pvcAnnotations := pvc.GetAnnotations()
	if pvcAnnotations == nil || len(pvcAnnotations) == 0 {
		return false
	}

	for key := range pvcAnnotations {
		if key == AnnSelectedNode {
			return true
		}
	}
	return false
}

// EnsureLvmResource 确保lvm资源
func (r *PersistentVolumeClaimReconciler) ensureLvmResource(pvc *corev1.PersistentVolumeClaim) error {
	//check ds资源是否存在
	//判断文件目录是否存在
	//判断lv是否创建
	//判断是否已经格式化成文件系统
	//判断是否已经mount
	//判断是否已经写入到fstab文件
	return nil
}

// EnsurePvResource  确保lvm资源
func (r *PersistentVolumeClaimReconciler) ensurePvResource(pvc *corev1.PersistentVolumeClaim) error {
	//确保创建pv不存在
	//创建pv
	//确保pv的状态是bound状态
	return nil
}

// SetupWithManager sets up the controllers with the Manager.
func (r *PersistentVolumeClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.Funcs{
		CreateFunc:  func(event.CreateEvent) bool { return false },
		DeleteFunc:  func(event.DeleteEvent) bool { return true },
		UpdateFunc:  func(event.UpdateEvent) bool { return true },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(pred).
		For(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
