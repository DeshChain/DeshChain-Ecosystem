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

import React, { useEffect, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
  TouchableOpacity,
  Dimensions,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

import { useAppSelector, useAppDispatch } from '@store/index';
import { COLORS, theme } from '@constants/theme';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { CulturalButton } from '@components/common/CulturalButton';
import { useFestival } from '@contexts/FestivalContext';
import { RootStackParamList } from '@navigation/AppNavigator';

const { width } = Dimensions.get('window');

type NavigationProp = StackNavigationProp<RootStackParamList, 'MainTabs'>;

export const HomeScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const dispatch = useAppDispatch();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival } = useFestival();
  const [refreshing, setRefreshing] = useState(false);
  const [dailyQuote, setDailyQuote] = useState({
    quote: 'वसुधैव कुटुम्बकम्',
    translation: 'The world is one family',
    author: 'Maha Upanishad',
  });

  const onRefresh = async () => {
    setRefreshing(true);
    // Refresh balances and data
    await new Promise((resolve) => setTimeout(resolve, 1500));
    setRefreshing(false);
  };

  const formatAddress = (address: string) => {
    if (!address) return '';
    return `${address.slice(0, 10)}...${address.slice(-8)}`;
  };

  const quickActions = [
    {
      id: 'send',
      title: 'Send',
      icon: 'send',
      color: COLORS.saffron,
      onPress: () => navigation.navigate('Send', {}),
    },
    {
      id: 'receive',
      title: 'Receive',
      icon: 'qrcode',
      color: COLORS.green,
      onPress: () => navigation.navigate('Receive', {}),
    },
    {
      id: 'dex',
      title: 'Money Order',
      icon: 'swap-horizontal',
      color: COLORS.navy,
      onPress: () => navigation.navigate('CreateMoneyOrder'),
    },
    {
      id: 'sikkebaaz',
      title: 'Launch Token',
      icon: 'rocket-launch',
      color: COLORS.festivalPrimary,
      onPress: () => navigation.navigate('CreateLaunch'),
    },
  ];

  const features = [
    {
      id: 'suraksha',
      title: 'Gram Suraksha',
      subtitle: '50% guaranteed returns',
      icon: 'shield-check',
      gradient: [COLORS.green, COLORS.darkGreen],
    },
    {
      id: 'krishi',
      title: 'Krishi Mitra',
      subtitle: 'Agricultural loans at 6-9%',
      icon: 'sprout',
      gradient: [COLORS.saffron, COLORS.darkSaffron],
    },
    {
      id: 'kshetra',
      title: 'Kshetra Coins',
      subtitle: 'Local community tokens',
      icon: 'map-marker-radius',
      gradient: [COLORS.navy, '#1E3A8A'],
    },
  ];

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            colors={[COLORS.saffron]}
          />
        }
      >
        {/* Header */}
        <LinearGradient
          colors={currentFestival ? theme.gradients.festivalGradient : theme.gradients.indianFlag}
          style={styles.header}
        >
          <View style={styles.headerContent}>
            <View>
              <Text style={styles.greeting}>Namaste! 🙏</Text>
              <Text style={styles.address}>
                {dhanPataAddress || formatAddress(currentAddress || '')}
              </Text>
            </View>
            <TouchableOpacity
              onPress={() => navigation.navigate('Profile')}
              style={styles.profileButton}
            >
              <Icon name="account-circle" size={40} color={COLORS.white} />
            </TouchableOpacity>
          </View>

          {/* Balance Card */}
          <View style={styles.balanceCard}>
            <Text style={styles.balanceLabel}>Total Balance</Text>
            <GradientText style={styles.balanceAmount}>
              ₹ 1,23,456.78
            </GradientText>
            <Text style={styles.namoBalance}>≈ 12,345 NAMO</Text>
          </View>

          {/* Festival Banner */}
          {currentFestival && (
            <View style={styles.festivalBanner}>
              <Icon name="party-popper" size={24} color={COLORS.festivalAccent} />
              <Text style={styles.festivalText}>
                {currentFestival.traditionalGreeting} {currentFestival.bonusRate}% bonus active!
              </Text>
            </View>
          )}
        </LinearGradient>

        {/* Quick Actions */}
        <View style={styles.quickActions}>
          {quickActions.map((action) => (
            <TouchableOpacity
              key={action.id}
              style={styles.actionButton}
              onPress={action.onPress}
              activeOpacity={0.8}
            >
              <View style={[styles.actionIcon, { backgroundColor: action.color }]}>
                <Icon name={action.icon} size={24} color={COLORS.white} />
              </View>
              <Text style={styles.actionText}>{action.title}</Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Daily Quote */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>आज का विचार</Text>
          <QuoteCard
            quote={dailyQuote.quote}
            author={dailyQuote.author}
            language="sa"
            variant={currentFestival ? 'festival' : 'default'}
          />
          <Text style={styles.translation}>{dailyQuote.translation}</Text>
        </View>

        {/* Features */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>DeshChain Products</Text>
          {features.map((feature) => (
            <TouchableOpacity
              key={feature.id}
              style={styles.featureCard}
              activeOpacity={0.8}
            >
              <LinearGradient
                colors={feature.gradient}
                style={styles.featureGradient}
                start={{ x: 0, y: 0 }}
                end={{ x: 1, y: 0 }}
              >
                <Icon name={feature.icon} size={32} color={COLORS.white} />
                <View style={styles.featureContent}>
                  <Text style={styles.featureTitle}>{feature.title}</Text>
                  <Text style={styles.featureSubtitle}>{feature.subtitle}</Text>
                </View>
                <Icon name="chevron-right" size={24} color={COLORS.white} />
              </LinearGradient>
            </TouchableOpacity>
          ))}
        </View>

        {/* Cultural Message */}
        <View style={styles.culturalMessage}>
          <Text style={styles.culturalText}>
            Building India's financial future with blockchain technology
          </Text>
          <Text style={styles.culturalSubtext}>
            जय हिंद! 🇮🇳
          </Text>
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
    paddingTop: theme.spacing.lg,
    paddingBottom: theme.spacing.xl,
    borderBottomLeftRadius: theme.borderRadius.xl,
    borderBottomRightRadius: theme.borderRadius.xl,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.lg,
  },
  greeting: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  address: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
    marginTop: 4,
  },
  profileButton: {
    padding: theme.spacing.sm,
  },
  balanceCard: {
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    marginHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.lg,
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    alignItems: 'center',
  },
  balanceLabel: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  balanceAmount: {
    fontSize: 36,
    fontWeight: 'bold',
    marginVertical: theme.spacing.sm,
  },
  namoBalance: {
    fontSize: 16,
    color: COLORS.white,
    opacity: 0.9,
  },
  festivalBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: theme.spacing.md,
    paddingHorizontal: theme.spacing.lg,
  },
  festivalText: {
    fontSize: 14,
    color: COLORS.white,
    marginLeft: theme.spacing.sm,
    fontWeight: '600',
  },
  quickActions: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingVertical: theme.spacing.lg,
    paddingHorizontal: theme.spacing.md,
    marginTop: -theme.spacing.lg,
  },
  actionButton: {
    alignItems: 'center',
  },
  actionIcon: {
    width: 56,
    height: 56,
    borderRadius: 28,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.medium,
  },
  actionText: {
    fontSize: 12,
    color: COLORS.gray700,
    marginTop: theme.spacing.sm,
  },
  section: {
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.md,
  },
  translation: {
    fontSize: 14,
    color: COLORS.gray600,
    textAlign: 'center',
    marginTop: theme.spacing.sm,
    fontStyle: 'italic',
  },
  featureCard: {
    marginBottom: theme.spacing.md,
  },
  featureGradient: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    ...theme.shadows.medium,
  },
  featureContent: {
    flex: 1,
    marginLeft: theme.spacing.md,
  },
  featureTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  featureSubtitle: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
    marginTop: 2,
  },
  culturalMessage: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
    paddingHorizontal: theme.spacing.lg,
  },
  culturalText: {
    fontSize: 16,
    color: COLORS.gray700,
    textAlign: 'center',
    fontStyle: 'italic',
  },
  culturalSubtext: {
    fontSize: 18,
    color: COLORS.saffron,
    marginTop: theme.spacing.sm,
    fontWeight: 'bold',
  },
});