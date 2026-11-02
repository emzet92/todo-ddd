package application_test

import (
	"context"
	"errors"
	"testing"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/domain"
	"example.com/todo-ddd/internal/projects/infrastructure/memory"
)

func TestTaskServiceProjectMustExist(t *testing.T) {
	projects := memory.NewProjectRepository()
	tasks := memory.NewTaskRepository()
	service := application.NewTaskService(projects, tasks)
	projectID, _ := domain.NewProjectID("prj_abcd1234")
	taskID, _ := domain.NewTaskID("tsk_abcd1234")
	title, _ := domain.NewTaskTitle("Some task")
	task, err := service.CreateTask(context.Background(), taskID, projectID, title)
	if task != nil || !errors.Is(err, application.ErrProjectNotFound) {
		t.Fatalf("task=%v err=%v", task, err)
	}
}

func TestTaskServiceCannotReadFromOtherProject(t *testing.T) {
	ctx := context.Background()
	projects := memory.NewProjectRepository()
	tasks := memory.NewTaskRepository()
	service := application.NewTaskService(projects, tasks)
	projectID, _ := domain.NewProjectID("prj_abcd1234")
	otherID, _ := domain.NewProjectID("prj_zzzzzzzz")
	name, _ := domain.NewProjectName("First")
	p, _ := domain.NewProject(projectID, name)
	if err := projects.Create(ctx, p); err != nil {
		t.Fatal(err)
	}
	taskID, _ := domain.NewTaskID("tsk_abcd1234")
	title, _ := domain.NewTaskTitle("Some task")
	if _, err := service.CreateTask(ctx, taskID, projectID, title); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetTask(ctx, otherID, taskID); !errors.Is(err, application.ErrTaskNotFound) {
		t.Fatalf("get: %v", err)
	}
	if _, err := service.StartTask(ctx, otherID, taskID); !errors.Is(err, application.ErrTaskNotFound) {
		t.Fatalf("start: %v", err)
	}
}
