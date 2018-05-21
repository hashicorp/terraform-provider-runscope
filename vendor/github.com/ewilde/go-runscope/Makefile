SOURCEDIR=.
SOURCES = $(shell find $(SOURCEDIR) -name '*.go')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
VERSION=$(git describe --always --tags)
BINARY=bin/runscope

bin: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY) command/*

build:
	go get github.com/golang/lint/golint
	go test $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	golint $(go list ./... | grep -v /vendor/)

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test -v ./...

.PHONY: build fmt test
