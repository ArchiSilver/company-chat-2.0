package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	taskRepo    domain.TaskRepository
	messageRepo domain.MessageRepository
}

func NewTaskService(taskRepo domain.TaskRepository, messageRepo domain.MessageRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		messageRepo: messageRepo,
	}
}

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

func (s *TaskService) Create(ctx context.Context, creatorID uuid.UUID, input CreateTaskInput) (*domain.Task, error) {
	task := &domain.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.TaskStatusTodo,
		CreatorID:   creatorID,
		AssigneeID:  input.AssigneeID,
		DueDate:     input.DueDate,
	}

	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

type CreateFromMessageInput struct {
	MessageID  int64      `json:"message_id" binding:"required"`
	AssigneeID *uuid.UUID `json:"assignee_id"`
	DueDate    *time.Time `json:"due_date"`
}

func (s *TaskService) CreateFromMessage(ctx context.Context, creatorID uuid.UUID, input CreateFromMessageInput) (*domain.Task, error) {
	messages, err := s.messageRepo.GetMessagesByChatID(ctx, uuid.Nil, 1, 0)
	if err != nil {
		return nil, err
	}

	var messageContent string
	for _, msg := range messages {
		if msg.ID == input.MessageID {
			messageContent = msg.Content
			break
		}
	}

	if messageContent == "" {
		messageContent = "Task from message"
	}

	task := &domain.Task{
		Title:           messageContent,
		Description:     "Created from chat message",
		Status:          domain.TaskStatusTodo,
		CreatorID:       creatorID,
		AssigneeID:      input.AssigneeID,
		DueDate:         input.DueDate,
		SourceMessageID: &input.MessageID,
	}

	err = s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetAll(ctx context.Context, assigneeID *uuid.UUID, status *domain.TaskStatus) ([]domain.Task, error) {
	return s.taskRepo.GetAll(ctx, assigneeID, status)
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

type UpdateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, input UpdateTaskInput) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	task.AssigneeID = input.AssigneeID
	task.DueDate = input.DueDate

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

type UpdateStatusInput struct {
	Status domain.TaskStatus `json:"status" binding:"required"`
}

func (s *TaskService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	return s.taskRepo.UpdateStatus(ctx, id, status)
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.Delete(ctx, id)
}
