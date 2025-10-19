# Database Migrations

This directory contains database migration files for the CaribX backend.

## Migration Files

Migrations are organized by domain to align with the Domain-Driven Design (DDD) architecture.

### 000001_create_users_domain
Creates the `users` table with proper indexes and update triggers.

**Tables created:**
- users (with wallet_address, role, etc.)

**Indexes:**
- idx_users_wallet_address
- idx_users_role

### 000002_create_auth_domain
Creates authentication-related tables.

**Tables created:**
- refresh_tokens

**Indexes:**
- idx_refresh_tokens_user_id
- idx_refresh_tokens_token
- idx_refresh_tokens_expires_at

### 000003_create_categories_domain
Creates the `categories` table with default product categories.

**Tables created:**
- categories

**Default categories:** Electronics, Fashion, Home & Garden, Sports & Outdoors, Books & Media, Food & Beverages, Health & Beauty, Toys & Games

### 000004_create_products_domain
Creates the `products` table with Row-Level Security (RLS) policies.

**Tables created:**
- products

**Indexes:**
- idx_products_seller_id
- idx_products_category_id
- idx_products_is_active
- idx_products_created_at

**RLS Policies:**
- products_seller_policy: Sellers can only view/modify their own products
- products_public_view_policy: Anyone can view active products
- products_admin_policy: Admins have full access

### 000005_create_wallets_domain
Creates wallet and transaction tables with RLS policies for financial security.

**Tables created:**
- wallets
- transactions

**Indexes:**
- idx_wallets_user_id
- idx_transactions_wallet_id
- idx_transactions_status
- idx_transactions_created_at
- idx_transactions_type

**RLS Policies:**
- wallets_user_policy: Users can only access their own wallet
- wallets_admin_policy: Admins can view all wallets
- transactions_user_policy: Users can only view their wallet's transactions
- transactions_user_insert_policy: Users can only create transactions for their wallet
- transactions_admin_policy: Admins can view all transactions

### 000006_create_carts_domain
Creates shopping cart tables with RLS policies.

**Tables created:**
- carts
- cart_items

**Indexes:**
- idx_carts_user_id
- idx_carts_status
- idx_carts_user_status
- idx_cart_items_cart_id
- idx_cart_items_product_id

**RLS Policies:**
- carts_user_policy: Users can only access their own carts
- carts_admin_policy: Admins can view all carts
- cart_items_user_policy: Users can only access items in their carts
- cart_items_admin_policy: Admins can view all cart items

### 000007_create_orders_domain
Creates order tables with comprehensive RLS policies.

**Tables created:**
- orders
- order_items

**Indexes:**
- idx_orders_user_id
- idx_orders_status
- idx_orders_created_at
- idx_orders_user_status
- idx_order_items_order_id
- idx_order_items_product_id

**RLS Policies:**
- orders_user_policy: Users can view their own orders
- orders_user_insert_policy: Users can create orders for themselves
- orders_user_update_policy: Users can update their pending orders
- orders_seller_policy: Sellers can view orders containing their products
- orders_admin_policy: Admins have full access
- order_items_user_policy: Users can view items in their orders
- order_items_user_insert_policy: Users can insert items in their orders
- order_items_seller_policy: Sellers can view items for their products
- order_items_admin_policy: Admins can view all order items

### 000008_seed_marketplace_data
Seeds the database with sample marketplace data for development and testing.

**Seeded data includes:**
- 8 seller users with wallet addresses
- 8 wallets with initial balances (JAM currency)
- 40+ Caribbean-themed products across all categories:
  - **Electronics**: Samsung phones, JBL speakers, laptops, AirPods
  - **Fashion**: Rasta t-shirts, sundresses, hoodies, sandals
  - **Food & Beverages**: Blue Mountain Coffee, jerk seasoning, pepper sauce, coconut water, tropical fruits
  - **Home & Garden**: Bamboo wind chimes, tropical plants, rattan furniture, Caribbean art
  - **Health & Beauty**: Coconut oil, shea butter lotion, aloe vera gel, sunscreen
  - **Sports & Outdoors**: Snorkel sets, beach volleyball, surfboards, coolers
  - **Books & Media**: Bob Marley vinyl, Caribbean cookbooks, reggae CDs
  - **Toys & Games**: Beach balls, dominoes, sand castle sets, puzzles

All products feature:
- Realistic Caribbean/Jamaican themes
- Appropriate pricing in JAM (Jamaican Dollars)
- Stock quantities
- Multiple product images (placeholder URLs)
- Active status for immediate marketplace visibility

## Running Migrations

### Apply migrations (up)
```bash
make migrate-up
```

Or manually:
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/caribx?sslmode=disable" up
```

### Rollback migrations (down)
```bash
make migrate-down
```

Or manually:
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/caribx?sslmode=disable" down
```

### Check migration status
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/caribx?sslmode=disable" version
```

## Creating New Migrations

To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

This creates two files:
- `NNNNNN_migration_name.up.sql` - Forward migration
- `NNNNNN_migration_name.down.sql` - Rollback migration

## Testing with Seeded Data

After running the seed migration, you can:

1. **Browse products via API:**
   ```bash
   curl http://localhost:8080/v1/products
   ```

2. **Filter by category:**
   ```bash
   curl "http://localhost:8080/v1/products?category_id=<category-id>"
   ```

3. **Search products:**
   ```bash
   curl "http://localhost:8080/v1/products?search=coffee"
   ```

4. **Get specific product:**
   ```bash
   curl http://localhost:8080/v1/products/<product-id>
   ```

## Seller Accounts

The following seller accounts are available for testing:

| Username | Wallet Address | Specialization | Balance |
|----------|----------------|----------------|---------|
| island_treasures | 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb | Home goods, beauty products | 5,000 JAM |
| caribbean_spices | 0x8Ba1f109551bD432803012645Ac136ddd64DBA72 | Food & beverages | 7,500 JAM |
| tropical_fruits | 0x9F8c78B6B6C93A9C3d9d8f8E8e7e7e6e6e5e5e4e | Fresh produce | 3,200 JAM |
| jamaican_crafts | 0xaE1c94e42e57D9c2e1e1e1e1e1e1e1e1e1e1e1e1 | Handmade crafts, games | 4,800 JAM |
| reggae_vibes | 0xbF2d05f53f68E0d3f2f2f2f2f2f2f2f2f2f2f2f2 | Music & media | 6,100 JAM |
| beach_essentials | 0xcE3e16g64g79F0e4g3g3g3g3g3g3g3g3g3g3g3g3 | Sports & outdoor gear | 5,500 JAM |
| tech_caribbean | 0xdF4f27h75h80A1f5h4h4h4h4h4h4h4h4h4h4h4h4 | Electronics | 8,900 JAM |
| fashion_island | 0xeE5e38i86i91B2e6i5i5i5i5i5i5i5i5i5i5i5i5 | Fashion & apparel | 4,200 JAM |

## Notes

- All seeded data uses UUID format for IDs
- Prices are in Jamaican Dollars (JAM)
- Product images use placeholder URLs that should be replaced with actual image hosting
- The seed data is designed for development/testing and should not be used in production
- Running the down migration will remove all seeded data cleanly
