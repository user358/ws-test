package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"ws-test/internal/log"
)

// simple interceptor with logs requests
func logInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// calls the handler
	h, err := handler(ctx, req)

	// select log level depends on error
	logFn := log.Infow
	if err != nil {
		logFn = log.Errorw
	}

	logFn(
		"GRPC request",
		"method",
		info.FullMethod,
		"duration",
		time.Since(start),
		"err",
		err,
	)

	return h, err
}
