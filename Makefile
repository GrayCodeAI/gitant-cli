.PHONY: build test install clean vet fmt

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -ldflags "-s -w \
	-X github.com/GrayCodeAI/gitant-cli/internal/version.Version=$(VERSION) \
	-X github.com/GrayCodeAI/gitant-cli/internal/version.Commit=$(COMMIT) \
	-X github.com/GrayCodeAI/gitant-cli/internal/version.BuildTime=$(BUILD_TIME)"

## build: Compile gt (gitant) and git-remote-gitant
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/gt ./cmd/gitant/
	cp bin/gt bin/gitant
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/git-remote-gitant ./cmd/git-remote-gitant/

## install: Install to $GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/gitant/
	go install $(LDFLAGS) ./cmd/git-remote-gitant/

## test: Run tests
test:
	go test ./... -race -count=1 -timeout=120s

## vet: Run go vet
vet:
	go vet ./...

## fmt: List non-gofmt files
fmt:
	gofmt -l .

## clean: Remove build artifacts
clean:
	rm -rf bin/
