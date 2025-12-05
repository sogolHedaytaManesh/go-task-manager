package utils

import (
	"sync"
	"task-manager/pkg/monitoring"
)

var (
	initMetricsOnce   sync.Once
	globalTaskMetrics *monitoring.TaskMetrics
)

// InitGlobalTaskMetrics initializes the task metrics only once and returns the instance
func InitGlobalTaskMetrics() *monitoring.TaskMetrics {
	initMetricsOnce.Do(func() {
		metricsManager := monitoring.NewMetricsManager()
		globalTaskMetrics = monitoring.InitTaskMetrics(metricsManager)
	})

	return globalTaskMetrics
}
