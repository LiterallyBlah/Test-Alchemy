package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"TestAlchemy/internal/database"
	"TestAlchemy/internal/models"
	"TestAlchemy/internal/session"
	"github.com/google/uuid"
)

type UserService struct {
	db      database.Service
	session *session.Store
}

func NewUserService(db database.Service, sessionStore *session.Store) *UserService {
	return &UserService{
		db:      db,
		session: sessionStore,
	}
}

type RegisterUserInput struct {
	Email    string
	Password string
}

type LoginUserInput struct {
	Email    string
	Password string
}

func (s *UserService) ValidateRegistration(input RegisterUserInput) error {
	// Validate email
	if input.Email == "" {
		return errors.New("email is required")
	}
	parsedEmail, err := mail.ParseAddress(input.Email)
	if err != nil {
		return errors.New("invalid email format")
	}
	if parsedEmail.Name != "" {
		return errors.New("email cannot contain a name")
	}
	if strings.Contains(parsedEmail.Address, ",") {
		return errors.New("email cannot contain a comma")
	}
	if strings.Contains(parsedEmail.Address, ";") {
		return errors.New("email cannot contain a semicolon")
	}

	// Additional domain validation
	parts := strings.Split(parsedEmail.Address, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	domain := parts[1]
	if !strings.Contains(domain, ".") {
		return errors.New("invalid email domain format")
	}
	domainParts := strings.Split(domain, ".")
	if len(domainParts[len(domainParts)-1]) < 2 {
		return errors.New("invalid top-level domain")
	}

	// Validate password
	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(input.Password) > 256 {
		return errors.New("password must be at most 256 characters long")
	}
	if !strings.ContainsAny(input.Password, "0123456789") {
		return errors.New("password must contain at least one number")
	}
	if strings.ToLower(input.Password) == input.Password {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !strings.ContainsAny(input.Password, "!@#$%^&*()_+-=[]{}|;:'\",<.>/?") {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func (s *UserService) RegisterUser(ctx context.Context, input RegisterUserInput) error {
	if err := s.ValidateRegistration(input); err != nil {
		return err
	}

	// Create new user with UUID
	user := &models.User{
		UserID: uuid.New(),
		Email:  input.Email,
	}
	if err := user.HashPassword(input.Password); err != nil {
		return err
	}

	err := s.db.Create(ctx, user)
	if err != nil {
		return errors.New("registration failed")
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, input LoginUserInput) (string, error) {
	var user models.User
	err := s.db.Read(ctx, &user, "email = ?", input.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.ValidatePassword(input.Password) {
		return "", errors.New("invalid email or password")
	}

	sessionID, err := s.session.CreateSession(ctx, user.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return sessionID, nil
}
