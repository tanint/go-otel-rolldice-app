package o11y

import (
	"context"
	"fmt"
	"time"

	"github.com/demo/rolldice/config"

	exceptions "github.com/demo/rolldice/pkg/exceptions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func createBasicAuthHeaders(authToken string) map[string]string {
	return map[string]string{
		"Authorization": "Basic " + authToken,
	}
}

func createResource(ctx context.Context, appName string) (*sdkresource.Resource, error) {
	return sdkresource.New(
		ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
			attribute.String("application", fmt.Sprintf("/%s", appName)),
		),
	)
}

func createTraceExporter(otlpEndpoint string, authToken string) (*otlptrace.Exporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(createBasicAuthHeaders(authToken)),
	)
	return otlptrace.New(context.Background(), client)
}

func createTraceProvider(resource *sdkresource.Resource, spanProcessor sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(spanProcessor),
	)
}

func createMetricExporter(ctx context.Context, otlpEndpoint string, authToken string) (*otlpmetrichttp.Exporter, error) {
	return otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(otlpEndpoint),
		otlpmetrichttp.WithHeaders(createBasicAuthHeaders(authToken)),
	)
}

func createMeterProvider(resource *sdkresource.Resource, metricExporter *otlpmetrichttp.Exporter) (*sdkmetric.MeterProvider, error) {
	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	), nil
}

func createLogExporter(ctx context.Context, otlpEndpoint string, authToken string) (*otlploghttp.Exporter, error) {
	return otlploghttp.New(
		ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithHeaders(createBasicAuthHeaders(authToken)),
	)
}

func createLoggerProvider(resource *sdkresource.Resource, exporter *otlploghttp.Exporter) *sdklog.LoggerProvider {
	processor := sdklog.NewBatchProcessor(exporter)
	return sdklog.NewLoggerProvider(
		sdklog.WithResource(resource),
		sdklog.WithProcessor(processor),
	)
}

type InitResult struct {
	Shutdown       func()
	LoggerProvider *sdklog.LoggerProvider
}

func InitOTel(config *config.InitOTelConfig) InitResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := createResource(ctx, config.AppName)
	exceptions.Print(err, "Error creating OTLP resource")

	// Create trace exporter and provider
	traceExporter, err := createTraceExporter(config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Trace exporter")
	traceProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := createTraceProvider(resource, traceProcessor)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Create metric exporter and provider
	metricExporter, err := createMetricExporter(ctx, config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Metric exporter")
	metricProvider, err := createMeterProvider(resource, metricExporter)
	exceptions.Print(err, "Error creating Metric provider")

	// Create log exporter and logger provider
	logExporter, err := createLogExporter(ctx, config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Log exporter")
	loggerProvider := createLoggerProvider(resource, logExporter)

	return InitResult{
		Shutdown: func() {
			exceptions.Print(traceProvider.Shutdown(ctx), "Error shutting down Trace provider")
			exceptions.Print(metricProvider.Shutdown(ctx), "Error shutting down Metric provider")
			exceptions.Print(loggerProvider.Shutdown(ctx), "Error shutting down Log provider")
			cancel()
		},
		LoggerProvider: loggerProvider,
	}
}
