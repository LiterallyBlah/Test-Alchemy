package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartKeyDBContainer() (func(context.Context) error, string, string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "eqalpha/keydb:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get container port: %v", err)
	}

	cleanup := func(ctx context.Context) error {
		return container.Terminate(ctx)
	}

	return cleanup, host, port.Port(), nil
}

func TestNewKeyDB(t *testing.T) {
	cleanup, host, port, err := mustStartKeyDBContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup(context.Background())

	// Set environment variables for the test
	t.Setenv("KEYDB_HOST", host)
	t.Setenv("KEYDB_PORT", port)

	// Test creating a new KeyDB service
	service := NewKeyDB()
	assert.NotNil(t, service)
	assert.NotNil(t, service.Client())

	// Test singleton pattern
	service2 := NewKeyDB()
	assert.Equal(t, service, service2)
}

func TestKeyDBHealth(t *testing.T) {
	cleanup, host, port, err := mustStartKeyDBContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup(context.Background())

	t.Setenv("KEYDB_HOST", host)
	t.Setenv("KEYDB_PORT", port)

	service := NewKeyDB()
	health := service.Health()
	assert.Equal(t, "healthy", health["status"])
}

func TestKeyDBOperations(t *testing.T) {
	cleanup, host, port, err := mustStartKeyDBContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup(context.Background())

	t.Setenv("KEYDB_HOST", host)
	t.Setenv("KEYDB_PORT", port)

	service := NewKeyDB()
	client := service.Client()
	ctx := context.Background()

	// Test Set operation
	err = client.Set(ctx, "test_key", "test_value", time.Minute).Err()
	assert.NoError(t, err)

	// Test Get operation
	val, err := client.Get(ctx, "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)

	// Test non-existent key
	_, err = client.Get(ctx, "non_existent_key").Result()
	assert.Equal(t, redis.Nil, err)
}

func TestKeyDBClose(t *testing.T) {
	cleanup, host, port, err := mustStartKeyDBContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup(context.Background())

	t.Setenv("KEYDB_HOST", host)
	t.Setenv("KEYDB_PORT", port)

	service := NewKeyDB()
	err = service.Close()
	assert.NoError(t, err)
}
