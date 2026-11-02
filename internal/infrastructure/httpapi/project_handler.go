package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/todo-ddd/internal/application"
	"example.com/todo-ddd/internal/domain"
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.ID) == "" ||
		strings.TrimSpace(req.Name) == "" {
		http.Error(w, "ID and name are required", http.StatusBadRequest)
		return
	}

	project, err := h.service.CreateProject(
		r.Context(),
		domain.ProjectID(req.ID),
		domain.ProjectName(req.Name),
	)
	if err != nil {
		http.Error(w, "could not create project", http.StatusInternalServerError)
		return
	}

	writeProject(w, http.StatusCreated, project)
}

func (h *ProjectHandler) get(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := domain.ProjectID(r.PathValue("id"))

	project, err := h.service.GetProject(r.Context(), id)

	if errors.Is(err, application.ErrProjectNotFound) {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
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

	json.NewEncoder(w).Encode(projectResponse{
		ID:   string(project.ID()),
		Name: string(project.Name()),
	})
}
