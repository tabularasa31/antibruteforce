package grpcserver

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

const (
	_defaultMaxConnectionIdle = 5 * time.Minute
	_defaultMaxConnectionAge  = 5 * time.Minute
	_defaultTimeout           = 15 * time.Second
	_defaultTime              = 5 * time.Minute
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
func (s *Server) Shutdown() error {
	s.Server.GracefulStop()
	s.logg.Info("grpc Server Exited Properly")
	return s.Listener.Close()
}

//-------------

//package grpcserver
//
//import (
//"context"
//"fmt"
//"net"
//
//"google.golang.org/grpc"
//)
//
//type Server struct {
//	Address string
//	Server  *grpc.Server
//}
//
//func NewServer(address string, server *grpc.Server) *Server {
//	return &Server{
//		Address: address,
//		Server:  server,
//	}
//}

//func (s *Server) Start() error {
//	listener, err := net.Listen("tcp", s.Address)
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		if err := s.Server.Serve(listener); err != nil {
//			fmt.Printf("gRPC server error: %v\n", err)
//		}
//	}()
//
//	fmt.Printf("gRPC server started on %s\n", s.Address)
//	return nil
//}

func (s *Server) Stop() {
	s.Server.Stop()
	fmt.Printf("gRPC server stopped")
}

func (s *Server) GracefulStop() {
	s.Server.GracefulStop()
	fmt.Printf("gRPC server gracefully stopped")
}

//func (s *Server) StopWithGracePeriod(gracePeriod int) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gracePeriod)*time.Second)
//	defer cancel()
//	s.Server.GracefulStop()
//
//	fmt.Printf("gRPC server stopped with grace period of %d seconds on %s\n", gracePeriod, s.cfg.Server.Port)
//}

//func (s *Server) StopWithGracePeriod(gracePeriod int) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gracePeriod)*time.Second)
//	defer cancel()
//
//	if s.GracefulStop() == nil {
//		fmt.Printf("gRPC server gracefully stopped on %s\n", s.Address)
//		return
//	}
//
//	<-ctx.Done()
//
//	fmt.Printf("gRPC server stopped with grace period of %d seconds on %s\n", gracePeriod, s.Address)
//}
