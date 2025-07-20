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
  Chip,
  Surface,
  IconButton,
  FAB,
  List,
  Divider,
  Searchbar,
  SegmentedButtons,
  Avatar,
  Badge,
} from 'react-native-paper';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';
import * as Animatable from 'react-native-animatable';

import { RootState } from '../../store';
import { theme, spacing, colors } from '../../theme';
import { P2PService } from '../../services/P2PService';

const { width } = Dimensions.get('window');

interface P2POrder {
  id: string;
  type: 'buy' | 'sell';
  amount: number;
  price: number;
  currency: string;
  paymentMethods: string[];
  minLimit: number;
  maxLimit: number;
  trader: {
    id: string;
    name: string;
    trustScore: number;
    completedTrades: number;
    responseTime: string;
    isOnline: boolean;
  };
  location: {
    city: string;
    state: string;
    pincode: string;
  };
  createdAt: Date;
  expiresAt: Date;
}

interface ActiveTrade {
  id: string;
  orderId: string;
  amount: number;
  price: number;
  total: number;
  status: 'pending' | 'payment_sent' | 'payment_confirmed' | 'completed' | 'disputed';
  counterparty: {
    name: string;
    trustScore: number;
  };
  paymentMethod: string;
  timeRemaining: number;
  escrowAmount: number;
}

export const P2PScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance } = useSelector((state: RootState) => state.wallet);
  const { language } = useSelector((state: RootState) => state.settings);
  
  const [refreshing, setRefreshing] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<'all' | 'buy' | 'sell'>('all');
  const [sortBy, setSortBy] = useState<'price' | 'trust' | 'time'>('price');
  const [orders, setOrders] = useState<P2POrder[]>([]);
  const [activeTrades, setActiveTrades] = useState<ActiveTrade[]>([]);
  const [selectedTab, setSelectedTab] = useState<'orders' | 'trades'>('orders');

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      // Fetch P2P orders
      const ordersData = await P2PService.getP2POrders({
        type: filterType === 'all' ? undefined : filterType,
        search: searchQuery,
        sortBy,
      });
      setOrders(ordersData);
      
      // Fetch active trades
      const tradesData = await P2PService.getActiveTrades();
      setActiveTrades(tradesData);
      
    } catch (error) {
      console.error('Error refreshing P2P data:', error);
      Alert.alert('Error', 'Failed to refresh P2P data. Please try again.');
    } finally {
      setRefreshing(false);
    }
  }, [filterType, searchQuery, sortBy]);

  useEffect(() => {
    onRefresh();
  }, [onRefresh]);

  const getTrustScoreBadge = (trustScore: number) => {
    if (trustScore >= 90) return { label: 'Diamond', color: colors.diamond };
    if (trustScore >= 80) return { label: 'Platinum', color: colors.platinum };
    if (trustScore >= 70) return { label: 'Gold', color: colors.goldBadge };
    if (trustScore >= 60) return { label: 'Silver', color: colors.silver };
    if (trustScore >= 50) return { label: 'Bronze', color: colors.bronze };
    return { label: 'New', color: colors.error };
  };

  const formatCurrency = (amount: number, currency: string = 'INR') => {
    if (currency === 'INR') {
      return `₹${amount.toLocaleString('en-IN')}`;
    }
    return `${amount.toLocaleString('en-IN')} ${currency}`;
  };

  const getPaymentMethodIcon = (method: string) => {
    switch (method.toLowerCase()) {
      case 'upi': return 'qr-code';
      case 'bank': return 'account-balance';
      case 'paytm': return 'payment';
      case 'phonepe': return 'phone-android';
      case 'googlepay': return 'google-pay';
      default: return 'payment';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return colors.warning;
      case 'payment_sent': return colors.info;
      case 'payment_confirmed': return colors.success;
      case 'completed': return colors.success;
      case 'disputed': return colors.error;
      default: return colors.textSecondary;
    }
  };

  const formatTimeRemaining = (minutes: number) => {
    if (minutes >= 60) {
      const hours = Math.floor(minutes / 60);
      const mins = minutes % 60;
      return `${hours}h ${mins}m`;
    }
    return `${minutes}m`;
  };

  const renderOrderCard = (order: P2POrder) => {
    const trustBadge = getTrustScoreBadge(order.trader.trustScore);
    
    return (
      <Card key={order.id} style={styles.orderCard}>
        <Card.Content>
          <View style={styles.orderHeader}>
            <View style={styles.orderType}>
              <Chip
                icon={order.type === 'buy' ? 'trending-up' : 'trending-down'}
                textStyle={{ color: 'white', fontWeight: 'bold' }}
                style={[
                  styles.orderTypeChip,
                  { backgroundColor: order.type === 'buy' ? colors.success : colors.error }
                ]}
              >
                {order.type === 'buy' ? 'BUY' : 'SELL'} {order.currency}
              </Chip>
              <View style={styles.orderPriceInfo}>
                <Text style={styles.orderPrice}>{formatCurrency(order.price)}</Text>
                <Text style={styles.orderAmount}>
                  {order.amount.toLocaleString()} {order.currency}
                </Text>
              </View>
            </View>
            
            <View style={styles.traderInfo}>
              <View style={styles.traderHeader}>
                <Avatar.Text
                  size={32}
                  label={order.trader.name.charAt(0)}
                  style={{ backgroundColor: trustBadge.color }}
                />
                {order.trader.isOnline && (
                  <Badge
                    style={styles.onlineBadge}
                    size={12}
                  />
                )}
              </View>
              <Text style={styles.traderName}>{order.trader.name}</Text>
              <Chip
                textStyle={{ fontSize: 10, color: 'white' }}
                style={[styles.trustChip, { backgroundColor: trustBadge.color }]}
              >
                {trustBadge.label} {order.trader.trustScore}
              </Chip>
            </View>
          </View>

          <Divider style={styles.orderDivider} />

          <View style={styles.orderDetails}>
            <View style={styles.orderLimits}>
              <Text style={styles.orderLabel}>
                {language === 'hi' ? 'सीमा' : 'Limits'}
              </Text>
              <Text style={styles.orderLimitText}>
                {formatCurrency(order.minLimit)} - {formatCurrency(order.maxLimit)}
              </Text>
            </View>
            
            <View style={styles.paymentMethods}>
              <Text style={styles.orderLabel}>
                {language === 'hi' ? 'भुगतान विधि' : 'Payment'}
              </Text>
              <View style={styles.paymentMethodsList}>
                {order.paymentMethods.slice(0, 3).map((method, index) => (
                  <Chip
                    key={index}
                    icon={getPaymentMethodIcon(method)}
                    compact
                    style={styles.paymentMethodChip}
                  >
                    {method}
                  </Chip>
                ))}
                {order.paymentMethods.length > 3 && (
                  <Chip compact style={styles.paymentMethodChip}>
                    +{order.paymentMethods.length - 3}
                  </Chip>
                )}
              </View>
            </View>

            <View style={styles.orderMeta}>
              <View style={styles.locationInfo}>
                <Icon name="location-on" size={14} color={colors.textSecondary} />
                <Text style={styles.locationText}>
                  {order.location.city}, {order.location.state}
                </Text>
              </View>
              <View style={styles.traderStats}>
                <Text style={styles.traderStatsText}>
                  {order.trader.completedTrades} trades • {order.trader.responseTime}
                </Text>
              </View>
            </View>
          </View>

          <View style={styles.orderActions}>
            <Button
              mode="outlined"
              style={styles.orderActionButton}
              onPress={() => navigation.navigate('TradeDetails' as never, { 
                tradeId: order.id 
              } as never)}
            >
              {language === 'hi' ? 'विवरण' : 'Details'}
            </Button>
            <Button
              mode="contained"
              style={styles.orderActionButton}
              onPress={() => {
                // Navigate to trade creation
                Alert.alert(
                  'Start Trade',
                  `Start a ${order.type} trade for ${order.currency}?`,
                  [
                    { text: 'Cancel', style: 'cancel' },
                    { text: 'Start Trade', onPress: () => console.log('Start trade:', order.id) },
                  ]
                );
              }}
            >
              {language === 'hi' ? 'ट्रेड करें' : 'Trade'}
            </Button>
          </View>
        </Card.Content>
      </Card>
    );
  };

  const renderActiveTradeCard = (trade: ActiveTrade) => {
    return (
      <Card key={trade.id} style={styles.tradeCard}>
        <Card.Content>
          <View style={styles.tradeHeader}>
            <View>
              <Text style={styles.tradeTitle}>
                Trade #{trade.id.slice(-6)}
              </Text>
              <Text style={styles.tradeCounterparty}>
                with {trade.counterparty.name}
              </Text>
            </View>
            <Chip
              textStyle={{ fontSize: 10, color: 'white' }}
              style={[
                styles.tradeStatusChip,
                { backgroundColor: getStatusColor(trade.status) }
              ]}
            >
              {trade.status.replace('_', ' ').toUpperCase()}
            </Chip>
          </View>

          <View style={styles.tradeDetails}>
            <View style={styles.tradeAmountInfo}>
              <Text style={styles.tradeAmount}>
                {trade.amount.toLocaleString()} NAMO
              </Text>
              <Text style={styles.tradeTotal}>
                {formatCurrency(trade.total)}
              </Text>
              <Text style={styles.tradePrice}>
                @ {formatCurrency(trade.price)}/NAMO
              </Text>
            </View>

            <View style={styles.tradeMetaInfo}>
              <View style={styles.tradeMetaItem}>
                <Icon name="payment" size={16} color={colors.textSecondary} />
                <Text style={styles.tradeMetaText}>{trade.paymentMethod}</Text>
              </View>
              <View style={styles.tradeMetaItem}>
                <Icon name="security" size={16} color={colors.textSecondary} />
                <Text style={styles.tradeMetaText}>
                  Escrow: {formatCurrency(trade.escrowAmount)}
                </Text>
              </View>
              <View style={styles.tradeMetaItem}>
                <Icon name="timer" size={16} color={colors.warning} />
                <Text style={[styles.tradeMetaText, { color: colors.warning }]}>
                  {formatTimeRemaining(trade.timeRemaining)}
                </Text>
              </View>
            </View>
          </View>

          <View style={styles.tradeActions}>
            <Button
              mode="outlined"
              style={styles.tradeActionButton}
              onPress={() => navigation.navigate('ChatDetails' as never, { 
                conversationId: trade.id 
              } as never)}
            >
              <Icon name="chat" size={16} />
              {language === 'hi' ? 'चैट' : 'Chat'}
            </Button>
            <Button
              mode="contained"
              style={styles.tradeActionButton}
              onPress={() => navigation.navigate('TradeDetails' as never, { 
                tradeId: trade.id 
              } as never)}
            >
              {language === 'hi' ? 'विवरण' : 'Manage'}
            </Button>
          </View>
        </Card.Content>
      </Card>
    );
  };

  return (
    <View style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {/* Header */}
        <Animatable.View animation="fadeInDown" duration={800} style={styles.header}>
          <View style={styles.headerContent}>
            <Text style={styles.headerTitle}>
              {language === 'hi' ? 'P2P ट्रेडिंग' : 'P2P Trading'}
            </Text>
            <IconButton
              icon="filter-list"
              size={24}
              iconColor="white"
              onPress={() => {
                // Show filter options
                Alert.alert(
                  'Filter Options',
                  'Choose filter criteria',
                  [
                    { text: 'All Orders', onPress: () => setFilterType('all') },
                    { text: 'Buy Orders', onPress: () => setFilterType('buy') },
                    { text: 'Sell Orders', onPress: () => setFilterType('sell') },
                    { text: 'Cancel', style: 'cancel' },
                  ]
                );
              }}
            />
          </View>
        </Animatable.View>

        {/* Tab Navigator */}
        <Animatable.View animation="fadeInUp" duration={800} delay={200}>
          <SegmentedButtons
            value={selectedTab}
            onValueChange={(value) => setSelectedTab(value as 'orders' | 'trades')}
            buttons={[
              { 
                value: 'orders', 
                label: language === 'hi' ? 'ऑर्डर्स' : 'Orders',
                icon: 'list'
              },
              { 
                value: 'trades', 
                label: language === 'hi' ? 'ट्रेड्स' : 'My Trades',
                icon: 'swap-horiz'
              },
            ]}
            style={styles.tabButtons}
          />
        </Animatable.View>

        {selectedTab === 'orders' && (
          <>
            {/* Search and Filters */}
            <Animatable.View animation="fadeInUp" duration={800} delay={400}>
              <View style={styles.searchContainer}>
                <Searchbar
                  placeholder={language === 'hi' ? 'ट्रेडर खोजें...' : 'Search traders...'}
                  onChangeText={setSearchQuery}
                  value={searchQuery}
                  style={styles.searchBar}
                />
                
                <View style={styles.filterChips}>
                  <Chip
                    selected={filterType === 'all'}
                    onPress={() => setFilterType('all')}
                    style={styles.filterChip}
                  >
                    All
                  </Chip>
                  <Chip
                    selected={filterType === 'buy'}
                    onPress={() => setFilterType('buy')}
                    style={styles.filterChip}
                  >
                    Buy
                  </Chip>
                  <Chip
                    selected={filterType === 'sell'}
                    onPress={() => setFilterType('sell')}
                    style={styles.filterChip}
                  >
                    Sell
                  </Chip>
                </View>
              </View>
            </Animatable.View>

            {/* Orders List */}
            <Animatable.View animation="fadeInUp" duration={800} delay={600}>
              {orders.length > 0 ? (
                <View style={styles.ordersList}>
                  {orders.map(renderOrderCard)}
                </View>
              ) : (
                <Card style={styles.emptyCard}>
                  <Card.Content style={styles.emptyContent}>
                    <Icon name="swap-horiz" size={64} color={colors.disabled} />
                    <Text style={styles.emptyText}>
                      {language === 'hi' ? 'कोई ऑर्डर नहीं मिला' : 'No orders found'}
                    </Text>
                    <Text style={styles.emptySubtext}>
                      {language === 'hi' 
                        ? 'नए ऑर्डर जल्द ही दिखाई देंगे'
                        : 'New orders will appear here soon'}
                    </Text>
                  </Card.Content>
                </Card>
              )}
            </Animatable.View>
          </>
        )}

        {selectedTab === 'trades' && (
          <Animatable.View animation="fadeInUp" duration={800} delay={400}>
            {activeTrades.length > 0 ? (
              <View style={styles.tradesList}>
                {activeTrades.map(renderActiveTradeCard)}
              </View>
            ) : (
              <Card style={styles.emptyCard}>
                <Card.Content style={styles.emptyContent}>
                  <Icon name="trending-up" size={64} color={colors.disabled} />
                  <Text style={styles.emptyText}>
                    {language === 'hi' ? 'कोई सक्रिय ट्रेड नहीं' : 'No active trades'}
                  </Text>
                  <Text style={styles.emptySubtext}>
                    {language === 'hi' 
                      ? 'ट्रेड शुरू करने के लिए ऑर्डर्स देखें'
                      : 'Check orders to start trading'}
                  </Text>
                </Card.Content>
              </Card>
            )}
          </Animatable.View>
        )}

        <View style={styles.bottomSpacer} />
      </ScrollView>

      {/* Floating Action Button */}
      <FAB
        icon="add"
        style={styles.fab}
        onPress={() => navigation.navigate('CreateOrder' as never)}
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
  tabButtons: {
    margin: spacing.md,
    backgroundColor: colors.surface,
  },
  searchContainer: {
    paddingHorizontal: spacing.md,
  },
  searchBar: {
    backgroundColor: colors.surface,
    marginBottom: spacing.md,
  },
  filterChips: {
    flexDirection: 'row',
    gap: spacing.sm,
    marginBottom: spacing.md,
  },
  filterChip: {
    backgroundColor: colors.surface,
  },
  // Orders List
  ordersList: {
    paddingHorizontal: spacing.md,
    gap: spacing.md,
  },
  orderCard: {
    backgroundColor: colors.card,
    elevation: 2,
  },
  orderHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.md,
  },
  orderType: {
    flex: 1,
  },
  orderTypeChip: {
    alignSelf: 'flex-start',
    marginBottom: spacing.sm,
  },
  orderPriceInfo: {
    alignItems: 'flex-start',
  },
  orderPrice: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.primary,
  },
  orderAmount: {
    fontSize: 14,
    color: colors.textSecondary,
  },
  traderInfo: {
    alignItems: 'center',
  },
  traderHeader: {
    position: 'relative',
  },
  onlineBadge: {
    position: 'absolute',
    top: -2,
    right: -2,
    backgroundColor: colors.success,
  },
  traderName: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.text,
    marginTop: spacing.xs,
  },
  trustChip: {
    marginTop: spacing.xs,
    height: 20,
  },
  orderDivider: {
    marginVertical: spacing.md,
  },
  orderDetails: {
    marginBottom: spacing.md,
  },
  orderLimits: {
    marginBottom: spacing.sm,
  },
  orderLabel: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  orderLimitText: {
    fontSize: 14,
    color: colors.text,
  },
  paymentMethods: {
    marginBottom: spacing.sm,
  },
  paymentMethodsList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.xs,
  },
  paymentMethodChip: {
    backgroundColor: colors.surface,
    height: 24,
  },
  orderMeta: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  locationInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  locationText: {
    fontSize: 12,
    color: colors.textSecondary,
    marginLeft: spacing.xs,
  },
  traderStats: {},
  traderStatsText: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  orderActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: spacing.md,
  },
  orderActionButton: {
    flex: 1,
  },
  // Active Trades
  tradesList: {
    paddingHorizontal: spacing.md,
    gap: spacing.md,
  },
  tradeCard: {
    backgroundColor: colors.card,
    elevation: 2,
  },
  tradeHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  tradeTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
  },
  tradeCounterparty: {
    fontSize: 14,
    color: colors.textSecondary,
  },
  tradeStatusChip: {
    height: 24,
  },
  tradeDetails: {
    marginBottom: spacing.md,
  },
  tradeAmountInfo: {
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  tradeAmount: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.primary,
  },
  tradeTotal: {
    fontSize: 16,
    color: colors.text,
    marginTop: spacing.xs,
  },
  tradePrice: {
    fontSize: 12,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  tradeMetaInfo: {
    gap: spacing.sm,
  },
  tradeMetaItem: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  tradeMetaText: {
    fontSize: 12,
    color: colors.textSecondary,
    marginLeft: spacing.sm,
  },
  tradeActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: spacing.md,
  },
  tradeActionButton: {
    flex: 1,
  },
  // Empty States
  emptyCard: {
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