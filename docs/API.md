# CaribX Backend API Documentation

## Overview

CaribX is a blockchain-backed money transfer and marketplace platform for Jamaica and the Caribbean. This document describes the REST API endpoints available for integration.

**Base URL**: `http://localhost:8080/v1`

**Authentication**: Most endpoints require authentication via SIWE (Sign-In With Ethereum) session cookies or JWT tokens.

---

## Authentication Endpoints

### Get Nonce for SIWE

Generate a nonce for Sign-In With Ethereum authentication.

**Endpoint**: `GET /v1/auth/nonce`

**Response**:
```json
{
  "nonce": "random-nonce-string",
  "expires_at": "2025-10-18T19:00:00Z"
}
```

### Authenticate with SIWE

Verify Ethereum signature and create/update user session.

**Endpoint**: `POST /v1/auth/siwe`

**Request Body**:
```json
{
  "message": "Sign-in message with nonce",
  "signature": "0x...",
  "wallet_address": "0x..."
}
```

**Response**:
```json
{
  "user": {
    "id": "uuid",
    "username": "user123",
    "wallet_address": "0x...",
    "role": "customer"
  },
  "session_token": "session-cookie-set"
}
```

### Get Current User

Retrieve authenticated user information.

**Endpoint**: `GET /v1/auth/me`

**Headers**: `Cookie: session=...`

**Response**:
```json
{
  "id": "uuid",
  "username": "user123",
  "wallet_address": "0x...",
  "role": "customer",
  "created_at": "2025-10-18T10:00:00Z"
}
```

---

## Wallet Endpoints

### Get Wallet Summary

Retrieve wallet balance and details.

**Endpoint**: `GET /v1/wallet`

**Headers**: `Cookie: session=...`

**Response**:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "balance": 1000.50,
  "currency": "JAM",
  "updated_at": "2025-10-18T12:00:00Z"
}
```

### Send Funds

Initiate an outgoing transfer.

**Endpoint**: `POST /v1/wallet/send`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "recipient_address": "0x...",
  "amount": 100.00,
  "reference": "Payment for Order #123"
}
```

**Response**:
```json
{
  "transaction_id": "uuid",
  "status": "pending",
  "amount": 100.00,
  "recipient": "0x...",
  "created_at": "2025-10-18T12:00:00Z"
}
```

### Get Transaction History

Retrieve wallet transaction ledger.

**Endpoint**: `GET /v1/wallet/transactions?page=1&page_size=20`

**Headers**: `Cookie: session=...`

**Response**:
```json
{
  "transactions": [
    {
      "id": "uuid",
      "type": "credit",
      "amount": 100.00,
      "reference": "Deposit",
      "status": "success",
      "created_at": "2025-10-18T12:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20
}
```

---

## Product Endpoints

### List Products

Browse marketplace products with optional filters.

**Endpoint**: `GET /v1/products?page=1&page_size=20&category_id=uuid&search=keyword`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20)
- `category_id` (optional): Filter by category
- `search` (optional): Search in title/description

**Response**:
```json
{
  "products": [
    {
      "id": "uuid",
      "seller_id": "uuid",
      "title": "Product Name",
      "description": "Product description",
      "price": 99.99,
      "quantity": 10,
      "images": ["url1", "url2"],
      "category_id": "uuid",
      "is_active": true,
      "created_at": "2025-10-18T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### Get Product Details

Retrieve single product information.

**Endpoint**: `GET /v1/products/:id`

**Response**:
```json
{
  "id": "uuid",
  "seller_id": "uuid",
  "title": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "quantity": 10,
  "images": ["url1", "url2"],
  "category_id": "uuid",
  "is_active": true,
  "created_at": "2025-10-18T10:00:00Z",
  "updated_at": "2025-10-18T11:00:00Z"
}
```

### Create Product (Seller Only)

Create a new product listing.

**Endpoint**: `POST /v1/products`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "title": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "quantity": 10,
  "images": ["url1", "url2"],
  "category_id": "uuid"
}
```

**Response**:
```json
{
  "id": "uuid",
  "seller_id": "uuid",
  "title": "Product Name",
  "is_active": true,
  "created_at": "2025-10-18T12:00:00Z"
}
```

### Update Product (Seller Only)

Update an existing product.

**Endpoint**: `PUT /v1/products/:id`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "title": "Updated Product Name",
  "price": 89.99,
  "quantity": 15
}
```

### Delete Product (Seller Only)

Delete a product listing.

**Endpoint**: `DELETE /v1/products/:id`

**Headers**: `Cookie: session=...`

**Response**: `204 No Content`

---

## Cart Endpoints

### Get Cart

Retrieve current user's active cart.

**Endpoint**: `GET /v1/cart`

**Headers**: `Cookie: session=...`

**Response**:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "status": "active",
  "total": 199.98,
  "items": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "quantity": 2,
      "price": 99.99,
      "product": {
        "title": "Product Name",
        "images": ["url1"]
      }
    }
  ]
}
```

### Add Item to Cart

Add or update a product in the cart.

**Endpoint**: `POST /v1/cart/items`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "product_id": "uuid",
  "quantity": 1
}
```

**Response**:
```json
{
  "id": "uuid",
  "cart_id": "uuid",
  "product_id": "uuid",
  "quantity": 1,
  "price": 99.99
}
```

### Update Cart Item

Modify quantity of an item in cart.

**Endpoint**: `PUT /v1/cart/items/:id`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "quantity": 3
}
```

### Remove Cart Item

Remove an item from cart.

**Endpoint**: `DELETE /v1/cart/items/:id`

**Headers**: `Cookie: session=...`

**Response**: `204 No Content`

---

## Order Endpoints

### Create Order (Checkout)

Convert cart to order and process payment.

**Endpoint**: `POST /v1/orders`

**Headers**: `Cookie: session=...`

**Request Body**:
```json
{
  "cart_id": "uuid",
  "payment_method": "wallet"
}
```

**Response**:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "status": "pending",
  "total": 199.98,
  "payment_ref": "tx-uuid",
  "created_at": "2025-10-18T12:00:00Z"
}
```

### Get User Orders

Retrieve order history.

**Endpoint**: `GET /v1/orders?page=1&page_size=20`

**Headers**: `Cookie: session=...`

**Response**:
```json
{
  "orders": [
    {
      "id": "uuid",
      "status": "completed",
      "total": 199.98,
      "created_at": "2025-10-18T12:00:00Z",
      "items": [
        {
          "product_id": "uuid",
          "quantity": 2,
          "price": 99.99
        }
      ]
    }
  ],
  "total": 10,
  "page": 1,
  "page_size": 20
}
```

---

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input parameters",
    "details": {
      "field": "email",
      "reason": "invalid format"
    }
  }
}
```

**Common HTTP Status Codes**:
- `200 OK`: Successful request
- `201 Created`: Resource created
- `204 No Content`: Successful deletion
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

## Rate Limiting

API endpoints are rate-limited to ensure fair usage:

- **Authentication endpoints**: 10 requests per minute per IP
- **Read endpoints**: 100 requests per minute per user
- **Write endpoints**: 30 requests per minute per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1697654400
```

---

## Pagination

List endpoints support pagination with the following query parameters:

- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20, max: 100)

Paginated responses include metadata:
```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```
