package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/todo-ddd/internal/bootstrap"
)

func call(t *testing.T, handler http.Handler, method, path, body string, status int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != status {
		t.Fatalf("%s %s: status %d, want %d, body %s", method, path, resp.Code, status, resp.Body.String())
	}
	if got := resp.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type: %s", got)
	}
	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	return payload
}

func TestTaskHTTPLifecycle(t *testing.T) {
	m := bootstrap.Init()
	base := "/projects/prj_abcd1234/tasks"
	call(t, m.HTTPHandler, http.MethodPost, "/projects", `{"id":"prj_abcd1234","name":"  Roadmap  "}`, 201)
	created := call(t, m.HTTPHandler, http.MethodPost, base, `{"id":"tsk_abcd1234","title":"  First task  "}`, 201)
	if created["status"] != "PLANNED" || created["title"] != "First task" || created["project_id"] != "prj_abcd1234" {
		t.Fatalf("created: %+v", created)
	}
	item := base + "/tsk_abcd1234"
	read := call(t, m.HTTPHandler, http.MethodGet, item, "", 200)
	if read["status"] != "PLANNED" {
		t.Fatalf("read: %+v", read)
	}
	call(t, m.HTTPHandler, http.MethodPost, item+"/complete", "", 409)
	started := call(t, m.HTTPHandler, http.MethodPost, item+"/start", "", 200)
	if started["status"] != "IN_PROGRESS" {
		t.Fatalf("started: %+v", started)
	}
	planned := call(t, m.HTTPHandler, http.MethodPost, item+"/plan", "", 200)
	if planned["status"] != "PLANNED" {
		t.Fatalf("planned: %+v", planned)
	}
	call(t, m.HTTPHandler, http.MethodPost, item+"/start", "", 200)
	completed := call(t, m.HTTPHandler, http.MethodPost, item+"/complete", "", 200)
	if completed["status"] != "DONE" {
		t.Fatalf("completed: %+v", completed)
	}
	call(t, m.HTTPHandler, http.MethodPost, item+"/start", "", 409)
	// Verify repository persisted status.
	read = call(t, m.HTTPHandler, http.MethodGet, item, "", 200)
	if read["status"] != "DONE" {
		t.Fatalf("read: %+v", read)
	}
}

func TestTaskHTTPList(t *testing.T) {
	m := bootstrap.Init()
	base := "/projects/prj_abcd1234/tasks"
	call(t, m.HTTPHandler, "POST", "/projects", `{"id":"prj_abcd1234","name":"Roadmap"}`, 201)
	call(t, m.HTTPHandler, "POST", base, `{"id":"tsk_zzzzzzzz","title":"Task Z"}`, 201)
	call(t, m.HTTPHandler, "POST", base, `{"id":"tsk_aaaaaaaa","title":"Task A"}`, 201)
	req := httptest.NewRequest("GET", base, nil)
	resp := httptest.NewRecorder()
	m.HTTPHandler.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatalf("list: %d: %s", resp.Code, resp.Body.String())
	}
	var list []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != "tsk_aaaaaaaa" || list[1].ID != "tsk_zzzzzzzz" {
		t.Fatalf("list: %+v", list)
	}
}

func TestTaskHTTPValidationAndNotFound(t *testing.T) {
	m := bootstrap.Init()
	base := "/projects/prj_abcd1234/tasks"
	call(t, m.HTTPHandler, "POST", "/projects", `{"id":"prj_abcd1234","name":"Roadmap"}`, 201)
	for _, tc := range []struct {
		method, path, body string
		status             int
		message            string
	}{
		{"POST", base, `{"id":"bad","title":"Something"}`, 400, "invalid task ID"},
		{"POST", base, `{"id":"tsk_abcd1234","title":" "}`, 400, "invalid task title"},
		{"POST", base, `{"id":"tsk_abcd1234","title":"OK","unknown":true}`, 400, "invalid JSON"},
		{"POST", base, `{"id":"tsk_abcd1234","title":"OK"} {}`, 400, "invalid JSON"},
		{"POST", "/projects/prj_zzzzzzzz/tasks", `{"id":"tsk_abcd1234","title":"OK"}`, 404, "project not found"},
		{"GET", base + "/tsk_abcd1234", "", 404, "task not found"},
		{"GET", "/projects/bad/tasks", "", 400, "invalid project ID"},
		{"GET", base + "/bad", "", 400, "invalid task ID"},
	} {
		got := call(t, m.HTTPHandler, tc.method, tc.path, tc.body, tc.status)
		if got["error"] != tc.message {
			t.Fatalf("%s %s: %+v", tc.method, tc.path, got)
		}
	}
	call(t, m.HTTPHandler, "POST", base, `{"id":"tsk_abcd1234","title":"OK"}`, 201)
	call(t, m.HTTPHandler, "POST", base, `{"id":"tsk_abcd1234","title":"OK"}`, 409)
	call(t, m.HTTPHandler, "POST", "/projects", `{"id":"prj_abcd1234","name":"Roadmap"}`, 409)
	call(t, m.HTTPHandler, "GET", "/projects/prj_zzzzzzzz/tasks/tsk_abcd1234", "", 404)
}

func TestBootstrapInitIsIndependent(t *testing.T) {
	one, two := bootstrap.Init(), bootstrap.Init()
	call(t, one.HTTPHandler, "POST", "/projects", `{"id":"prj_abcd1234","name":"First"}`, 201)
	call(t, two.HTTPHandler, "GET", "/projects/prj_abcd1234", "", 404)
}
