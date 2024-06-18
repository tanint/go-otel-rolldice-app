package main

import (
	"log"

	"github.com/demo/rolldice/config"
	"github.com/demo/rolldice/internal/rolldice/api"
	"github.com/demo/rolldice/internal/rolldice/services"
	"github.com/demo/rolldice/pkg/logger"
	"github.com/demo/rolldice/pkg/middlewares"
	"github.com/demo/rolldice/pkg/o11y"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func main() {
	e := echo.New()

	otelConfig, err := config.LoadOtelConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	otelservice := o11y.InitOTel(otelConfig)

	logger := logger.NewLogger(otelservice.LoggerProvider)

	defer otelservice.Shutdown()

	tracer := otel.Tracer("main")

	e.Use(middlewares.OtelMiddleware(otelConfig.AppName))

	rolldiceService := services.NewRollDiceService(tracer, logger)

	api.InitRolldiceHandler(e, rolldiceService)

	e.Logger.Fatal(e.Start(":8083"))
}
