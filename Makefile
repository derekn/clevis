APP_NAME := clevis
CURRENT_PLATFORM := $(shell printf '%s-%s' $$(go env GOOS GOARCH))
VERSION := $(shell date '+%Y.%-m.%-d')
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LIB_VERSION := $(shell go list -m github.com/anatol/clevis.go | awk '{split($$2,x,"-"); print x[2] "-" x[3]}')
PLATFORMS := $(sort darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm $(CURRENT_PLATFORM))

MAKEFLAGS += -j
.DEFAULT_GOAL := build
.PHONY: clean update build build-all release lint test $(PLATFORMS)

clean:
	@rm -rf dist/

update:
	@go get -u ./cmd
	@go mod tidy

linux-arm: export GOARM=5
$(PLATFORMS): OUTPUT=$(APP_NAME)-$@$(if $(findstring windows,$@),.exe,)
$(PLATFORMS): export GOOS=$(word 1,$(subst -, ,$@))
$(PLATFORMS): export GOARCH=$(word 2,$(subst -, ,$@))
$(PLATFORMS):
	@go build \
		-C cmd \
		-trimpath \
		-buildvcs=false \
		-ldflags '-s -w -X main.version=$(VERSION)+$(GIT_COMMIT) -X main.libVersion=$(LIB_VERSION)' \
		-o '../dist/$(OUTPUT)'
	@echo $(OUTPUT)

build: $(CURRENT_PLATFORM)

build-all: $(PLATFORMS)

release: lint clean build-all
	@find dist -type f ! -name '*.exe' | parallel 'xz -z9v {}'
	@find dist -type f -name '*.exe' | parallel 'zip -m9 {.}.zip {}'
	@rhash -r --printf '%{sha-256}  %f\n' dist > dist/SHA256SUMS
	@git tag -f 'v$(VERSION)'

lint:
	@go vet ./cmd
	@-golangci-lint run
	@gofmt -d ./cmd
