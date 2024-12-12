package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

type Store struct {
	client *redis.Client
}

type Session struct {
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewStore() (*Store, error) {
	addr := os.Getenv("KEYDB_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to KeyDB: %v", err)
	}

	return &Store{client: client}, nil
}

func (s *Store) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	session := Session{
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Sessions expire after 24 hours
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %v", err)
	}

	sessionID := generateSessionID()
	err = s.client.Set(ctx, sessionKey(sessionID), sessionData, 24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session: %v", err)
	}

	return sessionID, nil
}

func (s *Store) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	data, err := s.client.Get(ctx, sessionKey(sessionID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}

	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(ctx, sessionID)
		return nil, nil
	}

	return &session, nil
}

func (s *Store) DeleteSession(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, sessionKey(sessionID)).Err()
}

func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:32]
}
