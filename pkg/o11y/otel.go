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

func newResource(ctx context.Context, appName string) (*sdkresource.Resource, error) {
	return sdkresource.New(
		ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
			attribute.String("application", fmt.Sprintf("/%s", appName)),
		),
	)
}

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

func newTraceProvider(resource *sdkresource.Resource, spanProcessor sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(spanProcessor),
	)

	return tracerProvider
}

func newMeterProvider(resource *sdkresource.Resource, metricExporter *otlpmetrichttp.Exporter) (*sdkmetric.MeterProvider, error) {
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			metricExporter,
		)),
	)

	return mp, nil
}

func newMetricHttpExporter(ctx context.Context, otlpEndpoint string, authToken string) (*otlpmetrichttp.Exporter, error) {
	metricHttpExporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(otlpEndpoint),
		otlpmetrichttp.WithHeaders(map[string]string{
			"Authorization": "Basic " + authToken,
		}),
	)

	if err != nil {
		return nil, err
	}

	return metricHttpExporter, nil
}

func newLogHttpExporter(ctx context.Context, otlpEndpoint string, authToken string) (*otlploghttp.Exporter, error) {
	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithHeaders(map[string]string{
			"Authorization": "Basic " + authToken,
		}),
	)

	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func newLoggerProvider(res *sdkresource.Resource, exporter *otlploghttp.Exporter) *sdklog.LoggerProvider {
	processor := sdklog.NewBatchProcessor(exporter)
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	return provider
}

type InitResult struct {
	Shutdown       func()
	LoggerProvider *sdklog.LoggerProvider
}

func InitOTel(config *config.InitOTelConfig) InitResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := newResource(ctx, config.AppName)
	exceptions.Print(err, "Error creating OTLP resource")

	exporter, err := newHttpTraceExporter(config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Trace provider")
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tp := newTraceProvider(resource, batchSpanProcessor)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	metricHttpExporter, err := newMetricHttpExporter(ctx, config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Metric HTTP exporter")
	mp, err := newMeterProvider(resource, metricHttpExporter)
	exceptions.Print(err, "Error creating Metric provider")

	logHttpExporter, err := newLogHttpExporter(ctx, config.OtlpEndpoint, config.HttpExporterAuthToken)
	exceptions.Print(err, "Error creating Log HTTP exporter")
	lp := newLoggerProvider(resource, logHttpExporter)

	return InitResult{
		Shutdown: func() {
			exceptions.Print(tp.Shutdown(ctx), "Error shutting down Trace provider")
			exceptions.Print(mp.Shutdown(ctx), "Error shutting down Metric provider")
			exceptions.Print(lp.Shutdown(ctx), "Error shutting down Log provider")

			cancel()
		},
		LoggerProvider: lp,
	}
}
