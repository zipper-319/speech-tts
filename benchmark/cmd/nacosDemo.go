package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"io"
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/conf"
	"speech-tts/internal/pkg/nacos"
	"time"
)

func main() {
	config := &conf.Data{
		Nacos: &conf.Data_Nacos{
			Addr:                "nacos.region-dev-2.service.iamidata.com",
			Port:                31684,
			ContextPath:         "/nacos",
			NamespaceId:         "0021fbf85038b162b3d43794a1944bde39680a5d",
			Group:               "speech-tts-86",
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./runtime/nacos-logs",
			CacheDir:            "./runtime/nacos-cache",
			LogLevel:            "debug",
		},
	}
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///speech-tts.grpc"),
		grpc.WithDiscovery(nacos.NewRegister(config)))
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("state:%v", conn.GetState())
	conn.Connect()
	ttsV2Client := v2.NewCloudMindsTTSClient(conn)
	req := &v2.TtsReq{
		Text: "成都今天的天气",
	}
	now := time.Now()
	response, err := ttsV2Client.Call(context.Background(), req)
	if err != nil {
		log.Error(err)
	}
	for {
		temp, err := response.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			break
		}
		if temp.ErrorCode != 0 {
			log.Info("tts 内部服务错误：", temp.ErrorCode)
			break
		}
		log.Infof("receive message(Type %T)", temp)

		if audio, ok := temp.ResultOneof.(*v2.TtsRes_SynthesizedAudio); ok {
			log.Infof("pcm length:%d, status:", len(audio.SynthesizedAudio.Pcm))
		}
	}
	log.Infof("----------cost:%d\n\n", time.Since(now).Milliseconds())
}
