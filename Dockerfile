FROM harbor.cloudminds.com/library/asr-mkl-base:bionic.CM-Beta-1.3

ENV LOGPATH=/opt/speech/tts/runtime/logs
ENV PROJECT=speech-tts
ENV MODULE=tts-server
ENV dataServiceEnv=tts-data-service:9001

RUN apt update &&\
    apt install -d libcurl3 -y

EXPOSE 4012
EXPOSE 3012

WORKDIR /opt/speech/tts

COPY bin .
COPY run_speech_tts_srv.sh /etc/services.d/speech-tts/run
