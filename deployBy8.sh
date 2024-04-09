#!/bin/bash

AppName="speech-tts"
IMAGE_NAME="cloudminds-tts"

user="tts"
password="Q9dMZp_pGX"
CI_COMMIT_TAG=`git log --pretty=format:"%h" -1`
DOCKER_REGISTRY_HOST="harbor.cloudminds.com"
VERSION="v4.3.6"


expect -c '
  spawn scp  data@10.12.32.8:~/project/speech-tts/bin/speech-tts ./bin
  expect "*password"
  send "hankewei\r"
  interact
'

docker build  --no-cache -t harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG .
echo DOCKER_REGISTRY_USER=$user DOCKER_REGISTRY_PASSWORD=$password
echo $password |  docker login -u $user --password-stdin $DOCKER_REGISTRY_HOST >/dev/null 2>&1 && docker push harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG
