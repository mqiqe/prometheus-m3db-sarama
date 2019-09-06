GO              ?= GO15VENDOREXPERIMENT=1 go
GOPATH          := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
PROMU           ?= $(GOPATH)/bin/promu
GOLINTER        ?= $(GOPATH)/bin/gometalinter
pkgs            = $(shell $(GO) list ./... | grep -v /vendor/)
TARGET          ?= saramam3db

PREFIX          ?= $(shell pwd)
BIN_DIR         ?= $(shell pwd)

all: clean format build test

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

gometalinter: $(GOLINTER)
	@echo ">> linting code"
	@$(GOLINTER) --install --update > /dev/null
	@$(GOLINTER) --config=./.gometalinter.json ./...

build:
	@echo ">> building binaries"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o saramam3db

clean:
	@echo ">> Cleaning up"
	@find . -type f -name '*~' -exec rm -fv {} \;
	@rm -fv $(TARGET)


.PHONY: all format  build test promu clean  lint