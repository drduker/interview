package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"secure-artifact-manager/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrSessionExpired     = errors.New("session expired")
)

// AuthService handles authentication operations
type AuthService struct {
	// In a real implementation, this would use a database
	users    map[string]*models.User
	sessions map[string]*models.Session
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{
		users:    make(map[string]*models.User),
		sessions: make(map[string]*models.Session),
	}
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(username, password, role string) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// Generate ID
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	
	// Create user
	user := &models.User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}
	
	// Store user
	s.users[id] = user
	
	return user, nil
}

// Login creates a new session for a user
func (s *AuthService) Login(username, password string) (*models.Session, error) {
	// Find user by username
	var user *models.User
	for _, u := range s.users {
		if u.Username == username {
			user = u
			break
		}
	}
	
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	
	// Generate token
	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	
	// Create session
	session := &models.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	// Store session
	s.sessions[token] = session
	
	return session, nil
}

// Logout invalidates a session
func (s *AuthService) Logout(token string) error {
	delete(s.sessions, token)
	return nil
}

// ValidateSession validates a session token
func (s *AuthService) ValidateSession(token string) (*models.User, error) {
	session, exists := s.sessions[token]
	if !exists {
		return nil, ErrInvalidCredentials
	}
	
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return nil, ErrSessionExpired
	}
	
	user, exists := s.users[session.UserID]
	if !exists {
		delete(s.sessions, token)
		return nil, ErrInvalidCredentials
	}
	
	return user, nil
}

// Helper function to generate a random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}