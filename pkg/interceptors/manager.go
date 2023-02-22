package interceptors

import (
	"context"
	"fmt"
	"github.com/tabularasa31/antibruteforce/config"
	"go.uber.org/zap"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager struct {
	logger *zap.Logger
	cfg    *config.Config
}

// NewInterceptorManager constructor -.
func NewInterceptorManager(logger *zap.Logger, cfg *config.Config) *InterceptorManager {
	return &InterceptorManager{logger: logger, cfg: cfg}
}

// Logger Interceptor
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Info(fmt.Sprintf("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err))

	return reply, err
}
