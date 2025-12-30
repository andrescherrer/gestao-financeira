package observability

import (
	"context"
	"testing"
)

func TestInitTracing_Disabled(t *testing.T) {
	// Test with tracing disabled
	config := TracingConfig{
		Enabled:     false,
		ServiceName: "test-service",
		Environment: "test",
		JaegerURL:   "http://localhost:14268/api/traces",
	}

	shutdown, err := InitTracing(config)
	if err != nil {
		t.Fatalf("InitTracing() error = %v, want nil", err)
	}

	// Shutdown should be a no-op function
	if shutdown == nil {
		t.Error("InitTracing() shutdown function should not be nil")
	}

	// Call shutdown (should not panic)
	shutdown()
}

func TestStartSpan_NoTracer(t *testing.T) {
	// Test StartSpan when tracer is not initialized
	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")

	// Should return a no-op span
	if span == nil {
		t.Error("StartSpan() should return a span even when tracer is not initialized")
	}
}

