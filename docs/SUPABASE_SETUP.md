# Supabase Storage Setup Guide

This guide walks you through setting up Supabase Storage for product image uploads in the CaribEX Backend.

## Prerequisites

- A Supabase account (free tier is sufficient)
- Access to the Supabase dashboard

## Step-by-Step Setup

### 1. Create a Supabase Project

1. Go to [supabase.com](https://supabase.com)
2. Sign up or log in
3. Click **"New Project"**
4. Fill in the project details:
   - **Name**: CaribEX (or your preferred name)
   - **Database Password**: Choose a strong password (save this!)
   - **Region**: Select the closest region to your users
5. Click **"Create new project"**
6. Wait for the project to initialize (1-2 minutes)

### 2. Create a Storage Bucket

1. In your Supabase dashboard, navigate to **Storage** in the left sidebar
2. Click **"Create a new bucket"**
3. Enter the bucket details:
   - **Name**: `product-images`
   - **Public bucket**: ✅ **Enable** (required for public image access)
4. Click **"Create bucket"**

### 3. Configure Bucket Policies

By default, public buckets allow read access but restrict uploads. We need to configure policies to allow authenticated uploads.

#### Option A: Using the Supabase Dashboard (Recommended)

1. Click on the **"product-images"** bucket
2. Navigate to the **"Policies"** tab
3. Click **"New Policy"** and select **"Custom policy"**
4. Create the upload policy:
   - **Policy name**: `Allow authenticated uploads`
   - **Allowed operation**: `INSERT`
   - **Target roles**: `authenticated`
   - **USING expression**:
     ```sql
     bucket_id = 'product-images'
     ```
5. Click **"Review"** then **"Save policy"**

#### Option B: Using SQL Editor

1. Go to **SQL Editor** in the left sidebar
2. Run the following SQL:

```sql
-- Allow public read access (should be enabled by default for public buckets)
CREATE POLICY "Public Access"
ON storage.objects FOR SELECT
USING (bucket_id = 'product-images');

-- Allow authenticated uploads
CREATE POLICY "Authenticated uploads"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (bucket_id = 'product-images');

-- Optional: Allow users to update their own uploads
CREATE POLICY "Users can update their own uploads"
ON storage.objects FOR UPDATE
TO authenticated
USING (bucket_id = 'product-images')
WITH CHECK (bucket_id = 'product-images');

-- Optional: Allow users to delete their own uploads
CREATE POLICY "Users can delete their own uploads"
ON storage.objects FOR DELETE
TO authenticated
USING (bucket_id = 'product-images');
```

### 4. Get Your API Credentials

1. Navigate to **Settings** → **API** in the left sidebar
2. Copy the following values:
   - **Project URL**: `https://your-project-id.supabase.co`
   - **anon public** key: Long string starting with `eyJ...`

### 5. Configure Environment Variables

Add the credentials to your `.env` file:

```env
# Supabase Storage Configuration
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
SUPABASE_BUCKET=product-images
STORAGE_MAX_FILE_SIZE=5242880
```

**Important:** Replace the placeholders with your actual values:
- `your-project-id` → Your actual Supabase project ID
- `eyJ...` → Your actual anon/public API key

### 6. Verify the Setup

#### Test Upload via cURL

Once your backend is running, test the upload endpoint:

```bash
# First, authenticate to get a session cookie
# (Requires SIWE authentication flow)

# Then upload an image
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session-cookie" \
  -F "image=@test-image.jpg"
```

Expected response:
```json
{
  "url": "https://your-project-id.supabase.co/storage/v1/object/public/product-images/products/1697712345_test-image.jpg",
  "filename": "test-image.jpg"
}
```

#### Verify in Supabase Dashboard

1. Go to **Storage** → **product-images**
2. You should see a `products/` folder
3. Click into it to view uploaded images
4. Test the public URL by pasting it in a browser

## Configuration Options

### Adjust File Size Limit

Default: 5MB (5242880 bytes)

To change:
```env
STORAGE_MAX_FILE_SIZE=10485760  # 10MB
```

### Supported Image Types

The backend validates these MIME types:
- `image/jpeg`
- `image/jpg`
- `image/png`
- `image/gif`
- `image/webp`
- `image/svg+xml`

To add more types, edit `pkg/storage/storage.go`:

```go
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		// Add your types here
		"image/bmp",
	}
	// ...
}
```

## Folder Structure

Uploaded images are organized as:
```
product-images/
  └── products/
      ├── 1697712345_image1.jpg
      ├── 1697712346_image2.jpg
      └── 1697712347_image3.jpg
```

The timestamp prefix prevents filename collisions.

## Security Best Practices

### 1. Use Row Level Security (RLS)

For production, implement RLS to restrict access:

```sql
-- Only allow users to upload to their own folders
CREATE POLICY "Users own folder uploads"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (
  bucket_id = 'product-images' AND
  (storage.foldername(name))[1] = auth.uid()::text
);
```

Then modify the upload path in code:
```go
// In product_controller.go
userID := ctx.GetString("user_id")
folder := fmt.Sprintf("users/%s/products", userID)
url, err := c.storageService.UploadFile(ctx.Request.Context(), file, header, folder)
```

### 2. Enable Rate Limiting

Configure rate limits in Supabase:
1. Go to **Settings** → **API**
2. Scroll to **Rate Limiting**
3. Set limits for storage operations

### 3. Monitor Storage Usage

1. Go to **Settings** → **Billing**
2. Check your storage usage
3. Free tier includes 1GB storage
4. Consider upgrading if needed

### 4. Enable Content Moderation (Optional)

For production apps, consider:
- Client-side image validation (dimensions, size)
- Server-side malware scanning
- Content moderation APIs (AWS Rekognition, Google Cloud Vision)

## Troubleshooting

### Error: "Failed to upload file to storage"

**Possible causes:**
1. Invalid Supabase credentials
2. Bucket doesn't exist
3. Bucket is not public
4. Network connectivity issues

**Solutions:**
- Verify `SUPABASE_URL` and `SUPABASE_KEY` in `.env`
- Check bucket name matches exactly
- Ensure bucket is marked as **Public**
- Check Supabase service status

### Error: "Invalid file type"

**Cause:** File MIME type not in allowed list

**Solution:** Check the Content-Type header and ensure it's a supported image format

### Error: "File size exceeds maximum"

**Cause:** File larger than `STORAGE_MAX_FILE_SIZE`

**Solutions:**
- Compress the image
- Increase `STORAGE_MAX_FILE_SIZE` in `.env`
- Implement client-side compression

### Images Not Accessible

**Possible causes:**
1. Bucket is not public
2. RLS policies blocking access

**Solutions:**
- Make bucket public in Supabase dashboard
- Review and update RLS policies
- Check public URL format

### Slow Upload Performance

**Solutions:**
- Use a CDN (Supabase includes this)
- Compress images before upload
- Consider image optimization services
- Choose a closer Supabase region

## Advanced Features

### Image Transformations

Supabase Storage supports image transformations:

```
https://your-project.supabase.co/storage/v1/render/image/public/product-images/products/image.jpg?width=400&height=300
```

Parameters:
- `width`: Resize width
- `height`: Resize height
- `quality`: JPEG quality (1-100)

### Custom Storage Backend

To use AWS S3, Google Cloud Storage, or other providers:

1. Implement the `storage.Service` interface
2. Update initialization in `cmd/api-server/main.go`

Example:
```go
type CustomStorage struct {
    // Your implementation
}

func (s *CustomStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
    // Your upload logic
}
```

## Cost Estimation

**Supabase Free Tier:**
- Storage: 1GB
- Bandwidth: 2GB/month
- No credit card required

**Pro Plan ($25/month):**
- Storage: 8GB included
- Bandwidth: 50GB/month included
- Additional storage: $0.125/GB
- Additional bandwidth: $0.09/GB

**Example calculation for 1000 products:**
- Average 3 images per product, 500KB each
- Total storage: ~1.5GB
- Monthly views: 100k (assuming 150KB per view)
- Bandwidth: ~15GB
- Cost: **Free tier sufficient** or **Pro plan covers it**

## Next Steps

1. ✅ Complete this setup
2. Test uploads via the API
3. Integrate with your frontend
4. Monitor usage and performance
5. Consider enabling image CDN/optimization
6. Set up automated backups (optional)

## Support Resources

- [Supabase Storage Documentation](https://supabase.com/docs/guides/storage)
- [Supabase Storage API Reference](https://supabase.com/docs/reference/javascript/storage-from-upload)
- [CaribEX API Documentation](./API_REFERENCE_IMAGE_UPLOAD.md)
- [Supabase Community Discord](https://discord.supabase.com)

---

**Last Updated:** 2025-10-19
