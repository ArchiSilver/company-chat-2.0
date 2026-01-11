package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Description     string     `json:"description,omitempty" db:"description"`
	Status          TaskStatus `json:"status" db:"status"`
	CreatorID       uuid.UUID  `json:"creator_id" db:"creator_id"`
	AssigneeID      *uuid.UUID `json:"assignee_id,omitempty" db:"assignee_id"`
	DueDate         *time.Time `json:"due_date,omitempty" db:"due_date"`
	SourceMessageID *int64     `json:"source_message_id,omitempty" db:"source_message_id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	GetAll(ctx context.Context, assigneeID *uuid.UUID, status *TaskStatus) ([]Task, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
