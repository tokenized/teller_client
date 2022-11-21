BUILD_DATE = `date +%FT%T%z`
BUILD_USER = $(USER)@`hostname`
VERSION = `git describe --tags`

# command to build and run on the local OS.
GO_BUILD = go build

# tools
BINARY_CONTRACT_CLI=teller

GO_DIST = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -a -tags netgo -ldflags "-w -X main.buildVersion=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.buildUser=$(BUILD_USER)"

all: deps tools test dist-cli

deps:
	go get -t ./...

prepare:
	mkdir -p tmp

test: prepare
	go test -coverprofile=tmp/coverage.out ./...

test-race:
	go test -race ./...

dist-cli:
	mkdir -p dist
	$(GO_DIST) -o dist/$(BINARY_CONTRACT_CLI) cmd/$(BINARY_CONTRACT_CLI)/main.go
