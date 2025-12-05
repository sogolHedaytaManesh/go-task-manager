package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"task-manager/pkg/monitoring"
	"time"
)

// CORSMiddleware provides Cross-Origin Resource Sharing (CORS) support.
// This middleware allows the API to be accessed from different origins
// by setting appropriate headers. It also handles preflight OPTIONS requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// TaskMetricsMiddleware records Prometheus metrics for task-related endpoints.
// It captures request count and latency per HTTP method and response status.
// In a multi-pod deployment, each pod exposes its own metrics; Prometheus scrapes all pods
// and aggregates metrics at the cluster level.
func TaskMetricsMiddleware(metrics *monitoring.TaskMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency in milliseconds
		latencyMs := float64(time.Since(start).Milliseconds())

		status := c.Writer.Status()
		method := c.Request.Method

		// Record request count
		metrics.RequestCount.WithLabelValues(
			statusLabel(status),
			"handler_task",
		).Inc()

		// Record request latency
		metrics.RequestLatency.WithLabelValues(
			method,
			statusLabel(status),
			"handler_task",
		).Observe(latencyMs)
	}
}

// TracingMiddleware injects a unique trace ID into the request context for observability.
// This ID can be used for distributed tracing and correlating logs across services.
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), "traceID", traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// TraceIDFromContext retrieves the trace ID from the context, if available.
// Returns empty string if not found.
func TraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value("traceID"); v != nil {
		return v.(string)
	}
	return ""
}

// statusLabel converts HTTP status code to string for Prometheus labels.
// This ensures that metrics labels are consistent and compatible with Prometheus requirements.
func statusLabel(code int) string {
	return fmt.Sprintf("%d", code)
}
