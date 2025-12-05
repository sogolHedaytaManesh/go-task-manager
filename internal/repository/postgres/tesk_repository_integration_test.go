package postgres_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"task-manager/internal/repository"
	"task-manager/internal/repository/postgres"
	"task-manager/internal/utils"
	"task-manager/pkg/rest"
	"testing"
)

// TestCreateTaskIntegration tests that a new task can be successfully
// created in the database and all fields are returned correctly.
func TestCreateTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	taskRepository := repository.MakeNewTaskRepository()

	taskInstance := repository.RandomTask()

	res, err := taskRepository.Create(context.Background(), taskInstance)

	assert.Nil(t, err)
	assert.Equal(t, taskInstance.ID, res.ID)
	assert.Equal(t, taskInstance.Title, res.Title)
	assert.Equal(t, taskInstance.AssigneeID, res.AssigneeID)
	assert.Equal(t, taskInstance.Description, res.Description)
	assert.NotNil(t, taskInstance.CreatedAt)
	assert.NotNil(t, taskInstance.CreatedAt)
	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestUpdateTaskIntegration tests that an existing task can be updated
// and changes are persisted in the database.
func TestUpdateTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	createdTask := repository.CreateTestTask()

	taskRepository := repository.MakeNewTaskRepository()

	taskInstance := repository.RandomTask()

	assert.NotEqual(t, taskInstance.Title, createdTask.Title)

	createdTask.Title = taskInstance.Title

	res, err := taskRepository.Update(ctx, createdTask)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.NotZero(t, res.ID)
	assert.Equal(t, taskInstance.Title, res.Title)
	assert.Equal(t, createdTask.Description, res.Description)
	assert.Equal(t, createdTask.AssigneeID, res.AssigneeID)
	assert.Equal(t, createdTask.Status, res.Status)
	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestDeleteTaskIntegration tests deleting a task by ID. It also verifies
// that attempting to delete the same task again returns ErrTaskNotFound.
func TestDeleteTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	createdTask := repository.CreateTestTask()

	taskRepository := repository.MakeNewTaskRepository()

	err := taskRepository.Delete(ctx, createdTask.ID)

	require.NoError(t, err)

	err = taskRepository.Delete(ctx, createdTask.ID)

	require.Error(t, err)

	assert.Equal(t, postgres.ErrTaskNotFound, err)

	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestGetTaskByIDIntegration tests fetching a task by ID from the database.
// It ensures all fields are returned correctly.
func TestGetTaskByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	createdTask := repository.CreateTestTask()

	taskRepository := repository.MakeNewTaskRepository()

	res, err := taskRepository.GetByID(ctx, createdTask.ID)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.NotZero(t, res.ID)
	assert.Equal(t, createdTask.Title, res.Title)
	assert.Equal(t, createdTask.Description, res.Description)
	assert.Equal(t, createdTask.AssigneeID, res.AssigneeID)
	assert.Equal(t, createdTask.Status, res.Status)

	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestListTasksIntegration tests listing tasks with pagination.
// It ensures correct number of tasks and total count is returned.
func TestListTasksIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()
	_ = repository.CreateTestTask()

	taskRepository := repository.MakeNewTaskRepository()

	query := rest.Query{
		PaginationMeta: rest.PaginationMeta{
			Page:    1,
			PerPage: 5,
		},
	}

	res, total, err := taskRepository.List(ctx, query)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, len(res), 5)
	assert.Equal(t, total, 7)

	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}
