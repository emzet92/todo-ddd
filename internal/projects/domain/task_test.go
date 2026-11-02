package domain_test

import (
	"errors"
	"testing"

	"example.com/todo-ddd/internal/projects/domain"
)

func taskFixture(t *testing.T) *domain.Task {
	t.Helper()
	id, _ := domain.NewTaskID("tsk_abcd1234")
	projectID, _ := domain.NewProjectID("prj_abcd1234")
	title, _ := domain.NewTaskTitle("Write tests")
	task, err := domain.NewTask(id, projectID, title)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func TestTaskLifecycle(t *testing.T) {
	task := taskFixture(t)
	if task.Status() != domain.TaskPlanned {
		t.Fatal("new task must be PLANNED")
	}
	if err := task.Complete(); !errors.Is(err, domain.ErrInvalidTaskTransition) {
		t.Fatalf("complete planned: %v", err)
	}
	if task.Status() != domain.TaskPlanned {
		t.Fatal("failed transition changed status")
	}
	if err := task.Start(); err != nil {
		t.Fatal(err)
	}
	if task.Status() != domain.TaskInProgress {
		t.Fatal("expected IN_PROGRESS")
	}
	if err := task.Start(); !errors.Is(err, domain.ErrInvalidTaskTransition) {
		t.Fatalf("start again: %v", err)
	}
	if err := task.Plan(); err != nil {
		t.Fatal(err)
	}
	if task.Status() != domain.TaskPlanned {
		t.Fatal("expected PLANNED")
	}
	if err := task.Start(); err != nil {
		t.Fatal(err)
	}
	if err := task.Complete(); err != nil {
		t.Fatal(err)
	}
	if task.Status() != domain.TaskDone {
		t.Fatal("expected DONE")
	}
	if err := task.Plan(); !errors.Is(err, domain.ErrInvalidTaskTransition) {
		t.Fatalf("reopen done: %v", err)
	}
	if err := task.Start(); !errors.Is(err, domain.ErrInvalidTaskTransition) {
		t.Fatalf("restart done: %v", err)
	}
	if task.Status() != domain.TaskDone {
		t.Fatal("done should be terminal")
	}
}

func TestTaskRejectsZeroValueObjects(t *testing.T) {
	id, _ := domain.NewTaskID("tsk_abcd1234")
	projectID, _ := domain.NewProjectID("prj_abcd1234")
	title, _ := domain.NewTaskTitle("Write tests")
	for _, tc := range []struct {
		name      string
		id        domain.TaskID
		projectID domain.ProjectID
		title     domain.TaskTitle
		want      error
	}{
		{"empty ID", domain.TaskID{}, projectID, title, domain.ErrInvalidTaskID},
		{"empty project ID", id, domain.ProjectID{}, title, domain.ErrInvalidProjectID},
		{"empty title", id, projectID, domain.TaskTitle{}, domain.ErrInvalidTaskTitle},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewTask(tc.id, tc.projectID, tc.title)
			if got != nil || !errors.Is(err, tc.want) {
				t.Fatalf("got %v, %v", got, err)
			}
		})
	}
}
