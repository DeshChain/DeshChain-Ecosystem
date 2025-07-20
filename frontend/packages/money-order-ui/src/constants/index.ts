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

import { LanguageOption, FestivalInfo, ThemeConfig } from '../types';

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
  {
    code: 'hi',
    name: 'Hindi',
    nativeName: 'हिन्दी',
    region: 'North India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'en',
    name: 'English',
    nativeName: 'English',
    region: 'Pan India',
    supported: true,
    script: 'Latin'
  },
  {
    code: 'bn',
    name: 'Bengali',
    nativeName: 'বাংলা',
    region: 'East India',
    supported: true,
    script: 'Bengali'
  },
  {
    code: 'te',
    name: 'Telugu',
    nativeName: 'తెలుగు',
    region: 'South India',
    supported: true,
    script: 'Telugu'
  },
  {
    code: 'ta',
    name: 'Tamil',
    nativeName: 'தமிழ்',
    region: 'South India',
    supported: true,
    script: 'Tamil'
  },
  {
    code: 'mr',
    name: 'Marathi',
    nativeName: 'मराठी',
    region: 'West India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'ur',
    name: 'Urdu',
    nativeName: 'اردو',
    region: 'North India',
    supported: true,
    script: 'Arabic'
  },
  {
    code: 'gu',
    name: 'Gujarati',
    nativeName: 'ગુજરાતી',
    region: 'West India',
    supported: true,
    script: 'Gujarati'
  },
  {
    code: 'kn',
    name: 'Kannada',
    nativeName: 'ಕನ್ನಡ',
    region: 'South India',
    supported: true,
    script: 'Kannada'
  },
  {
    code: 'ml',
    name: 'Malayalam',
    nativeName: 'മലയാളം',
    region: 'South India',
    supported: true,
    script: 'Malayalam'
  },
  {
    code: 'pa',
    name: 'Punjabi',
    nativeName: 'ਪੰਜਾਬੀ',
    region: 'North India',
    supported: true,
    script: 'Gurmukhi'
  },
  {
    code: 'or',
    name: 'Odia',
    nativeName: 'ଓଡ଼ିଆ',
    region: 'East India',
    supported: true,
    script: 'Odia'
  },
  {
    code: 'as',
    name: 'Assamese',
    nativeName: 'অসমীয়া',
    region: 'Northeast India',
    supported: true,
    script: 'Bengali-Assamese'
  },
  {
    code: 'ks',
    name: 'Kashmiri',
    nativeName: 'کٲشُر',
    region: 'North India',
    supported: true,
    script: 'Arabic'
  },
  {
    code: 'ne',
    name: 'Nepali',
    nativeName: 'नेपाली',
    region: 'North India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'si',
    name: 'Sinhala',
    nativeName: 'සිංහල',
    region: 'Sri Lanka',
    supported: true,
    script: 'Sinhala'
  },
  {
    code: 'bho',
    name: 'Bhojpuri',
    nativeName: 'भोजपुरी',
    region: 'North India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'raj',
    name: 'Rajasthani',
    nativeName: 'राजस्थानी',
    region: 'West India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'mai',
    name: 'Maithili',
    nativeName: 'मैथिली',
    region: 'East India',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'sa',
    name: 'Sanskrit',
    nativeName: 'संस्कृत',
    region: 'Classical',
    supported: true,
    script: 'Devanagari'
  },
  {
    code: 'sd',
    name: 'Sindhi',
    nativeName: 'سنڌي',
    region: 'West India',
    supported: true,
    script: 'Arabic'
  },
  {
    code: 'kok',
    name: 'Konkani',
    nativeName: 'कोंकणी',
    region: 'West India',
    supported: true,
    script: 'Devanagari'
  }
];

export const FESTIVALS: FestivalInfo[] = [
  {
    festivalId: 'diwali',
    name: 'Diwali',
    description: 'Festival of Lights',
    startDate: '2024-11-01',
    endDate: '2024-11-05',
    bonusRate: 0.15,
    culturalTheme: 'prosperity',
    region: 'pan_india',
    significance: 'Victory of light over darkness, knowledge over ignorance',
    traditionalGreeting: 'दीपावली की शुभकामनाएं',
    active: false,
    daysRemaining: 290
  },
  {
    festivalId: 'holi',
    name: 'Holi',
    description: 'Festival of Colors',
    startDate: '2024-03-13',
    endDate: '2024-03-14',
    bonusRate: 0.10,
    culturalTheme: 'unity',
    region: 'north_india',
    significance: 'Celebration of spring, love, and new beginnings',
    traditionalGreeting: 'होली की शुभकामनाएं',
    active: false,
    daysRemaining: 53
  },
  {
    festivalId: 'independence_day',
    name: 'Independence Day',
    description: 'National Independence Day',
    startDate: '2024-08-15',
    endDate: '2024-08-15',
    bonusRate: 0.20,
    culturalTheme: 'patriotism',
    region: 'pan_india',
    significance: 'Celebrating Indias independence from British rule',
    traditionalGreeting: 'स्वतंत्रता दिवस की शुभकामनाएं',
    active: false,
    daysRemaining: 208
  },
  {
    festivalId: 'republic_day',
    name: 'Republic Day',
    description: 'National Republic Day',
    startDate: '2024-01-26',
    endDate: '2024-01-26',
    bonusRate: 0.20,
    culturalTheme: 'patriotism',
    region: 'pan_india',
    significance: 'Celebrating the Constitution of India',
    traditionalGreeting: 'गणतंत्र दिवस की शुभकामनाएं',
    active: true,
    daysRemaining: 0
  },
  {
    festivalId: 'eid',
    name: 'Eid ul-Fitr',
    description: 'Festival of Breaking the Fast',
    startDate: '2024-04-10',
    endDate: '2024-04-10',
    bonusRate: 0.12,
    culturalTheme: 'community',
    region: 'pan_india',
    significance: 'End of Ramadan fasting period',
    traditionalGreeting: 'ईद मुबारक',
    active: false,
    daysRemaining: 81
  },
  {
    festivalId: 'durga_puja',
    name: 'Durga Puja',
    description: 'Festival of Goddess Durga',
    startDate: '2024-10-10',
    endDate: '2024-10-14',
    bonusRate: 0.12,
    culturalTheme: 'power',
    region: 'east_india',
    significance: 'Celebrating divine feminine power',
    traditionalGreeting: 'दुर्गा पूजा की शुभकामनाएं',
    active: false,
    daysRemaining: 264
  },
  {
    festivalId: 'ganesh_chaturthi',
    name: 'Ganesh Chaturthi',
    description: 'Festival of Lord Ganesha',
    startDate: '2024-09-07',
    endDate: '2024-09-17',
    bonusRate: 0.10,
    culturalTheme: 'wisdom',
    region: 'west_india',
    significance: 'Celebrating the remover of obstacles',
    traditionalGreeting: 'गणेश चतुर्थी की शुभकामनाएं',
    active: false,
    daysRemaining: 231
  },
  {
    festivalId: 'navratri',
    name: 'Navratri',
    description: 'Nine Nights Festival',
    startDate: '2024-10-03',
    endDate: '2024-10-11',
    bonusRate: 0.08,
    culturalTheme: 'devotion',
    region: 'west_india',
    significance: 'Nine nights dedicated to Goddess Durga',
    traditionalGreeting: 'नवरात्रि की शुभकामनाएं',
    active: false,
    daysRemaining: 257
  }
];

export const CULTURAL_THEMES = [
  'independence',
  'prosperity',
  'community',
  'unity',
  'patriotism',
  'wisdom',
  'devotion',
  'power',
  'tradition',
  'family',
  'courage',
  'peace',
  'knowledge',
  'service',
  'truth'
];

export const DEFAULT_THEME_CONFIG: ThemeConfig = {
  primary: '#FF6B35', // Saffron
  secondary: '#138808', // Green
  accent: '#000080', // Navy Blue
  background: '#FAFAFA',
  surface: '#FFFFFF',
  text: '#212121',
  festivalColors: {
    diwali: {
      primary: '#FFD700', // Gold
      secondary: '#FF4500', // Orange Red
      accent: '#8B0000' // Dark Red
    },
    holi: {
      primary: '#FF69B4', // Hot Pink
      secondary: '#00FF00', // Lime Green
      accent: '#FF1493' // Deep Pink
    },
    independence_day: {
      primary: '#FF6B35', // Saffron
      secondary: '#138808', // Green
      accent: '#000080' // Navy Blue
    },
    republic_day: {
      primary: '#FF6B35', // Saffron
      secondary: '#138808', // Green
      accent: '#000080' // Navy Blue
    },
    eid: {
      primary: '#00A86B', // Jade Green
      secondary: '#FFD700', // Gold
      accent: '#4B0082' // Indigo
    },
    durga_puja: {
      primary: '#DC143C', // Crimson
      secondary: '#FFD700', // Gold
      accent: '#800080' // Purple
    },
    ganesh_chaturthi: {
      primary: '#FFA500', // Orange
      secondary: '#FF0000', // Red
      accent: '#FFFF00' // Yellow
    },
    navratri: {
      primary: '#FF69B4', // Hot Pink
      secondary: '#32CD32', // Lime Green
      accent: '#FF1493' // Deep Pink
    }
  },
  culturalElements: {
    borderStyle: 'traditional',
    pattern: 'mandala',
    iconSet: 'indian_classical'
  }
};

export const PATRIOTISM_LEVELS = [
  { level: 'bronze', minScore: 0, maxScore: 25, color: '#CD7F32', reward: '100' },
  { level: 'silver', minScore: 26, maxScore: 50, color: '#C0C0C0', reward: '250' },
  { level: 'gold', minScore: 51, maxScore: 75, color: '#FFD700', reward: '500' },
  { level: 'platinum', minScore: 76, maxScore: 100, color: '#E5E4E2', reward: '1000' }
];

export const TRANSACTION_PRIORITIES = [
  { 
    value: 'standard', 
    label: 'Standard', 
    description: 'Normal processing time (5-10 minutes)', 
    feeMultiplier: 1.0,
    estimatedTime: '5-10 min'
  },
  { 
    value: 'fast', 
    label: 'Fast', 
    description: 'Faster processing (2-5 minutes)', 
    feeMultiplier: 1.5,
    estimatedTime: '2-5 min'
  },
  { 
    value: 'instant', 
    label: 'Instant', 
    description: 'Immediate processing (< 1 minute)', 
    feeMultiplier: 2.0,
    estimatedTime: '< 1 min'
  }
];

export const POOL_TYPES = [
  {
    type: 'fixed_rate',
    name: 'Fixed Rate',
    description: 'Stable exchange rates for predictable transactions',
    icon: 'TrendingFlat',
    features: ['Predictable rates', 'Low slippage', 'Ideal for money orders']
  },
  {
    type: 'amm',
    name: 'AMM Pool',
    description: 'Automated market maker with dynamic pricing',
    icon: 'Timeline',
    features: ['Dynamic pricing', 'Yield farming', 'Higher returns']
  },
  {
    type: 'village',
    name: 'Village Pool',
    description: 'Community-driven pools with local governance',
    icon: 'Group',
    features: ['Community governed', 'Lower fees', 'Local benefits']
  },
  {
    type: 'concentrated',
    name: 'Concentrated',
    description: 'Capital-efficient liquidity provision',
    icon: 'CenterFocusStrong',
    features: ['Capital efficient', 'Professional trading', 'Advanced features']
  }
];

export const CULTURAL_CATEGORIES = [
  'wisdom',
  'motivation',
  'patriotism',
  'prosperity',
  'family',
  'community',
  'tradition',
  'spirituality',
  'courage',
  'peace',
  'unity',
  'service'
];

export const INDIAN_REGIONS = [
  { code: 'north_india', name: 'North India', states: ['Delhi', 'Punjab', 'Haryana', 'Himachal Pradesh', 'Jammu & Kashmir', 'Uttarakhand', 'Uttar Pradesh'] },
  { code: 'south_india', name: 'South India', states: ['Andhra Pradesh', 'Karnataka', 'Kerala', 'Tamil Nadu', 'Telangana'] },
  { code: 'west_india', name: 'West India', states: ['Gujarat', 'Maharashtra', 'Rajasthan', 'Goa'] },
  { code: 'east_india', name: 'East India', states: ['West Bengal', 'Odisha', 'Jharkhand', 'Bihar'] },
  { code: 'northeast_india', name: 'Northeast India', states: ['Assam', 'Arunachal Pradesh', 'Manipur', 'Meghalaya', 'Mizoram', 'Nagaland', 'Sikkim', 'Tripura'] },
  { code: 'central_india', name: 'Central India', states: ['Madhya Pradesh', 'Chhattisgarh'] }
];

export const API_ENDPOINTS = {
  pools: '/pools',
  moneyOrders: '/money-orders',
  swaps: '/swaps',
  liquidity: '/liquidity',
  cultural: '/cultural',
  analytics: '/analytics',
  receipts: '/receipts'
};

export const VALIDATION_RULES = {
  minAmount: '1',
  maxAmount: '10000000',
  maxSlippage: 0.05,
  maxMemoLength: 200,
  addressPattern: /^desh1[a-z0-9]{38}$/,
  postalCodePattern: /^[1-9][0-9]{5}$/
};

export const ANIMATION_DURATIONS = {
  fast: 150,
  normal: 300,
  slow: 500
};

export const BREAKPOINTS = {
  mobile: 576,
  tablet: 768,
  desktop: 992,
  wide: 1200
};