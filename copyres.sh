#!/bin/bash

source=/devepu/jenkins1/workspace/CmTts/jni/NATIVE/tts/TTS_SDK_*
dstLibs=`pwd`/internal/cgo/libs

rm -rf $dstLibs/*
rm -rf ./res

cp $source/lib/*  $dstLibs
cp $source/libCmTts.so $dstLibs
cp -R $source/res `pwd`/res