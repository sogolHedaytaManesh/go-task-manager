package service

import (
	"context"
	"task-manager/internal/entities"
	"task-manager/internal/repository/postgres"
	"task-manager/pkg/monitoring"
	"task-manager/pkg/rest"
	"time"
)

// TaskService
//
// Interface defining the available task-related operations.
// All operations are context-aware and return standard Go errors.
type TaskService interface {
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Update(ctx context.Context, task *entities.Task) (*entities.Task, error)
	GetByID(ctx context.Context, id int64) (*entities.Task, error)
	List(ctx context.Context, query rest.Query) ([]entities.Task, int, error)
	Delete(ctx context.Context, id int64) error
}

// Task
//
// Concrete implementation of TaskService. Wraps a repository and metrics
// to perform database operations and record Prometheus metrics for each request.
type Task struct {
	taskRepo postgres.TaskRepository
	metrics  *monitoring.TaskMetrics
}

// NewTaskService
//
// Constructs a new TaskService with the provided repository and metrics manager.
func NewTaskService(taskRepo postgres.TaskRepository, metrics *monitoring.TaskMetrics) TaskService {
	return &Task{
		taskRepo: taskRepo,
		metrics:  metrics,
	}
}

// Create
//
// Creates a new task in the database and increments metrics counters.
// Records the request latency in Prometheus.
func (t *Task) Create(ctx context.Context, task *entities.Task) (createdTask *entities.Task, err error) {
	start := time.Now()

	createdTask, err = t.taskRepo.Create(ctx, task)

	if err == nil {
		t.metrics.TasksCount.WithLabelValues("task_service").Inc()
	}

	t.metrics.RequestLatency.
		WithLabelValues("POST", statusLabel(err), "task_service").
		Observe(float64(time.Since(start).Milliseconds()))

	return
}

// GetByID
//
// Fetches a task by its ID from the repository.
// Records request latency in Prometheus.
func (t *Task) GetByID(ctx context.Context, id int64) (task *entities.Task, err error) {
	start := time.Now()

	task, err = t.taskRepo.GetByID(ctx, id)

	t.metrics.RequestLatency.
		WithLabelValues("GET", statusLabel(err), "task_service").
		Observe(float64(time.Since(start).Milliseconds()))

	return
}

// List
//
// Fetches a list of tasks with pagination and optional filters.
// Returns the tasks slice, total count, and an error if any.
// Records request latency in Prometheus.
func (t *Task) List(ctx context.Context, query rest.Query) (tasks []entities.Task, total int, err error) {
	start := time.Now()

	tasks, total, err = t.taskRepo.List(ctx, query)

	t.metrics.RequestLatency.
		WithLabelValues("GET", statusLabel(err), "task_service").
		Observe(float64(time.Since(start).Milliseconds()))

	return
}

// Update
//
// Updates an existing task and records request latency.
// Returns the updated task and an error if any.
func (t *Task) Update(ctx context.Context, task *entities.Task) (updatedTask *entities.Task, err error) {
	start := time.Now()

	updatedTask, err = t.taskRepo.Update(ctx, task)

	t.metrics.RequestLatency.
		WithLabelValues("PUT", statusLabel(err), "task_service").
		Observe(float64(time.Since(start).Milliseconds()))

	return
}

// Delete
//
// Deletes a task by ID. Updates metrics counters and records request latency.
func (t *Task) Delete(ctx context.Context, id int64) (err error) {
	start := time.Now()

	err = t.taskRepo.Delete(ctx, id)

	if err == nil {
		t.metrics.TasksCount.WithLabelValues("task_service").Desc()
	}

	t.metrics.RequestLatency.
		WithLabelValues("DELETE", statusLabel(err), "task_service").
		Observe(float64(time.Since(start).Milliseconds()))

	return
}

// statusLabel
//
// Helper function to map error presence to a Prometheus metric label.
func statusLabel(err error) string {
	if err != nil {
		return "error"
	}

	return "success"
}
