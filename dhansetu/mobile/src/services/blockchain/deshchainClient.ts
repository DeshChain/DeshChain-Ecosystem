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

import { 
  StargateClient, 
  SigningStargateClient,
  defaultRegistryTypes,
  coins,
  StdFee,
  DeliverTxResponse,
} from '@cosmjs/stargate';
import { Registry } from '@cosmjs/proto-signing';
import { TxRaw } from 'cosmjs-types/cosmos/tx/v1beta1/tx';
import axios, { AxiosInstance } from 'axios';
import { HDWallet } from '../wallet/hdWallet';

// Import DeshChain message types
import { 
  MsgSend,
  MsgDelegate,
  MsgUndelegate,
  MsgWithdrawDelegatorReward,
} from './messages';

export interface DeshChainConfig {
  rpcEndpoint: string;
  apiEndpoint: string;
  chainId: string;
  addressPrefix: string;
  gasPrice: string;
  gasDenom: string;
}

export interface Balance {
  denom: string;
  amount: string;
}

export interface Transaction {
  hash: string;
  height: number;
  timestamp: string;
  type: string;
  amount?: string;
  fee?: string;
  status: 'success' | 'failed' | 'pending';
  memo?: string;
  culturalQuote?: string;
}

export interface Validator {
  operatorAddress: string;
  consensusPubkey: string;
  jailed: boolean;
  status: string;
  tokens: string;
  delegatorShares: string;
  description: {
    moniker: string;
    identity: string;
    website: string;
    details: string;
  };
  commission: {
    rate: string;
    maxRate: string;
    maxChangeRate: string;
  };
}

export class DeshChainClient {
  private static instance: DeshChainClient;
  private client?: StargateClient;
  private signingClient?: SigningStargateClient;
  private apiClient: AxiosInstance;
  private config: DeshChainConfig;
  private registry: Registry;

  constructor(config?: DeshChainConfig) {
    if (!config) {
      // Default configuration for DeshChain
      config = {
        rpcEndpoint: 'https://rpc.deshchain.com',
        apiEndpoint: 'https://api.deshchain.com',
        chainId: 'deshchain-1',
        addressPrefix: 'desh',
        gasPrice: '0.025unamo',
        gasDenom: 'unamo',
      };
    }
    this.config = config;
    this.apiClient = axios.create({
      baseURL: config.apiEndpoint,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Initialize custom registry with DeshChain types
    this.registry = new Registry(defaultRegistryTypes);
    this.registerCustomTypes();
  }

  /**
   * Get singleton instance
   */
  static async getInstance(): Promise<DeshChainClient> {
    if (!DeshChainClient.instance) {
      DeshChainClient.instance = new DeshChainClient();
      await DeshChainClient.instance.connect();
    }
    return DeshChainClient.instance;
  }

  /**
   * Register custom DeshChain message types
   */
  private registerCustomTypes(): void {
    // Register NAMO token messages
    this.registry.register('/deshchain.namo.MsgTransferNAMO', MsgSend);
    
    // Register Money Order DEX messages
    this.registry.register('/deshchain.moneyorder.MsgCreateOrder', MsgSend);
    
    // Register Sikkebaaz messages
    this.registry.register('/deshchain.sikkebaaz.MsgCreateLaunch', MsgSend);
    
    // Register DhanSetu messages
    this.registry.register('/deshchain.dhansetu.MsgCreateDhanPata', MsgSend);
    
    // Add more custom types as needed
  }

  /**
   * Connect to DeshChain
   */
  async connect(): Promise<void> {
    try {
      this.client = await StargateClient.connect(this.config.rpcEndpoint);
      console.log('Connected to DeshChain');
    } catch (error) {
      console.error('Failed to connect to DeshChain:', error);
      throw error;
    }
  }

  /**
   * Connect with signing capability
   */
  async connectWithSigner(wallet: HDWallet, accountIndex: number = 0): Promise<void> {
    try {
      const account = await wallet.deriveDeshChainAccount(accountIndex);
      
      // Create offline signer
      const offlineSigner = {
        getAccounts: async () => [{
          address: account.address,
          algo: 'secp256k1' as const,
          pubkey: Buffer.from(account.publicKey, 'hex'),
        }],
        signDirect: async (signerAddress: string, signDoc: any) => {
          return wallet.signDeshChainTransaction(accountIndex, signDoc);
        },
      };

      this.signingClient = await SigningStargateClient.connectWithSigner(
        this.config.rpcEndpoint,
        offlineSigner,
        {
          registry: this.registry,
          gasPrice: {
            amount: this.config.gasPrice,
            denom: this.config.gasDenom,
          },
        }
      );
      
      console.log('Connected to DeshChain with signer');
    } catch (error) {
      console.error('Failed to connect with signer:', error);
      throw error;
    }
  }

  /**
   * Get account balance
   */
  async getBalance(address: string, denom: string = 'namo'): Promise<Balance> {
    if (!this.client) {
      throw new Error('Client not connected');
    }

    const balance = await this.client.getBalance(address, denom);
    return balance;
  }

  /**
   * Get all balances
   */
  async getAllBalances(address: string): Promise<Balance[]> {
    if (!this.client) {
      throw new Error('Client not connected');
    }

    const balances = await this.client.getAllBalances(address);
    return balances;
  }

  /**
   * Send NAMO tokens
   */
  async sendNAMO(
    fromAddress: string,
    toAddress: string,
    amount: string,
    memo?: string,
    culturalQuote?: string
  ): Promise<DeliverTxResponse> {
    if (!this.signingClient) {
      throw new Error('Signing client not connected');
    }

    const sendMsg = {
      typeUrl: '/cosmos.bank.v1beta1.MsgSend',
      value: {
        fromAddress,
        toAddress,
        amount: coins(amount, 'namo'),
      },
    };

    // Add cultural quote to memo
    const finalMemo = culturalQuote 
      ? `${memo || ''}\n\n"${culturalQuote}"`
      : memo;

    const fee: StdFee = {
      amount: coins('2500', 'namo'),
      gas: '100000',
    };

    return this.signingClient.signAndBroadcast(
      fromAddress,
      [sendMsg],
      fee,
      finalMemo
    );
  }

  /**
   * Get transaction history
   */
  async getTransactionHistory(address: string, limit: number = 50): Promise<Transaction[]> {
    try {
      const response = await this.apiClient.get(`/cosmos/tx/v1beta1/txs`, {
        params: {
          'events': `message.sender='${address}'`,
          'limit': limit,
          'order_by': 'ORDER_BY_DESC',
        },
      });

      const txs = response.data.tx_responses || [];
      
      return txs.map((tx: any) => this.parseTransaction(tx));
    } catch (error) {
      console.error('Failed to fetch transaction history:', error);
      return [];
    }
  }

  /**
   * Get validators
   */
  async getValidators(status?: string): Promise<Validator[]> {
    try {
      const response = await this.apiClient.get('/cosmos/staking/v1beta1/validators', {
        params: status ? { status } : {},
      });

      return response.data.validators || [];
    } catch (error) {
      console.error('Failed to fetch validators:', error);
      return [];
    }
  }

  /**
   * Delegate to validator
   */
  async delegate(
    delegatorAddress: string,
    validatorAddress: string,
    amount: string
  ): Promise<DeliverTxResponse> {
    if (!this.signingClient) {
      throw new Error('Signing client not connected');
    }

    const delegateMsg = {
      typeUrl: '/cosmos.staking.v1beta1.MsgDelegate',
      value: {
        delegatorAddress,
        validatorAddress,
        amount: {
          denom: 'namo',
          amount,
        },
      },
    };

    const fee: StdFee = {
      amount: coins('5000', 'namo'),
      gas: '200000',
    };

    return this.signingClient.signAndBroadcast(
      delegatorAddress,
      [delegateMsg],
      fee,
      'Staking NAMO tokens'
    );
  }

  /**
   * Get delegation info
   */
  async getDelegations(delegatorAddress: string): Promise<any[]> {
    try {
      const response = await this.apiClient.get(
        `/cosmos/staking/v1beta1/delegations/${delegatorAddress}`
      );
      return response.data.delegation_responses || [];
    } catch (error) {
      console.error('Failed to fetch delegations:', error);
      return [];
    }
  }

  /**
   * Get staking rewards
   */
  async getStakingRewards(delegatorAddress: string): Promise<any> {
    try {
      const response = await this.apiClient.get(
        `/cosmos/distribution/v1beta1/delegators/${delegatorAddress}/rewards`
      );
      return response.data;
    } catch (error) {
      console.error('Failed to fetch staking rewards:', error);
      return { rewards: [], total: [] };
    }
  }

  /**
   * Claim staking rewards
   */
  async claimRewards(
    delegatorAddress: string,
    validatorAddress: string
  ): Promise<DeliverTxResponse> {
    if (!this.signingClient) {
      throw new Error('Signing client not connected');
    }

    const withdrawMsg = {
      typeUrl: '/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward',
      value: {
        delegatorAddress,
        validatorAddress,
      },
    };

    const fee: StdFee = {
      amount: coins('3000', 'namo'),
      gas: '150000',
    };

    return this.signingClient.signAndBroadcast(
      delegatorAddress,
      [withdrawMsg],
      fee,
      'Claiming staking rewards'
    );
  }

  /**
   * Get current block height
   */
  async getCurrentHeight(): Promise<number> {
    if (!this.client) {
      throw new Error('Client not connected');
    }
    return this.client.getHeight();
  }

  /**
   * Get chain ID
   */
  getChainId(): string {
    return this.config.chainId;
  }

  /**
   * Disconnect client
   */
  disconnect(): void {
    if (this.client) {
      this.client.disconnect();
    }
    this.client = undefined;
    this.signingClient = undefined;
  }

  /**
   * Parse raw transaction
   */
  private parseTransaction(tx: any): Transaction {
    const msg = tx.tx?.body?.messages?.[0];
    const timestamp = new Date(tx.timestamp).toISOString();
    
    let type = 'Unknown';
    let amount = '0';
    
    if (msg?.['@type']?.includes('MsgSend')) {
      type = 'Send';
      amount = msg.amount?.[0]?.amount || '0';
    } else if (msg?.['@type']?.includes('MsgDelegate')) {
      type = 'Delegate';
      amount = msg.amount?.amount || '0';
    }

    // Extract cultural quote from memo
    const memo = tx.tx?.body?.memo || '';
    const quoteMatch = memo.match(/"([^"]+)"/);
    const culturalQuote = quoteMatch ? quoteMatch[1] : undefined;

    return {
      hash: tx.txhash,
      height: parseInt(tx.height),
      timestamp,
      type,
      amount,
      fee: tx.tx?.auth_info?.fee?.amount?.[0]?.amount || '0',
      status: tx.code === 0 ? 'success' : 'failed',
      memo: memo.replace(/\n\n"[^"]+"/, '').trim(),
      culturalQuote,
    };
  }

  /**
   * Estimate gas for transaction
   */
  async estimateGas(messages: any[], memo?: string): Promise<number> {
    // Simplified gas estimation
    const baseGas = 80000;
    const gasPerMessage = 20000;
    const gasPerByte = 10;
    
    const memoGas = memo ? memo.length * gasPerByte : 0;
    const totalGas = baseGas + (messages.length * gasPerMessage) + memoGas;
    
    return Math.ceil(totalGas * 1.3); // Add 30% buffer
  }
}