#!/bin/bash

PROJECT_NAME="speech-tts"
COMMIT=$(git describe --tags --always)
FILE=$(date +%F).log


PID=$(ps x| grep $PROJECT_NAME | grep -v grep | awk '{print $1}')
if [[ $PID ]]; then
  kill -9 $PID
fi

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:`pwd`/internal/cgo/libs
export GODEBUG=cgocheck=0


go mod tidy -compat=1.17
go mod download
go mod verify
mkdir -p bin/
go build  -ldflags "-s -w -X main.commit=$COMMIT" -o ./bin/$PROJECT_NAME  `pwd`/cmd/$PROJECT_NAME/... && ulimit -c unlimited && mkdir -p log/ && bin/$PROJECT_NAME