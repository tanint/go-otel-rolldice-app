package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/demo/rolldice/internal/notification/events"
	"github.com/demo/rolldice/internal/notification/models"
	"github.com/demo/rolldice/internal/notification/services"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type RollDiceResultNotiEventHandler struct {
	logger      *logrus.Logger
	tracer      trace.Tracer
	lineService *services.LineService
}

func RollDiceResultEventHandler(lineService *services.LineService, logger *logrus.Logger, tracer trace.Tracer) *RollDiceResultNotiEventHandler {
	return &RollDiceResultNotiEventHandler{logger, tracer, lineService}
}

func (h *RollDiceResultNotiEventHandler) Handle(ctx context.Context, event *events.RollEvent) error {
	ctx, span := h.tracer.Start(ctx, "Start processing RollEvent")
	defer span.End()

	payload := models.PushMessage{
		To: os.Getenv("LINE_BOT_RECEIVER_ID"),
		Messages: []models.Message{
			{
				Type: "text",
				Text: fmt.Sprintf("Rolled result: %d", event.Result),
			},
		},
	}

	h.lineService.SendPushMessage(ctx, payload)

	h.logger.WithContext(ctx).Infof("Send notification for rolled result: %d", event.Result)

	return nil
}
