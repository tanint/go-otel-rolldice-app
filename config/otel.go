package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type InitOTelConfig struct {
	OtlpEndpoint          string
	AppName               string
	HttpExporterAuthToken string
}

func LoadOtelConfig() (*InitOTelConfig, error) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config := &InitOTelConfig{
		OtlpEndpoint:          os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		AppName:               os.Getenv("SERVICE_NAME"),
		HttpExporterAuthToken: os.Getenv("OTEL_HTTP_EXPORTER_AUTH_TOKEN"),
	}

	if config.OtlpEndpoint == "" || config.AppName == "" || config.HttpExporterAuthToken == "" {
		return nil, fmt.Errorf("missing required configuration")
	}

	return config, nil
}
