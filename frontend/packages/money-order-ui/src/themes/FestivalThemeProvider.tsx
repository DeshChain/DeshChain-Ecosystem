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

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { ThemeProvider, createTheme, Theme } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { deepmerge } from '@mui/utils';

import { Festival, getCurrentFestival, FestivalTheme } from './festivals';
import { FestivalParticles } from '../components/festival/FestivalParticles';
import { useLanguage } from '../hooks/useLanguage';

interface FestivalThemeContextValue {
  currentFestival: Festival | null;
  isEnabled: boolean;
  toggleFestivalTheme: (enabled: boolean) => void;
  baseTheme: Theme;
  festivalTheme: Theme;
}

const FestivalThemeContext = createContext<FestivalThemeContextValue | undefined>(undefined);

export const useFestivalTheme = () => {
  const context = useContext(FestivalThemeContext);
  if (!context) {
    throw new Error('useFestivalTheme must be used within FestivalThemeProvider');
  }
  return context;
};

interface FestivalThemeProviderProps {
  children: ReactNode;
  defaultTheme?: Theme;
}

export const FestivalThemeProvider: React.FC<FestivalThemeProviderProps> = ({
  children,
  defaultTheme
}) => {
  const { currentLanguage } = useLanguage();
  const [currentFestival, setCurrentFestival] = useState<Festival | null>(null);
  const [isEnabled, setIsEnabled] = useState(true);

  // Default base theme
  const baseTheme = defaultTheme || createTheme({
    palette: {
      primary: {
        main: '#FF6B35',
        light: '#FF8F65',
        dark: '#E55100',
        contrastText: '#FFFFFF'
      },
      secondary: {
        main: '#138808',
        light: '#4CAF50',
        dark: '#0D5F07',
        contrastText: '#FFFFFF'
      },
      background: {
        default: '#FAFAFA',
        paper: '#FFFFFF'
      },
      text: {
        primary: '#212121',
        secondary: '#757575'
      }
    },
    typography: {
      fontFamily: '"Roboto", "Noto Sans", "Helvetica", "Arial", sans-serif',
      h1: {
        fontSize: '2.5rem',
        fontWeight: 700
      },
      h2: {
        fontSize: '2rem',
        fontWeight: 600
      },
      h3: {
        fontSize: '1.75rem',
        fontWeight: 600
      }
    },
    shape: {
      borderRadius: 8
    }
  });

  // Create festival-themed MUI theme
  const createFestivalTheme = (festival: Festival): Theme => {
    const festivalPalette = {
      primary: {
        main: festival.theme.primary,
        contrastText: '#FFFFFF'
      },
      secondary: {
        main: festival.theme.secondary,
        contrastText: '#FFFFFF'
      },
      background: {
        default: festival.theme.background,
        paper: '#FFFFFF'
      },
      text: {
        primary: festival.theme.textPrimary,
        secondary: festival.theme.textSecondary
      }
    };

    const festivalOverrides = {
      palette: festivalPalette,
      components: {
        MuiButton: {
          styleOverrides: {
            root: {
              background: festival.theme.backgroundGradient,
              '&:hover': {
                background: festival.theme.primary,
                transform: 'translateY(-2px)',
                boxShadow: '0 4px 20px rgba(0,0,0,0.1)'
              }
            }
          }
        },
        MuiCard: {
          styleOverrides: {
            root: {
              background: `linear-gradient(135deg, ${festival.theme.background} 0%, ${festival.theme.backgroundGradient})`,
              borderImage: `linear-gradient(45deg, ${festival.theme.primary}, ${festival.theme.secondary}) 1`,
              borderImageSlice: 1
            }
          }
        },
        MuiChip: {
          styleOverrides: {
            root: {
              background: festival.theme.accent,
              color: '#FFFFFF'
            }
          }
        }
      }
    };

    return createTheme(deepmerge(baseTheme, festivalOverrides));
  };

  // Check for current festival
  useEffect(() => {
    const checkFestival = () => {
      const festival = getCurrentFestival();
      setCurrentFestival(festival);
    };

    checkFestival();
    
    // Check every hour for festival changes
    const interval = setInterval(checkFestival, 60 * 60 * 1000);
    
    return () => clearInterval(interval);
  }, []);

  // Load festival theme preference
  useEffect(() => {
    const saved = localStorage.getItem('deshchain-festival-theme');
    if (saved !== null) {
      setIsEnabled(saved === 'true');
    }
  }, []);

  const toggleFestivalTheme = (enabled: boolean) => {
    setIsEnabled(enabled);
    localStorage.setItem('deshchain-festival-theme', enabled.toString());
  };

  const activeTheme = isEnabled && currentFestival 
    ? createFestivalTheme(currentFestival)
    : baseTheme;

  const contextValue: FestivalThemeContextValue = {
    currentFestival,
    isEnabled,
    toggleFestivalTheme,
    baseTheme,
    festivalTheme: activeTheme
  };

  return (
    <FestivalThemeContext.Provider value={contextValue}>
      <ThemeProvider theme={activeTheme}>
        <CssBaseline />
        
        {/* Global festival styles */}
        {isEnabled && currentFestival && (
          <style>
            {`
              @keyframes glow {
                0% { box-shadow: 0 0 5px ${currentFestival.theme.primary}40; }
                50% { box-shadow: 0 0 20px ${currentFestival.theme.primary}60, 0 0 30px ${currentFestival.theme.primary}40; }
                100% { box-shadow: 0 0 5px ${currentFestival.theme.primary}40; }
              }
              
              @keyframes sparkle {
                0% { opacity: 0; }
                50% { opacity: 1; }
                100% { opacity: 0; }
              }
              
              @keyframes splash {
                0% { transform: scale(0) rotate(0deg); opacity: 1; }
                100% { transform: scale(1.5) rotate(180deg); opacity: 0; }
              }
              
              @keyframes rainbow {
                0% { background-position: 0% 50%; }
                50% { background-position: 100% 50%; }
                100% { background-position: 0% 50%; }
              }
              
              @keyframes wave {
                0% { transform: translateX(-100%); }
                100% { transform: translateX(100%); }
              }
              
              @keyframes bloom {
                0% { transform: scale(0.8) rotate(-5deg); }
                50% { transform: scale(1.1) rotate(5deg); }
                100% { transform: scale(1) rotate(0deg); }
              }
              
              @keyframes dance {
                0%, 100% { transform: translateY(0) rotate(0deg); }
                25% { transform: translateY(-10px) rotate(-5deg); }
                75% { transform: translateY(-10px) rotate(5deg); }
              }
              
              .festival-glow {
                animation: glow 2s ease-in-out infinite;
              }
              
              .festival-sparkle {
                position: relative;
                overflow: hidden;
              }
              
              .festival-sparkle::before {
                content: '';
                position: absolute;
                top: -50%;
                left: -50%;
                width: 200%;
                height: 200%;
                background: linear-gradient(
                  45deg,
                  transparent 30%,
                  ${currentFestival.theme.accent}20 50%,
                  transparent 70%
                );
                animation: sparkle 3s linear infinite;
              }
              
              .festival-gradient-text {
                background: linear-gradient(
                  45deg,
                  ${currentFestival.theme.primary},
                  ${currentFestival.theme.secondary},
                  ${currentFestival.theme.accent}
                );
                -webkit-background-clip: text;
                -webkit-text-fill-color: transparent;
                background-size: 200% auto;
                animation: rainbow 3s ease-in-out infinite;
              }
              
              .festival-card {
                position: relative;
                overflow: hidden;
                background: ${currentFestival.theme.backgroundGradient};
              }
              
              .festival-card::before {
                content: '';
                position: absolute;
                top: 0;
                left: -100%;
                width: 100%;
                height: 100%;
                background: linear-gradient(
                  90deg,
                  transparent,
                  ${currentFestival.theme.accent}40,
                  transparent
                );
                animation: wave 3s linear infinite;
              }
              
              /* Festival-specific decorations */
              ${currentFestival.id === 'diwali' ? `
                .diwali-diya {
                  position: relative;
                }
                
                .diwali-diya::after {
                  content: '🪔';
                  position: absolute;
                  top: -10px;
                  right: -10px;
                  font-size: 20px;
                  animation: glow 2s ease-in-out infinite;
                }
              ` : ''}
              
              ${currentFestival.id === 'holi' ? `
                .holi-splash {
                  position: relative;
                }
                
                .holi-splash::before,
                .holi-splash::after {
                  content: '';
                  position: absolute;
                  width: 20px;
                  height: 20px;
                  border-radius: 50%;
                  animation: splash 3s linear infinite;
                }
                
                .holi-splash::before {
                  background: ${currentFestival.theme.primary};
                  top: 10%;
                  left: 10%;
                }
                
                .holi-splash::after {
                  background: ${currentFestival.theme.secondary};
                  bottom: 10%;
                  right: 10%;
                  animation-delay: 1.5s;
                }
              ` : ''}
            `}
          </style>
        )}
        
        {/* Festival particles */}
        {isEnabled && currentFestival && (
          <FestivalParticles festival={currentFestival} />
        )}
        
        {children}
      </ThemeProvider>
    </FestivalThemeContext.Provider>
  );
};