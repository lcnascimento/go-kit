GOPATH=$(shell go env GOPATH)
CONTRACTS=$(shell find . -name "contracts.go")

.PHONY: setup
setup: mocks lint vuln-check
	@echo "===> Installing Node dev dependencies"
	@npm install

	@echo "===> Installing husky"
	@npx husky install

	@echo "===> Setup concluded"

.PHONY: mocks
mocks:
	@echo "==> Installing mockgen"
	@go install go.uber.org/mock/mockgen@v0.4.0

	@echo "==> Generating mocks"
	@for file in $(CONTRACTS); do \
		dir=$$(dirname $$file); \
		rm -Rf "$$dir/mocks"; \
		mkdir -p "$$dir/mocks"; \
		$(GOPATH)/bin/mockgen -source=$$file -destination="$$dir/mocks/$$(basename $$file)" -package=mocks; \
	done

	@echo "==> Mock generation completed successfully"

.PHONY: lint
lint:
	@echo "===> Installing golangci-lint"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

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
