package database

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// KeyDBService represents a service that interacts with KeyDB
type KeyDBService interface {
	// Health returns the status of the KeyDB connection
	Health() map[string]string
	// Close terminates the KeyDB connection
	Close() error
	// Client returns the underlying Redis client
	Client() *redis.Client
	// Set stores a key-value pair in KeyDB
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get retrieves a value from KeyDB by key
	Get(ctx context.Context, key string) (string, error)
	// Update updates a key-value pair in KeyDB
	Update(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Delete removes a key-value pair from KeyDB
	Delete(ctx context.Context, key string) error
}

type keydbService struct {
	client *redis.Client
}

var (
	keydbHost     = getEnvOrDefault("KEYDB_HOST", "localhost")
	keydbPort     = getEnvOrDefault("KEYDB_PORT", "6379")
	keydbPassword = getEnvOrDefault("KEYDB_PASSWORD", "")
	keydbInstance *keydbService
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewKeyDB creates a new KeyDB service instance
func NewKeyDB() KeyDBService {
	// Reuse Connection
	if keydbInstance != nil {
		return keydbInstance
	}

	client := redis.NewClient(&redis.Options{
		Addr:     keydbHost + ":" + keydbPort,
		Password: keydbPassword,
		DB:       0,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to KeyDB: " + err.Error())
	}

	keydbInstance = &keydbService{
		client: client,
	}

	return keydbInstance
}

// Health implements KeyDBService
func (s *keydbService) Health() map[string]string {
	status := make(map[string]string)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.client.Ping(ctx).Err(); err != nil {
		status["status"] = "error"
		status["message"] = "KeyDB connection error: " + err.Error()
	} else {
		status["status"] = "healthy"
		status["message"] = "KeyDB connection is healthy"
	}

	return status
}

// Close implements KeyDBService
func (s *keydbService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// Client implements KeyDBService
func (s *keydbService) Client() *redis.Client {
	return s.client
}

// Set implements KeyDBService
func (s *keydbService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

// Get implements KeyDBService
func (s *keydbService) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Update implements KeyDBService
func (s *keydbService) Update(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.Set(ctx, key, value, expiration)
}

// Delete implements KeyDBService
func (s *keydbService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
