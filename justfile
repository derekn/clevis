current_platform := `printf '%s-%s' $(go env GOOS GOARCH)`
version := `gitversion -showvariable semver`
all-platforms := 'darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm'

_default:
	@just --list

clean:
	rm -rfv dist

update:
	go get -u ./cmd
	go mod tidy

build platform=current_platform:
	#!/bin/bash -eu
	IFS='-' read -r GOOS GOARCH <<< {{ platform }}
	export GOOS GOARCH
	go build \
		-C cmd \
		-trimpath \
		-buildvcs=false \
		-ldflags '-s -w -X main.version={{ version }}' \
		-o ../dist/clevis-{{ platform }}-{{ version }}
	echo clevis-{{ platform }}-{{ version }}

build-all *platforms=all-platforms:
	parallel 'just build {}' ::: {{ platforms }}

release: clean build-all
	parallel 'xz -zv {}' ::: dist/*
	@git tag -f 'v{{ version }}'
