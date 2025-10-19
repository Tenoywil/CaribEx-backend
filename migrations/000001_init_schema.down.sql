-- Active: 1760838300877@@aws-1-us-east-2.pooler.supabase.com@5432@postgres
-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
