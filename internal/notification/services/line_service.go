package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/demo/rolldice/internal/notification/models"
	"github.com/demo/rolldice/pkg/httpclient"
	"go.opentelemetry.io/otel/trace"
)

const (
	apiURL      = "https://api.line.me/v2/bot/message/push"
	contentType = "application/json"
)

type LineService struct {
	tracer     trace.Tracer
	HTTPClient httpclient.HTTPClient
	authToken  string
}

func NewLineService(tracer trace.Tracer, HTTPClient httpclient.HTTPClient, authToken string) *LineService {
	return &LineService{
		tracer,
		HTTPClient,
		authToken,
	}
}

func (s *LineService) SendPushMessage(ctx context.Context, payload models.PushMessage) error {
	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.authToken,
	}

	res, err := s.HTTPClient.Post(ctx, apiURL, contentType, payloadBytes, headers)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message. Status code: %d", res.StatusCode)
	}

	log.Println("Message sent successfully!")
	return nil
}
