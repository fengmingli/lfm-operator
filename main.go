package main

/**
 * @Author: LFM
 * @Date: 2022/5/3 4:52 下午
 * @Since: 1.0.0
 * @Desc: TODO
 */

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	v1 "lfm-operator/api/v1"
	"lfm-operator/controllers"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
}

func main() {
	//1.实例化 manager，参数 config
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
		// metrics export
		MetricsBindAddress: ":8888",
		// health 检测
		HealthProbeBindAddress: ":8899",
	})

	if err != nil {
		klog.Errorf("Unable to set up controllers manager: %v", err)
		os.Exit(1)
	}

	klog.Infof("Registering Components.")

	//2。注册 Controller 到 Manager
	if err1 := controllers.NewVgReconciler(mgr).SetupWithManager(mgr); err1 != nil {
		klog.Errorf("Unable to create controllers: %v", err1)
		os.Exit(1)
	}

	if err2 := controllers.NewPersistentVolumeReconciler(mgr).SetupWithManager(mgr); err2 != nil {
		klog.Errorf("Unable to create pv controllers: %v", err2)
		os.Exit(1)
	}

	if err2 := controllers.NewPersistentVolumeClaim(mgr).SetupWithManager(mgr); err2 != nil {
		klog.Errorf("Unable to create pvc controllers: %v", err2)
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Errorf("unable to set up health check", err)
		os.Exit(1)
	}

	klog.Infof("Starting the Cmd.")

	//3。 启动一个manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Errorf("Unable to start manager: %v", err)
		os.Exit(1)
	}
}
