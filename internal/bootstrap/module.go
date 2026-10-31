package bootstrap

import (
	"net/http"

	"example.com/todo-ddd/internal/application"
	"example.com/todo-ddd/internal/infrastructure/httpapi"
	"example.com/todo-ddd/internal/infrastructure/memory"
)

type Module struct {
	Projects    *application.ProjectService
	HTTPHandler http.Handler
}

func Init() *Module {
	// Repositories
	projectRepo := memory.NewProjectRepository()

	// Application services
	projectService := application.NewProjectService(projectRepo)

	// HTTP
	mux := http.NewServeMux()

	projectHandler := httpapi.NewProjectHandler(projectService)
	projectHandler.RegisterRoutes(mux)

	return &Module{
		Projects:    projectService,
		HTTPHandler: mux,
	}
}
