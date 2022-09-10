SOURCES = $(wildcard *.go) \
          $(wildcard */*/*.go)


.PHONY: get build test clean format update

.DEFAULT_GOAL := all

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

all: format build test

get: format
	@echo get depencies

test:
	@echo Testing gossiped
	@go test -v -cover ./...

build:
	@echo Building gossiped
	@go build

clean:
	@echo Cleaning
	@go clean

format:
	@echo Formating sources
	@gofmt -s -w $(SOURCES)

update: format
	go mod tidy
