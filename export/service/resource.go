package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	ttsData "speech-tts/export/service/proto"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var globalClientConn unsafe.Pointer
var lck sync.Mutex

const AddrDefault = "10.12.32.198:9001"
const ResPath = "./res"

type CallbackFn func(ttsData.ResType, ttsData.LanguageType, string)

func GetGrpcConn(ctx context.Context) (*grpc.ClientConn, error) {
	addr := os.Getenv("dataServiceAddr")
	if addr == "" {
		addr = AddrDefault
	}
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConn, unsafe.Pointer(conn))
	return conn, nil
}

func GetTTSResAndSave(ctx context.Context, resType ttsData.ResType, languageType ttsData.LanguageType) (string, error) {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		return "", err
	}
	client := ttsData.NewTtsDataClient(conn)
	resp, err := client.GetTtsData(ctx, &ttsData.GetTtsDataRequest{
		Resource: resType,
		Language: languageType,
	})

	return SaveResource(resp, resType, languageType)
}

func RegisterResService(ctx context.Context, serviceName, callbackUrl string) error {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		log.Errorf("register res service error:%v", err)
		return err
	}
	client := ttsData.NewTtsDataClient(conn)
	_, err = client.RegisterResService(ctx, &ttsData.RegisterResServiceRequest{
		ServiceName: serviceName,
		CallbackUrl: callbackUrl,
	})
	return err
}

func UnRegisterResService(ctx context.Context, serviceName, callbackUrl string) error {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		log.Errorf("register res service error:%v", err)
		return err
	}
	client := ttsData.NewTtsDataClient(conn)
	_, err = client.UnRegisterResService(ctx, &ttsData.UnRegisterResServiceRequest{
		ServiceName: serviceName,
		CallbackUrl: callbackUrl,
	})
	return err
}

func InitTTSResource(ctx context.Context, fn CallbackFn) error {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		return err
	}
	client := ttsData.NewTtsDataClient(conn)

	for v, _ := range ttsData.ResType_name {
		for lang, _ := range ttsData.LanguageType_name {
			resType := ttsData.ResType(v)
			languageType := ttsData.LanguageType(lang)
			resp, err := client.GetTtsData(ctx, &ttsData.GetTtsDataRequest{
				Resource: resType,
				Language: languageType,
			})
			if err != nil {
				log.Errorf("get tts resource; resourceType:%s,language:%s, error:%v", resType, languageType, err)
				continue
			}
			fileName, err := SaveResource(resp, resType, languageType)
			if err != nil {
				log.Errorf("Save tts resource; resourceType:%s,language:%s, error:%v", resType, languageType, err)
				continue
			}
			fn(resType, languageType, fileName)
		}
	}
	return nil
}

func InitTTSRes(ctx context.Context) error {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		return err
	}
	client := ttsData.NewTtsDataClient(conn)

	for v, _ := range ttsData.ResType_name {
		for lang, _ := range ttsData.LanguageType_name {
			resType := ttsData.ResType(v)
			languageType := ttsData.LanguageType(lang)
			resp, err := client.GetTtsData(ctx, &ttsData.GetTtsDataRequest{
				Resource: resType,
				Language: languageType,
			})
			if err != nil {
				log.Errorf("get tts resource; resourceType:%s,language:%s, error:%v", resType, languageType, err)
				continue
			}
			fileName, err := SaveResource(resp, resType, languageType)
			if err != nil {
				log.Errorf("Save tts resource; resourceType:%s,language:%s, error:%v", resType, languageType, err)
				continue
			}
			log.Info(resType, languageType, fileName)
		}
	}
	return nil
}

func SaveResource(resp *ttsData.GetTtsDataResponse, resType ttsData.ResType, languageType ttsData.LanguageType) (string, error) {
	fileName := fmt.Sprintf("%s/%s_%s.txt", ResPath, resType.String(), languageType.String())
	os.Rename(fileName, fileName+".bak"+time.Now().Format("20060102150405"))

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, v := range resp.Data {
		var n int
		n, err = f.WriteString(fmt.Sprintf("%s@@%s\n", v.Key, v.Value))
		log.Infof("write string:%d,%v", n, err)
	}
	return fileName, nil
}
