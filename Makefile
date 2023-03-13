.DEFAULT_GOAL := build
PROJECT_NAME = vault-plugin-headscale
SHELL = /bin/bash

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

export GOPATH

VERSION = $(shell git describe --tags --abbrev=5 --always)
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
.PHONY = environment build dependency tests race msan deploy_snap deploy_release
.RECIPEPREFIX = >


environment:
> @echo -e "system environment:\n"
> @env | sort -u;
> @echo -e "\ngo environment:\n"
> @go env

dependency:
> go get -v -d ./...

build: #dependency
> @echo -e "\nUsing $(GOPATH) as GOPATH";
> CGO_ENABLED=0 go build -v -o bin/$(PROJECT_NAME) \
>  -ldflags="-s -w -X 'github.com/adrienmrgn/vault-plugin-headscale/backend.Version=$(VERSION)'";

lint:
> go install github.com/mgechev/revive@v1.2.1
> revive -formatter friendly ${PKG_LIST}

tests:
> @echo -e "\nTesting Unit using $(GOPATH) as GOPATH";
> go test -short ${PKG_LIST}

race:
> @echo -e "\nTesting Race using $(GOPATH) as GOPATH";
> go test -race -short ${PKG_LIST}

destroy:
> @echo -e "\nDestroying local test infra with docker compose";
> docker-compose top;
> docker-compose down;

compose: destroy build
> @echo -e "\nDeploying local test infra with docker compose";
> docker-compose up --build -d
> @sleep 1
> @echo -e "\nEnabling Headscale plugin";
> @export VAULT_ADDR="http://127.0.0.1:8200"; \
>  vault login root; \
>  vault secrets enable -path=headscale vault-plugin-headscale