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

import { MD3LightTheme as DefaultTheme } from 'react-native-paper';

// Indian flag colors
export const COLORS = {
  // Primary colors
  saffron: '#FF9933',
  white: '#FFFFFF',
  green: '#138808',
  navy: '#000080',
  
  // Extended palette
  lightSaffron: '#FFB366',
  darkSaffron: '#CC7A29',
  lightGreen: '#1CAA0A',
  darkGreen: '#0E5905',
  
  // Neutral colors
  gray50: '#FAFAFA',
  gray100: '#F5F5F5',
  gray200: '#EEEEEE',
  gray300: '#E0E0E0',
  gray400: '#BDBDBD',
  gray500: '#9E9E9E',
  gray600: '#757575',
  gray700: '#616161',
  gray800: '#424242',
  gray900: '#212121',
  
  // Semantic colors
  success: '#4CAF50',
  warning: '#FF9800',
  error: '#F44336',
  info: '#2196F3',
  
  // Festival colors (can be dynamically changed)
  festivalPrimary: '#FF6B35',
  festivalSecondary: '#F7931E',
  festivalAccent: '#FFD700',
};

export const theme = {
  ...DefaultTheme,
  colors: {
    ...DefaultTheme.colors,
    primary: COLORS.saffron,
    secondary: COLORS.green,
    tertiary: COLORS.navy,
    background: COLORS.white,
    surface: COLORS.gray100,
    error: COLORS.error,
    onPrimary: COLORS.white,
    onSecondary: COLORS.white,
    onBackground: COLORS.gray900,
    onSurface: COLORS.gray900,
    onError: COLORS.white,
    
    // Custom colors
    saffron: COLORS.saffron,
    green: COLORS.green,
    navy: COLORS.navy,
    festivalPrimary: COLORS.festivalPrimary,
    festivalSecondary: COLORS.festivalSecondary,
    festivalAccent: COLORS.festivalAccent,
    
    // Gray scale
    gray50: COLORS.gray50,
    gray100: COLORS.gray100,
    gray200: COLORS.gray200,
    gray300: COLORS.gray300,
    gray400: COLORS.gray400,
    gray500: COLORS.gray500,
    gray600: COLORS.gray600,
    gray700: COLORS.gray700,
    gray800: COLORS.gray800,
    gray900: COLORS.gray900,
    
    // Semantic
    success: COLORS.success,
    warning: COLORS.warning,
    info: COLORS.info,
    text: COLORS.gray900,
    white: COLORS.white,
  },
  
  // Custom properties
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  
  borderRadius: {
    sm: 4,
    md: 8,
    lg: 16,
    xl: 24,
    full: 999,
  },
  
  fonts: {
    ...DefaultTheme.fonts,
    medium: {
      ...DefaultTheme.fonts.medium,
      fontFamily: 'Roboto-Medium',
    },
    regular: {
      ...DefaultTheme.fonts.regular,
      fontFamily: 'Roboto-Regular',
    },
    bold: {
      fontFamily: 'Roboto-Bold',
      fontWeight: 'bold' as const,
    },
  },
  
  gradients: {
    indianFlag: ['#FF9933', '#FFFFFF', '#138808'],
    saffronWhite: ['#FF9933', '#FFFFFF'],
    whiteGreen: ['#FFFFFF', '#138808'],
    festivalGradient: ['#FF6B35', '#F7931E', '#FFD700'],
  },
  
  shadows: {
    small: {
      shadowColor: '#000',
      shadowOffset: {
        width: 0,
        height: 2,
      },
      shadowOpacity: 0.1,
      shadowRadius: 3,
      elevation: 2,
    },
    medium: {
      shadowColor: '#000',
      shadowOffset: {
        width: 0,
        height: 4,
      },
      shadowOpacity: 0.15,
      shadowRadius: 6,
      elevation: 4,
    },
    large: {
      shadowColor: '#000',
      shadowOffset: {
        width: 0,
        height: 8,
      },
      shadowOpacity: 0.2,
      shadowRadius: 12,
      elevation: 8,
    },
  },
};

export type Theme = typeof theme;