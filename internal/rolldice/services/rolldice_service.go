package services

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/demo/rolldice/pkg/messaging/kafka"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RollDiceService struct {
	tracer   trace.Tracer
	logger   *logrus.Logger
	producer *kafka.KafkaProducer
}

type RollEvent struct {
	RollID    string `json:"roll_id"`
	Result    int    `json:"result"`
	Timestamp string `json:"timestamp"`
}

func NewRollDiceService(tracer trace.Tracer, logger *logrus.Logger, producer *kafka.KafkaProducer) *RollDiceService {
	return &RollDiceService{
		tracer,
		logger,
		producer,
	}
}

func (s *RollDiceService) Dice(ctx context.Context) (int, error) {
	ctx, span := s.tracer.Start(ctx, "Rolling")

	defer span.End()

	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)

	diceRoll := rnd.Intn(6) + 1

	rollId := strconv.Itoa(generateRollID())

	rollEvent := RollEvent{
		RollID:    rollId,
		Result:    diceRoll,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	value, _ := json.Marshal(rollEvent)

	err := s.producer.Publish(ctx, "poc.rolldice", string(value), rollId)

	if err != nil {
		return 0, err
	}

	s.logger.WithContext(ctx).Infof("Roll result = %d", diceRoll)

	span.SetAttributes(
		attribute.Int("app.roll.result", diceRoll),
	)

	return diceRoll, nil
}

// Counter to generate dynamic RollID
var rollIDCounter int32

func generateRollID() int {
	return int(atomic.AddInt32(&rollIDCounter, 1))
}
