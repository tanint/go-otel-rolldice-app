package services

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type RollDiceService struct {
	tracer trace.Tracer
}

func NewRollDiceService(tracer trace.Tracer) *RollDiceService {
	return &RollDiceService{
		tracer,
	}
}

func (s *RollDiceService) Dice(ctx context.Context) int {
	_, span := s.tracer.Start(ctx, "Rolling")
	defer span.End()

	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)

	diceRoll := rnd.Intn(6) + 1

	return diceRoll
}
