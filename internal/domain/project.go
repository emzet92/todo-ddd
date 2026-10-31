package domain

import (
	"errors"
	"strings"
)

type ProjectID string
type ProjectName string

type Project struct {
	id   ProjectID
	name ProjectName
}

func NewProject(id ProjectID, name ProjectName) (*Project, error) {
	name = ProjectName(strings.TrimSpace(string(name)))

	if id == "" {
		return nil, errors.New("project ID cannot be empty")
	}

	if name == "" {
		return nil, errors.New("project name cannot be empty")
	}

	return &Project{
		id:   id,
		name: name,
	}, nil
}

func (p *Project) ID() ProjectID {
	return p.id
}

func (p *Project) Name() ProjectName {
	return p.name
}
