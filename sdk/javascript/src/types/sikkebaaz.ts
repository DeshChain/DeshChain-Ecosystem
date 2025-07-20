/**
 * Sikkebaaz memecoin launchpad types for DeshChain
 */

export interface LaunchpadToken {
  symbol: string
  name: string
  description: string
  creator: string
  totalSupply: number
  currentSupply: number
  launchDate: string
  status: TokenStatus
  category: TokenCategory
  logoUrl: string
  websiteUrl?: string
  socialLinks: SocialLinks
  culturalTheme?: string
  region?: string
  antiPumpConfig: AntiPumpConfig
  bondingCurve: BondingCurveConfig
  tradingMetrics: TradingMetrics
  communityScore: number
}

export type TokenStatus = 
  | 'pending'
  | 'review'
  | 'approved'
  | 'launched'
  | 'trading'
  | 'graduated'
  | 'failed'
  | 'vetoed'

export type TokenCategory = 
  | 'cultural'
  | 'regional'
  | 'festival'
  | 'meme'
  | 'utility'
  | 'charity'
  | 'gaming'
  | 'art'
  | 'music'
  | 'sports'

export interface SocialLinks {
  twitter?: string
  telegram?: string
  discord?: string
  website?: string
  instagram?: string
  youtube?: string
}

export interface TokenLaunch {
  launchId: string
  symbol: string
  launchDate: string
  initialPrice: number
  targetRaise: number
  currentRaise: number
  participants: number
  launchDuration: number
  status: LaunchStatus
  milestones: LaunchMilestone[]
  terms: LaunchTerms
}

export type LaunchStatus = 
  | 'upcoming'
  | 'active'
  | 'successful'
  | 'failed'
  | 'cancelled'

export interface LaunchMilestone {
  name: string
  targetAmount: number
  currentAmount: number
  completed: boolean
  reward?: string
}

export interface LaunchTerms {
  minInvestment: number
  maxInvestment: number
  lockupPeriod: number
  vestingSchedule: VestingSchedule[]
  refundPolicy: string
}

export interface VestingSchedule {
  percentage: number
  releaseDate: string
  condition?: string
}

export interface AntiPumpConfig {
  maxWalletPercent: number
  maxTransactionPercent: number
  tradingDelay: number
  liquidityLockDuration: number
  priceImpactLimit: number
  volumeLimit24h: number
  honeypotProtection: boolean
  rugPullPrevention: boolean
  restrictions: TradingRestriction[]
}

export interface TradingRestriction {
  restrictionType: RestrictionType
  threshold: number
  duration: number
  penalty?: string
}

export type RestrictionType = 
  | 'whale_protection'
  | 'bot_prevention'
  | 'pump_prevention'
  | 'dump_prevention'
  | 'flash_loan_protection'

export interface BondingCurveConfig {
  curveType: 'linear' | 'exponential' | 'logarithmic' | 'sigmoid'
  parameters: {
    slope: number
    intercept: number
    maxPrice: number
    reserveRatio: number
  }
  liquidityTarget: number
  graduationThreshold: number
}

export interface TradingMetrics {
  price: number
  priceChange24h: number
  volume24h: number
  marketCap: number
  holders: number
  transactions24h: number
  liquidityUSD: number
  fdv: number
  ath: number
  atl: number
}

export interface CommunityVote {
  voteId: string
  symbol: string
  proposalType: VoteType
  description: string
  proposer: string
  votingPower: number
  votes: {
    yes: number
    no: number
    abstain: number
  }
  status: VoteStatus
  startTime: string
  endTime: string
  quorum: number
  threshold: number
}

export type VoteType = 
  | 'approval'
  | 'veto'
  | 'parameter_change'
  | 'upgrade'
  | 'emergency_stop'

export type VoteStatus = 
  | 'active'
  | 'passed'
  | 'rejected'
  | 'expired'
  | 'executed'

export interface CreatorRewards {
  creator: string
  totalTokensCreated: number
  successfulLaunches: number
  totalRaised: number
  reputationScore: number
  badges: CreatorBadge[]
  earnings: {
    launchFees: number
    tradingFees: number
    bonuses: number
    total: number
  }
  culturalImpact: number
}

export interface CreatorBadge {
  badgeId: string
  name: string
  description: string
  icon: string
  earnedDate: string
  rarity: 'common' | 'rare' | 'epic' | 'legendary'
}

export interface LaunchpadStats {
  totalTokensLaunched: number
  totalValueLocked: number
  activeLaunches: number
  successfulLaunches: number
  totalVolume24h: number
  uniqueTraders: number
  averageSuccessRate: number
}

export interface TokenSearchResult {
  symbol: string
  name: string
  creator: string
  marketCap: number
  volume24h: number
  priceChange24h: number
  holders: number
  category: TokenCategory
  culturalScore: number
  launchDate: string
}

export interface LaunchApplication {
  tokenInfo: {
    symbol: string
    name: string
    description: string
    totalSupply: number
    logoUrl: string
    websiteUrl?: string
    socialLinks: SocialLinks
  }
  creator: {
    address: string
    name: string
    experience: string
    previousLaunches: number
  }
  culturalInfo: {
    theme: string
    region: string
    category: TokenCategory
    significance: string
    communityBenefit: string
  }
  technical: {
    antiPumpEnabled: boolean
    liquidityLockDuration: number
    maxWalletPercent: number
    maxTxPercent: number
  }
  marketing: {
    strategy: string
    budget: number
    channels: string[]
    targetAudience: string
  }
}

export interface CulturalBonus {
  bonusType: CulturalBonusType
  percentage: number
  maxAmount: number
  conditions: string[]
  validUntil: string
  region?: string
  festival?: string
}

export type CulturalBonusType = 
  | 'regional_pride'
  | 'festival_special'
  | 'cultural_preservation'
  | 'community_benefit'
  | 'educational_value'

export interface TokenMetrics {
  symbol: string
  price: number
  priceHistory: PricePoint[]
  volume: VolumeData[]
  holders: HolderDistribution[]
  transactions: TransactionData[]
  liquidity: LiquidityData
  sentiment: SentimentData
}

export interface PricePoint {
  timestamp: string
  price: number
  volume: number
}

export interface VolumeData {
  timestamp: string
  volume: number
  transactions: number
}

export interface HolderDistribution {
  range: string
  count: number
  percentage: number
}

export interface TransactionData {
  timestamp: string
  type: 'buy' | 'sell'
  amount: number
  price: number
  trader: string
}

export interface LiquidityData {
  totalUSD: number
  tokenReserve: number
  namoReserve: number
  lpTokens: number
  apr: number
}

export interface SentimentData {
  bullish: number
  bearish: number
  neutral: number
  socialMentions: number
  communityGrowth: number
}