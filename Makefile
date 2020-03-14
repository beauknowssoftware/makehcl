export SHELL=/bin/bash
ifdef DEBUG
export .SHELLFLAGS = -xuec
else
export .SHELLFLAGS = -uec
endif
.EXPORT_ALL_VARIABLES:
.ONESHELL:
.SILENT:
.SUFFIXES:
.DEFAULT_GOAL := build

GOFILES = $(shell find . -name '*.go')

bin/%: cmd/% $(GOFILES) go.mod go.sum test lint
	go build -o $@ ./$<

.PHONY: import
import: ; goimports -w .

.PHONY: test
test: import ; go test ./...

.PHONY: build
build: bin/makehcl

.PHONY: lint
lint: ; golangci-lint run

.PHONY: clean
clean: ; git clean -f -fdX
