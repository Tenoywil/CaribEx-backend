-- Create products table (Product Domain)
-- Stores marketplace product listings with seller association
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    images TEXT[],
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for products table
CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_created_at ON products(created_at DESC);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Enable Row-Level Security (RLS) on products table
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

-- Policy: Sellers can only view and modify their own products
CREATE POLICY products_seller_policy ON products
    FOR ALL
    USING (seller_id = current_setting('app.current_user_id', true)::UUID);

-- Policy: Anyone can view active products (for browsing)
CREATE POLICY products_public_view_policy ON products
    FOR SELECT
    USING (is_active = true);

-- Policy: Admins can view and modify all products
CREATE POLICY products_admin_policy ON products
    FOR ALL
    USING (current_setting('app.current_user_role', true) = 'admin');
