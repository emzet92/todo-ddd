.PHONY: test test-report

test:
	go test ./...

test-report:
	@bash ./scripts/test-report.sh
