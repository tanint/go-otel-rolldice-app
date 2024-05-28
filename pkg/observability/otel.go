package observability

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	exceptions "github.com/demo/rolldice/pkg/exceptions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newExporter(ctx context.Context, otlpEndpoint string) (*otlptrace.Exporter, error) {

	conn, err := grpc.NewClient(
		otlpEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	exceptions.ReportError(err, "unable to reach GRPC OTLP endpoint")

	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func newHttpExporter(otlpEndpoint string) (*otlptrace.Exporter, error) {
	username := "otel"
	password := "1234"
	auth := username + ":" + password
	authBase64 := base64.StdEncoding.EncodeToString([]byte(auth))
	authHeader := "Basic " + authBase64

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": authHeader,
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

// func newMeterProvider(resource *resource.Resource) (*sdkmetric.MeterProvider, error) {
// 	metricExporter, err := stdoutmetric.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	meterProvider := sdkmetric.NewMeterProvider(
// 		sdkmetric.WithResource(resource),
// 		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
// 			metricExporter,
// 			sdkmetric.WithInterval(3*time.Second),
// 		)),
// 	)

// 	return meterProvider, nil
// }

func InitTracer(otlpEndpoint string, appName string) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := newResource(ctx, appName)
	exceptions.ReportError(err, "failed to create the OTLP resource")

	// exporter, err := newExporter(ctx, otlpEndpoint)
	exporter, err := newHttpExporter(otlpEndpoint)
	exceptions.ReportError(err, "failed to created the OTLP exporter")

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := newTraceProvider(resource, batchSpanProcessor)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		exceptions.ReportError(tracerProvider.Shutdown(ctx), "failed to gracefully shutdown the tracer provider")
		cancel()
	}
}

func InitialiseOpentelemetry(otlpEndpoint string, appName string) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := newResource(ctx, appName)
	exceptions.ReportError(err, "failed to create the OTLP resource")

	exporter, err := newHttpExporter(otlpEndpoint)
	exceptions.ReportError(err, "failed to created the OTLP exporter")

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := newTraceProvider(resource, batchSpanProcessor)

	// meterProvider, err := newMeterProvider(resource)
	exceptions.ReportError(err, "failed to initialise the meter provider")

	otel.SetTracerProvider(tracerProvider)
	// otel.SetMeterProvider(meterProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		exceptions.ReportError(tracerProvider.Shutdown(ctx), "failed to gracefully shutdown the tracer provider")
		// exceptions.ReportError(meterProvider.Shutdown(ctx), "failed to gracefully shutdown the meter provider")
		cancel()
	}
}
