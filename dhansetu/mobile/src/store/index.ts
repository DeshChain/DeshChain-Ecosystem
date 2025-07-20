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

import { configureStore } from '@reduxjs/toolkit';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

// Import reducers
import authReducer from './slices/authSlice';
import walletReducer from './slices/walletSlice';
import balanceReducer from './slices/balanceSlice';
import transactionReducer from './slices/transactionSlice';
import dexReducer from './slices/dexSlice';
import sikkebaazReducer from './slices/sikkebaazSlice';
import surakshaReducer from './slices/surakshaSlice';
import culturalReducer from './slices/culturalSlice';
import settingsReducer from './slices/settingsSlice';
import uiReducer from './slices/uiSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    wallet: walletReducer,
    balance: balanceReducer,
    transaction: transactionReducer,
    dex: dexReducer,
    sikkebaaz: sikkebaazReducer,
    suraksha: surakshaReducer,
    cultural: culturalReducer,
    settings: settingsReducer,
    ui: uiReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Ignore these action types
        ignoredActions: ['wallet/setWallet', 'transaction/addTransaction'],
        // Ignore these field paths in all actions
        ignoredActionPaths: ['payload.timestamp', 'payload.wallet'],
        // Ignore these paths in the state
        ignoredPaths: ['wallet.hdWallet', 'transaction.pendingTx'],
      },
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Typed hooks
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;