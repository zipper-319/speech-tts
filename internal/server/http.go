package server

import (
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/conf"
	"speech-tts/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, tts *service.CloudMindsTTSService, ttsV1 *service.CloudMindsTTSServiceV1, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterCloudMindsTTSHTTPServer(srv, ttsV1)
	v2.RegisterCloudMindsTTSHTTPServer(srv, tts)
	return srv
}
