package main

import (
	"github.com/demo/rolldice/internal/rolldice/api"
	"github.com/demo/rolldice/internal/rolldice/services"
	"github.com/demo/rolldice/pkg/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func main() {
	e := echo.New()

	shutdown := observability.InitialiseOpentelemetry("localhost:65479", "rolldice-app")
	defer shutdown()

	tracer := otel.Tracer("main")

	rolldiceService := services.NewRollDiceService(tracer)

	api.InitRolldiceHandler(e, rolldiceService)

	e.Logger.Fatal(e.Start(":8083"))
}
