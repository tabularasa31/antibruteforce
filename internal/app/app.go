package app

import (
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/pkg/logger"
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
	appLogger.Infof("Starting server")
}
