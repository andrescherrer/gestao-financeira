package eventbus

import (
	"errors"
	"testing"
	"time"

	"gestao-financeira/backend/internal/shared/domain/events"
)

// Mock event for testing
type mockEvent struct {
	events.BaseDomainEvent
}

func newMockEvent(eventType string) *mockEvent {
	return &mockEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(eventType, "test-id", "Test"),
	}
}

func (e *mockEvent) EventType() string {
	return e.BaseDomainEvent.EventType()
}

func TestEventBus_RetryLogic(t *testing.T) {
	tests := []struct {
		name           string
		handlerFails   int // Number of times handler should fail before succeeding
		maxRetries     int
		expectedResult bool // true if should succeed, false if should fail
	}{
		{
			name:           "Handler succeeds on first attempt",
			handlerFails:   0,
			maxRetries:     3,
			expectedResult: true,
		},
		{
			name:           "Handler succeeds after 1 retry",
			handlerFails:   1,
			maxRetries:     3,
			expectedResult: true,
		},
		{
			name:           "Handler succeeds after 2 retries",
			handlerFails:   2,
			maxRetries:     3,
			expectedResult: true,
		},
		{
			name:           "Handler fails after all retries",
			handlerFails:   4,
			maxRetries:     3,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			handler := func(event events.DomainEvent) error {
				attempts++
				if attempts <= tt.handlerFails {
					return errors.New("handler error")
				}
				return nil
			}

			retryConfig := RetryConfig{
				MaxRetries:        tt.maxRetries,
				InitialDelay:      10 * time.Millisecond,
				MaxDelay:          100 * time.Millisecond,
				BackoffMultiplier: 2.0,
			}

			eventBus := NewEventBusWithConfig(retryConfig)
			eventBus.SubscribeWithRetry("TestEvent", handler, retryConfig, "test-handler")

			event := newMockEvent("TestEvent")
			err := eventBus.Publish(event)

			if tt.expectedResult {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got success")
				}
			}
		})
	}
}

func TestEventBus_ErrorCallback(t *testing.T) {
	callbackCalled := false
	var callbackEventType string
	var callbackHandlerName string
	var callbackErr error

	errorCallback := func(eventType string, handlerName string, err error) {
		callbackCalled = true
		callbackEventType = eventType
		callbackHandlerName = handlerName
		callbackErr = err
	}

	eventBus := NewEventBus()
	eventBus.SetErrorCallback(errorCallback)

	handler := func(event events.DomainEvent) error {
		return errors.New("test error")
	}

	eventBus.SubscribeWithRetry("TestEvent", handler, DefaultRetryConfig(), "test-handler")

	event := newMockEvent("TestEvent")
	_ = eventBus.Publish(event)

	if !callbackCalled {
		t.Error("Error callback was not called")
	}

	if callbackEventType != "TestEvent" {
		t.Errorf("Expected event type 'TestEvent', got '%s'", callbackEventType)
	}

	if callbackHandlerName != "test-handler" {
		t.Errorf("Expected handler name 'test-handler', got '%s'", callbackHandlerName)
	}

	if callbackErr == nil {
		t.Error("Expected error in callback")
	}
}

func TestHandlerWrapper_ExponentialBackoff(t *testing.T) {
	attempts := 0
	handler := func(event events.DomainEvent) error {
		attempts++
		if attempts < 3 {
			return errors.New("handler error")
		}
		return nil
	}

	retryConfig := RetryConfig{
		MaxRetries:        3,
		InitialDelay:      10 * time.Millisecond,
		MaxDelay:          100 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	wrapper := NewHandlerWrapper(handler, retryConfig, "test-handler")

	start := time.Now()
	event := newMockEvent("TestEvent")
	err := wrapper.Execute(event)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}

	// Should have taken at least initial delay + (initial delay * 2) = 30ms
	// But allow some margin for execution time
	if duration < 20*time.Millisecond {
		t.Errorf("Expected exponential backoff delay, but duration was only %v", duration)
	}
}

func TestEventBus_SubscribeWithRetry(t *testing.T) {
	eventBus := NewEventBus()

	retryConfig := RetryConfig{
		MaxRetries:        2,
		InitialDelay:      10 * time.Millisecond,
		MaxDelay:          100 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	handler := func(event events.DomainEvent) error {
		return nil
	}

	eventBus.SubscribeWithRetry("TestEvent", handler, retryConfig, "custom-handler")

	if !eventBus.HasSubscribers("TestEvent") {
		t.Error("Expected subscribers for TestEvent")
	}

	if eventBus.GetSubscriberCount("TestEvent") != 1 {
		t.Errorf("Expected 1 subscriber, got %d", eventBus.GetSubscriberCount("TestEvent"))
	}
}
