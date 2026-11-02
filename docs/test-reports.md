# Go test reports

This repository uses [go-test-html-report](https://github.com/brpaz/go-test-html-report)
for a readable HTML test report, plus the built-in Go coverage viewer.
No testing library or reporting dependency is added to the application `go.mod`.

## Generate reports

From the repository root:

```bash
make test-report
```

On the first run, the script downloads the reporting command with `go install`
into `reports/.bin/` (requires internet access and Go). Later runs reuse it.
You can still run plain tests using `make test` or `go test ./...`.

Open the results on macOS:

```bash
open reports/tests.html
open reports/coverage.html
```

The generated files are:

- `reports/tests.html` — test runs, durations and failure output.
- `reports/coverage.html` — source-level coverage, created with `go tool cover`.
- `reports/tests.json` — machine-readable Go test events (`go test -json`).
- `reports/coverage.out` — machine-readable Go coverage profile.
- `reports/tests.stderr` — compiler/test runner diagnostics, if present.

**Important:** `make test-report` exits unsuccessfully if any Go test fails,
even when it successfully creates an HTML report. It runs tests with
`-count=1` to avoid cached test results. Coverage percentage is useful, but
high coverage does not by itself prove good test quality.

The `reports/` directory ignores generated output through its own `.gitignore`.
The script does not change the application runtime or source imports.
