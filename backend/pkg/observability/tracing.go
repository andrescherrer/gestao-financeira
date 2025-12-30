package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Tracer is the global tracer instance
	Tracer trace.Tracer
)

// TracingConfig holds configuration for OpenTelemetry tracing
type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Environment string
	JaegerURL   string // e.g., "http://localhost:14268/api/traces"
}

// InitTracing initializes OpenTelemetry tracing
func InitTracing(config TracingConfig) (func(), error) {
	if !config.Enabled {
		// Return a no-op shutdown function
		return func() {}, nil
	}

	if config.ServiceName == "" {
		config.ServiceName = "gestao-financeira-api"
	}

	// Create Jaeger exporter
	// Note: Jaeger exporter is deprecated, but we'll use it for now
	// In production, consider using OTLP exporter instead
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()), // Sample all traces for now
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Get tracer
	Tracer = otel.Tracer(config.ServiceName)

	// Return shutdown function
	shutdown := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			// Log error but don't fail
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}

	return shutdown, nil
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if Tracer == nil {
		// Return no-op span if tracing is not initialized
		return ctx, trace.SpanFromContext(ctx)
	}
	return Tracer.Start(ctx, name, opts...)
}

