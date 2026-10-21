GOBIN=$(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
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
		$(GOBIN)/mockgen -source=$$file -destination="$$dir/mocks/mocks.go" -package=mocks; \
	done

	@echo "==> Mock generation completed successfully"

.PHONY: lint
lint:
	@echo "===> Installing golangci-lint"
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0

	@echo "===> Running Go Linter"
	@for dir in $$(go list -f '{{.Dir}}' -m); do (cd $$dir && $(GOBIN)/golangci-lint run --timeout 10m ./...) || exit 1; done

	@echo "===> Installing actionlint"
	@go install github.com/rhysd/actionlint/cmd/actionlint@v1.7.12

	@echo "===> Running Github Actions Linter"
	@$(GOBIN)/actionlint -shellcheck=

	@echo "===> No lint issues found"

.PHONY: vuln-check
vuln-check:
	@echo "===> Installing vulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@v1.9.0

	@echo "===> Checking vulnerabilities"
	@for dir in $$(go list -f '{{.Dir}}' -m); do $(GOBIN)/govulncheck -C $$dir ./... || exit 1; done

	@echo "===> No vulnerabilities found"

.PHONY: test.unit
test.unit:
	@echo "===> Running Unit Tests"
	@go list -f '{{.Dir}}/...' -m | xargs go test -short -cover
	@echo "===> Unit Tests passed"

# Kafka's component tests stay out until CI runs a broker.
.PHONY: test.component
test.component:
	@echo "===> Running Component Tests"
	@go test -run Component ./postgres/...
	@echo "===> Component Tests passed"
