import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface SettingsState {
  language: string;
  currency: string;
  theme: 'light' | 'dark' | 'auto';
  notifications: {
    enabled: boolean;
    transactions: boolean;
    updates: boolean;
    marketing: boolean;
  };
  security: {
    biometricEnabled: boolean;
    autoLockTimeout: number; // in minutes
    hideBalances: boolean;
  };
  network: {
    rpcUrl: string;
    chainId: string;
    explorerUrl: string;
  };
}

const initialState: SettingsState = {
  language: 'en',
  currency: 'INR',
  theme: 'light',
  notifications: {
    enabled: true,
    transactions: true,
    updates: true,
    marketing: false,
  },
  security: {
    biometricEnabled: false,
    autoLockTimeout: 5,
    hideBalances: false,
  },
  network: {
    rpcUrl: 'https://rpc.deshchain.com',
    chainId: 'deshchain-1',
    explorerUrl: 'https://explorer.deshchain.com',
  },
};

const settingsSlice = createSlice({
  name: 'settings',
  initialState,
  reducers: {
    setLanguage: (state, action: PayloadAction<string>) => {
      state.language = action.payload;
    },
    setCurrency: (state, action: PayloadAction<string>) => {
      state.currency = action.payload;
    },
    setTheme: (state, action: PayloadAction<'light' | 'dark' | 'auto'>) => {
      state.theme = action.payload;
    },
    updateNotifications: (state, action: PayloadAction<Partial<SettingsState['notifications']>>) => {
      state.notifications = { ...state.notifications, ...action.payload };
    },
    updateSecurity: (state, action: PayloadAction<Partial<SettingsState['security']>>) => {
      state.security = { ...state.security, ...action.payload };
    },
    updateNetwork: (state, action: PayloadAction<Partial<SettingsState['network']>>) => {
      state.network = { ...state.network, ...action.payload };
    },
    resetSettings: () => initialState,
  },
});

export const {
  setLanguage,
  setCurrency,
  setTheme,
  updateNotifications,
  updateSecurity,
  updateNetwork,
  resetSettings,
} = settingsSlice.actions;

export default settingsSlice.reducer;