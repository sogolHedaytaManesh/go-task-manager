package service

import (
	"context"
	"task-manager/internal/repository/postgres"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"task-manager/internal/entities"
	"task-manager/internal/utils"
)

// MakeNewTaskService creates a new TaskService instance for testing or production.
// It initializes a TaskRepository with a test database connection and global metrics.
func MakeNewTaskService() TaskService {
	repo := postgres.NewTaskRepository(
		utils.CreateTestDatabaseConnection(),
	)

	return NewTaskService(repo, utils.InitGlobalTaskMetrics())
}

// CreateTestTask creates a random task in the database for testing purposes.
// Panics if creation fails. Used to pre-populate data in integration tests.
func CreateTestTask() *entities.Task {
	service := MakeNewTaskService()

	task, err := service.Create(
		context.Background(),
		RandomTask(),
	)
	if err != nil {
		panic(err)
	}

	return task
}

// RandomTask generates a random Task entity with realistic data using gofakeit.
// It is used for testing or seeding purposes.
func RandomTask() *entities.Task {
	statuses := []entities.TaskStatus{
		entities.TaskStatusPending,
		entities.TaskStatusInProgress,
		entities.TaskStatusDone,
	}

	return &entities.Task{
		ID:          int64(gofakeit.Number(1, 1000)),
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Paragraph(1, 2, 5, " "),
		Status:      statuses[gofakeit.Number(0, len(statuses)-1)],
		AssigneeID:  int64(gofakeit.Number(1, 50)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
