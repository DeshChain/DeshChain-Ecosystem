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

import { useMemo } from 'react';
import {
  Celebration as CelebrationIcon,
  Palette as PaletteIcon,
  Flag as FlagIcon,
  MilitaryTech as MilitaryTechIcon,
  Diversity1 as DiversityIcon,
  Favorite as FavoriteIcon,
  AutoStories as AutoStoriesIcon,
  Agriculture as AgricultureIcon,
  School as SchoolIcon,
} from '@mui/icons-material';
import { SvgIconComponent } from '@mui/icons-material';

import { useFestivalSync, FestivalInfo } from './useFestivalSync';

export interface FestivalColors {
  primary: string;
  secondary: string;
  accent: string;
}

export interface FestivalTheme {
  colors: FestivalColors;
  icon: SvgIconComponent;
  gradient: string;
  textColor: string;
  bgColor: string;
}

// Default festival colors
const defaultColors: FestivalColors = {
  primary: '#FF6B35',
  secondary: '#F7931E', 
  accent: '#FFD700',
};

// Festival theme mappings
const festivalThemes: Record<string, Partial<FestivalTheme>> = {
  // Diwali - Festival of Lights
  'lights_prosperity': {
    colors: {
      primary: '#FF6B35',
      secondary: '#F7931E',
      accent: '#FFD700',
    },
    icon: CelebrationIcon,
  },
  
  // Holi - Festival of Colors
  'colors_spring': {
    colors: {
      primary: '#E91E63',
      secondary: '#9C27B0',
      accent: '#FF9800',
    },
    icon: PaletteIcon,
  },
  
  // Independence Day / Republic Day - Patriotic festivals
  'patriotism_freedom': {
    colors: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808',
    },
    icon: FlagIcon,
  },
  
  'constitution_democracy': {
    colors: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808',
    },
    icon: FlagIcon,
  },
  
  // Dussehra - Victory of good over evil
  'victory_righteousness': {
    colors: {
      primary: '#DC143C',
      secondary: '#FFD700',
      accent: '#FF4500',
    },
    icon: MilitaryTechIcon,
  },
  
  // Eid - Harmony and brotherhood
  'harmony_brotherhood': {
    colors: {
      primary: '#00A86B',
      secondary: '#FFFFFF',
      accent: '#FFD700',
    },
    icon: DiversityIcon,
  },
  
  // Christmas - Peace and joy
  'peace_joy': {
    colors: {
      primary: '#DC143C',
      secondary: '#228B22',
      accent: '#FFD700',
    },
    icon: FavoriteIcon,
  },
  
  // Guru Nanak Jayanti - Wisdom and devotion
  'wisdom_devotion': {
    colors: {
      primary: '#FF9933',
      secondary: '#FFFFFF',
      accent: '#138808',
    },
    icon: AutoStoriesIcon,
  },
  
  // Harvest festivals
  'harvest_gratitude': {
    colors: {
      primary: '#8BC34A',
      secondary: '#FFC107',
      accent: '#FF5722',
    },
    icon: AgricultureIcon,
  },
  
  // Educational festivals
  'knowledge_wisdom': {
    colors: {
      primary: '#2196F3',
      secondary: '#FFC107',
      accent: '#9C27B0',
    },
    icon: SchoolIcon,
  },
};

// Festival ID to theme mapping
const festivalIdToTheme: Record<string, string> = {
  'diwali': 'lights_prosperity',
  'holi': 'colors_spring',
  'independence_day': 'patriotism_freedom',
  'republic_day': 'constitution_democracy',
  'dussehra': 'victory_righteousness',
  'eid_ul_fitr': 'harmony_brotherhood',
  'christmas': 'peace_joy',
  'guru_nanak_jayanti': 'wisdom_devotion',
};

export const useFestivalTheme = () => {
  const { activeFestivals, upcomingFestivals, getFestivalColors } = useFestivalSync();
  
  // Get current festival theme
  const currentTheme = useMemo((): FestivalTheme => {
    const activeFestival = activeFestivals[0];
    
    if (!activeFestival) {
      return createDefaultTheme();
    }
    
    const themeKey = activeFestival.culturalTheme || festivalIdToTheme[activeFestival.id];
    const theme = festivalThemes[themeKey];
    
    if (!theme) {
      return createDefaultTheme();
    }
    
    return createTheme(theme, activeFestival.colors);
  }, [activeFestivals]);
  
  // Get festival colors for a specific festival
  const getFestivalColorsById = (festivalId: string): FestivalColors => {
    const festival = [...activeFestivals, ...upcomingFestivals].find(f => f.id === festivalId);
    return festival?.colors || defaultColors;
  };
  
  // Get festival icon for a specific festival
  const getFestivalIcon = (festivalId: string): SvgIconComponent => {
    const festival = [...activeFestivals, ...upcomingFestivals].find(f => f.id === festivalId);
    const themeKey = festival?.culturalTheme || festivalIdToTheme[festivalId];
    const theme = festivalThemes[themeKey];
    
    return theme?.icon || CelebrationIcon;
  };
  
  // Get festival theme for a specific festival
  const getFestivalTheme = (festivalId: string): FestivalTheme => {
    const festival = [...activeFestivals, ...upcomingFestivals].find(f => f.id === festivalId);
    
    if (!festival) {
      return createDefaultTheme();
    }
    
    const themeKey = festival.culturalTheme || festivalIdToTheme[festivalId];
    const theme = festivalThemes[themeKey];
    
    if (!theme) {
      return createDefaultTheme();
    }
    
    return createTheme(theme, festival.colors);
  };
  
  // Create gradient string
  const createGradient = (colors: FestivalColors, direction = '135deg'): string => {
    return `linear-gradient(${direction}, ${colors.primary}, ${colors.secondary}, ${colors.accent})`;
  };
  
  // Create theme with gradient background
  const createThemeWithGradient = (colors: FestivalColors): string => {
    return `linear-gradient(135deg, ${colors.primary}20, ${colors.secondary}20, ${colors.accent}20)`;
  };
  
  // Get text color based on background
  const getTextColor = (backgroundColor: string): string => {
    // Simple implementation - can be enhanced with luminance calculation
    const lightColors = ['#FFFFFF', '#FFC107', '#FFD700', '#FFEB3B'];
    return lightColors.includes(backgroundColor) ? '#000000' : '#FFFFFF';
  };
  
  // Apply festival theme to component
  const applyFestivalTheme = (baseTheme: any, festivalId?: string) => {
    const theme = festivalId ? getFestivalTheme(festivalId) : currentTheme;
    
    return {
      ...baseTheme,
      palette: {
        ...baseTheme.palette,
        primary: {
          main: theme.colors.primary,
          light: theme.colors.accent,
          dark: theme.colors.secondary,
        },
        secondary: {
          main: theme.colors.secondary,
          light: theme.colors.accent,
          dark: theme.colors.primary,
        },
      },
      components: {
        ...baseTheme.components,
        MuiAppBar: {
          styleOverrides: {
            root: {
              background: theme.gradient,
            },
          },
        },
        MuiButton: {
          styleOverrides: {
            root: {
              '&.festival-themed': {
                background: createGradient(theme.colors),
                color: theme.textColor,
                '&:hover': {
                  background: createGradient(theme.colors, '45deg'),
                },
              },
            },
          },
        },
      },
    };
  };
  
  // Check if festival theming is active
  const isFestivalThemeActive = (): boolean => {
    return activeFestivals.length > 0;
  };
  
  return {
    // Current theme
    currentTheme,
    
    // Theme utilities
    getFestivalColors: getFestivalColorsById,
    getFestivalIcon,
    getFestivalTheme,
    
    // Gradient utilities
    createGradient,
    createThemeWithGradient,
    getTextColor,
    
    // Theme application
    applyFestivalTheme,
    isFestivalThemeActive,
    
    // Festival data
    activeFestivals,
    upcomingFestivals,
  };
};

// Helper functions
function createDefaultTheme(): FestivalTheme {
  return {
    colors: defaultColors,
    icon: CelebrationIcon,
    gradient: createGradientString(defaultColors),
    textColor: '#FFFFFF',
    bgColor: defaultColors.primary,
  };
}

function createTheme(partialTheme: Partial<FestivalTheme>, colors?: FestivalColors): FestivalTheme {
  const themeColors = colors || partialTheme.colors || defaultColors;
  
  return {
    colors: themeColors,
    icon: partialTheme.icon || CelebrationIcon,
    gradient: createGradientString(themeColors),
    textColor: getContrastColor(themeColors.primary),
    bgColor: themeColors.primary,
  };
}

function createGradientString(colors: FestivalColors, direction = '135deg'): string {
  return `linear-gradient(${direction}, ${colors.primary}, ${colors.secondary}, ${colors.accent})`;
}

function getContrastColor(hexColor: string): string {
  // Convert hex to RGB
  const r = parseInt(hexColor.slice(1, 3), 16);
  const g = parseInt(hexColor.slice(3, 5), 16);
  const b = parseInt(hexColor.slice(5, 7), 16);
  
  // Calculate luminance
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
  
  return luminance > 0.5 ? '#000000' : '#FFFFFF';
}

export default useFestivalTheme;