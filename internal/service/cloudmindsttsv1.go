package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	pb "speech-tts/api/tts/v1"
	"speech-tts/internal/cgo/service"
	jwtUtil "speech-tts/internal/pkg/jwt"
	"speech-tts/internal/pkg/pointer"
	"speech-tts/internal/pkg/trace"
	"speech-tts/internal/utils"
	"strings"
)

type CloudMindsTTSServiceV1 struct {
	pb.UnimplementedCloudMindsTTSServer
	log log.Logger
	uc  *service.TTSService
}

func NewCloudMindsTTSServiceV1(logger log.Logger, uc *service.TTSService) *CloudMindsTTSServiceV1 {
	return &CloudMindsTTSServiceV1{
		log: logger,
		uc:  uc,
	}
}

func (s *CloudMindsTTSServiceV1) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {
	ctx := conn.Context()
	spanCtx, span := trace.NewTraceSpan(ctx, "TTSService v1 call", nil)
	defer span.End()

	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()
	var identifier string
	if tokenInfo, ok := ctx.Value(jwtUtil.Identifier{}).(*jwtUtil.IdentityClaims); ok {
		identifier = tokenInfo.Account
	}

	logger := log.NewHelper(log.With(s.log, "traceId", req.TraceId, "rootTraceId", req.RootTraceId))
	logger.Infof("call TTSServiceV1;the req——————text:%s;speakerName:%s;Emotions:%s;Pitch:%s;identifier:%s",
		req.Text, req.ParameterSpeakerName, req.Emotions, req.Pitch, identifier)

	if req.ParameterSpeakerName == "" {
		req.ParameterSpeakerName = "DaXiaoFang"
	} else {
		temp := strings.Split(req.ParameterSpeakerName, "_")
		if len(temp) > 1 {
			req.ParameterSpeakerName = temp[0]
		}
	}
	//if !s.uc.IsLegalSpeaker(req.ParameterSpeakerName) {
	//	return errors.New("ParameterSpeakerName param is invalid")
	//}
	//if req.Emotions != "" && !s.uc.IsLegalEmotion(req.Emotions) {
	//	return errors.New("emotion param is invalid")
	//}
	//if req.Pitch != "" && !s.uc.IsLegalPitch(req.Pitch) {
	//	return errors.New("pitch param is invalid")
	//}

	object := s.uc.GeneHandlerObjectV1(spanCtx, req.ParameterSpeakerName, logger)
	pUserData, err := pointer.Save(object)
	if err != nil {
		return err
	}
	defer pointer.Unref(pUserData)

	logger.Infof("CallTTSServiceV1;pUserData:%v", pUserData)
	id, err := s.uc.CallTTSServiceV1(req, pUserData)
	logger.Infof("CallTTSServiceV1;pUserData:%v;id:%d", pUserData, id)
	if err != nil {
		return err
	}
	isInterrupted := false
	for response := range object.BackChan {
		if !isInterrupted {
			if err := conn.Send(&response); err != nil {
				logger.Errorf("send err:%v", err)
				span.SetStatus(codes.Error, fmt.Sprintf("Err send:%v", err))
				isInterrupted = true
				span.SetAttributes(attribute.Key("IsInterrupted").Bool(true))
				s.uc.CancelTTSService(id)
			}
		}
	}

	return nil
}
func (s *CloudMindsTTSServiceV1) GetVersion(ctx context.Context, req *pb.VerReq) (*pb.VerRsp, error) {
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

	return &pb.VerRsp{
		Version: version,
	}, nil
}
func (s *CloudMindsTTSServiceV1) MixCall(req *pb.MixTtsReq, conn pb.CloudMindsTTS_MixCallServer) error {
	for {
		err := conn.Send(&pb.TtsRes{})
		if err != nil {
			return err
		}
	}
}
func (s *CloudMindsTTSServiceV1) GetSpeaker(ctx context.Context, req *pb.VerReq) (*pb.SpeakerList, error) {
	speakerList := make([]*pb.SpeakerParameter, len(s.uc.Speakers))
	for i, speaker := range s.uc.Speakers {
		speakerList[i] = &pb.SpeakerParameter{
			SpeakerName:          speaker.SpeakerName,
			ParameterSpeakerName: speaker.ParameterSpeakerName,
		}
	}
	return &pb.SpeakerList{
		List: speakerList,
	}, nil
}
