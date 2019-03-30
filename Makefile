SOURCES = $(wildcard *.go) \
          $(wildcard */*/*.go)


.PHONY: get build generate test clean format update

.DEFAULT_GOAL := all

all: build test

get: format
	@echo get depencies
	@dep ensure

generate: get
	@echo Generating version.go
	@go generate ./...

test:
	@echo Testing goated
	@go test -v -cover ./...

build: generate
	@echo Building goated
	@go build

clean:
	@echo Cleaning
	@go clean

format:
	@echo Formating sources
	@gofmt -s -w $(SOURCES)

update: format
	dep ensure -update