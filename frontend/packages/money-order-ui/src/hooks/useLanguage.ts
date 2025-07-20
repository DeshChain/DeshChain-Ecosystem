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

import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  getCurrentLanguage, 
  changeLanguage, 
  getLanguageDirection,
  formatCurrency,
  formatNumber,
  formatDate,
  loadScriptFonts,
  SUPPORTED_LANGUAGES
} from '../i18n';
import { LanguageService, LANGUAGE_METADATA } from '../i18n/LanguageService';

interface UseLanguageReturn {
  // Current language
  currentLanguage: string;
  currentLanguageInfo: typeof SUPPORTED_LANGUAGES[0] | undefined;
  currentLanguageMetadata: typeof LANGUAGE_METADATA[string] | undefined;
  
  // Language operations
  changeLanguage: (languageCode: string) => Promise<void>;
  isRTL: boolean;
  
  // Formatting functions
  formatCurrency: (amount: number, currency?: string) => string;
  formatNumber: (number: number, options?: Intl.NumberFormatOptions) => string;
  formatDate: (date: Date | string, options?: Intl.DateTimeFormatOptions) => string;
  formatIndianNumber: (amount: number) => string;
  localizeNumber: (number: number | string) => string;
  
  // Cultural helpers
  getGreeting: () => string;
  getCulturalPhrase: (type: 'thanks' | 'welcome' | 'blessing') => string;
  validatePhoneNumber: (phoneNumber: string) => boolean;
  
  // UI helpers
  supportedLanguages: typeof SUPPORTED_LANGUAGES;
  getLanguageProgress: () => number;
}

export const useLanguage = (): UseLanguageReturn => {
  const { i18n } = useTranslation();
  const [currentLanguage, setCurrentLanguage] = useState(getCurrentLanguage());
  const [isRTL, setIsRTL] = useState(false);

  // Get current language info
  const currentLanguageInfo = SUPPORTED_LANGUAGES.find(lang => lang.code === currentLanguage);
  const currentLanguageMetadata = LANGUAGE_METADATA[currentLanguage];

  // Update language state when i18n changes
  useEffect(() => {
    const handleLanguageChange = (lng: string) => {
      setCurrentLanguage(lng);
      setIsRTL(getLanguageDirection(lng) === 'rtl');
      
      // Update document direction
      document.documentElement.dir = getLanguageDirection(lng);
      document.documentElement.lang = lng;
      
      // Load appropriate fonts
      loadScriptFonts(lng);
    };

    i18n.on('languageChanged', handleLanguageChange);
    
    // Initial setup
    handleLanguageChange(getCurrentLanguage());

    return () => {
      i18n.off('languageChanged', handleLanguageChange);
    };
  }, [i18n]);

  // Change language handler
  const handleChangeLanguage = useCallback(async (languageCode: string) => {
    await changeLanguage(languageCode);
    
    // Save to localStorage for persistence
    localStorage.setItem('deshchain-preferred-language', languageCode);
    
    // Emit custom event for other components
    window.dispatchEvent(new CustomEvent('deshchain-language-changed', {
      detail: { languageCode }
    }));
  }, []);

  // Format currency with Indian numbering
  const handleFormatCurrency = useCallback((amount: number, currency: string = 'INR') => {
    return formatCurrency(amount, currency);
  }, []);

  // Format number with localization
  const handleFormatNumber = useCallback((number: number, options?: Intl.NumberFormatOptions) => {
    return formatNumber(number, options);
  }, []);

  // Format date with localization
  const handleFormatDate = useCallback((date: Date | string, options?: Intl.DateTimeFormatOptions) => {
    return formatDate(date, options);
  }, []);

  // Format with Indian numbering system
  const handleFormatIndianNumber = useCallback((amount: number) => {
    return LanguageService.formatIndianNumber(amount, currentLanguage);
  }, [currentLanguage]);

  // Localize number digits
  const handleLocalizeNumber = useCallback((number: number | string) => {
    return LanguageService.localizeNumber(number, currentLanguage);
  }, [currentLanguage]);

  // Get time-based greeting
  const handleGetGreeting = useCallback(() => {
    return LanguageService.getGreeting(currentLanguage);
  }, [currentLanguage]);

  // Get cultural phrases
  const handleGetCulturalPhrase = useCallback((type: 'thanks' | 'welcome' | 'blessing') => {
    return LanguageService.getCulturalPhrase(currentLanguage, type);
  }, [currentLanguage]);

  // Validate phone number
  const handleValidatePhoneNumber = useCallback((phoneNumber: string) => {
    return LanguageService.validatePhoneNumber(phoneNumber, currentLanguage);
  }, [currentLanguage]);

  // Calculate language translation progress (mock - would connect to real data)
  const getLanguageProgress = useCallback(() => {
    // In a real app, this would check actual translation completion
    const translationProgress: Record<string, number> = {
      en: 100,
      hi: 100,
      bn: 100,
      te: 40,
      mr: 40,
      ta: 40,
      ur: 40,
      gu: 40,
      kn: 40,
      // Others at 20% (placeholder)
    };
    
    return translationProgress[currentLanguage] || 20;
  }, [currentLanguage]);

  return {
    // Current language
    currentLanguage,
    currentLanguageInfo,
    currentLanguageMetadata,
    
    // Language operations
    changeLanguage: handleChangeLanguage,
    isRTL,
    
    // Formatting functions
    formatCurrency: handleFormatCurrency,
    formatNumber: handleFormatNumber,
    formatDate: handleFormatDate,
    formatIndianNumber: handleFormatIndianNumber,
    localizeNumber: handleLocalizeNumber,
    
    // Cultural helpers
    getGreeting: handleGetGreeting,
    getCulturalPhrase: handleGetCulturalPhrase,
    validatePhoneNumber: handleValidatePhoneNumber,
    
    // UI helpers
    supportedLanguages: SUPPORTED_LANGUAGES,
    getLanguageProgress
  };
};

// Hook for RTL support
export const useRTL = () => {
  const { isRTL } = useLanguage();
  
  return {
    isRTL,
    textAlign: isRTL ? 'right' : 'left',
    direction: isRTL ? 'rtl' : 'ltr',
    // MUI theme direction helpers
    theme: {
      direction: isRTL ? 'rtl' : 'ltr'
    },
    // Utility classes
    classes: {
      alignStart: isRTL ? 'text-right' : 'text-left',
      alignEnd: isRTL ? 'text-left' : 'text-right',
      marginStart: isRTL ? 'mr' : 'ml',
      marginEnd: isRTL ? 'ml' : 'mr',
      paddingStart: isRTL ? 'pr' : 'pl',
      paddingEnd: isRTL ? 'pl' : 'pr'
    }
  };
};

// Hook for number localization
export const useLocalizedNumbers = () => {
  const { localizeNumber, formatIndianNumber, currentLanguage } = useLanguage();
  
  const localizeInput = useCallback((value: string): string => {
    // Convert localized digits back to Western for processing
    const numberSystem = LanguageService.getNumberSystem(currentLanguage);
    const westernDigits = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'];
    
    let normalized = value;
    numberSystem.forEach((localDigit, index) => {
      normalized = normalized.replace(new RegExp(localDigit, 'g'), westernDigits[index]);
    });
    
    return normalized;
  }, [currentLanguage]);
  
  const formatForDisplay = useCallback((value: number | string): string => {
    const numValue = typeof value === 'string' ? parseFloat(value) : value;
    if (isNaN(numValue)) return '';
    
    return formatIndianNumber(numValue);
  }, [formatIndianNumber]);
  
  return {
    localizeNumber,
    formatIndianNumber,
    localizeInput,
    formatForDisplay
  };
};