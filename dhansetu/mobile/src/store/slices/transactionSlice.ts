import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Transaction {
  hash: string;
  type: 'send' | 'receive' | 'swap' | 'delegate' | 'undelegate';
  status: 'pending' | 'success' | 'failed';
  amount: string;
  denom: string;
  from: string;
  to: string;
  fee: string;
  timestamp: number;
  memo?: string;
}

interface TransactionState {
  transactions: Transaction[];
  pendingTransactions: Transaction[];
  isLoading: boolean;
  error: string | null;
}

const initialState: TransactionState = {
  transactions: [],
  pendingTransactions: [],
  isLoading: false,
  error: null,
};

const transactionSlice = createSlice({
  name: 'transaction',
  initialState,
  reducers: {
    addTransaction: (state, action: PayloadAction<Transaction>) => {
      state.transactions.unshift(action.payload);
      if (action.payload.status === 'pending') {
        state.pendingTransactions.push(action.payload);
      }
    },
    updateTransactionStatus: (state, action: PayloadAction<{ hash: string; status: Transaction['status'] }>) => {
      const { hash, status } = action.payload;
      
      // Update in transactions array
      const txIndex = state.transactions.findIndex(tx => tx.hash === hash);
      if (txIndex !== -1) {
        state.transactions[txIndex].status = status;
      }
      
      // Remove from pending if no longer pending
      if (status !== 'pending') {
        state.pendingTransactions = state.pendingTransactions.filter(tx => tx.hash !== hash);
      }
    },
    setTransactions: (state, action: PayloadAction<Transaction[]>) => {
      state.transactions = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    resetTransactions: () => initialState,
  },
});

export const {
  addTransaction,
  updateTransactionStatus,
  setTransactions,
  setLoading,
  setError,
  resetTransactions,
} = transactionSlice.actions;

export default transactionSlice.reducer;