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
  TextInput,
  Alert,
  Dimensions,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import Animated, {
  useAnimatedStyle,
  withSpring,
  withTiming,
  interpolate,
  useSharedValue,
} from 'react-native-reanimated';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector, useAppDispatch } from '@store/index';
import { useFestival } from '@contexts/FestivalContext';

const { width } = Dimensions.get('window');

interface CollateralAsset {
  symbol: string;
  name: string;
  balance: string;
  price: number;
  ratio: number;
  color: string;
}

interface StabilityMetric {
  label: string;
  value: string;
  status: 'good' | 'warning' | 'danger';
  icon: string;
}

export const DINRScreen: React.FC = () => {
  const navigation = useNavigation();
  const { currentFestival } = useFestival();
  
  const [activeTab, setActiveTab] = useState<'mint' | 'burn' | 'stats'>('mint');
  const [amount, setAmount] = useState('');
  const [selectedCollateral, setSelectedCollateral] = useState<CollateralAsset | null>(null);
  const [collateralAmount, setCollateralAmount] = useState('');
  
  const animatedTabValue = useSharedValue(0);

  const collateralAssets: CollateralAsset[] = [
    {
      symbol: 'ETH',
      name: 'Ethereum',
      balance: '0.5432',
      price: 200000, // ₹2,00,000 per ETH
      ratio: 150, // 150% collateral ratio
      color: '#627EEA',
    },
    {
      symbol: 'BTC',
      name: 'Bitcoin',
      balance: '0.0123',
      price: 4000000, // ₹40,00,000 per BTC
      ratio: 150,
      color: '#F7931A',
    },
    {
      symbol: 'USDT',
      name: 'Tether',
      balance: '5000.00',
      price: 83, // ₹83 per USDT
      ratio: 110, // Lower ratio for stable assets
      color: '#26A17B',
    },
    {
      symbol: 'USDC',
      name: 'USD Coin',
      balance: '3200.00',
      price: 83, // ₹83 per USDC
      ratio: 110,
      color: '#2775CA',
    },
  ];

  const stabilityMetrics: StabilityMetric[] = [
    {
      label: 'DINR Price',
      value: '₹1.00',
      status: 'good',
      icon: 'currency-inr',
    },
    {
      label: 'Total Supply',
      value: '₹12.5 Cr',
      status: 'good',
      icon: 'chart-line',
    },
    {
      label: 'Collateral Ratio',
      value: '156%',
      status: 'good',
      icon: 'shield-check',
    },
    {
      label: 'Stability Fee',
      value: '0.1%',
      status: 'good',
      icon: 'percent',
    },
    {
      label: 'Yield APY',
      value: '5.2%',
      status: 'good',
      icon: 'trending-up',
    },
    {
      label: 'Oracle Status',
      value: 'Active',
      status: 'good',
      icon: 'database',
    },
  ];

  useEffect(() => {
    animatedTabValue.value = withSpring(activeTab === 'mint' ? 0 : activeTab === 'burn' ? 1 : 2);
  }, [activeTab, animatedTabValue]);

  const calculateMintAmount = () => {
    if (!selectedCollateral || !collateralAmount) return 0;
    const collateralValue = parseFloat(collateralAmount) * selectedCollateral.price;
    return (collateralValue / selectedCollateral.ratio) * 100;
  };

  const handleMint = () => {
    if (!selectedCollateral || !collateralAmount || !amount) {
      Alert.alert('Error', 'Please fill all required fields');
      return;
    }

    const requiredCollateral = (parseFloat(amount) * selectedCollateral.ratio) / 100 / selectedCollateral.price;
    
    if (parseFloat(collateralAmount) < requiredCollateral) {
      Alert.alert('Insufficient Collateral', 
        `You need at least ${requiredCollateral.toFixed(6)} ${selectedCollateral.symbol} to mint ${amount} DINR`);
      return;
    }

    Alert.alert(
      'Confirm Mint',
      `Mint ${amount} DINR using ${collateralAmount} ${selectedCollateral.symbol} as collateral?`,
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Mint', 
          onPress: () => {
            // Execute mint transaction
            Alert.alert('Success', 'DINR minted successfully!');
          }
        }
      ]
    );
  };

  const handleBurn = () => {
    if (!amount) {
      Alert.alert('Error', 'Please enter amount to burn');
      return;
    }

    Alert.alert(
      'Confirm Burn',
      `Burn ${amount} DINR to retrieve collateral?`,
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Burn', 
          onPress: () => {
            // Execute burn transaction
            Alert.alert('Success', 'DINR burned and collateral released!');
          }
        }
      ]
    );
  };

  const renderTabIndicator = () => {
    const animatedStyle = useAnimatedStyle(() => {
      return {
        transform: [
          {
            translateX: interpolate(
              animatedTabValue.value,
              [0, 1, 2],
              [0, width / 3, (width * 2) / 3]
            ),
          },
        ],
      };
    });

    return (
      <Animated.View style={[styles.tabIndicator, animatedStyle]} />
    );
  };

  const renderMintTab = () => (
    <ScrollView style={styles.tabContent}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Select Collateral</Text>
        <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.collateralScroll}>
          {collateralAssets.map((asset) => (
            <TouchableOpacity
              key={asset.symbol}
              style={[
                styles.collateralCard,
                selectedCollateral?.symbol === asset.symbol && styles.collateralCardSelected,
              ]}
              onPress={() => setSelectedCollateral(asset)}
            >
              <View style={[styles.collateralIcon, { backgroundColor: asset.color }]}>
                <Text style={styles.collateralSymbol}>{asset.symbol.charAt(0)}</Text>
              </View>
              <Text style={styles.collateralName}>{asset.name}</Text>
              <Text style={styles.collateralBalance}>{asset.balance}</Text>
              <Text style={styles.collateralRatio}>{asset.ratio}% ratio</Text>
            </TouchableOpacity>
          ))}
        </ScrollView>
      </View>

      {selectedCollateral && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Collateral Amount</Text>
          <View style={styles.inputContainer}>
            <TextInput
              style={styles.input}
              placeholder="0.00"
              value={collateralAmount}
              onChangeText={setCollateralAmount}
              keyboardType="numeric"
            />
            <Text style={styles.inputUnit}>{selectedCollateral.symbol}</Text>
          </View>
          <Text style={styles.inputHelper}>
            Available: {selectedCollateral.balance} {selectedCollateral.symbol}
          </Text>
        </View>
      )}

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>DINR Amount</Text>
        <View style={styles.inputContainer}>
          <TextInput
            style={styles.input}
            placeholder="0.00"
            value={amount}
            onChangeText={setAmount}
            keyboardType="numeric"
          />
          <Text style={styles.inputUnit}>DINR</Text>
        </View>
        {selectedCollateral && collateralAmount && (
          <Text style={styles.inputHelper}>
            Max mintable: {calculateMintAmount().toFixed(2)} DINR
          </Text>
        )}
      </View>

      <View style={styles.section}>
        <View style={styles.feeBreakdown}>
          <Text style={styles.feeTitle}>Transaction Summary</Text>
          <View style={styles.feeRow}>
            <Text style={styles.feeLabel}>Stability Fee (0.1%)</Text>
            <Text style={styles.feeValue}>₹{(parseFloat(amount) * 0.001).toFixed(2)}</Text>
          </View>
          <View style={styles.feeRow}>
            <Text style={styles.feeLabel}>Gas Fee (estimated)</Text>
            <Text style={styles.feeValue}>₹15</Text>
          </View>
          <View style={[styles.feeRow, styles.feeTotalRow]}>
            <Text style={styles.feeTotalLabel}>Total Cost</Text>
            <Text style={styles.feeTotalValue}>₹{(parseFloat(amount) * 0.001 + 15).toFixed(2)}</Text>
          </View>
        </View>
      </View>

      <CulturalButton
        title="Mint DINR"
        onPress={handleMint}
        style={styles.actionButton}
        disabled={!selectedCollateral || !collateralAmount || !amount}
      />
    </ScrollView>
  );

  const renderBurnTab = () => (
    <ScrollView style={styles.tabContent}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>DINR Balance</Text>
        <View style={styles.balanceCard}>
          <LinearGradient
            colors={[COLORS.navy, COLORS.saffron]}
            style={styles.balanceGradient}
          >
            <Text style={styles.balanceLabel}>Available DINR</Text>
            <Text style={styles.balanceAmount}>50,000.00</Text>
            <Text style={styles.balanceValue}>₹50,000.00</Text>
          </LinearGradient>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Burn Amount</Text>
        <View style={styles.inputContainer}>
          <TextInput
            style={styles.input}
            placeholder="0.00"
            value={amount}
            onChangeText={setAmount}
            keyboardType="numeric"
          />
          <Text style={styles.inputUnit}>DINR</Text>
        </View>
        <Text style={styles.inputHelper}>
          Available: 50,000.00 DINR
        </Text>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Collateral Recovery</Text>
        <Text style={styles.recoveryText}>
          Burning {amount || '0'} DINR will release approximately:
        </Text>
        <View style={styles.recoveryList}>
          {collateralAssets.map((asset) => (
            <View key={asset.symbol} style={styles.recoveryItem}>
              <View style={[styles.recoveryIcon, { backgroundColor: asset.color }]}>
                <Text style={styles.recoverySymbol}>{asset.symbol.charAt(0)}</Text>
              </View>
              <Text style={styles.recoveryAmount}>
                {(parseFloat(amount || '0') / asset.price).toFixed(6)} {asset.symbol}
              </Text>
            </View>
          ))}
        </View>
      </View>

      <CulturalButton
        title="Burn DINR"
        onPress={handleBurn}
        style={[styles.actionButton, styles.burnButton]}
        disabled={!amount}
      />
    </ScrollView>
  );

  const renderStatsTab = () => (
    <ScrollView style={styles.tabContent}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>DINR Stability Metrics</Text>
        <View style={styles.metricsGrid}>
          {stabilityMetrics.map((metric, index) => (
            <View key={index} style={styles.metricCard}>
              <View style={[styles.metricIcon, 
                { backgroundColor: metric.status === 'good' ? COLORS.success : 
                  metric.status === 'warning' ? COLORS.warning : COLORS.error }]}>
                <Icon name={metric.icon} size={20} color={COLORS.white} />
              </View>
              <Text style={styles.metricLabel}>{metric.label}</Text>
              <Text style={styles.metricValue}>{metric.value}</Text>
            </View>
          ))}
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Your DINR Position</Text>
        <View style={styles.positionCard}>
          <LinearGradient
            colors={currentFestival ? theme.gradients.festivalGradient : [COLORS.navy, COLORS.saffron]}
            style={styles.positionGradient}
          >
            <View style={styles.positionHeader}>
              <Text style={styles.positionTitle}>Total DINR Holdings</Text>
              <Text style={styles.positionAmount}>50,000.00 DINR</Text>
            </View>
            <View style={styles.positionStats}>
              <View style={styles.positionStat}>
                <Text style={styles.positionStatLabel}>Yield Earned</Text>
                <Text style={styles.positionStatValue}>2,125.50 DINR</Text>
              </View>
              <View style={styles.positionStat}>
                <Text style={styles.positionStatLabel}>APY</Text>
                <Text style={styles.positionStatValue}>5.2%</Text>
              </View>
            </View>
          </LinearGradient>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>About DINR</Text>
        <View style={styles.infoCard}>
          <Text style={styles.infoText}>
            DINR is DeshChain's algorithmic stablecoin pegged 1:1 to Indian Rupee (₹1.00).
            It's backed by multiple crypto assets and maintains stability through automated
            mechanisms and yield-generating strategies.
          </Text>
          
          <View style={styles.featureList}>
            <View style={styles.featureItem}>
              <Icon name="shield-check" size={16} color={COLORS.success} />
              <Text style={styles.featureText}>Multi-collateral backing for security</Text>
            </View>
            <View style={styles.featureItem}>
              <Icon name="trending-up" size={16} color={COLORS.success} />
              <Text style={styles.featureText}>5%+ yield on DINR holdings</Text>
            </View>
            <View style={styles.featureItem}>
              <Icon name="currency-inr" size={16} color={COLORS.success} />
              <Text style={styles.featureText}>Lowest fees: 0.1% capped at ₹100</Text>
            </View>
            <View style={styles.featureItem}>
              <Icon name="flash" size={16} color={COLORS.success} />
              <Text style={styles.featureText}>Instant transactions on DeshChain</Text>
            </View>
          </View>
        </View>
      </View>
    </ScrollView>
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity
          style={styles.backButton}
          onPress={() => navigation.goBack()}
        >
          <Icon name="arrow-left" size={24} color={COLORS.navy} />
        </TouchableOpacity>
        <Text style={styles.title}>DINR Stablecoin</Text>
        <View style={styles.headerRight} />
      </View>

      {/* Tabs */}
      <View style={styles.tabs}>
        {renderTabIndicator()}
        {(['mint', 'burn', 'stats'] as const).map((tab) => (
          <TouchableOpacity
            key={tab}
            style={styles.tab}
            onPress={() => setActiveTab(tab)}
          >
            <Text style={[
              styles.tabText,
              activeTab === tab && styles.activeTabText,
            ]}>
              {tab.charAt(0).toUpperCase() + tab.slice(1)}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {/* Tab Content */}
      {activeTab === 'mint' && renderMintTab()}
      {activeTab === 'burn' && renderBurnTab()}
      {activeTab === 'stats' && renderStatsTab()}
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
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  backButton: {
    padding: theme.spacing.sm,
  },
  title: {
    flex: 1,
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.navy,
    textAlign: 'center',
  },
  headerRight: {
    width: 40,
  },
  tabs: {
    flexDirection: 'row',
    position: 'relative',
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  tab: {
    flex: 1,
    paddingVertical: theme.spacing.md,
    alignItems: 'center',
  },
  tabText: {
    fontSize: 16,
    color: COLORS.gray600,
    fontWeight: '500',
  },
  activeTabText: {
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  tabIndicator: {
    position: 'absolute',
    bottom: 0,
    height: 2,
    width: width / 3,
    backgroundColor: COLORS.saffron,
  },
  tabContent: {
    flex: 1,
    padding: theme.spacing.lg,
  },
  section: {
    marginBottom: theme.spacing.xl,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.navy,
    marginBottom: theme.spacing.md,
  },
  collateralScroll: {
    marginBottom: theme.spacing.md,
  },
  collateralCard: {
    width: 120,
    padding: theme.spacing.md,
    marginRight: theme.spacing.md,
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    borderWidth: 1,
    borderColor: COLORS.gray200,
    alignItems: 'center',
  },
  collateralCardSelected: {
    borderColor: COLORS.saffron,
    backgroundColor: COLORS.saffron + '10',
  },
  collateralIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  collateralSymbol: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  collateralName: {
    fontSize: 12,
    color: COLORS.gray700,
    textAlign: 'center',
  },
  collateralBalance: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.navy,
    marginTop: 2,
  },
  collateralRatio: {
    fontSize: 10,
    color: COLORS.gray500,
    marginTop: 2,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    borderWidth: 1,
    borderColor: COLORS.gray200,
    paddingHorizontal: theme.spacing.md,
  },
  input: {
    flex: 1,
    fontSize: 18,
    color: COLORS.navy,
    paddingVertical: theme.spacing.md,
  },
  inputUnit: {
    fontSize: 16,
    color: COLORS.gray600,
    marginLeft: theme.spacing.sm,
  },
  inputHelper: {
    fontSize: 12,
    color: COLORS.gray500,
    marginTop: theme.spacing.sm,
  },
  feeBreakdown: {
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
  },
  feeTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.navy,
    marginBottom: theme.spacing.md,
  },
  feeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.sm,
  },
  feeLabel: {
    fontSize: 14,
    color: COLORS.gray700,
  },
  feeValue: {
    fontSize: 14,
    fontWeight: '500',
    color: COLORS.navy,
  },
  feeTotalRow: {
    borderTopWidth: 1,
    borderTopColor: COLORS.gray300,
    paddingTop: theme.spacing.sm,
    marginTop: theme.spacing.sm,
  },
  feeTotalLabel: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.navy,
  },
  feeTotalValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.saffron,
  },
  actionButton: {
    marginTop: theme.spacing.lg,
  },
  burnButton: {
    backgroundColor: COLORS.error,
  },
  balanceCard: {
    borderRadius: theme.borderRadius.md,
    overflow: 'hidden',
    ...theme.shadows.medium,
  },
  balanceGradient: {
    padding: theme.spacing.lg,
    alignItems: 'center',
  },
  balanceLabel: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  balanceAmount: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.white,
    marginVertical: theme.spacing.sm,
  },
  balanceValue: {
    fontSize: 16,
    color: COLORS.white,
    opacity: 0.9,
  },
  recoveryText: {
    fontSize: 14,
    color: COLORS.gray700,
    marginBottom: theme.spacing.md,
  },
  recoveryList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: theme.spacing.md,
  },
  recoveryItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    minWidth: '45%',
  },
  recoveryIcon: {
    width: 32,
    height: 32,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: theme.spacing.sm,
  },
  recoverySymbol: {
    fontSize: 14,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  recoveryAmount: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.navy,
  },
  metricsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: theme.spacing.md,
  },
  metricCard: {
    width: '48%',
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
  },
  metricIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  metricLabel: {
    fontSize: 12,
    color: COLORS.gray600,
    textAlign: 'center',
    marginBottom: theme.spacing.xs,
  },
  metricValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.navy,
    textAlign: 'center',
  },
  positionCard: {
    borderRadius: theme.borderRadius.md,
    overflow: 'hidden',
    ...theme.shadows.medium,
  },
  positionGradient: {
    padding: theme.spacing.lg,
  },
  positionHeader: {
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  positionTitle: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  positionAmount: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
    marginTop: theme.spacing.sm,
  },
  positionStats: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  positionStat: {
    alignItems: 'center',
  },
  positionStatLabel: {
    fontSize: 12,
    color: COLORS.white,
    opacity: 0.8,
  },
  positionStatValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.white,
    marginTop: theme.spacing.xs,
  },
  infoCard: {
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.md,
  },
  infoText: {
    fontSize: 14,
    color: COLORS.gray700,
    lineHeight: 20,
    marginBottom: theme.spacing.md,
  },
  featureList: {
    gap: theme.spacing.sm,
  },
  featureItem: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  featureText: {
    fontSize: 14,
    color: COLORS.gray700,
    marginLeft: theme.spacing.sm,
    flex: 1,
  },
});