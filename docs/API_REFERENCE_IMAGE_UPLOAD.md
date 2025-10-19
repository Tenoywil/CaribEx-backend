# Product Image Upload API Reference

This document provides detailed information about the product image upload endpoints in the CaribEX Backend API.

## Table of Contents
- [Overview](#overview)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Upload Product Image](#upload-product-image)
  - [Create Product with Images (Multipart)](#create-product-with-images-multipart)
  - [Create Product (JSON)](#create-product-json)
  - [List Products](#list-products)
  - [Get Product](#get-product)
  - [Update Product](#update-product)
  - [Delete Product](#delete-product)
  - [Get Categories](#get-categories)
- [Error Responses](#error-responses)
- [Storage Configuration](#storage-configuration)

---

## Overview

The CaribEX Backend provides a robust file upload system for product images using Supabase Storage (S3-compatible). All uploaded images are:
- Stored in a public S3-based bucket
- Validated for type and size
- Given unique, timestamped filenames to prevent collisions
- Accessible via public URLs

**Base URL:** `http://localhost:8080/v1`

---

## Authentication

All write operations (POST, PUT, DELETE) require authentication via SIWE (Sign-In With Ethereum). Include the session cookie obtained from the authentication flow.

### Authentication Flow:
1. **GET** `/v1/auth/nonce` - Get a nonce for signing
2. **POST** `/v1/auth/siwe` - Submit signed message to authenticate
3. Use the returned session cookie for subsequent requests

---

## Endpoints

### Upload Product Image

Upload a single product image to storage and receive the public URL.

**Endpoint:** `POST /v1/products/upload-image`

**Authentication:** Required

**Content-Type:** `multipart/form-data`

**Request Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| image | file | Yes | Image file (JPEG, PNG, GIF, WebP, SVG) |

**cURL Example:**
```bash
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session-cookie" \
  -F "image=@/path/to/image.jpg"
```

**Success Response (200 OK):**
```json
{
  "url": "https://your-project.supabase.co/storage/v1/object/public/product-images/products/1697712345_product_image.jpg",
  "filename": "product_image.jpg"
}
```

**Error Responses:**
- `400 Bad Request` - Missing image file or invalid file type
- `401 Unauthorized` - Missing or invalid authentication
- `413 Payload Too Large` - File exceeds maximum size (5MB default)
- `500 Internal Server Error` - Storage upload failed

---

### Create Product with Images (Multipart)

Create a new product listing with multiple images uploaded simultaneously.

**Endpoint:** `POST /v1/products/multipart`

**Authentication:** Required

**Content-Type:** `multipart/form-data`

**Request Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Product title |
| description | string | No | Product description |
| price | number | Yes | Product price (decimal) |
| quantity | integer | Yes | Available quantity |
| category_id | string | No | Category UUID |
| images | file[] | No | Multiple image files |

**cURL Example:**
```bash
curl -X POST http://localhost:8080/v1/products/multipart \
  -H "Cookie: session=your-session-cookie" \
  -F "title=Premium Coffee Beans" \
  -F "description=High-quality Arabica coffee from Blue Mountains" \
  -F "price=29.99" \
  -F "quantity=100" \
  -F "category_id=123e4567-e89b-12d3-a456-426614174000" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg"
```

**JavaScript/Fetch Example:**
```javascript
const formData = new FormData();
formData.append('title', 'Premium Coffee Beans');
formData.append('description', 'High-quality Arabica coffee');
formData.append('price', '29.99');
formData.append('quantity', '100');
formData.append('category_id', '123e4567-e89b-12d3-a456-426614174000');

// Add multiple images
formData.append('images', imageFile1);
formData.append('images', imageFile2);

const response = await fetch('http://localhost:8080/v1/products/multipart', {
  method: 'POST',
  credentials: 'include', // Include cookies
  body: formData
});

const product = await response.json();
```

**Success Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "seller_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Premium Coffee Beans",
  "description": "High-quality Arabica coffee from Blue Mountains",
  "price": 29.99,
  "quantity": 100,
  "images": [
    "https://your-project.supabase.co/storage/v1/object/public/product-images/products/1697712345_image1.jpg",
    "https://your-project.supabase.co/storage/v1/object/public/product-images/products/1697712346_image2.jpg"
  ],
  "category_id": "123e4567-e89b-12d3-a456-426614174000",
  "is_active": true,
  "created_at": "2025-10-19T12:30:00Z",
  "updated_at": "2025-10-19T12:30:00Z"
}
```

---

### Create Product (JSON)

Create a new product with pre-uploaded image URLs.

**Endpoint:** `POST /v1/products`

**Authentication:** Required

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "title": "Premium Coffee Beans",
  "description": "High-quality Arabica coffee from Blue Mountains",
  "price": 29.99,
  "quantity": 100,
  "category_id": "123e4567-e89b-12d3-a456-426614174000",
  "images": [
    "https://your-project.supabase.co/storage/v1/object/public/product-images/products/1697712345_image1.jpg"
  ]
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/v1/products \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session-cookie" \
  -d '{
    "title": "Premium Coffee Beans",
    "description": "High-quality Arabica coffee",
    "price": 29.99,
    "quantity": 100,
    "category_id": "123e4567-e89b-12d3-a456-426614174000",
    "images": ["https://your-project.supabase.co/storage/.../image.jpg"]
  }'
```

**Success Response (201 Created):** Same as multipart endpoint

---

### List Products

Retrieve a paginated list of products with optional filters.

**Endpoint:** `GET /v1/products`

**Authentication:** Not required

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | No | 1 | Page number (min: 1) |
| page_size | integer | No | 20 | Items per page (min: 1, max: 100) |
| category_id | string | No | - | Filter by category UUID |
| search | string | No | - | Search in title/description |
| sort_by | string | No | created_at | Sort field: created_at, updated_at, price, title |
| sort_order | string | No | desc | Sort order: asc, desc |

**cURL Example:**
```bash
# Basic list
curl -X GET "http://localhost:8080/v1/products"

# With filters and sorting
curl -X GET "http://localhost:8080/v1/products?page=1&page_size=10&category_id=123e4567-e89b-12d3-a456-426614174000&sort_by=price&sort_order=asc"

# Search by keyword
curl -X GET "http://localhost:8080/v1/products?search=coffee&sort_by=created_at&sort_order=desc"
```

**Success Response (200 OK):**
```json
{
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "seller_id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Premium Coffee Beans",
      "description": "High-quality Arabica coffee",
      "price": 29.99,
      "quantity": 100,
      "images": [
        "https://your-project.supabase.co/storage/.../image.jpg"
      ],
      "category_id": "123e4567-e89b-12d3-a456-426614174000",
      "category": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Food & Beverages"
      },
      "is_active": true,
      "created_at": "2025-10-19T12:30:00Z",
      "updated_at": "2025-10-19T12:30:00Z"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 10,
  "total_pages": 5
}
```

**Notes:**
- Products now include nested `category` object with `id` and `name`
- Sorting is case-insensitive
- Invalid sort fields default to `created_at DESC`
- Page size is capped at 100 items
- Search is case-insensitive and searches both title and description

---

### Get Product

Retrieve a single product by ID.

**Endpoint:** `GET /v1/products/:id`

**Authentication:** Not required

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Product UUID |

**cURL Example:**
```bash
curl -X GET http://localhost:8080/v1/products/550e8400-e29b-41d4-a716-446655440000
```

**Success Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "seller_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Premium Coffee Beans",
  "description": "High-quality Arabica coffee",
  "price": 29.99,
  "quantity": 100,
  "images": [
    "https://your-project.supabase.co/storage/.../image.jpg"
  ],
  "category_id": "123e4567-e89b-12d3-a456-426614174000",
  "category": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Food & Beverages"
  },
  "is_active": true,
  "created_at": "2025-10-19T12:30:00Z",
  "updated_at": "2025-10-19T12:30:00Z"
}
```

**Notes:**
- Product includes nested `category` object with full details
- Category is null if the product has no category assigned

**Error Response (404 Not Found):**
```json
{
  "error": "product not found"
}
```

---

### Update Product

Update an existing product.

**Endpoint:** `PUT /v1/products/:id`

**Authentication:** Required (must be the seller)

**Content-Type:** `application/json`

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Product UUID |

**Request Body:**
```json
{
  "title": "Updated Product Title",
  "description": "Updated description",
  "price": 34.99,
  "quantity": 75,
  "images": [
    "https://your-project.supabase.co/storage/.../new-image.jpg"
  ],
  "category_id": "123e4567-e89b-12d3-a456-426614174000",
  "is_active": true
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/v1/products/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session-cookie" \
  -d '{
    "title": "Updated Product Title",
    "price": 34.99,
    "quantity": 75
  }'
```

**Success Response (200 OK):** Returns updated product object

---

### Delete Product

Delete a product listing.

**Endpoint:** `DELETE /v1/products/:id`

**Authentication:** Required (must be the seller)

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Product UUID |

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/v1/products/550e8400-e29b-41d4-a716-446655440000 \
  -H "Cookie: session=your-session-cookie"
```

**Success Response (204 No Content):** Empty body

---

### Get Categories

Retrieve all product categories.

**Endpoint:** `GET /v1/categories`

**Authentication:** Not required

**cURL Example:**
```bash
curl -X GET http://localhost:8080/v1/categories
```

**Success Response (200 OK):**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Food & Beverages"
  },
  {
    "id": "223e4567-e89b-12d3-a456-426614174000",
    "name": "Electronics"
  }
]
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 413 Payload Too Large
```json
{
  "error": "file size exceeds maximum allowed size of 5242880 bytes"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error message"
}
```

---

## Storage Configuration

### Supported Image Types
- JPEG/JPG (`image/jpeg`)
- PNG (`image/png`)
- GIF (`image/gif`)
- WebP (`image/webp`)
- SVG (`image/svg+xml`)

### File Size Limits
- Maximum file size: **5MB** (configurable via `STORAGE_MAX_FILE_SIZE`)
- Maximum form size: **10MB**

### Environment Variables

Add the following to your `.env` file:

```env
# Supabase Storage Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
SUPABASE_BUCKET=product-images
STORAGE_MAX_FILE_SIZE=5242880  # 5MB in bytes
```

### Setting up Supabase Storage

1. **Create a Supabase project** at [supabase.com](https://supabase.com)

2. **Create a storage bucket:**
   - Go to Storage in your Supabase dashboard
   - Create a new bucket named `product-images`
   - Make it **public** for read access

3. **Configure bucket policies:**
   ```sql
   -- Allow public read access
   CREATE POLICY "Public Access"
   ON storage.objects FOR SELECT
   USING (bucket_id = 'product-images');

   -- Allow authenticated uploads
   CREATE POLICY "Authenticated uploads"
   ON storage.objects FOR INSERT
   WITH CHECK (bucket_id = 'product-images');
   ```

4. **Get your credentials:**
   - Copy your project URL (e.g., `https://abcdefgh.supabase.co`)
   - Copy your anon/public API key from Settings > API

---

## Workflow Examples

### Workflow 1: Upload Image Then Create Product

```bash
# Step 1: Upload image
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session" \
  -F "image=@product.jpg"

# Response: { "url": "https://...", "filename": "..." }

# Step 2: Create product with the URL
curl -X POST http://localhost:8080/v1/products \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session" \
  -d '{
    "title": "My Product",
    "price": 19.99,
    "quantity": 50,
    "images": ["https://..."]
  }'
```

### Workflow 2: Create Product with Images in One Request

```bash
curl -X POST http://localhost:8080/v1/products/multipart \
  -H "Cookie: session=your-session" \
  -F "title=My Product" \
  -F "price=19.99" \
  -F "quantity=50" \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg"
```

---

## Integration with Frontend (React/Next.js)

### Example: Upload with React

```javascript
import React, { useState } from 'react';

function ProductImageUpload() {
  const [uploading, setUploading] = useState(false);
  const [imageUrl, setImageUrl] = useState('');

  const handleFileUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    setUploading(true);
    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await fetch('http://localhost:8080/v1/products/upload-image', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      const data = await response.json();
      setImageUrl(data.url);
      console.log('Uploaded image URL:', data.url);
    } catch (error) {
      console.error('Upload failed:', error);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        accept="image/*"
        onChange={handleFileUpload}
        disabled={uploading}
      />
      {uploading && <p>Uploading...</p>}
      {imageUrl && <img src={imageUrl} alt="Uploaded" style={{ maxWidth: '200px' }} />}
    </div>
  );
}
```

### Example: Create Product with Images

```javascript
import React, { useState } from 'react';

function CreateProduct() {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: '',
    quantity: '',
    category_id: '',
  });
  const [images, setImages] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const data = new FormData();
    data.append('title', formData.title);
    data.append('description', formData.description);
    data.append('price', formData.price);
    data.append('quantity', formData.quantity);
    data.append('category_id', formData.category_id);

    // Add all selected images
    images.forEach(image => {
      data.append('images', image);
    });

    try {
      const response = await fetch('http://localhost:8080/v1/products/multipart', {
        method: 'POST',
        credentials: 'include',
        body: data,
      });

      const product = await response.json();
      console.log('Product created:', product);
    } catch (error) {
      console.error('Failed to create product:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Title"
        value={formData.title}
        onChange={(e) => setFormData({...formData, title: e.target.value})}
        required
      />
      <textarea
        placeholder="Description"
        value={formData.description}
        onChange={(e) => setFormData({...formData, description: e.target.value})}
      />
      <input
        type="number"
        step="0.01"
        placeholder="Price"
        value={formData.price}
        onChange={(e) => setFormData({...formData, price: e.target.value})}
        required
      />
      <input
        type="number"
        placeholder="Quantity"
        value={formData.quantity}
        onChange={(e) => setFormData({...formData, quantity: e.target.value})}
        required
      />
      <input
        type="file"
        multiple
        accept="image/*"
        onChange={(e) => setImages(Array.from(e.target.files))}
      />
      <button type="submit">Create Product</button>
    </form>
  );
}
```

---

## Testing

### Manual Testing with cURL

```bash
# 1. Get authentication nonce
NONCE=$(curl -s http://localhost:8080/v1/auth/nonce | jq -r '.nonce')

# 2. Authenticate (requires SIWE signature - use actual wallet)
# ... (SIWE authentication flow)

# 3. Upload an image
curl -X POST http://localhost:8080/v1/products/upload-image \
  -H "Cookie: session=your-session" \
  -F "image=@test-image.jpg" \
  -v

# 4. Create product with uploaded image
curl -X POST http://localhost:8080/v1/products \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session" \
  -d '{
    "title": "Test Product",
    "price": 9.99,
    "quantity": 10,
    "images": ["https://your-project.supabase.co/storage/.../test-image.jpg"]
  }'
```

---

## Best Practices

1. **Image Optimization**: Compress images before upload to reduce storage costs and improve loading times
2. **Multiple Images**: Upload multiple product angles for better user experience
3. **Error Handling**: Always handle upload failures gracefully in your frontend
4. **Progress Indicators**: Show upload progress for better UX
5. **Image Validation**: Validate image dimensions and aspect ratios on the frontend before upload
6. **Cleanup**: Consider implementing image cleanup when products are deleted
7. **CDN**: Supabase Storage includes CDN, but consider additional optimization for high-traffic scenarios

---

## Support

For issues or questions:
- Check the [main README](../README.md)
- Review [Supabase Storage documentation](https://supabase.com/docs/guides/storage)
- Open an issue on GitHub

---

**Last Updated:** 2025-10-19
**API Version:** v1
