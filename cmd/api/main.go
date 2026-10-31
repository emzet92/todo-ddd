package main

import (
	"context"
	"log"

	"example.com/todo-ddd/internal/bootstrap"
)

func main() {
	module := bootstrap.Init()

	project, err := module.Projects.CreateProject(
		context.Background(),
		"project-1",
		"My First Project",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created project: %s", project.Name())
}
