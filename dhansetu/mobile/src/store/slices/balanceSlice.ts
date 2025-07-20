import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';

interface Balance {
  denom: string;
  amount: string;
}

interface BalanceState {
  balances: Balance[];
  totalValue: string;
  isLoading: boolean;
  error: string | null;
}

const initialState: BalanceState = {
  balances: [],
  totalValue: '0',
  isLoading: false,
  error: null,
};

const balanceSlice = createSlice({
  name: 'balance',
  initialState,
  reducers: {
    setBalances: (state, action: PayloadAction<Balance[]>) => {
      state.balances = action.payload;
      state.error = null;
    },
    setTotalValue: (state, action: PayloadAction<string>) => {
      state.totalValue = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    resetBalance: () => initialState,
  },
});

export const { 
  setBalances, 
  setTotalValue, 
  setLoading, 
  setError, 
  resetBalance 
} = balanceSlice.actions;

export default balanceSlice.reducer;