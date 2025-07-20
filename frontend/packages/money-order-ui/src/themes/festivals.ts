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

export interface Festival {
  id: string;
  name: string;
  localNames: Record<string, string>;
  startDate: string; // ISO date
  endDate: string; // ISO date
  type: 'national' | 'religious' | 'regional' | 'cultural';
  regions: string[];
  theme: FestivalTheme;
  animations: FestivalAnimations;
  greetings: Record<string, string>;
  specialOffers?: SpecialOffer[];
}

export interface FestivalTheme {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  backgroundGradient: string;
  textPrimary: string;
  textSecondary: string;
  decorativeElements: string[];
  iconSet: string;
  fontFamily?: string;
}

export interface FestivalAnimations {
  particles: ParticleConfig[];
  transitions: string;
  backgroundAnimation?: string;
  iconAnimations?: Record<string, string>;
}

export interface ParticleConfig {
  type: 'diya' | 'flower' | 'rangoli' | 'firework' | 'star' | 'custom';
  count: number;
  size: { min: number; max: number };
  speed: { min: number; max: number };
  colors: string[];
  customImage?: string;
}

export interface SpecialOffer {
  type: 'fee_discount' | 'bonus_amount' | 'special_quote';
  value: number | string;
  message: Record<string, string>;
}

// Festival definitions for 2024-2025
export const FESTIVALS: Festival[] = [
  {
    id: 'diwali',
    name: 'Diwali',
    localNames: {
      hi: 'दीपावली',
      bn: 'দীপাবলি',
      te: 'దీపావళి',
      ta: 'தீபாவளி',
      gu: 'દિવાળી',
      mr: 'दिवाळी',
      ml: 'ദീപാവലി',
      kn: 'ದೀಪಾವಳಿ',
      or: 'ଦୀପାବଳି',
      pa: 'ਦਿਵਾਲੀ'
    },
    startDate: '2024-10-31',
    endDate: '2024-11-04',
    type: 'national',
    regions: ['all'],
    theme: {
      primary: '#FF6B35',
      secondary: '#FFD700',
      accent: '#FF1744',
      background: '#FFF3E0',
      backgroundGradient: 'linear-gradient(135deg, #FFF3E0 0%, #FFE0B2 100%)',
      textPrimary: '#5D4037',
      textSecondary: '#8D6E63',
      decorativeElements: ['diya', 'rangoli', 'fireworks'],
      iconSet: 'diwali'
    },
    animations: {
      particles: [
        {
          type: 'diya',
          count: 15,
          size: { min: 20, max: 40 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FFD700', '#FF6B35', '#FFA726']
        },
        {
          type: 'firework',
          count: 8,
          size: { min: 30, max: 60 },
          speed: { min: 1, max: 2 },
          colors: ['#FF1744', '#FFD700', '#4CAF50', '#2196F3']
        }
      ],
      transitions: 'glow',
      backgroundAnimation: 'sparkle'
    },
    greetings: {
      en: 'Happy Diwali! May the festival of lights brighten your life',
      hi: 'दीपावली की हार्दिक शुभकामनाएं! 🪔',
      bn: 'শুভ দীপাবলি! আলোর উৎসব আপনার জীবন আলোকিত করুক'
    },
    specialOffers: [
      {
        type: 'fee_discount',
        value: 25,
        message: {
          en: '25% off on transaction fees this Diwali!',
          hi: 'इस दिवाली लेनदेन शुल्क पर 25% की छूट!'
        }
      }
    ]
  },
  {
    id: 'holi',
    name: 'Holi',
    localNames: {
      hi: 'होली',
      bn: 'হোলি',
      te: 'హోళి',
      ta: 'ஹோலி',
      gu: 'હોળી',
      mr: 'होळी',
      pa: 'ਹੋਲੀ'
    },
    startDate: '2025-03-13',
    endDate: '2025-03-14',
    type: 'national',
    regions: ['all'],
    theme: {
      primary: '#E91E63',
      secondary: '#9C27B0',
      accent: '#00BCD4',
      background: '#FCE4EC',
      backgroundGradient: 'linear-gradient(135deg, #FCE4EC 0%, #F3E5F5 50%, #E1F5FE 100%)',
      textPrimary: '#4A148C',
      textSecondary: '#7B1FA2',
      decorativeElements: ['gulal', 'pichkari', 'colors'],
      iconSet: 'holi'
    },
    animations: {
      particles: [
        {
          type: 'custom',
          customImage: '/assets/gulal.png',
          count: 20,
          size: { min: 15, max: 35 },
          speed: { min: 1, max: 3 },
          colors: ['#E91E63', '#9C27B0', '#00BCD4', '#4CAF50', '#FF9800']
        }
      ],
      transitions: 'splash',
      backgroundAnimation: 'rainbow'
    },
    greetings: {
      en: 'Happy Holi! May your life be filled with colors of joy',
      hi: 'होली मुबारक! रंगों का त्योहार आपके जीवन में खुशियां लाए',
      bn: 'শুভ হোলি! রঙের উৎসব আপনার জীবনে আনন্দ নিয়ে আসুক'
    }
  },
  {
    id: 'eid',
    name: 'Eid al-Fitr',
    localNames: {
      ur: 'عید الفطر',
      hi: 'ईद उल-फितर',
      bn: 'ঈদ উল ফিতর',
      ml: 'പെരുന്നാൾ'
    },
    startDate: '2024-04-10',
    endDate: '2024-04-11',
    type: 'religious',
    regions: ['all'],
    theme: {
      primary: '#1B5E20',
      secondary: '#4CAF50',
      accent: '#FFD700',
      background: '#E8F5E9',
      backgroundGradient: 'linear-gradient(135deg, #E8F5E9 0%, #C8E6C9 100%)',
      textPrimary: '#1B5E20',
      textSecondary: '#388E3C',
      decorativeElements: ['crescent', 'star', 'mosque'],
      iconSet: 'eid'
    },
    animations: {
      particles: [
        {
          type: 'star',
          count: 15,
          size: { min: 10, max: 25 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FFD700', '#FFFFFF']
        }
      ],
      transitions: 'fade',
      backgroundAnimation: 'twinkle'
    },
    greetings: {
      en: 'Eid Mubarak! May Allah bless you with happiness',
      ur: 'عید مبارک! اللہ آپ کو خوشیوں سے نوازے',
      hi: 'ईद मुबारक! अल्लाह आपको खुशियों से नवाज़े'
    }
  },
  {
    id: 'durga_puja',
    name: 'Durga Puja',
    localNames: {
      bn: 'দুর্গা পূজা',
      hi: 'दुर्गा पूजा',
      as: 'দুৰ্গা পূজা'
    },
    startDate: '2024-10-09',
    endDate: '2024-10-13',
    type: 'religious',
    regions: ['West Bengal', 'Assam', 'Odisha', 'Tripura'],
    theme: {
      primary: '#D32F2F',
      secondary: '#FFC107',
      accent: '#FF5722',
      background: '#FFEBEE',
      backgroundGradient: 'linear-gradient(135deg, #FFEBEE 0%, #FFF9C4 100%)',
      textPrimary: '#B71C1C',
      textSecondary: '#D32F2F',
      decorativeElements: ['durga', 'dhak', 'lotus'],
      iconSet: 'durga_puja'
    },
    animations: {
      particles: [
        {
          type: 'flower',
          count: 12,
          size: { min: 20, max: 35 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#D32F2F', '#FFC107', '#FFFFFF']
        }
      ],
      transitions: 'devotion',
      backgroundAnimation: 'blessing'
    },
    greetings: {
      en: 'Shubho Durga Puja! May Maa Durga bless you',
      bn: 'শুভ দুর্গা পূজা! মা দুর্গা আপনাকে আশীর্বাদ করুন',
      hi: 'दुर्गा पूजा की शुभकामनाएं!'
    }
  },
  {
    id: 'onam',
    name: 'Onam',
    localNames: {
      ml: 'ഓണം',
      ta: 'ஓணம்',
      kn: 'ಓಣಂ'
    },
    startDate: '2024-08-29',
    endDate: '2024-09-07',
    type: 'regional',
    regions: ['Kerala'],
    theme: {
      primary: '#388E3C',
      secondary: '#FDD835',
      accent: '#FF6F00',
      background: '#F1F8E9',
      backgroundGradient: 'linear-gradient(135deg, #F1F8E9 0%, #FFF9C4 100%)',
      textPrimary: '#1B5E20',
      textSecondary: '#388E3C',
      decorativeElements: ['pookalam', 'boat', 'umbrella'],
      iconSet: 'onam'
    },
    animations: {
      particles: [
        {
          type: 'flower',
          count: 20,
          size: { min: 15, max: 30 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FF6F00', '#FDD835', '#388E3C', '#E91E63']
        }
      ],
      transitions: 'bloom',
      backgroundAnimation: 'wave'
    },
    greetings: {
      en: 'Happy Onam! May King Mahabali bless you',
      ml: 'ഓണാശംസകൾ! മാവേലി നാടുവാണീടും കാലം',
      ta: 'ஓணம் நல்வாழ்த்துக்கள்!'
    }
  },
  {
    id: 'pongal',
    name: 'Pongal',
    localNames: {
      ta: 'பொங்கல்',
      te: 'మకర సంక్రాంతి',
      kn: 'ಮಕರ ಸಂಕ್ರಾಂತಿ'
    },
    startDate: '2025-01-14',
    endDate: '2025-01-15',
    type: 'regional',
    regions: ['Tamil Nadu', 'Andhra Pradesh', 'Telangana', 'Karnataka'],
    theme: {
      primary: '#F57C00',
      secondary: '#4CAF50',
      accent: '#FFD700',
      background: '#FFF3E0',
      backgroundGradient: 'linear-gradient(135deg, #FFF3E0 0%, #F4FF81 100%)',
      textPrimary: '#E65100',
      textSecondary: '#F57C00',
      decorativeElements: ['sun', 'sugarcane', 'pot'],
      iconSet: 'pongal'
    },
    animations: {
      particles: [
        {
          type: 'custom',
          customImage: '/assets/sun.png',
          count: 8,
          size: { min: 25, max: 45 },
          speed: { min: 0.3, max: 0.8 },
          colors: ['#FFD700', '#F57C00']
        }
      ],
      transitions: 'sunrise',
      backgroundAnimation: 'shine'
    },
    greetings: {
      en: 'Happy Pongal! May the harvest bring prosperity',
      ta: 'பொங்கல் நல்வாழ்த்துக்கள்! இனிய பொங்கல் திருநாள்',
      te: 'మకర సంక్రాంతి శుభాకాంక్షలు!'
    }
  },
  {
    id: 'ganesh_chaturthi',
    name: 'Ganesh Chaturthi',
    localNames: {
      hi: 'गणेश चतुर्थी',
      mr: 'गणेश चतुर्थी',
      te: 'వినాయక చవితి',
      kn: 'ಗಣೇಶ ಚತುರ್ಥಿ'
    },
    startDate: '2024-09-07',
    endDate: '2024-09-17',
    type: 'religious',
    regions: ['Maharashtra', 'Karnataka', 'Andhra Pradesh', 'Tamil Nadu'],
    theme: {
      primary: '#FF5722',
      secondary: '#FFC107',
      accent: '#4CAF50',
      background: '#FBE9E7',
      backgroundGradient: 'linear-gradient(135deg, #FBE9E7 0%, #FFF9C4 100%)',
      textPrimary: '#BF360C',
      textSecondary: '#D84315',
      decorativeElements: ['ganesha', 'modak', 'drum'],
      iconSet: 'ganesh'
    },
    animations: {
      particles: [
        {
          type: 'flower',
          count: 15,
          size: { min: 20, max: 35 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FF5722', '#FFC107', '#FF1744']
        }
      ],
      transitions: 'celebration',
      backgroundAnimation: 'dance'
    },
    greetings: {
      en: 'Ganpati Bappa Morya! May Lord Ganesha remove all obstacles',
      hi: 'गणपति बप्पा मोरया! विघ्नहर्ता मंगलमूर्ति',
      mr: 'गणपती बाप्पा मोरया! पुढच्या वर्षी लवकर या'
    }
  },
  {
    id: 'baisakhi',
    name: 'Baisakhi',
    localNames: {
      pa: 'ਵਿਸਾਖੀ',
      hi: 'बैसाखी'
    },
    startDate: '2024-04-13',
    endDate: '2024-04-14',
    type: 'regional',
    regions: ['Punjab', 'Haryana'],
    theme: {
      primary: '#FF9800',
      secondary: '#4CAF50',
      accent: '#2196F3',
      background: '#FFF8E1',
      backgroundGradient: 'linear-gradient(135deg, #FFF8E1 0%, #E8F5E9 100%)',
      textPrimary: '#E65100',
      textSecondary: '#F57C00',
      decorativeElements: ['wheat', 'dhol', 'turban'],
      iconSet: 'baisakhi'
    },
    animations: {
      particles: [
        {
          type: 'custom',
          customImage: '/assets/wheat.png',
          count: 15,
          size: { min: 20, max: 35 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FFD700', '#FF9800']
        }
      ],
      transitions: 'harvest',
      backgroundAnimation: 'sway'
    },
    greetings: {
      en: 'Happy Baisakhi! May the harvest bring joy',
      pa: 'ਵਿਸਾਖੀ ਦੀਆਂ ਮੁਬਾਰਕਾਂ! ਰੱਬ ਤੁਹਾਨੂੰ ਖੁਸ਼ੀਆਂ ਦੇਵੇ',
      hi: 'बैसाखी की हार्दिक शुभकामनाएं!'
    }
  },
  {
    id: 'independence_day',
    name: 'Independence Day',
    localNames: {
      hi: 'स्वतंत्रता दिवस',
      bn: 'স্বাধীনতা দিবস',
      ta: 'சுதந்திர தினம்',
      te: 'స్వాతంత్ర్య దినోత్సవం'
    },
    startDate: '2024-08-15',
    endDate: '2024-08-15',
    type: 'national',
    regions: ['all'],
    theme: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808',
      background: '#FFF5E6',
      backgroundGradient: 'linear-gradient(135deg, #FF9933 0%, #FFFFFF 50%, #138808 100%)',
      textPrimary: '#000080',
      textSecondary: '#FF6600',
      decorativeElements: ['flag', 'ashoka_chakra', 'monument'],
      iconSet: 'independence'
    },
    animations: {
      particles: [
        {
          type: 'custom',
          customImage: '/assets/flag.png',
          count: 10,
          size: { min: 30, max: 50 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FF9933', '#FFFFFF', '#138808']
        }
      ],
      transitions: 'patriotic',
      backgroundAnimation: 'wave'
    },
    greetings: {
      en: 'Happy Independence Day! Jai Hind!',
      hi: 'स्वतंत्रता दिवस की शुभकामनाएं! जय हिंद!',
      bn: 'স্বাধীনতা দিবসের শুভেচ্ছা! জয় হিন্দ!',
      ta: 'சுதந்திர தின வாழ்த்துக்கள்! ஜெய் ஹிந்த்!'
    },
    specialOffers: [
      {
        type: 'fee_discount',
        value: 15,
        message: {
          en: '15% off to celebrate freedom!',
          hi: 'स्वतंत्रता का जश्न मनाने के लिए 15% की छूट!'
        }
      }
    ]
  },
  {
    id: 'republic_day',
    name: 'Republic Day',
    localNames: {
      hi: 'गणतंत्र दिवस',
      bn: 'প্রজাতন্ত্র দিবস',
      ta: 'குடியரசு தினம்'
    },
    startDate: '2025-01-26',
    endDate: '2025-01-26',
    type: 'national',
    regions: ['all'],
    theme: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808',
      background: '#FFF5E6',
      backgroundGradient: 'linear-gradient(45deg, #FF9933 0%, #FFFFFF 50%, #138808 100%)',
      textPrimary: '#000080',
      textSecondary: '#FF6600',
      decorativeElements: ['constitution', 'ashoka_pillar', 'flag'],
      iconSet: 'republic'
    },
    animations: {
      particles: [
        {
          type: 'star',
          count: 15,
          size: { min: 15, max: 30 },
          speed: { min: 0.5, max: 1.5 },
          colors: ['#FF9933', '#138808', '#000080']
        }
      ],
      transitions: 'salute',
      backgroundAnimation: 'march'
    },
    greetings: {
      en: 'Happy Republic Day! Celebrating our Constitution',
      hi: 'गणतंत्र दिवस की शुभकामनाएं! हमारे संविधान का सम्मान',
      bn: 'প্রজাতন্ত্র দিবসের শুভেচ্ছা!'
    }
  }
];

// Helper function to get current festival
export const getCurrentFestival = (): Festival | null => {
  const today = new Date();
  const todayStr = today.toISOString().split('T')[0];
  
  return FESTIVALS.find(festival => {
    const start = new Date(festival.startDate);
    const end = new Date(festival.endDate);
    const current = new Date(todayStr);
    
    return current >= start && current <= end;
  }) || null;
};

// Helper function to get upcoming festivals
export const getUpcomingFestivals = (limit: number = 5): Festival[] => {
  const today = new Date();
  const todayStr = today.toISOString().split('T')[0];
  
  return FESTIVALS
    .filter(festival => new Date(festival.startDate) > new Date(todayStr))
    .sort((a, b) => new Date(a.startDate).getTime() - new Date(b.startDate).getTime())
    .slice(0, limit);
};

// Helper function to get festival by region
export const getFestivalsByRegion = (region: string): Festival[] => {
  return FESTIVALS.filter(festival => 
    festival.regions.includes('all') || festival.regions.includes(region)
  );
};

// Helper function to check if today is a festival day
export const isFestivalToday = (): boolean => {
  return getCurrentFestival() !== null;
};

// Helper function to get festival greeting
export const getFestivalGreeting = (festivalId: string, language: string = 'en'): string => {
  const festival = FESTIVALS.find(f => f.id === festivalId);
  if (!festival) return '';
  
  return festival.greetings[language] || festival.greetings['en'] || '';
};

// Helper function to get days until next festival
export const getDaysUntilNextFestival = (): { festival: Festival; days: number } | null => {
  const today = new Date();
  const upcoming = getUpcomingFestivals(1);
  
  if (upcoming.length === 0) return null;
  
  const festival = upcoming[0];
  const festivalDate = new Date(festival.startDate);
  const days = Math.ceil((festivalDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
  
  return { festival, days };
};