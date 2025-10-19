-- Drop RLS policies for cart_items
DROP POLICY IF EXISTS cart_items_admin_policy ON cart_items;
DROP POLICY IF EXISTS cart_items_user_policy ON cart_items;

-- Disable RLS on cart_items
ALTER TABLE cart_items DISABLE ROW LEVEL SECURITY;

-- Drop trigger
DROP TRIGGER IF EXISTS update_cart_items_updated_at ON cart_items;

-- Drop cart_items table
DROP TABLE IF EXISTS cart_items CASCADE;

-- Drop RLS policies for carts
DROP POLICY IF EXISTS carts_admin_policy ON carts;
DROP POLICY IF EXISTS carts_user_policy ON carts;

-- Disable RLS on carts
ALTER TABLE carts DISABLE ROW LEVEL SECURITY;

-- Drop trigger
DROP TRIGGER IF EXISTS update_carts_updated_at ON carts;

-- Drop carts table
DROP TABLE IF EXISTS carts CASCADE;
