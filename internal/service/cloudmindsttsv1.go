package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"speech-tts/internal/cgo/service"
	"speech-tts/internal/utils"
	"strings"

	pb "speech-tts/api/tts/v1"
)

type CloudMindsTTSServiceV1 struct {
	pb.UnimplementedCloudMindsTTSServer
	log *log.Helper
	uc  *service.TTSService
}

func NewCloudMindsTTSServiceV1(logger log.Logger, uc *service.TTSService) *CloudMindsTTSServiceV1 {
	return &CloudMindsTTSServiceV1{
		log: log.NewHelper(logger),
		uc:  uc,
	}
}

func (s *CloudMindsTTSServiceV1) Call(req *pb.TtsReq, conn pb.CloudMindsTTS_CallServer) error {
	ctx := context.Background()

	if req.ParameterSpeakerName == "" {
		for _, speaker := range s.uc.SupportedSpeaker {
			if speaker.Id == int(req.Speaker) {
				req.ParameterSpeakerName = speaker.Name
				break
			}
		}
	}

	if v, exists := utils.SpeakerMap[strings.ToLower(req.ParameterSpeakerName)]; exists {
		req.ParameterSpeakerName = v
	}
	object := s.uc.GeneHandlerObjectV1(ctx, req.ParameterSpeakerName, s.log)
	if err := s.uc.CallTTSServiceV1(req, object); err != nil {
		return err
	}
	for response := range object.BackChan {
		err := conn.Send(&response)
		if err != nil {
			object.IsInterrupted = true
			return err
		}
	}
	return nil
}
func (s *CloudMindsTTSServiceV1) GetVersion(ctx context.Context, req *pb.VerReq) (*pb.VerRsp, error) {
	return &pb.VerRsp{}, nil
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
	return &pb.SpeakerList{}, nil
}
