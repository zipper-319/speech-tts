package benchmark

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
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
		Language:             "123",
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
				//if temp.Status == v1.PcmStatus_STATUS_END {
				//	log.Printf("cost:%d", time.Since(now).Milliseconds())
				//} else {
				//	log.Printf("pcm length:%d, status:%s", len(temp.Pcm), temp.Status)
				//}

			}
		}
	}()
	wg.Wait()
	trailerMD := response.Trailer()
	for key, value := range trailerMD {
		log.Printf("trailer key:%s, value:%s\n", key, value)
	}
	log.Printf("--------------------------------TestTTSV1----(%d); cost:%d\n\n", num, time.Since(now).Milliseconds())
	return nil
}

func TestTTSV2(ctx context.Context, user, addr, text, speaker, traceId, robotTraceId, movement, expression string, num int, isSaveFile bool) error {

	now := time.Now()
	conn, err := GetGrpcConn(addr, ctx)
	if err != nil {
		return err
	}
	ttsV2Client := v2.NewCloudMindsTTSClient(conn)

	flagSet := make(map[string]string, 3)
	flagSet["mouth"] = "false"
	if movement != "" {
		flagSet["movement"] = "true"
		flagSet["movementPara"] = movement
	} else {
		flagSet["movement"] = "false"
	}
	if expression != "" {
		flagSet["expressionPara"] = expression
		flagSet["expression"] = "true"
	} else {
		flagSet["expression"] = "false"
	}
	req := &v2.TtsReq{
		Text:                 text,
		ParameterSpeakerName: speaker,
		TraceId:              traceId,
		RootTraceId:          robotTraceId,
		ParameterFlag:        flagSet,
		Version:              v2.ClientVersion_Version,
		Language:             "zh",
		Userspace:            user,
	}

	log.Printf("----------------------TestTTSV2-----------(%d:%s)\n", num, text)
	response, err := ttsV2Client.Call(ctx, req)
	if err != nil {
		log.Printf("Text:%s, err;%v", text, err)
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var f *os.File
	var total int64

	if isSaveFile {
		f, err = os.OpenFile(fmt.Sprintf("./tmp/tts_%d.pcm", num), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
	}

	defer f.Close()
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
				var isFirstFrame bool
				if temp.ErrorCode != 0 {
					log.Println("tts 内部服务错误：", temp.ErrorCode)
				}
				//log.Printf("receive message(Type %T)", temp)

				if audio, ok := temp.ResultOneof.(*v2.TtsRes_SynthesizedAudio); ok {
					audioLength := int64(len(audio.SynthesizedAudio.Pcm))
					total += audioLength
					if total != 0 && total == audioLength {
						isFirstFrame = true
					}

					if isSaveFile && f != nil {
						n, err := f.Write(audio.SynthesizedAudio.Pcm)
						if err != nil {
							log.Println(err)
							return
						}
						log.Printf("pcm length:%d, total:%d, status:%d, write length:%d, isFirstFrame:%t, cost:%dms", audioLength, total, temp.Status, n, isFirstFrame, time.Since(now).Milliseconds())
					} else {
						log.Printf("pcm length:%d, total:%d, status:%d, isFirstFrame:%t, cost:%dms", audioLength, total, temp.Status, isFirstFrame, time.Since(now).Milliseconds())
					}
				}

			}
		}
	}()
	wg.Wait()
	trailerMD := response.Trailer()
	for key, value := range trailerMD {
		log.Printf("trailer key:%s, value:%s\n", key, value)
	}
	log.Printf("-------TestTTSV2---(%d:%s);client cost:%d,server cost:%d, first frame cost:%d\n\n", num, text, time.Since(now).Milliseconds())
	return nil
}
