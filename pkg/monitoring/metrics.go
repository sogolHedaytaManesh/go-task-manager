package monitoring

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"strings"
)

// TaskMetrics
//
// Defines all Prometheus metrics for the task service.
// Includes counters, gauges, and histograms to measure
// task creation, request latency, request count, and current task load.
type TaskMetrics struct {
	TasksCount     *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
	RequestCount   *prometheus.CounterVec
	TasksCurrent   *prometheus.GaugeVec
}

// InitTaskMetrics
//
// Initializes all metrics using a MetricsManager.
// Each metric is labeled for service identification and Prometheus scraping.
func InitTaskMetrics(m *MetricsManager) *TaskMetrics {
	return &TaskMetrics{
		TasksCount: m.RegisterCounter(
			"tasks_total",
			"service",
			"Total tasks created",
			"service",
		),

		RequestLatency: m.RegisterHistogram(
			"request_latency_ms",
			"http",
			"Request latency in milliseconds",
			getBuckets(),
			"method", "status", "service",
		),

		RequestCount: m.RegisterCounter(
			"http_requests_total",
			"http",
			"Total HTTP requests",
			"status", "service",
		),

		TasksCurrent: m.RegisterGauge(
			"tasks_current",
			"service",
			"Current number of tasks in the system",
			"service",
		),
	}
}

// InitialGinMetrics
//
// Sets up Prometheus metrics scraping for a Gin HTTP server.
// If user/password are provided, it configures basic auth.
//
// Important:
// 1. In a multi-pod environment (e.g., Kubernetes), each pod exposes its own metrics endpoint.
// 2. Prometheus scrapes each pod individually using a Service or Pod annotations.
// 3. Labels like 'service' or 'pod' are crucial to distinguish metrics from different instances.
func InitialGinMetrics(e *gin.Engine, metricsPath string, metricsPort int, user string, password string) *ginprometheus.Prometheus {
	ginProm := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem:          "gin",
		MetricsList:        nil,
		CustomLabels:       nil,
		DisableBodyReading: false,
	})

	// Override default metrics path
	if metricsPath != "" {
		ginProm.MetricsPath = metricsPath
	}

	// Override default listening port
	if metricsPort != 0 {
		ginProm.SetListenAddress(fmt.Sprintf(":%d", metricsPort))
	}

	// Basic authentication for metrics endpoint (optional)
	if user != "" && password != "" {
		fmt.Printf("Setting up metrics endpoint with basic authentication, user: %s, password: %s\n", user, password)
		ginProm.UseWithAuth(e, gin.Accounts{
			user: password,
		})
	} else {
		ginProm.Use(e)
	}

	// Replace actual values in URL path with param placeholders for consistent metrics labeling
	// e.g., /api/tasks/123 -> /api/tasks/:id
	ginProm.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		}
		return url
	}

	return ginProm
}

// getBuckets local func to return an array of thresholds
func getBuckets() []float64 {
	return []float64{
		0.005, // 5ms
		0.01,  // 10ms
		0.025, // 25ms
		0.05,  // 50ms
		0.1,   // 100ms
		0.25,  // 250ms
		0.5,   // 500ms
		1.0,   // 1s
		2.5,   // 2.5s
		5.0,   // 5s
		10.0,  // 10s
	}
}

//
// Notes on Multi-Pod Setup:
//
// - Each pod runs its own instance of the Gin server with its own /metrics endpoint.
// - Prometheus discovers pods using Kubernetes Service discovery and scrapes each pod individually.
// - Labels such as 'pod', 'namespace', 'service' are automatically attached in Kubernetes to differentiate metrics.
// - Avoid sharing in-memory metrics across pods; rely on Prometheus to aggregate.
//
