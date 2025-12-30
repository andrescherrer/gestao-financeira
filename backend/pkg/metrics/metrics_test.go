package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsInitialization(t *testing.T) {
	// Test that metrics are properly initialized
	Init()

	// Verify that metrics are registered
	registry := prometheus.NewRegistry()
	registry.MustRegister(HTTPRequestDuration)
	registry.MustRegister(HTTPRequestTotal)
	registry.MustRegister(HTTPRequestInFlight)
	registry.MustRegister(BusinessMetrics.TransactionsCreated)
	registry.MustRegister(BusinessMetrics.UsersRegistered)

	// Metrics should be registered without errors
	if registry == nil {
		t.Error("Registry should not be nil")
	}
}

func TestBusinessMetrics(t *testing.T) {
	// Test business metrics increment
	Init()

	// Test transaction metrics
	BusinessMetrics.TransactionsCreated.WithLabelValues("INCOME").Inc()
	BusinessMetrics.TransactionsCreated.WithLabelValues("EXPENSE").Inc()

	count := testutil.ToFloat64(BusinessMetrics.TransactionsCreated.WithLabelValues("INCOME"))
	if count != 1.0 {
		t.Errorf("Expected 1 transaction created, got %f", count)
	}

	// Test user registration
	BusinessMetrics.UsersRegistered.Inc()
	userCount := testutil.ToFloat64(BusinessMetrics.UsersRegistered)
	if userCount != 1.0 {
		t.Errorf("Expected 1 user registered, got %f", userCount)
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "health endpoint",
			input:    "/health",
			expected: "/health",
		},
		{
			name:     "metrics endpoint",
			input:    "/metrics",
			expected: "/metrics",
		},
		{
			name:     "API root",
			input:    "/api/v1",
			expected: "/api/v1",
		},
		{
			name:     "transaction path",
			input:    "/api/v1/transactions",
			expected: "/api/v1/transactions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("normalizePath(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

