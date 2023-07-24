#!/bin/bash

AppName="speech-tts"
IMAGE_NAME="cloudminds-tts"

user="tts"
password="Q9dMZp_pGX"
CI_COMMIT_TAG=`git log --pretty=format:"%h" -1`
DOCKER_REGISTRY_HOST="harbor.cloudminds.com"
VERSION="v4.1.15"


expect <<EOF
set timeout 30
spawn ssh 172.16.31.72 -p 10022 -l root
expect {
"*(yes/no)?" { send "yes\r",exp_continue }
"*password"  { send "123456\r" }
}
expect "*:~#" { send "cd speech-tts\r" }
expect "*:~/speech-tts#" { send "bash build.sh $VERSION $AppName $CI_COMMIT_TAG\r" }
expect "*:~/speech-tts#" { send "exit \r" }
expect eof
EOF

expect -c '
  spawn scp -P  10022 root@172.16.31.72:~/speech-tts/bin/speech-tts ./bin
  expect "*password"
  send "123456\r"
  interact
'

docker build  --no-cache -t harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG .
echo DOCKER_REGISTRY_USER=$user DOCKER_REGISTRY_PASSWORD=$password
echo $password |  docker login -u $user --password-stdin $DOCKER_REGISTRY_HOST >/dev/null 2>&1 && docker push harbor.cloudminds.com/$AppName/$IMAGE_NAME:$VERSION.$CI_COMMIT_TAG
