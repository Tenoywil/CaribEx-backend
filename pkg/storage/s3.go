package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// S3Service handles file uploads to S3-compatible storage
type S3Service struct {
	uploader *s3manager.Uploader
	s3Client *s3.S3
	bucket   string
}

// NewS3Service creates a new S3 service
func NewS3Service(uploader *s3manager.Uploader, s3Client *s3.S3, bucket string) *S3Service {
	return &S3Service{
		uploader: uploader,
		s3Client: s3Client,
		bucket:   bucket,
	}
}

// UploadFileResult contains the result of a file upload
type UploadFileResult struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// UploadFiles uploads multiple files to S3 and returns their keys and URLs
func (s *S3Service) UploadFiles(files []*multipart.FileHeader, prefix string) ([]UploadFileResult, error) {
	var results []UploadFileResult

	for _, fileHeader := range files {
		result, err := s.UploadFile(fileHeader, prefix)
		if err != nil {
			log.Error().
				Err(err).
				Str("filename", fileHeader.Filename).
				Msg("failed to upload file")
			return nil, fmt.Errorf("failed to upload %s: %w", fileHeader.Filename, err)
		}
		results = append(results, result)
	}

	log.Info().
		Int("count", len(results)).
		Str("prefix", prefix).
		Msg("successfully uploaded files to S3")

	return results, nil
}

// UploadFile uploads a single file to S3
func (s *S3Service) UploadFile(fileHeader *multipart.FileHeader, prefix string) (UploadFileResult, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return UploadFileResult{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Generate unique key
	ext := path.Ext(fileHeader.Filename)
	id := uuid.New().String()
	key := fmt.Sprintf("%s/%s%s", strings.TrimSuffix(prefix, "/"), id, ext)

	// Detect content type
	contentType, err := detectContentType(file)
	if err != nil {
		return UploadFileResult{}, fmt.Errorf("failed to detect content type: %w", err)
	}

	// Reset file reader after content type detection
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return UploadFileResult{}, fmt.Errorf("failed to reset file reader: %w", err)
		}
	}

	// Validate content type (security)
	if !isAllowedContentType(contentType) {
		return UploadFileResult{}, fmt.Errorf("content type not allowed: %s", contentType)
	}

	log.Debug().
		Str("key", key).
		Str("content_type", contentType).
		Int64("size", fileHeader.Size).
		Str("filename", fileHeader.Filename).
		Msg("uploading file to S3")

	// Upload to S3
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         aws.String("private"), // Keep files private, use presigned URLs
	})
	if err != nil {
		return UploadFileResult{}, fmt.Errorf("failed to upload to S3: %w", err)
	}

	log.Info().
		Str("key", key).
		Str("location", result.Location).
		Msg("file uploaded successfully")

	return UploadFileResult{
		Key:         key,
		URL:         result.Location,
		Size:        fileHeader.Size,
		ContentType: contentType,
	}, nil
}

// GeneratePresignedURL generates a presigned URL for downloading a file
func (s *S3Service) GeneratePresignedURL(key string, expirationMinutes int) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	// Generate presigned URL
	//convert to time
	expirationDuration := time.Duration(expirationMinutes) * time.Minute
	urlStr, err := req.Presign(expirationDuration)
	if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to generate presigned URL")
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.Debug().
		Str("key", key).
		Int("expiration_minutes", expirationMinutes).
		Msg("generated presigned URL")

	return urlStr, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to delete file from S3")
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Info().
		Str("key", key).
		Msg("file deleted successfully")

	return nil
}

// DeleteFiles deletes multiple files from S3
func (s *S3Service) DeleteFiles(keys []string) error {
	var objects []*s3.ObjectIdentifier
	for _, key := range keys {
		objects = append(objects, &s3.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	if len(objects) == 0 {
		return nil
	}

	_, err := s.s3Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Int("count", len(keys)).
			Msg("failed to delete files from S3")
		return fmt.Errorf("failed to delete files: %w", err)
	}

	log.Info().
		Int("count", len(keys)).
		Msg("files deleted successfully")

	return nil
}

// detectContentType detects the content type of a file
func detectContentType(file multipart.File) (string, error) {
	// Read first 512 bytes for detection
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buf[:n])
	return contentType, nil
}

// isAllowedContentType checks if the content type is allowed
func isAllowedContentType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"image/svg+xml":   true,
		"application/pdf": true,
		"video/mp4":       true,
		"video/mpeg":      true,
		"video/quicktime": true,
	}

	return allowedTypes[contentType]
}

// getDuration converts minutes to time.Duration
func getDuration(minutes int) int64 {
	return int64(minutes * 60) // seconds
}
