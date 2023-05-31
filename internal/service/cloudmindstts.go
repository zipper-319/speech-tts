package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/pkg/pointer"
	"speech-tts/internal/pkg/trace"
	"speech-tts/internal/utils"
	"strings"

	pb "speech-tts/api/tts/v2"
)

type CloudMindsTTSService struct {
	pb.UnimplementedCloudMindsTTSServer
	log log.Logger
	uc  *service.TTSService
}

func NewCloudMindsTTSService(logger log.Logger, uc *service.TTSService) *CloudMindsTTSService {
	return &CloudMindsTTSService{
		log: logger,
		uc:  uc,
	}
}

func (s *CloudMindsTTSService) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {

	spanCtx, span := trace.NewTraceSpan(conn.Context(), "TTSService v2 call", nil)

	if v, exists := utils.SpeakerMap[strings.ToLower(req.ParameterSpeakerName)]; exists {
		req.ParameterSpeakerName = v
	}
	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()

	logger := log.NewHelper(log.With(s.log, "traceId", req.TraceId, "rootTraceId", req.RootTraceId))
	object := s.uc.GeneHandlerObjectV2(spanCtx, req.ParameterSpeakerName, logger)
	PUserData := pointer.Save(object)

	audioLen := 0
	if err := s.uc.CallTTSServiceV2(req, PUserData); err != nil {
		return err
	}
	for response := range object.BackChan {
		if response.ResultOneof != nil {
			if audio, ok := response.ResultOneof.(*pb.TtsRes_SynthesizedAudio); ok {
				audioLen += len(audio.SynthesizedAudio.Pcm)
				span.SetAttributes(attribute.Key("audio.IsPunctuation").Int(int(audio.SynthesizedAudio.IsPunctuation)))
			}
		}
		err := conn.Send(&response)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Err send:%v", err))
			object.IsInterrupted = true
			span.SetAttributes(attribute.Key("IsInterrupted").Bool(object.IsInterrupted))
			return err
		}
	}
	span.SetAttributes(attribute.Key("response.audioPcm.len").Int(audioLen))
	return nil
}
func (s *CloudMindsTTSService) GetVersion(ctx context.Context, req *pb.VerVersionReq) (*pb.VerVersionRsp, error) {
	res := struct {
		ServerVersion     string
		TtsModelVersion   string
		ResServiceVersion string
	}{
		utils.GetServerVersion(),
		s.uc.GetSDKVersion(),
		s.uc.GetResServiceVersion(),
	}

	verStr, _ := json.Marshal(res)
	version := string(verStr)

	return &pb.VerVersionRsp{
		Version: version,
	}, nil
}
func (s *CloudMindsTTSService) GetTtsConfig(ctx context.Context, req *pb.VerReq) (*pb.RespGetTtsConfig, error) {
	speakerList := make([]*pb.SpeakerParameter, len(s.uc.Speakers))

	for i, speaker := range s.uc.Speakers {
		speakerList[i] = &pb.SpeakerParameter{
			SpeakerId:            int32(speaker.SpeakerId),
			SpeakerName:          speaker.SpeakerName,
			ParameterSpeakerName: speaker.ParameterSpeakerName,
			IsSupportEmotion:     speaker.IsSupportEmotion,
			IsSupportMixedVoice:  speaker.IsSupportMixedVoice,
		}
	}

	return &pb.RespGetTtsConfig{
		SpeakerList: &pb.SpeakerList{
			List: speakerList,
		},
		SpeedList:         s.uc.SupportedSpeed,
		VolumeList:        s.uc.SupportedVolume,
		PitchList:         s.uc.GetSupportedPitch(),
		EmotionList:       s.uc.GetSupportedEmotion(),
		DigitalPersonList: s.uc.GetSupportedDigitalPerson(),
	}, nil
}
