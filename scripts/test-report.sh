#!/usr/bin/env bash
set -euo pipefail

# Run from any directory; keep all generated artifacts outside tracked source files.
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_DIR="$ROOT_DIR/reports"
TOOL_BIN="$REPORT_DIR/.bin/go-test-html-report"

mkdir -p "$REPORT_DIR/.bin"

# Install the reporting CLI separately from the application module.
# This does not change go.mod or go.sum.
if [[ ! -x "$TOOL_BIN" ]]; then
    echo "Installing go-test-html-report (first run only)..."
    GOBIN="$REPORT_DIR/.bin" go install github.com/brpaz/go-test-html-report/cmd/go-test-html-report@latest
fi

cd "$ROOT_DIR"
echo "Running Go tests and collecting coverage..."

# Preserve the test exit status: generating an HTML report must never turn a failed
# test run into a green CI/build result.
set +e
go test -count=1 -json -coverprofile="$REPORT_DIR/coverage.out" ./... \
    > "$REPORT_DIR/tests.json" 2> "$REPORT_DIR/tests.stderr"
TEST_EXIT_CODE=$?
set -e

if [[ -s "$REPORT_DIR/tests.stderr" ]]; then
    cat "$REPORT_DIR/tests.stderr" >&2
fi

if ! "$TOOL_BIN" \
    -i "$REPORT_DIR/tests.json" \
    -o "$REPORT_DIR/tests.html" \
    -t "Todo DDD - Test Report"; then
    echo "Could not generate the HTML test report. Inspect reports/tests.json." >&2
    exit 1
fi

if [[ -s "$REPORT_DIR/coverage.out" ]]; then
    go tool cover -html="$REPORT_DIR/coverage.out" -o "$REPORT_DIR/coverage.html"
fi

printf '\nReports:\n  %s\n' "$REPORT_DIR/tests.html"
if [[ -f "$REPORT_DIR/coverage.html" ]]; then
    printf '  %s\n' "$REPORT_DIR/coverage.html"
fi

if (( TEST_EXIT_CODE != 0 )); then
    echo "Some tests failed (exit code: $TEST_EXIT_CODE). Check the report above." >&2
    exit "$TEST_EXIT_CODE"
fi

echo "All Go tests passed."
