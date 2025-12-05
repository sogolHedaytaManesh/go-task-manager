package service_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"task-manager/internal/repository/postgres"
	"task-manager/internal/service"
	"task-manager/internal/utils"
	"task-manager/pkg/rest"
	"testing"
)

// TestCreateTaskIntegration verifies that a task can be created successfully.
func TestCreateTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskService := service.MakeNewTaskService()
	task := service.RandomTask()

	res, err := taskService.Create(ctx, task)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.NotZero(t, res.ID)
	assert.Equal(t, task.Title, res.Title)
	assert.Equal(t, task.Description, res.Description)
	assert.Equal(t, task.AssigneeID, res.AssigneeID)
	assert.Equal(t, task.Status, res.Status)

	// Clean up database after the test
	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestUpdateTaskIntegration verifies that a task can be updated successfully.
func TestUpdateTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	createdTask := service.CreateTestTask()
	taskService := service.MakeNewTaskService()
	taskInstance := service.RandomTask()

	assert.NotEqual(t, taskInstance.Title, createdTask.Title)

	createdTask.Title = taskInstance.Title
	res, err := taskService.Update(ctx, createdTask)

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

// TestDeleteTaskIntegration verifies that a task can be deleted successfully
// and ensures that deleting a non-existent task returns the expected error.
func TestDeleteTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	createdTask := service.CreateTestTask()
	taskService := service.MakeNewTaskService()

	err := taskService.Delete(ctx, createdTask.ID)
	require.NoError(t, err)

	// Attempt to delete the same task again should fail
	err = taskService.Delete(ctx, createdTask.ID)
	require.Error(t, err)
	assert.Equal(t, postgres.ErrTaskNotFound, err)

	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}

// TestGetByIDTaskIntegration verifies that a task can be retrieved by its ID.
func TestGetByIDTaskIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	createdTask := service.CreateTestTask()
	taskService := service.MakeNewTaskService()

	res, err := taskService.GetByID(ctx, createdTask.ID)
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

// TestListTasksIntegration verifies that tasks can be listed with pagination.
func TestListTasksIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 7; i++ {
		_ = service.CreateTestTask()
	}

	taskService := service.MakeNewTaskService()
	query := rest.Query{
		PaginationMeta: rest.PaginationMeta{
			Page:    1,
			PerPage: 5,
		},
	}

	res, total, err := taskService.List(ctx, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, len(res), 5)
	assert.Equal(t, total, 7)

	t.Cleanup(func() {
		fmt.Println("🧹 Cleaning up after test...")
		utils.TruncateTables(t)
	})
}
