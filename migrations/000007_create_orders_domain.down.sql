-- Drop RLS policies for order_items
DROP POLICY IF EXISTS order_items_admin_policy ON order_items;
DROP POLICY IF EXISTS order_items_seller_policy ON order_items;
DROP POLICY IF EXISTS order_items_user_insert_policy ON order_items;
DROP POLICY IF EXISTS order_items_user_policy ON order_items;

-- Disable RLS on order_items
ALTER TABLE order_items DISABLE ROW LEVEL SECURITY;

-- Drop order_items table
DROP TABLE IF EXISTS order_items CASCADE;

-- Drop RLS policies for orders
DROP POLICY IF EXISTS orders_admin_policy ON orders;
DROP POLICY IF EXISTS orders_seller_policy ON orders;
DROP POLICY IF EXISTS orders_user_update_policy ON orders;
DROP POLICY IF EXISTS orders_user_insert_policy ON orders;
DROP POLICY IF EXISTS orders_user_policy ON orders;

-- Disable RLS on orders
ALTER TABLE orders DISABLE ROW LEVEL SECURITY;

-- Drop trigger
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop orders table
DROP TABLE IF EXISTS orders CASCADE;
