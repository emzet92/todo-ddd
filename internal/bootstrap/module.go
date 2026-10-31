package bootstrap

import (
	"example.com/todo-ddd/internal/application"
	"example.com/todo-ddd/internal/infrastructure/memory"
)

type Module struct {
	Projects *application.ProjectService
}

func Init() *Module {
	projectRepo := memory.NewProjectRepository()

	return &Module{
		Projects: application.NewProjectService(projectRepo),
	}
}
