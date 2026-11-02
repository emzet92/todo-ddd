package application

import (
	"context"
	"fmt"

	"example.com/todo-ddd/internal/projects/domain"
)

type TaskService struct {
	projects ProjectRepository
	tasks    TaskRepository
}

func NewTaskService(projects ProjectRepository, tasks TaskRepository) *TaskService {
	return &TaskService{projects: projects, tasks: tasks}
}

func (s *TaskService) CreateTask(ctx context.Context, id domain.TaskID, projectID domain.ProjectID, title domain.TaskTitle) (*domain.Task, error) {
	if _, err := s.projects.FindByID(ctx, projectID); err != nil {
		return nil, err
	}
	task, err := domain.NewTask(id, projectID, title)
	if err != nil {
		return nil, err
	}
	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, projectID domain.ProjectID, id domain.TaskID) (*domain.Task, error) {
	task, err := s.tasks.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.ProjectID() != projectID {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, projectID domain.ProjectID) ([]*domain.Task, error) {
	if _, err := s.projects.FindByID(ctx, projectID); err != nil {
		return nil, err
	}
	return s.tasks.FindByProjectID(ctx, projectID)
}

func (s *TaskService) StartTask(ctx context.Context, projectID domain.ProjectID, id domain.TaskID) (*domain.Task, error) {
	return s.changeStatus(ctx, projectID, id, (*domain.Task).Start)
}

func (s *TaskService) CompleteTask(ctx context.Context, projectID domain.ProjectID, id domain.TaskID) (*domain.Task, error) {
	return s.changeStatus(ctx, projectID, id, (*domain.Task).Complete)
}

func (s *TaskService) PlanTask(ctx context.Context, projectID domain.ProjectID, id domain.TaskID) (*domain.Task, error) {
	return s.changeStatus(ctx, projectID, id, (*domain.Task).Plan)
}

func (s *TaskService) changeStatus(ctx context.Context, projectID domain.ProjectID, id domain.TaskID, transition func(*domain.Task) error) (*domain.Task, error) {
	task, err := s.GetTask(ctx, projectID, id)
	if err != nil {
		return nil, err
	}
	if err := transition(task); err != nil {
		return nil, err
	}
	if err := s.tasks.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save task: %w", err)
	}
	return task, nil
}
