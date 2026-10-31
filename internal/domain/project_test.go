package domain

import "testing"

func TestNewProject(t *testing.T) {
	tests := []struct {
		name        string
		id          ProjectID
		projectName ProjectName
		wantName    ProjectName
		wantErr     bool
	}{
		{
			name:        "creates valid project",
			id:          "project-1",
			projectName: "My Project",
			wantName:    "My Project",
		},
		{
			name:        "trims project name",
			id:          "project-1",
			projectName: "  My Project  ",
			wantName:    "My Project",
		},
		{
			name:        "rejects empty ID",
			id:          "",
			projectName: "My Project",
			wantErr:     true,
		},
		{
			name:        "rejects empty name",
			id:          "project-1",
			projectName: "",
			wantErr:     true,
		},
		{
			name:        "rejects whitespace-only name",
			id:          "project-1",
			projectName: "   ",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := NewProject(tt.id, tt.projectName)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}

				if project != nil {
					t.Fatal("expected nil project")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if project == nil {
				t.Fatal("expected project, got nil")
			}

			if project.ID() != tt.id {
				t.Errorf(
					"ID() = %q, want %q",
					project.ID(),
					tt.id,
				)
			}

			if project.Name() != tt.wantName {
				t.Errorf(
					"Name() = %q, want %q",
					project.Name(),
					tt.wantName,
				)
			}
		})
	}
}
