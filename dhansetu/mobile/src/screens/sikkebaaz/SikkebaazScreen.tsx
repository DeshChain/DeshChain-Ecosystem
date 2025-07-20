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

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  FlatList,
  Image,
  TextInput,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import Animated, {
  useAnimatedStyle,
  withSpring,
  interpolate,
} from 'react-native-reanimated';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector } from '@store/index';
import { useFestival } from '@contexts/FestivalContext';

interface TokenLaunch {
  id: string;
  tokenName: string;
  tokenSymbol: string;
  description: string;
  imageUrl?: string;
  creator: string;
  creatorPincode: string;
  targetAmount: string;
  raisedAmount: string;
  progress: number;
  status: 'active' | 'successful' | 'failed' | 'vetoed';
  launchDate: string;
  endDate: string;
  participants: number;
  culturalCategory: string;
  antiPumpConfig: {
    maxWalletPercent: number;
    tradingDelayHours: number;
    liquidityLockMonths: number;
  };
  communityVeto?: {
    active: boolean;
    votes: number;
    threshold: number;
  };
}

export const SikkebaazScreen: React.FC = () => {
  const navigation = useNavigation();
  const { currentAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival, festivalBonusRate } = useFestival();
  
  const [activeTab, setActiveTab] = useState<'trending' | 'new' | 'myTokens'>('trending');
  const [launches, setLaunches] = useState<TokenLaunch[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const categories = [
    { id: 'all', name: 'All', icon: 'all-inclusive', color: COLORS.saffron },
    { id: 'bollywood', name: 'Bollywood', icon: 'movie', color: '#E91E63' },
    { id: 'cricket', name: 'Cricket', icon: 'cricket', color: '#4CAF50' },
    { id: 'festival', name: 'Festival', icon: 'party-popper', color: '#FF9800' },
    { id: 'regional', name: 'Regional', icon: 'map-marker', color: '#2196F3' },
    { id: 'meme', name: 'Meme', icon: 'emoticon-lol', color: '#9C27B0' },
  ];

  useEffect(() => {
    fetchLaunches();
  }, [activeTab, selectedCategory]);

  const fetchLaunches = async () => {
    setLoading(true);
    // Simulate API call
    setTimeout(() => {
      setLaunches(getMockLaunches());
      setLoading(false);
    }, 1000);
  };

  const filteredLaunches = launches.filter(launch => {
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      return (
        launch.tokenName.toLowerCase().includes(query) ||
        launch.tokenSymbol.toLowerCase().includes(query) ||
        launch.description.toLowerCase().includes(query)
      );
    }
    return true;
  });

  const renderCategoryFilter = () => (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      style={styles.categoryScroll}
      contentContainerStyle={styles.categoryContainer}
    >
      {categories.map((category) => (
        <TouchableOpacity
          key={category.id}
          style={[
            styles.categoryButton,
            selectedCategory === category.id && styles.categoryButtonActive,
          ]}
          onPress={() => setSelectedCategory(category.id)}
          activeOpacity={0.8}
        >
          <Icon
            name={category.icon}
            size={20}
            color={selectedCategory === category.id ? COLORS.white : category.color}
          />
          <Text
            style={[
              styles.categoryText,
              selectedCategory === category.id && styles.categoryTextActive,
            ]}
          >
            {category.name}
          </Text>
        </TouchableOpacity>
      ))}
    </ScrollView>
  );

  const renderLaunchCard = ({ item }: { item: TokenLaunch }) => {
    const progressAnimatedStyle = useAnimatedStyle(() => ({
      width: withSpring(`${item.progress}%`),
    }));

    return (
      <TouchableOpacity
        style={styles.launchCard}
        onPress={() => navigation.navigate('LaunchDetails', { launchId: item.id })}
        activeOpacity={0.9}
      >
        {/* Header with image and basic info */}
        <View style={styles.cardHeader}>
          <View style={styles.tokenInfo}>
            {item.imageUrl ? (
              <Image source={{ uri: item.imageUrl }} style={styles.tokenImage} />
            ) : (
              <View style={[styles.tokenImage, styles.tokenImagePlaceholder]}>
                <Text style={styles.tokenInitial}>
                  {item.tokenSymbol.charAt(0)}
                </Text>
              </View>
            )}
            <View style={styles.tokenDetails}>
              <Text style={styles.tokenName}>{item.tokenName}</Text>
              <Text style={styles.tokenSymbol}>${item.tokenSymbol}</Text>
            </View>
          </View>
          
          {/* Status badge */}
          <View style={[styles.statusBadge, styles[`status_${item.status}`]]}>
            <Text style={styles.statusText}>{item.status.toUpperCase()}</Text>
          </View>
        </View>

        {/* Description */}
        <Text style={styles.description} numberOfLines={2}>
          {item.description}
        </Text>

        {/* Cultural category */}
        <View style={styles.categoryTag}>
          <Icon name="tag" size={14} color={COLORS.gray600} />
          <Text style={styles.categoryTagText}>{item.culturalCategory}</Text>
          {item.creatorPincode && (
            <>
              <Icon name="map-marker" size={14} color={COLORS.gray600} style={{ marginLeft: 8 }} />
              <Text style={styles.categoryTagText}>{item.creatorPincode}</Text>
            </>
          )}
        </View>

        {/* Progress */}
        <View style={styles.progressSection}>
          <View style={styles.progressHeader}>
            <Text style={styles.raisedAmount}>₹{item.raisedAmount}</Text>
            <Text style={styles.targetAmount}>/ ₹{item.targetAmount}</Text>
          </View>
          <View style={styles.progressBar}>
            <Animated.View style={[styles.progressFill, progressAnimatedStyle]} />
          </View>
          <View style={styles.progressStats}>
            <Text style={styles.progressStat}>{item.progress}% funded</Text>
            <Text style={styles.progressStat}>{item.participants} participants</Text>
          </View>
        </View>

        {/* Anti-pump features */}
        <View style={styles.antiPumpSection}>
          <View style={styles.antiPumpItem}>
            <Icon name="wallet" size={16} color={COLORS.green} />
            <Text style={styles.antiPumpText}>
              Max {item.antiPumpConfig.maxWalletPercent}% per wallet
            </Text>
          </View>
          <View style={styles.antiPumpItem}>
            <Icon name="clock-outline" size={16} color={COLORS.saffron} />
            <Text style={styles.antiPumpText}>
              {item.antiPumpConfig.tradingDelayHours}h trading delay
            </Text>
          </View>
          <View style={styles.antiPumpItem}>
            <Icon name="lock" size={16} color={COLORS.navy} />
            <Text style={styles.antiPumpText}>
              {item.antiPumpConfig.liquidityLockMonths}mo liquidity lock
            </Text>
          </View>
        </View>

        {/* Community veto indicator */}
        {item.communityVeto?.active && (
          <View style={styles.vetoWarning}>
            <Icon name="alert-circle" size={16} color={COLORS.error} />
            <Text style={styles.vetoText}>
              Community veto in progress ({item.communityVeto.votes}/{item.communityVeto.threshold} votes)
            </Text>
          </View>
        )}

        {/* Festival bonus */}
        {currentFestival && (
          <View style={styles.festivalBonus}>
            <Icon name="party-popper" size={16} color={COLORS.festivalPrimary} />
            <Text style={styles.festivalBonusText}>
              {currentFestival.name} bonus: Extra {currentFestival.bonusRate}% tokens!
            </Text>
          </View>
        )}
      </TouchableOpacity>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <LinearGradient
        colors={currentFestival ? theme.gradients.festivalGradient : [COLORS.saffron, COLORS.darkSaffron]}
        style={styles.header}
      >
        <View style={styles.headerContent}>
          <View>
            <Text style={styles.headerTitle}>Sikkebaaz</Text>
            <Text style={styles.headerSubtitle}>
              Desi Memecoin Launchpad
            </Text>
          </View>
          <TouchableOpacity
            style={styles.createButton}
            onPress={() => navigation.navigate('CreateLaunch')}
          >
            <Icon name="rocket-launch" size={24} color={COLORS.saffron} />
          </TouchableOpacity>
        </View>

        {/* Stats */}
        <View style={styles.statsContainer}>
          <View style={styles.statItem}>
            <Text style={styles.statValue}>₹12.5 Cr</Text>
            <Text style={styles.statLabel}>Total Raised</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            <Text style={styles.statValue}>156</Text>
            <Text style={styles.statLabel}>Active Launches</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.statItem}>
            <Text style={styles.statValue}>89%</Text>
            <Text style={styles.statLabel}>Success Rate</Text>
          </View>
        </View>
      </LinearGradient>

      {/* Search */}
      <View style={styles.searchContainer}>
        <Icon name="magnify" size={20} color={COLORS.gray600} />
        <TextInput
          style={styles.searchInput}
          placeholder="Search tokens..."
          value={searchQuery}
          onChangeText={setSearchQuery}
          placeholderTextColor={COLORS.gray400}
        />
      </View>

      {/* Category Filter */}
      {renderCategoryFilter()}

      {/* Tabs */}
      <View style={styles.tabs}>
        {(['trending', 'new', 'myTokens'] as const).map((tab) => (
          <TouchableOpacity
            key={tab}
            style={[styles.tab, activeTab === tab && styles.activeTab]}
            onPress={() => setActiveTab(tab)}
          >
            <Icon
              name={tab === 'trending' ? 'fire' : tab === 'new' ? 'new-box' : 'account-star'}
              size={20}
              color={activeTab === tab ? COLORS.saffron : COLORS.gray600}
            />
            <Text style={[
              styles.tabText,
              activeTab === tab && styles.activeTabText,
            ]}>
              {tab === 'trending' ? 'Trending' : tab === 'new' ? 'New' : 'My Tokens'}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {/* Launch List */}
      {loading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={COLORS.saffron} />
        </View>
      ) : (
        <FlatList
          data={filteredLaunches}
          renderItem={renderLaunchCard}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.launchList}
          showsVerticalScrollIndicator={false}
          ListEmptyComponent={
            <View style={styles.emptyState}>
              <Icon name="rocket-launch-outline" size={64} color={COLORS.gray400} />
              <Text style={styles.emptyText}>No launches found</Text>
              <CulturalButton
                title="Create First Launch"
                onPress={() => navigation.navigate('CreateLaunch')}
                variant="outline"
                size="medium"
                style={styles.emptyButton}
              />
            </View>
          }
        />
      )}
    </SafeAreaView>
  );
};

// Mock data generator
const getMockLaunches = (): TokenLaunch[] => [
  {
    id: '1',
    tokenName: 'Bollywood Coin',
    tokenSymbol: 'BOLLY',
    description: 'The official memecoin for Bollywood fans worldwide! Earn rewards for movie predictions.',
    imageUrl: 'https://via.placeholder.com/100',
    creator: 'creator@dhan',
    creatorPincode: '400001',
    targetAmount: '1000000',
    raisedAmount: '750000',
    progress: 75,
    status: 'active',
    launchDate: new Date().toISOString(),
    endDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
    participants: 234,
    culturalCategory: 'Bollywood',
    antiPumpConfig: {
      maxWalletPercent: 2,
      tradingDelayHours: 24,
      liquidityLockMonths: 6,
    },
  },
  {
    id: '2',
    tokenName: 'Cricket Champions',
    tokenSymbol: 'WICKET',
    description: 'Support your favorite cricket team and earn rewards for match predictions!',
    creator: 'sports@dhan',
    creatorPincode: '110001',
    targetAmount: '500000',
    raisedAmount: '450000',
    progress: 90,
    status: 'active',
    launchDate: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    endDate: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString(),
    participants: 567,
    culturalCategory: 'Cricket',
    antiPumpConfig: {
      maxWalletPercent: 1.5,
      tradingDelayHours: 48,
      liquidityLockMonths: 12,
    },
    communityVeto: {
      active: true,
      votes: 45,
      threshold: 100,
    },
  },
];

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  header: {
    paddingBottom: theme.spacing.lg,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.lg,
    paddingTop: theme.spacing.lg,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  headerSubtitle: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  createButton: {
    width: 48,
    height: 48,
    backgroundColor: COLORS.white,
    borderRadius: 24,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.medium,
  },
  statsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginTop: theme.spacing.lg,
    paddingHorizontal: theme.spacing.lg,
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  statLabel: {
    fontSize: 12,
    color: COLORS.white,
    opacity: 0.8,
    marginTop: 2,
  },
  statDivider: {
    width: 1,
    height: 30,
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    marginHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.md,
    paddingHorizontal: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    height: 44,
  },
  searchInput: {
    flex: 1,
    fontSize: 16,
    color: COLORS.gray900,
    marginLeft: theme.spacing.sm,
  },
  categoryScroll: {
    marginTop: theme.spacing.md,
  },
  categoryContainer: {
    paddingHorizontal: theme.spacing.lg,
    gap: theme.spacing.sm,
  },
  categoryButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.full,
    backgroundColor: COLORS.gray100,
    marginRight: theme.spacing.sm,
  },
  categoryButtonActive: {
    backgroundColor: COLORS.saffron,
  },
  categoryText: {
    fontSize: 14,
    color: COLORS.gray700,
    marginLeft: theme.spacing.xs,
  },
  categoryTextActive: {
    color: COLORS.white,
    fontWeight: '600',
  },
  tabs: {
    flexDirection: 'row',
    paddingHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.md,
    marginBottom: theme.spacing.sm,
  },
  tab: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: theme.spacing.sm,
    gap: theme.spacing.xs,
  },
  activeTab: {
    borderBottomWidth: 2,
    borderBottomColor: COLORS.saffron,
  },
  tabText: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  activeTabText: {
    color: COLORS.saffron,
    fontWeight: '600',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  launchList: {
    padding: theme.spacing.lg,
  },
  launchCard: {
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    marginBottom: theme.spacing.md,
    ...theme.shadows.medium,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.md,
  },
  tokenInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  tokenImage: {
    width: 48,
    height: 48,
    borderRadius: 24,
    marginRight: theme.spacing.md,
  },
  tokenImagePlaceholder: {
    backgroundColor: COLORS.saffron,
    justifyContent: 'center',
    alignItems: 'center',
  },
  tokenInitial: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  tokenDetails: {
    flex: 1,
  },
  tokenName: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  tokenSymbol: {
    fontSize: 14,
    color: COLORS.gray600,
    marginTop: 2,
  },
  statusBadge: {
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: 4,
    borderRadius: theme.borderRadius.sm,
  },
  status_active: {
    backgroundColor: COLORS.success + '20',
  },
  status_successful: {
    backgroundColor: COLORS.info + '20',
  },
  status_failed: {
    backgroundColor: COLORS.error + '20',
  },
  status_vetoed: {
    backgroundColor: COLORS.warning + '20',
  },
  statusText: {
    fontSize: 10,
    fontWeight: 'bold',
    color: COLORS.gray800,
  },
  description: {
    fontSize: 14,
    color: COLORS.gray700,
    lineHeight: 20,
    marginBottom: theme.spacing.sm,
  },
  categoryTag: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  categoryTagText: {
    fontSize: 12,
    color: COLORS.gray600,
    marginLeft: 4,
  },
  progressSection: {
    marginBottom: theme.spacing.md,
  },
  progressHeader: {
    flexDirection: 'row',
    alignItems: 'baseline',
    marginBottom: theme.spacing.sm,
  },
  raisedAmount: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.saffron,
  },
  targetAmount: {
    fontSize: 16,
    color: COLORS.gray600,
    marginLeft: 4,
  },
  progressBar: {
    height: 8,
    backgroundColor: COLORS.gray200,
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: theme.spacing.sm,
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.saffron,
    borderRadius: 4,
  },
  progressStats: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  progressStat: {
    fontSize: 12,
    color: COLORS.gray600,
  },
  antiPumpSection: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: theme.spacing.sm,
    marginBottom: theme.spacing.sm,
  },
  antiPumpItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: 4,
    borderRadius: theme.borderRadius.sm,
  },
  antiPumpText: {
    fontSize: 11,
    color: COLORS.gray700,
    marginLeft: 4,
  },
  vetoWarning: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.error + '10',
    padding: theme.spacing.sm,
    borderRadius: theme.borderRadius.sm,
    marginBottom: theme.spacing.sm,
  },
  vetoText: {
    fontSize: 12,
    color: COLORS.error,
    marginLeft: theme.spacing.sm,
    flex: 1,
  },
  festivalBonus: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.festivalPrimary + '10',
    padding: theme.spacing.sm,
    borderRadius: theme.borderRadius.sm,
  },
  festivalBonusText: {
    fontSize: 12,
    color: COLORS.festivalPrimary,
    marginLeft: theme.spacing.sm,
    fontWeight: '600',
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xxl,
  },
  emptyText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray700,
    marginTop: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  emptyButton: {
    paddingHorizontal: theme.spacing.xl,
  },
});