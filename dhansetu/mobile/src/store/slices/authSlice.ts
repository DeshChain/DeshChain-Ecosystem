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
import * as LocalAuthentication from 'expo-local-authentication';
import { SecureStorage } from '@services/wallet/secureStorage';

interface AuthState {
  isAuthenticated: boolean;
  isBiometricEnabled: boolean;
  isBiometricAvailable: boolean;
  lastAuthTime: number | null;
  authMethod: 'pin' | 'biometric' | null;
  pinAttempts: number;
  isLocked: boolean;
  lockUntil: number | null;
  pinHash: string | null;
  biometricEnabled: boolean;
}

const initialState: AuthState = {
  isAuthenticated: false,
  isBiometricEnabled: false,
  isBiometricAvailable: false,
  lastAuthTime: null,
  authMethod: null,
  pinAttempts: 0,
  isLocked: false,
  lockUntil: null,
  pinHash: null,
  biometricEnabled: false,
};

const MAX_PIN_ATTEMPTS = 3;
const LOCK_DURATION = 5 * 60 * 1000; // 5 minutes

// Async thunks
export const checkBiometricAvailability = createAsyncThunk(
  'auth/checkBiometric',
  async () => {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    const isEnrolled = await LocalAuthentication.isEnrolledAsync();
    return hasHardware && isEnrolled;
  }
);

export const authenticateWithBiometric = createAsyncThunk(
  'auth/biometric',
  async () => {
    const result = await LocalAuthentication.authenticateAsync({
      promptMessage: 'Authenticate to access DhanSetu',
      cancelLabel: 'Cancel',
      fallbackLabel: 'Use PIN',
    });
    return result.success;
  }
);

export const verifyPin = createAsyncThunk(
  'auth/verifyPin',
  async (pin: string) => {
    const secureStorage = SecureStorage.getInstance();
    const storedPin = await secureStorage.getItem('pin', 'AUTH');
    return storedPin === pin;
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setAuthenticated: (state, action: PayloadAction<boolean>) => {
      state.isAuthenticated = action.payload;
      state.lastAuthTime = action.payload ? Date.now() : null;
      if (action.payload) {
        state.pinAttempts = 0;
        state.isLocked = false;
        state.lockUntil = null;
      }
    },
    setBiometricEnabled: (state, action: PayloadAction<boolean>) => {
      state.isBiometricEnabled = action.payload;
    },
    setAuthMethod: (state, action: PayloadAction<'pin' | 'biometric'>) => {
      state.authMethod = action.payload;
    },
    incrementPinAttempts: (state) => {
      state.pinAttempts += 1;
      if (state.pinAttempts >= MAX_PIN_ATTEMPTS) {
        state.isLocked = true;
        state.lockUntil = Date.now() + LOCK_DURATION;
      }
    },
    resetPinAttempts: (state) => {
      state.pinAttempts = 0;
    },
    checkLockStatus: (state) => {
      if (state.isLocked && state.lockUntil && Date.now() > state.lockUntil) {
        state.isLocked = false;
        state.lockUntil = null;
        state.pinAttempts = 0;
      }
    },
    logout: (state) => {
      state.isAuthenticated = false;
      state.lastAuthTime = null;
      state.authMethod = null;
    },
    setPinHash: (state, action: PayloadAction<string>) => {
      state.pinHash = action.payload;
    },
    reset: () => initialState,
  },
  extraReducers: (builder) => {
    builder
      .addCase(checkBiometricAvailability.fulfilled, (state, action) => {
        state.isBiometricAvailable = action.payload;
      })
      .addCase(authenticateWithBiometric.fulfilled, (state, action) => {
        if (action.payload) {
          state.isAuthenticated = true;
          state.lastAuthTime = Date.now();
          state.authMethod = 'biometric';
        }
      })
      .addCase(verifyPin.fulfilled, (state, action) => {
        if (action.payload) {
          state.isAuthenticated = true;
          state.lastAuthTime = Date.now();
          state.authMethod = 'pin';
          state.pinAttempts = 0;
        } else {
          state.pinAttempts += 1;
          if (state.pinAttempts >= MAX_PIN_ATTEMPTS) {
            state.isLocked = true;
            state.lockUntil = Date.now() + LOCK_DURATION;
          }
        }
      });
  },
});

export const {
  setAuthenticated,
  setBiometricEnabled,
  setAuthMethod,
  incrementPinAttempts,
  resetPinAttempts,
  checkLockStatus,
  logout,
  setPinHash,
  reset,
} = authSlice.actions;

export default authSlice.reducer;