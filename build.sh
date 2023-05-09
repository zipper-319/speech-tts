#!/bin/bash

COMMIT=$(git describe --tags --always)
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:`pwd`/internal/cgo/libs
export GODEBUG=cgocheck=0

go mod tidy -compat=1.17
go mod download
go mod verify
mkdir -p bin/
go build  -ldflags "-s -w -X main.commit=$COMMIT" -o ./bin/$PROJECT_NAME  `pwd`/cmd/$PROJECT_NAME/...