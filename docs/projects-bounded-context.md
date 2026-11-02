# Todo DDD — Projects bounded context (Go)

Przykład **jednego bounded contextu `projects`** z dwoma agregatami (`Project`, `Task`), walidowanymi Value Objectami, serwisami aplikacyjnymi, repozytoriami in-memory i REST API na standardowym `net/http`. Bez Springa, DI containerów i bibliotek zewnętrznych.

## Struktura

```text
internal/
  projects/                         # bounded context
    domain/                         # Project, Task, ProjectID, ProjectName, TaskID, TaskTitle
    application/                    # ProjectService, TaskService + interfejsy repozytoriów
    infrastructure/
      memory/                       # map + sync.RWMutex
      httpapi/                      # ProjectHandler, TaskHandler, błędy JSON
    module.go                       # lokalny composition root + RegisterRoutes
  bootstrap/
    module.go                       # złożenie modułów i jeden http.ServeMux
cmd/api/main.go                     # start HTTP i pomiar startu
```

### Zależności / DI

```go
// internal/projects/module.go
func NewModule() *Module {
    projects := memory.NewProjectRepository()
    tasks := memory.NewTaskRepository()
    return &Module{
        Projects: application.NewProjectService(projects),
        Tasks: application.NewTaskService(projects, tasks),
    }
}
```

`bootstrap.Init()` tworzy router i woła `projectsModule.RegisterRoutes(mux)`. Kolejny bounded context dodasz przez kolejny `NewModule()` i `RegisterRoutes(mux)`. Bez globalnych singletonów i automatycznego skanowania.

## Domena

- `ProjectID`: `prj_` + 8–32 małe znaki ASCII `a-z`/`0-9`.
- `ProjectName`: 1–100 znaków Unicode po trimowaniu, bez znaków sterujących.
- `TaskID`: `tsk_` + 8–32 małe znaki ASCII `a-z`/`0-9`.
- `TaskTitle`: 1–200 znaków Unicode po trimowaniu, bez znaków sterujących.
- Każdy Value Object to `struct` z nieeksportowanym polem, konstruktorem `New...` i metodami `String()`/`IsZero()`.
- `Task` przechowuje `ProjectID`, nie wskaźnik do `Project` — dwa osobne agregaty w jednym kontekście.
- Przejścia: `PLANNED → IN_PROGRESS → DONE`, a także `IN_PROGRESS → PLANNED`. `DONE` jest terminalne.

## API

Uruchom: `go run ./cmd/api` (port `8080`).

| Metoda | Ścieżka | Działanie |
|---|---|---|
| POST | `/projects` | Tworzy projekt (409 przy duplikacie ID) |
| GET | `/projects/{id}` | Pobiera projekt |
| POST | `/projects/{projectID}/tasks` | Tworzy task w istniejącym projekcie |
| GET | `/projects/{projectID}/tasks` | Lista tasków projektu (po ID) |
| GET | `/projects/{projectID}/tasks/{taskID}` | Pobiera task |
| POST | `/projects/{projectID}/tasks/{taskID}/start` | `PLANNED → IN_PROGRESS` |
| POST | `/projects/{projectID}/tasks/{taskID}/complete` | `IN_PROGRESS → DONE` |
| POST | `/projects/{projectID}/tasks/{taskID}/plan` | `IN_PROGRESS → PLANNED` |

Błędy: JSON `{"error":"..."}`, HTTP 400 (walidacja), 404 (brak zasobu), 409 (duplikat / niedozwolone przejście), 500 (inne).

```bash
curl -i -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"id":"prj_a1b2c3d4","name":"Mój projekt"}'

curl -i -X POST http://localhost:8080/projects/prj_a1b2c3d4/tasks \
  -H 'Content-Type: application/json' \
  -d '{"id":"tsk_a1b2c3d4","title":"Pierwsze zadanie"}'

curl -i -X POST http://localhost:8080/projects/prj_a1b2c3d4/tasks/tsk_a1b2c3d4/start
curl -i -X POST http://localhost:8080/projects/prj_a1b2c3d4/tasks/tsk_a1b2c3d4/complete
curl -i http://localhost:8080/projects/prj_a1b2c3d4/tasks
```

## Testy i budowanie

```bash
go test ./...
go test -race ./...
go vet ./...
go build -o bin/todo-api ./cmd/api
```

`go.mod` zachowuje Twoją wersję Go z poprzedniego projektu. Implementacja nie dodaje zależności spoza stdlib.

## Ograniczenia na obecnym etapie

- Dane są w pamięci procesu i znikają po restarcie. Dwa `bootstrap.Init()` dają osobne magazyny.
- `TaskService` zmienia status przez `FindByID → metoda domenowa → Save`. Mapa jest chroniona mutexem (brak data races), ale **konkurencyjne przejścia tego samego zadania nie są atomowe**. Przy bazie danych dodamy optimistic locking/version check lub transakcję.
- CreateTask sprawdza istnienie Project, ale pełna atomowość referencyjna (np. równoczesne usunięcie projektu) wymaga docelowej bazy danych i transakcji.
- `ProjectID{}`/`TaskID{}` to legalne zero values w Go; konstruktory encji je odrzucają.
- Dla czytelności backend nadal przyjmuje ID od klienta, a nie generuje UUID — to może być kolejny krok.
