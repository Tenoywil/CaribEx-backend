package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	storagego "github.com/supabase-community/storage-go"
)

// Service defines the interface for storage operations
type Service interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error)
	DeleteFile(ctx context.Context, path string) error
	GetPublicURL(path string) string
}

// SupabaseStorage implements the Service interface using Supabase Storage
type SupabaseStorage struct {
	client      *storagego.Client
	bucket      string
	baseURL     string
	maxFileSize int64
}

// Config holds the configuration for Supabase Storage
type Config struct {
	URL         string
	Key         string
	Bucket      string
	MaxFileSize int64
}

// NewSupabaseStorage creates a new Supabase storage service
func NewSupabaseStorage(cfg Config) (*SupabaseStorage, error) {
	client := storagego.NewClient(cfg.URL, cfg.Key, nil)

	// Set default max file size if not specified (5MB)
	maxFileSize := cfg.MaxFileSize
	if maxFileSize == 0 {
		maxFileSize = 5 * 1024 * 1024 // 5MB default
	}

	return &SupabaseStorage{
		client:      client,
		bucket:      cfg.Bucket,
		baseURL:     cfg.URL,
		maxFileSize: maxFileSize,
	}, nil
}

// UploadFile uploads a file to Supabase Storage and returns the public URL
func (s *SupabaseStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Validate file size
	if header.Size > s.maxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxFileSize)
	}

	// Validate file type (images only)
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type: %s. Only images are allowed", contentType)
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s/%d_%s%s", folder, timestamp, sanitizeFilename(header.Filename), ext)

	// Upload to Supabase Storage
	_, err = s.client.UploadFile(s.bucket, filename, bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Get public URL
	publicURL := s.GetPublicURL(filename)

	return publicURL, nil
}

// DeleteFile deletes a file from Supabase Storage
func (s *SupabaseStorage) DeleteFile(ctx context.Context, path string) error {
	// Extract the path from the full URL if needed
	cleanPath := extractPathFromURL(path, s.baseURL, s.bucket)

	_, err := s.client.RemoveFile(s.bucket, []string{cleanPath})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetPublicURL returns the public URL for a file
func (s *SupabaseStorage) GetPublicURL(path string) string {
	// Construct public URL for Supabase Storage
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.baseURL, s.bucket, path)
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// sanitizeFilename removes unsafe characters from filename
func sanitizeFilename(filename string) string {
	// Remove file extension for sanitization
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// Replace spaces and special characters
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, name)

	// Truncate if too long
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}

// extractPathFromURL extracts the file path from a full URL
func extractPathFromURL(url, baseURL, bucket string) string {
	// If it's already a path, return as is
	if !strings.HasPrefix(url, "http") {
		return url
	}

	// Extract path from public URL format
	prefix := fmt.Sprintf("%s/storage/v1/object/public/%s/", baseURL, bucket)
	return strings.TrimPrefix(url, prefix)
}
