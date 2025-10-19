# Blockchain Transaction Verification API

## Overview

The CaribEX backend provides on-chain transaction verification endpoints that allow the frontend to verify Ethereum transactions initiated via wagmi. This ensures authenticity and accurate transaction state before processing payments or transfers.

## Backend Implementation

### Configuration

Add the following to your `.env` file:

```bash
# Blockchain Configuration
RPC_URL=https://mainnet.infura.io/v3/YOUR_INFURA_KEY
```

For testing, you can use:
- Sepolia testnet: `https://sepolia.infura.io/v3/YOUR_INFURA_KEY`
- Polygon: `https://polygon-mainnet.infura.io/v3/YOUR_INFURA_KEY`

### Supported Networks

The backend validates the following chain IDs:
- `1` - Ethereum Mainnet
- `11155111` - Sepolia Testnet
- `137` - Polygon Mainnet
- `80001` - Mumbai Testnet (Polygon)

**Note:** Goerli testnet has been deprecated and is no longer supported.

## API Endpoints

### 1. Verify Transaction

Verifies an on-chain transaction and logs it to the database.

**Endpoint:** `POST /v1/wallet/verify-transaction`

**Authentication:** Required (JWT/Session)

**Request Body:**
```json
{
  "txHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
  "chainId": 1
}
```

**Success Response (200):**
```json
{
  "status": "verified",
  "txHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
  "message": "Transaction successfully verified",
  "from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
  "to": "0x1234567890123456789012345678901234567890",
  "chainId": 1
}
```

**Error Responses:**

- **400 Bad Request** - Invalid request data
```json
{
  "error": "Invalid request: txHash is required"
}
```

- **400 Bad Request** - Unsupported chain
```json
{
  "status": "failed",
  "txHash": "0xabc...",
  "error": "unsupported chain ID",
  "message": "Transaction verification failed"
}
```

- **400 Bad Request** - Transaction pending
```json
{
  "status": "failed",
  "txHash": "0xabc...",
  "error": "transaction is still pending",
  "message": "Transaction verification failed"
}
```

- **400 Bad Request** - Transaction failed on-chain
```json
{
  "status": "failed",
  "txHash": "0xabc...",
  "error": "transaction failed on-chain",
  "message": "Transaction verification failed"
}
```

- **401 Unauthorized** - User not authenticated
```json
{
  "error": "User not authenticated"
}
```

### 2. Get Transaction Status

Retrieves the status of a transaction without logging it.

**Endpoint:** `GET /v1/wallet/transaction-status`

**Authentication:** Required (JWT/Session)

**Query Parameters:**
- `txHash` (required): The transaction hash to check
- `chainId` (optional): The chain ID (defaults to 1 - Ethereum mainnet)

**Example Request:**
```
GET /v1/wallet/transaction-status?txHash=0xabc...&chainId=1
```

**Success Response (200):**
```json
{
  "status": "success",
  "txHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
  "message": "Transaction status retrieved",
  "from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
  "to": "0x1234567890123456789012345678901234567890",
  "value": "1000000000000000000",
  "chainId": 1,
  "isPending": false
}
```

**Pending Transaction Response (200):**
```json
{
  "status": "success",
  "txHash": "0xabc...",
  "message": "Transaction is pending",
  "from": "",
  "to": "",
  "value": "0",
  "chainId": 1,
  "isPending": true
}
```

## Frontend Integration

### Prerequisites

```bash
npm install wagmi viem @wagmi/core
```

### 1. Send Transaction and Verify (React + wagmi)

```typescript
import { useWalletClient, useWaitForTransaction } from 'wagmi';
import { useState } from 'react';

const SendAndVerifyTransaction = () => {
  const { data: walletClient } = useWalletClient();
  const [txHash, setTxHash] = useState<string | null>(null);
  const [isVerifying, setIsVerifying] = useState(false);

  // Wait for transaction to be mined
  const { isLoading, isSuccess } = useWaitForTransaction({
    hash: txHash as `0x${string}`,
    enabled: !!txHash,
  });

  const sendTransaction = async () => {
    try {
      if (!walletClient) {
        throw new Error('Wallet not connected');
      }

      // Send the transaction
      const hash = await walletClient.sendTransaction({
        to: '0x1234567890123456789012345678901234567890',
        value: parseEther('0.01'),
      });

      setTxHash(hash);
    } catch (error) {
      console.error('Transaction failed:', error);
    }
  };

  // Verify transaction after it's mined
  useEffect(() => {
    const verifyTransaction = async () => {
      if (isSuccess && txHash && !isVerifying) {
        setIsVerifying(true);
        
        try {
          const response = await fetch('/api/v1/wallet/verify-transaction', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              // Include your auth token
              'Authorization': `Bearer ${authToken}`,
            },
            credentials: 'include', // For cookie-based auth
            body: JSON.stringify({
              txHash: txHash,
              chainId: walletClient.chain.id,
            }),
          });

          if (!response.ok) {
            throw new Error('Verification failed');
          }

          const result = await response.json();
          console.log('Transaction verified:', result);
          
          // Handle successful verification (update UI, etc.)
          alert('Transaction verified successfully!');
        } catch (error) {
          console.error('Verification error:', error);
          alert('Failed to verify transaction');
        } finally {
          setIsVerifying(false);
        }
      }
    };

    verifyTransaction();
  }, [isSuccess, txHash, isVerifying]);

  return (
    <div>
      <button onClick={sendTransaction} disabled={isLoading || isVerifying}>
        {isLoading ? 'Sending...' : isVerifying ? 'Verifying...' : 'Send Transaction'}
      </button>
      {txHash && (
        <p>
          Transaction Hash: {txHash}
          {isLoading && ' (Pending...)'}
          {isSuccess && ' (Confirmed)'}
        </p>
      )}
    </div>
  );
};
```

### 2. Check Transaction Status (Polling)

```typescript
const useTransactionStatus = (txHash: string | null, chainId: number) => {
  const [status, setStatus] = useState<any>(null);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!txHash) return;

    const checkStatus = async () => {
      try {
        const response = await fetch(
          `/api/v1/wallet/transaction-status?txHash=${txHash}&chainId=${chainId}`,
          {
            headers: {
              'Authorization': `Bearer ${authToken}`,
            },
            credentials: 'include',
          }
        );

        if (!response.ok) {
          throw new Error('Failed to fetch status');
        }

        const data = await response.json();
        setStatus(data);

        // Stop polling if transaction is no longer pending
        if (!data.isPending) {
          return true; // Signal to stop polling
        }
        return false;
      } catch (err) {
        setError(err as Error);
        return true; // Stop polling on error
      }
    };

    // Poll every 5 seconds
    const intervalId = setInterval(async () => {
      const shouldStop = await checkStatus();
      if (shouldStop) {
        clearInterval(intervalId);
      }
    }, 5000);

    // Initial check
    checkStatus();

    return () => clearInterval(intervalId);
  }, [txHash, chainId]);

  return { status, error };
};
```

### 3. Redux-Saga Implementation

```typescript
// actions.ts
export const verifyTransaction = (txHash: string, chainId: number) => ({
  type: 'VERIFY_TRANSACTION_REQUEST',
  payload: { txHash, chainId },
});

// sagas.ts
import { call, put, takeLatest } from 'redux-saga/effects';

function* verifyTransactionSaga(action: ReturnType<typeof verifyTransaction>) {
  try {
    const { txHash, chainId } = action.payload;
    
    const response = yield call(fetch, '/api/v1/wallet/verify-transaction', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`,
      },
      credentials: 'include',
      body: JSON.stringify({ txHash, chainId }),
    });

    if (!response.ok) {
      throw new Error('Verification failed');
    }

    const data = yield response.json();

    yield put({
      type: 'VERIFY_TRANSACTION_SUCCESS',
      payload: data,
    });
  } catch (error) {
    yield put({
      type: 'VERIFY_TRANSACTION_FAILURE',
      error: error.message,
    });
  }
}

export function* watchVerifyTransaction() {
  yield takeLatest('VERIFY_TRANSACTION_REQUEST', verifyTransactionSaga);
}
```

### 4. Complete Flow Example

```typescript
import { useAccount, useSendTransaction, useWaitForTransaction } from 'wagmi';
import { parseEther } from 'viem';

const TransactionFlow = () => {
  const { address, chain } = useAccount();
  const [txHash, setTxHash] = useState<string | null>(null);
  
  const { sendTransaction } = useSendTransaction({
    onSuccess: (hash) => {
      setTxHash(hash);
    },
  });

  const { isLoading: isConfirming, isSuccess: isConfirmed } = useWaitForTransaction({
    hash: txHash as `0x${string}`,
  });

  useEffect(() => {
    if (isConfirmed && txHash && chain) {
      // Verify with backend
      verifyWithBackend(txHash, chain.id);
    }
  }, [isConfirmed, txHash, chain]);

  const verifyWithBackend = async (hash: string, chainId: number) => {
    try {
      const response = await fetch('/api/v1/wallet/verify-transaction', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ txHash: hash, chainId }),
      });

      const result = await response.json();
      
      if (result.status === 'verified') {
        // Transaction verified - update your app state
        console.log('✅ Transaction verified on backend');
      }
    } catch (error) {
      console.error('❌ Verification failed:', error);
    }
  };

  return (
    <button
      onClick={() => {
        sendTransaction({
          to: '0x...',
          value: parseEther('0.01'),
        });
      }}
      disabled={isConfirming}
    >
      {isConfirming ? 'Confirming...' : 'Send Transaction'}
    </button>
  );
};
```

## Security Considerations

1. **Chain ID Validation**: The backend validates that the chain ID is in the list of supported networks
2. **Authentication**: All verification endpoints require user authentication
3. **Rate Limiting**: Consider implementing rate limiting on the frontend to prevent abuse
4. **Transaction Confirmation**: Always wait for transaction confirmation before verifying
5. **Error Handling**: Handle all possible error states (pending, failed, invalid chain, etc.)

## Testing

### Test with Sepolia Testnet

1. Update your `.env`:
```bash
RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
```

2. Use Sepolia testnet in your frontend:
```typescript
chainId: 11155111 // Sepolia
```

3. Get test ETH from a faucet: https://sepoliafaucet.com/

### Mock Transaction for Development

```typescript
// Mock transaction hash for testing
const mockTxHash = '0x1234567890123456789012345678901234567890123456789012345678901234';

// Test verification endpoint
const testVerification = async () => {
  const response = await fetch('/api/v1/wallet/verify-transaction', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({
      txHash: mockTxHash,
      chainId: 11155111, // Sepolia
    }),
  });
  
  console.log(await response.json());
};
```

## Troubleshooting

### Common Issues

1. **"RPC client not initialized"**
   - Ensure `RPC_URL` is set in your `.env` file
   - Restart the backend server

2. **"Transaction not found"**
   - Transaction may still be propagating
   - Verify the transaction hash is correct
   - Check you're using the correct chain ID

3. **"Chain ID mismatch"**
   - Ensure frontend and backend are configured for the same network
   - Check the chain ID matches the RPC endpoint

4. **"Transaction is still pending"**
   - Wait for transaction to be mined
   - Use `useWaitForTransaction` to wait for confirmation

5. **"Unsupported chain ID"**
   - Check that the chain is in the supported list
   - Contact backend team to add support for additional chains

## Database Schema

The Transaction model includes blockchain-specific fields:

```sql
-- Added fields to transactions table
ALTER TABLE transactions ADD COLUMN tx_hash VARCHAR(66);
ALTER TABLE transactions ADD COLUMN chain_id BIGINT;
ALTER TABLE transactions ADD COLUMN from_address VARCHAR(42);
ALTER TABLE transactions ADD COLUMN to_address VARCHAR(42);
```

## Next Steps

1. Implement rate limiting on verification endpoints
2. Add webhook support for transaction notifications
3. Support for ERC-20 token transfers
4. Multi-signature wallet support
5. Gas estimation and optimization

## Support

For issues or questions, please open an issue on GitHub or contact the development team.
