import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
  duration?: number;
}

interface UIState {
  isLoading: boolean;
  loadingMessage: string;
  toasts: Toast[];
  activeModal: string | null;
  bottomSheetOpen: boolean;
  tabBarVisible: boolean;
  refreshing: {
    home: boolean;
    wallet: boolean;
    transactions: boolean;
  };
}

const initialState: UIState = {
  isLoading: false,
  loadingMessage: '',
  toasts: [],
  activeModal: null,
  bottomSheetOpen: false,
  tabBarVisible: true,
  refreshing: {
    home: false,
    wallet: false,
    transactions: false,
  },
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    setLoading: (state, action: PayloadAction<{ isLoading: boolean; message?: string }>) => {
      state.isLoading = action.payload.isLoading;
      state.loadingMessage = action.payload.message || '';
    },
    showToast: (state, action: PayloadAction<Omit<Toast, 'id'>>) => {
      const toast: Toast = {
        ...action.payload,
        id: Date.now().toString(),
        duration: action.payload.duration || 3000,
      };
      state.toasts.push(toast);
    },
    dismissToast: (state, action: PayloadAction<string>) => {
      state.toasts = state.toasts.filter(toast => toast.id !== action.payload);
    },
    setActiveModal: (state, action: PayloadAction<string | null>) => {
      state.activeModal = action.payload;
    },
    setBottomSheetOpen: (state, action: PayloadAction<boolean>) => {
      state.bottomSheetOpen = action.payload;
    },
    setTabBarVisible: (state, action: PayloadAction<boolean>) => {
      state.tabBarVisible = action.payload;
    },
    setRefreshing: (state, action: PayloadAction<{ screen: keyof UIState['refreshing']; value: boolean }>) => {
      state.refreshing[action.payload.screen] = action.payload.value;
    },
    resetUI: () => initialState,
  },
});

export const {
  setLoading,
  showToast,
  dismissToast,
  setActiveModal,
  setBottomSheetOpen,
  setTabBarVisible,
  setRefreshing,
  resetUI,
} = uiSlice.actions;

export default uiSlice.reducer;