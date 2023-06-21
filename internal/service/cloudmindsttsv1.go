package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	pb "speech-tts/api/tts/v1"
	"speech-tts/internal/cgo/service"
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

	spanCtx, span := trace.NewTraceSpan(conn.Context(), "TTSService v1 call", nil)
	defer span.End()

	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()
	logger := log.NewHelper(log.With(s.log, "traceId", req.TraceId, "rootTraceId", req.RootTraceId))
	logger.Infof("call TTSServiceV1;the req——————text:%s;speakerName:%s;Emotions:%s;Pitch:%s",
		req.Text, req.ParameterSpeakerName, req.Emotions, req.Pitch)

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

	object := s.uc.GeneHandlerObjectV1(spanCtx, req.ParameterSpeakerName, logger)
	PUserData := pointer.Save(object)
	defer pointer.Unref(PUserData)
	if err := s.uc.CallTTSServiceV1(req, PUserData); err != nil {
		return err
	}
	audioLen := 0
	for response := range object.BackChan {
		audioLen += len(response.Pcm)
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
