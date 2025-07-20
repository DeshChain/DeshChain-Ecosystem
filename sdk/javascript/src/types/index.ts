/**
 * Core DeshChain SDK types
 */

export interface DeshChainClientOptions {
  chainId: string
  prefix: string
  gasPrice: string
  rpcUrl?: string
  restUrl?: string
}

export interface ChainInfo {
  chainId: string
  nodeVersion: string
  blockHeight: number
  blockTime: Date
  validatorCount: number
  catchingUp: boolean
}

export interface NetworkStatus extends ChainInfo {
  activeValidators: number
  totalVotingPower: number
  tps: number
  networkHealth: 'healthy' | 'syncing' | 'error'
  culturalEvents: any[]
}

export interface Account {
  address: string
  accountNumber: number
  sequence: number
  pubKey?: any
}

export interface Balance {
  denom: string
  amount: string
}

export interface Transaction {
  txhash: string
  height: number
  timestamp: string
  fee: {
    amount: Balance[]
    gas: string
  }
  memo: string
  messages: any[]
  events: any[]
  logs: any[]
  gasUsed: string
  gasWanted: string
  success: boolean
}

export interface Block {
  header: {
    chainId: string
    height: number
    time: Date
    proposer: string
    hash: string
  }
  data: {
    txs: string[]
  }
  evidence: any[]
  lastCommit: any
}

export interface Validator {
  operatorAddress: string
  consensusPubkey: string
  jailed: boolean
  status: string
  tokens: string
  delegatorShares: string
  description: {
    moniker: string
    identity: string
    website: string
    details: string
  }
  unbondingHeight: number
  unbondingTime: Date
  commission: {
    rate: string
    maxRate: string
    maxChangeRate: string
    updateTime: Date
  }
  minSelfDelegation: string
}

export interface QueryOptions {
  height?: number
  prove?: boolean
}

export interface PaginationOptions {
  limit?: number
  offset?: number
  key?: string
  reverse?: boolean
}

export interface SearchResult {
  transactions: Transaction[]
  blocks: Block[]
  addresses: any[]
  loans: any[]
  tokens: any[]
}

// Error types
export class DeshChainError extends Error {
  constructor(message: string, public code?: string) {
    super(message)
    this.name = 'DeshChainError'
  }
}

export class NetworkError extends DeshChainError {
  constructor(message: string) {
    super(message, 'NETWORK_ERROR')
    this.name = 'NetworkError'
  }
}

export class TransactionError extends DeshChainError {
  constructor(message: string, public txHash?: string) {
    super(message, 'TRANSACTION_ERROR')
    this.name = 'TransactionError'
  }
}

export class ValidationError extends DeshChainError {
  constructor(message: string) {
    super(message, 'VALIDATION_ERROR')
    this.name = 'ValidationError'
  }
}

// Constants
export const CHAIN_IDS = {
  MAINNET: 'deshchain-1',
  TESTNET: 'deshchain-testnet-1',
  DEVNET: 'deshchain-devnet-1',
} as const

export const DENOMINATIONS = {
  NAMO: 'unamo',
  MICRO_NAMO: 'unamo',
} as const

export const DEFAULT_GAS_PRICE = '0.025unamo'
export const DEFAULT_GAS_LIMIT = 200000

// Network endpoints
export const MAINNET_RPC = 'https://rpc.deshchain.network'
export const MAINNET_REST = 'https://api.deshchain.network'
export const TESTNET_RPC = 'https://testnet-rpc.deshchain.network'
export const TESTNET_REST = 'https://testnet-api.deshchain.network'

export const ENDPOINTS = {
  [CHAIN_IDS.MAINNET]: {
    rpc: MAINNET_RPC,
    rest: MAINNET_REST,
  },
  [CHAIN_IDS.TESTNET]: {
    rpc: TESTNET_RPC,
    rest: TESTNET_REST,
  },
} as const