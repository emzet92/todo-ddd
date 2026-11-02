package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/domain"
)

type ProjectHandler struct {
	service *application.ProjectService
}

func NewProjectHandler(
	service *application.ProjectService,
) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /projects", h.create)
	mux.HandleFunc("GET /projects/{id}", h.get)
}

type createProjectRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *ProjectHandler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req createProjectRequest

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	id, err := domain.NewProjectID(req.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	name, err := domain.NewProjectName(req.Name)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project name")
		return
	}

	project, err := h.service.CreateProject(r.Context(), id, name)
	if errors.Is(err, application.ErrProjectAlreadyExists) {
		writeError(w, http.StatusConflict, "project already exists")
		return
	}
	if err != nil {
		if errors.Is(err, domain.ErrInvalidProjectID) || errors.Is(err, domain.ErrInvalidProjectName) {
			writeError(w, http.StatusBadRequest, "invalid project")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create project")
		return
	}

	writeProject(w, http.StatusCreated, project)
}

func (h *ProjectHandler) get(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := domain.NewProjectID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	project, err := h.service.GetProject(r.Context(), id)

	if errors.Is(err, application.ErrProjectNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeProject(w, http.StatusOK, project)
}

func writeProject(
	w http.ResponseWriter,
	status int,
	project *domain.Project,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(projectResponse{
		ID:   project.ID().String(),
		Name: project.Name().String(),
	})
}
