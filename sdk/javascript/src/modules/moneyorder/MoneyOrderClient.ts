import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { 
  MoneyOrder, 
  MoneyOrderStatus, 
  MoneyOrderStats,
  TradeOrder,
  LiquidityPool,
  TradingPair,
  OrderBook,
  TradeHistory,
  DEXStats,
  PriceData
} from '../../types/moneyorder'

/**
 * Client for interacting with Money Order DEX
 * Traditional money order concept reimagined for blockchain
 */
export class MoneyOrderClient {
  constructor(
    private readonly client: StargateClient,
    private readonly tmClient: Tendermint34Client
  ) {}

  /**
   * Money Order Management
   */

  /**
   * Get money order by ID
   */
  async getMoneyOrder(orderId: string): Promise<MoneyOrder | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_money_order: { order_id: orderId } }
      )
      return response.money_order
    } catch {
      return null
    }
  }

  /**
   * Get money orders by sender
   */
  async getMoneyOrdersBySender(senderAddress: string): Promise<MoneyOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_money_orders_by_sender: { sender: senderAddress } }
    )
    return response.money_orders || []
  }

  /**
   * Get money orders by recipient
   */
  async getMoneyOrdersByRecipient(recipientAddress: string): Promise<MoneyOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_money_orders_by_recipient: { recipient: recipientAddress } }
    )
    return response.money_orders || []
  }

  /**
   * Get pending money orders
   */
  async getPendingMoneyOrders(address?: string): Promise<MoneyOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_pending_money_orders: { address } }
    )
    return response.money_orders || []
  }

  /**
   * Get money order history
   */
  async getMoneyOrderHistory(
    address: string,
    limit: number = 10,
    offset: number = 0
  ): Promise<MoneyOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_money_order_history: { address, limit, offset } }
    )
    return response.money_orders || []
  }

  /**
   * Calculate money order fees
   */
  async calculateMoneyOrderFees(
    amount: number,
    sourcePin: string,
    destinationPin: string
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { calculate_money_order_fees: { amount, source_pin: sourcePin, destination_pin: destinationPin } }
    )
    return {
      baseFee: response.base_fee,
      distanceFee: response.distance_fee,
      serviceFee: response.service_fee,
      totalFee: response.total_fee,
      estimatedDelivery: response.estimated_delivery,
    }
  }

  /**
   * DEX Trading Functions
   */

  /**
   * Get all trading pairs
   */
  async getTradingPairs(): Promise<TradingPair[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_trading_pairs: {} }
    )
    return response.pairs || []
  }

  /**
   * Get trading pair info
   */
  async getTradingPair(baseToken: string, quoteToken: string): Promise<TradingPair | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_trading_pair: { base_token: baseToken, quote_token: quoteToken } }
      )
      return response.pair
    } catch {
      return null
    }
  }

  /**
   * Get order book for trading pair
   */
  async getOrderBook(baseToken: string, quoteToken: string): Promise<OrderBook> {
    const response = await this.client.queryContractSmart(
      '',
      { get_order_book: { base_token: baseToken, quote_token: quoteToken } }
    )
    return {
      bids: response.bids || [],
      asks: response.asks || [],
      lastPrice: response.last_price,
      spread: response.spread,
    }
  }

  /**
   * Get trade order by ID
   */
  async getTradeOrder(orderId: string): Promise<TradeOrder | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_trade_order: { order_id: orderId } }
      )
      return response.order
    } catch {
      return null
    }
  }

  /**
   * Get trade orders by user
   */
  async getTradeOrdersByUser(userAddress: string): Promise<TradeOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_trade_orders_by_user: { user: userAddress } }
    )
    return response.orders || []
  }

  /**
   * Get open orders for user
   */
  async getOpenOrders(userAddress: string): Promise<TradeOrder[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_open_orders: { user: userAddress } }
    )
    return response.orders || []
  }

  /**
   * Get trade history for pair
   */
  async getTradeHistory(
    baseToken: string,
    quoteToken: string,
    limit: number = 50
  ): Promise<TradeHistory[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_trade_history: { base_token: baseToken, quote_token: quoteToken, limit } }
    )
    return response.trades || []
  }

  /**
   * Get user trade history
   */
  async getUserTradeHistory(
    userAddress: string,
    limit: number = 50
  ): Promise<TradeHistory[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_user_trade_history: { user: userAddress, limit } }
    )
    return response.trades || []
  }

  /**
   * Calculate trade price and fees
   */
  async calculateTrade(
    baseToken: string,
    quoteToken: string,
    side: 'buy' | 'sell',
    amount: number,
    orderType: 'market' | 'limit',
    limitPrice?: number
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { 
        calculate_trade: { 
          base_token: baseToken,
          quote_token: quoteToken,
          side,
          amount,
          order_type: orderType,
          limit_price: limitPrice
        } 
      }
    )
    return {
      executionPrice: response.execution_price,
      totalCost: response.total_cost,
      fees: response.fees,
      priceImpact: response.price_impact,
      slippage: response.slippage,
    }
  }

  /**
   * Liquidity Pool Management
   */

  /**
   * Get all liquidity pools
   */
  async getLiquidityPools(): Promise<LiquidityPool[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_liquidity_pools: {} }
    )
    return response.pools || []
  }

  /**
   * Get liquidity pool by tokens
   */
  async getLiquidityPool(tokenA: string, tokenB: string): Promise<LiquidityPool | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_liquidity_pool: { token_a: tokenA, token_b: tokenB } }
      )
      return response.pool
    } catch {
      return null
    }
  }

  /**
   * Get user liquidity positions
   */
  async getUserLiquidityPositions(userAddress: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_user_liquidity_positions: { user: userAddress } }
    )
    return response.positions || []
  }

  /**
   * Calculate add liquidity
   */
  async calculateAddLiquidity(
    tokenA: string,
    tokenB: string,
    amountA: number,
    amountB?: number
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { calculate_add_liquidity: { token_a: tokenA, token_b: tokenB, amount_a: amountA, amount_b: amountB } }
    )
    return {
      optimalAmountA: response.optimal_amount_a,
      optimalAmountB: response.optimal_amount_b,
      lpTokens: response.lp_tokens,
      shareOfPool: response.share_of_pool,
    }
  }

  /**
   * Calculate remove liquidity
   */
  async calculateRemoveLiquidity(
    tokenA: string,
    tokenB: string,
    lpTokenAmount: number
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { calculate_remove_liquidity: { token_a: tokenA, token_b: tokenB, lp_token_amount: lpTokenAmount } }
    )
    return {
      amountA: response.amount_a,
      amountB: response.amount_b,
      shareOfPool: response.share_of_pool,
    }
  }

  /**
   * Price and Market Data
   */

  /**
   * Get current price for token pair
   */
  async getCurrentPrice(baseToken: string, quoteToken: string): Promise<number> {
    const response = await this.client.queryContractSmart(
      '',
      { get_current_price: { base_token: baseToken, quote_token: quoteToken } }
    )
    return response.price
  }

  /**
   * Get price history
   */
  async getPriceHistory(
    baseToken: string,
    quoteToken: string,
    interval: '1m' | '5m' | '1h' | '4h' | '1d' = '1h',
    limit: number = 100
  ): Promise<PriceData[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_price_history: { base_token: baseToken, quote_token: quoteToken, interval, limit } }
    )
    return response.prices || []
  }

  /**
   * Get 24h ticker data
   */
  async get24hTicker(baseToken: string, quoteToken: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_24h_ticker: { base_token: baseToken, quote_token: quoteToken } }
    )
    return {
      symbol: response.symbol,
      lastPrice: response.last_price,
      priceChange: response.price_change,
      priceChangePercent: response.price_change_percent,
      high: response.high,
      low: response.low,
      volume: response.volume,
      quoteVolume: response.quote_volume,
    }
  }

  /**
   * Get market depth
   */
  async getMarketDepth(baseToken: string, quoteToken: string, limit: number = 20) {
    const response = await this.client.queryContractSmart(
      '',
      { get_market_depth: { base_token: baseToken, quote_token: quoteToken, limit } }
    )
    return {
      bids: response.bids || [],
      asks: response.asks || [],
      lastUpdateTime: response.last_update_time,
    }
  }

  /**
   * Statistics and Analytics
   */

  /**
   * Get DEX statistics
   */
  async getDEXStats(): Promise<DEXStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_dex_stats: {} }
    )
    return {
      totalValueLocked: response.total_value_locked,
      totalVolume24h: response.total_volume_24h,
      totalTrades24h: response.total_trades_24h,
      activePairs: response.active_pairs,
      uniqueTraders24h: response.unique_traders_24h,
      totalFees24h: response.total_fees_24h,
    }
  }

  /**
   * Get money order statistics
   */
  async getMoneyOrderStats(): Promise<MoneyOrderStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_money_order_stats: {} }
    )
    return {
      totalOrders: response.total_orders,
      totalVolume: response.total_volume,
      averageAmount: response.average_amount,
      completionRate: response.completion_rate,
      averageDeliveryTime: response.average_delivery_time,
      popularRoutes: response.popular_routes || [],
    }
  }

  /**
   * Get volume statistics
   */
  async getVolumeStats(timeframe: '24h' | '7d' | '30d' = '24h') {
    const response = await this.client.queryContractSmart(
      '',
      { get_volume_stats: { timeframe } }
    )
    return response.stats
  }

  /**
   * Get top trading pairs
   */
  async getTopTradingPairs(limit: number = 10) {
    const response = await this.client.queryContractSmart(
      '',
      { get_top_trading_pairs: { limit } }
    )
    return response.pairs || []
  }

  /**
   * Cultural and Regional Features
   */

  /**
   * Get regional money order volume
   */
  async getRegionalVolume(state: string, district?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_volume: { state, district } }
    )
    return response.volume
  }

  /**
   * Get festival trading bonuses
   */
  async getFestivalTradingBonuses() {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_trading_bonuses: {} }
    )
    return response.bonuses || []
  }

  /**
   * Get cultural trading incentives
   */
  async getCulturalTradingIncentives(pincode: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_trading_incentives: { pincode } }
    )
    return response.incentives
  }

  /**
   * Utility Functions
   */

  /**
   * Get supported tokens
   */
  async getSupportedTokens(): Promise<string[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_supported_tokens: {} }
    )
    return response.tokens || []
  }

  /**
   * Get trading fees structure
   */
  async getTradingFees() {
    const response = await this.client.queryContractSmart(
      '',
      { get_trading_fees: {} }
    )
    return response.fees
  }

  /**
   * Get money order routes
   */
  async getMoneyOrderRoutes(sourcePin: string, destinationPin: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_money_order_routes: { source_pin: sourcePin, destination_pin: destinationPin } }
    )
    return response.routes || []
  }

  /**
   * Check money order status
   */
  async checkMoneyOrderStatus(orderId: string): Promise<MoneyOrderStatus> {
    const response = await this.client.queryContractSmart(
      '',
      { check_money_order_status: { order_id: orderId } }
    )
    return response.status
  }

  /**
   * Health check
   */
  async healthCheck() {
    const response = await this.client.queryContractSmart(
      '',
      { health_check: {} }
    )
    return {
      healthy: response.healthy,
      dexStatus: response.dex_status,
      moneyOrderStatus: response.money_order_status,
      lastUpdated: response.last_updated,
    }
  }
}