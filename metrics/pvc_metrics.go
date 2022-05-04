package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"time"
)

/**
 * @Author: LFM
 * @Date: 2022/5/4 10:09 下午
 * @Since: 1.0.0
 * @Desc: TODO
 */

// DeviceMetrics Device Metrics
type DeviceMetrics struct {
	nodeName    string
	FreeBytes   uint64
	TotalBytes  uint64
	DeviceGroup string
}

// VolumeMetrics volume Metrics
type VolumeMetrics struct {
	nodeName   string
	Volume     string
	TotalBytes uint64
	UsedBytes  float64
}

type metricsExporter struct {
	vgFreeBytes      *prometheus.GaugeVec
	vgTotalBytes     *prometheus.GaugeVec
	volumeTotalBytes *prometheus.GaugeVec
	volumeUsedBytes  *prometheus.GaugeVec
}

func (m metricsExporter) NeedLeaderElection() bool {
	return false
}

var _ manager.LeaderElectionRunnable = &metricsExporter{}

func NewMetricsExporter() manager.Runnable {
	klog.Errorf("===>>new metrics<<==")
	vgFreeBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "lvmsimple",
		Subsystem:   "devicegroup",
		Name:        "vg_free_bytes",
		Help:        "LVM VG free bytes",
		ConstLabels: prometheus.Labels{},
	}, []string{"node", "device_group"})

	vgTotalBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "lvmsimple",
		Subsystem:   "devicegroup",
		Name:        "vg_total_bytes",
		Help:        "LVM VG total bytes",
		ConstLabels: prometheus.Labels{},
	}, []string{"node", "device_group"})

	volumeTotalBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "lvmsimple",
		Subsystem:   "volume",
		Name:        "volume_total_bytes",
		Help:        "LVM Volume total bytes",
		ConstLabels: prometheus.Labels{},
	}, []string{"node", "volume"})

	volumeUsedBytes := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "lvmsimple",
		Subsystem:   "volume",
		Name:        "volume_used_bytes",
		Help:        "LVM volume used bytes",
		ConstLabels: prometheus.Labels{},
	}, []string{"node", "volume"})

	metrics.Registry.MustRegister(vgTotalBytes)
	metrics.Registry.MustRegister(vgFreeBytes)
	metrics.Registry.MustRegister(volumeTotalBytes)
	metrics.Registry.MustRegister(volumeUsedBytes)

	return &metricsExporter{
		vgFreeBytes:      vgFreeBytes,
		vgTotalBytes:     vgTotalBytes,
		volumeTotalBytes: volumeTotalBytes,
		volumeUsedBytes:  volumeUsedBytes,
	}
}

func (m *metricsExporter) Start(ctx context.Context) error {
	klog.Errorf("===>>start metrics<<==")
	metricsCh := make(chan DeviceMetrics)
	volumeCh := make(chan VolumeMetrics)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case met := <-metricsCh:
				m.vgTotalBytes.WithLabelValues(met.nodeName, met.DeviceGroup).Set(float64(met.TotalBytes))
				m.vgFreeBytes.WithLabelValues(met.nodeName, met.DeviceGroup).Set(float64(met.FreeBytes))
			case vc := <-volumeCh:
				m.volumeTotalBytes.WithLabelValues(vc.nodeName, vc.Volume).Set(float64(vc.TotalBytes))
				m.volumeUsedBytes.WithLabelValues(vc.nodeName, vc.Volume).Set(vc.UsedBytes)
			}
		}
	}()

	ticker := time.Tick(1 * time.Minute)
	for range ticker {
		dm, err := vgMetrics()
		if err == nil && len(dm) > 0 {
			for _, m := range dm {
				metricsCh <- m
			}
		}
		vm, err := volumeMetrics()
		if err == nil && len(vm) > 0 {
			for _, v := range vm {
				volumeCh <- v
			}
		}
	}
	return nil
}

func vgMetrics() ([]DeviceMetrics, error) {
	var metricsResult []DeviceMetrics
	metricsResult = append(metricsResult, DeviceMetrics{
		nodeName:    "test-1",
		FreeBytes:   uint64(10),
		TotalBytes:  uint64(10),
		DeviceGroup: "hxmysql",
	})
	return metricsResult, nil
}

func volumeMetrics() ([]VolumeMetrics, error) {
	var metricsResult []VolumeMetrics
	metricsResult = append(metricsResult, VolumeMetrics{
		nodeName:   "test-2",
		Volume:     "vol_1",
		TotalBytes: uint64(10),
		UsedBytes:  float64(10) * 1000 / 100,
	})

	return metricsResult, nil
}
