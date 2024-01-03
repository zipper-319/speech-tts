package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	pb "speech-tts/api/tts/v2"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/conf"
	"speech-tts/internal/data"
	jwtUtil "speech-tts/internal/pkg/jwt"
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
	path        string
	jwtKey      string
}

func NewCloudMindsTTSService(logger log.Logger, uc *service.TTSService, s *conf.Server) *CloudMindsTTSService {
	failTexts := make(map[string]struct{})
	return &CloudMindsTTSService{
		log:       logger,
		uc:        uc,
		failTexts: failTexts,
		path:      s.App.Path,
		jwtKey:    s.App.GetJwt().GetKey(),
	}
}

func (s *CloudMindsTTSService) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {
	ctx := conn.Context()
	spanCtx, span := trace.NewTraceSpan(ctx, "TTSService v2 call", nil)

	myTraceId := ctx.Value(jwtUtil.TraceId{})

	req.TraceId = fmt.Sprintf("sdk(%s)-%s", myTraceId, req.TraceId)

	span.SetAttributes(attribute.Key("speakerName").String(req.ParameterSpeakerName))
	span.SetAttributes(attribute.Key("traceId").String(req.TraceId))
	span.SetAttributes(attribute.Key("rootTraceId").String(req.RootTraceId))
	span.SetAttributes(attribute.Key("text").String(req.Text))
	defer span.End()
	var identifier string
	if tokenInfo, ok := ctx.Value(jwtUtil.Identifier{}).(*jwtUtil.IdentityClaims); ok {
		identifier = tokenInfo.Account
	}

	movement := req.ParameterFlag["movementPara"]
	expression := req.ParameterFlag["expressionPara"]

	logger := log.NewHelper(log.With(s.log, "traceId", req.TraceId, "rootTraceId", req.RootTraceId))
	logger.Infof("call TTSServiceV2;the req——————text:%s;speakerName:%s;Emotions:%s,DigitalPerson:%s,ParameterFlag:%v,Expression:%s,Movement:%s,clientVersion:%s, identifier:%s, userspace:%s",
		req.Text, req.ParameterSpeakerName, req.Emotions, req.ParameterDigitalPerson, req.ParameterFlag, expression, movement, req.Version, identifier, req.Userspace)

	if req.Text == "" {
		return errors.New("text param is null")
	}

	if req.Userspace == "" {
		req.Userspace = utils.DefaultUser
	}

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
	//if expression != "" && !s.uc.IsLegalExpression(expression) {
	//	return errors.New("expression param is invalid")
	//}
	//if movement != "" && !s.uc.IsLegalMovement(movement) {
	//	return errors.New("movement param is invalid")
	//}

	object := s.uc.GeneHandlerObjectV2(spanCtx, req.ParameterSpeakerName, logger)
	pUserData, err := pointer.Save(object)
	if err != nil {
		return err
	}
	defer pointer.Unref(pUserData)
	logger.Infof("pUserData:%d", pUserData)

	id, err := s.uc.CallTTSServiceV2(&data.Speaker{
		Text:                 req.Text,
		Speed:                req.Speed,
		Volume:               req.Volume,
		Pitch:                req.Pitch,
		Emotions:             req.Emotions,
		ParameterSpeakerName: req.ParameterSpeakerName,
		ParameterFlag:        req.ParameterFlag,
		Movement:             movement,
		Expression:           expression,
		Language:             req.Language,
		Userspace:            req.Userspace,
		AudioEncoding:        req.AudioEncoding,
	}, pUserData, fmt.Sprintf("%s_%s", req.RootTraceId, req.TraceId))
	logger.Infof("CallTTSServiceV2;pUserData:%v;id:%d", pUserData, id)
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

func (s *CloudMindsTTSService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	token, err := jwtUtil.GetToken(req.GetAccount(), int(req.Expire), s.jwtKey)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResp{Token: token}, nil
}

func (s *CloudMindsTTSService) GetUserSpeakers(ctx context.Context, req *pb.GetUserSpeakersRequest) (*pb.GetUserSpeakersResponse, error) {
	speakerList, err := s.uc.GetUserSpeakers(req.User)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserSpeakersResponse{
		Speakers: speakerList,
	}, nil
}

func (s *CloudMindsTTSService) GetTtsConfigByUser(ctx context.Context, req *pb.GetTtsConfigByUserRequest) (*pb.RespGetTtsConfig, error) {
	speakerList := make([]*pb.SpeakerParameter, 0, len(s.uc.Speakers))
	cloneSpeakerList, err := s.uc.GetUserSpeakers(req.User)
	if err != nil {
		return nil, err
	}

	for _, speaker := range s.uc.Speakers {
		speakerList = append(speakerList, &pb.SpeakerParameter{
			SpeakerName:          speaker.SpeakerName,
			ParameterSpeakerName: speaker.ParameterSpeakerName,
			IsSupportEmotion:     speaker.IsSupportEmotion,
			IsSupportMixedVoice:  speaker.IsSupportMixedVoice,
		})
	}
	log.Debugf("speakerList: %v", speakerList)
	for _, speaker := range cloneSpeakerList {
		speakerList = append(speakerList, &pb.SpeakerParameter{
			ParameterSpeakerName: speaker,
			IsBelongClone:        true,
		})
	}
	log.Debugf("cloneSpeakerList: %v", speakerList)

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
