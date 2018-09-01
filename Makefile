SOURCES = $(wildcard *.go) \
          $(wildcard */*/*.go)


.PHONY: get build generate test clean format

.DEFAULT_GOAL := all

all: build test

get: format
	@echo get gocui
	@go get -u github.com/askovpen/gocui
	@echo get transform
	@go get -u  golang.org/x/text/transform
	@echo get yaml
	@go get -u  gopkg.in/yaml.v2
	@echo get goblin
	@go get -u  github.com/franela/goblin

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