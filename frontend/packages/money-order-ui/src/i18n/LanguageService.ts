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

import i18n from './index';
import { SUPPORTED_LANGUAGES } from './index';

// Base translation structure
interface TranslationStructure {
  [key: string]: string | TranslationStructure;
}

// Language metadata
export interface LanguageInfo {
  code: string;
  name: string;
  nativeName: string;
  direction: 'ltr' | 'rtl';
  script: string;
  region: string[];
  speakers: number;
  official: boolean;
}

// Extended language information
export const LANGUAGE_METADATA: Record<string, LanguageInfo> = {
  en: {
    code: 'en',
    name: 'English',
    nativeName: 'English',
    direction: 'ltr',
    script: 'Latin',
    region: ['Pan-India', 'Official'],
    speakers: 125000000,
    official: true
  },
  hi: {
    code: 'hi',
    name: 'Hindi',
    nativeName: 'हिन्दी',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['North India', 'Central India'],
    speakers: 600000000,
    official: true
  },
  bn: {
    code: 'bn',
    name: 'Bengali',
    nativeName: 'বাংলা',
    direction: 'ltr',
    script: 'Bengali',
    region: ['West Bengal', 'Tripura'],
    speakers: 97000000,
    official: true
  },
  te: {
    code: 'te',
    name: 'Telugu',
    nativeName: 'తెలుగు',
    direction: 'ltr',
    script: 'Telugu',
    region: ['Andhra Pradesh', 'Telangana'],
    speakers: 82000000,
    official: true
  },
  mr: {
    code: 'mr',
    name: 'Marathi',
    nativeName: 'मराठी',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Maharashtra', 'Goa'],
    speakers: 83000000,
    official: true
  },
  ta: {
    code: 'ta',
    name: 'Tamil',
    nativeName: 'தமிழ்',
    direction: 'ltr',
    script: 'Tamil',
    region: ['Tamil Nadu', 'Puducherry'],
    speakers: 78000000,
    official: true
  },
  ur: {
    code: 'ur',
    name: 'Urdu',
    nativeName: 'اردو',
    direction: 'rtl',
    script: 'Arabic',
    region: ['North India', 'Jammu & Kashmir'],
    speakers: 52000000,
    official: true
  },
  gu: {
    code: 'gu',
    name: 'Gujarati',
    nativeName: 'ગુજરાતી',
    direction: 'ltr',
    script: 'Gujarati',
    region: ['Gujarat', 'Dadra and Nagar Haveli'],
    speakers: 56000000,
    official: true
  },
  kn: {
    code: 'kn',
    name: 'Kannada',
    nativeName: 'ಕನ್ನಡ',
    direction: 'ltr',
    script: 'Kannada',
    region: ['Karnataka'],
    speakers: 44000000,
    official: true
  },
  or: {
    code: 'or',
    name: 'Odia',
    nativeName: 'ଓଡ଼ିଆ',
    direction: 'ltr',
    script: 'Odia',
    region: ['Odisha'],
    speakers: 38000000,
    official: true
  },
  ml: {
    code: 'ml',
    name: 'Malayalam',
    nativeName: 'മലയാളം',
    direction: 'ltr',
    script: 'Malayalam',
    region: ['Kerala', 'Lakshadweep'],
    speakers: 35000000,
    official: true
  },
  pa: {
    code: 'pa',
    name: 'Punjabi',
    nativeName: 'ਪੰਜਾਬੀ',
    direction: 'ltr',
    script: 'Gurmukhi',
    region: ['Punjab'],
    speakers: 33000000,
    official: true
  },
  as: {
    code: 'as',
    name: 'Assamese',
    nativeName: 'অসমীয়া',
    direction: 'ltr',
    script: 'Bengali',
    region: ['Assam'],
    speakers: 15000000,
    official: true
  },
  mai: {
    code: 'mai',
    name: 'Maithili',
    nativeName: 'मैथिली',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Bihar', 'Jharkhand'],
    speakers: 14000000,
    official: true
  },
  sat: {
    code: 'sat',
    name: 'Santali',
    nativeName: 'ᱥᱟᱱᱛᱟᱲᱤ',
    direction: 'ltr',
    script: 'Ol Chiki',
    region: ['Jharkhand', 'West Bengal', 'Odisha'],
    speakers: 7500000,
    official: true
  },
  ks: {
    code: 'ks',
    name: 'Kashmiri',
    nativeName: 'کٲشُر',
    direction: 'rtl',
    script: 'Arabic',
    region: ['Jammu & Kashmir'],
    speakers: 7000000,
    official: true
  },
  ne: {
    code: 'ne',
    name: 'Nepali',
    nativeName: 'नेपाली',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Sikkim', 'West Bengal', 'Assam'],
    speakers: 3000000,
    official: true
  },
  kok: {
    code: 'kok',
    name: 'Konkani',
    nativeName: 'कोंकणी',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Goa', 'Karnataka', 'Maharashtra'],
    speakers: 2300000,
    official: true
  },
  sd: {
    code: 'sd',
    name: 'Sindhi',
    nativeName: 'सिन्धी',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Gujarat', 'Rajasthan'],
    speakers: 1700000,
    official: true
  },
  doi: {
    code: 'doi',
    name: 'Dogri',
    nativeName: 'डोगरी',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Jammu & Kashmir'],
    speakers: 2600000,
    official: true
  },
  mni: {
    code: 'mni',
    name: 'Manipuri',
    nativeName: 'মৈতৈলোন্',
    direction: 'ltr',
    script: 'Bengali',
    region: ['Manipur'],
    speakers: 1800000,
    official: true
  },
  sa: {
    code: 'sa',
    name: 'Sanskrit',
    nativeName: 'संस्कृतम्',
    direction: 'ltr',
    script: 'Devanagari',
    region: ['Pan-India', 'Religious'],
    speakers: 25000,
    official: true
  }
};

// Number formatting by language
export class LanguageService {
  // Get language-specific number system
  static getNumberSystem(languageCode: string): string[] {
    const numberSystems: Record<string, string[]> = {
      en: ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'],
      hi: ['०', '१', '२', '३', '४', '५', '६', '७', '८', '९'],
      bn: ['০', '১', '২', '৩', '৪', '৫', '৬', '৭', '৮', '৯'],
      te: ['౦', '౧', '౨', '౩', '౪', '౫', '౬', '౭', '౮', '౯'],
      mr: ['०', '१', '२', '३', '४', '५', '६', '७', '८', '९'],
      ta: ['௦', '௧', '௨', '௩', '௪', '௫', '௬', '௭', '௮', '௯'],
      ur: ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹'],
      gu: ['૦', '૧', '૨', '૩', '૪', '૫', '૬', '૭', '૮', '૯'],
      kn: ['೦', '೧', '೨', '೩', '೪', '೫', '೬', '೭', '೮', '೯'],
      or: ['୦', '୧', '୨', '୩', '୪', '୫', '୬', '୭', '୮', '୯'],
      ml: ['൦', '൧', '൨', '൩', '൪', '൫', '൬', '൭', '൮', '൯'],
      pa: ['੦', '੧', '੨', '੩', '੪', '੫', '੬', '੭', '੮', '੯']
    };

    return numberSystems[languageCode] || numberSystems.en;
  }

  // Convert number to localized string
  static localizeNumber(number: number | string, languageCode: string): string {
    const numStr = number.toString();
    const numberSystem = this.getNumberSystem(languageCode);
    
    if (numberSystem === this.getNumberSystem('en')) {
      return numStr;
    }

    return numStr.split('').map(digit => {
      const index = parseInt(digit);
      return isNaN(index) ? digit : numberSystem[index];
    }).join('');
  }

  // Get language-specific currency symbol
  static getCurrencySymbol(languageCode: string): string {
    // Most Indian languages use ₹ for Rupee
    return '₹';
  }

  // Format amount with proper separators (Indian numbering system)
  static formatIndianNumber(amount: number, languageCode: string): string {
    const parts = amount.toString().split('.');
    let integerPart = parts[0];
    const decimalPart = parts[1];

    // Indian numbering system: 1,00,00,000 (1 crore)
    const lastThree = integerPart.substring(integerPart.length - 3);
    const otherNumbers = integerPart.substring(0, integerPart.length - 3);
    
    if (otherNumbers !== '') {
      integerPart = otherNumbers.replace(/\B(?=(\d{2})+(?!\d))/g, ',') + ',' + lastThree;
    } else {
      integerPart = lastThree;
    }

    const formattedNumber = decimalPart ? `${integerPart}.${decimalPart}` : integerPart;
    
    // Convert to localized digits
    return this.localizeNumber(formattedNumber, languageCode);
  }

  // Get greeting based on time of day
  static getGreeting(languageCode: string): string {
    const hour = new Date().getHours();
    const greetings: Record<string, { morning: string; afternoon: string; evening: string; night: string }> = {
      en: {
        morning: 'Good Morning',
        afternoon: 'Good Afternoon',
        evening: 'Good Evening',
        night: 'Good Night'
      },
      hi: {
        morning: 'शुभ प्रभात',
        afternoon: 'नमस्ते',
        evening: 'शुभ संध्या',
        night: 'शुभ रात्रि'
      },
      bn: {
        morning: 'সুপ্রভাত',
        afternoon: 'শুভ অপরাহ্ণ',
        evening: 'শুভ সন্ধ্যা',
        night: 'শুভ রাত্রি'
      }
      // Add more languages as needed
    };

    const timeGreetings = greetings[languageCode] || greetings.en;
    
    if (hour < 12) return timeGreetings.morning;
    if (hour < 17) return timeGreetings.afternoon;
    if (hour < 21) return timeGreetings.evening;
    return timeGreetings.night;
  }

  // Get cultural phrases
  static getCulturalPhrase(languageCode: string, type: 'thanks' | 'welcome' | 'blessing'): string {
    const phrases: Record<string, Record<string, string>> = {
      en: {
        thanks: 'Thank you',
        welcome: 'You\'re welcome',
        blessing: 'May you prosper'
      },
      hi: {
        thanks: 'धन्यवाद 🙏',
        welcome: 'आपका स्वागत है',
        blessing: 'जय हिंद! आपकी यात्रा शुभ हो'
      },
      bn: {
        thanks: 'ধন্যবাদ 🙏',
        welcome: 'আপনাকে স্বাগতম',
        blessing: 'জয় বাংলা! আপনার যাত্রা শুভ হোক'
      },
      ta: {
        thanks: 'நன்றி 🙏',
        welcome: 'உங்களை வரவேற்கிறோம்',
        blessing: 'வாழ்க தமிழ்! உங்கள் பயணம் சிறப்பாக அமைய வாழ்த்துகள்'
      }
      // Add more languages
    };

    const languagePhrases = phrases[languageCode] || phrases.en;
    return languagePhrases[type] || phrases.en[type];
  }

  // Validate phone number format for different regions
  static validatePhoneNumber(phoneNumber: string, languageCode: string): boolean {
    // Remove all non-digit characters
    const digits = phoneNumber.replace(/\D/g, '');
    
    // Indian phone numbers: 10 digits starting with 6-9
    const indianPhoneRegex = /^[6-9]\d{9}$/;
    
    // Check if it's a valid Indian number (without country code)
    if (indianPhoneRegex.test(digits)) {
      return true;
    }
    
    // Check if it includes country code +91
    if (digits.startsWith('91') && digits.length === 12) {
      return indianPhoneRegex.test(digits.substring(2));
    }
    
    return false;
  }

  // Get language-specific date format
  static getDateFormat(languageCode: string): string {
    const formats: Record<string, string> = {
      en: 'DD/MM/YYYY',
      hi: 'DD/MM/YYYY',
      bn: 'DD/MM/YYYY',
      ur: 'DD/MM/YYYY', // Urdu uses same format but RTL
      // Add more if different
    };
    
    return formats[languageCode] || formats.en;
  }

  // Generate placeholder translations for missing languages
  static generatePlaceholderTranslations(baseTranslations: any): Record<string, any> {
    const placeholderTranslations: Record<string, any> = {};
    
    // For languages without full translations, create basic placeholders
    const languagesWithoutFullTranslations = ['te', 'mr', 'ta', 'ur', 'gu', 'kn', 'or', 'ml', 'pa', 'as', 'mai', 'sat', 'ks', 'ne', 'kok', 'sd', 'doi', 'mni', 'sa'];
    
    languagesWithoutFullTranslations.forEach(lang => {
      placeholderTranslations[lang] = this.createBasicTranslations(lang, baseTranslations);
    });
    
    return placeholderTranslations;
  }

  // Create basic translations with key terms translated
  private static createBasicTranslations(languageCode: string, baseTranslations: any): any {
    // Key terms in different languages
    const keyTerms: Record<string, Record<string, string>> = {
      te: {
        'Money Order': 'మనీ ఆర్డర్',
        'Send': 'పంపు',
        'Receive': 'స్వీకరించు',
        'Amount': 'మొత్తం',
        'Home': 'హోమ్'
      },
      ta: {
        'Money Order': 'பணக்கட்டளை',
        'Send': 'அனுப்பு',
        'Receive': 'பெறு',
        'Amount': 'தொகை',
        'Home': 'முகப்பு'
      },
      // Add more key terms for other languages
    };

    // Clone base translations and replace key terms
    const translations = JSON.parse(JSON.stringify(baseTranslations));
    
    // Simple replacement logic - in production, use proper translation service
    if (keyTerms[languageCode]) {
      Object.entries(keyTerms[languageCode]).forEach(([english, translated]) => {
        this.replaceInObject(translations, english, translated);
      });
    }
    
    return translations;
  }

  // Helper to replace text in nested objects
  private static replaceInObject(obj: any, searchText: string, replaceText: string): void {
    Object.keys(obj).forEach(key => {
      if (typeof obj[key] === 'string') {
        obj[key] = obj[key].replace(new RegExp(searchText, 'g'), replaceText);
      } else if (typeof obj[key] === 'object' && obj[key] !== null) {
        this.replaceInObject(obj[key], searchText, replaceText);
      }
    });
  }
}