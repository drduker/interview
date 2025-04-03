package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"secure-artifact-manager/models"
)

// ArtifactService handles artifact operations
type ArtifactService struct {
	StoragePath string
}

// NewArtifactService creates a new artifact service
func NewArtifactService(storagePath string) (*ArtifactService, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	
	return &ArtifactService{
		StoragePath: storagePath,
	}, nil
}

// StoreArtifact saves an artifact to storage
func (s *ArtifactService) StoreArtifact(name, description string, fileSize int64, reader io.Reader, uploaderID string) (*models.Artifact, error) {
	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	
	// Create file path
	filePath := filepath.Join(s.StoragePath, id)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	// Create a TeeReader to calculate hash while writing file
	hashReader := io.TeeReader(reader, file)
	
	// Calculate SHA256
	sha256, err := models.CalculateSHA256(hashReader)
	if err != nil {
		// Clean up file if hash calculation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}
	
	// Create artifact
	artifact := &models.Artifact{
		ID:          id,
		Name:        name,
		Description: description,
		FileSize:    fileSize,
		SHA256:      sha256,
		UploadTime:  time.Now(),
		UploaderID:  uploaderID,
		FilePath:    filePath,
	}
	
	return artifact, nil
}

// GetArtifact retrieves an artifact by ID
func (s *ArtifactService) GetArtifact(id string) (*models.Artifact, error) {
	// In a real implementation, this would query a database
	// For this example, we'll return an error
	return nil, fmt.Errorf("artifact not found")
}

// DeleteArtifact removes an artifact by ID
func (s *ArtifactService) DeleteArtifact(id string) error {
	// In a real implementation, this would delete from database and storage
	// For this example, we'll return an error
	return fmt.Errorf("artifact not found")
}

// OpenArtifactForRead opens an artifact file for reading
func (s *ArtifactService) OpenArtifactForRead(id string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.StoragePath, id)
	return os.Open(filePath)
}

// ListArtifacts returns all artifacts
func (s *ArtifactService) ListArtifacts() ([]*models.Artifact, error) {
	// In a real implementation, this would query a database
	// For this example, we'll return an empty list
	return []*models.Artifact{}, nil
}

// Helper function to generate a random ID
func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}