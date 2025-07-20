/**
 * Festival utilities for DeshChain SDK
 */

export interface FestivalInfo {
  name: string
  nameHindi: string
  date: Date
  duration: number
  category: string
  significance: string
  bonusMultiplier: number
  regions: string[]
}

export class FestivalUtils {
  /**
   * Get current active festivals
   */
  static getCurrentFestivals(date: Date = new Date()): FestivalInfo[] {
    const festivals = FestivalUtils.getAllFestivals()
    const currentDate = new Date(date)
    
    return festivals.filter(festival => {
      const festivalStart = new Date(festival.date)
      const festivalEnd = new Date(festivalStart.getTime() + (festival.duration * 24 * 60 * 60 * 1000))
      
      return currentDate >= festivalStart && currentDate <= festivalEnd
    })
  }

  /**
   * Get upcoming festivals
   */
  static getUpcomingFestivals(
    daysAhead: number = 30,
    date: Date = new Date()
  ): FestivalInfo[] {
    const festivals = FestivalUtils.getAllFestivals()
    const currentDate = new Date(date)
    const futureDate = new Date(currentDate.getTime() + (daysAhead * 24 * 60 * 60 * 1000))
    
    return festivals.filter(festival => {
      const festivalDate = new Date(festival.date)
      return festivalDate > currentDate && festivalDate <= futureDate
    }).sort((a, b) => a.date.getTime() - b.date.getTime())
  }

  /**
   * Get festival by name
   */
  static getFestivalByName(name: string): FestivalInfo | null {
    const festivals = FestivalUtils.getAllFestivals()
    return festivals.find(f => 
      f.name.toLowerCase() === name.toLowerCase() || 
      f.nameHindi === name
    ) || null
  }

  /**
   * Check if date is during a festival
   */
  static isFestivalDay(date: Date = new Date()): { 
    isFestival: boolean
    festivals: FestivalInfo[]
  } {
    const activeFestivals = FestivalUtils.getCurrentFestivals(date)
    return {
      isFestival: activeFestivals.length > 0,
      festivals: activeFestivals
    }
  }

  /**
   * Get festival bonus multiplier for date
   */
  static getFestivalBonus(date: Date = new Date()): number {
    const activeFestivals = FestivalUtils.getCurrentFestivals(date)
    
    if (activeFestivals.length === 0) return 1.0
    
    // Return the highest bonus multiplier if multiple festivals are active
    return Math.max(...activeFestivals.map(f => f.bonusMultiplier))
  }

  /**
   * Get festivals by region
   */
  static getFestivalsByRegion(region: string): FestivalInfo[] {
    const festivals = FestivalUtils.getAllFestivals()
    return festivals.filter(festival => 
      festival.regions.includes(region) || 
      festival.regions.includes('All India')
    )
  }

  /**
   * Get festivals by category
   */
  static getFestivalsByCategory(category: string): FestivalInfo[] {
    const festivals = FestivalUtils.getAllFestivals()
    return festivals.filter(festival => 
      festival.category.toLowerCase() === category.toLowerCase()
    )
  }

  /**
   * Calculate days until next festival
   */
  static getDaysUntilNextFestival(date: Date = new Date()): {
    festival: FestivalInfo | null
    days: number
  } {
    const upcoming = FestivalUtils.getUpcomingFestivals(365, date)
    
    if (upcoming.length === 0) {
      return { festival: null, days: -1 }
    }

    const nextFestival = upcoming[0]
    const currentDate = new Date(date)
    const festivalDate = new Date(nextFestival.date)
    const diffTime = festivalDate.getTime() - currentDate.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

    return { festival: nextFestival, days: diffDays }
  }

  /**
   * Get festival calendar for month
   */
  static getMonthlyFestivalCalendar(
    year: number,
    month: number
  ): FestivalInfo[] {
    const festivals = FestivalUtils.getAllFestivals()
    
    return festivals.filter(festival => {
      const festivalDate = new Date(festival.date)
      return festivalDate.getFullYear() === year && 
             festivalDate.getMonth() === month - 1
    }).sort((a, b) => a.date.getTime() - b.date.getTime())
  }

  /**
   * Generate festival greeting
   */
  static getFestivalGreeting(festivalName: string): {
    english: string
    hindi: string
  } {
    const greetings: Record<string, { english: string; hindi: string }> = {
      'Diwali': {
        english: 'Happy Diwali! May the festival of lights bring joy and prosperity.',
        hindi: 'दीपावली की हार्दिक शुभकामनाएं!'
      },
      'Holi': {
        english: 'Happy Holi! May your life be filled with colors of joy.',
        hindi: 'होली की शुभकामनाएं!'
      },
      'Dussehra': {
        english: 'Happy Dussehra! May good triumph over evil.',
        hindi: 'दशहरा की शुभकामनाएं!'
      },
      'Independence Day': {
        english: 'Happy Independence Day! Jai Hind!',
        hindi: 'स्वतंत्रता दिवस की शुभकामनाएं! जय हिन्द!'
      },
      'Republic Day': {
        english: 'Happy Republic Day! Jai Hind!',
        hindi: 'गणतंत्र दिवस की शुभकामनाएं! जय हिन्द!'
      }
    }

    return greetings[festivalName] || {
      english: `Happy ${festivalName}!`,
      hindi: `${festivalName} की शुभकामनाएं!`
    }
  }

  /**
   * Get all festivals (static data)
   */
  private static getAllFestivals(): FestivalInfo[] {
    const currentYear = new Date().getFullYear()
    
    return [
      {
        name: 'Republic Day',
        nameHindi: 'गणतंत्र दिवस',
        date: new Date(currentYear, 0, 26), // January 26
        duration: 1,
        category: 'National',
        significance: 'Celebrates the adoption of the Constitution of India',
        bonusMultiplier: 1.5,
        regions: ['All India']
      },
      {
        name: 'Holi',
        nameHindi: 'होली',
        date: new Date(currentYear, 2, 13), // March 13 (approximate)
        duration: 2,
        category: 'Religious',
        significance: 'Festival of colors celebrating the victory of good over evil',
        bonusMultiplier: 1.3,
        regions: ['All India']
      },
      {
        name: 'Independence Day',
        nameHindi: 'स्वतंत्रता दिवस',
        date: new Date(currentYear, 7, 15), // August 15
        duration: 1,
        category: 'National',
        significance: 'Celebrates India\'s independence from British rule',
        bonusMultiplier: 1.5,
        regions: ['All India']
      },
      {
        name: 'Gandhi Jayanti',
        nameHindi: 'गांधी जयंती',
        date: new Date(currentYear, 9, 2), // October 2
        duration: 1,
        category: 'National',
        significance: 'Birthday of Mahatma Gandhi',
        bonusMultiplier: 1.3,
        regions: ['All India']
      },
      {
        name: 'Dussehra',
        nameHindi: 'दशहरा',
        date: new Date(currentYear, 9, 15), // October 15 (approximate)
        duration: 1,
        category: 'Religious',
        significance: 'Celebrates the victory of Lord Rama over Ravana',
        bonusMultiplier: 1.3,
        regions: ['All India']
      },
      {
        name: 'Diwali',
        nameHindi: 'दीपावली',
        date: new Date(currentYear, 10, 12), // November 12 (approximate)
        duration: 5,
        category: 'Religious',
        significance: 'Festival of lights celebrating the return of Lord Rama',
        bonusMultiplier: 1.5,
        regions: ['All India']
      },
      {
        name: 'Durga Puja',
        nameHindi: 'दुर्गा पूजा',
        date: new Date(currentYear, 9, 10), // October 10 (approximate)
        duration: 10,
        category: 'Religious',
        significance: 'Worship of Goddess Durga',
        bonusMultiplier: 1.4,
        regions: ['West Bengal', 'Assam', 'Bihar', 'Odisha']
      },
      {
        name: 'Ganesh Chaturthi',
        nameHindi: 'गणेश चतुर्थी',
        date: new Date(currentYear, 7, 22), // August 22 (approximate)
        duration: 11,
        category: 'Religious',
        significance: 'Birthday of Lord Ganesha',
        bonusMultiplier: 1.3,
        regions: ['Maharashtra', 'Karnataka', 'Andhra Pradesh', 'Tamil Nadu']
      },
      {
        name: 'Onam',
        nameHindi: 'ओणम',
        date: new Date(currentYear, 8, 8), // September 8 (approximate)
        duration: 10,
        category: 'Regional',
        significance: 'Harvest festival of Kerala',
        bonusMultiplier: 1.3,
        regions: ['Kerala']
      },
      {
        name: 'Pongal',
        nameHindi: 'पोंगल',
        date: new Date(currentYear, 0, 14), // January 14
        duration: 4,
        category: 'Regional',
        significance: 'Harvest festival of Tamil Nadu',
        bonusMultiplier: 1.3,
        regions: ['Tamil Nadu']
      },
      {
        name: 'Baisakhi',
        nameHindi: 'बैसाखी',
        date: new Date(currentYear, 3, 13), // April 13
        duration: 1,
        category: 'Regional',
        significance: 'Harvest festival and Sikh New Year',
        bonusMultiplier: 1.3,
        regions: ['Punjab', 'Haryana']
      },
      {
        name: 'Karva Chauth',
        nameHindi: 'करवा चौथ',
        date: new Date(currentYear, 9, 20), // October 20 (approximate)
        duration: 1,
        category: 'Cultural',
        significance: 'Festival of married women for their husbands\' long life',
        bonusMultiplier: 1.2,
        regions: ['North India']
      }
    ]
  }

  /**
   * Get festival theme colors
   */
  static getFestivalTheme(festivalName: string): {
    primaryColor: string
    secondaryColor: string
    accentColor: string
  } {
    const themes: Record<string, { primaryColor: string; secondaryColor: string; accentColor: string }> = {
      'Diwali': { primaryColor: '#FFD700', secondaryColor: '#FF6B35', accentColor: '#8B0000' },
      'Holi': { primaryColor: '#FF69B4', secondaryColor: '#00FF00', accentColor: '#FFD700' },
      'Independence Day': { primaryColor: '#FF9933', secondaryColor: '#FFFFFF', accentColor: '#138808' },
      'Republic Day': { primaryColor: '#FF9933', secondaryColor: '#FFFFFF', accentColor: '#138808' },
      'Dussehra': { primaryColor: '#FF4500', secondaryColor: '#FFD700', accentColor: '#8B0000' },
      'Ganesh Chaturthi': { primaryColor: '#FF6347', secondaryColor: '#FFD700', accentColor: '#8B4513' },
      'Durga Puja': { primaryColor: '#DC143C', secondaryColor: '#FFD700', accentColor: '#8B0000' }
    }

    return themes[festivalName] || { 
      primaryColor: '#FF6B35', 
      secondaryColor: '#FFD700', 
      accentColor: '#8B0000' 
    }
  }

  /**
   * Check if festival is national holiday
   */
  static isNationalHoliday(festivalName: string): boolean {
    const nationalHolidays = [
      'Independence Day',
      'Republic Day',
      'Gandhi Jayanti'
    ]
    
    return nationalHolidays.includes(festivalName)
  }

  /**
   * Get festival-specific transaction bonuses
   */
  static getFestivalTransactionBonus(
    festivalName: string,
    transactionType: string
  ): number {
    const bonuses: Record<string, Record<string, number>> = {
      'Diwali': {
        'send': 1.5,
        'lending': 1.3,
        'trading': 1.2,
        'donation': 2.0
      },
      'Independence Day': {
        'send': 1.3,
        'lending': 1.2,
        'trading': 1.1,
        'donation': 1.8
      },
      'Holi': {
        'send': 1.2,
        'lending': 1.1,
        'trading': 1.2,
        'donation': 1.5
      }
    }

    return bonuses[festivalName]?.[transactionType] || 1.0
  }
}