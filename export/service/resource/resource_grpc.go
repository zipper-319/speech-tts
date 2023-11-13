package resource

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	ttsData "speech-tts/export/service/proto"
	"sync"
	"sync/atomic"
	"unsafe"
)

var globalClientConn unsafe.Pointer
var lck sync.Mutex

func GetGrpcConn(ctx context.Context) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	conn, err := grpc.DialContext(ctx, GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConn, unsafe.Pointer(conn))
	return conn, nil
}

func GetTTSResByGrpc(ctx context.Context, resType ttsData.ResType, languageType ttsData.LanguageType) ([]*ttsData.GetTtsDataResponse_TTSData, error) {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		return nil, err
	}
	client := ttsData.NewTtsDataClient(conn)
	resp, err := client.GetTtsData(ctx, &ttsData.GetTtsDataRequest{
		Resource: resType,
		Language: languageType,
	})
	if err != nil {
		return nil, err
	}
	return resp.DataList, nil
}

func RegisterResServiceByGrpc(ctx context.Context, serviceName, callbackUrl string) error {
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

func UnRegisterResServiceByGrpc(ctx context.Context, serviceName, callbackUrl string) error {
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

func GetSpeakerModelByGrpc(ctx context.Context) ([]*ttsData.GetSpeakerModelResult_SpeakerModel, error) {
	conn, err := GetGrpcConn(ctx)
	if err != nil {
		return nil, err
	}
	client := ttsData.NewTtsDataClient(conn)
	resp, err := client.GetSpeakerModel(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.SpeakerModels, nil
}
