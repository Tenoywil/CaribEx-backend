# Product Image Upload - Quick Reference

## Overview

This implementation adds S3-based image upload functionality to the CaribEX Backend using Supabase Storage. Images are stored in a public bucket and served via CDN.

## Quick Start

### 1. Configure Supabase

```bash
# Add to .env file
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
SUPABASE_BUCKET=product-images
STORAGE_MAX_FILE_SIZE=5242880
```

See [SUPABASE_SETUP.md](./SUPABASE_SETUP.md) for detailed setup instructions.

### 2. Test Upload

```bash
# Upload a single image
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session" \
  -F "image=@product.jpg"

# Response
{
  "url": "https://your-project.supabase.co/storage/v1/object/public/product-images/products/1697712345_product.jpg",
  "filename": "product.jpg"
}
```

## API Endpoints

### Image Upload Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/products/upload-image` | POST | ✓ | Upload single product image |
| `/v1/products/multipart` | POST | ✓ | Create product with images |
| `/v1/products` | POST | ✓ | Create product (JSON with URLs) |

### Product Management Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/products` | GET | ✗ | List products (paginated) |
| `/v1/products/:id` | GET | ✗ | Get product details |
| `/v1/products/:id` | PUT | ✓ | Update product |
| `/v1/products/:id` | DELETE | ✓ | Delete product |
| `/v1/categories` | GET | ✗ | List categories |

## Usage Examples

### Upload Image (JavaScript)

```javascript
const uploadImage = async (file) => {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch('http://localhost:8080/v1/products/upload-image', {
    method: 'POST',
    credentials: 'include',
    body: formData,
  });

  const { url } = await response.json();
  return url;
};
```

### Create Product with Images (JavaScript)

```javascript
const createProduct = async (productData, imageFiles) => {
  const formData = new FormData();
  formData.append('title', productData.title);
  formData.append('description', productData.description);
  formData.append('price', productData.price);
  formData.append('quantity', productData.quantity);
  
  imageFiles.forEach(file => {
    formData.append('images', file);
  });

  const response = await fetch('http://localhost:8080/v1/products/multipart', {
    method: 'POST',
    credentials: 'include',
    body: formData,
  });

  return await response.json();
};
```

### Create Product with Pre-uploaded URLs (JavaScript)

```javascript
const createProductWithURLs = async (productData) => {
  const response = await fetch('http://localhost:8080/v1/products', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(productData),
  });

  return await response.json();
};
```

## File Specifications

### Supported Formats

- JPEG/JPG (`image/jpeg`)
- PNG (`image/png`)
- GIF (`image/gif`)
- WebP (`image/webp`)
- SVG (`image/svg+xml`)

### Limits

- **Max File Size**: 5MB (configurable via `STORAGE_MAX_FILE_SIZE`)
- **Max Form Size**: 10MB
- **Multiple Files**: Unlimited (within form size limit)

## Security Features

- ✅ Authentication required for uploads
- ✅ File type validation (images only)
- ✅ File size validation
- ✅ Filename sanitization (prevents path traversal)
- ✅ Unique filenames (timestamp-based)
- ✅ Public read-only bucket access

## Workflow Options

### Option 1: Upload Then Create

1. Upload image → Get URL
2. Create product with URL

**Best for:** Image preview before product creation

```javascript
// Step 1: Upload
const imageUrl = await uploadImage(file);

// Step 2: Create product
const product = await createProductWithURLs({
  title: 'Product',
  price: 29.99,
  quantity: 100,
  images: [imageUrl]
});
```

### Option 2: Upload and Create Together

1. Create product with images in single request

**Best for:** Simpler flow, faster UX

```javascript
const product = await createProduct(
  { title: 'Product', price: 29.99, quantity: 100 },
  [imageFile1, imageFile2]
);
```

## Error Handling

### Common Errors

| Error | Status | Cause | Solution |
|-------|--------|-------|----------|
| `image file is required` | 400 | Missing file | Include file in request |
| `invalid file type` | 400 | Unsupported format | Use supported image format |
| `file size exceeds maximum` | 413 | File too large | Compress image or increase limit |
| `Unauthorized` | 401 | Not authenticated | Authenticate via SIWE |
| `failed to upload file` | 500 | Storage error | Check Supabase config |

### Error Response Format

```json
{
  "error": "Error message describing the issue"
}
```

## Testing

### Unit Tests

```bash
# Run storage package tests
go test -v ./pkg/storage/...

# Run all tests
make test
```

### Manual Testing

```bash
# 1. Start server
make run-dev

# 2. Authenticate (get session cookie)
# ... SIWE flow ...

# 3. Test upload
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session" \
  -F "image=@test.jpg" \
  -v
```

## Performance Tips

1. **Compress images** before upload (client-side)
2. Use **WebP format** for better compression
3. Implement **lazy loading** for product images
4. Use Supabase **image transformations** for thumbnails:
   ```
   https://your-project.supabase.co/storage/v1/render/image/public/product-images/products/image.jpg?width=200
   ```

## Troubleshooting

### Build Issues

```bash
# Ensure dependencies are up to date
go mod tidy
go mod download

# Rebuild
make build
```

### Upload Failures

1. Check Supabase credentials in `.env`
2. Verify bucket is public
3. Check bucket name matches exactly
4. Review Supabase dashboard for errors

### File Not Accessible

1. Ensure bucket is marked as **Public**
2. Check RLS policies
3. Verify URL format is correct

## Documentation

- **[API_REFERENCE_IMAGE_UPLOAD.md](./API_REFERENCE_IMAGE_UPLOAD.md)** - Complete API documentation
- **[SUPABASE_SETUP.md](./SUPABASE_SETUP.md)** - Supabase configuration guide

## Architecture

```
Frontend (Next.js)
    ↓
Product Controller
    ↓
Storage Service (pkg/storage)
    ↓
Supabase Storage (S3-compatible)
    ↓
Public CDN URL
```

## Dependencies Added

```go
github.com/supabase-community/storage-go v0.8.1
github.com/aws/aws-sdk-go-v2 v1.39.3
github.com/aws/aws-sdk-go-v2/config v1.31.13
github.com/aws/aws-sdk-go-v2/service/s3 v1.88.5
github.com/aws/aws-sdk-go-v2/credentials v1.18.17
```

## Future Enhancements

- [ ] Image optimization/compression
- [ ] Multiple bucket support
- [ ] Image deletion when product is deleted
- [ ] Image cropping/resizing
- [ ] Content moderation
- [ ] Cloudflare Images integration
- [ ] Duplicate detection

## Support

For issues or questions:
- Check the documentation links above
- Review Supabase Storage docs
- Open an issue on GitHub

---

**Created:** 2025-10-19  
**Status:** ✅ Production Ready
