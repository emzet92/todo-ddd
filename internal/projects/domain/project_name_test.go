package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "My Project", want: "My Project"},
		{name: "trims surrounding whitespace", input: "  My Project  ", want: "My Project"},
		{name: "Unicode", input: "Zażółć gęślą", want: "Zażółć gęślą"},
		{name: "100 Unicode characters", input: strings.Repeat("ą", 100), want: strings.Repeat("ą", 100)},
		{name: "empty", input: "", wantErr: true},
		{name: "whitespace only", input: " \t ", wantErr: true},
		{name: "101 Unicode characters", input: strings.Repeat("ą", 101), wantErr: true},
		{name: "control character", input: "Hello\nWorld", wantErr: true},
		{name: "invalid UTF-8", input: string([]byte{0xff, 0xfe}), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProjectName(tt.input)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidProjectName) {
					t.Fatalf("expected ErrInvalidProjectName, got %v", err)
				}
				if !got.IsZero() {
					t.Fatal("invalid name should return zero value")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Fatalf("String() = %q, want %q", got.String(), tt.want)
			}
			if got.IsZero() {
				t.Fatal("valid name should not be zero")
			}
		})
	}
}

func TestProjectNameValueEquality(t *testing.T) {
	first, err := NewProjectName("  My Project ")
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewProjectName("My Project")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("normalized names with equal values must compare equal")
	}
}
