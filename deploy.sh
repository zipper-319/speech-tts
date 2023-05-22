#!/bin/bash

AppName="speech-tts"
IMAGE_NAME="cloudminds-tts"
user="tts"
password="Q9dMZp_pGX"
CI_COMMIT_TAG=`git log --pretty=format:"%h" -1`
DOCKER_REGISTRY_HOST="harbor.cloudminds.com"
expect -c '
  spawn scp -P -r 10022 root@172.16.31.72:~/speech-tts/bin/* ./bin
  expect "*password"
  send "123456\r"
  interact
'
VERSION="v4.1.0"

docker build  --no-cache -t harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG .
echo DOCKER_REGISTRY_USER=$user DOCKER_REGISTRY_PASSWORD=$password
echo $password |  docker login -u $user --password-stdin $DOCKER_REGISTRY_HOST >/dev/null 2>&1 && docker push harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG