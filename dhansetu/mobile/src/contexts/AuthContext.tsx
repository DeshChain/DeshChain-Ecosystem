import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useAppDispatch, useAppSelector } from '@store/index';
import { setAuthenticated, setPinHash, setBiometricEnabled } from '@store/slices/authSlice';
import * as LocalAuthentication from 'expo-local-authentication';
import * as SecureStore from 'expo-secure-store';

interface AuthContextType {
  isAuthenticated: boolean;
  isPinSet: boolean;
  isBiometricEnabled: boolean;
  authenticate: (pin: string) => Promise<boolean>;
  authenticateWithBiometrics: () => Promise<boolean>;
  setupPin: (pin: string) => Promise<void>;
  enableBiometrics: () => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const dispatch = useAppDispatch();
  const { isAuthenticated, pinHash, biometricEnabled } = useAppSelector(state => state.auth);
  const [isPinSet, setIsPinSet] = useState(false);

  useEffect(() => {
    checkPinStatus();
  }, []);

  const checkPinStatus = async () => {
    const storedPin = await SecureStore.getItemAsync('pinHash');
    setIsPinSet(!!storedPin);
  };

  const authenticate = async (pin: string): Promise<boolean> => {
    const storedPin = await SecureStore.getItemAsync('pinHash');
    if (storedPin === pin) {
      dispatch(setAuthenticated(true));
      return true;
    }
    return false;
  };

  const authenticateWithBiometrics = async (): Promise<boolean> => {
    try {
      const hasHardware = await LocalAuthentication.hasHardwareAsync();
      if (!hasHardware) return false;

      const isEnrolled = await LocalAuthentication.isEnrolledAsync();
      if (!isEnrolled) return false;

      const result = await LocalAuthentication.authenticateAsync({
        promptMessage: 'Authenticate to access DhanSetu',
        disableDeviceFallback: false,
      });

      if (result.success) {
        dispatch(setAuthenticated(true));
        return true;
      }
      return false;
    } catch (error) {
      console.error('Biometric authentication error:', error);
      return false;
    }
  };

  const setupPin = async (pin: string): Promise<void> => {
    await SecureStore.setItemAsync('pinHash', pin);
    dispatch(setPinHash(pin));
    setIsPinSet(true);
  };

  const enableBiometrics = async (): Promise<void> => {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    if (hasHardware) {
      dispatch(setBiometricEnabled(true));
      await SecureStore.setItemAsync('biometricEnabled', 'true');
    }
  };

  const logout = () => {
    dispatch(setAuthenticated(false));
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isPinSet,
        isBiometricEnabled: biometricEnabled,
        authenticate,
        authenticateWithBiometrics,
        setupPin,
        enableBiometrics,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};