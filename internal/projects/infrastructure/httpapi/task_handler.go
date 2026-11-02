package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"example.com/todo-ddd/internal/projects/application"
	"example.com/todo-ddd/internal/projects/domain"
)

type TaskHandler struct {
	service *application.TaskService
}

func NewTaskHandler(service *application.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /projects/{projectID}/tasks", h.create)
	mux.HandleFunc("GET /projects/{projectID}/tasks", h.list)
	mux.HandleFunc("GET /projects/{projectID}/tasks/{taskID}", h.get)
	mux.HandleFunc("POST /projects/{projectID}/tasks/{taskID}/start", h.start)
	mux.HandleFunc("POST /projects/{projectID}/tasks/{taskID}/complete", h.complete)
	mux.HandleFunc("POST /projects/{projectID}/tasks/{taskID}/plan", h.plan)
}

type createTaskRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type taskResponse struct {
	ID        string            `json:"id"`
	ProjectID string            `json:"project_id"`
	Title     string            `json:"title"`
	Status    domain.TaskStatus `json:"status"`
}

func taskToResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID: t.ID().String(), ProjectID: t.ProjectID().String(),
		Title: t.Title().String(), Status: t.Status(),
	}
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func projectIDFromPath(w http.ResponseWriter, r *http.Request) (domain.ProjectID, bool) {
	id, err := domain.NewProjectID(r.PathValue("projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return domain.ProjectID{}, false
	}
	return id, true
}

func taskIDFromPath(w http.ResponseWriter, r *http.Request) (domain.TaskID, bool) {
	id, err := domain.NewTaskID(r.PathValue("taskID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return domain.TaskID{}, false
	}
	return id, true
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	projectID, ok := projectIDFromPath(w, r)
	if !ok {
		return
	}

	var req createTaskRequest
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

	id, err := domain.NewTaskID(req.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}
	title, err := domain.NewTaskTitle(req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task title")
		return
	}
	task, err := h.service.CreateTask(r.Context(), id, projectID, title)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, taskToResponse(task))
}

func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	projectID, ok := projectIDFromPath(w, r)
	if !ok {
		return
	}
	id, ok := taskIDFromPath(w, r)
	if !ok {
		return
	}
	task, err := h.service.GetTask(r.Context(), projectID, id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, taskToResponse(task))
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	projectID, ok := projectIDFromPath(w, r)
	if !ok {
		return
	}
	tasks, err := h.service.ListTasks(r.Context(), projectID)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	response := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		response = append(response, taskToResponse(t))
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *TaskHandler) start(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.service.StartTask)
}
func (h *TaskHandler) complete(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.service.CompleteTask)
}
func (h *TaskHandler) plan(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.service.PlanTask)
}

func (h *TaskHandler) transition(w http.ResponseWriter, r *http.Request, action func(ctx context.Context, projectID domain.ProjectID, id domain.TaskID) (*domain.Task, error)) {
	projectID, ok := projectIDFromPath(w, r)
	if !ok {
		return
	}
	id, ok := taskIDFromPath(w, r)
	if !ok {
		return
	}
	task, err := action(r.Context(), projectID, id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, taskToResponse(task))
}

func writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrProjectNotFound):
		writeError(w, http.StatusNotFound, "project not found")
	case errors.Is(err, application.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "task not found")
	case errors.Is(err, application.ErrTaskAlreadyExists):
		writeError(w, http.StatusConflict, "task already exists")
	case errors.Is(err, domain.ErrInvalidTaskTransition):
		writeError(w, http.StatusConflict, "invalid task status transition")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
