package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"net"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	nethttp "net/http"
)

// CheckTokenMiddleWare Check Token middleware
func LogFile(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {

				var (
					code      int32
					reason    string
					kind      string
					operation string
				)

				startTime := time.Now()
				kind = tr.Kind().String()
				operation = tr.Operation()

				httpTr := tr.(*http.Transport)

				reply, err = handler(ctx, req)
				if se := errors.FromError(err); se != nil {
					code = se.Code
					reason = se.Reason
				}

				body, _ := json.Marshal(req)

				logHelper(ctx, err, logger,
					"kind", "server",
					"component", kind,
					"url", httpTr.Request().URL.String(),
					"operation", operation,
					"method", httpTr.Request().Method,
					"body", string(body),
					"IP", getIP(httpTr.Request()),
					"code", code,
					"reason", reason,
					"latency", fmt.Sprintf("%dms", time.Since(startTime).Milliseconds()))
			}
			return
		}
	}
}

func getIP(r *nethttp.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return ""
	}
	return remoteIP.String()
}

func logHelper(ctx context.Context, err error, logger log.Logger, keyvals ...interface{}) {
	helper := log.NewHelper(log.WithContext(ctx, logger))
	if err != nil {
		helper.Errorw(append(keyvals, "stack", fmt.Sprintf("%+v", err))...)
	} else {
		helper.Infow(keyvals...)
	}
}

// CheckTokenMiddleWare Check Token middleware
func ClientLogFile(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}

			level, stack := extractError(err)
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"kind", "server",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
