package benchmark

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var globalClientConn unsafe.Pointer
var lck sync.Mutex

func GetGrpcConn(addr string, ctx context.Context) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConn, unsafe.Pointer(conn))
	return conn, nil
}

func TestTTSV1(addr, text, speaker, traceId, robotTraceId string, num int) error {
	md := metadata.Pairs(
		"authorization", "Bearer some-secret-token",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	now := time.Now()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV1Client := v1.NewCloudMindsTTSClient(conn)
	req := &v1.TtsReq{
		Text:                 fmt.Sprintf("%s,测试第%d句话", text, num),
		ParameterSpeakerName: speaker,
		TraceId:              traceId,
		RootTraceId:          robotTraceId,
	}

	response, err := ttsV1Client.Call(ctx, req)
	if err != nil {
		log.Println(err)
		return err
	}
	for {
		temp, err := response.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		if temp.Error != v1.TtsErr_TTS_ERR_OK {
			log.Printf("tts 内部服务错误：%v", temp.Error)
		}
		if temp.Status == v1.PcmStatus_STATUS_END {
			log.Printf("cost:%d", time.Since(now).Milliseconds())
		} else {
			log.Printf("pcm length:%d, status:%s", len(temp.Pcm), temp.Status)
		}
	}
	return nil
}

func TestTTSV2(addr, text, speaker, traceId, robotTraceId, movement, expression string, num int) error {
	md := metadata.Pairs(
		"authorization", "Bearer some-secret-token",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV2Client := v2.NewCloudMindsTTSClient(conn)
	flagSet := make(map[string]string, 3)
	flagSet["mouth"] = "true"
	if movement != "" {
		flagSet["movement"] = "true"
	}
	if expression != "" {
		flagSet["expression"] = "true"
	}
	req := &v2.TtsReq{
		Text:                 fmt.Sprintf("%s,测试第%d句话", text, num),
		ParameterSpeakerName: speaker,
		TraceId:              traceId,
		RootTraceId:          robotTraceId,
		Movement:             movement,
		Expression:           expression,
		ParameterFlag:        flagSet,
	}
	now := time.Now()
	response, err := ttsV2Client.Call(ctx, req)
	if err != nil {
		log.Println(err)
		return err
	}
	for {
		temp, err := response.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		if temp.ErrorCode != 0 {
			log.Println("tts 内部服务错误：", temp.ErrorCode)
			break
		}
		log.Printf("receive message(Type %T)", temp)

		if audio, ok := temp.ResultOneof.(*v2.TtsRes_SynthesizedAudio); ok {
			log.Printf("pcm length:%d, status:%d", len(audio.SynthesizedAudio.Pcm), temp.Status)
		}
	}
	log.Printf("----------cost:%d\n\n", time.Since(now).Milliseconds())
	return nil
}
