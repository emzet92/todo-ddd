package domain_test

import (
	"errors"
	"strings"
	"testing"

	"example.com/todo-ddd/internal/projects/domain"
)

func TestTaskTitleValidation(t *testing.T) {
	for _, tc := range []struct {
		name, raw, want string
		valid           bool
	}{
		{"trim", "  Fix login  ", "Fix login", true},
		{"unicode", "  Zażółć 🐸  ", "Zażółć 🐸", true},
		{"max", strings.Repeat("x", 200), strings.Repeat("x", 200), true},
		{"empty", "", "", false},
		{"blank", " \t ", "", false},
		{"too long", strings.Repeat("x", 201), "", false},
		{"newline", "Line\nBreak", "", false},
		{"invalid utf8", string([]byte{0xff}), "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			title, err := domain.NewTaskTitle(tc.raw)
			if tc.valid {
				if err != nil || title.String() != tc.want || title.IsZero() {
					t.Fatalf("title=%q err=%v", title, err)
				}
			} else if !errors.Is(err, domain.ErrInvalidTaskTitle) || !title.IsZero() {
				t.Fatalf("expected invalid title, got title=%q err=%v", title, err)
			}
		})
	}
}
