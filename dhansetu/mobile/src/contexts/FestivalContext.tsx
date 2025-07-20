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

import React, { createContext, useContext, useEffect, useState } from 'react';
import axios from 'axios';
import { COLORS } from '@constants/theme';

export interface Festival {
  id: string;
  name: string;
  description: string;
  date: string;
  endDate?: string;
  type: string;
  region: string;
  bonusRate: string;
  isActive: boolean;
  traditionalGreeting: string;
  culturalTheme: string;
  daysRemaining: number;
  significance: string;
  colors?: {
    primary: string;
    secondary: string;
    accent: string;
  };
}

interface FestivalContextValue {
  activeFestivals: Festival[];
  upcomingFestivals: Festival[];
  currentFestival: Festival | null;
  festivalBonusRate: number;
  isLoading: boolean;
  error: string | null;
  refreshFestivals: () => Promise<void>;
  getFestivalColors: (festivalId: string) => {
    primary: string;
    secondary: string;
    accent: string;
  };
}

const FestivalContext = createContext<FestivalContextValue | undefined>(undefined);

export const useFestival = () => {
  const context = useContext(FestivalContext);
  if (!context) {
    throw new Error('useFestival must be used within FestivalProvider');
  }
  return context;
};

interface FestivalProviderProps {
  children: React.ReactNode;
}

export const FestivalProvider: React.FC<FestivalProviderProps> = ({ children }) => {
  const [activeFestivals, setActiveFestivals] = useState<Festival[]>([]);
  const [upcomingFestivals, setUpcomingFestivals] = useState<Festival[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchFestivals = async () => {
    setIsLoading(true);
    setError(null);

    try {
      // Fetch festival status from DeshChain
      const response = await axios.get(
        `${process.env.EXPO_PUBLIC_API_URL}/cultural/festivals/status`
      );

      const { activeFestivals: active, upcomingFestivals: upcoming } = response.data;
      
      setActiveFestivals(active || []);
      setUpcomingFestivals(upcoming || []);
    } catch (err) {
      console.error('Failed to fetch festivals:', err);
      setError('Failed to fetch festival information');
      
      // Use mock data as fallback
      setActiveFestivals(getMockActiveFestivals());
      setUpcomingFestivals(getMockUpcomingFestivals());
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchFestivals();
    
    // Refresh every hour
    const interval = setInterval(fetchFestivals, 60 * 60 * 1000);
    
    return () => clearInterval(interval);
  }, []);

  const currentFestival = activeFestivals[0] || null;
  
  const festivalBonusRate = currentFestival 
    ? parseFloat(currentFestival.bonusRate) / 100 
    : 0;

  const getFestivalColors = (festivalId: string) => {
    const festival = [...activeFestivals, ...upcomingFestivals].find(
      (f) => f.id === festivalId
    );
    
    return festival?.colors || {
      primary: COLORS.saffron,
      secondary: COLORS.white,
      accent: COLORS.green,
    };
  };

  const value: FestivalContextValue = {
    activeFestivals,
    upcomingFestivals,
    currentFestival,
    festivalBonusRate,
    isLoading,
    error,
    refreshFestivals: fetchFestivals,
    getFestivalColors,
  };

  return (
    <FestivalContext.Provider value={value}>
      {children}
    </FestivalContext.Provider>
  );
};

// Mock data for development
const getMockActiveFestivals = (): Festival[] => {
  const today = new Date();
  return [
    {
      id: 'diwali_2024',
      name: 'Diwali',
      description: 'Festival of Lights - Celebrating prosperity and victory of light over darkness',
      date: today.toISOString(),
      endDate: new Date(today.getTime() + 5 * 24 * 60 * 60 * 1000).toISOString(),
      type: 'religious',
      region: 'pan_india',
      bonusRate: '15',
      isActive: true,
      traditionalGreeting: 'Shubh Deepawali!',
      culturalTheme: 'lights_prosperity',
      daysRemaining: 5,
      significance: 'Most important Hindu festival celebrating light, knowledge, and prosperity',
      colors: {
        primary: '#FF6B35',
        secondary: '#F7931E',
        accent: '#FFD700',
      },
    },
  ];
};

const getMockUpcomingFestivals = (): Festival[] => {
  const futureDate = new Date();
  futureDate.setDate(futureDate.getDate() + 30);
  
  return [
    {
      id: 'holi_2025',
      name: 'Holi',
      description: 'Festival of Colors - Celebrating spring and victory of good over evil',
      date: futureDate.toISOString(),
      type: 'religious',
      region: 'pan_india',
      bonusRate: '12',
      isActive: false,
      traditionalGreeting: 'Holi Hai!',
      culturalTheme: 'colors_spring',
      daysRemaining: 30,
      significance: 'Spring festival celebrating colors, love, and new beginnings',
      colors: {
        primary: '#E91E63',
        secondary: '#9C27B0',
        accent: '#FF9800',
      },
    },
  ];
};