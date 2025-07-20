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

import React, { createContext, ReactNode, useState, useEffect } from 'react';
import { CulturalQuoteData, FestivalInfo, PatriotismScore, UIPreferences } from '../types';

interface CulturalContextType {
  currentLanguage: string;
  setCurrentLanguage: (language: string) => void;
  currentQuote: CulturalQuoteData | null;
  setCurrentQuote: (quote: CulturalQuoteData | null) => void;
  currentFestival: FestivalInfo | null;
  setCurrentFestival: (festival: FestivalInfo | null) => void;
  patriotismScore: PatriotismScore | null;
  setPatriotismScore: (score: PatriotismScore | null) => void;
  culturalPreferences: UIPreferences;
  setCulturalPreferences: (preferences: UIPreferences) => void;
  enabledFeatures: Record<string, boolean>;
  setEnabledFeatures: (features: Record<string, boolean>) => void;
}

export const CulturalContext = createContext<CulturalContextType | null>(null);

interface CulturalProviderProps {
  children: ReactNode;
  initialLanguage?: string;
  enabledFeatures?: Record<string, boolean>;
  initialPreferences?: Partial<UIPreferences>;
}

const defaultPreferences: UIPreferences = {
  language: 'en',
  theme: 'prosperity',
  festivalThemes: true,
  culturalQuotes: true,
  patriotismFeatures: true,
  animations: true,
  density: 'comfortable',
  notifications: {
    transactions: true,
    festivals: true,
    rewards: true
  }
};

const defaultEnabledFeatures = {
  cultural: true,
  festivals: true,
  patriotism: true,
  quotes: true,
  themes: true,
  animations: true
};

export const CulturalProvider: React.FC<CulturalProviderProps> = ({
  children,
  initialLanguage = 'en',
  enabledFeatures: userEnabledFeatures = {},
  initialPreferences = {}
}) => {
  const [currentLanguage, setCurrentLanguage] = useState(initialLanguage);
  const [currentQuote, setCurrentQuote] = useState<CulturalQuoteData | null>(null);
  const [currentFestival, setCurrentFestival] = useState<FestivalInfo | null>(null);
  const [patriotismScore, setPatriotismScore] = useState<PatriotismScore | null>(null);
  const [culturalPreferences, setCulturalPreferences] = useState<UIPreferences>({
    ...defaultPreferences,
    ...initialPreferences,
    language: initialLanguage
  });
  const [enabledFeatures, setEnabledFeatures] = useState({
    ...defaultEnabledFeatures,
    ...userEnabledFeatures
  });

  // Load preferences from localStorage on mount
  useEffect(() => {
    try {
      const savedPreferences = localStorage.getItem('deshchain-cultural-preferences');
      if (savedPreferences) {
        const parsed = JSON.parse(savedPreferences);
        setCulturalPreferences(prev => ({ ...prev, ...parsed }));
        setCurrentLanguage(parsed.language || initialLanguage);
      }

      const savedFeatures = localStorage.getItem('deshchain-enabled-features');
      if (savedFeatures) {
        const parsed = JSON.parse(savedFeatures);
        setEnabledFeatures(prev => ({ ...prev, ...parsed }));
      }
    } catch (error) {
      console.error('Error loading cultural preferences:', error);
    }
  }, [initialLanguage]);

  // Save preferences to localStorage when they change
  useEffect(() => {
    try {
      localStorage.setItem('deshchain-cultural-preferences', JSON.stringify(culturalPreferences));
    } catch (error) {
      console.error('Error saving cultural preferences:', error);
    }
  }, [culturalPreferences]);

  useEffect(() => {
    try {
      localStorage.setItem('deshchain-enabled-features', JSON.stringify(enabledFeatures));
    } catch (error) {
      console.error('Error saving enabled features:', error);
    }
  }, [enabledFeatures]);

  // Update language in preferences when currentLanguage changes
  useEffect(() => {
    setCulturalPreferences(prev => ({
      ...prev,
      language: currentLanguage
    }));
  }, [currentLanguage]);

  const contextValue: CulturalContextType = {
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
  };

  return (
    <CulturalContext.Provider value={contextValue}>
      {children}
    </CulturalContext.Provider>
  );
};