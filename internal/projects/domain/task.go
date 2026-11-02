package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidTaskTransition = errors.New("invalid task status transition")

type TaskStatus string

const (
	TaskPlanned    TaskStatus = "PLANNED"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskDone       TaskStatus = "DONE"
)

// Task is its own aggregate; it references a Project by ID, not by a pointer.
type Task struct {
	id        TaskID
	projectID ProjectID
	title     TaskTitle
	status    TaskStatus
}

func NewTask(id TaskID, projectID ProjectID, title TaskTitle) (*Task, error) {
	if id.IsZero() {
		return nil, ErrInvalidTaskID
	}
	if projectID.IsZero() {
		return nil, ErrInvalidProjectID
	}
	if title.IsZero() {
		return nil, ErrInvalidTaskTitle
	}
	return &Task{id: id, projectID: projectID, title: title, status: TaskPlanned}, nil
}

func (t *Task) Start() error {
	if t.status != TaskPlanned {
		return t.invalidTransition(TaskInProgress)
	}
	t.status = TaskInProgress
	return nil
}

func (t *Task) Complete() error {
	if t.status != TaskInProgress {
		return t.invalidTransition(TaskDone)
	}
	t.status = TaskDone
	return nil
}

func (t *Task) Plan() error {
	if t.status != TaskInProgress {
		return t.invalidTransition(TaskPlanned)
	}
	t.status = TaskPlanned
	return nil
}

func (t *Task) invalidTransition(target TaskStatus) error {
	return fmt.Errorf("%w: %s -> %s", ErrInvalidTaskTransition, t.status, target)
}

func (t *Task) ID() TaskID           { return t.id }
func (t *Task) ProjectID() ProjectID { return t.projectID }
func (t *Task) Title() TaskTitle     { return t.title }
func (t *Task) Status() TaskStatus   { return t.status }
