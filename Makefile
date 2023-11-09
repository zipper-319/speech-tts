GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
PROJECT_NAME = speech-tts
COMMIT := $(shell git describe --tags --always)
FILE = $(date +%F).log
VERSION = v4.2.5
DSUrl = 172.16.23.15:31637
DSAddr = 10.12.32.198:9001
pwd := $(shell pwd)


ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,cmd\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif


.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go get -u github.com/envoyproxy/protoc-gen-validate
	go get -u github.com/google/wire/cmd/wire@v0.5.0
	go get -u github.com/favadi/protoc-go-inject-tag
	go get github.com/gogo/protobuf/protoc-gen-gofast

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
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --doc_out=html,api.html:./api/tts   \
	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
 	       --validate_out=paths=source_relative,lang=go:./api \
 	       --openapiv2_out ./api \
           --openapiv2_opt logtostderr=true \
           --openapiv2_opt json_names_for_fields=false \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	go mod download
	go mod verify
	mkdir -p bin/
	export dataServiceAddr=$(DSAddr) && export dataServiceEnv=$(DSUrl) && export LD_LIBRARY_PATH=$(LD_LIBRARY_PATH):$(pwd)/internal/cgo/libs  && \
	go build  -ldflags "-s -w -X main.Commit=$(COMMIT) -X main.Version=$(VERSION)  -X main.Name=$(PROJECT_NAME)" -o ./bin/$(PROJECT_NAME)  $(pwd)/cmd/$(PROJECT_NAME)/...

buildso:
	go build -buildmode=c-shared -o export/libs/libttsgo.so  export/main.go
	cp export/libs/libttsgo.so internal/cgo/libs

start: build
	export dataServiceEnv=$(DSUrl) && export LD_LIBRARY_PATH=$(LD_LIBRARY_PATH):$(pwd)/internal/cgo/libs  && \
	ulimit -c unlimited && bin/$(PROJECT_NAME)


.PHONY: generate
# generate
generate:
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

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


