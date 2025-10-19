package storage

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestIsValidImageType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"JPEG", "image/jpeg", true},
		{"JPG", "image/jpg", true},
		{"PNG", "image/png", true},
		{"GIF", "image/gif", true},
		{"WebP", "image/webp", true},
		{"SVG", "image/svg+xml", true},
		{"PDF", "application/pdf", false},
		{"Text", "text/plain", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidImageType(tt.contentType); got != tt.want {
				t.Errorf("isValidImageType(%q) = %v, want %v", tt.contentType, got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "Simple filename",
			filename: "test.jpg",
			want:     "test",
		},
		{
			name:     "Filename with spaces",
			filename: "my product image.jpg",
			want:     "my_product_image",
		},
		{
			name:     "Filename with special characters",
			filename: "product@2023!.jpg",
			want:     "product_2023_",
		},
		{
			name:     "Already clean filename",
			filename: "product_123.jpg",
			want:     "product_123",
		},
		{
			name:     "Unicode characters",
			filename: "productñ©.jpg",
			want:     "product__",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeFilename(tt.filename); got != tt.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestExtractPathFromURL(t *testing.T) {
	baseURL := "https://project.supabase.co"
	bucket := "product-images"

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Full URL",
			url:  "https://project.supabase.co/storage/v1/object/public/product-images/products/123_image.jpg",
			want: "products/123_image.jpg",
		},
		{
			name: "Path only",
			url:  "products/123_image.jpg",
			want: "products/123_image.jpg",
		},
		{
			name: "Nested path",
			url:  "https://project.supabase.co/storage/v1/object/public/product-images/products/subfolder/image.jpg",
			want: "products/subfolder/image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPathFromURL(tt.url, baseURL, bucket); got != tt.want {
				t.Errorf("extractPathFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestNewSupabaseStorage(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Valid configuration",
			config: Config{
				URL:         "https://test.supabase.co",
				Key:         "test-key",
				Bucket:      "test-bucket",
				MaxFileSize: 5242880,
			},
			wantErr: false,
		},
		{
			name: "Default max file size",
			config: Config{
				URL:    "https://test.supabase.co",
				Key:    "test-key",
				Bucket: "test-bucket",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewSupabaseStorage(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSupabaseStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if storage == nil {
					t.Error("NewSupabaseStorage() returned nil storage")
				}
				if storage.maxFileSize == 0 {
					t.Error("NewSupabaseStorage() maxFileSize should be set to default")
				}
			}
		})
	}
}

func TestGetPublicURL(t *testing.T) {
	storage, _ := NewSupabaseStorage(Config{
		URL:    "https://project.supabase.co",
		Key:    "test-key",
		Bucket: "product-images",
	})

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Simple path",
			path: "products/image.jpg",
			want: "https://project.supabase.co/storage/v1/object/public/product-images/products/image.jpg",
		},
		{
			name: "Nested path",
			path: "products/subfolder/image.jpg",
			want: "https://project.supabase.co/storage/v1/object/public/product-images/products/subfolder/image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storage.GetPublicURL(tt.path); got != tt.want {
				t.Errorf("GetPublicURL(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

// Helper function to create a mock multipart file
func createMockFile(t *testing.T, filename, contentType string, content []byte) (multipart.File, *multipart.FileHeader) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(int64(len(content)) + 1024)
	if err != nil {
		t.Fatal(err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no files in form")
	}

	file, err := files[0].Open()
	if err != nil {
		t.Fatal(err)
	}

	return file, files[0]
}

func TestUploadFile_FileSizeValidation(t *testing.T) {
	storage, _ := NewSupabaseStorage(Config{
		URL:         "https://test.supabase.co",
		Key:         "test-key",
		Bucket:      "test-bucket",
		MaxFileSize: 100, // Very small for testing
	})

	// Create a file larger than max size
	largeContent := make([]byte, 200)
	file, header := createMockFile(t, "large.jpg", "image/jpeg", largeContent)
	defer file.Close()

	// This should fail due to size
	_, err := storage.UploadFile(nil, file, header, "test")
	if err == nil {
		t.Error("UploadFile() should fail for oversized file")
	}
}

func TestUploadFile_ContentTypeValidation(t *testing.T) {
	storage, _ := NewSupabaseStorage(Config{
		URL:         "https://test.supabase.co",
		Key:         "test-key",
		Bucket:      "test-bucket",
		MaxFileSize: 5242880,
	})

	tests := []struct {
		name        string
		contentType string
		wantErr     bool
	}{
		{"Invalid PDF", "application/pdf", true},
		{"Invalid text", "text/plain", true},
		{"Empty type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := []byte("fake content")
			file, header := createMockFile(t, "test.jpg", tt.contentType, content)
			defer file.Close()

			// Should fail at content type validation before actual upload
			_, err := storage.UploadFile(nil, file, header, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() with %s error = %v, wantErr %v", tt.contentType, err, tt.wantErr)
			}
		})
	}
}
