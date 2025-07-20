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

import React from 'react';
import type { AppProps } from 'next/app';
import Head from 'next/head';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import {
  MoneyOrderProvider,
  CulturalProvider,
  FestivalThemeProvider
} from '@deshchain/money-order-ui';

import { config } from '../config';
import { Layout } from '../components/Layout';
import '../styles/globals.css';

// Create MUI theme with Indian cultural colors
const theme = createTheme({
  palette: {
    primary: {
      main: '#FF6B35', // Saffron
      light: '#FF8F65',
      dark: '#E85D25'
    },
    secondary: {
      main: '#138808', // Green
      light: '#4CAF50',
      dark: '#0F6605'
    },
    info: {
      main: '#000080' // Navy Blue
    },
    background: {
      default: '#FAFAFA',
      paper: '#FFFFFF'
    }
  },
  typography: {
    fontFamily: '"Inter", "Noto Sans", "Noto Sans Devanagari", sans-serif',
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
    borderRadius: 12
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 600,
          borderRadius: 8
        },
        contained: {
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(0,0,0,0.15)'
          }
        }
      }
    },
    MuiCard: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
          '&:hover': {
            boxShadow: '0 4px 16px rgba(0,0,0,0.12)'
          }
        }
      }
    }
  }
});

// Create React Query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000 // 5 minutes
    }
  }
});

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>DeshChain Money Order - Simple & Cultural Money Transfers</title>
        <meta name="description" content="Send money orders with cultural integration, supporting 22 Indian languages and festival celebrations" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
        
        {/* Preload fonts for better performance */}
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Noto+Sans+Devanagari:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </Head>

      <QueryClientProvider client={queryClient}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <MoneyOrderProvider config={config.moneyOrder}>
            <CulturalProvider 
              initialLanguage={config.defaultLanguage}
              enabledFeatures={config.features}
            >
              <FestivalThemeProvider>
                <Layout>
                  <Component {...pageProps} />
                  <Toaster
                    position="top-right"
                    toastOptions={{
                      duration: 4000,
                      style: {
                        background: theme.palette.background.paper,
                        color: theme.palette.text.primary,
                        boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                        borderRadius: '8px'
                      },
                      success: {
                        iconTheme: {
                          primary: theme.palette.success.main,
                          secondary: '#fff'
                        }
                      },
                      error: {
                        iconTheme: {
                          primary: theme.palette.error.main,
                          secondary: '#fff'
                        }
                      }
                    }}
                  />
                </Layout>
              </FestivalThemeProvider>
            </CulturalProvider>
          </MoneyOrderProvider>
        </ThemeProvider>
      </QueryClientProvider>
    </>
  );
}

export default MyApp;