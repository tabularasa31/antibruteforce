package grpcserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	proto "github.com/tabularasa31/antibruteforce/api"
	grpcv1 "github.com/tabularasa31/antibruteforce/internal/controller/grpc/v1"
	"github.com/tabularasa31/antibruteforce/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	_defaultMaxConnectionIdle = 5 * time.Minute
	_defaultMaxConnectionAge  = 5 * time.Minute
	_defaultTimeout           = 15 * time.Second
	_defaultTime              = 5 * time.Minute
)

type Server struct {
	Server   *grpc.Server
	Listener net.Listener
	logg     *zap.Logger
	notify   chan error
}

// Start server -.
func New(
	useCases *usecase.UseCases,
	lis net.Listener,
	logg *zap.Logger,
) *Server {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: _defaultMaxConnectionIdle,
			Timeout:           _defaultTimeout,
			MaxConnectionAge:  _defaultMaxConnectionAge,
			Time:              _defaultTime,
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
	logg.Info(fmt.Sprintf("Monitoring export start listening %s", ":9091"))
	httpserver := &http.Server{
		Addr:         ":9091",
		Handler:      promhttp.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		err := httpserver.ListenAndServe()
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
