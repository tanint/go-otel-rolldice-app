package o11y

import (
	"context"
	"fmt"
	"time"

	"github.com/demo/rolldice/config"
	exceptions "github.com/demo/rolldice/pkg/exceptions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newHttpExporter(otlpEndpoint string, authToken string) (*otlptrace.Exporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Basic " + authToken,
		}),
	)

	return otlptrace.New(context.Background(), client)
}

func newResource(ctx context.Context, appName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
			attribute.String("application", fmt.Sprintf("/%s", appName)),
		),
	)
}

func newTraceProvider(resource *resource.Resource, spanProcessor sdktrace.SpanProcessor) *sdktrace.TracerProvider {

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(spanProcessor),
	)

	return tracerProvider
}

func InitOTel(c *config.InitOTelConfig) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := newResource(ctx, c.AppName)
	exceptions.Print(err, "failed to create the OTLP resource")

	exporter, err := newHttpExporter(c.OtlpEndpoint, c.HttpExporterAuthToken)
	exceptions.Print(err, "failed to created the OTLP exporter")

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := newTraceProvider(resource, batchSpanProcessor)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		exceptions.Print(tracerProvider.Shutdown(ctx), "failed to gracefully shutdown the tracer provider")

		cancel()
	}
}
