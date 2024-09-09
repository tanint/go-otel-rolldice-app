package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
	tracer   trace.Tracer
}

func NewKafkaProducer(brokers []string, username, password string, logger *logrus.Logger, tracer trace.Tracer) (*KafkaProducer, error) {
	config := createProducerConfig(username, password)

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.WithError(err).Error("Failed to create Kafka SyncProducer")
		return nil, fmt.Errorf("failed to create Kafka SyncProducer: %w", err)
	}

	wrappedProducer := otelsarama.WrapSyncProducer(config, producer)

	return &KafkaProducer{
		producer: wrappedProducer,
		logger:   logger,
		tracer:   tracer,
	}, nil
}

func createProducerConfig(username, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.TLS.Enable = true

	return config
}

func (p *KafkaProducer) Publish(ctx context.Context, topic, value, key string) error {
	_, span := p.tracer.Start(ctx, "publish to kafka")
	defer span.End()

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(producerMessage))

	partition, offset, err := p.producer.SendMessage(producerMessage)
	if err != nil {
		p.logError(ctx, topic, key, value, err)
		return fmt.Errorf("failed to publish message to Kafka: %w", err)
	}

	p.logSuccess(ctx, topic, key, value, partition, offset)

	return nil
}

func (p *KafkaProducer) logSuccess(ctx context.Context, topic, key, value string, partition int32, offset int64) {
	p.logger.WithContext(ctx).WithFields(logrus.Fields{
		"topic":     topic,
		"value":     value,
		"key":       key,
		"partition": partition,
		"offset":    offset,
	}).Info("Message published to Kafka topic")
}

func (p *KafkaProducer) logError(ctx context.Context, topic, key, value string, err error) {
	p.logger.WithContext(ctx).WithFields(logrus.Fields{
		"topic": topic,
		"value": value,
		"key":   key,
	}).WithError(err).Error("Failed to publish message to Kafka")
}
