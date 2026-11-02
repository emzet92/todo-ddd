package application

import (
	"context"
	"errors"

	"example.com/todo-ddd/internal/projects/domain"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	FindByID(ctx context.Context, id domain.TaskID) (*domain.Task, error)
	FindByProjectID(ctx context.Context, projectID domain.ProjectID) ([]*domain.Task, error)
	Save(ctx context.Context, task *domain.Task) error
}
