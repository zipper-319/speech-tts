#!/bin/bash

MODEL_Dir=`pwd`
MODULE_DIR=/devepu/jenkins1/workspace/CmTts/jni/NATIVE/tts
sourceDir=$MODULE_DIR/TTS_SDK_*
cd $MODULE_DIR
VERSION=`ls | grep TTS_SDK_ | awk -F_ '{print $3}'`
echo 'version is' $VERSION
MODEL_PATH="speech-tts-model-out"
cd $MODEL_Dir
rm -rf $MODEL_PATH
modelx init $MODEL_PATH
mkdir -p $MODEL_PATH/lib
mkdir -p $MODEL_PATH/res

cp $sourceDir/lib/*  $MODEL_PATH/lib/
cp  /usr/lib/x86_64-linux-gnu/libcurl.so.4.4.0  $MODEL_PATH/lib/libcurl.so.4.4.0
cp $sourceDir/libCmTts.so $MODEL_PATH/libCmTts.so.online_voicetuning
cp -r $MODULE_DIR/res/* $MODEL_PATH/res
rm -rf $MODEL_PATH/res/animation
mkdir -p $MODEL_PATH/res/animation
cp -r $MODULE_DIR/../animation/res/* $MODEL_PATH/res/animation
cd $MODEL_Dir/$MODEL_PATH

/usr/bin/expect login.exp

modelx list modelx/speech-tts/ttsv2-model
modelx push  modelx/speech-tts/ttsv2-model@${VERSION}