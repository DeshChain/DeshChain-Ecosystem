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

import { useState, useCallback, useContext, useEffect } from 'react';
import { CulturalQuoteData, FestivalInfo, PatriotismScore, LanguageOption } from '../types';
import { CulturalContext } from '../providers/CulturalProvider';
import { SUPPORTED_LANGUAGES, FESTIVALS } from '../constants';

export const useCulturalContext = () => {
  const context = useContext(CulturalContext);
  
  if (!context) {
    throw new Error('useCulturalContext must be used within CulturalProvider');
  }

  const {
    currentLanguage,
    setCurrentLanguage,
    currentQuote,
    setCurrentQuote,
    currentFestival,
    setCurrentFestival,
    patriotismScore,
    setPatriotismScore,
    culturalPreferences,
    setCulturalPreferences,
    enabledFeatures,
    setEnabledFeatures
  } = context;

  const [isLoadingQuote, setIsLoadingQuote] = useState(false);
  const [isLoadingFestival, setIsLoadingFestival] = useState(false);

  // Get contextual quote based on category and occasion
  const getContextualQuote = useCallback(async (
    category?: string,
    occasion?: string,
    language?: string
  ): Promise<CulturalQuoteData | null> => {
    setIsLoadingQuote(true);
    
    try {
      const queryParams = new URLSearchParams();
      if (category) queryParams.append('category', category);
      if (occasion) queryParams.append('occasion', occasion);
      if (language || currentLanguage) queryParams.append('language', language || currentLanguage);

      const response = await fetch(`/api/cultural/quotes?${queryParams}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch cultural quote');
      }

      const result = await response.json();
      const quote = result.quotes[0]; // Get the first quote from the response
      
      if (quote) {
        setCurrentQuote(quote);
        return quote;
      }
      
      return null;
    } catch (error) {
      console.error('Error fetching cultural quote:', error);
      return null;
    } finally {
      setIsLoadingQuote(false);
    }
  }, [currentLanguage, setCurrentQuote]);

  // Get random quote from a specific category
  const getRandomQuote = useCallback(async (category: string): Promise<CulturalQuoteData | null> => {
    return getContextualQuote(category);
  }, [getContextualQuote]);

  // Get current active festivals
  const getActiveFestivals = useCallback(async (): Promise<FestivalInfo[]> => {
    setIsLoadingFestival(true);
    
    try {
      const response = await fetch('/api/cultural/festivals?current=true');
      
      if (!response.ok) {
        throw new Error('Failed to fetch festivals');
      }

      const result = await response.json();
      
      if (result.current_active && result.current_active.length > 0) {
        setCurrentFestival(result.festivals.find((f: FestivalInfo) => f.active));
      }
      
      return result.festivals;
    } catch (error) {
      console.error('Error fetching festivals:', error);
      // Fallback to local festival data
      const today = new Date();
      const activeFestivals = FESTIVALS.filter(festival => {
        const startDate = new Date(festival.startDate);
        const endDate = new Date(festival.endDate);
        return today >= startDate && today <= endDate;
      });
      
      if (activeFestivals.length > 0) {
        setCurrentFestival(activeFestivals[0]);
      }
      
      return activeFestivals;
    } finally {
      setIsLoadingFestival(false);
    }
  }, [setCurrentFestival]);

  // Update patriotism score
  const updatePatriotismScore = useCallback(async (userId: string): Promise<PatriotismScore | null> => {
    try {
      const response = await fetch(`/api/cultural/patriotism/${userId}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch patriotism score');
      }

      const score = await response.json();
      setPatriotismScore(score);
      return score;
    } catch (error) {
      console.error('Error fetching patriotism score:', error);
      return null;
    }
  }, [setPatriotismScore]);

  // Change language
  const changeLanguage = useCallback((languageCode: string) => {
    const language = SUPPORTED_LANGUAGES.find(lang => lang.code === languageCode);
    if (language && language.supported) {
      setCurrentLanguage(languageCode);
      // Refresh current quote in new language
      if (currentQuote) {
        getContextualQuote(currentQuote.category, currentQuote.occasion, languageCode);
      }
    }
  }, [setCurrentLanguage, currentQuote, getContextualQuote]);

  // Get language native name
  const getLanguageNativeName = useCallback((languageCode: string): string => {
    const language = SUPPORTED_LANGUAGES.find(lang => lang.code === languageCode);
    return language?.nativeName || languageCode;
  }, []);

  // Get supported languages for a region
  const getLanguagesByRegion = useCallback((region: string): LanguageOption[] => {
    return SUPPORTED_LANGUAGES.filter(lang => 
      lang.region.toLowerCase().includes(region.toLowerCase()) || 
      lang.region === 'Pan India'
    );
  }, []);

  // Check if cultural features are enabled
  const isFeatureEnabled = useCallback((feature: string): boolean => {
    return enabledFeatures[feature] === true;
  }, [enabledFeatures]);

  const isCulturalFeaturesEnabled = isFeatureEnabled('cultural');
  const isFestivalFeaturesEnabled = isFeatureEnabled('festivals');
  const isPatriotismEnabled = isFeatureEnabled('patriotism');

  // Auto-fetch active festivals on mount
  useEffect(() => {
    if (isFestivalFeaturesEnabled) {
      getActiveFestivals();
    }
  }, [isFestivalFeaturesEnabled, getActiveFestivals]);

  // Auto-fetch contextual quote when language changes
  useEffect(() => {
    if (isCulturalFeaturesEnabled && culturalPreferences.includeQuotes) {
      getContextualQuote('general', 'general', currentLanguage);
    }
  }, [currentLanguage, isCulturalFeaturesEnabled, culturalPreferences.includeQuotes, getContextualQuote]);

  return {
    // Current state
    currentLanguage,
    currentQuote,
    currentFestival,
    patriotismScore,
    culturalPreferences,
    enabledFeatures,
    
    // Loading states
    isLoadingQuote,
    isLoadingFestival,
    
    // Feature flags
    isCulturalFeaturesEnabled,
    isFestivalFeaturesEnabled,
    isPatriotismEnabled,
    
    // Quote operations
    getContextualQuote,
    getRandomQuote,
    
    // Festival operations
    getActiveFestivals,
    
    // Patriotism operations
    updatePatriotismScore,
    
    // Language operations
    changeLanguage,
    getLanguageNativeName,
    getLanguagesByRegion,
    
    // Preferences
    setCulturalPreferences,
    setEnabledFeatures,
    
    // Utilities
    isFeatureEnabled,
    
    // Constants
    supportedLanguages: SUPPORTED_LANGUAGES,
    availableFestivals: FESTIVALS
  };
};