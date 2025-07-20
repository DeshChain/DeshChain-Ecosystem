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

import * as bip39 from 'bip39';
import { BIP32Factory } from 'bip32';
import * as ecc from 'tiny-secp256k1';
import { ethers } from 'ethers';
import { 
  DirectSecp256k1HdWallet,
  DirectSecp256k1Wallet,
  makeCosmoshubPath,
} from '@cosmjs/proto-signing';
import {
  Secp256k1,
  Secp256k1Keypair,
  sha256,
} from '@cosmjs/crypto';
import { fromHex, toHex } from '@cosmjs/encoding';
import { SecureStorage } from './secureStorage';

export interface WalletAccount {
  address: string;
  publicKey: string;
  privateKey: string;
  path: string;
  coin: string;
}

export interface HDWalletConfig {
  mnemonic: string;
  password?: string;
  accounts?: number;
}

export class HDWallet {
  private mnemonic: string;
  private seed: Buffer;
  private masterKey: any; // BIP32Interface
  private accounts: Map<string, WalletAccount[]> = new Map();
  private secureStorage: SecureStorage;
  private static bip32: any;

  // Derivation paths
  private static readonly PATHS = {
    DESHCHAIN: "m/44'/118'/0'/0", // Cosmos standard
    ETHEREUM: "m/44'/60'/0'/0",
    BITCOIN: "m/44'/0'/0'/0",
  };

  constructor(private config: HDWalletConfig) {
    this.mnemonic = config.mnemonic;
    this.seed = bip39.mnemonicToSeedSync(this.mnemonic, config.password);
    if (!HDWallet.bip32) {
      HDWallet.bip32 = BIP32Factory(ecc);
    }
    this.masterKey = HDWallet.bip32.fromSeed(this.seed);
    this.secureStorage = SecureStorage.getInstance();
  }

  /**
   * Generate a new HD wallet with random mnemonic
   */
  static async generate(strength: number = 256): Promise<HDWallet> {
    const mnemonic = bip39.generateMnemonic(strength);
    return new HDWallet({ mnemonic });
  }

  /**
   * Restore HD wallet from mnemonic
   */
  static async fromMnemonic(mnemonic: string, password?: string): Promise<HDWallet> {
    if (!bip39.validateMnemonic(mnemonic)) {
      throw new Error('Invalid mnemonic phrase');
    }
    return new HDWallet({ mnemonic, password });
  }

  /**
   * Get mnemonic phrase
   */
  getMnemonic(): string {
    return this.mnemonic;
  }

  /**
   * Derive DeshChain account
   */
  async deriveDeshChainAccount(index: number = 0): Promise<WalletAccount> {
    const path = `${HDWallet.PATHS.DESHCHAIN}/${index}`;
    
    // Use CosmJS for DeshChain (Cosmos SDK based)
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
      this.mnemonic,
      {
        prefix: 'desh',
        hdPaths: [makeCosmoshubPath(index)],
      }
    );

    const [account] = await wallet.getAccounts();
    const privateKey = await this.getPrivateKeyForPath(path);

    const walletAccount: WalletAccount = {
      address: account.address,
      publicKey: toHex(account.pubkey),
      privateKey: toHex(privateKey),
      path,
      coin: 'DESHCHAIN',
    };

    this.cacheAccount('DESHCHAIN', walletAccount);
    return walletAccount;
  }

  /**
   * Derive Ethereum account
   */
  deriveEthereumAccount(index: number = 0): WalletAccount {
    const path = `${HDWallet.PATHS.ETHEREUM}/${index}`;
    const child = this.masterKey.derivePath(path);
    
    if (!child.privateKey) {
      throw new Error('Failed to derive private key');
    }

    const privateKey = child.privateKey;
    const wallet = new ethers.Wallet(privateKey);

    const walletAccount: WalletAccount = {
      address: wallet.address,
      publicKey: wallet.publicKey,
      privateKey: wallet.privateKey,
      path,
      coin: 'ETHEREUM',
    };

    this.cacheAccount('ETHEREUM', walletAccount);
    return walletAccount;
  }

  /**
   * Derive Bitcoin account (simplified for demo)
   */
  deriveBitcoinAccount(index: number = 0): WalletAccount {
    const path = `${HDWallet.PATHS.BITCOIN}/${index}`;
    const child = this.masterKey.derivePath(path);
    
    if (!child.privateKey) {
      throw new Error('Failed to derive private key');
    }

    // Simplified Bitcoin address generation
    // In production, use proper Bitcoin libraries
    const privateKey = child.privateKey;
    const publicKey = child.publicKey;

    const walletAccount: WalletAccount = {
      address: `bc1${toHex(sha256(publicKey)).substring(0, 40)}`, // Simplified
      publicKey: toHex(publicKey),
      privateKey: toHex(privateKey),
      path,
      coin: 'BITCOIN',
    };

    this.cacheAccount('BITCOIN', walletAccount);
    return walletAccount;
  }

  /**
   * Get all derived accounts
   */
  getAllAccounts(): Map<string, WalletAccount[]> {
    return this.accounts;
  }

  /**
   * Get accounts for specific coin
   */
  getAccountsForCoin(coin: string): WalletAccount[] {
    return this.accounts.get(coin) || [];
  }

  /**
   * Sign transaction for DeshChain
   */
  async signDeshChainTransaction(
    accountIndex: number,
    signDoc: any
  ): Promise<any> {
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
      this.mnemonic,
      {
        prefix: 'desh',
        hdPaths: [makeCosmoshubPath(accountIndex)],
      }
    );

    return wallet.signDirect(
      (await wallet.getAccounts())[0].address,
      signDoc
    );
  }

  /**
   * Sign Ethereum transaction
   */
  async signEthereumTransaction(
    accountIndex: number,
    transaction: ethers.TransactionRequest
  ): Promise<string> {
    const account = this.deriveEthereumAccount(accountIndex);
    const wallet = new ethers.Wallet(account.privateKey);
    return wallet.signTransaction(transaction);
  }

  /**
   * Export wallet (encrypted)
   */
  async export(password: string): Promise<string> {
    const walletData = {
      mnemonic: this.mnemonic,
      accounts: Array.from(this.accounts.entries()),
    };

    return this.secureStorage.encryptData(
      JSON.stringify(walletData),
      password
    );
  }

  /**
   * Import wallet from encrypted export
   */
  static async import(encryptedData: string, password: string): Promise<HDWallet> {
    const secureStorage = SecureStorage.getInstance();
    const decrypted = await secureStorage.decryptData(encryptedData, password);
    const walletData = JSON.parse(decrypted);

    const wallet = new HDWallet({ mnemonic: walletData.mnemonic });
    
    // Restore cached accounts
    wallet.accounts = new Map(walletData.accounts);
    
    return wallet;
  }

  /**
   * Clear sensitive data
   */
  clear(): void {
    this.mnemonic = '';
    this.seed = Buffer.alloc(0);
    this.accounts.clear();
  }

  /**
   * Private helper methods
   */
  private async getPrivateKeyForPath(path: string): Promise<Uint8Array> {
    const child = this.masterKey.derivePath(path);
    if (!child.privateKey) {
      throw new Error('Failed to derive private key');
    }
    return child.privateKey;
  }

  private cacheAccount(coin: string, account: WalletAccount): void {
    const existing = this.accounts.get(coin) || [];
    const exists = existing.some(a => a.address === account.address);
    
    if (!exists) {
      existing.push(account);
      this.accounts.set(coin, existing);
    }
  }
}