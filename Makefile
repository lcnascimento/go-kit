GOPATH=$(shell go env GOPATH)

.PHONY: lint
lint:
	@echo "===> Installing golangci-lint"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1

	@echo "===> Running Go Linter"
	@go list -f '{{.Dir}}' -m | xargs $(GOPATH)/bin/golangci-lint run --timeout 10m

	@echo "===> Installing actionlint"
	@go install github.com/rhysd/actionlint/cmd/actionlint@v1.6.26

	@echo "===> Running Github Actions Linter"
	@$(GOPATH)/bin/actionlint -shellcheck=

	@echo "===> No lint issues found"

.PHONY: vuln-check
vuln-check:
	@echo "===> Installing vulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest

	@echo "===> Checking vulnerabilities"
	@go list -m | $(GOPATH)/bin/govulncheck

	@echo "===> No vulnerabilities found"

.PHONY: test.unit
test.unit:
	@echo "===> Running Unit Tests"
	@go list -m | xargs go test -cover
	@echo "===> Unit Tests passed"
