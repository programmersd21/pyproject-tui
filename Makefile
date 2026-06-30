# Detect OS for cross-platform commands
ifeq ($(OS),Windows_NT)
    RM_FILE   = cmd /c del /q /f 2>nul
    RM_DIR    = cmd /c rd /s /q 2>nul
    MKDIR     = cmd /c mkdir 2>nul
    NULL      = nul
    DATE_CMD  = powershell -Command "[DateTime]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ssZ')"
else
    RM_FILE   = rm -f
    RM_DIR    = rm -rf
    MKDIR     = mkdir -p
    NULL      = /dev/null
    DATE_CMD  = date -u +%Y-%m-%dT%H:%M:%SZ
endif

BINARY     := pyproject-tui
MODULE     := github.com/programmersd21/pyproject-tui
VERSION    := $(strip $(file < VERSION))
ifeq ($(VERSION),)
    VERSION := dev
endif
COMMIT     := $(shell git rev-parse --short HEAD 2>$(NULL) || echo none)
DATE       := $(shell $(DATE_CMD) 2>$(NULL) || echo unknown)
LDFLAGS    := -s -w \
              -X '$(MODULE)/cmd/pyproject-tui.version=$(VERSION)' \
              -X '$(MODULE)/cmd/pyproject-tui.commit=$(COMMIT)' \
              -X '$(MODULE)/cmd/pyproject-tui.date=$(DATE)'

GO         := go
GOFLAGS    := -trimpath
DIST       := dist

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: $(DIST) ## Build binary to dist/
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(DIST)/$(BINARY) ./cmd/pyproject-tui

$(DIST):
	$(MKDIR) $(DIST)

.PHONY: install
install: ## Install binary to GOPATH/bin
	$(GO) install $(GOFLAGS) -ldflags="$(LDFLAGS)" ./cmd/pyproject-tui

.PHONY: run
run: ## Build and run with ./pyproject.toml
	$(GO) run ./cmd/pyproject-tui

.PHONY: test
test: ## Run tests with race detector and coverage
	$(GO) test -v -race -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: test ## Open coverage report in browser
	$(GO) tool cover -html=coverage.out

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix ./...

.PHONY: fmt
fmt: ## Format all Go source files
	gofmt -s -w .
	$(GO) mod tidy

.PHONY: version
version: ## Print the repository version
	@echo $(VERSION)

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy and verify go.mod
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: clean
clean: ## Remove build artifacts
	$(RM_DIR) $(DIST)
	$(RM_FILE) coverage.out
	$(RM_FILE) coverage.html

.PHONY: check
check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
