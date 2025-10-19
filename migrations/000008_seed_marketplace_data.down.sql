-- Remove seeded marketplace data
-- This migration removes all seeded products, users, and wallets

-- Delete products (cascade will handle related data)
DELETE FROM products WHERE seller_id IN (
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e',
    'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f',
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a',
    'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b',
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c',
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d',
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'
);

-- Delete wallets for seeded users
DELETE FROM wallets WHERE user_id IN (
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e',
    'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f',
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a',
    'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b',
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c',
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d',
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'
);

-- Delete seeded users
DELETE FROM users WHERE id IN (
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
    'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e',
    'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f',
    'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a',
    'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b',
    'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c',
    'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d',
    'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e'
);
