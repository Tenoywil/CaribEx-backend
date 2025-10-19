-- Seed marketplace data with Caribbean-themed products
-- This migration seeds the database with sample users, categories, and products

-- Create sample seller users
INSERT INTO users (id, username, wallet_address, role, created_at, updated_at) VALUES
('a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID, 'island_treasures', '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb', 'seller', NOW(), NOW()),
('b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID, 'caribbean_spices', '0x8Ba1f109551bD432803012645Ac136ddd64DBA72', 'seller', NOW(), NOW()),
('c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f'::UUID, 'tropical_fruits', '0x9F8c78B6B6C93A9C3d9d8f8E8e7e7e6e6e5e5e4e', 'seller', NOW(), NOW()),
('d4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID, 'jamaican_crafts', '0xaE1c94e42e57D9c2e1e1e1e1e1e1e1e1e1e1e1e1', 'seller', NOW(), NOW()),
('e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b'::UUID, 'reggae_vibes', '0xbF2d05f53f68E0d3f2f2f2f2f2f2f2f2f2f2f2f2', 'seller', NOW(), NOW()),
('f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID, 'beach_essentials', '0xcE3e16g64g79F0e4g3g3g3g3g3g3g3g3g3g3g3g3', 'seller', NOW(), NOW()),
('a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID, 'tech_caribbean', '0xdF4f27h75h80A1f5h4h4h4h4h4h4h4h4h4h4h4h4', 'seller', NOW(), NOW()),
('b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID, 'fashion_island', '0xeE5e38i86i91B2e6i5i5i5i5i5i5i5i5i5i5i5i5', 'seller', NOW(), NOW());

-- Create wallets for sellers
INSERT INTO wallets (id, user_id, balance, currency, updated_at) VALUES
(gen_random_uuid(), 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID, 5000.00, 'JAM', NOW()),
(gen_random_uuid(), 'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID, 7500.00, 'JAM', NOW()),
(gen_random_uuid(), 'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f'::UUID, 3200.00, 'JAM', NOW()),
(gen_random_uuid(), 'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID, 4800.00, 'JAM', NOW()),
(gen_random_uuid(), 'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b'::UUID, 6100.00, 'JAM', NOW()),
(gen_random_uuid(), 'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID, 5500.00, 'JAM', NOW()),
(gen_random_uuid(), 'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID, 8900.00, 'JAM', NOW()),
(gen_random_uuid(), 'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID, 4200.00, 'JAM', NOW());

-- Get category IDs (assuming they were inserted in the initial migration)
-- We'll store them in variables for reference

-- Insert products for Electronics category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Samsung Galaxy A54 Smartphone',
    'Latest Samsung Galaxy A54 with 128GB storage, 6GB RAM, and stunning 6.4" AMOLED display. Perfect for capturing island memories.',
    45000.00,
    15,
    ARRAY['https://images.example.com/samsung-a54-1.jpg', 'https://images.example.com/samsung-a54-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'JBL Bluetooth Speaker - Waterproof',
    'Portable JBL speaker perfect for beach parties. Waterproof and dust-proof with 12 hours battery life.',
    8500.00,
    25,
    ARRAY['https://images.example.com/jbl-speaker-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Apple AirPods Pro 2nd Gen',
    'Active noise cancellation, spatial audio, and sweat-resistant design. Experience premium sound quality.',
    32000.00,
    10,
    ARRAY['https://images.example.com/airpods-pro-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'HP Laptop 15.6" - Core i5',
    'HP Pavilion laptop with Intel Core i5, 8GB RAM, 256GB SSD. Perfect for work and study.',
    65000.00,
    8,
    ARRAY['https://images.example.com/hp-laptop-1.jpg', 'https://images.example.com/hp-laptop-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Fashion category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Rasta Colors T-Shirt - Unisex',
    'Premium cotton t-shirt featuring the iconic Rasta colors. Comfortable fit for everyday wear.',
    1200.00,
    50,
    ARRAY['https://images.example.com/rasta-tshirt-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Fashion' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Caribbean Floral Sundress',
    'Light and breezy sundress with tropical floral print. Perfect for warm island days.',
    3500.00,
    30,
    ARRAY['https://images.example.com/sundress-1.jpg', 'https://images.example.com/sundress-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Fashion' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Bob Marley Graphic Hoodie',
    'Comfortable hoodie featuring Bob Marley artwork. Made from soft cotton blend.',
    4200.00,
    20,
    ARRAY['https://images.example.com/marley-hoodie-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Fashion' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Beach Sandals - Leather',
    'Handcrafted leather sandals perfect for beach walks. Durable and stylish.',
    2800.00,
    40,
    ARRAY['https://images.example.com/sandals-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Fashion' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Food & Beverages category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Jamaican Blue Mountain Coffee - 1lb',
    'Premium Jamaican Blue Mountain Coffee beans. Smooth, rich flavor with subtle chocolate notes.',
    6500.00,
    100,
    ARRAY['https://images.example.com/blue-mountain-coffee-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Food & Beverages' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Jerk Seasoning Mix - Authentic',
    'Traditional Jamaican jerk seasoning blend. Add authentic Caribbean flavor to your meals.',
    800.00,
    200,
    ARRAY['https://images.example.com/jerk-seasoning-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Food & Beverages' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Scotch Bonnet Pepper Sauce - Hot',
    'Fiery hot sauce made with authentic Scotch Bonnet peppers. Handle with care!',
    650.00,
    150,
    ARRAY['https://images.example.com/pepper-sauce-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Food & Beverages' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e'::UUID,
    'Coconut Water - Fresh 12-Pack',
    'Fresh coconut water straight from Jamaican coconuts. Natural electrolytes and refreshing taste.',
    1800.00,
    80,
    ARRAY['https://images.example.com/coconut-water-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Food & Beverages' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f'::UUID,
    'Tropical Fruit Box - Mixed',
    'Fresh seasonal tropical fruits including mangoes, pineapples, and guavas. Farm-fresh delivery.',
    2500.00,
    60,
    ARRAY['https://images.example.com/fruit-box-1.jpg', 'https://images.example.com/fruit-box-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Food & Beverages' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Home & Garden category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Bamboo Wind Chimes',
    'Handcrafted bamboo wind chimes producing soothing Caribbean melodies. Perfect for patios.',
    1500.00,
    35,
    ARRAY['https://images.example.com/wind-chimes-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Home & Garden' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Tropical Plant - Bird of Paradise',
    'Beautiful Bird of Paradise plant in decorative pot. Adds tropical flair to any space.',
    3200.00,
    25,
    ARRAY['https://images.example.com/bird-paradise-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Home & Garden' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Rattan Furniture Set - Outdoor',
    'Complete outdoor furniture set with rattan chairs and table. Weather-resistant and stylish.',
    28000.00,
    5,
    ARRAY['https://images.example.com/rattan-set-1.jpg', 'https://images.example.com/rattan-set-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Home & Garden' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Caribbean Art Print - Sunset',
    'Beautiful canvas print of Caribbean sunset. Ready to hang, 24x36 inches.',
    4500.00,
    15,
    ARRAY['https://images.example.com/sunset-print-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Home & Garden' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Health & Beauty category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Coconut Oil - Extra Virgin 16oz',
    'Pure extra virgin coconut oil. Perfect for cooking, skin, and hair care.',
    1200.00,
    100,
    ARRAY['https://images.example.com/coconut-oil-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Health & Beauty' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Shea Butter Body Lotion',
    'Moisturizing body lotion with shea butter and tropical scents. Nourishes dry skin.',
    1800.00,
    75,
    ARRAY['https://images.example.com/shea-lotion-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Health & Beauty' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'Aloe Vera Gel - Pure',
    'Pure aloe vera gel for soothing sunburns and skin irritation. Natural healing properties.',
    950.00,
    120,
    ARRAY['https://images.example.com/aloe-gel-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Health & Beauty' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::UUID,
    'SPF 50 Sunscreen - Reef Safe',
    'Reef-safe sunscreen protecting against UVA/UVB rays. Perfect for Caribbean beaches.',
    1650.00,
    90,
    ARRAY['https://images.example.com/sunscreen-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Health & Beauty' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Sports & Outdoors category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID,
    'Snorkel Set - Adult',
    'Complete snorkeling set with mask, snorkel, and fins. Explore Caribbean underwater life.',
    3800.00,
    40,
    ARRAY['https://images.example.com/snorkel-set-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Sports & Outdoors' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID,
    'Beach Volleyball Set',
    'Professional beach volleyball with net and ball. Perfect for beach games.',
    5500.00,
    20,
    ARRAY['https://images.example.com/volleyball-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Sports & Outdoors' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID,
    'Surfboard - Beginner 7ft',
    'Foam surfboard perfect for beginners. Stable and easy to learn on.',
    18000.00,
    12,
    ARRAY['https://images.example.com/surfboard-1.jpg', 'https://images.example.com/surfboard-2.jpg'],
    (SELECT id FROM categories WHERE name = 'Sports & Outdoors' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c'::UUID,
    'Cooler Box - 48 Quart',
    'Large cooler perfect for beach trips and picnics. Keeps ice for up to 3 days.',
    4200.00,
    30,
    ARRAY['https://images.example.com/cooler-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Sports & Outdoors' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Books & Media category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Bob Marley - Legend Vinyl',
    'Classic Bob Marley Legend album on vinyl. Remastered edition with iconic reggae hits.',
    3500.00,
    25,
    ARRAY['https://images.example.com/marley-vinyl-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Books & Media' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b'::UUID,
    'Caribbean Cookbook - Traditional Recipes',
    'Comprehensive cookbook featuring traditional Caribbean recipes from across the islands.',
    2200.00,
    50,
    ARRAY['https://images.example.com/cookbook-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Books & Media' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Jamaica History & Culture Book',
    'Comprehensive guide to Jamaican history, culture, and traditions. Educational and engaging.',
    1800.00,
    40,
    ARRAY['https://images.example.com/history-book-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Books & Media' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b'::UUID,
    'Reggae Music Collection - 3 CD Set',
    'Ultimate reggae music collection featuring top artists. 50 classic tracks.',
    2800.00,
    35,
    ARRAY['https://images.example.com/reggae-cds-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Books & Media' LIMIT 1),
    true,
    NOW(),
    NOW();

-- Insert products for Toys & Games category
INSERT INTO products (seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Beach Ball Set - Colorful',
    'Set of 3 colorful beach balls in various sizes. Perfect for family beach fun.',
    900.00,
    60,
    ARRAY['https://images.example.com/beach-balls-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Toys & Games' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Dominoes - Caribbean Edition',
    'Traditional domino set with Caribbean-themed designs. Family game night favorite.',
    1200.00,
    45,
    ARRAY['https://images.example.com/dominoes-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Toys & Games' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Sand Castle Building Set',
    'Complete sand castle building kit with molds and tools. Hours of beach entertainment.',
    1500.00,
    50,
    ARRAY['https://images.example.com/sand-castle-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Toys & Games' LIMIT 1),
    true,
    NOW(),
    NOW()
UNION ALL SELECT 
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a'::UUID,
    'Caribbean Jigsaw Puzzle - 1000pc',
    '1000-piece jigsaw puzzle featuring beautiful Caribbean beach scene.',
    1800.00,
    30,
    ARRAY['https://images.example.com/puzzle-1.jpg'],
    (SELECT id FROM categories WHERE name = 'Toys & Games' LIMIT 1),
    true,
    NOW(),
    NOW();
