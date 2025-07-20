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

import React, { useEffect } from 'react';
import { StatusBar } from 'react-native';
import { Provider } from 'react-redux';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { NavigationContainer } from '@react-navigation/native';
import { ThemeProvider } from 'react-native-paper';
import SplashScreen from 'react-native-splash-screen';
import { GestureHandlerRootView } from 'react-native-gesture-handler';

import { store } from '@store/index';
import { AppNavigator } from '@navigation/AppNavigator';
import { theme } from '@constants/theme';
import { LocalizationProvider } from '@contexts/LocalizationContext';
import { AuthProvider } from '@contexts/AuthContext';
import { WalletProvider } from '@contexts/WalletContext';
import { FestivalProvider } from '@contexts/FestivalContext';
import initializeApp from '@services/app/initialization';
import ErrorBoundary from '@components/common/ErrorBoundary';
import LoadingOverlay from '@components/common/LoadingOverlay';

// Import node polyfills
import 'react-native-get-random-values';
import '@ethersproject/shims';
import { Buffer } from 'buffer';
global.Buffer = Buffer;

export default function App() {
  useEffect(() => {
    const init = async () => {
      try {
        // Initialize app services
        await initializeApp();
        
        // Hide splash screen
        SplashScreen.hide();
      } catch (error) {
        console.error('App initialization failed:', error);
        // Show error screen or fallback
      }
    };
    
    init();
  }, []);

  return (
    <ErrorBoundary>
      <GestureHandlerRootView style={{ flex: 1 }}>
        <Provider store={store}>
          <LocalizationProvider>
            <AuthProvider>
              <WalletProvider>
                <FestivalProvider>
                  <ThemeProvider theme={theme}>
                    <SafeAreaProvider>
                      <NavigationContainer>
                        <StatusBar
                          backgroundColor={theme.colors.primary}
                          barStyle="light-content"
                        />
                        <AppNavigator />
                        <LoadingOverlay isVisible={false} />
                      </NavigationContainer>
                    </SafeAreaProvider>
                  </ThemeProvider>
                </FestivalProvider>
              </WalletProvider>
            </AuthProvider>
          </LocalizationProvider>
        </Provider>
      </GestureHandlerRootView>
    </ErrorBoundary>
  );
}