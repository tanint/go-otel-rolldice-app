package logger

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type Formatter struct {
	ecslogrus.Formatter
}

func NewLogger(provider *sdklog.LoggerProvider) *logrus.Logger {
	hook := otellogrus.NewHook("rolldice-service", otellogrus.WithLoggerProvider(provider))

	logger := logrus.New()
	logger.SetFormatter(&Formatter{
		Formatter: ecslogrus.Formatter{},
	})
	logger.SetLevel(logrus.InfoLevel)
	logger.AddHook(hook)

	return logger
}
