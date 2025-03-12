APP_NAME := clevis
CURRENT_PLATFORM := $(shell printf '%s-%s' $$(go env GOOS GOARCH))
VERSION := $(shell date '+%Y.%-m.%-d')-$(shell git rev-parse --short HEAD)
LIB_VERSION := $(shell go list -m github.com/anatol/clevis.go | awk '{split($$2,x,"-"); print x[2] "-" x[3]}')
PLATFORMS := $(sort darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm $(CURRENT_PLATFORM))

MAKEFLAGS += -j
.DEFAULT_GOAL := $(CURRENT_PLATFORM)
.PHONY: clean update build release lint

clean:
	@rm -rf dist/

update:
	@go get -u ./cmd
	@go mod tidy

$(PLATFORMS): OUTPUT=$(APP_NAME)-$@-$(VERSION)$(if $(findstring windows,$@),.exe,)
$(PLATFORMS): export GOOS=$(word 1,$(subst -, ,$@))
$(PLATFORMS): export GOARCH=$(word 2,$(subst -, ,$@))
$(PLATFORMS):
	@echo $(OUTPUT)
	@go build \
		-C cmd \
		-trimpath \
		-buildvcs=false \
		-ldflags '-s -w -X main.version=$(VERSION) -X main.libVersion=$(LIB_VERSION)' \
		-o '../dist/$(OUTPUT)'

build: $(PLATFORMS)

release: lint clean build
	@parallel 'xz -zv {}' ::: dist/$(APP_NAME)-*
	@git tag -f 'v$(VERSION)'

lint:
	@go vet ./cmd
	@-golangci-lint run
	@gofmt -d ./cmd
