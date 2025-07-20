/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { HDWallet, WalletAccount } from '@services/wallet/hdWallet';
import { SecureStorage } from '@services/wallet/secureStorage';

interface WalletState {
  hasWallet: boolean;
  isWalletCreated: boolean;
  address: string | null;
  mnemonic: string | null;
  currentAddress: string | null;
  currentCoin: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN';
  accounts: {
    deshchain: WalletAccount[];
    ethereum: WalletAccount[];
    bitcoin: WalletAccount[];
  };
  dhanPataAddress: string | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: WalletState = {
  hasWallet: false,
  isWalletCreated: false,
  address: null,
  mnemonic: null,
  currentAddress: null,
  currentCoin: 'DESHCHAIN',
  accounts: {
    deshchain: [],
    ethereum: [],
    bitcoin: [],
  },
  dhanPataAddress: null,
  isLoading: false,
  error: null,
};

// Async thunks
export const createWallet = createAsyncThunk(
  'wallet/create',
  async (password?: string) => {
    const wallet = await HDWallet.generate();
    const secureStorage = SecureStorage.getInstance();
    
    // Store encrypted mnemonic
    await secureStorage.setItem('mnemonic', wallet.getMnemonic(), 'WALLET');
    
    // Derive initial accounts
    const deshAccount = await wallet.deriveDeshChainAccount(0);
    const ethAccount = wallet.deriveEthereumAccount(0);
    const btcAccount = wallet.deriveBitcoinAccount(0);
    
    return {
      accounts: {
        deshchain: [deshAccount],
        ethereum: [ethAccount],
        bitcoin: [btcAccount],
      },
      currentAddress: deshAccount.address,
    };
  }
);

export const importWallet = createAsyncThunk(
  'wallet/import',
  async ({ mnemonic, password }: { mnemonic: string; password?: string }) => {
    const wallet = await HDWallet.fromMnemonic(mnemonic, password);
    const secureStorage = SecureStorage.getInstance();
    
    // Store encrypted mnemonic
    await secureStorage.setItem('mnemonic', mnemonic, 'WALLET');
    
    // Derive initial accounts
    const deshAccount = await wallet.deriveDeshChainAccount(0);
    const ethAccount = wallet.deriveEthereumAccount(0);
    const btcAccount = wallet.deriveBitcoinAccount(0);
    
    return {
      accounts: {
        deshchain: [deshAccount],
        ethereum: [ethAccount],
        bitcoin: [btcAccount],
      },
      currentAddress: deshAccount.address,
    };
  }
);

export const loadWallet = createAsyncThunk(
  'wallet/load',
  async () => {
    const secureStorage = SecureStorage.getInstance();
    const mnemonic = await secureStorage.getItem('mnemonic', 'WALLET');
    
    if (!mnemonic) {
      throw new Error('No wallet found');
    }
    
    const wallet = await HDWallet.fromMnemonic(mnemonic);
    
    // Load accounts
    const deshAccount = await wallet.deriveDeshChainAccount(0);
    const ethAccount = wallet.deriveEthereumAccount(0);
    const btcAccount = wallet.deriveBitcoinAccount(0);
    
    return {
      accounts: {
        deshchain: [deshAccount],
        ethereum: [ethAccount],
        bitcoin: [btcAccount],
      },
      currentAddress: deshAccount.address,
    };
  }
);

export const createDhanPataAddress = createAsyncThunk(
  'wallet/createDhanPata',
  async ({ username, address }: { username: string; address: string }) => {
    // This would call the DeshChain API to register DhanPata address
    // For now, return the formatted address
    return `${username}@dhan`;
  }
);

const walletSlice = createSlice({
  name: 'wallet',
  initialState,
  reducers: {
    setCurrentCoin: (state, action: PayloadAction<'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN'>) => {
      state.currentCoin = action.payload;
      
      // Update current address based on selected coin
      switch (action.payload) {
        case 'DESHCHAIN':
          state.currentAddress = state.accounts.deshchain[0]?.address || null;
          break;
        case 'ETHEREUM':
          state.currentAddress = state.accounts.ethereum[0]?.address || null;
          break;
        case 'BITCOIN':
          state.currentAddress = state.accounts.bitcoin[0]?.address || null;
          break;
      }
    },
    setCurrentAddress: (state, action: PayloadAction<string>) => {
      state.currentAddress = action.payload;
    },
    clearWallet: (state) => {
      return initialState;
    },
    setWalletCreated: (state, action: PayloadAction<boolean>) => {
      state.isWalletCreated = action.payload;
      state.hasWallet = action.payload;
    },
    setAddress: (state, action: PayloadAction<string>) => {
      state.address = action.payload;
      state.currentAddress = action.payload;
    },
    setMnemonic: (state, action: PayloadAction<string>) => {
      state.mnemonic = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      // Create wallet
      .addCase(createWallet.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createWallet.fulfilled, (state, action) => {
        state.isLoading = false;
        state.hasWallet = true;
        state.accounts = action.payload.accounts;
        state.currentAddress = action.payload.currentAddress;
      })
      .addCase(createWallet.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to create wallet';
      })
      
      // Import wallet
      .addCase(importWallet.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(importWallet.fulfilled, (state, action) => {
        state.isLoading = false;
        state.hasWallet = true;
        state.accounts = action.payload.accounts;
        state.currentAddress = action.payload.currentAddress;
      })
      .addCase(importWallet.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to import wallet';
      })
      
      // Load wallet
      .addCase(loadWallet.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loadWallet.fulfilled, (state, action) => {
        state.isLoading = false;
        state.hasWallet = true;
        state.accounts = action.payload.accounts;
        state.currentAddress = action.payload.currentAddress;
      })
      .addCase(loadWallet.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to load wallet';
      })
      
      // Create DhanPata
      .addCase(createDhanPataAddress.fulfilled, (state, action) => {
        state.dhanPataAddress = action.payload;
      });
  },
});

export const { 
  setCurrentCoin, 
  setCurrentAddress, 
  clearWallet,
  setWalletCreated,
  setAddress,
  setMnemonic,
} = walletSlice.actions;

export default walletSlice.reducer;