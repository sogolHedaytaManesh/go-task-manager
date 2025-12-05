package repository

import (
	"context"
	"task-manager/internal/repository/postgres"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"task-manager/internal/entities"
	"task-manager/internal/utils"
)

// MakeNewTaskRepository initializes a new TaskRepository using the test database connection.
func MakeNewTaskRepository() postgres.TaskRepository {
	return postgres.NewTaskRepository(utils.CreateTestDatabaseConnection())
}

// CreateTestTask generates a random task and saves it to the test database.
// Returns the created Task entity. Panics if creation fails.
func CreateTestTask() *entities.Task {
	repo := MakeNewTaskRepository()

	task := RandomTask()
	createdTask, err := repo.Create(context.Background(), task)
	if err != nil {
		panic(err)
	}

	return createdTask
}

// RandomTask generates a new Task entity with randomized values for testing purposes.
// The ID, Title, Description, Status, and AssigneeID fields are populated randomly.
// CreatedAt and UpdatedAt are set to the current timestamp.
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
