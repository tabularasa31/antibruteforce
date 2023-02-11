package grpcserver

import (
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
	"github.com/tabularasa31/antibruteforce/pkg/interceptors"
	"github.com/tabularasa31/antibruteforce/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	Server   *grpc.Server
	Listener net.Listener
	buckets  *repo.BucketRepo
	lists    *repo.ListRepo
	logger   logger.Logger
	cfg      *config.Config
}

// NewServer Server constructor -.
func NewServer(buckets *repo.BucketRepo, lists *repo.ListRepo, logger logger.Logger, cfg *config.Config) *Server {
	return &Server{buckets: buckets, lists: lists, logger: logger, cfg: cfg}
}

// Start server -.
func (s *Server) Start() error {
	im := interceptors.NewInterceptorManager(s.logger, s.cfg)

	lis, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		s.logger.Errorf("app - Run - net.Listen: %v", err)
	}
	s.Listener = lis

	s.Server = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Server.Time * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(s.Server)
	}

	return nil
}

// Shutdown -.
func (s *Server) Shutdown() error {
	s.Server.GracefulStop()
	s.logger.Info("grpc Server Exited Properly")
	return s.Listener.Close()
}
