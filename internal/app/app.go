package app

import (
	"fmt"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
	"github.com/tabularasa31/antibruteforce/pkg/boltdb"
	"github.com/tabularasa31/antibruteforce/pkg/grpcserver"
	"github.com/tabularasa31/antibruteforce/pkg/logger"
	"github.com/tabularasa31/antibruteforce/pkg/redis"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config")

	// Redis
	newRedis := redis.NewRedis(cfg)
	appLogger.Infof("............redis successfully connected")

	// Bucket repo
	bucketRepo := repo.NewBucketRepo(newRedis)

	// BoltDB
	// Open a connection to the BoltDB database
	err := boltdb.Init(cfg.DB)
	if err != nil {
		fmt.Println(err)
	}
	defer boltdb.Close()
	appLogger.Infof("............boltdb successfully connected")

	// Get a reference to the BoltDB instance
	db := boltdb.GetDB()

	// White and black lists
	listRepo := repo.NewListRepo(db)

	// GRPC Server
	appLogger.Infof("Starting server...")
	s := grpcserver.NewServer(bucketRepo, listRepo, appLogger, cfg)

	if err := s.Start(); err != nil {
		appLogger.Fatal("Failed ti start GRPC server: %v", err)
	}
	defer s.Shutdown()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}
