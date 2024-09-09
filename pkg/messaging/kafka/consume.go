package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
)

// toggleConsumptionFlow pauses or resumes consumption based on the current state
func toggleConsumptionFlow(client sarama.ConsumerGroup, isPause *bool) {
	if *isPause {
		client.ResumeAll()
	} else {
		client.PauseAll()
	}

	*isPause = !*isPause
}

// createConsumerConfig creates a Kafka consumer configuration with SASL and TLS
func createConsumerConfig(username, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_5_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	// Configure SASL and TLS for secure connections
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.TLS.Enable = true

	return config
}

// StartConsumption starts Kafka consumer group and handles message processing
func StartConsumption(
	brokers []string,
	topics []string,
	clientId string,
	groupId string,
	username string,
	password string,
	handlerFunc func(*sarama.ConsumerMessage) error,
) error {
	keepRunning := true

	// Create Kafka consumer configuration
	config := createConsumerConfig(username, password)

	// Create the KafkaConsumerGroupHandler to process messages
	consumer := KafkaConsumerGroupHandler{
		ready:       make(chan bool),
		handlerFunc: handlerFunc,
	}

	// Wrap the consumer with OpenTelemetry for tracing
	wrapped := otelsarama.WrapConsumerGroupHandler(&consumer)

	// Create a context to handle cancellation
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, groupId, config)
	if err != nil {
		log.Panicf("error creating consumer client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Start the Kafka consumption in a separate goroutine
	go func() {
		defer wg.Done()
		for {
			// Start consuming messages
			if err := client.Consume(ctx, topics, wrapped); err != nil {
				log.Panicf("error initiating consumption: %v", err)
			}

			// Exit if the context is done
			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	// Wait for the consumer to be ready
	<-consumer.ready
	log.Println("Sarama consumer up and running!...")

	// Handle signals for pausing, resuming, and termination
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			keepRunning = false
		case <-sigterm:
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Panicf("error closing the client: %v", err)
	}

	return nil
}
