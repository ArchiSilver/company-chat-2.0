package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `INSERT INTO tasks.tasks (title, description, status, creator_id, assignee_id, due_date, source_message_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query,
		task.Title, task.Description, task.Status, task.CreatorID, task.AssigneeID, task.DueDate, task.SourceMessageID,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	query := `SELECT id, title, description, status, creator_id, assignee_id, due_date, source_message_id, created_at, updated_at
              FROM tasks.tasks WHERE id = $1`
	err := r.db.GetContext(ctx, &task, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (r *TaskRepository) GetAll(ctx context.Context, assigneeID *uuid.UUID, status *domain.TaskStatus) ([]domain.Task, error) {
	var tasks []domain.Task
	query := `SELECT id, title, description, status, creator_id, assignee_id, due_date, source_message_id, created_at, updated_at
              FROM tasks.tasks WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if assigneeID != nil {
		query += ` AND assignee_id = $` + string(rune('0'+argIndex))
		args = append(args, *assigneeID)
		argIndex++
	}
	if status != nil {
		query += ` AND status = $` + string(rune('0'+argIndex))
		args = append(args, *status)
	}

	query += ` ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &tasks, query, args...)
	return tasks, err
}

func (r *TaskRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]domain.Task, error) {
	var tasks []domain.Task
	query := `SELECT id, title, description, status, creator_id, assignee_id, due_date, source_message_id, created_at, updated_at
              FROM tasks.tasks 
              WHERE (creator_id = $1 OR assignee_id = $1) 
              AND created_at >= $2 AND created_at <= $3
              ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &tasks, query, userID, from, to)
	return tasks, err
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `UPDATE tasks.tasks SET title = $1, description = $2, status = $3, assignee_id = $4, due_date = $5, updated_at = $6
              WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.AssigneeID, task.DueDate, time.Now(), task.ID)
	return err
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	query := `UPDATE tasks.tasks SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks.tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
