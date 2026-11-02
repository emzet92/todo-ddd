// Package projects is the composition root of the Projects bounded context.
package projects

import (
	"net/http"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/infrastructure/httpapi"
	"example.com/todo-ddd/internal/projects/infrastructure/memory"
)

// Module exposes use cases, not repositories or implementation details.
type Module struct {
	Projects *application.ProjectService
	Tasks    *application.TaskService
}

func NewModule() *Module {
	projects := memory.NewProjectRepository()
	tasks := memory.NewTaskRepository()
	return &Module{
		Projects: application.NewProjectService(projects),
		Tasks:    application.NewTaskService(projects, tasks),
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	httpapi.NewProjectHandler(m.Projects).RegisterRoutes(mux)
	httpapi.NewTaskHandler(m.Tasks).RegisterRoutes(mux)
}
