SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell basename ${CURDIR})
.DEFAULT_GOAL:=$(TARGET)


# These will be provided to the target
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.1)
BUILD := $(shell git rev-parse HEAD)

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Golang Flags
GOFLAGS ?= $(GOFLAGS:)
GO=go

PLATFORMS := linux/amd64 # linux/arm64 windows/amd64 darwin/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

#Set zip file name
BASE_ARCHIVE = binaries_*.tar
ARCHIVE = binaries_$(os).tar

WD = $(subst $(BSLASH),$(FLASH),$(shell pwd))
#Directories
DISTDIR = $(WD)/dist

.PHONY: all build clean install uninstall run release $(PLATFORMS)

release: clean $(PLATFORMS) _gzip

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o '$(DISTDIR)/bin/$(TARGET)_$(os)_$(arch)' $(GOFLAGS) $(GO_LINKER_FLAGS) $(LDFLAGS)
	@cd "$(DISTDIR)/bin/" && tar -uvf "../$(ARCHIVE)" "./$(TARGET)_$(os)_$(arch)"

_gzip:
	@cd "$(DISTDIR)" && gzip ./$(BASE_ARCHIVE)

$(TARGET): $(SRC)
	@$(GO) build $(GOFLAGS) $(GO_LINKER_FLAGS) $(LDFLAGS)

build: $(TARGET)
	@true

clean:
	@$(GO) clean
	@rm -f $(TARGET)
	@rm -rf dist/*

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which $(TARGET))

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

run: build
	./$(TARGET)
