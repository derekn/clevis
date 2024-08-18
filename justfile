current_platform := `printf '%s-%s' $(go env GOOS GOARCH)`
version := `gitversion -showvariable semver`
commit_sha := `git rev-parse HEAD`
go_version := `go version | awk '{print $3}'`
lib_version := `go list -m github.com/anatol/clevis.go | awk '{split($2,x,"-"); print x[2] "-" x[3]}'`
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
		-ldflags '-s -w -X main.version={{ version }} -X main.commitSha={{ commit_sha }} -X main.goVersion={{ go_version }} -X main.libVersion={{ lib_version }}' \
		-o ../dist/clevis-{{ platform }}-{{ version }}
	echo clevis-{{ platform }}-{{ version }}

build-all *platforms=all-platforms:
	parallel 'just build {}' ::: {{ platforms }}

release: clean build-all
	parallel 'xz -zv {}' ::: dist/*
	@git tag -f 'v{{ version }}'
