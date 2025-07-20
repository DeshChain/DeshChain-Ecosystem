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
import { createStackNavigator } from '@react-navigation/stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { useSelector } from 'react-redux';
import { RootState } from '@store/index';

// Import screens
import { OnboardingScreen } from '@screens/onboarding/OnboardingScreen';
import { CreateWalletScreen } from '@screens/onboarding/CreateWalletScreen';
import { ImportWalletScreen } from '@screens/onboarding/ImportWalletScreen';
import { PinSetupScreen } from '@screens/onboarding/PinSetupScreen';

import { HomeScreen } from '@screens/home/HomeScreen';
import { WalletScreen } from '@screens/wallet/WalletScreen';
import { DexScreen } from '@screens/dex/DexScreen';
import { SikkebaazScreen } from '@screens/sikkebaaz/SikkebaazScreen';
import { SurakshaScreen } from '@screens/suraksha/SurakshaScreen';

import { SendScreen } from '@screens/send/SendScreen';
import { ReceiveScreen } from '@screens/receive/ReceiveScreen';
import { TransactionDetailsScreen } from '@screens/transaction/TransactionDetailsScreen';

import { SettingsScreen } from '@screens/settings/SettingsScreen';
import { ProfileScreen } from '@screens/profile/ProfileScreen';

// Import tab bar
import { CustomTabBar } from './CustomTabBar';

export type RootStackParamList = {
  // Onboarding
  Onboarding: undefined;
  CreateWallet: undefined;
  ImportWallet: undefined;
  PinSetup: { isNewWallet: boolean };
  
  // Main
  MainTabs: undefined;
  
  // Transaction
  Send: { coinType?: string };
  Receive: { coinType?: string };
  TransactionDetails: { txHash: string };
  
  // DEX
  CreateMoneyOrder: undefined;
  MoneyOrderDetails: { orderId: string };
  
  // Sikkebaaz
  CreateLaunch: undefined;
  LaunchDetails: { launchId: string };
  
  // Suraksha
  EnrollSuraksha: undefined;
  SurakshaDetails: { accountId: string };
  
  // Profile
  Profile: undefined;
  DhanPataSetup: undefined;
  
  // Settings
  Settings: undefined;
  Security: undefined;
  Language: undefined;
  About: undefined;
};

export type MainTabParamList = {
  Home: undefined;
  Wallet: undefined;
  DEX: undefined;
  Sikkebaaz: undefined;
  Suraksha: undefined;
};

const Stack = createStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<MainTabParamList>();

const MainTabs = () => {
  return (
    <Tab.Navigator
      tabBar={(props) => <CustomTabBar {...props} />}
      screenOptions={{
        headerShown: false,
      }}
    >
      <Tab.Screen name="Home" component={HomeScreen} />
      <Tab.Screen name="Wallet" component={WalletScreen} />
      <Tab.Screen name="DEX" component={DexScreen} />
      <Tab.Screen name="Sikkebaaz" component={SikkebaazScreen} />
      <Tab.Screen name="Suraksha" component={SurakshaScreen} />
    </Tab.Navigator>
  );
};

export const AppNavigator = () => {
  const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated);
  const hasWallet = useSelector((state: RootState) => state.wallet.hasWallet);

  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        cardStyleInterpolator: ({ current: { progress } }) => ({
          cardStyle: {
            opacity: progress,
          },
        }),
      }}
    >
      {!hasWallet ? (
        <>
          <Stack.Screen name="Onboarding" component={OnboardingScreen} />
          <Stack.Screen name="CreateWallet" component={CreateWalletScreen} />
          <Stack.Screen name="ImportWallet" component={ImportWalletScreen} />
          <Stack.Screen name="PinSetup" component={PinSetupScreen} />
        </>
      ) : !isAuthenticated ? (
        <Stack.Screen name="PinSetup" component={PinSetupScreen} />
      ) : (
        <>
          <Stack.Screen name="MainTabs" component={MainTabs} />
          
          {/* Transaction Screens */}
          <Stack.Screen name="Send" component={SendScreen} />
          <Stack.Screen name="Receive" component={ReceiveScreen} />
          <Stack.Screen name="TransactionDetails" component={TransactionDetailsScreen} />
          
          {/* DEX Screens */}
          <Stack.Screen name="CreateMoneyOrder" component={CreateMoneyOrderScreen} />
          <Stack.Screen name="MoneyOrderDetails" component={MoneyOrderDetailsScreen} />
          
          {/* Sikkebaaz Screens */}
          <Stack.Screen name="CreateLaunch" component={CreateLaunchScreen} />
          <Stack.Screen name="LaunchDetails" component={LaunchDetailsScreen} />
          
          {/* Suraksha Screens */}
          <Stack.Screen name="EnrollSuraksha" component={EnrollSurakshaScreen} />
          <Stack.Screen name="SurakshaDetails" component={SurakshaDetailsScreen} />
          
          {/* Profile & Settings */}
          <Stack.Screen name="Profile" component={ProfileScreen} />
          <Stack.Screen name="DhanPataSetup" component={DhanPataSetupScreen} />
          <Stack.Screen name="Settings" component={SettingsScreen} />
          <Stack.Screen name="Security" component={SecurityScreen} />
          <Stack.Screen name="Language" component={LanguageScreen} />
          <Stack.Screen name="About" component={AboutScreen} />
        </>
      )}
    </Stack.Navigator>
  );
};

// Placeholder imports (these will be created)
const CreateMoneyOrderScreen = () => null;
const MoneyOrderDetailsScreen = () => null;
const CreateLaunchScreen = () => null;
const LaunchDetailsScreen = () => null;
const EnrollSurakshaScreen = () => null;
const SurakshaDetailsScreen = () => null;
const DhanPataSetupScreen = () => null;
const SecurityScreen = () => null;
const LanguageScreen = () => null;
const AboutScreen = () => null;