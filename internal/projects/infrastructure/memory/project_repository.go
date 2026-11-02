package memory

import (
	"context"
	"errors"
	"sync"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/domain"
)

type ProjectRepository struct {
	mu       sync.RWMutex
	projects map[domain.ProjectID]domain.Project
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		projects: make(map[domain.ProjectID]domain.Project),
	}
}

func (r *ProjectRepository) Create(
	_ context.Context,
	project *domain.Project,
) error {
	if project == nil {
		return errors.New("project cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.projects[project.ID()]; exists {
		return application.ErrProjectAlreadyExists
	}
	r.projects[project.ID()] = *project

	return nil
}

func (r *ProjectRepository) FindByID(
	_ context.Context,
	id domain.ProjectID,
) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, exists := r.projects[id]
	if !exists {
		return nil, application.ErrProjectNotFound
	}

	return &project, nil
}
