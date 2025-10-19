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

## Product Listing Features

### Pagination

All list endpoints support pagination with validation:

```javascript
// Get page 2 with 50 items
const response = await fetch('/v1/products?page=2&page_size=50');
```

- **page**: Page number (min: 1, default: 1)
- **page_size**: Items per page (min: 1, max: 100, default: 20)

### Sorting

Products can be sorted by multiple fields:

```javascript
// Sort by price ascending
const response = await fetch('/v1/products?sort_by=price&sort_order=asc');

// Sort by newest first (default)
const response = await fetch('/v1/products?sort_by=created_at&sort_order=desc');
```

**Available sort fields:**
- `created_at` - When product was created (default)
- `updated_at` - When product was last updated
- `price` - Product price
- `title` - Product title (alphabetical)

**Sort order:**
- `asc` - Ascending (lowest to highest)
- `desc` - Descending (highest to lowest, default)

### Filtering

Combine filters for precise results:

```javascript
// Category filter
const response = await fetch('/v1/products?category_id=123e4567-...');

// Search filter
const response = await fetch('/v1/products?search=coffee');

// Combined: Search in category, sorted by price
const response = await fetch('/v1/products?category_id=123e4567-...&search=premium&sort_by=price&sort_order=asc');
```

### Category Information

Products now include full category details:

```json
{
  "id": "550e8400-...",
  "title": "Premium Coffee",
  "category_id": "123e4567-...",
  "category": {
    "id": "123e4567-...",
    "name": "Food & Beverages"
  }
}
```

### Complete Example

```javascript
// Fetch products with all options
const fetchProducts = async (page = 1, categoryId = null) => {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: '20',
    sort_by: 'created_at',
    sort_order: 'desc'
  });
  
  if (categoryId) {
    params.append('category_id', categoryId);
  }
  
  const response = await fetch(`/v1/products?${params}`);
  const data = await response.json();
  
  return {
    products: data.products,
    total: data.total,
    totalPages: data.total_pages,
    currentPage: data.page
  };
};
```

## Security Features

- ✅ Authentication required for uploads
- ✅ File type validation (images only)
- ✅ File size validation
- ✅ Filename sanitization (prevents path traversal)
- ✅ Unique filenames (timestamp-based)
- ✅ Public read-only bucket access
- ✅ Pagination limits enforced (max 100 items per page)

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
