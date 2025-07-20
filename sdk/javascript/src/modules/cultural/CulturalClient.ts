import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { 
  Festival, 
  CulturalQuote, 
  CulturalEvent, 
  FestivalBonusInfo,
  CulturalCalendar,
  RegionalCelebration 
} from '../../types/cultural'

/**
 * Client for interacting with DeshChain cultural heritage features
 */
export class CulturalClient {
  constructor(
    private readonly client: StargateClient,
    private readonly tmClient: Tendermint34Client
  ) {}

  /**
   * Festival Management
   */

  /**
   * Get currently active festival
   */
  async getCurrentFestival(): Promise<Festival | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_current_festival: {} }
      )
      return response.festival
    } catch {
      return null
    }
  }

  /**
   * Get all active festivals
   */
  async getActiveFestivals(): Promise<Festival[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_active_festivals: {} }
    )
    return response.festivals || []
  }

  /**
   * Get upcoming festivals
   */
  async getUpcomingFestivals(days: number = 30): Promise<Festival[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_upcoming_festivals: { days } }
    )
    return response.festivals || []
  }

  /**
   * Get festival by ID
   */
  async getFestival(festivalId: string): Promise<Festival | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_festival: { festival_id: festivalId } }
      )
      return response.festival
    } catch {
      return null
    }
  }

  /**
   * Get festival calendar for a year
   */
  async getFestivalCalendar(year: number): Promise<CulturalCalendar> {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_calendar: { year } }
    )
    return response.calendar
  }

  /**
   * Cultural Quotes
   */

  /**
   * Get daily cultural quote
   */
  async getDailyQuote(): Promise<CulturalQuote> {
    const response = await this.client.queryContractSmart(
      '',
      { get_daily_quote: {} }
    )
    return response.quote
  }

  /**
   * Get random quote by category
   */
  async getQuoteByCategory(category: string, language: string = 'en'): Promise<CulturalQuote> {
    const response = await this.client.queryContractSmart(
      '',
      { get_quote_by_category: { category, language } }
    )
    return response.quote
  }

  /**
   * Get all quotes by author
   */
  async getQuotesByAuthor(author: string, language: string = 'en'): Promise<CulturalQuote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_quotes_by_author: { author, language } }
    )
    return response.quotes || []
  }

  /**
   * Search quotes
   */
  async searchQuotes(query: string, language: string = 'en'): Promise<CulturalQuote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { search_quotes: { query, language } }
    )
    return response.quotes || []
  }

  /**
   * Get quote categories
   */
  async getQuoteCategories(): Promise<string[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_quote_categories: {} }
    )
    return response.categories || []
  }

  /**
   * Cultural Events
   */

  /**
   * Get upcoming cultural events
   */
  async getUpcomingEvents(days: number = 30): Promise<CulturalEvent[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_upcoming_events: { days } }
    )
    return response.events || []
  }

  /**
   * Get events by region
   */
  async getEventsByRegion(state: string, district?: string): Promise<CulturalEvent[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_events_by_region: { state, district } }
    )
    return response.events || []
  }

  /**
   * Get historical events
   */
  async getHistoricalEvents(date: string): Promise<CulturalEvent[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_historical_events: { date } }
    )
    return response.events || []
  }

  /**
   * Festival Bonuses and Incentives
   */

  /**
   * Get current festival bonuses
   */
  async getCurrentFestivalBonuses(): Promise<FestivalBonusInfo[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_current_festival_bonuses: {} }
    )
    return response.bonuses || []
  }

  /**
   * Calculate festival bonus for transaction
   */
  async calculateFestivalBonus(
    amount: number,
    transactionType: string
  ): Promise<{ bonus: number; percentage: number; festival: string }> {
    const response = await this.client.queryContractSmart(
      '',
      { calculate_festival_bonus: { amount, transaction_type: transactionType } }
    )
    return response
  }

  /**
   * Get festival bonus history
   */
  async getFestivalBonusHistory(address: string, limit: number = 10) {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_bonus_history: { address, limit } }
    )
    return response.history || []
  }

  /**
   * Regional and Cultural Customization
   */

  /**
   * Get regional celebrations
   */
  async getRegionalCelebrations(state: string): Promise<RegionalCelebration[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_celebrations: { state } }
    )
    return response.celebrations || []
  }

  /**
   * Get cultural preferences by pincode
   */
  async getCulturalPreferences(pincode: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_preferences: { pincode } }
    )
    return response.preferences
  }

  /**
   * Get supported languages
   */
  async getSupportedLanguages(): Promise<string[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_supported_languages: {} }
    )
    return response.languages || []
  }

  /**
   * Cultural Statistics
   */

  /**
   * Get cultural engagement stats
   */
  async getCulturalStats() {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_stats: {} }
    )
    return response.stats
  }

  /**
   * Get festival participation metrics
   */
  async getFestivalParticipation(festivalId: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_participation: { festival_id: festivalId } }
    )
    return response.participation
  }

  /**
   * Get regional festival impact
   */
  async getRegionalFestivalImpact(state: string, festivalId: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_festival_impact: { state, festival_id: festivalId } }
    )
    return response.impact
  }

  /**
   * Cultural Heritage Preservation
   */

  /**
   * Get heritage sites
   */
  async getHeritageSites(state?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_heritage_sites: { state } }
    )
    return response.sites || []
  }

  /**
   * Get cultural artifacts
   */
  async getCulturalArtifacts(category?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_artifacts: { category } }
    )
    return response.artifacts || []
  }

  /**
   * Get traditional practices
   */
  async getTraditionalPractices(region: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_traditional_practices: { region } }
    )
    return response.practices || []
  }

  /**
   * Utility Methods
   */

  /**
   * Check if date is a festival
   */
  async isFestivalDate(date: string): Promise<boolean> {
    const response = await this.client.queryContractSmart(
      '',
      { is_festival_date: { date } }
    )
    return response.is_festival
  }

  /**
   * Get next festival
   */
  async getNextFestival(): Promise<Festival | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_next_festival: {} }
      )
      return response.festival
    } catch {
      return null
    }
  }

  /**
   * Get cultural context for transaction
   */
  async getCulturalContext(
    address: string,
    transactionType: string
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_context: { address, transaction_type: transactionType } }
    )
    return response.context
  }

  /**
   * Get localized content
   */
  async getLocalizedContent(
    contentId: string,
    language: string = 'en'
  ) {
    const response = await this.client.queryContractSmart(
      '',
      { get_localized_content: { content_id: contentId, language } }
    )
    return response.content
  }
}