-- Drop RLS policies for transactions
DROP POLICY IF EXISTS transactions_admin_policy ON transactions;
DROP POLICY IF EXISTS transactions_user_insert_policy ON transactions;
DROP POLICY IF EXISTS transactions_user_policy ON transactions;

-- Disable RLS on transactions
ALTER TABLE transactions DISABLE ROW LEVEL SECURITY;

-- Drop transactions table
DROP TABLE IF EXISTS transactions CASCADE;

-- Drop RLS policies for wallets
DROP POLICY IF EXISTS wallets_admin_policy ON wallets;
DROP POLICY IF EXISTS wallets_user_policy ON wallets;

-- Disable RLS on wallets
ALTER TABLE wallets DISABLE ROW LEVEL SECURITY;

-- Drop wallets table
DROP TABLE IF EXISTS wallets CASCADE;
