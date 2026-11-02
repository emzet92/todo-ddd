package bootstrap_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/todo-ddd/internal/bootstrap"
)

func TestProjectHTTP(t *testing.T) {
	module := bootstrap.Init()

	body := strings.NewReader(`{
        "id": "project-1",
        "name": "My Project"
    }`)

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/projects",
		body,
	)

	createResponse := httptest.NewRecorder()
	module.HTTPHandler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResponse.Code)
	}

	getRequest := httptest.NewRequest(
		http.MethodGet,
		"/projects/project-1",
		nil,
	)

	getResponse := httptest.NewRecorder()
	module.HTTPHandler.ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResponse.Code)
	}

	if !strings.Contains(getResponse.Body.String(), "My Project") {
		t.Fatal("expected project in response")
	}
}
