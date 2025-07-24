package logger

import (
	"go.uber.org/zap"
	"log"
)

func New(appEnv string) *zap.Logger {
	var logger *zap.Logger
	var err error

	if appEnv == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return logger
}
