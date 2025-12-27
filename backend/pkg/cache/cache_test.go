package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheService_NewCacheService(t *testing.T) {
	// This test requires Redis to be running
	// Skip if Redis is not available
	redisURL := "redis://localhost:6379"

	service, err := NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}

	require.NotNil(t, service)
	defer service.Close()
}

func TestCacheService_SetAndGet(t *testing.T) {
	redisURL := "redis://localhost:6379"

	service, err := NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer service.Close()

	// Test Set and Get
	key := "test:key:1"
	value := []byte("test value")

	err = service.Set(key, value, 1*time.Minute)
	require.NoError(t, err)

	retrieved, err := service.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, retrieved)
}

func TestCacheService_SetJSONAndGetJSON(t *testing.T) {
	redisURL := "redis://localhost:6379"

	service, err := NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer service.Close()

	// Test SetJSON and GetJSON
	key := "test:json:1"
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	input := TestStruct{
		Name:  "test",
		Value: 42,
	}

	err = service.SetJSON(key, input, 1*time.Minute)
	require.NoError(t, err)

	var output TestStruct
	found, err := service.GetJSON(key, &output)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, input.Name, output.Name)
	assert.Equal(t, input.Value, output.Value)
}

func TestCacheService_Delete(t *testing.T) {
	redisURL := "redis://localhost:6379"

	service, err := NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer service.Close()

	// Set a value
	key := "test:delete:1"
	value := []byte("test value")
	err = service.Set(key, value, 1*time.Minute)
	require.NoError(t, err)

	// Verify it exists
	retrieved, err := service.Get(key)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Delete it
	err = service.Delete(key)
	require.NoError(t, err)

	// Verify it's gone
	retrieved, err = service.Get(key)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}
