VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_USER ?= $(shell whoami)@$(shell hostname)
PKG := github.com/user-cube/kontext
LDFLAGS := -ldflags "-X $(PKG)/cmd.Version=$(VERSION) -X $(PKG)/cmd.BuildDate=$(BUILD_DATE) -X $(PKG)/cmd.GitCommit=$(COMMIT) -X $(PKG)/cmd.BuildUser=$(BUILD_USER)"

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "Building kontext $(VERSION) ($(COMMIT))"
	@go build $(LDFLAGS) -o kontext main.go

.PHONY: install
install:
	@echo "Installing kontext $(VERSION) to GOPATH"
	@go install $(LDFLAGS)

.PHONY: clean
clean:
	@echo "Cleaning build artifacts"
	@rm -f kontext
	@rm -rf dist

.PHONY: test
test:
	@echo "Running tests"
	@go test -v ./...

.PHONY: lint
lint:
	@echo "Running linters"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi

.PHONY: release
release: clean
	@echo "Creating release with GoReleaser $(VERSION)"
	@if ! command -v goreleaser > /dev/null; then \
		echo "Error: goreleaser not found. Install with 'go install github.com/goreleaser/goreleaser@latest'"; \
		exit 1; \
	fi
	@VERSION=$(VERSION) GIT_COMMIT=$(COMMIT) BUILD_DATE=$(BUILD_DATE) goreleaser release --clean

.PHONY: release-snapshot
release-snapshot: clean
	@echo "Creating snapshot release with GoReleaser (no publish)"
	@if ! command -v goreleaser > /dev/null; then \
		echo "Error: goreleaser not found. Install with 'go install github.com/goreleaser/goreleaser@latest'"; \
		exit 1; \
	fi
	@VERSION=$(VERSION) GIT_COMMIT=$(COMMIT) BUILD_DATE=$(BUILD_DATE) goreleaser release --snapshot --clean

.PHONY: build-release
build-release: clean
	@echo "Building release version of kontext $(VERSION) ($(COMMIT))"
	@go build $(LDFLAGS) -o kontext main.go
	@echo "Built kontext binary with release information"
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Run ./kontext version to verify"

.PHONY: help
help:
	@echo "Kontext Makefile"
	@echo "---------------"
	@echo "Available targets:"
	@echo "  all              - Clean and build kontext"
	@echo "  build            - Build the kontext binary"
	@echo "  install          - Install kontext to your GOPATH/bin"
	@echo "  clean            - Remove built binary and dist directory"
	@echo "  test             - Run tests"
	@echo "  lint             - Run linters (requires golangci-lint)"
	@echo "  release          - Create a full release using GoReleaser"
	@echo "  release-snapshot - Create a local release snapshot for testing (no publish)"
	@echo "  build-release    - Build kontext binary with release information"
	@echo "  help             - Show this help message"