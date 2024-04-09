FROM harbor.cloudminds.com/library/asr-mkl-base:bionic.CM-Beta-1.3

ENV LOGPATH=/opt/speech/tts/runtime/logs \
    PROJECT=speech-tts   \
    MODULE=tts-server    \
    dataServiceEnv=tts-data-service:9001  \
    dataServiceAddr=tts-data-server  \
    IsOpenGrpc=true \
    TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive


RUN echo ${TZ} > /etc/timezone

RUN apt update && apt install -y  -d libcurl3 libssl1.0.0

EXPOSE 4012
EXPOSE 3012

WORKDIR /opt/speech/tts

COPY bin/* ./bin/
#COPY libs/libCmTts.so  ./lib_interface/
COPY run_speech_tts_srv.sh /etc/services.d/speech-tts/run
