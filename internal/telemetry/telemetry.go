package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func SetupOtelSDK(ctx context.Context) (func(context.Context) error, error) {
	otlpClient := otlptracehttp.NewClient()

	exporter, err := otlptrace.New(ctx, otlpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	resource, err := resource.New(
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.DeploymentEnvironment("production"),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithBatcher(exporter),
		trace.WithSpanProcessor(
			baggagecopy.NewSpanProcessor(
				baggagecopy.AllowAllMembers,
			),
		),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider.Shutdown, nil
}
