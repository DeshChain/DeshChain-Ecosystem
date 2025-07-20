import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { 
  LaunchpadToken, 
  TokenLaunch, 
  LaunchpadStats, 
  AntiPumpConfig,
  CommunityVote,
  CreatorRewards,
  TradingMetrics,
  TokenSearchResult,
  LaunchApplication
} from '../../types/sikkebaaz'

/**
 * Client for interacting with Sikkebaaz memecoin launchpad
 * Features anti-pump & dump protection and cultural integration
 */
export class SikkebaazClient {
  constructor(
    private readonly client: StargateClient,
    private readonly tmClient: Tendermint34Client
  ) {}

  /**
   * Token Launch Management
   */

  /**
   * Get all active token launches
   */
  async getActiveTokens(): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_active_tokens: {} }
    )
    return response.tokens || []
  }

  /**
   * Get featured tokens
   */
  async getFeaturedTokens(): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_featured_tokens: {} }
    )
    return response.tokens || []
  }

  /**
   * Get token by symbol
   */
  async getToken(symbol: string): Promise<LaunchpadToken | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_token: { symbol } }
      )
      return response.token
    } catch {
      return null
    }
  }

  /**
   * Get token launch details
   */
  async getTokenLaunch(launchId: string): Promise<TokenLaunch | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_token_launch: { launch_id: launchId } }
      )
      return response.launch
    } catch {
      return null
    }
  }

  /**
   * Get upcoming token launches
   */
  async getUpcomingLaunches(): Promise<TokenLaunch[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_upcoming_launches: {} }
    )
    return response.launches || []
  }

  /**
   * Get token launches by creator
   */
  async getCreatorTokens(creatorAddress: string): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_creator_tokens: { creator: creatorAddress } }
    )
    return response.tokens || []
  }

  /**
   * Search tokens
   */
  async searchTokens(query: string): Promise<TokenSearchResult[]> {
    const response = await this.client.queryContractSmart(
      '',
      { search_tokens: { query } }
    )
    return response.results || []
  }

  /**
   * Anti-Pump & Dump Protection
   */

  /**
   * Get anti-pump configuration for token
   */
  async getAntiPumpConfig(symbol: string): Promise<AntiPumpConfig | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_anti_pump_config: { symbol } }
      )
      return response.config
    } catch {
      return null
    }
  }

  /**
   * Check if address can trade token
   */
  async canTrade(address: string, symbol: string, amount: number): Promise<{
    canTrade: boolean
    reason?: string
    delayUntil?: string
    maxAllowed?: number
  }> {
    const response = await this.client.queryContractSmart(
      '',
      { can_trade: { address, symbol, amount } }
    )
    return response
  }

  /**
   * Get trading restrictions for address
   */
  async getTradingRestrictions(address: string, symbol: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_trading_restrictions: { address, symbol } }
    )
    return response.restrictions
  }

  /**
   * Get wallet limits for token
   */
  async getWalletLimits(symbol: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_wallet_limits: { symbol } }
    )
    return response.limits
  }

  /**
   * Community Governance
   */

  /**
   * Get community votes for token
   */
  async getCommunityVotes(symbol: string): Promise<CommunityVote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_community_votes: { symbol } }
    )
    return response.votes || []
  }

  /**
   * Get veto power for address
   */
  async getVetoPower(address: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_veto_power: { address } }
    )
    return response.power
  }

  /**
   * Get token approval status
   */
  async getTokenApprovalStatus(symbol: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_approval_status: { symbol } }
    )
    return response.status
  }

  /**
   * Check if token is community approved
   */
  async isCommunityApproved(symbol: string): Promise<boolean> {
    const response = await this.client.queryContractSmart(
      '',
      { is_community_approved: { symbol } }
    )
    return response.approved
  }

  /**
   * Creator Rewards & Metrics
   */

  /**
   * Get creator rewards
   */
  async getCreatorRewards(creatorAddress: string): Promise<CreatorRewards> {
    const response = await this.client.queryContractSmart(
      '',
      { get_creator_rewards: { creator: creatorAddress } }
    )
    return response.rewards
  }

  /**
   * Get trading metrics for token
   */
  async getTradingMetrics(symbol: string): Promise<TradingMetrics> {
    const response = await this.client.queryContractSmart(
      '',
      { get_trading_metrics: { symbol } }
    )
    return response.metrics
  }

  /**
   * Get top performing tokens
   */
  async getTopPerformingTokens(timeframe: '24h' | '7d' | '30d' = '24h'): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_top_performing: { timeframe } }
    )
    return response.tokens || []
  }

  /**
   * Get creator leaderboard
   */
  async getCreatorLeaderboard(limit: number = 10) {
    const response = await this.client.queryContractSmart(
      '',
      { get_creator_leaderboard: { limit } }
    )
    return response.leaderboard || []
  }

  /**
   * Cultural Integration
   */

  /**
   * Get tokens by cultural category
   */
  async getTokensByCategory(category: string): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_tokens_by_category: { category } }
    )
    return response.tokens || []
  }

  /**
   * Get regional tokens
   */
  async getRegionalTokens(state: string, district?: string): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_tokens: { state, district } }
    )
    return response.tokens || []
  }

  /**
   * Get festival themed tokens
   */
  async getFestivalTokens(festivalId?: string): Promise<LaunchpadToken[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_tokens: { festival_id: festivalId } }
    )
    return response.tokens || []
  }

  /**
   * Get cultural bonus for token launch
   */
  async getCulturalBonus(
    pincode: string,
    category: string,
    festivalActive: boolean
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_bonus: { pincode, category, festival_active: festivalActive } }
    )
    return response.bonus
  }

  /**
   * Launchpad Statistics
   */

  /**
   * Get overall launchpad statistics
   */
  async getLaunchpadStats(): Promise<LaunchpadStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_launchpad_stats: {} }
    )
    return {
      totalTokensLaunched: response.total_tokens_launched,
      totalValueLocked: response.total_value_locked,
      activeLaunches: response.active_launches,
      successfulLaunches: response.successful_launches,
      totalVolume24h: response.total_volume_24h,
      uniqueTraders: response.unique_traders,
      averageSuccessRate: response.average_success_rate,
    }
  }

  /**
   * Get daily statistics
   */
  async getDailyStats(days: number = 30) {
    const response = await this.client.queryContractSmart(
      '',
      { get_daily_stats: { days } }
    )
    return response.stats || []
  }

  /**
   * Get regional statistics
   */
  async getRegionalStats(state?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_stats: { state } }
    )
    return response.stats
  }

  /**
   * Token Validation
   */

  /**
   * Validate token launch application
   */
  async validateLaunchApplication(application: LaunchApplication) {
    const response = await this.client.queryContractSmart(
      '',
      { validate_launch_application: { application } }
    )
    return {
      valid: response.valid,
      errors: response.errors || [],
      warnings: response.warnings || [],
      culturalScore: response.cultural_score,
    }
  }

  /**
   * Check token symbol availability
   */
  async isSymbolAvailable(symbol: string): Promise<boolean> {
    const response = await this.client.queryContractSmart(
      '',
      { is_symbol_available: { symbol } }
    )
    return response.available
  }

  /**
   * Get recommended launch parameters
   */
  async getRecommendedLaunchParams(
    category: string,
    targetAudience: string,
    pincode?: string
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { get_recommended_params: { category, target_audience: targetAudience, pincode } }
    )
    return response.params
  }

  /**
   * Liquidity Management
   */

  /**
   * Get liquidity lock info
   */
  async getLiquidityLockInfo(symbol: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_liquidity_lock_info: { symbol } }
    )
    return response.info
  }

  /**
   * Get bonding curve details
   */
  async getBondingCurveDetails(symbol: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_bonding_curve: { symbol } }
    )
    return response.curve
  }

  /**
   * Calculate token price
   */
  async calculateTokenPrice(symbol: string, amount: number) {
    const response = await this.client.queryContractSmart(
      '',
      { calculate_price: { symbol, amount } }
    )
    return {
      price: response.price,
      priceImpact: response.price_impact,
      fees: response.fees,
      slippage: response.slippage,
    }
  }

  /**
   * Utility Methods
   */

  /**
   * Get token categories
   */
  async getTokenCategories(): Promise<string[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_token_categories: {} }
    )
    return response.categories || []
  }

  /**
   * Get supported regions
   */
  async getSupportedRegions() {
    const response = await this.client.queryContractSmart(
      '',
      { get_supported_regions: {} }
    )
    return response.regions || []
  }

  /**
   * Get launchpad parameters
   */
  async getLaunchpadParams() {
    const response = await this.client.queryContractSmart(
      '',
      { get_launchpad_params: {} }
    )
    return response.params
  }

  /**
   * Health check for launchpad
   */
  async healthCheck() {
    const response = await this.client.queryContractSmart(
      '',
      { health_check: {} }
    )
    return {
      healthy: response.healthy,
      status: response.status,
      lastUpdated: response.last_updated,
    }
  }
}