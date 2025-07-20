import React, { useEffect, useState, useCallback } from 'react';
import {
  View,
  ScrollView,
  StyleSheet,
  RefreshControl,
  Dimensions,
  Alert,
} from 'react-native';
import {
  Text,
  Card,
  Button,
  Avatar,
  Chip,
  Surface,
  IconButton,
  FAB,
  List,
  Divider,
  ProgressBar,
} from 'react-native-paper';
import { PieChart, BarChart } from 'react-native-chart-kit';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';
import * as Animatable from 'react-native-animatable';

import { RootState } from '../../store';
import { theme, spacing, colors } from '../../theme';
import { WalletService } from '../../services/WalletService';
import { PriceService } from '../../services/PriceService';

const { width } = Dimensions.get('window');

interface Portfolio {
  symbol: string;
  name: string;
  balance: string;
  value: number;
  change24h: number;
  percentage: number;
}

interface StakingData {
  amount: string;
  rewards: string;
  apy: number;
  validator: string;
  unbondingDays: number;
}

export const WalletScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance, transactions } = useSelector((state: RootState) => state.wallet);
  const { language } = useSelector((state: RootState) => state.settings);
  
  const [refreshing, setRefreshing] = useState(false);
  const [portfolio, setPortfolio] = useState<Portfolio[]>([]);
  const [stakingData, setStakingData] = useState<StakingData | null>(null);
  const [totalValue, setTotalValue] = useState(0);
  const [showBalance, setShowBalance] = useState(true);

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      // Refresh wallet data
      await WalletService.refreshBalance();
      
      // Get portfolio breakdown
      const portfolioData = await WalletService.getPortfolioBreakdown();
      setPortfolio(portfolioData);
      
      // Get staking information
      const staking = await WalletService.getStakingData();
      setStakingData(staking);
      
      // Calculate total portfolio value
      const total = portfolioData.reduce((sum, asset) => sum + asset.value, 0);
      setTotalValue(total);
      
    } catch (error) {
      console.error('Error refreshing wallet data:', error);
      Alert.alert('Error', 'Failed to refresh wallet data. Please try again.');
    } finally {
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    onRefresh();
  }, [onRefresh]);

  const formatCurrency = (amount: string | number, currency: string = 'NAMO') => {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount;
    if (currency === 'NAMO') {
      return `${num.toLocaleString('en-IN', { maximumFractionDigits: 2 })} NAMO`;
    }
    return `₹${num.toLocaleString('en-IN', { maximumFractionDigits: 2 })}`;
  };

  const getAssetIcon = (symbol: string) => {
    switch (symbol) {
      case 'NAMO': return 'account-balance';
      case 'BTC': return 'currency-btc';
      case 'ETH': return 'currency-eth';
      default: return 'monetization-on';
    }
  };

  const chartConfig = {
    backgroundColor: colors.background,
    backgroundGradientFrom: colors.surface,
    backgroundGradientTo: colors.surface,
    decimalPlaces: 0,
    color: (opacity = 1) => `rgba(255, 153, 51, ${opacity})`,
    labelColor: (opacity = 1) => `rgba(117, 117, 117, ${opacity})`,
    style: {
      borderRadius: 16,
    },
  };

  const pieData = portfolio.map((asset, index) => ({
    name: asset.symbol,
    population: asset.percentage,
    color: [colors.primary, colors.secondary, colors.info, colors.warning, colors.success][index % 5],
    legendFontColor: colors.text,
    legendFontSize: 12,
  }));

  return (
    <View style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {/* Wallet Header */}
        <Animatable.View animation="fadeInDown" duration={800} style={styles.header}>
          <View style={styles.headerContent}>
            <Text style={styles.headerTitle}>
              {language === 'hi' ? 'डिजिटल वॉलेट' : 'Digital Wallet'}
            </Text>
            <IconButton
              icon={showBalance ? 'visibility' : 'visibility-off'}
              size={24}
              iconColor="white"
              onPress={() => setShowBalance(!showBalance)}
            />
          </View>
        </Animatable.View>

        {/* Total Balance Card */}
        <Animatable.View animation="fadeInUp" duration={800} delay={200}>
          <Card style={styles.balanceCard}>
            <Card.Content>
              <View style={styles.balanceHeader}>
                <Text style={styles.balanceLabel}>
                  {language === 'hi' ? 'कुल पोर्टफोलियो मूल्य' : 'Total Portfolio Value'}
                </Text>
                <Chip icon="trending-up" style={styles.performanceChip}>
                  +12.5%
                </Chip>
              </View>
              <Text style={styles.balanceAmount}>
                {showBalance ? formatCurrency(totalValue) : '••••••••'}
              </Text>
              <Text style={styles.balanceSubtext}>
                ≈ ₹{showBalance ? (totalValue * 85).toLocaleString('en-IN') : '••••••••'}
              </Text>
              
              {/* Balance Breakdown */}
              <View style={styles.balanceBreakdown}>
                <View style={styles.balanceItem}>
                  <Text style={styles.balanceItemLabel}>Available</Text>
                  <Text style={styles.balanceItemValue}>
                    {showBalance ? formatCurrency(balance?.available || '0') : '••••••'}
                  </Text>
                </View>
                <View style={styles.balanceItem}>
                  <Text style={styles.balanceItemLabel}>In Trading</Text>
                  <Text style={styles.balanceItemValue}>
                    {showBalance ? formatCurrency(balance?.inTrading || '0') : '••••••'}
                  </Text>
                </View>
                <View style={styles.balanceItem}>
                  <Text style={styles.balanceItemLabel}>Staked</Text>
                  <Text style={styles.balanceItemValue}>
                    {showBalance ? formatCurrency(stakingData?.amount || '0') : '••••••'}
                  </Text>
                </View>
              </View>
            </Card.Content>
          </Card>
        </Animatable.View>

        {/* Quick Actions */}
        <Animatable.View animation="fadeInUp" duration={800} delay={400}>
          <View style={styles.quickActions}>
            <Surface style={[styles.actionButton, { backgroundColor: colors.success }]}>
              <Button
                mode="text"
                onPress={() => navigation.navigate('ReceiveMoney' as never)}
                contentStyle={styles.actionButtonContent}
              >
                <View style={styles.actionButtonInner}>
                  <Icon name="qr-code" size={24} color="white" />
                  <Text style={styles.actionButtonText}>
                    {language === 'hi' ? 'प्राप्त करें' : 'Receive'}
                  </Text>
                </View>
              </Button>
            </Surface>
            
            <Surface style={[styles.actionButton, { backgroundColor: colors.primary }]}>
              <Button
                mode="text"
                onPress={() => navigation.navigate('SendMoney' as never)}
                contentStyle={styles.actionButtonContent}
              >
                <View style={styles.actionButtonInner}>
                  <Icon name="send" size={24} color="white" />
                  <Text style={styles.actionButtonText}>
                    {language === 'hi' ? 'भेजें' : 'Send'}
                  </Text>
                </View>
              </Button>
            </Surface>
            
            <Surface style={[styles.actionButton, { backgroundColor: colors.info }]}>
              <Button
                mode="text"
                onPress={() => navigation.navigate('CreateOrder' as never)}
                contentStyle={styles.actionButtonContent}
              >
                <View style={styles.actionButtonInner}>
                  <Icon name="swap-horiz" size={24} color="white" />
                  <Text style={styles.actionButtonText}>
                    {language === 'hi' ? 'ट्रेड' : 'Trade'}
                  </Text>
                </View>
              </Button>
            </Surface>
          </View>
        </Animatable.View>

        {/* Portfolio Allocation */}
        {portfolio.length > 0 && (
          <Animatable.View animation="fadeInUp" duration={800} delay={600}>
            <Text style={styles.sectionTitle}>
              {language === 'hi' ? 'पोर्टफोलियो आवंटन' : 'Portfolio Allocation'}
            </Text>
            <Card style={styles.portfolioCard}>
              <Card.Content>
                <PieChart
                  data={pieData}
                  width={width - 80}
                  height={200}
                  chartConfig={chartConfig}
                  accessor="population"
                  backgroundColor="transparent"
                  paddingLeft="15"
                  center={[0, 0]}
                />
              </Card.Content>
            </Card>
          </Animatable.View>
        )}

        {/* Assets List */}
        <Animatable.View animation="fadeInUp" duration={800} delay={800}>
          <Text style={styles.sectionTitle}>
            {language === 'hi' ? 'संपत्ति' : 'Assets'}
          </Text>
          <Card style={styles.assetsCard}>
            <Card.Content style={styles.assetsContent}>
              {portfolio.map((asset, index) => (
                <View key={asset.symbol}>
                  <List.Item
                    title={asset.name}
                    description={asset.symbol}
                    left={() => (
                      <Avatar.Icon
                        size={40}
                        icon={getAssetIcon(asset.symbol)}
                        style={{ backgroundColor: colors.primary }}
                      />
                    )}
                    right={() => (
                      <View style={styles.assetInfo}>
                        <Text style={styles.assetBalance}>
                          {showBalance ? formatCurrency(asset.balance, asset.symbol) : '••••••'}
                        </Text>
                        <Text style={styles.assetValue}>
                          ≈ ₹{showBalance ? asset.value.toLocaleString('en-IN') : '••••••'}
                        </Text>
                        <Chip
                          textStyle={{ fontSize: 10 }}
                          style={[
                            styles.changeChip,
                            {
                              backgroundColor: asset.change24h >= 0 ? colors.success : colors.error
                            }
                          ]}
                        >
                          {asset.change24h >= 0 ? '+' : ''}{asset.change24h.toFixed(2)}%
                        </Chip>
                      </View>
                    )}
                    onPress={() => {
                      // Navigate to asset details
                      console.log('Navigate to asset details:', asset.symbol);
                    }}
                  />
                  {index < portfolio.length - 1 && <Divider />}
                </View>
              ))}
            </Card.Content>
          </Card>
        </Animatable.View>

        {/* Staking Information */}
        {stakingData && (
          <Animatable.View animation="fadeInUp" duration={800} delay={1000}>
            <Text style={styles.sectionTitle}>
              {language === 'hi' ? 'स्टेकिंग रिवार्ड्स' : 'Staking Rewards'}
            </Text>
            <Card style={styles.stakingCard}>
              <Card.Content>
                <View style={styles.stakingHeader}>
                  <View>
                    <Text style={styles.stakingAmount}>
                      {showBalance ? formatCurrency(stakingData.amount) : '••••••••'}
                    </Text>
                    <Text style={styles.stakingLabel}>Staked Amount</Text>
                  </View>
                  <View style={styles.stakingApy}>
                    <Text style={styles.apyText}>{stakingData.apy}%</Text>
                    <Text style={styles.apyLabel}>APY</Text>
                  </View>
                </View>
                
                <Divider style={styles.stakingDivider} />
                
                <View style={styles.stakingDetails}>
                  <View style={styles.stakingDetailItem}>
                    <Text style={styles.stakingDetailLabel}>Pending Rewards</Text>
                    <Text style={styles.stakingDetailValue}>
                      {showBalance ? formatCurrency(stakingData.rewards) : '••••••'}
                    </Text>
                  </View>
                  <View style={styles.stakingDetailItem}>
                    <Text style={styles.stakingDetailLabel}>Validator</Text>
                    <Text style={styles.stakingDetailValue}>{stakingData.validator}</Text>
                  </View>
                  <View style={styles.stakingDetailItem}>
                    <Text style={styles.stakingDetailLabel}>Unbonding Period</Text>
                    <Text style={styles.stakingDetailValue}>{stakingData.unbondingDays} days</Text>
                  </View>
                </View>
                
                <View style={styles.stakingActions}>
                  <Button
                    mode="outlined"
                    style={styles.stakingActionButton}
                    onPress={() => console.log('Claim rewards')}
                  >
                    Claim Rewards
                  </Button>
                  <Button
                    mode="contained"
                    style={styles.stakingActionButton}
                    onPress={() => console.log('Manage staking')}
                  >
                    Manage
                  </Button>
                </View>
              </Card.Content>
            </Card>
          </Animatable.View>
        )}

        <View style={styles.bottomSpacer} />
      </ScrollView>

      {/* Floating Action Button */}
      <FAB
        icon="add"
        style={styles.fab}
        onPress={() => {
          Alert.alert(
            'Add Asset',
            'Choose how to add assets to your wallet',
            [
              { text: 'Buy with Fiat', onPress: () => console.log('Buy with fiat') },
              { text: 'Receive from Another Wallet', onPress: () => navigation.navigate('ReceiveMoney' as never) },
              { text: 'Cancel', style: 'cancel' },
            ]
          );
        }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    backgroundColor: colors.primary,
    paddingTop: 20,
    paddingBottom: 20,
    paddingHorizontal: spacing.md,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: 'white',
  },
  balanceCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  balanceHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  balanceLabel: {
    fontSize: 14,
    color: colors.textSecondary,
  },
  performanceChip: {
    backgroundColor: colors.success,
  },
  balanceAmount: {
    fontSize: 36,
    fontWeight: 'bold',
    color: colors.primary,
    marginBottom: spacing.xs,
  },
  balanceSubtext: {
    fontSize: 16,
    color: colors.textSecondary,
    marginBottom: spacing.md,
  },
  balanceBreakdown: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  balanceItem: {
    alignItems: 'center',
  },
  balanceItemLabel: {
    fontSize: 12,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  balanceItemValue: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
  },
  quickActions: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingHorizontal: spacing.md,
    marginBottom: spacing.md,
  },
  actionButton: {
    width: (width - spacing.md * 4) / 3,
    borderRadius: 12,
    elevation: 4,
  },
  actionButtonContent: {
    height: 60,
    justifyContent: 'center',
  },
  actionButtonInner: {
    alignItems: 'center',
  },
  actionButtonText: {
    color: 'white',
    fontSize: 12,
    fontWeight: '600',
    marginTop: spacing.xs,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.text,
    marginHorizontal: spacing.md,
    marginTop: spacing.lg,
    marginBottom: spacing.md,
  },
  portfolioCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  assetsCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  assetsContent: {
    padding: 0,
  },
  assetInfo: {
    alignItems: 'flex-end',
  },
  assetBalance: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
  },
  assetValue: {
    fontSize: 12,
    color: colors.textSecondary,
    marginTop: 2,
  },
  changeChip: {
    height: 20,
    marginTop: 4,
  },
  stakingCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  stakingHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  stakingAmount: {
    fontSize: 24,
    fontWeight: 'bold',
    color: colors.primary,
  },
  stakingLabel: {
    fontSize: 12,
    color: colors.textSecondary,
    marginTop: 2,
  },
  stakingApy: {
    alignItems: 'flex-end',
  },
  apyText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.success,
  },
  apyLabel: {
    fontSize: 12,
    color: colors.textSecondary,
    marginTop: 2,
  },
  stakingDivider: {
    marginVertical: spacing.md,
  },
  stakingDetails: {
    marginBottom: spacing.md,
  },
  stakingDetailItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.sm,
  },
  stakingDetailLabel: {
    fontSize: 14,
    color: colors.textSecondary,
  },
  stakingDetailValue: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
  },
  stakingActions: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  stakingActionButton: {
    flex: 0.45,
  },
  fab: {
    position: 'absolute',
    margin: 16,
    right: 0,
    bottom: 0,
    backgroundColor: colors.primary,
  },
  bottomSpacer: {
    height: 100,
  },
});