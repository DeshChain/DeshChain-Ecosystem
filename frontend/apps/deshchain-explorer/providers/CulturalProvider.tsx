'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

interface Festival {
  id: string
  name: string
  nameHindi: string
  date: Date
  description: string
  significance: string
  colors: {
    primary: string
    secondary: string
    accent: string
  }
  symbol: string
  isActive: boolean
  bonusPercentage: number
}

interface CulturalQuote {
  id: string
  text: string
  textHindi: string
  author: string
  authorHindi: string
  category: 'wisdom' | 'motivation' | 'peace' | 'dharma' | 'karma'
  source: string
}

interface CulturalEvent {
  id: string
  name: string
  description: string
  date: Date
  type: 'festival' | 'historical' | 'cultural' | 'modern'
  region: string
  importance: 'national' | 'regional' | 'local'
}

interface CulturalState {
  currentFestival: Festival | null
  upcomingFestivals: Festival[]
  dailyQuote: CulturalQuote | null
  culturalEvents: CulturalEvent[]
  selectedLanguage: string
  theme: 'light' | 'dark' | 'festival'
}

interface CulturalContextType extends CulturalState {
  setLanguage: (language: string) => void
  setTheme: (theme: 'light' | 'dark' | 'festival') => void
  getQuotesByCategory: (category: string) => CulturalQuote[]
  getFestivalByDate: (date: Date) => Festival | null
  refreshCulturalData: () => void
}

const CulturalContext = createContext<CulturalContextType | undefined>(undefined)

// Sample cultural data - in production this would come from the blockchain
const festivals: Festival[] = [
  {
    id: 'diwali',
    name: 'Diwali',
    nameHindi: 'दीवाली',
    date: new Date('2024-11-01'),
    description: 'Festival of Lights celebrating the triumph of good over evil',
    significance: 'Symbolizes the victory of light over darkness, knowledge over ignorance',
    colors: {
      primary: '#FFD700',
      secondary: '#FF6B35',
      accent: '#8B0000'
    },
    symbol: '🪔',
    isActive: false,
    bonusPercentage: 2.5
  },
  {
    id: 'holi',
    name: 'Holi',
    nameHindi: 'होली',
    date: new Date('2024-03-25'),
    description: 'Festival of Colors celebrating spring and love',
    significance: 'Celebrates the eternal love of Radha Krishna and the arrival of spring',
    colors: {
      primary: '#FF69B4',
      secondary: '#00FF00',
      accent: '#FFD700'
    },
    symbol: '🎨',
    isActive: false,
    bonusPercentage: 1.5
  },
  {
    id: 'independence_day',
    name: 'Independence Day',
    nameHindi: 'स्वतंत्रता दिवस',
    date: new Date('2024-08-15'),
    description: 'Celebrating India\'s independence from British rule',
    significance: 'Commemorates the freedom struggle and national sovereignty',
    colors: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808'
    },
    symbol: '🇮🇳',
    isActive: false,
    bonusPercentage: 3.0
  }
]

const culturalQuotes: CulturalQuote[] = [
  {
    id: '1',
    text: 'The best way to find yourself is to lose yourself in the service of others.',
    textHindi: 'अपने आप को खोजने का सबसे अच्छा तरीका दूसरों की सेवा में खुद को खो देना है।',
    author: 'Mahatma Gandhi',
    authorHindi: 'महात्मा गांधी',
    category: 'dharma',
    source: 'Gandhi\'s teachings'
  },
  {
    id: '2',
    text: 'You have the right to perform your actions, but you are not entitled to the fruits of your actions.',
    textHindi: 'कर्मण्येवाधिकारस्ते मा फलेषु कदाचन।',
    author: 'Bhagavad Gita',
    authorHindi: 'भगवद गीता',
    category: 'karma',
    source: 'Bhagavad Gita 2.47'
  },
  {
    id: '3',
    text: 'The mind is everything. What you think you become.',
    textHindi: 'मन ही सब कुछ है। जो आप सोचते हैं, आप वही बन जाते हैं।',
    author: 'Buddha',
    authorHindi: 'बुद्ध',
    category: 'wisdom',
    source: 'Buddhist teachings'
  }
]

const culturalEvents: CulturalEvent[] = [
  {
    id: '1',
    name: 'Gandhi Jayanti',
    description: 'Birth anniversary of Mahatma Gandhi',
    date: new Date('2024-10-02'),
    type: 'historical',
    region: 'National',
    importance: 'national'
  },
  {
    id: '2',
    name: 'Dussehra',
    description: 'Victory of good over evil',
    date: new Date('2024-10-24'),
    type: 'festival',
    region: 'National',
    importance: 'national'
  }
]

interface CulturalProviderProps {
  children: ReactNode
}

export function CulturalProvider({ children }: CulturalProviderProps) {
  const [state, setState] = useState<CulturalState>({
    currentFestival: null,
    upcomingFestivals: [],
    dailyQuote: null,
    culturalEvents: [],
    selectedLanguage: 'en',
    theme: 'light'
  })

  const setLanguage = (language: string) => {
    setState(prev => ({ ...prev, selectedLanguage: language }))
    localStorage.setItem('deshchain-language', language)
  }

  const setTheme = (theme: 'light' | 'dark' | 'festival') => {
    setState(prev => ({ ...prev, theme }))
    localStorage.setItem('deshchain-theme', theme)
  }

  const getQuotesByCategory = (category: string) => {
    return culturalQuotes.filter(quote => quote.category === category)
  }

  const getFestivalByDate = (date: Date) => {
    return festivals.find(festival => 
      festival.date.toDateString() === date.toDateString()
    ) || null
  }

  const refreshCulturalData = () => {
    const today = new Date()
    
    // Find current festival (within 7 days)
    const currentFestival = festivals.find(festival => {
      const daysDiff = Math.abs(festival.date.getTime() - today.getTime()) / (1000 * 60 * 60 * 24)
      return daysDiff <= 7
    })

    // Get upcoming festivals (next 3 months)
    const upcomingFestivals = festivals.filter(festival => {
      const monthsDiff = (festival.date.getTime() - today.getTime()) / (1000 * 60 * 60 * 24 * 30)
      return monthsDiff > 0 && monthsDiff <= 3
    }).sort((a, b) => a.date.getTime() - b.date.getTime())

    // Get random daily quote
    const dailyQuote = culturalQuotes[Math.floor(Math.random() * culturalQuotes.length)]

    // Filter upcoming cultural events
    const upcomingEvents = culturalEvents.filter(event => 
      event.date.getTime() > today.getTime()
    ).sort((a, b) => a.date.getTime() - b.date.getTime())

    setState(prev => ({
      ...prev,
      currentFestival: currentFestival || null,
      upcomingFestivals,
      dailyQuote,
      culturalEvents: upcomingEvents
    }))
  }

  // Initialize cultural data
  useEffect(() => {
    // Load saved preferences
    const savedLanguage = localStorage.getItem('deshchain-language')
    const savedTheme = localStorage.getItem('deshchain-theme')

    if (savedLanguage) {
      setState(prev => ({ ...prev, selectedLanguage: savedLanguage }))
    }

    if (savedTheme) {
      setState(prev => ({ ...prev, theme: savedTheme as any }))
    }

    // Load cultural data
    refreshCulturalData()

    // Refresh daily at midnight
    const now = new Date()
    const tomorrow = new Date(now)
    tomorrow.setDate(tomorrow.getDate() + 1)
    tomorrow.setHours(0, 0, 0, 0)
    
    const msUntilMidnight = tomorrow.getTime() - now.getTime()
    
    const timeout = setTimeout(() => {
      refreshCulturalData()
      
      // Set daily refresh interval
      const interval = setInterval(refreshCulturalData, 24 * 60 * 60 * 1000)
      return () => clearInterval(interval)
    }, msUntilMidnight)

    return () => clearTimeout(timeout)
  }, [])

  const value: CulturalContextType = {
    ...state,
    setLanguage,
    setTheme,
    getQuotesByCategory,
    getFestivalByDate,
    refreshCulturalData
  }

  return (
    <CulturalContext.Provider value={value}>
      {children}
    </CulturalContext.Provider>
  )
}

export function useCultural() {
  const context = useContext(CulturalContext)
  if (context === undefined) {
    throw new Error('useCultural must be used within a CulturalProvider')
  }
  return context
}