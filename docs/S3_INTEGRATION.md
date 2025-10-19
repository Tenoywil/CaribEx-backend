# S3-Compatible Storage Integration Complete

## Summary
Successfully integrated S3-compatible storage (Supabase/MinIO/AWS S3) into CaribEX backend with enhanced logging and security features.

## What Was Added

### 1. S3 Service Package (`pkg/storage/s3.go`)
A complete S3 service implementation with:
- **File Upload**: Single and multiple file uploads with automatic UUID naming
- **Content Type Detection**: Automatic MIME type detection with validation
- **Security**: 
  - Content type whitelist (images, PDFs, videos only)
  - Private ACL by default
  - File size validation support
- **Presigned URLs**: Generate temporary download URLs (no public access needed)
- **File Deletion**: Single and batch file deletion
- **Rich Logging**: Detailed logs for all operations using zerolog

**Allowed Content Types:**
- Images: jpeg, png, gif, webp, svg+xml
- Documents: PDF
- Videos: mp4, mpeg, quicktime

### 2. Configuration Fields Added (`pkg/config/config.go`)
New S3-specific environment variables:
```bash
SUPABASE_S3_ACCESS_KEY_ID       # S3 access key ID
SUPABASE_S3_SECRET_ACCESS_KEY   # S3 secret access key
SUPABASE_STORAGE_URL            # S3 endpoint URL
SUPABASE_REGION                 # AWS region (e.g., us-east-1)
```

### 3. Main Application Integration (`cmd/api-server/main.go`)
- AWS SDK v1 imported (aws-sdk-go)
- S3 session initialization with:
  - Static credentials
  - Custom endpoint support (Supabase/MinIO)
  - Path-style S3 addressing
- S3Service instance created and ready for controllers

### 4. Enhanced CORS Middleware (`pkg/middleware/cors.go`)
Improved with:
- **Richer Logging**: Request timestamp, remote IP, request ID, User-Agent, Referer
- **Scheme-Insensitive Matching**: Automatically handles http/https differences
- **Better Headers**: Added `Access-Control-Expose-Headers` and `Access-Control-Max-Age`
- **Explicit Preflight Handling**: Returns 403 for disallowed origins (easier debugging)
- **Security**: Credentials header only for explicit origin matches

## Usage Examples

### Upload Files in a Controller
```go
import "github.com/Tenoywil/CaribEx-backend/pkg/storage"

// In controller with s3Service injected
func (c *ProductController) UploadImages(ctx *gin.Context) {
    // Get multipart files
    form, _ := ctx.MultipartForm()
    files := form.File["images"]
    
    // Upload to S3
    results, err := c.s3Service.UploadFiles(files, "products/images")
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // results contains: Key, URL, Size, ContentType for each file
    ctx.JSON(200, gin.H{"files": results})
}
```

### Generate Presigned Download URL
```go
// Generate a URL valid for 15 minutes
url, err := s3Service.GeneratePresignedURL("products/images/abc-123.jpg", 15)
if err != nil {
    return err
}
// Return url to client for secure download
```

### Delete Files
```go
// Delete single file
err := s3Service.DeleteFile("products/images/old-image.jpg")

// Delete multiple files
keys := []string{"file1.jpg", "file2.png"}
err := s3Service.DeleteFiles(keys)
```

## Environment Setup

Add to your `.env`:
```bash
# Supabase S3-Compatible Storage
SUPABASE_S3_ACCESS_KEY_ID=your_access_key_here
SUPABASE_S3_SECRET_ACCESS_KEY=your_secret_key_here
SUPABASE_STORAGE_URL=https://your-project.supabase.co/storage/v1/s3
SUPABASE_REGION=us-east-1
SUPABASE_BUCKET=your-bucket-name
```

For local MinIO:
```bash
SUPABASE_S3_ACCESS_KEY_ID=minioadmin
SUPABASE_S3_SECRET_ACCESS_KEY=minioadmin
SUPABASE_STORAGE_URL=http://localhost:9000
SUPABASE_REGION=us-east-1
SUPABASE_BUCKET=uploads
```

## Security Features

1. **Private Files**: All uploads use ACL=private (no public access)
2. **Presigned URLs**: Generate temporary URLs for downloads
3. **Content Type Validation**: Only allowed file types can be uploaded
4. **Unique Filenames**: UUIDs prevent filename collisions and path traversal
5. **Tenant Isolation**: Use prefix parameter (e.g., `uploads/user-123/`)
6. **Size Limits**: Configure max file size in Gin middleware

## Next Steps to Use S3 in Controllers

1. **Update Controller Constructor**: Add `s3Service *storage.S3Service` parameter
2. **Update Routes**: Pass `s3Service` when creating controllers
3. **Add Upload Endpoint**: Use `s3Service.UploadFiles()` in handler
4. **Store Keys in DB**: Save S3 keys (not URLs) in product/user tables
5. **Generate URLs**: Use presigned URLs when returning data to clients

## Benefits

✅ **Production-Ready**: Supports Supabase, AWS S3, MinIO  
✅ **Secure**: Private files, presigned URLs, content validation  
✅ **Observable**: Rich logging with zerolog  
✅ **Scalable**: S3-compatible storage handles millions of files  
✅ **Flexible**: Easy to switch between providers  

## Files Modified

1. `pkg/storage/s3.go` (NEW)
2. `pkg/config/config.go` (UPDATED)
3. `cmd/api-server/main.go` (UPDATED)
4. `pkg/middleware/cors.go` (ENHANCED)
5. `go.mod` (UPDATED - added aws-sdk-go)

## Dependencies Added

```bash
go get github.com/aws/aws-sdk-go
```

**Note**: AWS SDK v1 is deprecated but still works. Consider migrating to aws-sdk-go-v2 in the future for better performance and features.
