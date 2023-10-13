FROM harbor.cloudminds.com/library/asr-mkl-base:bionic.CM-Beta-1.3

ENV LOGPATH=/opt/speech/tts/runtime/logs \
    PROJECT=speech-tts   \
    MODULE=tts-server    \
    dataServiceEnv=tts-data-service:9001  \
    TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

#RUN apt update \
#    && apt install -y tzdata \
#    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
#    && echo ${TZ} > /etc/timezone\
#    && dpkg-reconfigure --frontend noninteractive tzdata \
#    && rm -rf /var/lib/apt/lists/*

RUN echo ${TZ} > /etc/timezone

RUN apt update && apt install -y  -d libcurl3 libssl1.0.0

EXPOSE 4012
EXPOSE 3012

WORKDIR /opt/speech/tts

COPY bin/* ./bin/
COPY run_speech_tts_srv.sh /etc/services.d/speech-tts/run
