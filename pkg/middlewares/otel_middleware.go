package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func OtelMiddleware(appName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracerProvider := otel.GetTracerProvider()
			tracer := tracerProvider.Tracer("main")

			request := c.Request()
			savedCtx := request.Context()

			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()

			propagator := otel.GetTextMapPropagator()
			ctx := propagator.Extract(savedCtx, propagation.HeaderCarrier(request.Header))

			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(httpconv.ServerRequest(appName, request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			if path := c.Path(); path != "" {
				rAttr := semconv.HTTPRoute(path)
				opts = append(opts, oteltrace.WithAttributes(rAttr))
			}
			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			c.SetRequest(request.WithContext(ctx))

			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				status := c.Response().Status
				span.SetStatus(httpconv.ServerStatus(status))
				c.Error(err)
			}

			status := c.Response().Status
			span.SetStatus(httpconv.ServerStatus(status))
			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}

			return nil
		}
	}
}
