GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	INTERNAL_PROTO_FILES := $(shell powershell -Command "Get-ChildItem -Recurse internal -Filter *.proto | Resolve-Path -Relative")
    API_PROTO_FILES := $(shell powershell -Command "Get-ChildItem -Recurse api -Filter *.proto | Resolve-Path -Relative")
    INTERNAL_PROTO_FILES := $(subst \,/, $(INTERNAL_PROTO_FILES))
    API_PROTO_FILES := $(subst \,/, $(API_PROTO_FILES))
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# initializer env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

# install cli tools
cli:
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/gnostic@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api/proto \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api/gen \
 	       --go-http_out=paths=source_relative:./api/gen \
 	       --go-grpc_out=paths=source_relative:./api/gen \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

#在api中已经生成了不需要这个
.PHONY: errors
# generate errors proto
errors:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make wire;
	@echo '==============success=============='

 .PHONY: wire
wire:
	wire ./cmd/... ./internal/biz/... ./internal/data/... ./internal/service/...

.PHONY: test
# run all tests
test:
	go test -v ./internal/test/...

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
