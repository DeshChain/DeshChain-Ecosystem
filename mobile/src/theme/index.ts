import { DefaultTheme } from 'react-native-paper';
import { Theme } from '@react-navigation/native';

// DeshChain brand colors inspired by Indian culture
export const colors = {
  // Primary colors - Saffron gradient
  primary: '#FF9933', // Saffron
  primaryDark: '#E6751A',
  primaryLight: '#FFB366',
  
  // Secondary colors - Green for prosperity
  secondary: '#138808', // Indian Green
  secondaryDark: '#0D5A05',
  secondaryLight: '#4CA442',
  
  // Accent colors
  accent: '#000080', // Navy Blue (Chakra Blue)
  gold: '#FFD700', // Gold for premium features
  
  // Status colors
  success: '#4CAF50',
  warning: '#FF9800',
  error: '#F44336',
  info: '#2196F3',
  
  // Trust score colors
  diamond: '#9C27B0', // Purple
  platinum: '#607D8B', // Blue Grey
  goldBadge: '#FF9800', // Orange
  silver: '#9E9E9E', // Grey
  bronze: '#795548', // Brown
  
  // Neutral colors
  background: '#FFFFFF',
  surface: '#F5F5F5',
  card: '#FFFFFF',
  text: '#212121',
  textSecondary: '#757575',
  disabled: '#BDBDBD',
  placeholder: '#9E9E9E',
  border: '#E0E0E0',
  
  // Dark theme variants
  dark: {
    background: '#121212',
    surface: '#1E1E1E',
    card: '#2D2D2D',
    text: '#FFFFFF',
    textSecondary: '#B3B3B3',
    border: '#333333',
  },
  
  // Cultural colors
  cultural: {
    saffron: '#FF9933',
    white: '#FFFFFF',
    green: '#138808',
    lotus: '#FFB6C1',
    marigold: '#FFA500',
    turmeric: '#E6BE00',
    vermillion: '#E34234',
    peacock: '#005F69',
  },
};

// Typography scale following Material Design with Indian font preferences
export const typography = {
  fontFamily: {
    regular: 'Roboto-Regular',
    medium: 'Roboto-Medium',
    bold: 'Roboto-Bold',
    light: 'Roboto-Light',
    // Support for Devanagari script
    hindi: 'NotoSansDevanagari-Regular',
    hindiBold: 'NotoSansDevanagari-Bold',
  },
  fontSize: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 18,
    xl: 20,
    xxl: 24,
    xxxl: 32,
  },
  lineHeight: {
    xs: 16,
    sm: 20,
    md: 24,
    lg: 28,
    xl: 32,
    xxl: 36,
    xxxl: 48,
  },
};

// Spacing scale
export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
  xxxl: 64,
};

// Border radius scale
export const borderRadius = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  round: 50,
};

// Shadow styles for elevation
export const shadows = {
  small: {
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  medium: {
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 8,
    elevation: 4,
  },
  large: {
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.2,
    shadowRadius: 16,
    elevation: 8,
  },
};

// React Native Paper theme
export const theme = {
  ...DefaultTheme,
  colors: {
    ...DefaultTheme.colors,
    primary: colors.primary,
    accent: colors.accent,
    background: colors.background,
    surface: colors.surface,
    text: colors.text,
    disabled: colors.disabled,
    placeholder: colors.placeholder,
    backdrop: 'rgba(0, 0, 0, 0.5)',
    onSurface: colors.text,
    notification: colors.error,
    
    // Custom colors
    secondary: colors.secondary,
    success: colors.success,
    warning: colors.warning,
    error: colors.error,
    info: colors.info,
    gold: colors.gold,
    
    // Trust score colors
    diamond: colors.diamond,
    platinum: colors.platinum,
    goldBadge: colors.goldBadge,
    silver: colors.silver,
    bronze: colors.bronze,
  },
  fonts: {
    ...DefaultTheme.fonts,
    regular: {
      fontFamily: typography.fontFamily.regular,
      fontWeight: 'normal' as const,
    },
    medium: {
      fontFamily: typography.fontFamily.medium,
      fontWeight: '500' as const,
    },
    light: {
      fontFamily: typography.fontFamily.light,
      fontWeight: '300' as const,
    },
    thin: {
      fontFamily: typography.fontFamily.light,
      fontWeight: '100' as const,
    },
  },
  roundness: borderRadius.md,
  animation: {
    scale: 1.0,
  },
};

// Navigation theme
export const navigationTheme: Theme = {
  dark: false,
  colors: {
    primary: colors.primary,
    background: colors.background,
    card: colors.card,
    text: colors.text,
    border: colors.border,
    notification: colors.error,
  },
};

// Dark theme variant
export const darkTheme = {
  ...theme,
  colors: {
    ...theme.colors,
    primary: colors.primary,
    background: colors.dark.background,
    surface: colors.dark.surface,
    text: colors.dark.text,
    onSurface: colors.dark.text,
    disabled: colors.disabled,
    placeholder: colors.placeholder,
  },
};

// Cultural theme variants for festivals
export const culturalThemes = {
  diwali: {
    ...theme,
    colors: {
      ...theme.colors,
      primary: colors.cultural.marigold,
      secondary: colors.cultural.vermillion,
      accent: colors.gold,
    },
  },
  holi: {
    ...theme,
    colors: {
      ...theme.colors,
      primary: colors.cultural.peacock,
      secondary: colors.cultural.lotus,
      accent: colors.cultural.turmeric,
    },
  },
  independence: {
    ...theme,
    colors: {
      ...theme.colors,
      primary: colors.cultural.saffron,
      secondary: colors.cultural.green,
      accent: colors.cultural.white,
    },
  },
};

// Component-specific styles
export const componentStyles = {
  button: {
    primary: {
      backgroundColor: colors.primary,
      borderRadius: borderRadius.md,
      paddingVertical: spacing.md,
      paddingHorizontal: spacing.lg,
    },
    secondary: {
      backgroundColor: 'transparent',
      borderColor: colors.primary,
      borderWidth: 1,
      borderRadius: borderRadius.md,
      paddingVertical: spacing.md,
      paddingHorizontal: spacing.lg,
    },
    fab: {
      backgroundColor: colors.primary,
      borderRadius: borderRadius.round,
      width: 56,
      height: 56,
      ...shadows.medium,
    },
  },
  card: {
    default: {
      backgroundColor: colors.card,
      borderRadius: borderRadius.lg,
      padding: spacing.md,
      ...shadows.small,
    },
    elevated: {
      backgroundColor: colors.card,
      borderRadius: borderRadius.lg,
      padding: spacing.md,
      ...shadows.medium,
    },
  },
  input: {
    default: {
      borderRadius: borderRadius.md,
      borderWidth: 1,
      borderColor: colors.border,
      paddingVertical: spacing.md,
      paddingHorizontal: spacing.md,
      fontSize: typography.fontSize.md,
    },
  },
};

export default theme;