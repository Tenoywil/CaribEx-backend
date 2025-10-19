-- Drop RLS policies
DROP POLICY IF EXISTS products_admin_policy ON products;
DROP POLICY IF EXISTS products_public_view_policy ON products;
DROP POLICY IF EXISTS products_seller_policy ON products;

-- Disable RLS
ALTER TABLE products DISABLE ROW LEVEL SECURITY;

-- Drop trigger
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

-- Drop table
DROP TABLE IF EXISTS products CASCADE;
