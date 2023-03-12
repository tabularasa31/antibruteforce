package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
	"github.com/tabularasa31/antibruteforce/internal/usecase"
	"github.com/tabularasa31/antibruteforce/pkg/grpcserver"
	"github.com/tabularasa31/antibruteforce/pkg/logger"
	"github.com/tabularasa31/antibruteforce/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Logger
	logg, err := logger.GetLogger(cfg)
	if err != nil {
		log.Fatalf("unable to load logger: %v", err)
	}
	defer func() {
		_ = logg.Sync()
	}()

	logg.Info("...config successfully parsed")

	logg.Info("starting redis...")
	// Redis
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}

	newRedis := redis.NewClient(opt)
	if er := newRedis.Ping(); er.String() != "ping: PONG" {
		logg.Error(fmt.Sprintf("client Redis ping connection error: %v", err))
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	} else {
		logg.Info("...redis successfully connected")
	}

	// Bucket repo
	bucketRepo := repo.NewBucketRepo(newRedis, &cfg.App)

	// Postgres db create
	db, err := postgres.New(cfg)
	if err != nil {
		logg.Error(fmt.Sprintf("app - Run - postgres.New: %v", err))
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	} else {
		logg.Info("...postgres successfully connected")
	}
	defer db.Close()

	// White and black lists
	listRepo := repo.NewListRepo(db)
	if er := listRepo.Up(); er != nil {
		logg.Error(fmt.Sprintf("app - Run - listRepo.Up: %v", err))
	}

	// Use cases
	useCases := usecase.New(bucketRepo, listRepo)

	// GRPC Server
	logg.Info("Starting grpc server...")

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		logg.Error(fmt.Sprintf("app - Run - net.Listen: %v", err))
	}

	grpcServer := grpcserver.New(useCases, lis, logg)

	s := <-interrupt
	logg.Info("app - Run - signal: " + s.String())

	grpcServer.Shutdown()
}
