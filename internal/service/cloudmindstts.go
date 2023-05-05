package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"speech-tts/internal/cgo/service"
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

	if v, exists := utils.SpeakerMap[strings.ToLower(req.ParameterSpeakerName)]; exists {
		req.ParameterSpeakerName = v
	}
	object := s.uc.GeneHandlerObjectV2(ctx, req.ParameterSpeakerName)
	if err := s.uc.CallTTSServiceV2(req, object); err != nil {
		return err
	}
	for response := range object.BackChan {
		err := conn.Send(&response)
		if err != nil {
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
