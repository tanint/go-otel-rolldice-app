package config

import (
	"fmt"
	"log"
	"os"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/joho/godotenv"
)

type InitOTelConfig struct {
	OtlpEndpoint          string
	AppName               string
	HttpExporterAuthToken string
	TracingSampler        sdktrace.Sampler
}

func LoadOtelConfig() (*InitOTelConfig, error) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	appEnv := os.Getenv("APP_ENV")

	config := &InitOTelConfig{
		OtlpEndpoint:          os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		AppName:               os.Getenv("SERVICE_NAME"),
		HttpExporterAuthToken: os.Getenv("OTEL_HTTP_EXPORTER_AUTH_TOKEN"),
	}

	if appEnv == "production" {
		config.TracingSampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
	} else {
		config.TracingSampler = sdktrace.AlwaysSample()
	}

	if config.OtlpEndpoint == "" || config.AppName == "" || config.HttpExporterAuthToken == "" {
		return nil, fmt.Errorf("missing required configuration")
	}

	return config, nil
}
