#!/bin/bash

MODEL_Dir=`pwd`
MODULE_DIR=/devepu/jenkins1/workspace/CmTts
sourceDir=$MODULE_DIR/TTS_SDK_*
cd $MODULE_DIR
VERSION=`ls | grep TTS_SDK_ | grep -v zip| awk -F_ '{print $3}'`
echo 'version is' $VERSION
MODEL_PATH="speech-tts-model-out"
cd $MODEL_Dir
rm -rf $MODEL_PATH
modelx init $MODEL_PATH
mkdir -p $MODEL_PATH/lib
mkdir -p $MODEL_PATH/res

cp $sourceDir/lib/*  $MODEL_PATH/lib/
cp $sourceDir/libCmTts.so $MODEL_PATH/libCmTts.so.online_voicetuning
cp -r $sourceDir/res/* $MODEL_PATH/res
cd $MODEL_Dir/$MODEL_PATH

/usr/bin/expect -c '
  set timeout 30
  set TOKEN "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ6aXBwZXIuemhhbyIsImV4cCI6MTA4OTM3NTMwNzh9.Ji1YzUKOFTBSxRSa_BhPspRWD5A9mxFvbz1_CPuDTKK81QdU0xveEGp8GgUYJ6JMHFvQXFNjLo-kaQCgfICbiJjgoU67hC5OYf_5r9Au-4--XWabkWYbBiB10HnjBmmQP8_GbDgoa3sp3S0tKuIs-o4WGB8rYGx4M7_85TGYNcMlF6NwM8wg4UyeiQL-zbQnNqcdy7k8Kl3K-yN_95ZB8VAtSNnwIdXm6b1HotRkjVC2NO29AMgqdSIoCjxQNaad9mXGJkpRYCxmL90OV3x-HOSmVUkrlKYf5W5I_yd7wAwAtUpDIPcP5O7ay4plzNoH70t_9TFmN_6ihZPQ4O3w8w"
  puts $TOKEN
  spawn modelx login modelx
  expect "Token:"
  send "$TOKEN\r"
  interact
'

modelx list modelx/speech-tts/ttsv2-model
modelx push  modelx/speech-tts/ttsv2-model@${VERSION}