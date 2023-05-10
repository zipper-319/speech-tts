package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/trace"
	"speech-tts/internal/utils"
	"strings"

	pb "speech-tts/api/tts/v2"
)

type CloudMindsTTSService struct {
	pb.UnimplementedCloudMindsTTSServer
	log *log.Helper
	uc  *service.TTSService
}

func NewCloudMindsTTSService(logger log.Logger, uc *service.TTSService) *CloudMindsTTSService {
	return &CloudMindsTTSService{
		log: log.NewHelper(logger),
		uc:  uc,
	}
}

func (s *CloudMindsTTSService) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {
	ctx := context.Background()
	spanCtx, span := trace.NewTraceSpan(ctx, "TTSService v1 call")

	if v, exists := utils.SpeakerMap[strings.ToLower(req.ParameterSpeakerName)]; exists {
		req.ParameterSpeakerName = v
	}
	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()
	object := s.uc.GeneHandlerObjectV2(spanCtx, req.ParameterSpeakerName)
	if err := s.uc.CallTTSServiceV2(req, object); err != nil {
		return err
	}
	for response := range object.BackChan {
		if response.ResultOneof != nil {
			if audio,ok :=response.ResultOneof.(*pb.TtsRes_SynthesizedAudio); ok {
				span.SetAttributes(attribute.Key("response.audioPcm.len").Int(len(audio.SynthesizedAudio.Pcm)))
				span.SetAttributes(attribute.Key("audio.IsPunctuation").Int(int(audio.SynthesizedAudio.IsPunctuation)))
			}
		}
		span.SetAttributes(attribute.Key("response.status").Int(int(response.Status)))
		err := conn.Send(&response)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Err send:%v", err))
			object.IsInterrupted = true
			span.SetAttributes(attribute.Key("IsInterrupted").Bool(object.IsInterrupted))
			return err
		}
	}
	return nil
}
func (s *CloudMindsTTSService) GetVersion(ctx context.Context, req *pb.VerVersionReq) (*pb.VerVersionRsp, error) {
	return &pb.VerVersionRsp{}, nil
}
func (s *CloudMindsTTSService) GetTtsConfig(ctx context.Context, req *pb.VerReq) (*pb.RespGetTtsConfig, error) {
	return &pb.RespGetTtsConfig{}, nil
}
