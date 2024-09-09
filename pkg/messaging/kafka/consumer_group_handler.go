package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaConsumerGroupHandler struct {
	ready       chan bool
	handlerFunc func(*sarama.ConsumerMessage) error
}

func (cg *KafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cg *KafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	close(cg.ready)
	return nil
}

func (cg *KafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			if err := cg.handlerFunc(message); err != nil {
				return nil
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
