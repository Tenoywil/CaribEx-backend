-- Create categories table (Product Domain - Supporting Table)
-- Categories are used to organize products
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default categories
INSERT INTO categories (name) VALUES 
    ('Electronics'),
    ('Fashion'),
    ('Home & Garden'),
    ('Sports & Outdoors'),
    ('Books & Media'),
    ('Food & Beverages'),
    ('Health & Beauty'),
    ('Toys & Games')
ON CONFLICT (name) DO NOTHING;
