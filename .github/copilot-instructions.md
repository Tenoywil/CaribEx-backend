# CaribX — Backend Blueprint (Go 1.25.x)

**Project:** CaribX — Blockchain Money Transfer & Marketplace for Jamaica & the Caribbean
**Purpose:** Rapid, secure, hackathon-ready backend enabling P2P transfers, marketplace listings, wallet management, and order processing — built with DDD, strong security, and two-level caching.

---

## 1) Executive Summary & Conversation Context

You (Tenoy) are soloing a hackathon to deliver a blockchain-backed money transfer + marketplace MVP. Requirements collected across the conversation include:

* A Go backend following **DDD** (domain, repository, usecase, controller, routes).
* PostgreSQL (pgx + pgxpool) with **Row-Level Security (RLS)** and appropriate indexes.
* Two-level cache-aside: **L1 in-memory** (single-instance fast cache) + **L2 Redis** (shared cache).
* Auth via **SIWE (Sign-In With Ethereum)** and server sessions (HTTP-only cookies) plus optional JWT for m2m.
* Strict route validation, sanitization, rate limiting, CORS, circuit breaker for external providers, structured logging, monitoring (Prometheus + OTel), and no PII leakage.
* Frontend will be Next.js 15 + Redux Toolkit + Redux-Saga; interactions include send/receive, browse/list products, manage carts/orders, wallet management.

This document is a developer-oriented blueprint so Copilot or a human dev can scaffold and start coding immediately.

---

## 2) Goals & Non-Functional Requirements

* **Security-first**: redact PII in logs, RLS for tenant isolation, validated envs.
* **Resilience**: circuit breakers, rate limiting, metrics/alerts.
* **Performance**: L1 + L2 caching, singleflight/coalescing to avoid stampedes.
* **Extensibility**: DDD structure, interfaces for payment/wallet providers.
* **Hackathon deliverable**: working flows for wallet auth (SIWE), send/receive, list/buy products, manage carts and orders.

---

## 3) Project Structure (scaffold)

```
/cmd/api-server/main.go
/internal/domain/
  user/
  product/
  cart/
  order/
  wallet/
/internal/repository/
  postgres/
/internal/usecase/
/internal/controller/
/internal/routes/
/pkg/middleware/
pkg/config/
pkg/cache/
pkg/logger/
pkg/monitoring/
/migrations/
/tests/
Dockerfile
docker-compose.yml
Makefile

```

**Notes:** domain packages contain only business logic and types; infra adapters live in `internal/repository/postgres` and `pkg/cache`.

---

## 4) Database Schema (tables + key columns)

Use UUID PKs, timestamps, and tenant-aware columns where needed.

**users**

* id uuid PK
* username varchar
* wallet_address varchar (indexed)
* role enum: (customer,seller,admin)
* created_at, updated_at

**products**

* id uuid PK
* seller_id fk users(id)
* title, description, price (numeric), quantity
* images text[]
* category_id fk categories(id)
* is_active boolean
* created_at, updated_at

**categories**

* id uuid PK
* name varchar

**carts**

* id uuid PK
* user_id fk users(id)
* status enum (active, checked_out)
* total numeric

**cart_items**

* id uuid PK
* cart_id fk carts(id)
* product_id fk products(id)
* qty int
* price numeric

**orders**

* id uuid PK
* user_id fk users(id)
* cart_id fk carts(id)
* status enum (pending, paid, shipped, completed, cancelled)
* total numeric
* payment_ref varchar
* created_at, updated_at

**order_items**

* id uuid PK
* order_id fk orders(id)
* product_id fk products(id)
* qty int
* price numeric

**wallets**

* id uuid PK
* user_id fk users(id)
* balance numeric (store as integer of smallest unit if preferred)
* currency varchar (e.g., JAM, USD, USDC)
* updated_at

**transactions** (ledger)

* id uuid PK
* wallet_id fk wallets(id)
* type enum (credit,debit)
* amount numeric
* reference varchar
* status enum (pending,success,failed)
* created_at

**refresh_tokens**

* id uuid PK
* user_id fk users(id)
* token varchar(hash)
* expires_at

**indexes & policies**

* Index: products(seller_id), products(category_id), users(wallet_address)
* RLS: enable on multi-tenant tables (set app.current_tenant via session)

---

## 5) Key Domain Models & Aggregates (brief)

* **User**: identity (wallet), role, KYC flag (optional)
* **Wallet**: aggregate root that handles balance changes via transactions (atomic DB tx)
* **Product**: listing aggregate with price, inventory
* **Cart**: transient aggregate wiring user & cart_items
* **Order**: record of checkout with order_items and payment state

---

## 6) API Surface (selected endpoints)

**Auth & Profile**

* `GET /v1/auth/nonce` → issue SIWE nonce
* `POST /v1/auth/siwe` → verify signature, create/update user, set session cookie
* `GET /v1/auth/me` → current authenticated user

**Wallet**

* `GET /v1/wallet` → wallet summary
* `POST /v1/wallet/send` → submit outgoing transfer (create pending tx -> on-chain signing or server-assisted)
* `POST /v1/wallet/receive` → (optional) register incoming off-chain transfer
* `GET /v1/wallet/transactions` → ledger

**Products & Marketplace**

* `GET /v1/products` → browse (supports paging, filters)
* `GET /v1/products/:id`
* `POST /v1/products` (seller)
* `PUT /v1/products/:id`
* `DELETE /v1/products/:id`

**Cart & Orders**

* `GET /v1/cart`
* `POST /v1/cart/items`
* `PUT /v1/cart/items/:id`
* `DELETE /v1/cart/items/:id`
* `POST /v1/orders` → checkout
* `GET /v1/orders` → user orders

**Admin/Internal**

* `POST /internal/cache/invalidate` → invalidate keys (protected)

---

## 7) Middleware & Security Details

* **Validation**: `go-playground/validator` on DTOs + sanitization (trim, normalize)
* **Rate limiting**: per-IP + per-wallet token bucket; in-process limiter + Redis-backed for cluster
* **CORS**: whitelist origins; `Access-Control-Allow-Credentials` when using cookies
* **Session**: HTTP-only Secure SameSite cookie with short TTL; refresh via `refresh_tokens`
* **PII Handling**: never return email/phone in public endpoints; use hashed IDs in logs
* **Audit**: write audit events for money-moving operations to append-only `transactions`
* **Circuit Breaker**: wrap external RPCs (price feeds, payment provider) with `sony/gobreaker`

---

## 8) Cache Implementation Details

**L1 (In-memory)**

* Use `ristretto` or custom map + `singleflight` for request coalescing.
* Store small, hot objects: product list pages, recent wallet balances (but short TTL)

**L2 (Redis)**

* Use `go-redis/redis/v8` with JSON marshalling.
* Keys: `products:page:<p>`, `product:<id>`, `wallet:<user_id>`, `listings:recommended:<user_id>`
* TTL tuning: products 60–300s; wallets balances 5–15s; listings recommendations 30s–2m

**Write path**

* On mutation: update DB in transaction, then delete keys in L2 and L1.
* Optionally asynchronously repopulate L2 for hot keys.

**Cache stampede protection**

* Use `singleflight` to collapse concurrent L2->DB loads.

---

## 9) Observability & Monitoring

* **Logging:** `zap`/`zerolog` — JSON structured logs, include `request_id`, `user_hash`.
* **Metrics:** Prometheus counters (requests, errors), histograms (latency), gauges (cache sizes), circuit breaker states.
* **Tracing:** OpenTelemetry: instrument HTTP handlers, DB calls, cache calls.
* **Healthchecks:** `/healthz` and `/readyz` including DB & Redis connectivity

---

## 10) Testing Strategy (minimal but critical)

* **Unit tests**: domain logic (wallet debit/credit), validation.
* **Integration tests**: repositories + pgx against ephemeral Docker Postgres; use testcontainers or docker-compose.
* **Controller tests**: `httptest` for HTTP behavior.
* **Smoke E2E**: SIWE login -> create product -> add to cart -> checkout -> wallet transfer (use testnet / mocked chain or simulate ledger updates).

Test cases (minimal list):

* CreateWallet_DepositAndWithdraw_ShouldBalanceCorrectly
* CreateProduct_SellerOwnsProduct
* Checkout_CreatesOrderAndDecrementsInventory
* SendFunds_InsufficientBalance_ShouldFail

---

## 11) CI/CD & Dev Experience

* **Makefile**: `make build`, `make test`, `make lint`, `make run-dev`
* **Docker Compose**: Postgres, Redis, local dev server.
* **GitHub Actions**: lint, unit tests, build image, push to registry on main.
* **Secrets**: use GitHub Secrets or Vault for production keys.

---

## 12) Scaffolding Commands & Prompts (copy/paste)

```bash
# init module
mkdir caribx-backend && cd caribx-backend
go mod init github.com/<you>/caribx-backend
# install libs
go get github.com/jackc/pgx/v5 github.com/spf13/viper github.com/go-redis/redis/v8 github.com/sony/gobreaker github.com/go-playground/validator/v10 github.com/rs/zerolog github.com/prometheus/client_golang golang.org/x/sync/singleflight github.com/stretchr/testify
# create folders
mkdir -p cmd/internal/{domain,repository,usecase,controller,routes} pkg/{cache,config,middleware,logger,monitoring}
```

**AI/Copilot prompt (backend scaffolding):**

> "Scaffold a Go 1.25 DDD backend for CaribX with pgxpool, Redis cache, SIWE auth endpoints, domain packages for user/product/cart/order/wallet, middleware chain (request-id, validator, rate-limit, auth), and Prometheus metrics. Create sample DTOs, a wallet debit/credit usecase, and unit tests for wallet logic."

---

## 13) Next Steps & Prioritized TODO (for 1-day hackathon)

1. Scaffold project + Docker Compose (Postgres, Redis) — 30 min
2. Implement SIWE auth flow + session cookie — 60–90 min
3. Implement Wallet aggregate + transactions ledger + `POST /v1/wallet/send` — 60 min
4. Implement Product listing read endpoints + L1/L2 cache for product lists — 45 min
5. Implement Cart and Checkout (writes order, decrements inventory) — 60 min
6. Wire frontend (Next.js) with wallet auth and browse/purchase flow — 90 min
7. Tests & deploy to dev host or present local demo — remaining time

---

## 14) Appendix: Example SQL snippets

**Enable RLS on products**

```sql
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
CREATE POLICY products_tenant_policy ON products USING (tenant_id = current_setting('app.current_tenant')::uuid);
```

**Create wallets table (example)**

```sql
CREATE TABLE wallets (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references users(id),
  balance numeric not null default 0,
  currency text not null default 'JAM',
  updated_at timestamptz default now()
);
```

---

End of backend blueprint. This canvas is intended to be copied into the repo `README.md` and used as a primary onboarding document for any developer or code-assistant.
