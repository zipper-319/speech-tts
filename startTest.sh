#!/bin/bash

PROJECT_NAME="speech-tts-test"
COMMIT=$(git describe --tags --always)
FILE=$(date +%F).log


PID=$(ps x| grep $PROJECT_NAME | grep -v grep | awk '{print $1}')
if [[ $PID ]]; then
  kill -9 $PID
fi

export LD_LIBRARY_PATH=`pwd`/internal/cgo/libs:$LD_LIBRARY_PATH
export GODEBUG=cgocheck=0
export dataServiceEnv=172.16.23.15:31637


go mod download
go mod verify
mkdir -p bin/
go build  -ldflags "-s -w -X main.commit=$COMMIT" -o ./bin/$PROJECT_NAME  `pwd`/benchmark/cmd/demo.go && ulimit -c unlimited && mkdir -p log/ && bin/$PROJECT_NAME