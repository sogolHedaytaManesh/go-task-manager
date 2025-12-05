package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/entities"
	"task-manager/pkg/rest"
)

// MockTaskService
//
// A testify-based mock implementation of the TaskService interface.
// This mock is intended for unit testing HTTP handlers or other services
// that depend on TaskService without touching the real database.
type MockTaskService struct {
	mock.Mock
}

// Create mocks TaskService.Create
//
// Simulates creating a new task and returns the created *entities.Task
// and an error as configured in tests.
func (m *MockTaskService) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*entities.Task), args.Error(1)
}

// Update mocks TaskService.Update
//
// Simulates updating an existing task and returns the updated *entities.Task
// and an error as defined in test expectations.
func (m *MockTaskService) Update(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*entities.Task), args.Error(1)
}

// Delete mocks TaskService.Delete
//
// Simulates deleting a task by ID. Returns only an error (nil or configured)
// depending on test setup.
func (m *MockTaskService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetByID mocks TaskService.GetByID
//
// Returns a task that matches the provided ID or an error, depending
// on how the mock was configured in tests.
func (m *MockTaskService) GetByID(ctx context.Context, id int64) (*entities.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Task), args.Error(1)
}

// List mocks TaskService.List
//
// Returns a paginated list of tasks, the total count, and an error as
// defined in test expectations.
func (m *MockTaskService) List(ctx context.Context, query rest.Query) ([]entities.Task, int, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]entities.Task), args.Int(1), args.Error(2)
}

// ListByStatus mocks TaskService.ListByStatus
//
// Returns all tasks filtered by a given status and an error.
// Useful for testing handler logic that lists tasks by status.
func (m *MockTaskService) ListByStatus(ctx context.Context, status entities.TaskStatus) ([]entities.Task, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]entities.Task), args.Error(1)
}
