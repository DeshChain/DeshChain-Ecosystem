import React, { useEffect } from 'react';
import { StatusBar, LogBox } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { Provider as PaperProvider } from 'react-native-paper';
import { Provider as ReduxProvider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import SplashScreen from 'react-native-splash-screen';
import { enableScreens } from 'react-native-screens';
import 'react-native-get-random-values';

import { store, persistor } from './src/store';
import { AppNavigator } from './src/navigation/AppNavigator';
import { theme } from './src/theme';
import { LoadingScreen } from './src/components/LoadingScreen';
import { ErrorBoundary } from './src/components/ErrorBoundary';
import { NotificationManager } from './src/services/NotificationManager';
import { BiometricManager } from './src/services/BiometricManager';
import { CrashReportingService } from './src/services/CrashReportingService';

// Enable native screens for better performance
enableScreens();

// Ignore specific warnings for development
LogBox.ignoreLogs([
  'Non-serializable values were found in the navigation state',
  'Require cycle:',
]);

const App: React.FC = () => {
  useEffect(() => {
    const initializeApp = async () => {
      try {
        // Initialize services
        await NotificationManager.initialize();
        await BiometricManager.initialize();
        await CrashReportingService.initialize();
        
        // Hide splash screen after initialization
        setTimeout(() => {
          SplashScreen.hide();
        }, 1500);
      } catch (error) {
        console.error('App initialization error:', error);
        SplashScreen.hide();
      }
    };

    initializeApp();
  }, []);

  return (
    <ErrorBoundary>
      <ReduxProvider store={store}>
        <PersistGate loading={<LoadingScreen />} persistor={persistor}>
          <PaperProvider theme={theme}>
            <NavigationContainer theme={theme}>
              <StatusBar
                barStyle="light-content"
                backgroundColor={theme.colors.primary}
                translucent={false}
              />
              <AppNavigator />
            </NavigationContainer>
          </PaperProvider>
        </PersistGate>
      </ReduxProvider>
    </ErrorBoundary>
  );
};

export default App;