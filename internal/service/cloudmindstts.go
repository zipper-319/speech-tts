package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	pb "speech-tts/api/tts/v2"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/pkg/pointer"
	"speech-tts/internal/pkg/trace"
	"speech-tts/internal/utils"
	"strings"
)

type CloudMindsTTSService struct {
	pb.UnimplementedCloudMindsTTSServer
	log         log.Logger
	uc          *service.TTSService
	sdkFailTime int
	failTexts   map[string]struct{}
}

func NewCloudMindsTTSService(logger log.Logger, uc *service.TTSService) *CloudMindsTTSService {
	failTexts := make(map[string]struct{})
	return &CloudMindsTTSService{
		log:       logger,
		uc:        uc,
		failTexts: failTexts,
	}
}

func (s *CloudMindsTTSService) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {

	spanCtx, span := trace.NewTraceSpan(conn.Context(), "TTSService v2 call", nil)

	if req.TraceId == "" {
		uuidNum, _ := uuid.NewRandom()
		req.TraceId = fmt.Sprintf("%s-%s", "sdk", uuidNum.String())
	}

	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()

	logger := log.NewHelper(log.With(s.log, "traceId", req.TraceId, "rootTraceId", req.RootTraceId))
	logger.Infof("call TTSServiceV2;the req——————text:%s;speakerName:%s;Emotions:%s,DigitalPerson:%s,ParameterFlag:%v,Expression:%s,Movement:%s ",
		req.Text, req.ParameterSpeakerName, req.Emotions, req.ParameterDigitalPerson, req.ParameterFlag, req.Expression, req.Movement)

	if req.Text == "" {
		return errors.New("text param is null")
	}

	if req.ParameterSpeakerName == "" {
		req.ParameterSpeakerName = "DaXiaoFang"
	} else {
		temp := strings.Split(req.ParameterSpeakerName, "_")
		if len(temp) > 1 {
			req.ParameterSpeakerName = temp[0]
		}
	}

	if !s.uc.IsLegalSpeaker(req.ParameterSpeakerName) {
		return errors.New("ParameterSpeakerName param is invalid")
	}
	if req.Emotions != "" && !s.uc.IsLegalEmotion(req.Emotions) {
		return errors.New("emotion param is invalid")
	}
	if req.Pitch != "" && !s.uc.IsLegalPitch(req.Pitch) {
		return errors.New("pitch param is invalid")
	}
	if req.Expression != "" && !s.uc.IsLegalExpression(req.Expression) {
		return errors.New("expression param is invalid")
	}
	if req.Movement != "" && !s.uc.IsLegalMovement(req.Movement) {
		return errors.New("movement param is invalid")
	}

	object := s.uc.GeneHandlerObjectV2(spanCtx, req.ParameterSpeakerName, logger)
	pUserData, err := pointer.Save(object)
	if err != nil {
		return err
	}
	defer pointer.Unref(pUserData)
	logger.Infof("CallTTSServiceV2;pUserData:%v", pUserData)
	id, err := s.uc.CallTTSServiceV2(req, pUserData)
	logger.Infof("CallTTSServiceV2;pUserData:%v;id:%d", pUserData, id)
	if err != nil {
		return err
	}
	isInterrupted := false
	for response := range object.BackChan {
		if !isInterrupted {
			if err := conn.Send(&response); err != nil {
				log.Errorf("send err:%v", err)
				span.SetStatus(codes.Error, fmt.Sprintf("Err send:%v", err))
				isInterrupted = true
				span.SetAttributes(attribute.Key("IsInterrupted").Bool(true))
				s.uc.CancelTTSService(id)

			}
		}

	}

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
		SpeedList:      s.uc.SupportedSpeed,
		VolumeList:     s.uc.SupportedVolume,
		PitchList:      s.uc.GetSupportedPitch(),
		EmotionList:    s.uc.GetSupportedEmotion(),
		MovementList:   s.uc.GetSupportedMovement(),
		ExpressionList: s.uc.GetSupportedExpression(),
	}, nil
}
