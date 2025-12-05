package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"task-manager/pkg/rest"

	"github.com/cockroachdb/errors"
	"task-manager/internal/entities"
	"task-manager/pkg/db"
)

// -----------------------------------------------------------------------------
// Errors
// -----------------------------------------------------------------------------

// ErrTaskNotFound is returned when a task with the given ID does not exist.
var ErrTaskNotFound = fmt.Errorf("task not found")

// -----------------------------------------------------------------------------
// Interfaces
// -----------------------------------------------------------------------------

// TaskRepository defines all the operations required for interacting with tasks.
type TaskRepository interface {
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Update(ctx context.Context, task *entities.Task) (*entities.Task, error)
	GetByID(ctx context.Context, id int64) (*entities.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, query rest.Query) ([]entities.Task, int, error)
}

// -----------------------------------------------------------------------------
// Repository implementation
// -----------------------------------------------------------------------------

// Task implements TaskRepository using a SQL database.
type Task struct {
	db db.DB
}

// NewTaskRepository returns a new Task repository.
func NewTaskRepository(db db.DB) *Task {
	return &Task{db: db}
}

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

// Create inserts a new task and returns the created entity with generated fields.
func (r *Task) Create(ctx context.Context, t *entities.Task) (*entities.Task, error) {
	query := `
        INSERT INTO tasks (title, description, status, assignee_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
	err := r.db.GetContext(ctx, t, query,
		t.Title,
		t.Description,
		t.Status,
		t.AssigneeID,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

// -----------------------------------------------------------------------------
// GetByID
// -----------------------------------------------------------------------------

// GetByID fetches a task by its ID.
// Returns ErrTaskNotFound if no rows are returned.
func (r *Task) GetByID(ctx context.Context, id int64) (*entities.Task, error) {
	var task entities.Task

	query := `
        SELECT id, title, description, status, assignee_id, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

// -----------------------------------------------------------------------------
// List (with pagination + filters)
// -----------------------------------------------------------------------------

// List returns a list of tasks matching filters, along with the total count.
// Supports filtering by: status, assignee_id, title.
func (r *Task) List(ctx context.Context, query rest.Query) ([]entities.Task, int, error) {
	baseQuery := `
        SELECT id, title, description, status, assignee_id, created_at, updated_at
        FROM tasks
    `
	countQuery := `SELECT COUNT(*) FROM tasks`

	var (
		args       []interface{}
		conditions []string
		i          = 1
	)

	// Build WHERE filters dynamically
	for field, value := range query.Filter {
		switch field {
		case "status", "assignee_id", "title":
			conditions = append(conditions, fmt.Sprintf("%s = $%d", field, i))
			args = append(args, value)
			i++
		}
	}

	if len(conditions) > 0 {
		where := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += where
		countQuery += where
	}

	// Fetch total count
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (query.Page - 1) * query.PerPage
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, query.PerPage, offset)

	var tasks []entities.Task
	if err := r.db.SelectContext(ctx, &tasks, baseQuery, args...); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

// Update modifies an existing task and returns the updated entity.
// If the task does not exist, ErrTaskNotFound is returned.
func (r *Task) Update(ctx context.Context, t *entities.Task) (*entities.Task, error) {
	query := `
        UPDATE tasks
        SET title = $1,
            description = $2,
            status = $3,
            updated_at = now()
        WHERE id = $4
        RETURNING id, title, description, status, assignee_id, created_at, updated_at
    `

	err := r.db.GetContext(ctx, t, query,
		t.Title,
		t.Description,
		t.Status,
		t.ID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return t, nil
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

// Delete removes a task by ID.
// Returns ErrTaskNotFound if no record was deleted.
func (r *Task) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTaskNotFound
	}

	return nil
}
