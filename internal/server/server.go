package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"
	"reflect"
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
			log.Infof("request:%s", extractArgs(req))
			reply, err = handler(ctx, req)
			var result interface{}
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
				"result", result,
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
