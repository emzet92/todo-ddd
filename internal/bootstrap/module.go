package bootstrap

import (
	"net/http"

	"example.com/todo-ddd/internal/projects"
)

type Module struct {
	Projects    *projects.Module
	HTTPHandler http.Handler
}

func Init() *Module {
	mux := http.NewServeMux()
	projectsModule := projects.NewModule()
	projectsModule.RegisterRoutes(mux)

	return &Module{
		Projects:    projectsModule,
		HTTPHandler: mux,
	}
}
