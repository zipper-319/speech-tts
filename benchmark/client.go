package benchmark

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
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

func TestTTSV1(addr, text, speaker string) error {
	ctx := context.Background()
	now := time.Now()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV1Clinet := v1.NewCloudMindsTTSClient(conn)
	req := &v1.TtsReq{
		Text: text,
		//ParameterSpeakerName: speaker,
	}

	response, err := ttsV1Clinet.Call(ctx, req)
	if err != nil {
		log.Error(err)
	}
	for {
		temp, err := response.Recv()
		if err != nil {
			break
		}
		if temp.Error != v1.TtsErr_TTS_ERR_OK {
			log.Info("tts 内部服务错误：", temp.Error)
		}
		if temp.Status == v1.PcmStatus_STATUS_END {
			log.Infof("cost:%d", time.Since(now).Milliseconds())
		} else {
			log.Infof("pcm length:%d, status:%s", len(temp.Pcm), temp.Status)
		}
	}
	return nil
}

func TestTTSV2(addr string, text string) error {
	ctx := context.Background()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV2Client := v2.NewCloudMindsTTSClient(conn)
	req := &v2.TtsReq{
		Text: text,
	}
	response, err := ttsV2Client.Call(ctx, req)
	if err != nil {
		log.Error(err)
	}
	for {
		temp, err := response.Recv()
		if err != nil {
			break
		}
		if temp.ErrorCode != 0 {
			log.Info("tts 内部服务错误：", temp.ErrorCode)
		}
		log.Info(temp)
	}
	return nil
}
