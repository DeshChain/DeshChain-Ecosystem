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
import { useQuery, useQueryClient } from '@tanstack/react-query';

export interface FestivalInfo {
  id: string;
  name: string;
  description: string;
  date: string;
  endDate?: string;
  type: 'national' | 'religious' | 'regional' | 'seasonal' | 'cultural';
  region: string;
  bonusRate: number;
  isActive: boolean;
  traditionalGreeting: string;
  culturalTheme: string;
  colors: {
    primary: string;
    secondary: string;
    accent: string;
  };
  daysRemaining: number;
  significance: string;
  tags: string[];
  isVerified: boolean;
}

export interface FestivalStatus {
  activeFestivals: string[];
  upcomingFestivals: string[];
  currentBonusRate: number;
  nextFestival?: FestivalInfo;
  lastUpdated: string;
  metadata: Record<string, string>;
}

interface FestivalSyncState {
  isLoading: boolean;
  error: string | null;
  activeFestivals: FestivalInfo[];
  upcomingFestivals: FestivalInfo[];
  currentBonusRate: number;
  nextFestival?: FestivalInfo;
  lastSynced: Date | null;
}

interface FestivalAPI {
  getActiveFestivals: () => Promise<FestivalInfo[]>;
  getUpcomingFestivals: () => Promise<FestivalInfo[]>;
  getFestivalStatus: () => Promise<FestivalStatus>;
  getFestivalById: (id: string) => Promise<FestivalInfo>;
}

// Mock API implementation for development
const createMockAPI = (): FestivalAPI => {
  const mockFestivals: FestivalInfo[] = [
    {
      id: 'diwali',
      name: 'Diwali',
      description: 'Festival of Lights celebrating prosperity and good fortune',
      date: new Date(new Date().getFullYear(), 10, 1).toISOString(), // November 1st
      type: 'national',
      region: 'all_india',
      bonusRate: 0.15,
      isActive: false,
      traditionalGreeting: 'दीपावली की शुभकामनाएं! (Deepavali ki Shubhkamnayein!)',
      culturalTheme: 'lights_prosperity',
      colors: {
        primary: '#FF6B35',
        secondary: '#F7931E',
        accent: '#FFD700',
      },
      daysRemaining: 30,
      significance: 'Celebrates the victory of light over darkness and good over evil',
      tags: ['prosperity', 'lights', 'tradition'],
      isVerified: true,
    },
    {
      id: 'holi',
      name: 'Holi',
      description: 'Festival of Colors celebrating spring and new beginnings',
      date: new Date(new Date().getFullYear() + 1, 2, 14).toISOString(), // March 14th next year
      type: 'national',
      region: 'all_india',
      bonusRate: 0.12,
      isActive: false,
      traditionalGreeting: 'होली की शुभकामनाएं! (Holi ki Shubhkamnayein!)',
      culturalTheme: 'colors_spring',
      colors: {
        primary: '#E91E63',
        secondary: '#9C27B0',
        accent: '#FF9800',
      },
      daysRemaining: 120,
      significance: 'Celebrates the triumph of good over evil and the arrival of spring',
      tags: ['colors', 'spring', 'joy'],
      isVerified: true,
    },
  ];

  return {
    getActiveFestivals: async () => {
      await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay
      return mockFestivals.filter(f => f.isActive);
    },

    getUpcomingFestivals: async () => {
      await new Promise(resolve => setTimeout(resolve, 500));
      return mockFestivals.filter(f => !f.isActive && f.daysRemaining <= 30);
    },

    getFestivalStatus: async () => {
      await new Promise(resolve => setTimeout(resolve, 300));
      const activeFestivals = mockFestivals.filter(f => f.isActive);
      const upcomingFestivals = mockFestivals.filter(f => !f.isActive && f.daysRemaining <= 30);
      
      return {
        activeFestivals: activeFestivals.map(f => f.id),
        upcomingFestivals: upcomingFestivals.map(f => f.id),
        currentBonusRate: activeFestivals.reduce((max, f) => Math.max(max, f.bonusRate), 0),
        nextFestival: upcomingFestivals.sort((a, b) => a.daysRemaining - b.daysRemaining)[0],
        lastUpdated: new Date().toISOString(),
        metadata: {},
      };
    },

    getFestivalById: async (id: string) => {
      await new Promise(resolve => setTimeout(resolve, 200));
      const festival = mockFestivals.find(f => f.id === id);
      if (!festival) {
        throw new Error(`Festival with id ${id} not found`);
      }
      return festival;
    },
  };
};

// Real API implementation
const createRealAPI = (baseUrl: string): FestivalAPI => {
  const request = async <T>(endpoint: string): Promise<T> => {
    const response = await fetch(`${baseUrl}${endpoint}`, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }
    
    return response.json();
  };

  return {
    getActiveFestivals: () => request<FestivalInfo[]>('/cultural/festivals/active'),
    getUpcomingFestivals: () => request<FestivalInfo[]>('/cultural/festivals/upcoming'),
    getFestivalStatus: () => request<FestivalStatus>('/cultural/festivals/status'),
    getFestivalById: (id: string) => request<FestivalInfo>(`/cultural/festivals/${id}`),
  };
};

export const useFestivalSync = (options?: {
  enablePolling?: boolean;
  pollingInterval?: number;
  baseUrl?: string;
}) => {
  const {
    enablePolling = true,
    pollingInterval = 60000, // 1 minute
    baseUrl = process.env.REACT_APP_API_BASE_URL || '',
  } = options || {};

  const queryClient = useQueryClient();
  const [state, setState] = useState<FestivalSyncState>({
    isLoading: true,
    error: null,
    activeFestivals: [],
    upcomingFestivals: [],
    currentBonusRate: 0,
    lastSynced: null,
  });

  // Create API instance based on environment
  const api = baseUrl ? createRealAPI(baseUrl) : createMockAPI();

  // Query for festival status
  const { data: festivalStatus, error: statusError, isLoading: statusLoading } = useQuery({
    queryKey: ['festival-status'],
    queryFn: api.getFestivalStatus,
    refetchInterval: enablePolling ? pollingInterval : false,
    staleTime: 30000, // Consider data stale after 30 seconds
    gcTime: 300000, // Keep in cache for 5 minutes
  });

  // Query for active festivals
  const { data: activeFestivals, error: activeError, isLoading: activeLoading } = useQuery({
    queryKey: ['active-festivals'],
    queryFn: api.getActiveFestivals,
    refetchInterval: enablePolling ? pollingInterval : false,
    staleTime: 30000,
    gcTime: 300000,
  });

  // Query for upcoming festivals
  const { data: upcomingFestivals, error: upcomingError, isLoading: upcomingLoading } = useQuery({
    queryKey: ['upcoming-festivals'],
    queryFn: api.getUpcomingFestivals,
    refetchInterval: enablePolling ? pollingInterval : false,
    staleTime: 60000, // Upcoming festivals can be stale for longer
    gcTime: 600000, // Keep in cache for 10 minutes
  });

  // Update state when queries resolve
  useEffect(() => {
    const isLoading = statusLoading || activeLoading || upcomingLoading;
    const error = statusError?.message || activeError?.message || upcomingError?.message || null;

    setState(prevState => ({
      ...prevState,
      isLoading,
      error,
      activeFestivals: activeFestivals || [],
      upcomingFestivals: upcomingFestivals || [],
      currentBonusRate: festivalStatus?.currentBonusRate || 0,
      nextFestival: festivalStatus?.nextFestival,
      lastSynced: error ? prevState.lastSynced : new Date(),
    }));
  }, [
    statusLoading, activeLoading, upcomingLoading,
    statusError, activeError, upcomingError,
    festivalStatus, activeFestivals, upcomingFestivals
  ]);

  // Manual sync function
  const syncFestivals = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['festival-status'] }),
        queryClient.invalidateQueries({ queryKey: ['active-festivals'] }),
        queryClient.invalidateQueries({ queryKey: ['upcoming-festivals'] }),
      ]);
      
      setState(prev => ({ 
        ...prev, 
        isLoading: false, 
        lastSynced: new Date() 
      }));
    } catch (error) {
      setState(prev => ({ 
        ...prev, 
        isLoading: false, 
        error: error instanceof Error ? error.message : 'Sync failed' 
      }));
    }
  }, [queryClient]);

  // Get festival by ID
  const getFestivalById = useCallback(async (id: string): Promise<FestivalInfo | null> => {
    try {
      return await api.getFestivalById(id);
    } catch (error) {
      console.error(`Failed to fetch festival ${id}:`, error);
      return null;
    }
  }, [api]);

  // Check if a specific festival is active
  const isFestivalActive = useCallback((festivalId: string): boolean => {
    return state.activeFestivals.some(festival => festival.id === festivalId);
  }, [state.activeFestivals]);

  // Get current festival bonus rate
  const getCurrentBonusRate = useCallback((): number => {
    return state.currentBonusRate;
  }, [state.currentBonusRate]);

  // Get festival colors for theming
  const getFestivalColors = useCallback((festivalId?: string) => {
    if (!festivalId) {
      const activeFestival = state.activeFestivals[0];
      return activeFestival?.colors || {
        primary: '#FF6B35',
        secondary: '#F7931E',
        accent: '#FFD700',
      };
    }

    const festival = [...state.activeFestivals, ...state.upcomingFestivals]
      .find(f => f.id === festivalId);
    
    return festival?.colors || {
      primary: '#FF6B35',
      secondary: '#F7931E',
      accent: '#FFD700',
    };
  }, [state.activeFestivals, state.upcomingFestivals]);

  // Get festival greeting
  const getFestivalGreeting = useCallback((festivalId?: string): string => {
    if (!festivalId) {
      const activeFestival = state.activeFestivals[0];
      return activeFestival?.traditionalGreeting || 'नमस्ते! (Namaste!)';
    }

    const festival = [...state.activeFestivals, ...state.upcomingFestivals]
      .find(f => f.id === festivalId);
    
    return festival?.traditionalGreeting || 'नमस्ते! (Namaste!)';
  }, [state.activeFestivals, state.upcomingFestivals]);

  return {
    // State
    ...state,
    
    // Festival data
    allFestivals: [...state.activeFestivals, ...state.upcomingFestivals],
    hasActiveFestivals: state.activeFestivals.length > 0,
    hasUpcomingFestivals: state.upcomingFestivals.length > 0,
    
    // Actions
    syncFestivals,
    getFestivalById,
    
    // Utilities
    isFestivalActive,
    getCurrentBonusRate,
    getFestivalColors,
    getFestivalGreeting,
    
    // Status
    isOnline: !state.error,
    lastSyncAge: state.lastSynced ? Date.now() - state.lastSynced.getTime() : null,
  };
};

export default useFestivalSync;