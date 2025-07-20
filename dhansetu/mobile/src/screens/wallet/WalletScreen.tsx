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
  RefreshControl,
  FlatList,
  Image,
  Alert,
  Clipboard,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import QRCode from 'react-native-qrcode-svg';
import { useNavigation } from '@react-navigation/native';
import Animated, {
  useAnimatedStyle,
  withSpring,
  withTiming,
  interpolate,
} from 'react-native-reanimated';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector, useAppDispatch } from '@store/index';
import { setCurrentCoin } from '@store/slices/walletSlice';
import { useFestival } from '@contexts/FestivalContext';

interface TokenBalance {
  symbol: string;
  name: string;
  balance: string;
  value: string;
  change24h: number;
  icon?: string;
  color: string;
}

interface Transaction {
  id: string;
  type: 'send' | 'receive' | 'stake' | 'unstake' | 'reward';
  amount: string;
  symbol: string;
  address: string;
  timestamp: string;
  status: 'success' | 'pending' | 'failed';
  memo?: string;
  culturalQuote?: string;
}

export const WalletScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  const { currentAddress, currentCoin, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival } = useFestival();
  
  const [refreshing, setRefreshing] = useState(false);
  const [showQR, setShowQR] = useState(false);
  const [selectedToken, setSelectedToken] = useState<TokenBalance | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [activeTab, setActiveTab] = useState<'tokens' | 'nfts' | 'activity'>('tokens');

  // Mock token balances
  const [tokenBalances, setTokenBalances] = useState<TokenBalance[]>([
    {
      symbol: 'NAMO',
      name: 'NAMO Token',
      balance: '12,345.67',
      value: '₹1,23,456.70',
      change24h: 5.67,
      color: COLORS.saffron,
    },
    {
      symbol: 'ETH',
      name: 'Ethereum',
      balance: '0.5432',
      value: '₹1,08,640',
      change24h: -2.34,
      color: '#627EEA',
    },
    {
      symbol: 'BTC',
      name: 'Bitcoin',
      balance: '0.0123',
      value: '₹49,200',
      change24h: 3.21,
      color: '#F7931A',
    },
  ]);

  useEffect(() => {
    fetchTransactions();
  }, [currentCoin]);

  const onRefresh = async () => {
    setRefreshing(true);
    await fetchTransactions();
    await fetchBalances();
    setRefreshing(false);
  };

  const fetchTransactions = async () => {
    // Simulate API call
    setTimeout(() => {
      setTransactions(getMockTransactions());
    }, 1000);
  };

  const fetchBalances = async () => {
    // Fetch real balances from blockchain
    // For now using mock data
  };

  const copyAddress = (address: string) => {
    Clipboard.setString(address);
    Alert.alert('Copied!', 'Address copied to clipboard');
  };

  const renderAddressCard = () => {
    const displayAddress = dhanPataAddress || currentAddress || '';
    const shortAddress = displayAddress.includes('@dhan') 
      ? displayAddress 
      : `${displayAddress.slice(0, 10)}...${displayAddress.slice(-8)}`;

    return (
      <LinearGradient
        colors={currentFestival ? theme.gradients.festivalGradient : theme.gradients.indianFlag}
        style={styles.addressCard}
      >
        <View style={styles.addressHeader}>
          <Text style={styles.addressLabel}>Your {currentCoin} Address</Text>
          <TouchableOpacity
            style={styles.qrButton}
            onPress={() => setShowQR(!showQR)}
          >
            <Icon name={showQR ? 'close' : 'qrcode'} size={24} color={COLORS.white} />
          </TouchableOpacity>
        </View>

        {showQR ? (
          <View style={styles.qrContainer}>
            <View style={styles.qrCode}>
              <QRCode
                value={currentAddress || ''}
                size={150}
                backgroundColor="white"
                color={COLORS.navy}
              />
            </View>
            <Text style={styles.scanText}>Scan to receive {currentCoin}</Text>
          </View>
        ) : (
          <>
            <TouchableOpacity
              style={styles.addressContainer}
              onPress={() => copyAddress(displayAddress)}
              activeOpacity={0.8}
            >
              <Text style={styles.address}>{shortAddress}</Text>
              <Icon name="content-copy" size={20} color={COLORS.white} />
            </TouchableOpacity>

            {dhanPataAddress && (
              <View style={styles.dhanPataContainer}>
                <Icon name="at" size={16} color={COLORS.white} />
                <Text style={styles.dhanPataText}>
                  DhanPata ID: {dhanPataAddress}
                </Text>
              </View>
            )}

            {/* Total Balance */}
            <View style={styles.balanceSection}>
              <Text style={styles.totalLabel}>Total Balance</Text>
              <GradientText style={styles.totalBalance}>
                ₹2,81,296.70
              </GradientText>
              {currentFestival && (
                <View style={styles.festivalIndicator}>
                  <Icon name="party-popper" size={16} color={COLORS.white} />
                  <Text style={styles.festivalText}>
                    {currentFestival.bonusRate}% {currentFestival.name} bonus active
                  </Text>
                </View>
              )}
            </View>
          </>
        )}

        {/* Action Buttons */}
        <View style={styles.actionButtons}>
          <TouchableOpacity
            style={styles.actionButton}
            onPress={() => navigation.navigate('Send', { coinType: currentCoin })}
          >
            <Icon name="send" size={24} color={COLORS.white} />
            <Text style={styles.actionText}>Send</Text>
          </TouchableOpacity>
          
          <TouchableOpacity
            style={styles.actionButton}
            onPress={() => navigation.navigate('Receive', { coinType: currentCoin })}
          >
            <Icon name="qrcode-scan" size={24} color={COLORS.white} />
            <Text style={styles.actionText}>Receive</Text>
          </TouchableOpacity>
          
          <TouchableOpacity
            style={styles.actionButton}
            onPress={() => {/* Navigate to swap */}}
          >
            <Icon name="swap-horizontal" size={24} color={COLORS.white} />
            <Text style={styles.actionText}>Swap</Text>
          </TouchableOpacity>
          
          <TouchableOpacity
            style={styles.actionButton}
            onPress={() => {/* Navigate to stake */}}
          >
            <Icon name="lock" size={24} color={COLORS.white} />
            <Text style={styles.actionText}>Stake</Text>
          </TouchableOpacity>
        </View>
      </LinearGradient>
    );
  };

  const renderTokenItem = ({ item }: { item: TokenBalance }) => {
    const isPositive = item.change24h >= 0;
    
    return (
      <TouchableOpacity
        style={styles.tokenItem}
        onPress={() => {
          setSelectedToken(item);
          dispatch(setCurrentCoin(item.symbol as any));
        }}
        activeOpacity={0.7}
      >
        <View style={styles.tokenLeft}>
          <View style={[styles.tokenIcon, { backgroundColor: item.color }]}>
            <Text style={styles.tokenSymbolIcon}>
              {item.symbol.charAt(0)}
            </Text>
          </View>
          <View style={styles.tokenInfo}>
            <Text style={styles.tokenName}>{item.name}</Text>
            <Text style={styles.tokenSymbol}>{item.symbol}</Text>
          </View>
        </View>
        
        <View style={styles.tokenRight}>
          <Text style={styles.tokenBalance}>{item.balance}</Text>
          <View style={styles.tokenValueRow}>
            <Text style={styles.tokenValue}>{item.value}</Text>
            <Text style={[
              styles.tokenChange,
              { color: isPositive ? COLORS.success : COLORS.error }
            ]}>
              {isPositive ? '+' : ''}{item.change24h}%
            </Text>
          </View>
        </View>
      </TouchableOpacity>
    );
  };

  const renderTransactionItem = ({ item }: { item: Transaction }) => {
    const getIcon = () => {
      switch (item.type) {
        case 'send': return 'arrow-up-circle';
        case 'receive': return 'arrow-down-circle';
        case 'stake': return 'lock';
        case 'unstake': return 'lock-open';
        case 'reward': return 'gift';
        default: return 'circle';
      }
    };

    const getColor = () => {
      switch (item.type) {
        case 'send': return COLORS.error;
        case 'receive': return COLORS.success;
        case 'stake': return COLORS.saffron;
        case 'unstake': return COLORS.warning;
        case 'reward': return COLORS.festivalPrimary;
        default: return COLORS.gray600;
      }
    };

    return (
      <TouchableOpacity
        style={styles.transactionItem}
        onPress={() => navigation.navigate('TransactionDetails', { txHash: item.id })}
      >
        <View style={styles.transactionLeft}>
          <Icon name={getIcon()} size={32} color={getColor()} />
          <View style={styles.transactionInfo}>
            <Text style={styles.transactionType}>
              {item.type.charAt(0).toUpperCase() + item.type.slice(1)}
            </Text>
            <Text style={styles.transactionAddress} numberOfLines={1}>
              {item.type === 'send' ? 'To: ' : 'From: '}{item.address}
            </Text>
          </View>
        </View>
        
        <View style={styles.transactionRight}>
          <Text style={[
            styles.transactionAmount,
            { color: item.type === 'send' ? COLORS.error : COLORS.success }
          ]}>
            {item.type === 'send' ? '-' : '+'}{item.amount} {item.symbol}
          </Text>
          <Text style={styles.transactionTime}>
            {new Date(item.timestamp).toLocaleDateString()}
          </Text>
        </View>
      </TouchableOpacity>
    );
  };

  const renderTokensTab = () => (
    <FlatList
      data={tokenBalances}
      renderItem={renderTokenItem}
      keyExtractor={(item) => item.symbol}
      contentContainerStyle={styles.tokenList}
      ListEmptyComponent={
        <View style={styles.emptyState}>
          <Icon name="coin" size={48} color={COLORS.gray400} />
          <Text style={styles.emptyText}>No tokens yet</Text>
        </View>
      }
    />
  );

  const renderActivityTab = () => (
    <FlatList
      data={transactions}
      renderItem={renderTransactionItem}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.transactionList}
      ListEmptyComponent={
        <View style={styles.emptyState}>
          <Icon name="history" size={48} color={COLORS.gray400} />
          <Text style={styles.emptyText}>No transactions yet</Text>
        </View>
      }
    />
  );

  const renderNFTsTab = () => (
    <View style={styles.emptyState}>
      <Icon name="image-multiple" size={48} color={COLORS.gray400} />
      <Text style={styles.emptyText}>NFTs coming soon!</Text>
      <Text style={styles.emptySubtext}>
        Collect cultural NFTs and festival rewards
      </Text>
    </View>
  );

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
        {/* Chain Selector */}
        <View style={styles.chainSelector}>
          {(['DESHCHAIN', 'ETHEREUM', 'BITCOIN'] as const).map((chain) => (
            <TouchableOpacity
              key={chain}
              style={[
                styles.chainButton,
                currentCoin === chain && styles.chainButtonActive,
              ]}
              onPress={() => dispatch(setCurrentCoin(chain))}
            >
              <Text style={[
                styles.chainText,
                currentCoin === chain && styles.chainTextActive,
              ]}>
                {chain}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Address Card */}
        {renderAddressCard()}

        {/* Tabs */}
        <View style={styles.tabs}>
          {(['tokens', 'nfts', 'activity'] as const).map((tab) => (
            <TouchableOpacity
              key={tab}
              style={[styles.tab, activeTab === tab && styles.activeTab]}
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
        <View style={styles.tabContent}>
          {activeTab === 'tokens' && renderTokensTab()}
          {activeTab === 'nfts' && renderNFTsTab()}
          {activeTab === 'activity' && renderActivityTab()}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

// Mock data
const getMockTransactions = (): Transaction[] => [
  {
    id: '1',
    type: 'receive',
    amount: '1000',
    symbol: 'NAMO',
    address: 'ramesh@dhan',
    timestamp: new Date().toISOString(),
    status: 'success',
    culturalQuote: 'Unity in diversity',
  },
  {
    id: '2',
    type: 'send',
    amount: '500',
    symbol: 'NAMO',
    address: 'suresh@dhan',
    timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    status: 'success',
    memo: 'Thanks for lunch!',
  },
  {
    id: '3',
    type: 'stake',
    amount: '5000',
    symbol: 'NAMO',
    address: 'Validator: DeshChain Foundation',
    timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    status: 'success',
  },
  {
    id: '4',
    type: 'reward',
    amount: '50',
    symbol: 'NAMO',
    address: 'Staking Rewards',
    timestamp: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
    status: 'success',
  },
];

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  chainSelector: {
    flexDirection: 'row',
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    gap: theme.spacing.sm,
  },
  chainButton: {
    flex: 1,
    paddingVertical: theme.spacing.sm,
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
  },
  chainButtonActive: {
    backgroundColor: COLORS.saffron,
  },
  chainText: {
    fontSize: 12,
    fontWeight: '600',
    color: COLORS.gray700,
  },
  chainTextActive: {
    color: COLORS.white,
  },
  addressCard: {
    marginHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    ...theme.shadows.medium,
  },
  addressHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  addressLabel: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  qrButton: {
    padding: theme.spacing.sm,
  },
  addressContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.md,
  },
  address: {
    flex: 1,
    fontSize: 16,
    color: COLORS.white,
    fontWeight: '500',
  },
  dhanPataContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  dhanPataText: {
    fontSize: 14,
    color: COLORS.white,
    marginLeft: theme.spacing.xs,
  },
  qrContainer: {
    alignItems: 'center',
    paddingVertical: theme.spacing.lg,
  },
  qrCode: {
    padding: theme.spacing.md,
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.md,
  },
  scanText: {
    fontSize: 14,
    color: COLORS.white,
    marginTop: theme.spacing.md,
  },
  balanceSection: {
    alignItems: 'center',
    marginBottom: theme.spacing.lg,
  },
  totalLabel: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  totalBalance: {
    fontSize: 36,
    fontWeight: 'bold',
    marginVertical: theme.spacing.sm,
  },
  festivalIndicator: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.full,
  },
  festivalText: {
    fontSize: 12,
    color: COLORS.white,
    marginLeft: theme.spacing.xs,
  },
  actionButtons: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  actionButton: {
    alignItems: 'center',
    padding: theme.spacing.sm,
  },
  actionText: {
    fontSize: 12,
    color: COLORS.white,
    marginTop: 4,
  },
  tabs: {
    flexDirection: 'row',
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  tab: {
    flex: 1,
    paddingVertical: theme.spacing.md,
    alignItems: 'center',
  },
  activeTab: {
    borderBottomWidth: 2,
    borderBottomColor: COLORS.saffron,
  },
  tabText: {
    fontSize: 14,
    color: COLORS.gray600,
    fontWeight: '500',
  },
  activeTabText: {
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  tabContent: {
    flex: 1,
    minHeight: 300,
  },
  tokenList: {
    padding: theme.spacing.lg,
  },
  tokenItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.sm,
  },
  tokenLeft: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  tokenIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: theme.spacing.md,
  },
  tokenSymbolIcon: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  tokenInfo: {},
  tokenName: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.gray900,
  },
  tokenSymbol: {
    fontSize: 12,
    color: COLORS.gray600,
  },
  tokenRight: {
    alignItems: 'flex-end',
  },
  tokenBalance: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.gray900,
  },
  tokenValueRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: theme.spacing.sm,
  },
  tokenValue: {
    fontSize: 12,
    color: COLORS.gray600,
  },
  tokenChange: {
    fontSize: 12,
    fontWeight: '600',
  },
  transactionList: {
    padding: theme.spacing.lg,
  },
  transactionItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  transactionLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  transactionInfo: {
    marginLeft: theme.spacing.md,
    flex: 1,
  },
  transactionType: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.gray900,
  },
  transactionAddress: {
    fontSize: 12,
    color: COLORS.gray600,
    marginTop: 2,
  },
  transactionRight: {
    alignItems: 'flex-end',
  },
  transactionAmount: {
    fontSize: 16,
    fontWeight: '600',
  },
  transactionTime: {
    fontSize: 12,
    color: COLORS.gray500,
    marginTop: 2,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xxl,
  },
  emptyText: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.gray700,
    marginTop: theme.spacing.md,
  },
  emptySubtext: {
    fontSize: 14,
    color: COLORS.gray500,
    marginTop: theme.spacing.sm,
  },
});