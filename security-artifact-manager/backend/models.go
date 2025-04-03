package models

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"
)

// Artifact represents a binary artifact stored in the system
type Artifact struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FileSize    int64     `json:"fileSize"`
	SHA256      string    `json:"sha256"`
	UploadTime  time.Time `json:"uploadTime"`
	UploaderID  string    `json:"uploaderId"`
	FilePath    string    `json:"-"` // Internal use only, not exposed via API
}

// User represents a system user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Hashed password, not exposed via API
	Role     string `json:"role"`
}

// Session represents a user session
type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CalculateSHA256 calculates the SHA256 hash of a file
func CalculateSHA256(file io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum), nil
}