package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"os"
	v1 "speech-tts/api/tts/v1"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/data"
	"strconv"
	"time"
	"unsafe"
)

func main() {
	speaker := "DaXiaoFang"
	path := "./res"
	speakerSetting, err := data.NewSpeakerSetting(path)
	if err != nil {
		panic(fmt.Sprintf("fail to call speakerSetting;err:%v", err))
	}
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	ttsService := service.NewTTSService(path, speakerSetting)
	text := "成都今天的天气"
	i := 0
	for {
		now := time.Now()
		object := ttsService.GeneHandlerObjectV1(context.Background(), speaker, log.NewHelper(logger))
		log.NewHelper(logger).Info("pUserData-------------object", unsafe.Pointer(object))
		req := &v1.TtsReq{
			Text:    text,
			TraceId: strconv.Itoa(i),
		}
		if err := ttsService.CallTTSServiceV1(req, object); err != nil {
			panic(err)
		}
		for response := range object.BackChan {
			log.Info(response)
		}
		//	for {
		//		select {
		//		case response := <-object.BackChan:
		//			if response.Status == v1.PcmStatus_STATUS_END {
		//				goto TTSEnd
		//			}
		//		case <-time.After(200 * time.Millisecond):
		//			log.Info("timeout")
		//			goto TTSEnd
		//		}
		//
		//	}
		//TTSEnd:
		//time.Sleep(200 * time.Millisecond)

		log.NewHelper(logger).Info("pUserData-------------object", unsafe.Pointer(object))
		i += 1
		log.NewHelper(logger).Infof("finish %d time; cost:%dms", i, time.Since(now).Milliseconds())
	}
}
