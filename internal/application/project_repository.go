package application

import (
	"context"
	"errors"

	"example.com/todo-ddd/internal/domain"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectRepository interface {
	Save(ctx context.Context, project *domain.Project) error

	FindByID(
		ctx context.Context,
		id domain.ProjectID,
	) (*domain.Project, error)
}
