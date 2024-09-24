package telemetry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(ctx context.Context, resource *resource.Resource) (*trace.TracerProvider, error) {
	otlpClient := otlptracehttp.NewClient()

	exporter, err := otlptrace.New(ctx, otlpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
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

	return traceProvider, nil
}

func InitMeter(ctx context.Context, resource *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// exporter, err := prometheus.New()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	// }

	buyerDurationView := metric.NewView(
		metric.Instrument{
			Name: "buyer.buy_duration",
			Kind: metric.InstrumentKindHistogram,
		},
		metric.Stream{
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 100, 200, 300, 400, 500},
			},
		},
	)

	defautView := metric.NewView(
		metric.Instrument{
			Name: "*",
			Kind: metric.InstrumentKindCounter,
		},
		metric.Stream{},
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource),
		//metric.WithReader(exporter),
		metric.WithReader(metric.NewPeriodicReader(
			exporter,
			metric.WithInterval(3*time.Second),
		)),
		metric.WithView(buyerDurationView, defautView),
	)

	return meterProvider, nil
}

func SetupOtelSDK(ctx context.Context) (func(context.Context) error, error) {
	var (
		shutdownFunc []func(context.Context) error
		err          error
	)

	shutdown := func(ctx context.Context) error {
		var err error
		for _, f := range shutdownFunc {
			newErr := f(ctx)
			err = errors.Join(err, newErr)
		}
		shutdownFunc = nil
		return err
	}

	handle := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	resource, err := resource.New(
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.DeploymentEnvironment("production"),
		),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithTelemetrySDK(),
		resource.WithProcessRuntimeDescription(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceProvider, err := InitTracer(ctx, resource)
	if err != nil {
		handle(err)
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}
	shutdownFunc = append(shutdownFunc, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)

	metricProvider, err := InitMeter(ctx, resource)
	if err != nil {
		handle(err)
		return nil, fmt.Errorf("failed to initialize meter: %w", err)
	}
	shutdownFunc = append(shutdownFunc, metricProvider.Shutdown)
	otel.SetMeterProvider(metricProvider)

	return shutdown, nil
}
