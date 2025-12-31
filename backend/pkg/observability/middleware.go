package observability

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates a middleware that adds OpenTelemetry tracing to requests
func TracingMiddleware() fiber.Handler {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.Tracer("gestao-financeira-api")

	return func(c *fiber.Ctx) error {
		// Extract trace context from headers
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})
		ctx := propagator.Extract(c.Context(), &headerCarrier{headers: headers})

		// Start span
		ctx, span := tracer.Start(ctx, c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.path", c.Path()),
				attribute.String("http.route", c.Route().Path),
				attribute.String("http.user_agent", c.Get("User-Agent")),
				attribute.String("http.request_id", c.Get("X-Request-ID")),
			),
		)
		defer span.End()

		// Store context in locals
		c.Locals("trace_context", ctx)

		// Process request
		err := c.Next()

		// Add response attributes
		statusCode := c.Response().StatusCode()
		statusText := http.StatusText(statusCode)
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.String("http.status_text", statusText),
		)

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			if err != nil {
				span.RecordError(err)
			}
			span.SetStatus(codes.Error, statusText)
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		return err
	}
}

// headerCarrier implements the TextMapCarrier interface for Fiber headers
type headerCarrier struct {
	headers map[string]string
}

func (h *headerCarrier) Get(key string) string {
	return h.headers[key]
}

func (h *headerCarrier) Set(key, value string) {
	h.headers[key] = value
}

func (h *headerCarrier) Keys() []string {
	keys := make([]string, 0, len(h.headers))
	for k := range h.headers {
		keys = append(keys, k)
	}
	return keys
}
