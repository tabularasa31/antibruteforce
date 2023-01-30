package logger

import (
	"os"

	"github.com/tabularasa31/antibruteforce/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger methods interface .
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// Logger -.
type APILogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

// Logger constructor -.
func NewAPILogger(cfg *config.Config) *APILogger {
	return &APILogger{cfg: cfg}
}

// For mapping config logger to email_service logger levels .
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *APILogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// Init logger .
func (l *APILogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Mode == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder = zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

// Logger methods

func (l *APILogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *APILogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *APILogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *APILogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *APILogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *APILogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *APILogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *APILogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *APILogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *APILogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *APILogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *APILogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *APILogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *APILogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
