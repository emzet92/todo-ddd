package domain

import "errors"

var ErrInvalidTaskID = errors.New("invalid task ID")

// TaskID is an immutable value object: tsk_ + 8..32 lowercase ASCII letters/digits.
type TaskID struct {
	value string
}

func NewTaskID(raw string) (TaskID, error) {
	const prefix = "tsk_"
	if len(raw) < len(prefix) || raw[:len(prefix)] != prefix {
		return TaskID{}, ErrInvalidTaskID
	}

	suffix := raw[len(prefix):]
	if len(suffix) < 8 || len(suffix) > 32 {
		return TaskID{}, ErrInvalidTaskID
	}
	for _, c := range suffix {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return TaskID{}, ErrInvalidTaskID
		}
	}
	return TaskID{value: raw}, nil
}

func (id TaskID) String() string { return id.value }
func (id TaskID) IsZero() bool   { return id.value == "" }
