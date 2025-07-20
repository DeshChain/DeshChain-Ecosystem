import React, { useEffect, useState } from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { createDrawerNavigator } from '@react-navigation/drawer';
import Icon from 'react-native-vector-icons/MaterialIcons';
import { useSelector } from 'react-redux';
import { Alert } from 'react-native';

import { RootState } from '../store';
import { BiometricManager } from '../services/BiometricManager';

// Auth Screens
import { WelcomeScreen } from '../screens/auth/WelcomeScreen';
import { SignInScreen } from '../screens/auth/SignInScreen';
import { SignUpScreen } from '../screens/auth/SignUpScreen';
import { WalletSetupScreen } from '../screens/auth/WalletSetupScreen';
import { BiometricSetupScreen } from '../screens/auth/BiometricSetupScreen';

// Main Screens
import { HomeScreen } from '../screens/main/HomeScreen';
import { WalletScreen } from '../screens/main/WalletScreen';
import { TradingScreen } from '../screens/main/TradingScreen';
import { P2PScreen } from '../screens/main/P2PScreen';
import { SevaMitraScreen } from '../screens/main/SevaMitraScreen';
import { ChatScreen } from '../screens/main/ChatScreen';
import { ProfileScreen } from '../screens/main/ProfileScreen';

// Transaction Screens
import { SendMoneyScreen } from '../screens/transactions/SendMoneyScreen';
import { ReceiveMoneyScreen } from '../screens/transactions/ReceiveMoneyScreen';
import { QRScannerScreen } from '../screens/transactions/QRScannerScreen';
import { TransactionDetailsScreen } from '../screens/transactions/TransactionDetailsScreen';
import { TransactionHistoryScreen } from '../screens/transactions/TransactionHistoryScreen';

// P2P Trading Screens
import { CreateOrderScreen } from '../screens/p2p/CreateOrderScreen';
import { OrderBookScreen } from '../screens/p2p/OrderBookScreen';
import { TradeDetailsScreen } from '../screens/p2p/TradeDetailsScreen';
import { ChatDetailsScreen } from '../screens/p2p/ChatDetailsScreen';

// Seva Mitra Screens
import { SevaMitraMapScreen } from '../screens/sevamitra/SevaMitraMapScreen';
import { SevaMitraDetailsScreen } from '../screens/sevamitra/SevaMitraDetailsScreen';
import { SevaMitraRegistrationScreen } from '../screens/sevamitra/SevaMitraRegistrationScreen';

// Settings Screens
import { SettingsScreen } from '../screens/settings/SettingsScreen';
import { SecurityScreen } from '../screens/settings/SecurityScreen';
import { LanguageScreen } from '../screens/settings/LanguageScreen';
import { NotificationScreen } from '../screens/settings/NotificationScreen';
import { AboutScreen } from '../screens/settings/AboutScreen';

// Onboarding
import { OnboardingScreen } from '../screens/onboarding/OnboardingScreen';

// Custom Drawer
import { CustomDrawerContent } from '../components/navigation/CustomDrawerContent';

export type RootStackParamList = {
  // Auth Flow
  Welcome: undefined;
  SignIn: undefined;
  SignUp: undefined;
  WalletSetup: { mnemonic?: string };
  BiometricSetup: undefined;
  Onboarding: undefined;
  
  // Main App
  MainApp: undefined;
  
  // Transaction Flow
  SendMoney: { recipientAddress?: string };
  ReceiveMoney: undefined;
  QRScanner: { onScan: (data: string) => void };
  TransactionDetails: { transactionId: string };
  TransactionHistory: undefined;
  
  // P2P Trading
  CreateOrder: { orderType?: 'buy' | 'sell' };
  OrderBook: undefined;
  TradeDetails: { tradeId: string };
  ChatDetails: { conversationId: string };
  
  // Seva Mitra
  SevaMitraMap: undefined;
  SevaMitraDetails: { mitraId: string };
  SevaMitraRegistration: undefined;
  
  // Settings
  Settings: undefined;
  Security: undefined;
  Language: undefined;
  Notifications: undefined;
  About: undefined;
};

export type TabParamList = {
  Home: undefined;
  Wallet: undefined;
  Trading: undefined;
  P2P: undefined;
  SevaMitra: undefined;
  Chat: undefined;
  Profile: undefined;
};

const Stack = createNativeStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<TabParamList>();
const Drawer = createDrawerNavigator();

// Bottom Tab Navigator
const TabNavigator: React.FC = () => {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: string;

          switch (route.name) {
            case 'Home':
              iconName = 'home';
              break;
            case 'Wallet':
              iconName = 'account-balance-wallet';
              break;
            case 'Trading':
              iconName = 'trending-up';
              break;
            case 'P2P':
              iconName = 'swap-horiz';
              break;
            case 'SevaMitra':
              iconName = 'location-on';
              break;
            case 'Chat':
              iconName = 'chat';
              break;
            case 'Profile':
              iconName = 'person';
              break;
            default:
              iconName = 'home';
          }

          return <Icon name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: '#FF9933',
        tabBarInactiveTintColor: '#757575',
        tabBarStyle: {
          backgroundColor: '#FFFFFF',
          borderTopColor: '#E0E0E0',
          borderTopWidth: 1,
          height: 60,
          paddingBottom: 8,
          paddingTop: 8,
        },
        headerShown: false,
      })}
    >
      <Tab.Screen 
        name="Home" 
        component={HomeScreen}
        options={{ title: 'होम' }} // Hindi for Home
      />
      <Tab.Screen 
        name="Wallet" 
        component={WalletScreen}
        options={{ title: 'वॉलेट' }} // Hindi for Wallet
      />
      <Tab.Screen 
        name="Trading" 
        component={TradingScreen}
        options={{ title: 'ट्रेडिंग' }} // Hindi for Trading
      />
      <Tab.Screen 
        name="P2P" 
        component={P2PScreen}
        options={{ title: 'P2P' }}
      />
      <Tab.Screen 
        name="SevaMitra" 
        component={SevaMitraScreen}
        options={{ title: 'सेवा मित्र' }} // Hindi for Seva Mitra
      />
      <Tab.Screen 
        name="Chat" 
        component={ChatScreen}
        options={{ title: 'चैट' }} // Hindi for Chat
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen}
        options={{ title: 'प्रोफ़ाइल' }} // Hindi for Profile
      />
    </Tab.Navigator>
  );
};

// Drawer Navigator for main app
const DrawerNavigator: React.FC = () => {
  return (
    <Drawer.Navigator
      drawerContent={(props) => <CustomDrawerContent {...props} />}
      screenOptions={{
        headerShown: false,
        drawerStyle: {
          backgroundColor: '#FFFFFF',
          width: 280,
        },
        drawerActiveTintColor: '#FF9933',
        drawerInactiveTintColor: '#757575',
      }}
    >
      <Drawer.Screen name="MainTabs" component={TabNavigator} />
    </Drawer.Navigator>
  );
};

// Auth Navigator
const AuthNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="Welcome" component={WelcomeScreen} />
      <Stack.Screen name="SignIn" component={SignInScreen} />
      <Stack.Screen name="SignUp" component={SignUpScreen} />
      <Stack.Screen name="WalletSetup" component={WalletSetupScreen} />
      <Stack.Screen name="BiometricSetup" component={BiometricSetupScreen} />
      <Stack.Screen name="Onboarding" component={OnboardingScreen} />
    </Stack.Navigator>
  );
};

// Main App Navigator
export const AppNavigator: React.FC = () => {
  const { isAuthenticated, hasCompletedOnboarding, biometricEnabled } = useSelector(
    (state: RootState) => state.auth
  );
  const [biometricVerified, setBiometricVerified] = useState(false);
  const [isCheckingBiometric, setIsCheckingBiometric] = useState(true);

  useEffect(() => {
    const checkBiometricAuth = async () => {
      if (isAuthenticated && biometricEnabled && !biometricVerified) {
        try {
          const isAvailable = await BiometricManager.isBiometricAvailable();
          if (isAvailable) {
            const result = await BiometricManager.authenticate({
              title: 'DeshChain Authentication',
              subtitle: 'Verify your identity to access the app',
              description: 'Place your finger on the sensor or look at the camera',
              fallbackLabel: 'Use PIN',
              negativeLabel: 'Cancel',
            });

            if (result.success) {
              setBiometricVerified(true);
            } else {
              Alert.alert(
                'Authentication Failed',
                'Please try again or use alternative authentication method',
                [
                  { text: 'Retry', onPress: checkBiometricAuth },
                  { text: 'Use PIN', onPress: () => setBiometricVerified(true) },
                ]
              );
            }
          } else {
            setBiometricVerified(true);
          }
        } catch (error) {
          console.error('Biometric authentication error:', error);
          setBiometricVerified(true);
        }
      } else {
        setBiometricVerified(true);
      }
      setIsCheckingBiometric(false);
    };

    if (isAuthenticated) {
      checkBiometricAuth();
    } else {
      setIsCheckingBiometric(false);
    }
  }, [isAuthenticated, biometricEnabled, biometricVerified]);

  // Show loading while checking biometric
  if (isCheckingBiometric) {
    return <WelcomeScreen />; // Show splash while checking
  }

  // Show auth flow if not authenticated
  if (!isAuthenticated || (biometricEnabled && !biometricVerified)) {
    return <AuthNavigator />;
  }

  // Show onboarding if first time user
  if (!hasCompletedOnboarding) {
    return (
      <Stack.Navigator screenOptions={{ headerShown: false }}>
        <Stack.Screen name="Onboarding" component={OnboardingScreen} />
      </Stack.Navigator>
    );
  }

  // Main app navigation
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="MainApp" component={DrawerNavigator} />
      
      {/* Transaction Screens */}
      <Stack.Screen
        name="SendMoney"
        component={SendMoneyScreen}
        options={{
          headerShown: true,
          title: 'Send Money',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="ReceiveMoney"
        component={ReceiveMoneyScreen}
        options={{
          headerShown: true,
          title: 'Receive Money',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="QRScanner"
        component={QRScannerScreen}
        options={{
          headerShown: true,
          title: 'Scan QR Code',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="TransactionDetails"
        component={TransactionDetailsScreen}
        options={{
          headerShown: true,
          title: 'Transaction Details',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="TransactionHistory"
        component={TransactionHistoryScreen}
        options={{
          headerShown: true,
          title: 'Transaction History',
          headerBackTitle: 'Back',
        }}
      />
      
      {/* P2P Trading Screens */}
      <Stack.Screen
        name="CreateOrder"
        component={CreateOrderScreen}
        options={{
          headerShown: true,
          title: 'Create Order',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="OrderBook"
        component={OrderBookScreen}
        options={{
          headerShown: true,
          title: 'Order Book',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="TradeDetails"
        component={TradeDetailsScreen}
        options={{
          headerShown: true,
          title: 'Trade Details',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="ChatDetails"
        component={ChatDetailsScreen}
        options={{
          headerShown: true,
          title: 'Trade Chat',
          headerBackTitle: 'Back',
        }}
      />
      
      {/* Seva Mitra Screens */}
      <Stack.Screen
        name="SevaMitraMap"
        component={SevaMitraMapScreen}
        options={{
          headerShown: true,
          title: 'Find Seva Mitra',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="SevaMitraDetails"
        component={SevaMitraDetailsScreen}
        options={{
          headerShown: true,
          title: 'Seva Mitra Details',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="SevaMitraRegistration"
        component={SevaMitraRegistrationScreen}
        options={{
          headerShown: true,
          title: 'Become Seva Mitra',
          headerBackTitle: 'Back',
        }}
      />
      
      {/* Settings Screens */}
      <Stack.Screen
        name="Settings"
        component={SettingsScreen}
        options={{
          headerShown: true,
          title: 'Settings',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="Security"
        component={SecurityScreen}
        options={{
          headerShown: true,
          title: 'Security',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="Language"
        component={LanguageScreen}
        options={{
          headerShown: true,
          title: 'Language / भाषा',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="Notifications"
        component={NotificationScreen}
        options={{
          headerShown: true,
          title: 'Notifications',
          headerBackTitle: 'Back',
        }}
      />
      <Stack.Screen
        name="About"
        component={AboutScreen}
        options={{
          headerShown: true,
          title: 'About DeshChain',
          headerBackTitle: 'Back',
        }}
      />
    </Stack.Navigator>
  );
};