package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewProjectID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "minimum length", input: "prj_a1b2c3d4"},
		{name: "maximum length", input: "prj_" + strings.Repeat("a", 32)},
		{name: "empty", input: "", wantErr: true},
		{name: "invalid prefix", input: "project_a1b2c3d4", wantErr: true},
		{name: "too short", input: "prj_abc", wantErr: true},
		{name: "too long", input: "prj_" + strings.Repeat("a", 33), wantErr: true},
		{name: "uppercase", input: "prj_A1b2c3d4", wantErr: true},
		{name: "special characters", input: "prj_abc-1234", wantErr: true},
		{name: "Unicode", input: "prj_abcde123ż", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProjectID(tt.input)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidProjectID) {
					t.Fatalf("expected ErrInvalidProjectID, got %v", err)
				}
				if !got.IsZero() {
					t.Fatal("invalid ID should return zero value")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.input {
				t.Fatalf("String() = %q, want %q", got.String(), tt.input)
			}
			if got.IsZero() {
				t.Fatal("valid ID should not be zero")
			}
		})
	}
}

func TestProjectIDValueEquality(t *testing.T) {
	first, err := NewProjectID("prj_a1b2c3d4")
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewProjectID("prj_a1b2c3d4")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("IDs with equal values must compare equal")
	}
}
