package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var Logger *zap.Logger

func Init(appEnv string) {
	var cfg zap.Config

	if appEnv == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.MessageKey = "message"
		cfg.EncoderConfig.CallerKey = "caller"
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "T"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.LevelKey = "L"
		cfg.EncoderConfig.MessageKey = "M"
		cfg.EncoderConfig.CallerKey = "C"
		cfg.Encoding = "console"
	}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
}
