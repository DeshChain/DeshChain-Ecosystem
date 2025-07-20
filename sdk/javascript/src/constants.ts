/**
 * Constants for DeshChain SDK
 */

// Chain Information
export const CHAIN_IDS = {
  MAINNET: 'deshchain-1',
  TESTNET: 'deshchain-testnet-1',
  DEVNET: 'deshchain-devnet-1',
} as const

// Network Endpoints
export const ENDPOINTS = {
  [CHAIN_IDS.MAINNET]: {
    rpc: 'https://rpc.deshchain.network',
    rest: 'https://api.deshchain.network',
    explorer: 'https://explorer.deshchain.network',
  },
  [CHAIN_IDS.TESTNET]: {
    rpc: 'https://testnet-rpc.deshchain.network',
    rest: 'https://testnet-api.deshchain.network',
    explorer: 'https://testnet-explorer.deshchain.network',
  },
  [CHAIN_IDS.DEVNET]: {
    rpc: 'http://localhost:26657',
    rest: 'http://localhost:1317',
    explorer: 'http://localhost:3000',
  },
} as const

// Denominations
export const DENOMINATIONS = {
  NAMO: 'unamo',
  MICRO_NAMO: 'unamo',
} as const

// Gas Configuration
export const GAS = {
  DEFAULT_GAS_PRICE: '0.025unamo',
  DEFAULT_GAS_LIMIT: 200000,
  MAX_GAS_LIMIT: 10000000,
  MIN_GAS_PRICE: '0.001unamo',
  MAX_GAS_PRICE: '1.0unamo',
} as const

// Module Names
export const MODULES = {
  BANK: 'bank',
  STAKING: 'staking',
  GOVERNANCE: 'gov',
  DISTRIBUTION: 'distribution',
  SLASHING: 'slashing',
  NAMO: 'namo',
  CULTURAL: 'cultural',
  LENDING: 'lending',
  KRISHI_MITRA: 'krishimitra',
  VYAVASAYA_MITRA: 'vyavasayamitra',
  SHIKSHA_MITRA: 'shikshamitra',
  SIKKEBAAZ: 'sikkebaaz',
  MONEY_ORDER: 'moneyorder',
} as const

// Transaction Types
export const TX_TYPES = {
  SEND: '/cosmos.bank.v1beta1.MsgSend',
  DELEGATE: '/cosmos.staking.v1beta1.MsgDelegate',
  UNDELEGATE: '/cosmos.staking.v1beta1.MsgUndelegate',
  REDELEGATE: '/cosmos.staking.v1beta1.MsgBeginRedelegate',
  WITHDRAW_REWARDS: '/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward',
  VOTE: '/cosmos.gov.v1beta1.MsgVote',
  SUBMIT_PROPOSAL: '/cosmos.gov.v1beta1.MsgSubmitProposal',
  DEPOSIT: '/cosmos.gov.v1beta1.MsgDeposit',
  
  // DeshChain specific
  BURN_NAMO: '/deshchain.namo.v1.MsgBurnNAMO',
  VEST_NAMO: '/deshchain.namo.v1.MsgVestNAMO',
  APPLY_LOAN: '/deshchain.krishimitra.v1.MsgApplyLoan',
  LAUNCH_TOKEN: '/deshchain.sikkebaaz.v1.MsgLaunchToken',
  CREATE_MONEY_ORDER: '/deshchain.moneyorder.v1.MsgCreateMoneyOrder',
} as const

// API Endpoints
export const API_ENDPOINTS = {
  BLOCKS: '/cosmos/base/tendermint/v1beta1/blocks',
  BLOCK_BY_HEIGHT: (height: number) => `/cosmos/base/tendermint/v1beta1/blocks/${height}`,
  VALIDATORS: '/cosmos/staking/v1beta1/validators',
  DELEGATIONS: (address: string) => `/cosmos/staking/v1beta1/delegations/${address}`,
  BALANCES: (address: string) => `/cosmos/bank/v1beta1/balances/${address}`,
  PROPOSALS: '/cosmos/gov/v1beta1/proposals',
  PROPOSAL: (id: string) => `/cosmos/gov/v1beta1/proposals/${id}`,
  VOTES: (id: string) => `/cosmos/gov/v1beta1/proposals/${id}/votes`,
  
  // DeshChain specific
  CULTURAL_FESTIVALS: '/deshchain/cultural/v1/festivals',
  CULTURAL_QUOTES: '/deshchain/cultural/v1/quotes',
  LENDING_STATS: '/deshchain/lending/v1/stats',
  SIKKEBAAZ_TOKENS: '/deshchain/sikkebaaz/v1/tokens',
  MONEY_ORDER_STATS: '/deshchain/moneyorder/v1/stats',
} as const

// Error Codes
export const ERROR_CODES = {
  NETWORK_ERROR: 'NETWORK_ERROR',
  TRANSACTION_ERROR: 'TRANSACTION_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  INSUFFICIENT_FUNDS: 'INSUFFICIENT_FUNDS',
  INVALID_ADDRESS: 'INVALID_ADDRESS',
  INVALID_AMOUNT: 'INVALID_AMOUNT',
  GAS_ESTIMATION_FAILED: 'GAS_ESTIMATION_FAILED',
  BROADCAST_FAILED: 'BROADCAST_FAILED',
  TIMEOUT: 'TIMEOUT',
} as const

// Cultural Constants
export const CULTURAL = {
  SUPPORTED_LANGUAGES: [
    'en', 'hi', 'bn', 'te', 'ta', 'mr', 'gu', 'kn', 'ml', 'or',
    'pa', 'as', 'ur', 'ne', 'sd', 'kok', 'brx', 'doi', 'mni', 'sat',
    'bho', 'mai'
  ],
  
  FESTIVAL_CATEGORIES: [
    'religious', 'harvest', 'seasonal', 'national', 'regional', 'cultural', 'historical'
  ],
  
  QUOTE_CATEGORIES: [
    'wisdom', 'patriotism', 'philosophy', 'spirituality', 'motivation', 
    'peace', 'unity', 'progress', 'culture', 'education'
  ],
  
  INDIAN_STATES: [
    'Andhra Pradesh', 'Arunachal Pradesh', 'Assam', 'Bihar', 'Chhattisgarh',
    'Goa', 'Gujarat', 'Haryana', 'Himachal Pradesh', 'Jammu and Kashmir',
    'Jharkhand', 'Karnataka', 'Kerala', 'Madhya Pradesh', 'Maharashtra',
    'Manipur', 'Meghalaya', 'Mizoram', 'Nagaland', 'Odisha', 'Punjab',
    'Rajasthan', 'Sikkim', 'Tamil Nadu', 'Telangana', 'Tripura',
    'Uttar Pradesh', 'Uttarakhand', 'West Bengal', 'Delhi'
  ],
} as const

// Lending Constants
export const LENDING = {
  INTEREST_RATES: {
    KRISHI_MITRA: { MIN: 6, MAX: 9 }, // Agriculture: 6-9%
    VYAVASAYA_MITRA: { MIN: 8, MAX: 12 }, // Business: 8-12%
    SHIKSHA_MITRA: { MIN: 4, MAX: 7 }, // Education: 4-7%
  },
  
  LOAN_STATUSES: [
    'pending', 'approved', 'disbursed', 'active', 'completed', 'defaulted', 'rejected'
  ],
  
  KYC_STATUSES: ['pending', 'verified', 'rejected', 'expired'],
  
  CREDIT_SCORE_RANGES: {
    EXCELLENT: { MIN: 750, MAX: 900 },
    GOOD: { MIN: 650, MAX: 749 },
    FAIR: { MIN: 550, MAX: 649 },
    POOR: { MIN: 300, MAX: 549 },
  },
} as const

// Sikkebaaz Constants
export const SIKKEBAAZ = {
  TOKEN_CATEGORIES: [
    'cultural', 'regional', 'festival', 'meme', 'utility', 
    'charity', 'gaming', 'art', 'music', 'sports'
  ],
  
  TOKEN_STATUSES: [
    'pending', 'review', 'approved', 'launched', 'trading', 'graduated', 'failed', 'vetoed'
  ],
  
  ANTI_PUMP_LIMITS: {
    MAX_WALLET_PERCENT: 5, // 5% max wallet size
    MAX_TRANSACTION_PERCENT: 1, // 1% max transaction size
    TRADING_DELAY: 60, // 60 seconds between trades
    LIQUIDITY_LOCK_DURATION: 365, // 365 days liquidity lock
  },
} as const

// Money Order Constants
export const MONEY_ORDER = {
  STATUSES: [
    'created', 'pending', 'processing', 'in_transit', 'delivered', 'failed', 'cancelled', 'refunded'
  ],
  
  ORDER_TYPES: ['market', 'limit', 'stop'],
  
  TRADING_SIDES: ['buy', 'sell'],
  
  FEE_STRUCTURE: {
    BASE_FEE: 0.1, // 0.1% base fee
    DISTANCE_FEE_PER_KM: 0.001, // 0.001% per km
    URGENCY_MULTIPLIER: 1.5, // 1.5x for urgent delivery
  },
} as const

// Governance Constants
export const GOVERNANCE = {
  PROPOSAL_TYPES: [
    'text', 'parameter_change', 'software_upgrade', 'community_pool_spend',
    'cancel_software_upgrade', 'founder_protection', 'emergency'
  ],
  
  PROPOSAL_STATUSES: [
    'deposit_period', 'voting_period', 'passed', 'rejected', 'failed', 'invalid'
  ],
  
  VOTE_OPTIONS: ['yes', 'no', 'abstain', 'no_with_veto'],
  
  VETO_TYPES: [
    'founder_protection', 'inheritance_protection', 'revenue_protection', 'emergency_veto'
  ],
  
  VALIDATOR_STATUSES: ['bonded', 'unbonded', 'unbonding'],
} as const

// Time Constants
export const TIME = {
  BLOCK_TIME: 6, // 6 seconds average block time
  UNBONDING_PERIOD: 21 * 24 * 60 * 60 * 1000, // 21 days in milliseconds
  VOTING_PERIOD: 14 * 24 * 60 * 60 * 1000, // 14 days in milliseconds
  DEPOSIT_PERIOD: 7 * 24 * 60 * 60 * 1000, // 7 days in milliseconds
} as const

// Precision Constants
export const PRECISION = {
  NAMO_DECIMALS: 6,
  PERCENTAGE_DECIMALS: 2,
  PRICE_DECIMALS: 6,
  AMOUNT_DECIMALS: 6,
} as const

// Limits
export const LIMITS = {
  MAX_MEMO_LENGTH: 512,
  MAX_VALIDATORS_PER_DELEGATOR: 100,
  MAX_PROPOSAL_TITLE_LENGTH: 140,
  MAX_PROPOSAL_DESCRIPTION_LENGTH: 10000,
  MIN_DEPOSIT_AMOUNT: 1000000, // 1 NAMO in micro units
  MAX_DEPOSIT_AMOUNT: 1000000000000, // 1M NAMO in micro units
} as const

// Revenue Distribution (Platform Model)
export const REVENUE_DISTRIBUTION = {
  DEVELOPMENT: 0.30, // 30% to development fund
  COMMUNITY: 0.25, // 25% to community rewards
  LIQUIDITY: 0.20, // 20% to liquidity provision
  NGO: 0.10, // 10% to NGO donations
  EMERGENCY: 0.10, // 10% to emergency fund
  FOUNDER: 0.05, // 5% to founder (reduced from 20%)
} as const

// Cultural Impact Scoring
export const CULTURAL_SCORING = {
  WEIGHTS: {
    FESTIVAL_IMPORTANCE: 0.3,
    REGIONAL_PARTICIPATION: 0.3,
    TRADITIONAL_VALUE: 0.2,
    MODERN_RELEVANCE: 0.2,
  },
  
  BONUS_MULTIPLIERS: {
    NATIONAL_FESTIVAL: 1.5,
    REGIONAL_FESTIVAL: 1.3,
    LOCAL_CELEBRATION: 1.1,
    CULTURAL_EVENT: 1.2,
  },
} as const

// Export all constants as default
export default {
  CHAIN_IDS,
  ENDPOINTS,
  DENOMINATIONS,
  GAS,
  MODULES,
  TX_TYPES,
  API_ENDPOINTS,
  ERROR_CODES,
  CULTURAL,
  LENDING,
  SIKKEBAAZ,
  MONEY_ORDER,
  GOVERNANCE,
  TIME,
  PRECISION,
  LIMITS,
  REVENUE_DISTRIBUTION,
  CULTURAL_SCORING,
}