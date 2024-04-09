#!/bin/bash

if [ -z "$SERVER_ID" ]; then
    HOSTNAME="$(hostname -s)"
    HOST_ID="$(echo ${HOSTNAME} | sed -ne 's/^.*-\([0-9][0-9]*\)$/\1/p')"
    SERVER_ID=$(( HOST_ID + 1 ))
fi

if [ -z "$SERVER_ID" ]; then
   echo "NO Server id"
   exit 1
fi

PROJECT_NAME="speech-tts"
apt-get install libcurl3 tree -y

cd /opt/speech/tts/
ln -s /data third
ln -s /data/res
rm -rf res/config.json
cp conf/sdkconfig.json res/config.json

#mkdir lib_interface
#cd lib_interface
#ln -s /data/libCmTts.so.online_voicetuning libCmTts.so
cd /opt/speech/tts
tree
ls -l  lib_interface

export LD_LIBRARY_PATH=./third/lib:./lib_interface:${LD_LIBRARY_PATH}

# exec bin/tts_svrv2 -id "${SERVER_ID}"
stdbuf -oL env GOTRACEBACK=crash ./bin/$PROJECT_NAME id "${SERVER_ID}"
