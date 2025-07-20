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

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Switch,
  Alert,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import AsyncStorage from '@react-native-async-storage/async-storage';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { useAppSelector, useAppDispatch } from '@store/index';
import { logout } from '@store/slices/authSlice';
import { clearWallet } from '@store/slices/walletSlice';

interface SettingItem {
  id: string;
  title: string;
  subtitle?: string;
  icon: string;
  iconColor?: string;
  type: 'navigate' | 'toggle' | 'action';
  value?: boolean;
  onPress?: () => void;
  dangerous?: boolean;
}

export const SettingsScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  
  // Settings state
  const [biometricEnabled, setBiometricEnabled] = useState(false);
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [festivalNotifications, setFestivalNotifications] = useState(true);
  const [autoBackup, setAutoBackup] = useState(false);
  const [developerMode, setDeveloperMode] = useState(false);
  
  const handleToggle = (settingId: string, value: boolean) => {
    switch (settingId) {
      case 'biometric':
        setBiometricEnabled(value);
        // Save to secure storage
        break;
      case 'notifications':
        setNotificationsEnabled(value);
        break;
      case 'festivalNotifications':
        setFestivalNotifications(value);
        break;
      case 'autoBackup':
        setAutoBackup(value);
        break;
      case 'developerMode':
        setDeveloperMode(value);
        break;
    }
  };
  
  const handleBackupWallet = () => {
    Alert.alert(
      'Backup Wallet',
      'This will show your recovery phrase. Make sure no one is looking at your screen.',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Continue', onPress: () => navigation.navigate('BackupWallet') },
      ]
    );
  };
  
  const handleResetApp = () => {
    Alert.alert(
      'Reset App',
      'This will delete all app data including your wallet. Make sure you have backed up your recovery phrase!',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Reset',
          style: 'destructive',
          onPress: async () => {
            await AsyncStorage.clear();
            dispatch(clearWallet());
            dispatch(logout());
          },
        },
      ]
    );
  };
  
  const settingSections = [
    {
      title: 'Account',
      items: [
        {
          id: 'profile',
          title: 'Profile',
          subtitle: dhanPataAddress || 'Set up your DhanPata ID',
          icon: 'account-circle',
          iconColor: COLORS.saffron,
          type: 'navigate',
          onPress: () => navigation.navigate('Profile'),
        },
        {
          id: 'backup',
          title: 'Backup Wallet',
          subtitle: 'Save your recovery phrase',
          icon: 'shield-key',
          iconColor: COLORS.green,
          type: 'action',
          onPress: handleBackupWallet,
        },
      ],
    },
    {
      title: 'Security',
      items: [
        {
          id: 'biometric',
          title: 'Biometric Authentication',
          subtitle: 'Use fingerprint or face to unlock',
          icon: 'fingerprint',
          iconColor: COLORS.navy,
          type: 'toggle',
          value: biometricEnabled,
        },
        {
          id: 'changePin',
          title: 'Change PIN',
          subtitle: 'Update your security PIN',
          icon: 'lock-reset',
          iconColor: COLORS.gray700,
          type: 'navigate',
          onPress: () => navigation.navigate('ChangePin'),
        },
        {
          id: 'autoLock',
          title: 'Auto-Lock',
          subtitle: 'Lock wallet after 5 minutes',
          icon: 'timer-outline',
          iconColor: COLORS.warning,
          type: 'navigate',
          onPress: () => navigation.navigate('AutoLock'),
        },
      ],
    },
    {
      title: 'Preferences',
      items: [
        {
          id: 'language',
          title: 'Language',
          subtitle: 'English',
          icon: 'translate',
          iconColor: COLORS.info,
          type: 'navigate',
          onPress: () => navigation.navigate('Language'),
        },
        {
          id: 'currency',
          title: 'Display Currency',
          subtitle: 'INR (₹)',
          icon: 'currency-inr',
          iconColor: COLORS.success,
          type: 'navigate',
          onPress: () => navigation.navigate('Currency'),
        },
        {
          id: 'notifications',
          title: 'Push Notifications',
          subtitle: 'Receive transaction alerts',
          icon: 'bell',
          iconColor: COLORS.festivalPrimary,
          type: 'toggle',
          value: notificationsEnabled,
        },
        {
          id: 'festivalNotifications',
          title: 'Festival Notifications',
          subtitle: 'Get notified about festival bonuses',
          icon: 'party-popper',
          iconColor: COLORS.festivalSecondary,
          type: 'toggle',
          value: festivalNotifications,
        },
      ],
    },
    {
      title: 'Advanced',
      items: [
        {
          id: 'network',
          title: 'Network',
          subtitle: 'DeshChain Mainnet',
          icon: 'web',
          iconColor: COLORS.gray600,
          type: 'navigate',
          onPress: () => navigation.navigate('Network'),
        },
        {
          id: 'autoBackup',
          title: 'Auto Backup',
          subtitle: 'Backup to encrypted cloud',
          icon: 'cloud-upload',
          iconColor: COLORS.sky,
          type: 'toggle',
          value: autoBackup,
        },
        {
          id: 'developerMode',
          title: 'Developer Mode',
          subtitle: 'Show advanced options',
          icon: 'code-tags',
          iconColor: COLORS.gray700,
          type: 'toggle',
          value: developerMode,
        },
      ],
    },
    {
      title: 'About',
      items: [
        {
          id: 'help',
          title: 'Help & Support',
          subtitle: 'Get help and report issues',
          icon: 'help-circle',
          iconColor: COLORS.info,
          type: 'navigate',
          onPress: () => navigation.navigate('Help'),
        },
        {
          id: 'terms',
          title: 'Terms of Service',
          icon: 'file-document',
          iconColor: COLORS.gray600,
          type: 'navigate',
          onPress: () => navigation.navigate('Terms'),
        },
        {
          id: 'privacy',
          title: 'Privacy Policy',
          icon: 'shield-lock',
          iconColor: COLORS.gray600,
          type: 'navigate',
          onPress: () => navigation.navigate('Privacy'),
        },
        {
          id: 'about',
          title: 'About DhanSetu',
          subtitle: 'Version 1.0.0',
          icon: 'information',
          iconColor: COLORS.saffron,
          type: 'navigate',
          onPress: () => navigation.navigate('About'),
        },
      ],
    },
    {
      title: 'Danger Zone',
      items: [
        {
          id: 'logout',
          title: 'Lock Wallet',
          subtitle: 'Require PIN to access',
          icon: 'logout',
          iconColor: COLORS.warning,
          type: 'action',
          onPress: () => {
            dispatch(logout());
            navigation.reset({
              index: 0,
              routes: [{ name: 'PinSetup' as any }],
            });
          },
        },
        {
          id: 'reset',
          title: 'Reset App',
          subtitle: 'Delete all data and wallet',
          icon: 'delete-forever',
          iconColor: COLORS.error,
          type: 'action',
          onPress: handleResetApp,
          dangerous: true,
        },
      ],
    },
  ];
  
  const renderSettingItem = (item: SettingItem) => {
    if (item.type === 'toggle') {
      return (
        <View key={item.id} style={styles.settingItem}>
          <View style={styles.settingLeft}>
            <View style={[styles.iconContainer, { backgroundColor: item.iconColor + '20' }]}>
              <Icon name={item.icon} size={24} color={item.iconColor} />
            </View>
            <View style={styles.settingInfo}>
              <Text style={styles.settingTitle}>{item.title}</Text>
              {item.subtitle && (
                <Text style={styles.settingSubtitle}>{item.subtitle}</Text>
              )}
            </View>
          </View>
          <Switch
            value={item.value}
            onValueChange={(value) => handleToggle(item.id, value)}
            trackColor={{ false: COLORS.gray300, true: COLORS.saffron }}
            thumbColor={item.value ? COLORS.white : COLORS.gray500}
          />
        </View>
      );
    }
    
    return (
      <TouchableOpacity
        key={item.id}
        style={[styles.settingItem, item.dangerous && styles.dangerousItem]}
        onPress={item.onPress}
        activeOpacity={0.7}
      >
        <View style={styles.settingLeft}>
          <View style={[styles.iconContainer, { backgroundColor: item.iconColor + '20' }]}>
            <Icon name={item.icon} size={24} color={item.iconColor} />
          </View>
          <View style={styles.settingInfo}>
            <Text style={[styles.settingTitle, item.dangerous && styles.dangerousText]}>
              {item.title}
            </Text>
            {item.subtitle && (
              <Text style={styles.settingSubtitle}>{item.subtitle}</Text>
            )}
          </View>
        </View>
        <Icon name="chevron-right" size={20} color={COLORS.gray400} />
      </TouchableOpacity>
    );
  };
  
  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Icon name="arrow-left" size={24} color={COLORS.gray900} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Settings</Text>
        <View style={{ width: 24 }} />
      </View>
      
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Developer Mode Banner */}
        {developerMode && (
          <View style={styles.developerBanner}>
            <Icon name="code-tags" size={20} color={COLORS.warning} />
            <Text style={styles.developerText}>Developer Mode Active</Text>
          </View>
        )}
        
        {/* Settings Sections */}
        {settingSections.map((section, index) => (
          <View key={section.title} style={styles.section}>
            <Text style={styles.sectionTitle}>{section.title}</Text>
            <View style={styles.sectionContent}>
              {section.items.map(renderSettingItem)}
            </View>
          </View>
        ))}
        
        {/* App Info */}
        <View style={styles.appInfo}>
          <Text style={styles.appInfoText}>DhanSetu by DeshChain Foundation</Text>
          <Text style={styles.appInfoSubtext}>Building the future of Indian finance</Text>
          <Text style={styles.appVersion}>v1.0.0 (Build 100)</Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  developerBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.warning + '10',
    paddingVertical: theme.spacing.sm,
    gap: theme.spacing.sm,
  },
  developerText: {
    fontSize: 14,
    color: COLORS.warning,
    fontWeight: '600',
  },
  section: {
    marginTop: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.gray600,
    textTransform: 'uppercase',
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.sm,
  },
  sectionContent: {
    backgroundColor: COLORS.white,
  },
  settingItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray100,
  },
  dangerousItem: {
    backgroundColor: COLORS.error + '05',
  },
  settingLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  iconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: theme.spacing.md,
  },
  settingInfo: {
    flex: 1,
  },
  settingTitle: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  dangerousText: {
    color: COLORS.error,
  },
  settingSubtitle: {
    fontSize: 14,
    color: COLORS.gray600,
    marginTop: 2,
  },
  appInfo: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
    marginTop: theme.spacing.xl,
  },
  appInfoText: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray700,
  },
  appInfoSubtext: {
    fontSize: 14,
    color: COLORS.gray600,
    marginTop: theme.spacing.xs,
  },
  appVersion: {
    fontSize: 12,
    color: COLORS.gray500,
    marginTop: theme.spacing.sm,
  },
});