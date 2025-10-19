-- Create orders table (Order Domain)
-- Stores customer orders after checkout
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cart_id UUID REFERENCES carts(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'paid', 'shipped', 'completed', 'cancelled')),
    total NUMERIC(12, 2) NOT NULL CHECK (total >= 0),
    payment_ref VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for orders table
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_user_status ON orders(user_id, status);

-- Create order_items table (Order Domain)
-- Stores items within an order (snapshot at time of purchase)
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for order_items table
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Enable Row-Level Security (RLS) on orders table
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view their own orders
CREATE POLICY orders_user_policy ON orders
    FOR SELECT
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- Policy: Users can only create orders for themselves
CREATE POLICY orders_user_insert_policy ON orders
    FOR INSERT
    WITH CHECK (user_id = current_setting('app.current_user_id', true)::UUID);

-- Policy: Users can update their own pending orders (for cancellation)
CREATE POLICY orders_user_update_policy ON orders
    FOR UPDATE
    USING (
        user_id = current_setting('app.current_user_id', true)::UUID 
        AND status = 'pending'
    )
    WITH CHECK (
        user_id = current_setting('app.current_user_id', true)::UUID
    );

-- Policy: Sellers can view orders containing their products
CREATE POLICY orders_seller_policy ON orders
    FOR SELECT
    USING (
        id IN (
            SELECT DISTINCT oi.order_id 
            FROM order_items oi
            JOIN products p ON oi.product_id = p.id
            WHERE p.seller_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Admins can view and modify all orders
CREATE POLICY orders_admin_policy ON orders
    FOR ALL
    USING (current_setting('app.current_user_role', true) = 'admin');

-- Enable Row-Level Security (RLS) on order_items table
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view order items for their own orders
CREATE POLICY order_items_user_policy ON order_items
    FOR SELECT
    USING (
        order_id IN (
            SELECT id FROM orders WHERE user_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Users can insert order items for their own orders
CREATE POLICY order_items_user_insert_policy ON order_items
    FOR INSERT
    WITH CHECK (
        order_id IN (
            SELECT id FROM orders WHERE user_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Sellers can view order items for their products
CREATE POLICY order_items_seller_policy ON order_items
    FOR SELECT
    USING (
        product_id IN (
            SELECT id FROM products WHERE seller_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Admins can view all order items
CREATE POLICY order_items_admin_policy ON order_items
    FOR SELECT
    USING (current_setting('app.current_user_role', true) = 'admin');
