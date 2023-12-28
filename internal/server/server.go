package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"reflect"
	"runtime/debug"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
	jwtUtil "speech-tts/internal/pkg/jwt"
	"speech-tts/internal/pkg/trace"
	"strings"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

func server(logger log.Logger, timeout int64) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				kind      string
				operation string
			)
			status := 0
			message := "SUCCESS"
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			var version string
			if se := errors.FromError(err); se != nil {
				status = int(se.Code)
				message = se.Reason
			}
			isTimeout := false
			latency := time.Since(startTime).Milliseconds()
			if latency > timeout {
				isTimeout = true
			}
			level, stack := extractError(err)
			if isTimeout {
				level = log.LevelWarn
			}
			log.WithContext(ctx, logger).Log(level,
				"component", kind,
				"version", version,
				"operation", operation,
				"args", reflect.ValueOf(req).Elem().Interface(),
				"status", status,
				"message", message,
				"stack", stack,
				"result", reflect.ValueOf(reply).Elem().Interface(),
				"isTimeout", fmt.Sprintf("timeout is %t", isTimeout),
				"latency", latency,
			)
			return
		}
	}
}

func hasFieldName(i interface{}, name string) bool {
	_, result := reflect.TypeOf(i).Elem().FieldByName(name)
	return result
}

func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%#v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

// wrappedStream wraps around the embedded grpc.ServerStream,and intercepts the RecvMsg and SendMsg method call.
// SendMsg method call.
type wrappedStream struct {
	ctx context.Context
	grpc.ServerStream
	log.Logger
	firstTime    time.Time
	sendTimes    int
	sendAudioLen int
	*ttsReq
}

type ttsReq struct {
	speakerName string
	text        string
	traceId     string
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

type validator interface {
	Validate() error
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	var ttsReqPram *ttsReq
	if m != nil {
		switch req := m.(type) {
		case *v1.TtsReq:
			ttsReqPram = &ttsReq{
				speakerName: req.ParameterSpeakerName,
				text:        req.Text,
				traceId:     fmt.Sprintf("%s_%s", req.TraceId, req.RootTraceId),
			}
		case *v2.TtsReq:
			ttsReqPram = &ttsReq{
				speakerName: req.ParameterSpeakerName,
				text:        req.Text,
				traceId:     fmt.Sprintf("%s_%s", req.TraceId, req.RootTraceId),
			}
		}
		w.ttsReq = ttsReqPram
		log.NewHelper(w.Logger).Infof("Receive a message (Type: %T), ttsReq:%v", m, ttsReqPram)
	}
	return err
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	w.sendTimes += 1
	audioLength := 0
	var status int32
	var traceId string
	if w.ttsReq != nil {
		traceId = w.ttsReq.traceId
	}
	if m != nil {
		if resp, ok := m.(*v1.TtsRes); ok {
			audioLength = len(resp.Pcm)
			w.sendAudioLen += audioLength
			status = int32(resp.Status)
			if resp.Status == v1.PcmStatus_STATUS_END {
				log.NewHelper(w.Logger).Infof("traceId:%s;TtsRes status is 3, resp:{%+v}", traceId, resp)
			}

		}
		if resp, ok := m.(*v2.TtsRes); ok {
			status = resp.Status
			switch result := resp.ResultOneof.(type) {
			case *v2.TtsRes_SynthesizedAudio:
				audioLength = len(result.SynthesizedAudio.Pcm)
				w.sendAudioLen += audioLength
			case *v2.TtsRes_ActionElement:
				log.NewHelper(w.Logger).Infof("trace:%s; ActionElement:%v", traceId, result.ActionElement)
			case *v2.TtsRes_BodyMovement:
				log.NewHelper(w.Logger).Infof("trace:%s; BodyMovement(FrameSize:%d, StartTimeMs:%f)", traceId, result.BodyMovement.FrameSize, result.BodyMovement.StartTimeMs)
			case *v2.TtsRes_ConfigText:
				log.NewHelper(w.Logger).Infof("trace:%s; ConfigText:%v", traceId, result.ConfigText)
			case *v2.TtsRes_TimeMouthShapes:
				log.NewHelper(w.Logger).Infof("trace:%s; TimeMouthShapes(StartTimeMs:%f)", traceId, result.TimeMouthShapes.StartTimeMs)
			case *v2.TtsRes_Expression:
				log.NewHelper(w.Logger).Infof("trace:%s; Expression(FrameSize:%d, StartTimeMs:%f)", traceId, result.Expression.FrameSize, result.Expression.StartTimeMs)
			}
		}
	}
	_, span := trace.NewTraceSpan(w.ctx, "SendMsg", nil)
	defer span.End()
	span.SetAttributes(attribute.Key("SendMsg times").Int(w.sendTimes))
	span.SetAttributes(attribute.Key("SendMsg length").Int(audioLength))

	log.NewHelper(w.Logger).Infof("trace:%s;Send %d message (Type: %T) after %dms; the length of audio is %d; the total length is %d; status:%d",
		traceId, w.sendTimes, m, time.Since(w.firstTime).Milliseconds(), audioLength, w.sendAudioLen, status)
	if status == 3 {
		md := metadata.Pairs("cost", fmt.Sprintf("%d", time.Since(w.firstTime).Milliseconds()))
		w.SendHeader(md)
		w.SetTrailer(md)
	}
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream, logger log.Logger, ctx context.Context) grpc.ServerStream {
	return &wrappedStream{
		ctx:          ctx,
		ServerStream: s,
		Logger:       logger,
		firstTime:    time.Now(),
		sendTimes:    0,
	}
}

// valid validates the authorization
func valid(authorization string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	//Perform the token validation here.For the sake of this example,the code
	//here forgoes any of usual OAuth2 token validation and instead checks for
	// for token matching an arbitrary string.
	return token == "some-secret-token"
}

func streamInterceptor(logger log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// authentication (token verification)
		var err error
		code := codes.OK
		fullMethod := info.FullMethod
		now := time.Now()
		log.NewHelper(logger).Infof("---------------------start FullMethod:%s --------------------", fullMethod)
		defer func() {
			if e := recover(); e != nil {
				debug.PrintStack()
				err = status.Errorf(codes.Internal, "Panic err: %v", e)
			}
			log.NewHelper(logger).Infof("-----------------------end FullMethod:%s; cost:%.3fs; err:%v----------------", fullMethod, time.Since(now).Seconds(), err)
		}()

		ctx := ss.Context()

		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			code = codes.Internal
			return status.Errorf(code, "from incoming context err")
		}
		ctx, span := trace.NewTraceSpan(ctx, "TTSService", tr.RequestHeader())
		defer func() {
			span.SetAttributes(attribute.Key("grpc_code").Int(int(code)))
			if err != nil {
				span.SetAttributes(attribute.Key("err").String(err.Error()))
			}
			span.End()
		}()

		identifier, err := jwtUtil.IsValidity(logger, tr)
		if err != nil {
			code = codes.Unauthenticated
			return err
		}
		if identifier != nil {
			ctx = context.WithValue(ctx, jwtUtil.Identifier{}, identifier)
		}

		if err = handler(srv, newWrappedStream(ss, logger, ctx)); err != nil {
			code = codes.Internal
			log.NewHelper(logger).Errorf("------------RPC failed with error: %v", err)
			return status.Errorf(code, err.Error())
		}
		md := metadata.Pairs("cost", fmt.Sprintf("%d", time.Since(now).Milliseconds()))
		ss.SetTrailer(md)
		return err
	}
}
