#!/bin/bash

source=~/project/tts_test/CmTts/jni/NATIVE/tts/TTS_SDK_*
dstLibs=`pwd`/internal/cgo/libs
includePath=`pwd`/internal/cgo/include

rm -rf $dstLibs/*
rm -rf ./res

cp $source/lib/*  $dstLibs
cp $source/libCmTts.so $dstLibs
cp -R $source/res `pwd`/res
cp $source/include/* $includePath
