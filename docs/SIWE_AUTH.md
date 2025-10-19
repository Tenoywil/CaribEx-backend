# SIWE Authentication Setup Guide

## Overview

CaribEX backend uses **Sign-In With Ethereum (SIWE)** for Web3 authentication, compatible with wagmi on the frontend. This document explains the authentication flow and integration.

## Architecture

### Components

1. **Auth Domain** (`internal/domain/auth/`)
   - `Session`: User session model
   - `Nonce`: SIWE nonce model
   - `SessionRepository`: Interface for session storage

2. **Session Repository** (`internal/repository/redis/`)
   - Redis-based session and nonce storage
   - Automatic expiration handling
   - TTL-based cleanup

3. **Auth Use Case** (`internal/usecase/auth_usecase.go`)
   - Nonce generation
   - SIWE message verification
   - Session management
   - User creation/retrieval

4. **Auth Controller** (`internal/controller/auth_controller.go`)
   - HTTP handlers for auth endpoints
   - Cookie-based session management

5. **Auth Middleware** (`pkg/middleware/auth.go`)
   - Session validation
   - User context injection
   - Protected route enforcement

## Authentication Flow

### 1. Frontend Requests Nonce

```typescript
// Using wagmi
const response = await fetch('http://localhost:8080/v1/auth/nonce');
const { nonce, expires_at } = await response.json();
```

**Endpoint**: `GET /v1/auth/nonce`

**Response**:
```json
{
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "Fri, 18 Oct 2025 19:00:00 GMT"
}
```

### 2. Frontend Signs Message with Wallet

```typescript
import { useSignMessage } from 'wagmi';
import { SiweMessage } from 'siwe';

const { signMessageAsync } = useSignMessage();

// Create SIWE message
const message = new SiweMessage({
  domain: window.location.host,
  address: address,
  statement: 'Sign in to CaribEX',
  uri: window.location.origin,
  version: '1',
  chainId: chainId,
  nonce: nonce,
});

const preparedMessage = message.prepareMessage();
const signature = await signMessageAsync({ message: preparedMessage });
```

### 3. Frontend Sends Signature to Backend

```typescript
const response = await fetch('http://localhost:8080/v1/auth/siwe', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include', // Important for cookies
  body: JSON.stringify({
    message: preparedMessage,
    signature: signature,
  }),
});

const { user, session_id } = await response.json();
```

**Endpoint**: `POST /v1/auth/siwe`

**Request Body**:
```json
{
  "message": "localhost:3000 wants you to sign in with your Ethereum account:\n0x...",
  "signature": "0x..."
}
```

**Response**:
```json
{
  "user": {
    "id": "uuid",
    "username": "user_0x123456",
    "wallet_address": "0x...",
    "role": "customer"
  },
  "session_id": "session-uuid"
}
```

**Sets Cookie**: `session_id=<session-uuid>; Path=/; HttpOnly`

### 4. Authenticated Requests

All subsequent requests include the session cookie automatically:

```typescript
const response = await fetch('http://localhost:8080/v1/wallet', {
  credentials: 'include',
});
```

### 5. Get Current User

```typescript
const response = await fetch('http://localhost:8080/v1/auth/me', {
  credentials: 'include',
});

const user = await response.json();
```

**Endpoint**: `GET /v1/auth/me`

**Response**:
```json
{
  "user_id": "uuid",
  "wallet_address": "0x..."
}
```

### 6. Logout

```typescript
await fetch('http://localhost:8080/v1/auth/logout', {
  method: 'POST',
  credentials: 'include',
});
```

**Endpoint**: `POST /v1/auth/logout`

## Configuration

### Environment Variables

Add to `.env`:

```bash
# SIWE Configuration
SIWE_DOMAIN=localhost:3000  # Frontend domain (without protocol)

# Session Configuration
SESSION_SECRET=your-secret-key-change-in-production
SESSION_DURATION=24h

# Redis (required for session storage)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Production Settings

For production, update:

```bash
SIWE_DOMAIN=yourdomain.com
SESSION_SECRET=<strong-random-secret>
```

Also update cookie settings in `auth_controller.go`:
- Set `secure: true` (requires HTTPS)
- Set appropriate `domain` value

## API Endpoints

### Public Endpoints

- `GET /v1/auth/nonce` - Generate nonce
- `POST /v1/auth/siwe` - Authenticate with signature

### Protected Endpoints (Require Authentication)

- `GET /v1/auth/me` - Get current user
- `POST /v1/auth/logout` - Logout
- `GET /v1/wallet` - Get wallet
- `POST /v1/wallet/send` - Send funds
- `GET /v1/cart` - Get cart
- `POST /v1/orders` - Create order
- All user, cart, order, and wallet endpoints

### Public Read Endpoints

- `GET /v1/products` - List products
- `GET /v1/products/:id` - Get product
- `GET /v1/categories` - List categories

## Frontend Integration (wagmi)

### Install Dependencies

```bash
npm install wagmi viem siwe
```

### Configure wagmi

```typescript
import { WagmiConfig, createConfig, configureChains } from 'wagmi';
import { mainnet, polygon } from 'wagmi/chains';
import { publicProvider } from 'wagmi/providers/public';

const { chains, publicClient } = configureChains(
  [mainnet, polygon],
  [publicProvider()]
);

const config = createConfig({
  autoConnect: true,
  publicClient,
});

function App() {
  return (
    <WagmiConfig config={config}>
      {/* Your app */}
    </WagmiConfig>
  );
}
```

### Create Auth Hook

```typescript
import { useAccount, useSignMessage } from 'wagmi';
import { SiweMessage } from 'siwe';

export function useAuth() {
  const { address, chainId } = useAccount();
  const { signMessageAsync } = useSignMessage();

  const login = async () => {
    // 1. Get nonce
    const nonceRes = await fetch('http://localhost:8080/v1/auth/nonce');
    const { nonce } = await nonceRes.json();

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

    const preparedMessage = message.prepareMessage();
    const signature = await signMessageAsync({ message: preparedMessage });

    // 3. Authenticate
    const authRes = await fetch('http://localhost:8080/v1/auth/siwe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        message: preparedMessage,
        signature: signature,
      }),
    });

    return await authRes.json();
  };

  const logout = async () => {
    await fetch('http://localhost:8080/v1/auth/logout', {
      method: 'POST',
      credentials: 'include',
    });
  };

  return { login, logout };
}
```

### Usage in Component

```typescript
import { useAccount } from 'wagmi';
import { useAuth } from './hooks/useAuth';

function LoginButton() {
  const { address, isConnected } = useAccount();
  const { login, logout } = useAuth();

  if (!isConnected) {
    return <button>Connect Wallet</button>;
  }

  return (
    <div>
      <p>Connected: {address}</p>
      <button onClick={login}>Sign In</button>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

## Security Considerations

### Backend

1. **HTTPS in Production**: Always use HTTPS in production
2. **Secure Cookies**: Set `secure: true` for cookies in production
3. **CORS Configuration**: Configure allowed origins properly
4. **Session Expiration**: Sessions expire after 24 hours by default
5. **Nonce Expiration**: Nonces expire after 10 minutes
6. **One-Time Nonces**: Nonces are deleted after use

### Frontend

1. **Credentials**: Always include `credentials: 'include'` in fetch requests
2. **HTTPS**: Use HTTPS in production
3. **Domain Matching**: Ensure SIWE domain matches your frontend domain
4. **Error Handling**: Handle authentication errors gracefully

## Testing

### Manual Testing with cURL

```bash
# 1. Get nonce
curl http://localhost:8080/v1/auth/nonce

# 2. Sign message with wallet (use MetaMask or similar)

# 3. Authenticate
curl -X POST http://localhost:8080/v1/auth/siwe \
  -H "Content-Type: application/json" \
  -d '{
    "message": "...",
    "signature": "0x..."
  }' \
  -c cookies.txt

# 4. Make authenticated request
curl http://localhost:8080/v1/auth/me \
  -b cookies.txt
```

## Troubleshooting

### "Domain mismatch" Error

- Ensure `SIWE_DOMAIN` matches your frontend domain
- Frontend domain should not include protocol (http:// or https://)
- Example: `localhost:3000` not `http://localhost:3000`

### "Invalid or expired nonce" Error

- Nonces expire after 10 minutes
- Each nonce can only be used once
- Request a new nonce if authentication fails

### "Session expired" Error

- Sessions expire after 24 hours (configurable)
- User needs to sign in again
- Frontend should handle 401 responses and redirect to login

### CORS Issues

- Ensure frontend origin is in `ALLOWED_ORIGIN` environment variable
- Include `credentials: 'include'` in all fetch requests
- Backend must set appropriate CORS headers

## Migration from Traditional Auth

If migrating from username/password:

1. Keep existing user table structure
2. Add `wallet_address` column (already exists)
3. SIWE creates users automatically on first login
4. Users are identified by wallet address
5. Username is auto-generated but can be updated

## Next Steps

1. **Install dependencies**: `go mod tidy`
2. **Configure environment**: Copy `.env.example` to `.env` and update
3. **Start Redis**: `docker-compose up redis -d`
4. **Run migrations**: `make migrate-up`
5. **Start server**: `make run-dev`
6. **Test authentication**: Use the frontend or cURL

## Resources

- [SIWE Specification](https://eips.ethereum.org/EIPS/eip-4361)
- [wagmi Documentation](https://wagmi.sh/)
- [siwe-go Library](https://github.com/spruceid/siwe-go)
