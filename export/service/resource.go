package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"os"
	ttsData "speech-tts/export/service/proto"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var globalClientConn unsafe.Pointer
var lck sync.Mutex

const AddrDefault = "127.0.0.1:8080"

func GetGrpcConn(ctx context.Context) (*grpc.ClientConn, error) {
	addr := os.Getenv("dataServiceEnv")
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
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
	)
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
	fileName := fmt.Sprintf("%s_%s.txt", resType.String(), languageType.String())
	os.Rename(fileName, fileName+".bak"+time.Now().Format("20060102150405"))
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, v := range resp.Data {
		var n int
		n, err = f.WriteString(fmt.Sprintf("%s:%s\n", v.Key, v.Value))
		log.Infof("write string:%d,%v", n, err)
	}
	return fileName, err
}

func RegisterResService(ctx context.Context, serviceName, callbackUrl string) {

}
