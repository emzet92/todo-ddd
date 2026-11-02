package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/todo-ddd/internal/bootstrap"
)

func TestProjectHTTP(t *testing.T) {
	module := bootstrap.Init()

	body := strings.NewReader(`{
        "id": "prj_a1b2c3d4",
        "name": "  My Project  "
    }`)

	createRequest := httptest.NewRequest(http.MethodPost, "/projects", body)
	createResponse := httptest.NewRecorder()
	module.HTTPHandler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", createResponse.Code, createResponse.Body.String())
	}

	var created struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID != "prj_a1b2c3d4" || created.Name != "My Project" {
		t.Fatalf("unexpected created project: %+v", created)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/projects/prj_a1b2c3d4", nil)
	getResponse := httptest.NewRecorder()
	module.HTTPHandler.ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", getResponse.Code, getResponse.Body.String())
	}
	if !strings.Contains(getResponse.Body.String(), "My Project") {
		t.Fatal("expected project in response")
	}
}

func TestProjectHTTPValidationErrorsAreJSON(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "invalid ID on create",
			method:     http.MethodPost,
			path:       "/projects",
			body:       `{"id":"project-1","name":"My Project"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid project ID",
		},
		{
			name:       "invalid name on create",
			method:     http.MethodPost,
			path:       "/projects",
			body:       `{"id":"prj_a1b2c3d4","name":"    "}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid project name",
		},
		{
			name:       "invalid ID on get",
			method:     http.MethodGet,
			path:       "/projects/invalid",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid project ID",
		},
		{
			name:       "missing project",
			method:     http.MethodGet,
			path:       "/projects/prj_zzzzzzzz",
			wantStatus: http.StatusNotFound,
			wantError:  "project not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module := bootstrap.Init()
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			response := httptest.NewRecorder()
			module.HTTPHandler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			var payload struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.Error != tt.wantError {
				t.Fatalf("error = %q, want %q", payload.Error, tt.wantError)
			}
		})
	}
}
