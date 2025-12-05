package mock

import (
	"context"
	"task-manager/internal/entities"
	"task-manager/pkg/rest"

	"github.com/stretchr/testify/mock"
)

// MockTaskRepository
//
// A testify-based mock implementation of the TaskRepository interface.
// This mock is used for unit testing service and handler layers without
// touching the real database.
type MockTaskRepository struct {
	mock.Mock
}

// Create mocks TaskRepository.Create
//
// It returns the created *entities.Task and an error based on what the test
// has configured using: mockRepo.On("Create", ...).Return(...)
func (m *MockTaskRepository) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	args := m.Called(ctx, task)

	var t *entities.Task
	if args.Get(0) != nil {
		t = args.Get(0).(*entities.Task)
	}

	return t, args.Error(1)
}

// Update mocks TaskRepository.Update
//
// It simulates updating an existing task and returns the updated *entities.Task
// or an error, depending on the expected output defined in the test.
func (m *MockTaskRepository) Update(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	args := m.Called(ctx, task)

	var t *entities.Task
	if args.Get(0) != nil {
		t = args.Get(0).(*entities.Task)
	}

	return t, args.Error(1)
}

// GetByID mocks TaskRepository.GetByID
//
// It returns a task that matches the provided ID or an error.
// Test examples:
//
//	mockRepo.On("GetByID", mock.Anything, int64(10)).Return(task, nil)
func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*entities.Task, error) {
	args := m.Called(ctx, id)

	var t *entities.Task
	if args.Get(0) != nil {
		t = args.Get(0).(*entities.Task)
	}

	return t, args.Error(1)
}

// Delete mocks TaskRepository.Delete
//
// It simulates removing a task by ID.
// Returns only an error (nil or error), depending on configured expectations.
func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks TaskRepository.List
//
// It simulates fetching a paginated list of tasks.
// Returns:
//   - []entities.Task : list of tasks
//   - int             : total count
//   - error           : optional error
//
// Example mock setup:
//
//	mockRepo.
//	    On("List", mock.Anything, query).
//	    Return([]entities.Task{...}, 10, nil)
func (m *MockTaskRepository) List(ctx context.Context, query rest.Query) ([]entities.Task, int, error) {
	args := m.Called(ctx, query)

	var tasks []entities.Task
	if args.Get(0) != nil {
		tasks = args.Get(0).([]entities.Task)
	}

	return tasks, args.Int(1), args.Error(2)
}
