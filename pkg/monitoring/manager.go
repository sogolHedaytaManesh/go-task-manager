package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsManager
//
// Central manager for registering and keeping track of Prometheus metrics.
// Provides convenient methods to register counters, histograms, and gauges.
//
// In a multi-pod environment (e.g., Kubernetes):
//   - Each pod has its own in-memory metrics. Prometheus scrapes metrics from each pod separately.
//   - Labels such as 'service', 'pod', or 'instance' can be added when registering metrics
//     to differentiate between pods during aggregation.
type MetricsManager struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
	gauges     map[string]*prometheus.GaugeVec
	summaries  map[string]*prometheus.SummaryVec
}

// NewMetricsManager creates a new empty MetricsManager instance.
func NewMetricsManager() *MetricsManager {
	return &MetricsManager{
		counters:   make(map[string]*prometheus.CounterVec),
		histograms: make(map[string]*prometheus.HistogramVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
		summaries:  make(map[string]*prometheus.SummaryVec),
	}
}

// RegisterCounter creates and registers a new CounterVec with Prometheus.
// The 'labels' argument allows tagging metrics for multi-instance/pod environments.
func (m *MetricsManager) RegisterCounter(name, subsystem, help string, labels ...string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cleanNamespace(),
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	prometheus.MustRegister(counter)
	m.counters[name] = counter
	return counter
}

// RegisterHistogram creates and registers a new HistogramVec with Prometheus.
// Histograms are useful to observe request latencies or durations.
func (m *MetricsManager) RegisterHistogram(name, subsystem, help string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}

	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: cleanNamespace(),
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels,
	)
	prometheus.MustRegister(h)
	m.histograms[name] = h
	return h
}

// RegisterGauge creates and registers a new GaugeVec with Prometheus.
// Gauges represent current values (e.g., number of tasks in system) that can go up/down.
func (m *MetricsManager) RegisterGauge(name, subsystem, help string, labels ...string) *prometheus.GaugeVec {
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: cleanNamespace(),
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	prometheus.MustRegister(g)
	m.gauges[name] = g
	return g
}

// cleanNamespace returns the metrics namespace for this service.
// It's used consistently across all metrics to differentiate from other services.
func cleanNamespace() string {
	return "task_manager"
}

//
// Notes for Multi-Pod / Kubernetes:
//
// 1. Each pod runs a separate instance of MetricsManager; in-memory counters are not shared.
// 2. Prometheus scrapes each pod separately and aggregates metrics across pods.
// 3. To differentiate metrics per pod, you can add a label like 'pod' or 'instance' when registering metrics:
//      m.RegisterCounter("tasks_total", "service", "Total tasks created", "service", "pod")
// 4. Avoid trying to synchronize in-memory metrics across pods; let Prometheus handle aggregation.
//
