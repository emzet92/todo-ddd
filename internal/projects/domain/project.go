package domain

type Project struct {
	id   ProjectID
	name ProjectName
}

func NewProject(id ProjectID, name ProjectName) (*Project, error) {
	// Go allows zero-value structs even when fields are unexported.
	if id.IsZero() {
		return nil, ErrInvalidProjectID
	}

	if name.IsZero() {
		return nil, ErrInvalidProjectName
	}

	return &Project{id: id, name: name}, nil
}

func (p *Project) ID() ProjectID {
	return p.id
}

func (p *Project) Name() ProjectName {
	return p.name
}
