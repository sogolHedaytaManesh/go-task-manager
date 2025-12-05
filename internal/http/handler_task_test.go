package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/entities"
	HTTPhandler "task-manager/internal/http"
	"task-manager/internal/service"
	"task-manager/pkg/rest"
	"testing"
	"time"
)

var (
	errUnknown = errors.New("unknown error occurred")

	httpHandler *HTTPhandler.Handler

	stubBody = HTTPhandler.CreateTaskRequest{
		Title:       "Test",
		Description: "desc",
		Status:      entities.TaskStatusPending,
		AssigneeID:  10,
	}
	stubTask = entities.Task{
		ID:          12,
		Title:       stubBody.Title,
		Description: stubBody.Description,
		Status:      stubBody.Status,
		AssigneeID:  stubBody.AssigneeID,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	stubUpdateRequest = HTTPhandler.UpdateTaskRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      entities.TaskStatusInProgress,
		AssigneeID:  11,
	}
	stubUpdatedTask = entities.Task{
		ID:          12,
		Title:       stubUpdateRequest.Title,
		Description: stubUpdateRequest.Description,
		Status:      stubUpdateRequest.Status,
		AssigneeID:  stubUpdateRequest.AssigneeID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)

// TestTaskCreate_Success_ShouldReturnCreated tests successful creation of a task.
func TestTaskCreate_Success_ShouldReturnCreated(t *testing.T) {
	taskService := service.MockTaskService{}

	taskService.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("*entities.Task"),
	).Return(&stubTask, nil)

	httpHandler = HTTPhandler.SetupHandler(&taskService)
	router := httpHandler.SetupRouter()

	jsonBytes, _ := json.Marshal(stubBody)
	req, err := http.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var res rest.StandardResponse
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, float64(stubTask.ID), res.Data.(map[string]interface{})["id"])
	assert.Equal(t, stubTask.Title, res.Data.(map[string]interface{})["title"])
	assert.Equal(t, stubTask.Description, res.Data.(map[string]interface{})["description"])
	assert.Equal(t, string(stubTask.Status), res.Data.(map[string]interface{})["status"])
	assert.Equal(t, float64(stubTask.AssigneeID), res.Data.(map[string]interface{})["assignee_id"])

	taskService.AssertExpectations(t)
}

// TestTaskCreate_WithInvalidJSON_ShouldReturnBadRequest tests creation with invalid JSON payload.
func TestTaskCreate_WithInvalidJSON_ShouldReturnBadRequest(t *testing.T) {
	taskService := service.MockTaskService{}
	httpHandler = HTTPhandler.SetupHandler(&taskService)
	router := httpHandler.SetupRouter()

	req, err := http.NewRequest(http.MethodPost, "/api/tasks/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response rest.StandardResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, rest.Fail, response.Status)
}

// TestTaskCreate_InternalError_ShouldReturnInternalServerError tests handling of service errors during creation.
func TestTaskCreate_InternalError_ShouldReturnInternalServerError(t *testing.T) {
	taskService := service.MockTaskService{}

	taskService.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("*entities.Task"),
	).Return(&entities.Task{}, errUnknown)

	httpHandler = HTTPhandler.SetupHandler(&taskService)
	router := httpHandler.SetupRouter()

	jsonBytes, _ := json.Marshal(stubBody)
	req, err := http.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var res rest.StandardResponse
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, rest.Fail, res.Status)

	taskService.AssertExpectations(t)
}

// TestTaskCreate_InvalidStatus_ShouldReturnBadRequest tests task creation with an invalid status.
func TestTaskCreate_InvalidStatus_ShouldReturnBadRequest(t *testing.T) {
	taskService := service.MockTaskService{}
	httpHandler = HTTPhandler.SetupHandler(&taskService)
	router := httpHandler.SetupRouter()

	stubBody.Status = "invalid status"
	jsonBytes, _ := json.Marshal(stubBody)
	req, err := http.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var res rest.StandardResponse
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, rest.Fail, res.Status)
}
