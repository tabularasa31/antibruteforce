package logger

import (
	"github.com/tabularasa31/antibruteforce/config"
	"go.uber.org/zap"
)

func GetLogger(cfg *config.Config) (*zap.Logger, error) {
	var err error
	var l *zap.Logger

	switch cfg.Server.Mode {
	case "production":
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case "development":
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	default:
		l = zap.NewExample()
	}

	return l, nil

}
