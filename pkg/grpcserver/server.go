package grpcserver

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	proto "github.com/tabularasa31/antibruteforce/api"
	"github.com/tabularasa31/antibruteforce/config"
	grpcv1 "github.com/tabularasa31/antibruteforce/internal/controller/grpc/v1"
	"github.com/tabularasa31/antibruteforce/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	Server   *grpc.Server
	Listener net.Listener
	useCases *usecase.UseCases
	logg     *zap.Logger
	cfg      *config.Config
	notify   chan error
}

// Start server -.
func New(
	useCases *usecase.UseCases,
	lis net.Listener,
	logg *zap.Logger,
	cfg *config.Config) *Server {

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.Server.MaxConnectionIdle * time.Minute,
			Timeout:           cfg.Server.Timeout * time.Second,
			MaxConnectionAge:  cfg.Server.MaxConnectionAge * time.Minute,
			Time:              cfg.Server.Time * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logg),
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	reflection.Register(grpcServer)

	s := &Server{
		Server: grpcServer,
		notify: make(chan error, 1),
		logg:   logg,
	}

	srv := grpcv1.NewAntibruteforceService(*useCases, *logg)
	proto.RegisterAntiBruteforceServer(s.Server, srv)

	// start monitoring
	grpc_prometheus.Register(grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()
	logg.Info(fmt.Sprintf("Monitoring export listen %s", ":9091"))
	go func() {
		err := http.ListenAndServe(":9091", promhttp.Handler())
		if err != nil {
			logg.Error(err.Error())
		}
		http.Handle("/metrics", promhttp.Handler())
	}()

	// start server
	s.Start(lis)

	return s
}

func (s *Server) Start(lis net.Listener) {
	go func() {
		s.notify <- s.Server.Serve(lis)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() {
	s.Server.GracefulStop()
	s.logg.Info("gRPC server gracefully stopped")
}
