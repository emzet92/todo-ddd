package domain

import "errors"

var ErrInvalidProjectID = errors.New("invalid project ID")

// ProjectID is a comparable, immutable value object.
// Its format is prj_ followed by 8-32 lowercase ASCII letters or digits.
type ProjectID struct {
	value string
}

func NewProjectID(raw string) (ProjectID, error) {
	const prefix = "prj_"

	if len(raw) < len(prefix) || raw[:len(prefix)] != prefix {
		return ProjectID{}, ErrInvalidProjectID
	}

	suffix := raw[len(prefix):]
	if len(suffix) < 8 || len(suffix) > 32 {
		return ProjectID{}, ErrInvalidProjectID
	}

	for _, c := range suffix {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return ProjectID{}, ErrInvalidProjectID
		}
	}

	return ProjectID{value: raw}, nil
}

func (id ProjectID) String() string {
	return id.value
}

func (id ProjectID) IsZero() bool {
	return id.value == ""
}
