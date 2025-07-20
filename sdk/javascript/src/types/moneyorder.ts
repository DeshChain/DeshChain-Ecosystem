/**
 * Money Order DEX types for DeshChain
 */

export interface MoneyOrder {
  orderId: string
  sender: string
  recipient: string
  amount: number
  currency: string
  fees: MoneyOrderFees
  status: MoneyOrderStatus
  sourceLocation: {
    pincode: string
    district: string
    state: string
  }
  destinationLocation: {
    pincode: string
    district: string
    state: string
  }
  estimatedDelivery: string
  actualDelivery?: string
  createdAt: string
  completedAt?: string
  memo?: string
  trackingId: string
}

export interface MoneyOrderFees {
  baseFee: number
  distanceFee: number
  serviceFee: number
  urgencyFee?: number
  totalFee: number
}

export type MoneyOrderStatus = 
  | 'created'
  | 'pending'
  | 'processing'
  | 'in_transit'
  | 'delivered'
  | 'failed'
  | 'cancelled'
  | 'refunded'

export interface MoneyOrderStats {
  totalOrders: number
  totalVolume: number
  averageAmount: number
  completionRate: number
  averageDeliveryTime: number
  popularRoutes: PopularRoute[]
}

export interface PopularRoute {
  sourceState: string
  destinationState: string
  volume: number
  averageAmount: number
  deliveryTime: number
}

export interface TradeOrder {
  orderId: string
  trader: string
  pair: string
  side: 'buy' | 'sell'
  orderType: 'market' | 'limit' | 'stop'
  amount: number
  price?: number
  stopPrice?: number
  filled: number
  remaining: number
  status: OrderStatus
  createdAt: string
  updatedAt: string
  fees: TradingFees
}

export type OrderStatus = 
  | 'pending'
  | 'partial'
  | 'filled'
  | 'cancelled'
  | 'expired'
  | 'rejected'

export interface TradingFees {
  makerFee: number
  takerFee: number
  totalFee: number
  discount?: number
}

export interface TradingPair {
  symbol: string
  baseToken: string
  quoteToken: string
  status: 'active' | 'inactive' | 'delisted'
  minAmount: number
  maxAmount: number
  priceDecimals: number
  amountDecimals: number
  fees: {
    maker: number
    taker: number
  }
  restrictions?: string[]
}

export interface OrderBook {
  bids: OrderBookEntry[]
  asks: OrderBookEntry[]
  lastPrice: number
  spread: number
}

export interface OrderBookEntry {
  price: number
  amount: number
  total: number
  orders: number
}

export interface TradeHistory {
  tradeId: string
  symbol: string
  side: 'buy' | 'sell'
  amount: number
  price: number
  value: number
  fees: number
  timestamp: string
  maker: string
  taker: string
}

export interface LiquidityPool {
  poolId: string
  tokenA: string
  tokenB: string
  reserveA: number
  reserveB: number
  totalShares: number
  apy: number
  volume24h: number
  fees24h: number
  providers: number
  status: PoolStatus
}

export type PoolStatus = 'active' | 'paused' | 'migrating' | 'deprecated'

export interface DEXStats {
  totalValueLocked: number
  totalVolume24h: number
  totalTrades24h: number
  activePairs: number
  uniqueTraders24h: number
  totalFees24h: number
}

export interface PriceData {
  timestamp: string
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export interface MarketDepth {
  price: number
  amount: number
  total: number
  side: 'bid' | 'ask'
}

export interface LiquidityPosition {
  positionId: string
  provider: string
  poolId: string
  tokenA: string
  tokenB: string
  amountA: number
  amountB: number
  shares: number
  sharePercent: number
  value: number
  impermanentLoss: number
  feesEarned: number
  createdAt: string
}

export interface SwapRoute {
  inputToken: string
  outputToken: string
  path: string[]
  inputAmount: number
  outputAmount: number
  priceImpact: number
  fees: number
  slippage: number
  minOutput: number
}

export interface VolumeStats {
  timeframe: string
  totalVolume: number
  moneyOrderVolume: number
  dexVolume: number
  topPairs: Array<{
    pair: string
    volume: number
    percentage: number
  }>
  growth: {
    daily: number
    weekly: number
    monthly: number
  }
}

export interface RegionalTradingData {
  region: string
  volume24h: number
  transactions24h: number
  uniqueTraders: number
  popularPairs: string[]
  averageTradeSize: number
  growthRate: number
}

export interface FestivalTradingBonus {
  festival: string
  bonusType: 'fee_discount' | 'volume_bonus' | 'cashback'
  percentage: number
  maxBenefit: number
  eligiblePairs: string[]
  validFrom: string
  validTo: string
  conditions: string[]
  claimed: number
  totalPool: number
}

export interface TradingIncentive {
  incentiveId: string
  name: string
  description: string
  type: IncentiveType
  reward: number
  conditions: IncentiveCondition[]
  eligibleUsers: string[]
  startDate: string
  endDate: string
  claimed: number
  budget: number
}

export type IncentiveType = 
  | 'volume_milestone'
  | 'trading_streak'
  | 'first_trade'
  | 'referral'
  | 'cultural_bonus'
  | 'regional_pride'

export interface IncentiveCondition {
  type: 'min_volume' | 'min_trades' | 'time_period' | 'pair_specific' | 'region_specific'
  value: any
  description: string
}

export interface ArbitrageOpportunity {
  pair: string
  buyExchange: string
  sellExchange: string
  buyPrice: number
  sellPrice: number
  profit: number
  profitPercent: number
  volume: number
  confidence: number
  lastUpdated: string
}

export interface TradingAnalytics {
  trader: string
  totalVolume: number
  totalTrades: number
  winRate: number
  avgTradeSize: number
  totalPnL: number
  bestPair: string
  tradingStreak: number
  rank: number
  badges: string[]
  culturalBonus: number
}

export interface MoneyOrderRoute {
  routeId: string
  sourcePin: string
  destinationPin: string
  distance: number
  estimatedTime: number
  baseFee: number
  serviceFee: number
  totalFee: number
  availability: 'available' | 'limited' | 'unavailable'
  restrictions?: string[]
}