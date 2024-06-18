package o11y

import (
	"context"
	"fmt"
	"time"

	"github.com/demo/rolldice/config"
	exceptions "github.com/demo/rolldice/pkg/exceptions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newHttpTraceExporter(otlpEndpoint string, authToken string) (*otlptrace.Exporter, error) {
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

func newMeterProvider(resource *resource.Resource, metricExporter *otlpmetrichttp.Exporter) (*sdkmetric.MeterProvider, error) {
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			metricExporter,
		)),
	)

	return mp, nil
}

func InitOTel(c *config.InitOTelConfig) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	metricExporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(c.OtlpEndpoint),
		otlpmetrichttp.WithHeaders(map[string]string{
			"Authorization": "Basic " + c.HttpExporterAuthToken,
		}),
	)
	exceptions.Print(err, "Error creating Metric exporter")

	resource, err := newResource(ctx, c.AppName)
	exceptions.Print(err, "Error creating OTLP resource")

	exporter, err := newHttpTraceExporter(c.OtlpEndpoint, c.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Trace provider")

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tp := newTraceProvider(resource, batchSpanProcessor)

	mp, err := newMeterProvider(resource, metricExporter)
	exceptions.Print(err, "Error creating Metric provider")

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		exceptions.Print(tp.Shutdown(ctx), "Error shutting down Trace provider")
		exceptions.Print(mp.Shutdown(ctx), "Error shutting down Metric provider")

		cancel()
	}
}
