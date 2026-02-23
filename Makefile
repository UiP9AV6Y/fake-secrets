include .bingo/Variables.mk
include .changes/changelog.mk

GOLANGCI_LINT_VERSION ?= $(shell $(GOLANGCI_LINT) version --short)
GOLANG_VERSION ?= $(shell $(GO) env GOVERSION)
BINGO_VERSION ?= $(shell $(BINGO) version)

PREFIX ?= /usr/local
bin_PROGRAM = $(notdir $(wildcard cmd/*))

INSTALL ?= install

.PHONY: default
default: all

.PHONY: all
all: lint fmt test build

.tool-versions: $(GOLANGCI_LINT) $(BINGO)
	@echo golang $(GOLANG_VERSION:go%=%) > $@
	@echo golangci-lint $(GOLANGCI_LINT_VERSION) >> $@
	@echo bingo $(BINGO_VERSION:v%=%) >> $@

.PHONY: tools
tools: $(BINGO)
	$(BINGO) get -l

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

.PHONY: dependencies
dependencies:
	$(GO) mod download -x

.PHONY: install
install:
	$(INSTALL) -m 0755 -D -t $(DESTDIR)$(PREFIX)/bin $(addprefix $(GOBIN)/, $(bin_PROGRAM))

CHANGELOG.md: $(CHANGE_DIR)/CHANGELOG-v.md.tmp
	sed -n -e '/## /q;p' $@ > $@.tmp
	cat $^ >> $@.tmp
	sed -n -e '/## /,$$p' $@ >> $@.tmp
	mv $@.tmp $@
