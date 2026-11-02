package domain_test

import (
	"errors"
	"testing"

	"example.com/todo-ddd/internal/projects/domain"
)

func TestTaskIDValidation(t *testing.T) {
	for _, tc := range []struct {
		name, raw string
		valid     bool
	}{
		{"valid", "tsk_abcd1234", true},
		{"min", "tsk_12345678", true},
		{"max", "tsk_" + "12345678901234567890123456789012", true},
		{"empty", "", false},
		{"prefix", "prj_abcd1234", false},
		{"short", "tsk_123", false},
		{"long", "tsk_123456789012345678901234567890123", false},
		{"uppercase", "tsk_ABCD1234", false},
		{"underscore", "tsk_abcd_123", false},
		{"unicode", "tsk_żąbka123", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id, err := domain.NewTaskID(tc.raw)
			if tc.valid {
				if err != nil || id.String() != tc.raw || id.IsZero() {
					t.Fatalf("got id=%q err=%v", id, err)
				}
			} else if !errors.Is(err, domain.ErrInvalidTaskID) || !id.IsZero() {
				t.Fatalf("expected invalid task id, got id=%q err=%v", id, err)
			}
		})
	}
}

func TestTaskIDComparable(t *testing.T) {
	first, _ := domain.NewTaskID("tsk_abcd1234")
	second, _ := domain.NewTaskID("tsk_abcd1234")
	if first != second {
		t.Fatal("equal value objects must compare equal")
	}
}
