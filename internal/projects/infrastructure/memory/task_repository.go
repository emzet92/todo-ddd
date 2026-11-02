package memory

import (
	"context"
	"errors"
	"sort"
	"sync"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/domain"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[domain.TaskID]domain.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: make(map[domain.TaskID]domain.Task)}
}

func (r *TaskRepository) Create(_ context.Context, task *domain.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID()]; exists {
		return application.ErrTaskAlreadyExists
	}
	r.tasks[task.ID()] = *task
	return nil
}

func (r *TaskRepository) FindByID(_ context.Context, id domain.TaskID) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return nil, application.ErrTaskNotFound
	}
	return &task, nil
}

func (r *TaskRepository) FindByProjectID(_ context.Context, projectID domain.ProjectID) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.ProjectID() == projectID {
			copyOfTask := task
			tasks = append(tasks, &copyOfTask)
		}
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID().String() < tasks[j].ID().String() })
	return tasks, nil
}

func (r *TaskRepository) Save(_ context.Context, task *domain.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID()]; !exists {
		return application.ErrTaskNotFound
	}
	r.tasks[task.ID()] = *task
	return nil
}

var _ application.TaskRepository = (*TaskRepository)(nil)
