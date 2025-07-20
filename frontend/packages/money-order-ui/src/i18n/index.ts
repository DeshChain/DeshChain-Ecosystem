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

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// Import all language files
import en from './locales/en.json';
import hi from './locales/hi.json';
import bn from './locales/bn.json';
import te from './locales/te.json';
import mr from './locales/mr.json';
import ta from './locales/ta.json';
import ur from './locales/ur.json';
import gu from './locales/gu.json';
import kn from './locales/kn.json';
import or from './locales/or.json';
import ml from './locales/ml.json';
import pa from './locales/pa.json';
import as from './locales/as.json';
import mai from './locales/mai.json';
import sat from './locales/sat.json';
import ks from './locales/ks.json';
import ne from './locales/ne.json';
import kok from './locales/kok.json';
import sd from './locales/sd.json';
import doi from './locales/doi.json';
import mni from './locales/mni.json';
import sa from './locales/sa.json';

// Language configuration
export const SUPPORTED_LANGUAGES = [
  { code: 'en', name: 'English', nativeName: 'English', direction: 'ltr' },
  { code: 'hi', name: 'Hindi', nativeName: 'हिन्दी', direction: 'ltr' },
  { code: 'bn', name: 'Bengali', nativeName: 'বাংলা', direction: 'ltr' },
  { code: 'te', name: 'Telugu', nativeName: 'తెలుగు', direction: 'ltr' },
  { code: 'mr', name: 'Marathi', nativeName: 'मराठी', direction: 'ltr' },
  { code: 'ta', name: 'Tamil', nativeName: 'தமிழ்', direction: 'ltr' },
  { code: 'ur', name: 'Urdu', nativeName: 'اردو', direction: 'rtl' },
  { code: 'gu', name: 'Gujarati', nativeName: 'ગુજરાતી', direction: 'ltr' },
  { code: 'kn', name: 'Kannada', nativeName: 'ಕನ್ನಡ', direction: 'ltr' },
  { code: 'or', name: 'Odia', nativeName: 'ଓଡ଼ିଆ', direction: 'ltr' },
  { code: 'ml', name: 'Malayalam', nativeName: 'മലയാളം', direction: 'ltr' },
  { code: 'pa', name: 'Punjabi', nativeName: 'ਪੰਜਾਬੀ', direction: 'ltr' },
  { code: 'as', name: 'Assamese', nativeName: 'অসমীয়া', direction: 'ltr' },
  { code: 'mai', name: 'Maithili', nativeName: 'मैथिली', direction: 'ltr' },
  { code: 'sat', name: 'Santali', nativeName: 'ᱥᱟᱱᱛᱟᱲᱤ', direction: 'ltr' },
  { code: 'ks', name: 'Kashmiri', nativeName: 'کٲشُر', direction: 'rtl' },
  { code: 'ne', name: 'Nepali', nativeName: 'नेपाली', direction: 'ltr' },
  { code: 'kok', name: 'Konkani', nativeName: 'कोंकणी', direction: 'ltr' },
  { code: 'sd', name: 'Sindhi', nativeName: 'सिन्धी', direction: 'ltr' },
  { code: 'doi', name: 'Dogri', nativeName: 'डोगरी', direction: 'ltr' },
  { code: 'mni', name: 'Manipuri', nativeName: 'মৈতৈলোন্', direction: 'ltr' },
  { code: 'sa', name: 'Sanskrit', nativeName: 'संस्कृतम्', direction: 'ltr' }
];

// Resources object with all translations
const resources = {
  en: { translation: en },
  hi: { translation: hi },
  bn: { translation: bn },
  te: { translation: te },
  mr: { translation: mr },
  ta: { translation: ta },
  ur: { translation: ur },
  gu: { translation: gu },
  kn: { translation: kn },
  or: { translation: or },
  ml: { translation: ml },
  pa: { translation: pa },
  as: { translation: as },
  mai: { translation: mai },
  sat: { translation: sat },
  ks: { translation: ks },
  ne: { translation: ne },
  kok: { translation: kok },
  sd: { translation: sd },
  doi: { translation: doi },
  mni: { translation: mni },
  sa: { translation: sa }
};

// Initialize i18n
i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    debug: false,
    
    interpolation: {
      escapeValue: false // React already escapes values
    },

    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage']
    },

    react: {
      useSuspense: true
    }
  });

// Language helper functions
export const getCurrentLanguage = (): string => {
  return i18n.language;
};

export const changeLanguage = async (languageCode: string): Promise<void> => {
  await i18n.changeLanguage(languageCode);
  
  // Update document direction for RTL languages
  const language = SUPPORTED_LANGUAGES.find(lang => lang.code === languageCode);
  if (language) {
    document.documentElement.dir = language.direction;
    document.documentElement.lang = languageCode;
  }
  
  // Save preference
  localStorage.setItem('deshchain-language', languageCode);
};

export const getLanguageDirection = (languageCode?: string): 'ltr' | 'rtl' => {
  const code = languageCode || getCurrentLanguage();
  const language = SUPPORTED_LANGUAGES.find(lang => lang.code === code);
  return language?.direction || 'ltr';
};

export const formatCurrency = (amount: number, currency: string = 'INR'): string => {
  const locale = getCurrentLanguage() === 'en' ? 'en-IN' : `${getCurrentLanguage()}-IN`;
  
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  }).format(amount);
};

export const formatNumber = (number: number, options?: Intl.NumberFormatOptions): string => {
  const locale = getCurrentLanguage() === 'en' ? 'en-IN' : `${getCurrentLanguage()}-IN`;
  return new Intl.NumberFormat(locale, options).format(number);
};

export const formatDate = (date: Date | string, options?: Intl.DateTimeFormatOptions): string => {
  const locale = getCurrentLanguage() === 'en' ? 'en-IN' : `${getCurrentLanguage()}-IN`;
  const dateObj = typeof date === 'string' ? new Date(date) : date;
  
  const defaultOptions: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    ...options
  };
  
  return new Intl.DateTimeFormat(locale, defaultOptions).format(dateObj);
};

// Script-specific font loading
export const loadScriptFonts = async (languageCode: string): Promise<void> => {
  const fontMap: Record<string, string[]> = {
    hi: ['Noto Sans Devanagari'],
    bn: ['Noto Sans Bengali'],
    te: ['Noto Sans Telugu'],
    mr: ['Noto Sans Devanagari'],
    ta: ['Noto Sans Tamil'],
    ur: ['Noto Nastaliq Urdu', 'Noto Sans Arabic'],
    gu: ['Noto Sans Gujarati'],
    kn: ['Noto Sans Kannada'],
    or: ['Noto Sans Oriya'],
    ml: ['Noto Sans Malayalam'],
    pa: ['Noto Sans Gurmukhi'],
    as: ['Noto Sans Bengali'],
    mai: ['Noto Sans Devanagari'],
    sat: ['Noto Sans Ol Chiki'],
    ks: ['Noto Nastaliq Urdu', 'Noto Sans Arabic'],
    ne: ['Noto Sans Devanagari'],
    kok: ['Noto Sans Devanagari'],
    sd: ['Noto Sans Devanagari'],
    doi: ['Noto Sans Devanagari'],
    mni: ['Noto Sans Bengali'],
    sa: ['Noto Sans Devanagari']
  };

  const fonts = fontMap[languageCode];
  if (fonts) {
    // Dynamically load fonts
    const WebFont = (await import('webfontloader')).default;
    
    WebFont.load({
      google: {
        families: fonts
      }
    });
  }
};

// Automatically load fonts when language changes
i18n.on('languageChanged', (lng) => {
  loadScriptFonts(lng);
});

export default i18n;