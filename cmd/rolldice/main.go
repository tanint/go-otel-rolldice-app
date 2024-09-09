package main

import (
	"log"
	"os"

	"github.com/demo/rolldice/config"
	"github.com/demo/rolldice/internal/rolldice/api"
	"github.com/demo/rolldice/internal/rolldice/services"
	"github.com/demo/rolldice/pkg/logger"
	"github.com/demo/rolldice/pkg/messaging/kafka"
	"github.com/demo/rolldice/pkg/middlewares"
	"github.com/demo/rolldice/pkg/o11y"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func main() {
	e := echo.New()

	otelConfig, err := config.LoadOtelConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	otelservice := o11y.InitOTel(otelConfig)

	logger := logger.NewLogger(otelservice.LoggerProvider)

	defer otelservice.Shutdown()

	tracer := otel.Tracer("main")

	e.Use(middlewares.OtelMiddleware(otelConfig.AppName))

	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	kafkaBroker := os.Getenv("ENV_KAFKA_BROKERS")

	brokers := []string{kafkaBroker}

	kafkaProducer, err := kafka.NewKafkaProducer(brokers, kafkaUsername, kafkaPassword, logger, tracer)

	if err != nil {
		log.Fatal(err)
	}

	rolldiceService := services.NewRollDiceService(tracer, logger, kafkaProducer)

	api.InitRolldiceHandler(e, rolldiceService)

	e.Logger.Fatal(e.Start(":8083"))
}
