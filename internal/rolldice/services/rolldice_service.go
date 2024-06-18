package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RollDiceService struct {
	tracer trace.Tracer
	logger *logrus.Logger
}

func NewRollDiceService(tracer trace.Tracer, logger *logrus.Logger) *RollDiceService {
	return &RollDiceService{
		tracer,
		logger,
	}
}

func (s *RollDiceService) Dice(ctx context.Context) int {
	ctx, span := s.tracer.Start(ctx, "Rolling")

	defer span.End()

	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)

	diceRoll := rnd.Intn(6) + 1

	s.logger.WithContext(ctx).Infof("Roll result = %d", diceRoll)

	span.SetAttributes(
		attribute.Int("app.roll.result", diceRoll),
	)

	return diceRoll
}
