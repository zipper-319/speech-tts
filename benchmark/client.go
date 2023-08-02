package benchmark

import (
	"context"
	"google.golang.org/grpc"
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

func TestTTSV1(ctx context.Context, addr, text, speaker, traceId, robotTraceId string, num int) error {

	now := time.Now()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV1Client := v1.NewCloudMindsTTSClient(conn)
	req := &v1.TtsReq{
		Text:                 text,
		ParameterSpeakerName: speaker,
		TraceId:              traceId,
		RootTraceId:          robotTraceId,
	}

	response, err := ttsV1Client.Call(ctx, req)

	if err != nil {
		log.Println(err)
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Printf("TestTTSV2;traceId:%s,return after cancel", traceId)
				return
			default:
				temp, err := response.Recv()
				if err == io.EOF {
					log.Println("finished to receive ")
					return
				}
				if err != nil {
					log.Println(err)
					return
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
		}
	}()
	wg.Wait()
	log.Printf("--------------------------------TestTTSV1----(%d); cost:%d\n\n", num, time.Since(now).Milliseconds())
	return nil
}

func TestTTSV2(ctx context.Context, addr, text, speaker, traceId, robotTraceId, movement, expression string, num int) error {

	now := time.Now()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV2Client := v2.NewCloudMindsTTSClient(conn)

	flagSet := make(map[string]string, 3)
	flagSet["mouth"] = "true"
	if movement != "" {
		flagSet["movement"] = movement
	}
	if expression != "" {
		flagSet["expression"] = expression
	}
	req := &v2.TtsReq{
		Text:                 text,
		ParameterSpeakerName: speaker,
		TraceId:              traceId,
		RootTraceId:          robotTraceId,
		ParameterFlag:        flagSet,
		Version:              v2.ClientVersion_version,
	}

	response, err := ttsV2Client.Call(ctx, req)
	if err != nil {
		log.Printf("Text:%s, err;%v", text, err)
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Printf("traceId:%s;  TestTTSV2 return after cancel\n", traceId)
				return
			default:
				temp, err := response.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Println(err)
					return
				}
				if temp.ErrorCode != 0 {
					log.Println("tts 内部服务错误：", temp.ErrorCode)
				}
				log.Printf("receive message(Type %T)", temp)

				if audio, ok := temp.ResultOneof.(*v2.TtsRes_SynthesizedAudio); ok {
					log.Printf("pcm length:%d, status:%d", len(audio.SynthesizedAudio.Pcm), temp.Status)
				}
			}
		}
	}()
	wg.Wait()
	log.Printf("-------------------------,TestTTSV2---(%d);cost:%d\n\n", num, time.Since(now).Milliseconds())
	return nil
}
