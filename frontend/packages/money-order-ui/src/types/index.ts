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

export interface MoneyOrderFormData {
  sender: string;
  receiver: string;
  amount: {
    value: string;
    denom: string;
  };
  poolId: string;
  memo?: string;
  priority: 'standard' | 'fast' | 'instant';
  culturalPreferences: {
    language: string;
    theme: string;
    includeQuote: boolean;
  };
}

export interface PoolInfo {
  poolId: string;
  type: 'amm' | 'fixed_rate' | 'village' | 'concentrated';
  tokenA: string;
  tokenB: string;
  reserveA: {
    denom: string;
    amount: string;
  };
  reserveB: {
    denom: string;
    amount: string;
  };
  exchangeRate?: string;
  swapFee: string;
  culturalTheme: string;
  trustScore: number;
  priceStability?: number;
  villageName?: string;
  postalCode?: string;
  members?: string[];
  coordinator?: string;
  patriotismQuote?: string;
  status: 'active' | 'paused' | 'deprecated';
  apy?: number;
  volume24h?: string;
  monthlyTransactions?: number;
}

export interface CulturalQuoteData {
  quoteId: string;
  text: string;
  author: string;
  category: string;
  language: string;
  occasion: string;
  translation?: string;
  context?: string;
  weight: number;
  active: boolean;
}

export interface FestivalInfo {
  festivalId: string;
  name: string;
  description: string;
  startDate: string;
  endDate: string;
  bonusRate: number;
  culturalTheme: string;
  region: string;
  significance?: string;
  traditionalGreeting?: string;
  active: boolean;
  daysRemaining?: number;
}

export interface LanguageOption {
  code: string;
  name: string;
  nativeName: string;
  region: string;
  supported: boolean;
  script?: string;
}

export interface ThemeConfig {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  surface: string;
  text: string;
  festivalColors?: {
    [festivalName: string]: {
      primary: string;
      secondary: string;
      accent: string;
    };
  };
  culturalElements: {
    borderStyle: string;
    pattern: string;
    iconSet: string;
  };
}

export interface PatriotismScore {
  userId: string;
  score: number; // 0-100 scale
  socialImpact: number;
  communitySupport: number;
  culturalEngagement: number;
  level: 'bronze' | 'silver' | 'gold' | 'platinum';
  badges: string[];
  rewardsEarned: {
    denom: string;
    amount: string;
  };
}

export interface ReceiptData {
  receiptId: string;
  orderId: string;
  sender: string;
  receiver: string;
  amount: {
    denom: string;
    amount: string;
  };
  fee: {
    denom: string;
    amount: string;
  };
  exchangeRate: string;
  culturalQuote: CulturalQuoteData;
  patriotismBonus?: {
    denom: string;
    amount: string;
  };
  verificationCode: string;
  qrCode: string;
  timestamp: string;
  status: TransactionStatus;
  digitalSignature: string;
  blockchainConfirmation?: {
    transactionHash: string;
    blockHeight: number;
    confirmations: number;
  };
}

export type TransactionStatus = 
  | 'created' 
  | 'processing' 
  | 'completed' 
  | 'failed' 
  | 'cancelled' 
  | 'expired';

export interface MoneyOrderConfig {
  apiUrl: string;
  chainId: string;
  defaultLanguage: string;
  enableCulturalFeatures: boolean;
  enableFestivalThemes: boolean;
  enablePatriotismRewards: boolean;
  maxSlippage: number;
  defaultPoolType: string;
  autoSelectPool: boolean;
  theme: ThemeConfig;
}

export interface LiquidityPosition {
  positionId: string;
  poolId: string;
  owner: string;
  shares: string;
  tokenAAmount: {
    denom: string;
    amount: string;
  };
  tokenBAmount: {
    denom: string;
    amount: string;
  };
  rewardsAccumulated: Array<{
    denom: string;
    amount: string;
  }>;
  culturalBonus: {
    denom: string;
    amount: string;
  };
  patriotismReward: {
    denom: string;
    amount: string;
  };
  communityContribution: number;
  createdAt: string;
  lastClaimAt: string;
  status: 'active' | 'closed';
}

export interface SwapQuote {
  tokenIn: {
    denom: string;
    amount: string;
  };
  tokenOut: {
    denom: string;
    amount: string;
  };
  exchangeRate: string;
  priceImpact: number;
  fee: {
    denom: string;
    amount: string;
  };
  minimumReceived: string;
  route: Array<{
    poolId: string;
    tokenIn: string;
    tokenOut: string;
  }>;
  culturalBonus?: {
    denom: string;
    amount: string;
  };
  validUntil: string;
}

export interface VillagePoolStats {
  poolId: string;
  villageName: string;
  postalCode: string;
  memberCount: number;
  totalLiquidity: string;
  monthlyVolume: string;
  averageTransactionSize: string;
  trustScore: number;
  communityImpact: number;
  localEconomyBoost: number;
  culturalEngagement: number;
  topContributors: Array<{
    address: string;
    contribution: string;
    patriotismScore: number;
  }>;
}

export interface CommunityMetrics {
  totalPools: number;
  villagePools: number;
  totalUsers: number;
  activeUsers24h: number;
  totalVolume: string;
  culturalEngagementRate: number;
  quotesServedToday: number;
  festivalsCelebrated: number;
  averagePatriotismScore: number;
  communityPoolsGrowth: number;
  regionalBreakdown: {
    [region: string]: {
      pools: number;
      volume: string;
      users: number;
    };
  };
}

export interface ErrorInfo {
  code: string;
  message: string;
  details?: Record<string, any>;
  culturalMessage?: string;
  timestamp: string;
  requestId: string;
}

export interface ValidationError {
  field: string;
  message: string;
  code: string;
}

export interface FormState {
  isSubmitting: boolean;
  isValid: boolean;
  errors: ValidationError[];
  touched: Record<string, boolean>;
}

export interface UIPreferences {
  language: string;
  theme: string;
  festivalThemes: boolean;
  culturalQuotes: boolean;
  patriotismFeatures: boolean;
  animations: boolean;
  density: 'comfortable' | 'compact' | 'spacious';
  notifications: {
    transactions: boolean;
    festivals: boolean;
    rewards: boolean;
  };
}

export interface Analytics {
  poolPerformance: {
    poolId: string;
    volume24h: string;
    transactions: number;
    uniqueUsers: number;
    apy: number;
    priceStability: number;
    culturalEngagement: number;
  };
  userMetrics: {
    totalTransactions: number;
    totalVolume: string;
    patriotismScore: number;
    culturalEngagement: number;
    favoriteThemes: string[];
    preferredLanguages: string[];
  };
  culturalMetrics: {
    quotesDisplayed: number;
    festivalBonusesDistributed: string;
    patriotismRewards: string;
    communityImpactScore: number;
  };
}