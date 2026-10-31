package application

import (
	"context"
	"fmt"

	"example.com/todo-ddd/internal/domain"
)

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	id domain.ProjectID,
	name domain.ProjectName,
) (*domain.Project, error) {

	project, err := domain.NewProject(id, name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, fmt.Errorf("save project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) GetProject(
	ctx context.Context,
	id domain.ProjectID,
) (*domain.Project, error) {

	return s.repo.FindByID(ctx, id)
}
