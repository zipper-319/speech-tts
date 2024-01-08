#!/bin/bash

AppName="speech-tts"
IMAGE_NAME="cloudminds-tts"

user="tts"
password="Q9dMZp_pGX"
CI_COMMIT_TAG=`git log --pretty=format:"%h" -1`
DOCKER_REGISTRY_HOST="harbor.cloudminds.com"
VERSION="v4.4.3"


expect -c '
  spawn scp -P  10022 root@172.16.33.21:/devepu/jenkins1/workspace/speech-tts/bin/speech-tts ./bin
  expect "*password"
  send "123456\r"
  interact
'

docker build  --no-cache -t harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG .
echo DOCKER_REGISTRY_USER=$user DOCKER_REGISTRY_PASSWORD=$password
echo $password |  docker login -u $user --password-stdin $DOCKER_REGISTRY_HOST >/dev/null 2>&1 && docker push harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG
