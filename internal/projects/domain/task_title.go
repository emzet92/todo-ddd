package domain

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidTaskTitle = errors.New("invalid task title")

const maxTaskTitleLength = 200

// TaskTitle is an immutable value object containing 1..200 Unicode code points.
type TaskTitle struct {
	value string
}

func NewTaskTitle(raw string) (TaskTitle, error) {
	if !utf8.ValidString(raw) {
		return TaskTitle{}, ErrInvalidTaskTitle
	}

	value := strings.TrimSpace(raw)
	if value == "" || utf8.RuneCountInString(value) > maxTaskTitleLength {
		return TaskTitle{}, ErrInvalidTaskTitle
	}

	for _, c := range value {
		if unicode.IsControl(c) {
			return TaskTitle{}, ErrInvalidTaskTitle
		}
	}
	return TaskTitle{value: value}, nil
}

func (title TaskTitle) String() string { return title.value }
func (title TaskTitle) IsZero() bool   { return title.value == "" }
