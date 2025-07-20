/**
 * Cultural utilities for DeshChain SDK
 */

export class CulturalUtils {
  /**
   * Get state code from pincode
   */
  static getStateFromPincode(pincode: string): string {
    const firstDigit = pincode.charAt(0)
    
    const stateMap: Record<string, string> = {
      '1': 'Delhi',
      '2': 'Haryana/Himachal Pradesh',
      '3': 'Punjab/Jammu & Kashmir',
      '4': 'Rajasthan',
      '5': 'Uttar Pradesh/Uttarakhand',
      '6': 'Bihar/Jharkhand',
      '7': 'West Bengal/Sikkim',
      '8': 'Odisha',
      '9': 'Assam/Arunachal Pradesh/Manipur/Meghalaya/Mizoram/Nagaland/Tripura',
    }

    return stateMap[firstDigit] || 'Unknown'
  }

  /**
   * Get regional language from state
   */
  static getRegionalLanguage(state: string): string {
    const languageMap: Record<string, string> = {
      'Andhra Pradesh': 'Telugu',
      'Arunachal Pradesh': 'Hindi',
      'Assam': 'Assamese',
      'Bihar': 'Hindi',
      'Chhattisgarh': 'Hindi',
      'Goa': 'Konkani',
      'Gujarat': 'Gujarati',
      'Haryana': 'Hindi',
      'Himachal Pradesh': 'Hindi',
      'Jammu and Kashmir': 'Kashmiri',
      'Jharkhand': 'Hindi',
      'Karnataka': 'Kannada',
      'Kerala': 'Malayalam',
      'Madhya Pradesh': 'Hindi',
      'Maharashtra': 'Marathi',
      'Manipur': 'Manipuri',
      'Meghalaya': 'English',
      'Mizoram': 'Mizo',
      'Nagaland': 'English',
      'Odisha': 'Odia',
      'Punjab': 'Punjabi',
      'Rajasthan': 'Hindi',
      'Sikkim': 'Nepali',
      'Tamil Nadu': 'Tamil',
      'Telangana': 'Telugu',
      'Tripura': 'Bengali',
      'Uttar Pradesh': 'Hindi',
      'Uttarakhand': 'Hindi',
      'West Bengal': 'Bengali',
      'Delhi': 'Hindi',
    }

    return languageMap[state] || 'Hindi'
  }

  /**
   * Get major festivals by state
   */
  static getMajorFestivals(state: string): string[] {
    const festivalMap: Record<string, string[]> = {
      'West Bengal': ['Durga Puja', 'Kali Puja', 'Poila Boishakh'],
      'Tamil Nadu': ['Pongal', 'Diwali', 'Navaratri'],
      'Kerala': ['Onam', 'Vishu', 'Thrissur Pooram'],
      'Maharashtra': ['Ganesh Chaturthi', 'Gudi Padwa', 'Navratri'],
      'Gujarat': ['Navratri', 'Kite Festival', 'Diwali'],
      'Punjab': ['Baisakhi', 'Karva Chauth', 'Lohri'],
      'Rajasthan': ['Teej', 'Gangaur', 'Desert Festival'],
      'Karnataka': ['Mysore Dasara', 'Ugadi', 'Karaga'],
      'Andhra Pradesh': ['Ugadi', 'Sankranti', 'Bonalu'],
      'Odisha': ['Jagannath Puri Rath Yatra', 'Kali Puja', 'Durga Puja'],
      'Bihar': ['Chhath Puja', 'Bihula', 'Sonepur Mela'],
      'Assam': ['Bihu', 'Durga Puja', 'Kali Puja'],
    }

    return festivalMap[state] || ['Diwali', 'Holi', 'Dussehra']
  }

  /**
   * Get current season based on date
   */
  static getCurrentSeason(date: Date = new Date()): string {
    const month = date.getMonth() + 1 // 1-12

    if (month >= 3 && month <= 5) return 'Spring'
    if (month >= 6 && month <= 8) return 'Monsoon'
    if (month >= 9 && month <= 11) return 'Autumn'
    return 'Winter'
  }

  /**
   * Check if date is a major festival
   */
  static isMajorFestival(date: Date): { isFestival: boolean; festival?: string } {
    const month = date.getMonth() + 1
    const day = date.getDate()

    // Major festivals (approximate dates)
    const festivals: Record<string, { month: number; day: number }> = {
      'Independence Day': { month: 8, day: 15 },
      'Republic Day': { month: 1, day: 26 },
      'Gandhi Jayanti': { month: 10, day: 2 },
      'Diwali': { month: 10, day: 24 }, // Approximate
      'Holi': { month: 3, day: 13 }, // Approximate
      'Dussehra': { month: 10, day: 15 }, // Approximate
    }

    for (const [festivalName, festivalDate] of Object.entries(festivals)) {
      if (festivalDate.month === month && festivalDate.day === day) {
        return { isFestival: true, festival: festivalName }
      }
    }

    return { isFestival: false }
  }

  /**
   * Get cultural bonus multiplier based on festival and region
   */
  static getCulturalBonusMultiplier(
    festival: string,
    state: string,
    baseMultiplier: number = 1.0
  ): number {
    const regionalFestivals = CulturalUtils.getMajorFestivals(state)
    
    if (regionalFestivals.includes(festival)) {
      return baseMultiplier * 1.5 // 50% bonus for regional festivals
    }

    const nationalFestivals = ['Diwali', 'Holi', 'Independence Day', 'Republic Day']
    if (nationalFestivals.includes(festival)) {
      return baseMultiplier * 1.2 // 20% bonus for national festivals
    }

    return baseMultiplier
  }

  /**
   * Generate greeting in regional language
   */
  static getRegionalGreeting(state: string, timeOfDay: 'morning' | 'afternoon' | 'evening' = 'morning'): string {
    const greetings: Record<string, Record<string, string>> = {
      'West Bengal': {
        morning: 'শুভ সকাল (Shubho Shokal)',
        afternoon: 'শুভ দুপুর (Shubho Dupur)',
        evening: 'শুভ সন্ধ্যা (Shubho Sandhya)'
      },
      'Tamil Nadu': {
        morning: 'வணக்கம் (Vanakkam)',
        afternoon: 'வணக্கம் (Vanakkam)',
        evening: 'வணக্கம் (Vanakkam)'
      },
      'Maharashtra': {
        morning: 'सुप्रभात (Suprabhat)',
        afternoon: 'नमस्कार (Namaskaar)',
        evening: 'शुभ संध्या (Shubh Sandhya)'
      },
      'Punjab': {
        morning: 'ਸਤ ਸ੍ਰੀ ਅਕਾਲ (Sat Sri Akal)',
        afternoon: 'ਸਤ ਸ੍ਰੀ ਅਕਾਲ (Sat Sri Akal)',
        evening: 'ਸਤ ਸ੍ਰੀ ਅਕਾਲ (Sat Sri Akal)'
      }
    }

    const defaultGreeting = {
      morning: 'सुप्रभात (Suprabhat)',
      afternoon: 'नमस्ते (Namaste)',
      evening: 'शुभ संध्या (Shubh Sandhya)'
    }

    return greetings[state]?.[timeOfDay] || defaultGreeting[timeOfDay]
  }

  /**
   * Convert text to Devanagari (simplified)
   */
  static toDevanagari(text: string): string {
    // This is a simplified mapping for common words
    const transliterationMap: Record<string, string> = {
      'namaste': 'नमस्ते',
      'dhanyawad': 'धन्यवाद',
      'swagat': 'स्वागत',
      'deshchain': 'देशचेन',
      'namo': 'नमो',
      'bharat': 'भारत',
      'seva': 'सेवा',
      'samman': 'सम्मान',
      'sanskriti': 'संस्कृति',
      'parampara': 'परम्परा'
    }

    return transliterationMap[text.toLowerCase()] || text
  }

  /**
   * Get auspicious time (muhurat) for transactions
   */
  static getAuspiciousTime(date: Date = new Date()): { isAuspicious: boolean; reason: string } {
    const hour = date.getHours()
    
    // Morning hours (6 AM - 10 AM) are generally auspicious
    if (hour >= 6 && hour <= 10) {
      return { isAuspicious: true, reason: 'Brahma Muhurat - Most auspicious time' }
    }
    
    // Evening hours (5 PM - 7 PM) are also good
    if (hour >= 17 && hour <= 19) {
      return { isAuspicious: true, reason: 'Sandhya Kaal - Evening auspicious time' }
    }
    
    // Avoid Rahu Kaal (typically 4:30 PM - 6:00 PM)
    if (hour >= 16 && hour <= 18) {
      return { isAuspicious: false, reason: 'Rahu Kaal - Inauspicious time' }
    }
    
    return { isAuspicious: true, reason: 'Normal time' }
  }

  /**
   * Calculate cultural impact score
   */
  static calculateCulturalImpact(
    festivalName: string,
    regionParticipation: number,
    traditionalValue: number,
    modernRelevance: number
  ): number {
    // Weights for different factors
    const weights = {
      festival: 0.3,
      participation: 0.3,
      tradition: 0.2,
      relevance: 0.2
    }

    const festivalScore = this.getFestivalScore(festivalName)
    
    const totalScore = 
      (festivalScore * weights.festival) +
      (regionParticipation * weights.participation) +
      (traditionalValue * weights.tradition) +
      (modernRelevance * weights.relevance)

    return Math.min(100, Math.max(0, totalScore))
  }

  /**
   * Get festival importance score
   */
  private static getFestivalScore(festivalName: string): number {
    const scores: Record<string, number> = {
      'Diwali': 100,
      'Holi': 95,
      'Dussehra': 90,
      'Durga Puja': 95,
      'Ganesh Chaturthi': 90,
      'Navratri': 85,
      'Onam': 90,
      'Pongal': 85,
      'Baisakhi': 80,
      'Karva Chauth': 75,
      'Independence Day': 100,
      'Republic Day': 100,
      'Gandhi Jayanti': 95
    }

    return scores[festivalName] || 50
  }

  /**
   * Format currency in Indian style
   */
  static formatIndianCurrency(amount: number): string {
    const formatter = new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })

    return formatter.format(amount)
  }

  /**
   * Convert number to Indian numbering system (Lakh, Crore)
   */
  static toIndianNumbering(amount: number): string {
    if (amount >= 10000000) {
      return `${(amount / 10000000).toFixed(2)} Cr`
    } else if (amount >= 100000) {
      return `${(amount / 100000).toFixed(2)} L`
    } else if (amount >= 1000) {
      return `${(amount / 1000).toFixed(2)} K`
    }
    return amount.toString()
  }
}