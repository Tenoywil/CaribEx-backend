-- Create carts table (Cart Domain)
-- Stores shopping carts for users
CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'checked_out')),
    total NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for carts table
CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_status ON carts(status);
CREATE INDEX idx_carts_user_status ON carts(user_id, status);

-- Create cart_items table (Cart Domain)
-- Stores items within a shopping cart
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cart_id, product_id)
);

-- Create indexes for cart_items table
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- Add triggers to update updated_at timestamp
CREATE TRIGGER update_carts_updated_at
BEFORE UPDATE ON carts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cart_items_updated_at
BEFORE UPDATE ON cart_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Enable Row-Level Security (RLS) on carts table
ALTER TABLE carts ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view and modify their own carts
CREATE POLICY carts_user_policy ON carts
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- Policy: Admins can view all carts
CREATE POLICY carts_admin_policy ON carts
    FOR SELECT
    USING (current_setting('app.current_user_role', true) = 'admin');

-- Enable Row-Level Security (RLS) on cart_items table
ALTER TABLE cart_items ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view and modify items in their own carts
CREATE POLICY cart_items_user_policy ON cart_items
    FOR ALL
    USING (
        cart_id IN (
            SELECT id FROM carts WHERE user_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Admins can view all cart items
CREATE POLICY cart_items_admin_policy ON cart_items
    FOR SELECT
    USING (current_setting('app.current_user_role', true) = 'admin');
