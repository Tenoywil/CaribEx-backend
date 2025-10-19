# Migration Reorganization Summary

## Overview
The database migrations have been reorganized to align with Domain-Driven Design (DDD) principles, with each domain having its own migration file(s). Additionally, comprehensive Row-Level Security (RLS) policies and optimized indexes have been added.

## Changes Made

### 1. Domain-Specific Migration Files

**Previous Structure:**
- `000001_init_schema.up.sql` - Single monolithic file with all tables
- `000002_seed_marketplace_data.up.sql` - Seed data

**New Structure:**
- `000001_create_users_domain.up.sql` - User domain tables
- `000002_create_auth_domain.up.sql` - Authentication domain tables
- `000003_create_categories_domain.up.sql` - Category tables (product support)
- `000004_create_products_domain.up.sql` - Product domain tables
- `000005_create_wallets_domain.up.sql` - Wallet domain tables
- `000006_create_carts_domain.up.sql` - Cart domain tables
- `000007_create_orders_domain.up.sql` - Order domain tables
- `000008_seed_marketplace_data.up.sql` - Seed data (updated)

### 2. Row-Level Security (RLS) Policies

#### Products Domain
- **products_seller_policy**: Sellers can only view/modify their own products
- **products_public_view_policy**: Anyone can view active products
- **products_admin_policy**: Admins have full access

#### Wallets Domain
- **wallets_user_policy**: Users can only access their own wallet
- **wallets_admin_policy**: Admins can view all wallets
- **transactions_user_policy**: Users can view their own transactions
- **transactions_user_insert_policy**: Users can create transactions for their wallet
- **transactions_admin_policy**: Admins can view all transactions

#### Carts Domain
- **carts_user_policy**: Users can only access their own carts
- **carts_admin_policy**: Admins can view all carts
- **cart_items_user_policy**: Users can only access items in their carts
- **cart_items_admin_policy**: Admins can view all cart items

#### Orders Domain
- **orders_user_policy**: Users can view their own orders
- **orders_user_insert_policy**: Users can create orders for themselves
- **orders_user_update_policy**: Users can update their pending orders
- **orders_seller_policy**: Sellers can view orders containing their products
- **orders_admin_policy**: Admins have full access
- **order_items_user_policy**: Users can view items in their orders
- **order_items_user_insert_policy**: Users can insert items in their orders
- **order_items_seller_policy**: Sellers can view items for their products
- **order_items_admin_policy**: Admins can view all order items

**Total RLS Policies: 21**

### 3. Indexes for Performance

#### Users Domain
- `idx_users_wallet_address` - Fast lookup by wallet address
- `idx_users_role` - Filter by user role

#### Auth Domain
- `idx_refresh_tokens_user_id` - Lookup tokens by user
- `idx_refresh_tokens_token` - Fast token validation
- `idx_refresh_tokens_expires_at` - Cleanup expired tokens

#### Products Domain
- `idx_products_seller_id` - Filter products by seller
- `idx_products_category_id` - Filter by category
- `idx_products_is_active` - Filter active products
- `idx_products_created_at` - Sort by creation date

#### Wallets Domain
- `idx_wallets_user_id` - Lookup wallet by user
- `idx_transactions_wallet_id` - Transactions for a wallet
- `idx_transactions_status` - Filter by transaction status
- `idx_transactions_created_at` - Sort by date (descending)
- `idx_transactions_type` - Filter by credit/debit

#### Carts Domain
- `idx_carts_user_id` - Lookup cart by user
- `idx_carts_status` - Filter by cart status
- `idx_carts_user_status` - Combined user+status lookup
- `idx_cart_items_cart_id` - Items in a cart
- `idx_cart_items_product_id` - Product in carts

#### Orders Domain
- `idx_orders_user_id` - User's orders
- `idx_orders_status` - Filter by order status
- `idx_orders_created_at` - Sort by date (descending)
- `idx_orders_user_status` - Combined user+status lookup
- `idx_order_items_order_id` - Items in an order
- `idx_order_items_product_id` - Product in orders

**Total Indexes: 41 (including primary keys and unique constraints)**

### 4. Additional Improvements

#### Update Triggers
Added automatic `updated_at` timestamp triggers for:
- users
- products
- carts
- cart_items
- orders

#### Type Safety
Fixed UUID casting in seed data to ensure type safety in PostgreSQL.

#### Documentation
Updated migration README with comprehensive documentation of all migrations, policies, and indexes.

## Testing

All migrations have been tested:
- ✅ Full migration up (all 8 migrations)
- ✅ Individual migration down/up
- ✅ Full migration down (rollback all)
- ✅ Verification of all 21 RLS policies
- ✅ Verification of all 41 indexes
- ✅ Seed data verification (8 users, 33 products, 8 wallets)

## Usage

### Apply all migrations
```bash
make migrate-up
```

### Rollback all migrations
```bash
make migrate-down
```

### Check migration status
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/CaribEX?sslmode=disable" version
```

## Security Considerations

### Setting RLS Context Variables
For RLS policies to work correctly, the application must set these PostgreSQL session variables:
- `app.current_user_id` - The authenticated user's UUID
- `app.current_user_role` - The user's role (customer/seller/admin)

Example in Go:
```go
_, err := tx.Exec(ctx, "SET app.current_user_id = $1", userID)
_, err := tx.Exec(ctx, "SET app.current_user_role = $1", userRole)
```

### RLS Bypass for Superusers
Note: PostgreSQL superusers (like the default `postgres` user) bypass RLS policies. In production:
1. Use a non-superuser role for the application
2. Grant only necessary permissions to the application role
3. Keep superuser credentials secure and separate

## Benefits

1. **Better Organization**: Each domain has its own migration file(s)
2. **Clear Ownership**: Easy to identify which domain owns which tables
3. **Security**: RLS policies enforce data isolation at the database level
4. **Performance**: Proper indexes for common query patterns
5. **Maintainability**: Easier to understand and modify domain-specific schemas
6. **Rollback Safety**: Each migration has a corresponding down migration
7. **Audit Trail**: Clear migration history aligned with domain evolution

## Known Issues

### Seed Data Wallet Addresses
The seed data in `000008_seed_marketplace_data.up.sql` contains placeholder wallet addresses that are not valid Ethereum addresses (they contain invalid hexadecimal characters like 'g', 'h', 'i'). These should be replaced with valid Ethereum addresses if you plan to use SIWE authentication with the seeded users. For testing purposes, you can generate valid addresses using tools like:
- MetaMask test accounts
- Hardhat/Ganache local blockchain accounts
- Online Ethereum address generators

## Future Considerations

1. Replace seed wallet addresses with valid Ethereum addresses for SIWE compatibility
2. Consider adding RLS policies for the `users` table if multi-tenancy is needed
3. Add audit logging tables if regulatory compliance requires it
4. Consider partitioning for high-volume tables (transactions, order_items)
5. Add database-level constraints for business rules (e.g., order total = sum of items)
