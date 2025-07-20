import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

export class AppInitializationService {
  private static instance: AppInitializationService;

  private constructor() {}

  static getInstance(): AppInitializationService {
    if (!AppInitializationService.instance) {
      AppInitializationService.instance = new AppInitializationService();
    }
    return AppInitializationService.instance;
  }

  async initialize(): Promise<void> {
    try {
      // Check if this is first launch
      const isFirstLaunch = await this.checkFirstLaunch();
      
      // Initialize secure storage
      await this.initializeSecureStorage();
      
      // Load user preferences
      await this.loadUserPreferences();
      
      // Setup default values if first launch
      if (isFirstLaunch) {
        await this.setupDefaults();
      }
      
      // Initialize network configuration
      await this.initializeNetwork();
      
    } catch (error) {
      console.error('App initialization failed:', error);
      throw error;
    }
  }

  private async checkFirstLaunch(): Promise<boolean> {
    try {
      const hasLaunched = await AsyncStorage.getItem('hasLaunched');
      if (!hasLaunched) {
        await AsyncStorage.setItem('hasLaunched', 'true');
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to check first launch:', error);
      return false;
    }
  }

  private async initializeSecureStorage(): Promise<void> {
    // Platform-specific secure storage initialization
    if (Platform.OS === 'ios') {
      // iOS specific initialization
    } else if (Platform.OS === 'android') {
      // Android specific initialization
    }
  }

  private async loadUserPreferences(): Promise<void> {
    try {
      const preferences = await AsyncStorage.getItem('userPreferences');
      if (preferences) {
        // Load preferences into app state
        const prefs = JSON.parse(preferences);
        // Dispatch to Redux store or context
      }
    } catch (error) {
      console.error('Failed to load user preferences:', error);
    }
  }

  private async setupDefaults(): Promise<void> {
    const defaults = {
      language: 'en',
      currency: 'INR',
      theme: 'light',
      biometricEnabled: false,
    };
    
    await AsyncStorage.setItem('userPreferences', JSON.stringify(defaults));
  }

  private async initializeNetwork(): Promise<void> {
    // Initialize blockchain network configuration
    const networkConfig = {
      rpcUrl: 'https://rpc.deshchain.com',
      chainId: 'deshchain-1',
      explorerUrl: 'https://explorer.deshchain.com',
    };
    
    await AsyncStorage.setItem('networkConfig', JSON.stringify(networkConfig));
  }

  async reset(): Promise<void> {
    // Clear all stored data
    await AsyncStorage.clear();
    await SecureStore.deleteItemAsync('wallet');
    await SecureStore.deleteItemAsync('pinHash');
    await SecureStore.deleteItemAsync('biometricEnabled');
  }
}

const initializeApp = async () => {
  const instance = AppInitializationService.getInstance();
  await instance.initialize();
};

export default initializeApp;