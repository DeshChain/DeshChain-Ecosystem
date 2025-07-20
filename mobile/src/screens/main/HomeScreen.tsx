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
  Portal,
  FAB,
} from 'react-native-paper';
import { LineChart, PieChart } from 'react-native-chart-kit';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';
import * as Animatable from 'react-native-animatable';

import { RootState } from '../../store';
import { theme, spacing, colors } from '../../theme';
import { WalletService } from '../../services/WalletService';
import { PriceService } from '../../services/PriceService';
import { CulturalQuoteService } from '../../services/CulturalQuoteService';
import { NotificationService } from '../../services/NotificationService';

const { width } = Dimensions.get('window');

interface QuickAction {
  id: string;
  title: string;
  titleHindi: string;
  icon: string;
  color: string;
  route: string;
  params?: any;
}

interface RecentTransaction {
  id: string;
  type: 'send' | 'receive' | 'trade' | 'seva_mitra';
  amount: string;
  currency: string;
  counterparty: string;
  timestamp: Date;
  status: 'completed' | 'pending' | 'failed';
}

interface MarketData {
  price: number;
  change24h: number;
  chartData: number[];
  volume24h: number;
}

export const HomeScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  
  const { user, isAuthenticated } = useSelector((state: RootState) => state.auth);
  const { balance, transactions } = useSelector((state: RootState) => state.wallet);
  const { language } = useSelector((state: RootState) => state.settings);
  
  const [refreshing, setRefreshing] = useState(false);
  const [marketData, setMarketData] = useState<MarketData | null>(null);
  const [culturalQuote, setCulturalQuote] = useState<string>('');
  const [fabOpen, setFabOpen] = useState(false);
  const [recentTransactions, setRecentTransactions] = useState<RecentTransaction[]>([]);

  // Quick actions for home screen
  const quickActions: QuickAction[] = [
    {
      id: 'send',
      title: 'Send Money',
      titleHindi: 'पैसे भेजें',
      icon: 'send',
      color: colors.primary,
      route: 'SendMoney',
    },
    {
      id: 'receive',
      title: 'Receive',
      titleHindi: 'प्राप्त करें',
      icon: 'qr-code',
      color: colors.secondary,
      route: 'ReceiveMoney',
    },
    {
      id: 'trade',
      title: 'P2P Trade',
      titleHindi: 'P2P ट्रेड',
      icon: 'swap-horiz',
      color: colors.info,
      route: 'CreateOrder',
    },
    {
      id: 'seva_mitra',
      title: 'Find Seva Mitra',
      titleHindi: 'सेवा मित्र खोजें',
      icon: 'location-on',
      color: colors.gold,
      route: 'SevaMitraMap',
    },
  ];

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      // Refresh wallet balance
      await WalletService.refreshBalance();
      
      // Refresh market data
      const data = await PriceService.getMarketData('NAMO');
      setMarketData(data);
      
      // Get new cultural quote
      const quote = await CulturalQuoteService.getDailyQuote(language);
      setCulturalQuote(quote);
      
      // Refresh transactions
      const txns = await WalletService.getRecentTransactions(5);
      setRecentTransactions(txns);
      
    } catch (error) {
      console.error('Error refreshing data:', error);
      Alert.alert('Error', 'Failed to refresh data. Please try again.');
    } finally {
      setRefreshing(false);
    }
  }, [language]);

  useEffect(() => {
    // Initial data load
    onRefresh();
    
    // Set up periodic refresh
    const interval = setInterval(() => {
      if (!refreshing) {
        onRefresh();
      }
    }, 30000); // Refresh every 30 seconds

    return () => clearInterval(interval);
  }, [onRefresh, refreshing]);

  const getGreeting = () => {
    const hour = new Date().getHours();
    const isHindi = language === 'hi';
    
    if (hour < 12) {
      return isHindi ? 'सुप्रभात' : 'Good Morning';
    } else if (hour < 17) {
      return isHindi ? 'नमस्ते' : 'Good Afternoon';
    } else {
      return isHindi ? 'शुभ संध्या' : 'Good Evening';
    }
  };

  const getTrustScoreBadge = (trustScore: number) => {
    if (trustScore >= 90) return { label: 'Diamond', color: colors.diamond };
    if (trustScore >= 80) return { label: 'Platinum', color: colors.platinum };
    if (trustScore >= 70) return { label: 'Gold', color: colors.goldBadge };
    if (trustScore >= 60) return { label: 'Silver', color: colors.silver };
    if (trustScore >= 50) return { label: 'Bronze', color: colors.bronze };
    return { label: 'New', color: colors.error };
  };

  const formatCurrency = (amount: string, currency: string = 'NAMO') => {
    const num = parseFloat(amount);
    if (currency === 'NAMO') {
      return `${num.toLocaleString('en-IN')} NAMO`;
    }
    return `₹${num.toLocaleString('en-IN')}`;
  };

  const getTransactionIcon = (type: string, status: string) => {
    const iconColor = status === 'completed' ? colors.success : 
                     status === 'pending' ? colors.warning : colors.error;
    
    switch (type) {
      case 'send': return <Icon name="arrow-upward" size={24} color={iconColor} />;
      case 'receive': return <Icon name="arrow-downward" size={24} color={iconColor} />;
      case 'trade': return <Icon name="swap-horiz" size={24} color={iconColor} />;
      case 'seva_mitra': return <Icon name="location-on" size={24} color={iconColor} />;
      default: return <Icon name="payment" size={24} color={iconColor} />;
    }
  };

  const chartConfig = {
    backgroundColor: colors.background,
    backgroundGradientFrom: colors.surface,
    backgroundGradientTo: colors.surface,
    decimalPlaces: 2,
    color: (opacity = 1) => `rgba(255, 153, 51, ${opacity})`,
    labelColor: (opacity = 1) => `rgba(117, 117, 117, ${opacity})`,
    style: {
      borderRadius: 16,
    },
    propsForDots: {
      r: "6",
      strokeWidth: "2",
      stroke: colors.primary,
    },
  };

  const trustBadge = getTrustScoreBadge(user?.trustScore || 50);

  return (
    <View style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {/* Header Section */}
        <Animatable.View animation="fadeInDown" duration={800} style={styles.header}>
          <View style={styles.headerContent}>
            <View style={styles.userInfo}>
              <Avatar.Text
                size={50}
                label={user?.displayName?.charAt(0) || 'U'}
                style={{ backgroundColor: colors.primary }}
              />
              <View style={styles.userDetails}>
                <Text style={styles.greeting}>{getGreeting()}</Text>
                <Text style={styles.userName}>{user?.displayName || 'DeshChain User'}</Text>
                <Chip
                  icon="star"
                  textStyle={{ color: 'white', fontSize: 12 }}
                  style={[styles.trustBadge, { backgroundColor: trustBadge.color }]}
                >
                  {trustBadge.label} {user?.trustScore || 50}
                </Chip>
              </View>
            </View>
            <IconButton
              icon="notifications"
              size={24}
              onPress={() => NotificationService.showLocalNotification('Test', 'Notification test')}
            />
          </View>
        </Animatable.View>

        {/* Cultural Quote */}
        {culturalQuote && (
          <Animatable.View animation="fadeInLeft" duration={800} delay={200}>
            <Card style={styles.quoteCard}>
              <Card.Content>
                <Text style={styles.quoteText}>"{culturalQuote}"</Text>
                <Text style={styles.quoteAuthor}>— DeshChain Daily Inspiration</Text>
              </Card.Content>
            </Card>
          </Animatable.View>
        )}

        {/* Balance Card */}
        <Animatable.View animation="fadeInUp" duration={800} delay={400}>
          <Card style={styles.balanceCard}>
            <Card.Content>
              <Text style={styles.balanceLabel}>
                {language === 'hi' ? 'कुल बैलेंस' : 'Total Balance'}
              </Text>
              <Text style={styles.balanceAmount}>
                {formatCurrency(balance?.total || '0')}
              </Text>
              <View style={styles.balanceDetails}>
                <View style={styles.balanceItem}>
                  <Text style={styles.balanceSubLabel}>Available</Text>
                  <Text style={styles.balanceSubAmount}>
                    {formatCurrency(balance?.available || '0')}
                  </Text>
                </View>
                <View style={styles.balanceItem}>
                  <Text style={styles.balanceSubLabel}>In Trading</Text>
                  <Text style={styles.balanceSubAmount}>
                    {formatCurrency(balance?.inTrading || '0')}
                  </Text>
                </View>
              </View>
            </Card.Content>
          </Card>
        </Animatable.View>

        {/* Quick Actions */}
        <Animatable.View animation="fadeInUp" duration={800} delay={600}>
          <Text style={styles.sectionTitle}>
            {language === 'hi' ? 'त्वरित कार्य' : 'Quick Actions'}
          </Text>
          <View style={styles.quickActions}>
            {quickActions.map((action, index) => (
              <Animatable.View
                key={action.id}
                animation="bounceIn"
                delay={800 + index * 100}
              >
                <Surface style={[styles.actionCard, { backgroundColor: action.color }]}>
                  <Button
                    mode="text"
                    onPress={() => navigation.navigate(action.route as never, action.params as never)}
                    contentStyle={styles.actionButton}
                    labelStyle={styles.actionLabel}
                  >
                    <View style={styles.actionContent}>
                      <Icon name={action.icon} size={32} color="white" />
                      <Text style={styles.actionText}>
                        {language === 'hi' ? action.titleHindi : action.title}
                      </Text>
                    </View>
                  </Button>
                </Surface>
              </Animatable.View>
            ))}
          </View>
        </Animatable.View>

        {/* Market Data */}
        {marketData && (
          <Animatable.View animation="fadeInUp" duration={800} delay={1000}>
            <Text style={styles.sectionTitle}>
              {language === 'hi' ? 'मार्केट डेटा' : 'Market Data'}
            </Text>
            <Card style={styles.marketCard}>
              <Card.Content>
                <View style={styles.marketHeader}>
                  <View>
                    <Text style={styles.marketPrice}>₹{marketData.price.toFixed(2)}</Text>
                    <Text style={[
                      styles.marketChange,
                      { color: marketData.change24h >= 0 ? colors.success : colors.error }
                    ]}>
                      {marketData.change24h >= 0 ? '+' : ''}
                      {marketData.change24h.toFixed(2)}%
                    </Text>
                  </View>
                  <View style={styles.marketVolume}>
                    <Text style={styles.marketVolumeLabel}>24h Volume</Text>
                    <Text style={styles.marketVolumeValue}>
                      ₹{(marketData.volume24h / 1000000).toFixed(1)}M
                    </Text>
                  </View>
                </View>
                {marketData.chartData.length > 0 && (
                  <LineChart
                    data={{
                      labels: ['6h', '5h', '4h', '3h', '2h', '1h', 'Now'],
                      datasets: [{ data: marketData.chartData }],
                    }}
                    width={width - 80}
                    height={150}
                    chartConfig={chartConfig}
                    bezier
                    style={styles.chart}
                  />
                )}
              </Card.Content>
            </Card>
          </Animatable.View>
        )}

        {/* Recent Transactions */}
        <Animatable.View animation="fadeInUp" duration={800} delay={1200}>
          <View style={styles.sectionHeader}>
            <Text style={styles.sectionTitle}>
              {language === 'hi' ? 'हाल की गतिविधि' : 'Recent Activity'}
            </Text>
            <Button
              mode="text"
              onPress={() => navigation.navigate('TransactionHistory' as never)}
            >
              {language === 'hi' ? 'सभी देखें' : 'View All'}
            </Button>
          </View>
          
          {recentTransactions.length > 0 ? (
            <Card style={styles.transactionsCard}>
              <Card.Content style={styles.transactionsContent}>
                {recentTransactions.map((tx, index) => (
                  <View key={tx.id} style={styles.transactionItem}>
                    <View style={styles.transactionIcon}>
                      {getTransactionIcon(tx.type, tx.status)}
                    </View>
                    <View style={styles.transactionDetails}>
                      <Text style={styles.transactionTitle}>
                        {tx.type === 'send' ? 'Sent to' : 
                         tx.type === 'receive' ? 'Received from' :
                         tx.type === 'trade' ? 'P2P Trade with' :
                         'Seva Mitra Service'}
                      </Text>
                      <Text style={styles.transactionCounterparty}>
                        {tx.counterparty}
                      </Text>
                      <Text style={styles.transactionTime}>
                        {tx.timestamp.toLocaleDateString()}
                      </Text>
                    </View>
                    <View style={styles.transactionAmount}>
                      <Text style={[
                        styles.transactionAmountText,
                        { color: tx.type === 'receive' ? colors.success : colors.text }
                      ]}>
                        {tx.type === 'receive' ? '+' : '-'}
                        {formatCurrency(tx.amount, tx.currency)}
                      </Text>
                      <Chip
                        textStyle={{ fontSize: 10 }}
                        style={[
                          styles.statusChip,
                          {
                            backgroundColor: 
                              tx.status === 'completed' ? colors.success :
                              tx.status === 'pending' ? colors.warning : colors.error
                          }
                        ]}
                      >
                        {tx.status}
                      </Chip>
                    </View>
                  </View>
                ))}
              </Card.Content>
            </Card>
          ) : (
            <Card style={styles.emptyTransactions}>
              <Card.Content style={styles.emptyContent}>
                <Icon name="receipt" size={48} color={colors.disabled} />
                <Text style={styles.emptyText}>
                  {language === 'hi' ? 'कोई हाल की गतिविधि नहीं' : 'No recent activity'}
                </Text>
                <Text style={styles.emptySubtext}>
                  {language === 'hi' 
                    ? 'आपके लेनदेन यहाँ दिखाई देंगे' 
                    : 'Your transactions will appear here'}
                </Text>
              </Card.Content>
            </Card>
          )}
        </Animatable.View>

        <View style={styles.bottomSpacer} />
      </ScrollView>

      {/* Floating Action Button */}
      <Portal>
        <FAB.Group
          open={fabOpen}
          icon={fabOpen ? 'close' : 'add'}
          actions={[
            {
              icon: 'qr-code-scanner',
              label: 'Scan QR',
              onPress: () => navigation.navigate('QRScanner' as never, {
                onScan: (data: string) => console.log('Scanned:', data)
              } as never),
            },
            {
              icon: 'send',
              label: 'Send Money',
              onPress: () => navigation.navigate('SendMoney' as never),
            },
            {
              icon: 'swap-horiz',
              label: 'P2P Trade',
              onPress: () => navigation.navigate('CreateOrder' as never),
            },
          ]}
          onStateChange={({ open }) => setFabOpen(open)}
          onPress={() => setFabOpen(!fabOpen)}
          fabStyle={{ backgroundColor: colors.primary }}
        />
      </Portal>
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
    borderBottomLeftRadius: 20,
    borderBottomRightRadius: 20,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  userInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  userDetails: {
    marginLeft: spacing.md,
  },
  greeting: {
    color: 'white',
    fontSize: 14,
    opacity: 0.9,
  },
  userName: {
    color: 'white',
    fontSize: 18,
    fontWeight: 'bold',
    marginTop: 2,
  },
  trustBadge: {
    marginTop: 4,
    height: 24,
  },
  quoteCard: {
    margin: spacing.md,
    backgroundColor: colors.cultural.lotus,
  },
  quoteText: {
    fontSize: 16,
    fontStyle: 'italic',
    textAlign: 'center',
    color: colors.text,
    lineHeight: 24,
  },
  quoteAuthor: {
    fontSize: 12,
    textAlign: 'center',
    marginTop: spacing.sm,
    color: colors.textSecondary,
  },
  balanceCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  balanceLabel: {
    fontSize: 14,
    color: colors.textSecondary,
    textAlign: 'center',
  },
  balanceAmount: {
    fontSize: 32,
    fontWeight: 'bold',
    color: colors.primary,
    textAlign: 'center',
    marginTop: spacing.sm,
  },
  balanceDetails: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginTop: spacing.md,
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  balanceItem: {
    alignItems: 'center',
  },
  balanceSubLabel: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  balanceSubAmount: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.text,
    marginTop: 2,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.text,
    marginHorizontal: spacing.md,
    marginTop: spacing.lg,
    marginBottom: spacing.md,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginHorizontal: spacing.md,
    marginTop: spacing.lg,
  },
  quickActions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.md,
  },
  actionCard: {
    width: (width - spacing.md * 3) / 2,
    marginBottom: spacing.md,
    borderRadius: 12,
    elevation: 4,
  },
  actionButton: {
    height: 80,
    justifyContent: 'center',
  },
  actionContent: {
    alignItems: 'center',
  },
  actionText: {
    color: 'white',
    fontSize: 12,
    fontWeight: '600',
    marginTop: spacing.xs,
    textAlign: 'center',
  },
  actionLabel: {
    margin: 0,
  },
  marketCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  marketHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  marketPrice: {
    fontSize: 24,
    fontWeight: 'bold',
    color: colors.text,
  },
  marketChange: {
    fontSize: 16,
    fontWeight: '600',
  },
  marketVolume: {
    alignItems: 'flex-end',
  },
  marketVolumeLabel: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  marketVolumeValue: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.text,
  },
  chart: {
    marginVertical: spacing.sm,
    borderRadius: 16,
  },
  transactionsCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  transactionsContent: {
    padding: 0,
  },
  transactionItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  transactionIcon: {
    marginRight: spacing.md,
  },
  transactionDetails: {
    flex: 1,
  },
  transactionTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
  },
  transactionCounterparty: {
    fontSize: 12,
    color: colors.textSecondary,
    marginTop: 2,
  },
  transactionTime: {
    fontSize: 10,
    color: colors.textSecondary,
    marginTop: 2,
  },
  transactionAmount: {
    alignItems: 'flex-end',
  },
  transactionAmountText: {
    fontSize: 14,
    fontWeight: '600',
  },
  statusChip: {
    height: 20,
    marginTop: 4,
  },
  emptyTransactions: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  emptyContent: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
  },
  emptyText: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.textSecondary,
    marginTop: spacing.md,
  },
  emptySubtext: {
    fontSize: 14,
    color: colors.textSecondary,
    marginTop: spacing.sm,
    textAlign: 'center',
  },
  bottomSpacer: {
    height: 100,
  },
});