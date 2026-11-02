package application

import (
	"context"
	"errors"

	"example.com/todo-ddd/internal/projects/domain"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
)

type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error

	FindByID(
		ctx context.Context,
		id domain.ProjectID,
	) (*domain.Project, error)
}
