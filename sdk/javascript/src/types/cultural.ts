/**
 * Cultural heritage and festival types for DeshChain
 */

export interface Festival {
  festivalId: string
  name: string
  nameHindi: string
  description: string
  startDate: string
  endDate: string
  isActive: boolean
  category: FestivalCategory
  regions: string[]
  significance: string
  traditions: string[]
  bonusPercentage: number
  culturalImpact: number
  historicalImportance: string
}

export type FestivalCategory = 
  | 'religious'
  | 'harvest'
  | 'seasonal'
  | 'national'
  | 'regional'
  | 'cultural'
  | 'historical'

export interface CulturalQuote {
  quoteId: string
  text: string
  textHindi?: string
  author: string
  authorHindi?: string
  category: QuoteCategory
  language: string
  source: string
  significance: string
  popularity: number
  tags: string[]
  region?: string
  era?: string
}

export type QuoteCategory = 
  | 'wisdom'
  | 'patriotism'
  | 'philosophy'
  | 'spirituality'
  | 'motivation'
  | 'peace'
  | 'unity'
  | 'progress'
  | 'culture'
  | 'education'

export interface CulturalEvent {
  eventId: string
  name: string
  date: string
  type: EventType
  location: {
    state: string
    district?: string
    city?: string
  }
  description: string
  significance: string
  participants?: number
  duration: number
  status: 'upcoming' | 'ongoing' | 'completed'
}

export type EventType = 
  | 'festival'
  | 'ceremony'
  | 'celebration'
  | 'commemoration'
  | 'cultural_program'
  | 'historical_event'

export interface FestivalBonusInfo {
  festival: string
  bonusType: BonusType
  percentage: number
  maxAmount: number
  validFrom: string
  validTo: string
  eligibleTransactions: string[]
  regions: string[]
  conditions: string[]
}

export type BonusType = 
  | 'transaction_fee_discount'
  | 'lending_rate_discount'
  | 'trading_bonus'
  | 'staking_bonus'
  | 'cultural_reward'

export interface CulturalCalendar {
  year: number
  months: MonthData[]
  majorFestivals: Festival[]
  regionalEvents: CulturalEvent[]
  nationalHolidays: Holiday[]
}

export interface MonthData {
  month: number
  monthName: string
  monthNameHindi: string
  festivals: Festival[]
  events: CulturalEvent[]
  lunarEvents: LunarEvent[]
  seasonalInfo: SeasonalInfo
}

export interface Holiday {
  date: string
  name: string
  nameHindi: string
  type: 'national' | 'religious' | 'cultural'
  isPublicHoliday: boolean
  significance: string
}

export interface LunarEvent {
  date: string
  eventType: 'new_moon' | 'full_moon' | 'eclipse'
  name: string
  nameHindi: string
  significance: string
}

export interface SeasonalInfo {
  season: 'spring' | 'summer' | 'monsoon' | 'autumn' | 'winter'
  seasonHindi: string
  characteristics: string[]
  traditionalActivities: string[]
  crops: string[]
}

export interface RegionalCelebration {
  celebrationId: string
  name: string
  state: string
  district?: string
  description: string
  timeOfYear: string
  duration: number
  traditions: string[]
  significance: string
  participants: number
  economicImpact: number
}

export interface CulturalPreferences {
  pincode: string
  state: string
  district: string
  primaryLanguage: string
  secondaryLanguages: string[]
  majorFestivals: string[]
  culturalTraditions: string[]
  regionalCuisine: string[]
  artForms: string[]
  musicStyles: string[]
  danceStyles: string[]
}

export interface CulturalStats {
  totalQuotes: number
  totalFestivals: number
  activeFestivals: number
  supportedLanguages: number
  culturalEvents: number
  userEngagement: {
    dailyQuoteViews: number
    festivalParticipation: number
    culturalBonusClaimed: number
    regionalActiveUsers: Record<string, number>
  }
  popularContent: {
    topQuotes: CulturalQuote[]
    topFestivals: Festival[]
    trendingEvents: CulturalEvent[]
  }
}

export interface FestivalParticipation {
  festivalId: string
  participants: number
  transactions: number
  totalVolume: number
  bonusDistributed: number
  popularActivities: string[]
  regionalBreakdown: Record<string, number>
  engagementScore: number
}

export interface CulturalArtifact {
  artifactId: string
  name: string
  nameHindi: string
  category: ArtifactCategory
  origin: {
    state: string
    region: string
    period: string
  }
  description: string
  significance: string
  materials: string[]
  techniques: string[]
  images: string[]
  conservationStatus: string
  currentLocation: string
}

export type ArtifactCategory = 
  | 'sculpture'
  | 'painting'
  | 'textile'
  | 'jewelry'
  | 'pottery'
  | 'manuscript'
  | 'architecture'
  | 'musical_instrument'
  | 'tool'
  | 'weapon'

export interface HeritageSite {
  siteId: string
  name: string
  nameHindi: string
  type: SiteType
  location: {
    state: string
    district: string
    coordinates: {
      latitude: number
      longitude: number
    }
  }
  period: string
  significance: string
  description: string
  features: string[]
  conservationStatus: string
  accessibility: string
  visitingHours: string
  nearbyAttractions: string[]
}

export type SiteType = 
  | 'temple'
  | 'fort'
  | 'palace'
  | 'monument'
  | 'archaeological'
  | 'natural'
  | 'museum'
  | 'library'

export interface TraditionalPractice {
  practiceId: string
  name: string
  nameHindi: string
  region: string
  category: PracticeCategory
  description: string
  procedure: string[]
  significance: string
  materials: string[]
  occasions: string[]
  practitioners: string
  transmissionMethod: string
  currentStatus: 'active' | 'endangered' | 'extinct'
}

export type PracticeCategory = 
  | 'ritual'
  | 'craft'
  | 'music'
  | 'dance'
  | 'theater'
  | 'cuisine'
  | 'agriculture'
  | 'medicine'
  | 'storytelling'

export interface CulturalImpactScore {
  festivalId: string
  economicImpact: number
  socialCohesion: number
  culturalPreservation: number
  touristAttraction: number
  educationalValue: number
  overallScore: number
  factors: {
    participation: number
    mediaAttention: number
    economicActivity: number
    culturalAuthenticity: number
  }
}