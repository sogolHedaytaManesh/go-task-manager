package entities

import (
	"time"
)

// TaskStatus represents the allowed statuses for a Task.
type TaskStatus string

// -------------------------------
// Task Status Constants
// -------------------------------
// Define all valid statuses for tasks.
const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCanceled   TaskStatus = "canceled"
)

// Task represents the task entity, corresponding to the `tasks` table in the database.
type Task struct {
	ID          int64      `db:"id"`                    // Primary key
	Title       string     `db:"title"`                 // Task title
	Description string     `db:"description,omitempty"` // Optional task description
	Status      TaskStatus `db:"status"`                // Current task status
	AssigneeID  int64      `db:"assignee_id"`           // User assigned to the task
	CreatedAt   time.Time  `db:"created_at"`            // Timestamp when the task was created
	UpdatedAt   time.Time  `db:"updated_at"`            // Timestamp when the task was last updated
}

// -------------------------------
// TaskStatus Methods
// -------------------------------

// IsValid checks if the TaskStatus value is one of the allowed statuses.
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusDone, TaskStatusCanceled:
		return true
	default:
		return false
	}
}
