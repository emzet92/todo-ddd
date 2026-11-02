package domain

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidProjectName = errors.New("invalid project name")

const maxProjectNameLength = 100

// ProjectName is a comparable, immutable value object.
// A name has 1-100 Unicode code points after trimming surrounding whitespace.
type ProjectName struct {
	value string
}

func NewProjectName(raw string) (ProjectName, error) {
	if !utf8.ValidString(raw) {
		return ProjectName{}, ErrInvalidProjectName
	}

	value := strings.TrimSpace(raw)
	if value == "" || utf8.RuneCountInString(value) > maxProjectNameLength {
		return ProjectName{}, ErrInvalidProjectName
	}

	for _, c := range value {
		if unicode.IsControl(c) {
			return ProjectName{}, ErrInvalidProjectName
		}
	}

	return ProjectName{value: value}, nil
}

func (name ProjectName) String() string {
	return name.value
}

func (name ProjectName) IsZero() bool {
	return name.value == ""
}
