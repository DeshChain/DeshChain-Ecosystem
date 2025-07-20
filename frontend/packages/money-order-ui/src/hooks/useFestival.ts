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
import { 
  Festival, 
  getCurrentFestival, 
  getUpcomingFestivals,
  getFestivalsByRegion,
  isFestivalToday,
  getFestivalGreeting,
  getDaysUntilNextFestival,
  FESTIVALS
} from '../themes/festivals';
import { useFestivalTheme } from '../themes/FestivalThemeProvider';
import { useLanguage } from './useLanguage';

interface UseFestivalReturn {
  // Current festival
  currentFestival: Festival | null;
  isFestivalActive: boolean;
  
  // Festival data
  upcomingFestivals: Festival[];
  regionalFestivals: Festival[];
  allFestivals: Festival[];
  
  // Festival utilities
  getFestivalGreeting: (festivalId?: string) => string;
  getDaysUntilNextFestival: () => { festival: Festival; days: number } | null;
  
  // Theme controls
  festivalThemeEnabled: boolean;
  toggleFestivalTheme: (enabled: boolean) => void;
  
  // Festival bonuses
  calculateFestivalBonus: (amount: number) => {
    discount: number;
    bonus: number;
    totalBenefit: number;
    message: string;
  };
  
  // User preferences
  userRegion: string;
  setUserRegion: (region: string) => void;
  favoritesFestivals: string[];
  toggleFavoriteFestival: (festivalId: string) => void;
}

export const useFestival = (): UseFestivalReturn => {
  const { currentFestival: themeFestival, isEnabled, toggleFestivalTheme } = useFestivalTheme();
  const { currentLanguage } = useLanguage();
  
  const [currentFestival, setCurrentFestival] = useState<Festival | null>(null);
  const [upcomingFestivals, setUpcomingFestivals] = useState<Festival[]>([]);
  const [userRegion, setUserRegion] = useState<string>('all');
  const [favoritesFestivals, setFavoritesFestivals] = useState<string[]>([]);

  // Load user preferences
  useEffect(() => {
    const savedRegion = localStorage.getItem('deshchain-user-region');
    if (savedRegion) {
      setUserRegion(savedRegion);
    }
    
    const savedFavorites = localStorage.getItem('deshchain-favorite-festivals');
    if (savedFavorites) {
      setFavoritesFestivals(JSON.parse(savedFavorites));
    }
  }, []);

  // Update current festival
  useEffect(() => {
    const updateFestival = () => {
      const festival = getCurrentFestival();
      setCurrentFestival(festival);
      
      const upcoming = getUpcomingFestivals(10);
      setUpcomingFestivals(upcoming);
    };

    updateFestival();
    
    // Check every hour
    const interval = setInterval(updateFestival, 60 * 60 * 1000);
    
    return () => clearInterval(interval);
  }, []);

  // Save region preference
  const handleSetUserRegion = useCallback((region: string) => {
    setUserRegion(region);
    localStorage.setItem('deshchain-user-region', region);
  }, []);

  // Toggle favorite festival
  const toggleFavoriteFestival = useCallback((festivalId: string) => {
    setFavoritesFestivals(prev => {
      const newFavorites = prev.includes(festivalId)
        ? prev.filter(id => id !== festivalId)
        : [...prev, festivalId];
      
      localStorage.setItem('deshchain-favorite-festivals', JSON.stringify(newFavorites));
      return newFavorites;
    });
  }, []);

  // Get festival greeting
  const handleGetFestivalGreeting = useCallback((festivalId?: string) => {
    const id = festivalId || currentFestival?.id;
    if (!id) return '';
    
    return getFestivalGreeting(id, currentLanguage);
  }, [currentFestival, currentLanguage]);

  // Calculate festival bonus
  const calculateFestivalBonus = useCallback((amount: number) => {
    if (!currentFestival || !currentFestival.specialOffers) {
      return {
        discount: 0,
        bonus: 0,
        totalBenefit: 0,
        message: ''
      };
    }

    const feeDiscount = currentFestival.specialOffers.find(offer => offer.type === 'fee_discount');
    const bonusAmount = currentFestival.specialOffers.find(offer => offer.type === 'bonus_amount');
    
    const baseFee = amount * 0.01; // 1% base fee
    const discount = feeDiscount 
      ? baseFee * ((feeDiscount.value as number) / 100)
      : 0;
    const bonus = bonusAmount ? (bonusAmount.value as number) : 0;
    const totalBenefit = discount + bonus;
    
    const message = feeDiscount?.message[currentLanguage] || 
                   feeDiscount?.message.en || 
                   'Festival offer active!';

    return {
      discount,
      bonus,
      totalBenefit,
      message
    };
  }, [currentFestival, currentLanguage]);

  // Get regional festivals
  const regionalFestivals = getFestivalsByRegion(userRegion);

  return {
    // Current festival
    currentFestival: themeFestival || currentFestival,
    isFestivalActive: isFestivalToday(),
    
    // Festival data
    upcomingFestivals,
    regionalFestivals,
    allFestivals: FESTIVALS,
    
    // Festival utilities
    getFestivalGreeting: handleGetFestivalGreeting,
    getDaysUntilNextFestival,
    
    // Theme controls
    festivalThemeEnabled: isEnabled,
    toggleFestivalTheme,
    
    // Festival bonuses
    calculateFestivalBonus,
    
    // User preferences
    userRegion,
    setUserRegion: handleSetUserRegion,
    favoritesFestivals,
    toggleFavoriteFestival
  };
};