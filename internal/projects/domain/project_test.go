package domain

import (
	"errors"
	"testing"
)

func TestNewProject(t *testing.T) {
	id, err := NewProjectID("prj_a1b2c3d4")
	if err != nil {
		t.Fatal(err)
	}
	name, err := NewProjectName("My Project")
	if err != nil {
		t.Fatal(err)
	}

	project, err := NewProject(id, name)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID() != id {
		t.Fatalf("ID() = %q, want %q", project.ID(), id)
	}
	if project.Name() != name {
		t.Fatalf("Name() = %q, want %q", project.Name(), name)
	}
}

func TestNewProjectRejectsZeroValues(t *testing.T) {
	id, err := NewProjectID("prj_a1b2c3d4")
	if err != nil {
		t.Fatal(err)
	}
	name, err := NewProjectName("My Project")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("zero ID", func(t *testing.T) {
		project, err := NewProject(ProjectID{}, name)
		if !errors.Is(err, ErrInvalidProjectID) || project != nil {
			t.Fatalf("expected nil project and ErrInvalidProjectID, got %v, %v", project, err)
		}
	})

	t.Run("zero name", func(t *testing.T) {
		project, err := NewProject(id, ProjectName{})
		if !errors.Is(err, ErrInvalidProjectName) || project != nil {
			t.Fatalf("expected nil project and ErrInvalidProjectName, got %v, %v", project, err)
		}
	})
}
