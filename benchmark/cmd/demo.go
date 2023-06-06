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
	"speech-tts/internal/pkg/pointer"
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
	ttsService := service.NewTTSService(path, speakerSetting, logger)
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
		userData := pointer.Save(object)
		if err := ttsService.CallTTSServiceV1(req, userData); err != nil {
			panic(err)
		}
		log.NewHelper(logger).Info("---------end to CallTTSServiceV1-----------")
		for response := range object.BackChan {
			log.NewHelper(logger).Infof("response: status:%d; pcm length:%d;", response.Status, len(response.Pcm))
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
		pointer.Unref(userData)
		log.NewHelper(logger).Info("--------finish to call ---------", unsafe.Pointer(object))
		i += 1
		log.NewHelper(logger).Infof("finish %d time; cost:%dms", i, time.Since(now).Milliseconds())
	}
}
