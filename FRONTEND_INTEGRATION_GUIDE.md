# Frontend Integration Guide - Blockchain Transaction Verification

## Overview

This guide provides step-by-step instructions for integrating blockchain transaction verification into your Next.js 15 + Redux Toolkit + Redux-Saga frontend.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Backend Configuration](#backend-configuration)
4. [Frontend Implementation](#frontend-implementation)
5. [API Integration](#api-integration)
6. [Testing](#testing)
7. [Deployment](#deployment)

---

## Prerequisites

- Node.js 18+ and npm/yarn
- CaribEX Backend running with RPC_URL configured
- MetaMask or compatible Web3 wallet
- Basic understanding of wagmi and Redux

---

## Installation

### 1. Install Required Dependencies

```bash
npm install wagmi viem @wagmi/core @wagmi/chains
# or
yarn add wagmi viem @wagmi/core @wagmi/chains
```

### 2. Configure wagmi in Your Next.js App

```typescript
// app/providers.tsx
'use client';

import { WagmiConfig, createConfig, configureChains } from 'wagmi';
import { mainnet, sepolia, polygon } from 'wagmi/chains';
import { publicProvider } from 'wagmi/providers/public';
import { MetaMaskConnector } from 'wagmi/connectors/metaMask';

const { chains, publicClient, webSocketPublicClient } = configureChains(
  [mainnet, sepolia, polygon],
  [publicProvider()]
);

const config = createConfig({
  autoConnect: true,
  connectors: [
    new MetaMaskConnector({ chains }),
  ],
  publicClient,
  webSocketPublicClient,
});

export function Web3Provider({ children }: { children: React.ReactNode }) {
  return <WagmiConfig config={config}>{children}</WagmiConfig>;
}
```

---

## Backend Configuration

### 1. Environment Variables

Add to your `.env` file (backend):

```bash
RPC_URL=https://mainnet.infura.io/v3/YOUR_INFURA_KEY
# For testing:
# RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
```

### 2. Start Backend Server

```bash
cd CaribEx-backend
make run-dev
```

Verify the server is running at `http://localhost:8080`

---

## Frontend Implementation

### 1. Redux Slice for Blockchain Transactions

```typescript
// store/slices/blockchainSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface TransactionState {
  txHash: string | null;
  status: 'idle' | 'pending' | 'confirming' | 'verifying' | 'verified' | 'failed';
  error: string | null;
  verificationData: {
    from?: string;
    to?: string;
    value?: string;
    chainId?: number;
  } | null;
}

const initialState: TransactionState = {
  txHash: null,
  status: 'idle',
  error: null,
  verificationData: null,
};

const blockchainSlice = createSlice({
  name: 'blockchain',
  initialState,
  reducers: {
    setTransactionHash: (state, action: PayloadAction<string>) => {
      state.txHash = action.payload;
      state.status = 'pending';
    },
    setTransactionConfirming: (state) => {
      state.status = 'confirming';
    },
    setTransactionVerifying: (state) => {
      state.status = 'verifying';
    },
    setTransactionVerified: (state, action: PayloadAction<any>) => {
      state.status = 'verified';
      state.verificationData = action.payload;
    },
    setTransactionFailed: (state, action: PayloadAction<string>) => {
      state.status = 'failed';
      state.error = action.payload;
    },
    resetTransaction: (state) => {
      return initialState;
    },
  },
});

export const {
  setTransactionHash,
  setTransactionConfirming,
  setTransactionVerifying,
  setTransactionVerified,
  setTransactionFailed,
  resetTransaction,
} = blockchainSlice.actions;

export default blockchainSlice.reducer;
```

### 2. Redux Saga for Transaction Verification

```typescript
// store/sagas/blockchainSaga.ts
import { call, put, takeLatest, delay } from 'redux-saga/effects';
import { PayloadAction } from '@reduxjs/toolkit';
import {
  setTransactionVerifying,
  setTransactionVerified,
  setTransactionFailed,
} from '../slices/blockchainSlice';

// Action creator
export const verifyTransaction = (txHash: string, chainId: number) => ({
  type: 'blockchain/verifyTransaction',
  payload: { txHash, chainId },
});

// API call
async function verifyTransactionAPI(txHash: string, chainId: number) {
  const response = await fetch('/api/v1/wallet/verify-transaction', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Important for cookie-based auth
    body: JSON.stringify({ txHash, chainId }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Verification failed');
  }

  return response.json();
}

// Saga worker
function* verifyTransactionSaga(
  action: PayloadAction<{ txHash: string; chainId: number }>
): Generator<any, void, any> {
  try {
    const { txHash, chainId } = action.payload;
    
    yield put(setTransactionVerifying());
    
    // Call backend API
    const result = yield call(verifyTransactionAPI, txHash, chainId);
    
    yield put(setTransactionVerified(result));
    
    // Show success notification
    console.log('✅ Transaction verified:', result);
  } catch (error: any) {
    yield put(setTransactionFailed(error.message));
    console.error('❌ Verification failed:', error);
  }
}

// Saga watcher
export function* watchVerifyTransaction() {
  yield takeLatest('blockchain/verifyTransaction', verifyTransactionSaga);
}

// Root saga
export function* blockchainSaga() {
  yield watchVerifyTransaction();
}
```

### 3. Transaction Component with wagmi

```typescript
// components/SendTransaction.tsx
'use client';

import { useState, useEffect } from 'react';
import { useAccount, useSendTransaction, useWaitForTransaction } from 'wagmi';
import { parseEther } from 'viem';
import { useDispatch, useSelector } from 'react-redux';
import { setTransactionHash, verifyTransaction } from '@/store/slices/blockchainSlice';
import { RootState } from '@/store';

export default function SendTransaction() {
  const dispatch = useDispatch();
  const { address, chain } = useAccount();
  const [recipient, setRecipient] = useState('');
  const [amount, setAmount] = useState('');
  
  const { status, txHash, verificationData } = useSelector(
    (state: RootState) => state.blockchain
  );

  // Send transaction
  const { sendTransaction, isLoading: isSending } = useSendTransaction({
    onSuccess: (hash) => {
      dispatch(setTransactionHash(hash));
    },
    onError: (error) => {
      console.error('Transaction failed:', error);
    },
  });

  // Wait for transaction confirmation
  const { isLoading: isConfirming, isSuccess: isConfirmed } = useWaitForTransaction({
    hash: txHash as `0x${string}`,
    enabled: !!txHash,
  });

  // Verify transaction after confirmation
  useEffect(() => {
    if (isConfirmed && txHash && chain) {
      // Dispatch saga action to verify
      dispatch(verifyTransaction(txHash, chain.id));
    }
  }, [isConfirmed, txHash, chain, dispatch]);

  const handleSend = () => {
    if (!recipient || !amount) {
      alert('Please fill in all fields');
      return;
    }

    sendTransaction({
      to: recipient as `0x${string}`,
      value: parseEther(amount),
    });
  };

  return (
    <div className="p-6 max-w-md mx-auto bg-white rounded-xl shadow-md space-y-4">
      <h2 className="text-2xl font-bold">Send Transaction</h2>
      
      <div>
        <label className="block text-sm font-medium text-gray-700">
          Recipient Address
        </label>
        <input
          type="text"
          value={recipient}
          onChange={(e) => setRecipient(e.target.value)}
          placeholder="0x..."
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700">
          Amount (ETH)
        </label>
        <input
          type="text"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          placeholder="0.01"
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
        />
      </div>

      <button
        onClick={handleSend}
        disabled={isSending || isConfirming || status === 'verifying'}
        className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {isSending && 'Sending...'}
        {isConfirming && 'Confirming...'}
        {status === 'verifying' && 'Verifying...'}
        {!isSending && !isConfirming && status !== 'verifying' && 'Send Transaction'}
      </button>

      {/* Transaction Status */}
      {txHash && (
        <div className="space-y-2">
          <div className="text-sm">
            <span className="font-medium">TX Hash:</span>{' '}
            <a
              href={`https://etherscan.io/tx/${txHash}`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-500 hover:underline"
            >
              {txHash.substring(0, 10)}...{txHash.substring(txHash.length - 8)}
            </a>
          </div>

          <div className="text-sm">
            <span className="font-medium">Status:</span>{' '}
            <span className={`capitalize ${
              status === 'verified' ? 'text-green-600' : 
              status === 'failed' ? 'text-red-600' : 
              'text-yellow-600'
            }`}>
              {status}
            </span>
          </div>

          {status === 'verified' && verificationData && (
            <div className="bg-green-50 p-3 rounded text-sm space-y-1">
              <div>✅ Transaction Verified!</div>
              <div>From: {verificationData.from}</div>
              <div>To: {verificationData.to}</div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
```

### 4. Transaction Status Polling (Alternative Approach)

```typescript
// hooks/useTransactionStatus.ts
import { useState, useEffect } from 'react';

interface TransactionStatus {
  status: string;
  isPending: boolean;
  from?: string;
  to?: string;
  value?: string;
}

export function useTransactionStatus(txHash: string | null, chainId: number) {
  const [status, setStatus] = useState<TransactionStatus | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isPolling, setIsPolling] = useState(false);

  useEffect(() => {
    if (!txHash) return;

    setIsPolling(true);
    let intervalId: NodeJS.Timeout;

    const checkStatus = async () => {
      try {
        const response = await fetch(
          `/api/v1/wallet/transaction-status?txHash=${txHash}&chainId=${chainId}`,
          {
            credentials: 'include',
          }
        );

        if (!response.ok) {
          throw new Error('Failed to fetch status');
        }

        const data = await response.json();
        setStatus(data);

        // Stop polling if transaction is confirmed
        if (!data.isPending) {
          setIsPolling(false);
          clearInterval(intervalId);
        }
      } catch (err) {
        setError(err as Error);
        setIsPolling(false);
        clearInterval(intervalId);
      }
    };

    // Initial check
    checkStatus();

    // Poll every 5 seconds
    intervalId = setInterval(checkStatus, 5000);

    return () => {
      clearInterval(intervalId);
      setIsPolling(false);
    };
  }, [txHash, chainId]);

  return { status, error, isPolling };
}

// Usage in component
function TransactionMonitor({ txHash, chainId }: { txHash: string; chainId: number }) {
  const { status, error, isPolling } = useTransactionStatus(txHash, chainId);

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  if (isPolling) {
    return <div>Checking transaction status...</div>;
  }

  if (!status) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      {status.isPending ? (
        <div>Transaction pending...</div>
      ) : (
        <div>
          <div>Transaction confirmed!</div>
          <div>From: {status.from}</div>
          <div>To: {status.to}</div>
          <div>Value: {status.value} wei</div>
        </div>
      )}
    </div>
  );
}
```

---

## API Integration

### 1. Create API Service

```typescript
// services/blockchainService.ts
import { getAuthToken } from '@/utils/auth';

export interface VerifyTransactionRequest {
  txHash: string;
  chainId: number;
}

export interface VerifyTransactionResponse {
  status: string;
  txHash: string;
  message: string;
  from?: string;
  to?: string;
  chainId?: number;
}

class BlockchainService {
  private baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  async verifyTransaction(
    request: VerifyTransactionRequest
  ): Promise<VerifyTransactionResponse> {
    const response = await fetch(`${this.baseURL}/v1/wallet/verify-transaction`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`,
      },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Verification failed');
    }

    return response.json();
  }

  async getTransactionStatus(
    txHash: string,
    chainId: number
  ): Promise<any> {
    const response = await fetch(
      `${this.baseURL}/v1/wallet/transaction-status?txHash=${txHash}&chainId=${chainId}`,
      {
        headers: {
          'Authorization': `Bearer ${getAuthToken()}`,
        },
        credentials: 'include',
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch transaction status');
    }

    return response.json();
  }
}

export const blockchainService = new BlockchainService();
```

### 2. Environment Configuration

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Testing

### 1. Test with Sepolia Testnet

1. Configure backend:
```bash
RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
```

2. Get test ETH:
   - Visit https://sepoliafaucet.com/
   - Connect your wallet
   - Request test ETH

3. Update frontend to use Sepolia:
```typescript
const { chains } = configureChains(
  [sepolia], // Use Sepolia for testing
  [publicProvider()]
);
```

### 2. Mock Transaction for Development

```typescript
// For testing without real transactions
const mockVerification = async () => {
  return {
    status: 'verified',
    txHash: '0x1234567890123456789012345678901234567890123456789012345678901234',
    message: 'Transaction successfully verified',
    from: '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1',
    to: '0x1234567890123456789012345678901234567890',
    chainId: 11155111,
  };
};
```

### 3. Integration Tests

```typescript
// __tests__/blockchain.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { store } from '@/store';
import SendTransaction from '@/components/SendTransaction';

describe('SendTransaction', () => {
  it('should verify transaction after confirmation', async () => {
    render(
      <Provider store={store}>
        <SendTransaction />
      </Provider>
    );

    // Fill in form
    const recipientInput = screen.getByPlaceholderText('0x...');
    const amountInput = screen.getByPlaceholderText('0.01');
    
    fireEvent.change(recipientInput, {
      target: { value: '0x1234567890123456789012345678901234567890' }
    });
    fireEvent.change(amountInput, { target: { value: '0.01' } });

    // Click send button
    const sendButton = screen.getByText('Send Transaction');
    fireEvent.click(sendButton);

    // Wait for verification
    await waitFor(
      () => {
        expect(screen.getByText(/Transaction Verified/i)).toBeInTheDocument();
      },
      { timeout: 30000 }
    );
  });
});
```

---

## Deployment

### 1. Production Environment Variables

```bash
# Backend .env
RPC_URL=https://mainnet.infura.io/v3/YOUR_PRODUCTION_KEY
ENV=production

# Frontend .env.production
NEXT_PUBLIC_API_URL=https://api.caribex.com
```

### 2. Security Checklist

- [ ] Use HTTPS for all API calls
- [ ] Validate all user inputs
- [ ] Never expose private keys
- [ ] Implement rate limiting
- [ ] Use environment variables for sensitive data
- [ ] Enable CORS only for trusted origins
- [ ] Implement proper error handling
- [ ] Log security events

### 3. Monitoring

Add error tracking:

```typescript
// utils/errorTracking.ts
export function logError(error: Error, context: string) {
  // Send to your error tracking service (Sentry, etc.)
  console.error(`[${context}]`, error);
  
  // In production, send to backend
  if (process.env.NODE_ENV === 'production') {
    fetch('/api/log-error', {
      method: 'POST',
      body: JSON.stringify({
        error: error.message,
        stack: error.stack,
        context,
      }),
    }).catch(console.error);
  }
}
```

---

## Troubleshooting

### Common Issues

1. **"RPC client not initialized"**
   - Ensure `RPC_URL` is set in backend `.env`
   - Restart backend server

2. **CORS errors**
   - Check `ALLOWED_ORIGINS` in backend configuration
   - Ensure frontend URL is whitelisted

3. **Transaction not found**
   - Wait for transaction to propagate (5-30 seconds)
   - Verify correct chain ID

4. **Authentication errors**
   - Ensure SIWE authentication is complete
   - Check session cookie is being sent

---

## Support

For additional help:
- Backend API: See `docs/BLOCKCHAIN_VERIFICATION.md`
- Full API Reference: See `docs/API.md`
- SIWE Auth: See `docs/SIWE_AUTH.md`

---

## Summary

This integration enables:
✅ Send transactions via wagmi
✅ Wait for on-chain confirmation
✅ Verify with backend API
✅ Log verified transactions to database
✅ Multi-chain support
✅ Comprehensive error handling
✅ Production-ready architecture

Your frontend now has full blockchain transaction verification integrated with the CaribEX backend!
