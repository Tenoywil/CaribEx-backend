-- Create wallets table (Wallet Domain)
-- Stores user wallet balances for different currencies
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(20, 8) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'JAM',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for wallets table
CREATE INDEX idx_wallets_user_id ON wallets(user_id);

-- Create transactions table (Wallet Domain - Ledger)
-- Append-only ledger for all wallet transactions
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('credit', 'debit')),
    amount NUMERIC(20, 8) NOT NULL CHECK (amount >= 0),
    reference VARCHAR(255),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'success', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for transactions table
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_type ON transactions(type);

-- Enable Row-Level Security (RLS) on wallets table
ALTER TABLE wallets ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view and modify their own wallet
CREATE POLICY wallets_user_policy ON wallets
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- Policy: Admins can view all wallets
CREATE POLICY wallets_admin_policy ON wallets
    FOR SELECT
    USING (current_setting('app.current_user_role', true) = 'admin');

-- Enable Row-Level Security (RLS) on transactions table
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only view transactions for their own wallet
CREATE POLICY transactions_user_policy ON transactions
    FOR SELECT
    USING (
        wallet_id IN (
            SELECT id FROM wallets WHERE user_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Users can only insert transactions for their own wallet
CREATE POLICY transactions_user_insert_policy ON transactions
    FOR INSERT
    WITH CHECK (
        wallet_id IN (
            SELECT id FROM wallets WHERE user_id = current_setting('app.current_user_id', true)::UUID
        )
    );

-- Policy: Admins can view all transactions
CREATE POLICY transactions_admin_policy ON transactions
    FOR SELECT
    USING (current_setting('app.current_user_role', true) = 'admin');
