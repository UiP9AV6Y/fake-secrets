include .bingo/Variables.mk

GOLANGCI_LINT_VERSION ?= $(shell $(GOLANGCI_LINT) version --short)
GOLANG_VERSION ?= $(shell $(GO) env GOVERSION)
BINGO_VERSION ?= $(shell $(BINGO) version)

.PHONY: default
default: all

.PHONY: all
all: lint fmt test build

.tool-versions: $(GOLANGCI_LINT) $(BINGO)
	@echo golang $(GOLANG_VERSION:go%=%) > $@
	@echo golangci-lint $(GOLANGCI_LINT_VERSION) >> $@
	@echo bingo $(BINGO_VERSION:v%=%) >> $@

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: fmt
fmt: $(GOIMPORTS)
	$(GOIMPORTS) -w .

.PHONY: test
test:
	$(GO) test -v -race -timeout 20m ./...

.PHONY: build
build: generate
	$(GO) install -ldflags="-s -w" ./...

.PHONY: generate
generate:
	$(GO) generate ./...
