package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/demo/rolldice/config"
	"github.com/demo/rolldice/internal/notification/events"
	"github.com/demo/rolldice/internal/notification/events/handlers"
	"github.com/demo/rolldice/internal/notification/services"
	"github.com/demo/rolldice/pkg/httpclient"
	"github.com/demo/rolldice/pkg/logger"
	"github.com/demo/rolldice/pkg/messaging/kafka"
	"github.com/demo/rolldice/pkg/o11y"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
)

func main() {

	otelConfig, err := config.LoadOtelConfig()
	otelConfig.AppName = os.Getenv("NOTIFICATION_SERVICE_NAME")

	if err != nil {
		log.Fatalf(err.Error())
	}

	otelservice := o11y.InitOTel(otelConfig)

	defer otelservice.Shutdown()

	tracer := otel.Tracer("main")
	logger := logger.NewLogger(otelservice.LoggerProvider)

	kafkaBroker := os.Getenv("ENV_KAFKA_BROKERS")
	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	lineBotApiAuthToken := os.Getenv("LINE_BOT_API_AUTH_TOKEN")

	brokers := []string{kafkaBroker}

	httpClient := httpclient.NewClient(tracer)
	lineService := services.NewLineService(tracer, httpClient, lineBotApiAuthToken)

	eventHandler := handlers.RollDiceResultEventHandler(lineService, logger, tracer)

	log.Println("Notification service is starting...")

	if err := kafka.StartConsumption(
		brokers,
		[]string{"poc.rolldice"},
		"poc-project",
		"poc-group",
		kafkaUsername,
		kafkaPassword,
		func(message *sarama.ConsumerMessage) error {
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

			var rolledEvent events.RollEvent
			if err := json.Unmarshal(message.Value, &rolledEvent); err != nil {
				log.Fatal(err)
			}

			if err := eventHandler.Handle(ctx, &rolledEvent); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	); err != nil {
		log.Fatal(err)
	}

}
