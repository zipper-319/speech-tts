#!/bin/bash

COMMIT=$(git describe --tags --always)
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:`pwd`/internal/cgo/libs
export GODEBUG=cgocheck=0

if [ $# -lt 2 ];then
  echo "please input version and project name"
  exit
fi
version=$1
PROJECT_NAME=$2

go mod tidy -compat=1.17
go mod download
go mod verify
mkdir -p bin/
go build  -ldflags "-s -w -X main.Commit=$COMMIT -X main.Version=$version" -o ./bin/$PROJECT_NAME  `pwd`/cmd/$PROJECT_NAME/...
echo "build success"
