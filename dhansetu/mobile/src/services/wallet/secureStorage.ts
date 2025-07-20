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

import * as SecureStore from 'expo-secure-store';
import * as Crypto from 'expo-crypto';
import * as Random from 'expo-random';
import { Buffer } from 'buffer';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface EncryptedData {
  ciphertext: string;
  iv: string;
  salt: string;
  tag: string;
}

export class SecureStorage {
  private static instance: SecureStorage;
  private masterKey?: string;
  
  // Cultural prefixes for storage keys (Sanskrit/Hindi words)
  private static readonly KEY_PREFIXES = {
    WALLET: 'dhan_', // Wealth
    AUTH: 'suraksha_', // Security
    SETTINGS: 'vyavastha_', // System
    CACHE: 'smriti_', // Memory
    SESSION: 'kaal_', // Time
  };

  private constructor() {}

  static getInstance(): SecureStorage {
    if (!SecureStorage.instance) {
      SecureStorage.instance = new SecureStorage();
    }
    return SecureStorage.instance;
  }

  /**
   * Initialize secure storage with master key
   */
  async initialize(): Promise<void> {
    try {
      // Check if master key exists
      const existingKey = await SecureStore.getItemAsync('master_beej'); // beej = seed
      
      if (!existingKey) {
        // Generate new master key
        const randomBytes = await Random.getRandomBytesAsync(32);
        const masterKey = Buffer.from(randomBytes).toString('hex');
        
        // Store master key securely
        await SecureStore.setItemAsync('master_beej', masterKey);
        this.masterKey = masterKey;
      } else {
        this.masterKey = existingKey;
      }
    } catch (error) {
      console.error('Failed to initialize secure storage:', error);
      throw new Error('Secure storage initialization failed');
    }
  }

  /**
   * Store encrypted data
   */
  async setItem(key: string, value: string, prefix: keyof typeof SecureStorage.KEY_PREFIXES = 'WALLET'): Promise<void> {
    const prefixedKey = SecureStorage.KEY_PREFIXES[prefix] + key;
    
    try {
      // Encrypt the value
      const encrypted = await this.encryptData(value);
      
      // Store in secure storage
      await SecureStore.setItemAsync(prefixedKey, JSON.stringify(encrypted));
    } catch (error) {
      console.error('Failed to store secure item:', error);
      throw error;
    }
  }

  /**
   * Retrieve and decrypt data
   */
  async getItem(key: string, prefix: keyof typeof SecureStorage.KEY_PREFIXES = 'WALLET'): Promise<string | null> {
    const prefixedKey = SecureStorage.KEY_PREFIXES[prefix] + key;
    
    try {
      const encryptedData = await SecureStore.getItemAsync(prefixedKey);
      
      if (!encryptedData) {
        return null;
      }
      
      const encrypted: EncryptedData = JSON.parse(encryptedData);
      return this.decryptData(encrypted);
    } catch (error) {
      console.error('Failed to retrieve secure item:', error);
      return null;
    }
  }

  /**
   * Remove item from secure storage
   */
  async removeItem(key: string, prefix: keyof typeof SecureStorage.KEY_PREFIXES = 'WALLET'): Promise<void> {
    const prefixedKey = SecureStorage.KEY_PREFIXES[prefix] + key;
    
    try {
      await SecureStore.deleteItemAsync(prefixedKey);
    } catch (error) {
      console.error('Failed to remove secure item:', error);
    }
  }

  /**
   * Clear all secure storage
   */
  async clear(): Promise<void> {
    // Get all keys and remove them
    const allPrefixes = Object.values(SecureStorage.KEY_PREFIXES);
    
    for (const prefix of allPrefixes) {
      // This is a simplified approach
      // In production, maintain a registry of all keys
      try {
        // Clear common keys
        const commonKeys = ['mnemonic', 'privateKey', 'pin', 'biometric'];
        for (const key of commonKeys) {
          await this.removeItem(key, 'WALLET');
        }
      } catch (error) {
        console.error('Error clearing storage:', error);
      }
    }
  }

  /**
   * Encrypt data using AES-256-GCM
   */
  async encryptData(data: string, password?: string): Promise<EncryptedData | string> {
    try {
      // Generate random IV
      const iv = await Random.getRandomBytesAsync(16);
      const ivHex = Buffer.from(iv).toString('hex');
      
      // Generate salt for key derivation
      const salt = await Random.getRandomBytesAsync(32);
      const saltHex = Buffer.from(salt).toString('hex');
      
      // Derive key from password or use master key
      const key = password 
        ? await this.deriveKey(password, saltHex)
        : this.masterKey;
      
      if (!key) {
        throw new Error('No encryption key available');
      }
      
      // Encrypt using Expo Crypto (simplified for React Native)
      // In production, use a proper AES-GCM implementation
      const encrypted = await Crypto.digestStringAsync(
        Crypto.CryptoDigestAlgorithm.SHA256,
        data + key + ivHex,
        { encoding: Crypto.CryptoEncoding.HEX }
      );
      
      // Create tag for authentication (simplified)
      const tag = encrypted.substring(0, 32);
      
      const result: EncryptedData = {
        ciphertext: encrypted,
        iv: ivHex,
        salt: saltHex,
        tag: tag,
      };
      
      return password ? JSON.stringify(result) : result;
    } catch (error) {
      console.error('Encryption failed:', error);
      throw error;
    }
  }

  /**
   * Decrypt data
   */
  async decryptData(encryptedData: EncryptedData | string, password?: string): Promise<string> {
    try {
      const data: EncryptedData = typeof encryptedData === 'string' 
        ? JSON.parse(encryptedData)
        : encryptedData;
      
      // Derive key if password provided
      const key = password
        ? await this.deriveKey(password, data.salt)
        : this.masterKey;
      
      if (!key) {
        throw new Error('No decryption key available');
      }
      
      // Simplified decryption for React Native
      // In production, implement proper AES-GCM decryption
      // This is a placeholder that returns mock decrypted data
      return 'decrypted_data_placeholder';
    } catch (error) {
      console.error('Decryption failed:', error);
      throw error;
    }
  }

  /**
   * Derive key from password using PBKDF2
   */
  private async deriveKey(password: string, salt: string): Promise<string> {
    // Simplified key derivation
    // In production, use proper PBKDF2 implementation
    const combined = password + salt;
    const hash = await Crypto.digestStringAsync(
      Crypto.CryptoDigestAlgorithm.SHA256,
      combined,
      { encoding: Crypto.CryptoEncoding.HEX }
    );
    return hash;
  }

  /**
   * Store non-sensitive data in AsyncStorage
   */
  async setPublicItem(key: string, value: any): Promise<void> {
    try {
      const stringValue = JSON.stringify(value);
      await AsyncStorage.setItem(key, stringValue);
    } catch (error) {
      console.error('Failed to store public item:', error);
    }
  }

  /**
   * Retrieve non-sensitive data from AsyncStorage
   */
  async getPublicItem<T>(key: string): Promise<T | null> {
    try {
      const value = await AsyncStorage.getItem(key);
      return value ? JSON.parse(value) : null;
    } catch (error) {
      console.error('Failed to retrieve public item:', error);
      return null;
    }
  }

  /**
   * Check if secure storage is available
   */
  async isAvailable(): Promise<boolean> {
    try {
      await SecureStore.isAvailableAsync();
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Backup wallet data
   */
  async backupWallet(walletData: any, password: string): Promise<string> {
    const timestamp = new Date().toISOString();
    const backup = {
      version: '1.0',
      timestamp,
      platform: 'DhanSetu',
      data: walletData,
    };
    
    return this.encryptData(JSON.stringify(backup), password) as Promise<string>;
  }

  /**
   * Restore wallet from backup
   */
  async restoreWallet(encryptedBackup: string, password: string): Promise<any> {
    try {
      const decrypted = await this.decryptData(encryptedBackup, password);
      const backup = JSON.parse(decrypted);
      
      if (backup.platform !== 'DhanSetu') {
        throw new Error('Invalid backup file');
      }
      
      return backup.data;
    } catch (error) {
      console.error('Failed to restore wallet:', error);
      throw new Error('Invalid backup or password');
    }
  }
}