#!/bin/bash

source=~/project/CmTts1/TTS_SDK_v2.1.2.03a37cf
dstLibs=`pwd`/internal/cgo/libs
includePath=`pwd`/internal/cgo/include

rm -rf $dstLibs/*
rm -rf ./res

cp $source/lib/*  $dstLibs
cp $source/libCmTts.so $dstLibs
cp -R $source/res `pwd`/res
cp $source/include/* $includePath
