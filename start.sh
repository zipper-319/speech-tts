#!/bin/bash

PROJECT_NAME="speech-tts"
COMMIT=$(git describe --tags --always)
FILE=$(date +%F).log
VERSION="v4.2.6"


PID=$(ps x| grep $PROJECT_NAME | grep -v grep | awk '{print $1}')
if [[ $PID ]]; then
  kill -9 $PID
fi

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:`pwd`/internal/cgo/libs
export GODEBUG=cgocheck=0
export dataServiceEnv=172.16.23.15:31637
export IsOpenGrpc=true
export dataServiceAddr=10.12.32.198

go mod download
go mod verify
mkdir -p bin/
go build  -ldflags "-s -w -X main.Commit=$COMMIT -X main.Version=$VERSION -X main.Name=$PROJECT_NAME" -o ./bin/$PROJECT_NAME  `pwd`/cmd/$PROJECT_NAME/... && ulimit -c unlimited && mkdir -p log/ &&nohup bin/$PROJECT_NAME   > log/$FILE 2>&1 &