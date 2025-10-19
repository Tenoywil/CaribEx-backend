# ✅ SIWE Authentication Setup Complete

The CaribEX backend is now fully configured for **wagmi SIWE-based authentication**.

## 📦 What Was Added

### 1. Dependencies
- ✅ `github.com/spruceid/siwe-go` - SIWE message verification
- ✅ `github.com/redis/go-redis/v9` - Redis client for session storage

### 2. Domain Layer (`internal/domain/auth/`)
- ✅ `session.go` - Session and Nonce models
- ✅ `repository.go` - Session repository interface

### 3. Repository Layer (`internal/repository/redis/`)
- ✅ `session_repository.go` - Redis-based session storage with TTL

### 4. Use Case Layer (`internal/usecase/`)
- ✅ `auth_usecase.go` - SIWE verification, nonce generation, session management

### 5. Controller Layer (`internal/controller/`)
- ✅ `auth_controller.go` - HTTP handlers for auth endpoints

### 6. Middleware (`pkg/middleware/`)
- ✅ `auth.go` - Session validation middleware for protected routes

### 7. Infrastructure (`pkg/redis/`)
- ✅ `client.go` - Redis client initialization

### 8. Routes (`internal/routes/routes.go`)
- ✅ Auth endpoints added (`/v1/auth/nonce`, `/v1/auth/siwe`, `/v1/auth/me`, `/v1/auth/logout`)
- ✅ Protected routes configured with auth middleware
- ✅ Public product listing maintained

### 9. Configuration
- ✅ `SIWE_DOMAIN` added to config
- ✅ `.env.example` updated

### 10. Documentation
- ✅ `docs/SIWE_AUTH.md` - Complete authentication guide
- ✅ `docs/INTEGRATION_EXAMPLE.md` - Integration example with main.go

## 🚀 Quick Start

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` and set:
```bash
SIWE_DOMAIN=localhost:3000  # Your frontend domain (no protocol)
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 3. Start Redis

```bash
docker-compose up redis -d
```

### 4. Update main.go

Add the auth components to your main.go (see `docs/INTEGRATION_EXAMPLE.md` for complete example):

```go
// Initialize Redis
redis, err := redisClient.NewClient(cfg.Redis)

// Initialize repositories
sessionRepo := redisRepo.NewSessionRepository(redis)

// Initialize use cases
authUseCase := usecase.NewAuthUseCase(sessionRepo, userUseCase, cfg.Auth.SIWEDomain)

// Initialize controllers
authController := controller.NewAuthController(authUseCase)

// Setup routes with auth
routes.SetupRoutes(
    router,
    authController,
    authUseCase,
    userController,
    productController,
    walletController,
    cartController,
    orderController,
)
```

### 5. Start Server

```bash
make run-dev
```

## 📡 API Endpoints

### Public Auth Endpoints
- `GET /v1/auth/nonce` - Generate SIWE nonce
- `POST /v1/auth/siwe` - Authenticate with wallet signature

### Protected Auth Endpoints
- `GET /v1/auth/me` - Get current user
- `POST /v1/auth/logout` - Logout

### Protected Routes (Require Authentication)
- `/v1/users/*` - User management
- `/v1/wallet/*` - Wallet operations
- `/v1/cart/*` - Shopping cart
- `/v1/orders/*` - Order management
- `POST /v1/products` - Create product
- `PUT /v1/products/:id` - Update product
- `DELETE /v1/products/:id` - Delete product

### Public Routes
- `GET /v1/products` - List products
- `GET /v1/products/:id` - Get product
- `GET /v1/categories` - List categories

## 🔐 Authentication Flow

1. **Frontend requests nonce**: `GET /v1/auth/nonce`
2. **User signs message** with wallet (wagmi)
3. **Frontend sends signature**: `POST /v1/auth/siwe`
4. **Backend verifies** signature and creates session
5. **Session cookie** set automatically
6. **All subsequent requests** include cookie for authentication

## 🎯 Frontend Integration (wagmi)

```typescript
import { useSignMessage } from 'wagmi';
import { SiweMessage } from 'siwe';

// 1. Get nonce
const { nonce } = await fetch('/v1/auth/nonce').then(r => r.json());

// 2. Create and sign message
const message = new SiweMessage({
  domain: window.location.host,
  address: address,
  statement: 'Sign in to CaribEX',
  uri: window.location.origin,
  version: '1',
  chainId: chainId,
  nonce: nonce,
});

const signature = await signMessageAsync({ 
  message: message.prepareMessage() 
});

// 3. Authenticate
await fetch('/v1/auth/siwe', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include', // Important!
  body: JSON.stringify({
    message: message.prepareMessage(),
    signature: signature,
  }),
});
```

## 🧪 Testing

### Test Nonce Generation

```bash
curl http://localhost:8080/v1/auth/nonce
```

Expected response:
```json
{
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "Fri, 18 Oct 2025 19:00:00 GMT"
}
```

### Test Health Check

```bash
curl http://localhost:8080/healthz
```

## 📚 Documentation

- **Complete Auth Guide**: `docs/SIWE_AUTH.md`
- **Integration Example**: `docs/INTEGRATION_EXAMPLE.md`
- **API Documentation**: `docs/API.md`
- **Architecture**: `docs/ARCHITECTURE.md`

## ⚙️ Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SIWE_DOMAIN` | Frontend domain for SIWE | `localhost:3000` |
| `SESSION_SECRET` | Session encryption key | `change-me-in-production` |
| `SESSION_DURATION` | Session lifetime | `24h` |
| `REDIS_HOST` | Redis server host | `localhost` |
| `REDIS_PORT` | Redis server port | `6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `REDIS_DB` | Redis database number | `0` |

## 🔒 Security Features

- ✅ **SIWE Standard**: EIP-4361 compliant
- ✅ **Nonce Protection**: One-time use, 10-minute expiration
- ✅ **Session Management**: Redis-based with TTL
- ✅ **HTTP-Only Cookies**: XSS protection
- ✅ **Domain Verification**: Prevents phishing
- ✅ **Signature Verification**: Cryptographic authentication
- ✅ **Automatic Expiration**: Sessions expire after 24 hours

## 🛠️ Next Steps

1. **Install dependencies**: `go mod tidy`
2. **Update main.go**: Add auth components (see `docs/INTEGRATION_EXAMPLE.md`)
3. **Configure CORS**: Ensure `Access-Control-Allow-Credentials: true`
4. **Test locally**: Use cURL or frontend
5. **Deploy**: Update production config with secure values

## 🐛 Troubleshooting

### Redis Connection Error
- Ensure Redis is running: `docker-compose up redis -d`
- Check `REDIS_HOST` and `REDIS_PORT` in `.env`

### Domain Mismatch Error
- Set `SIWE_DOMAIN` to match frontend (without `http://` or `https://`)
- Example: `localhost:3000` not `http://localhost:3000`

### CORS Issues
- Set `ALLOWED_ORIGIN` to frontend URL
- Include `credentials: 'include'` in frontend requests
- Don't use wildcard `*` for origin with credentials

### Session Not Found
- Ensure cookies are being sent with requests
- Check `credentials: 'include'` in fetch options
- Verify session hasn't expired (24h default)

## 📞 Support

For issues or questions:
- Check `docs/SIWE_AUTH.md` for detailed documentation
- Review `docs/INTEGRATION_EXAMPLE.md` for integration help
- See GitHub Issues for known problems

---

**Status**: ✅ Ready for integration and testing
**Compatible with**: wagmi, RainbowKit, ConnectKit, Web3Modal
**Standards**: EIP-4361 (SIWE), EIP-191 (Signed Data)
