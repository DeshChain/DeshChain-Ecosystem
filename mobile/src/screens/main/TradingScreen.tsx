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
} from 'react-native-paper';
import { LineChart, CandlestickChart } from 'react-native-chart-kit';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';
import * as Animatable from 'react-native-animatable';

import { RootState } from '../../store';
import { theme, spacing, colors } from '../../theme';
import { TradingService } from '../../services/TradingService';
import { PriceService } from '../../services/PriceService';

const { width } = Dimensions.get('window');

interface TradingPair {
  symbol: string;
  baseAsset: string;
  quoteAsset: string;
  price: number;
  change24h: number;
  volume24h: number;
  high24h: number;
  low24h: number;
}

interface OrderBookEntry {
  price: number;
  amount: number;
  total: number;
}

interface RecentTrade {
  id: string;
  price: number;
  amount: number;
  timestamp: Date;
  side: 'buy' | 'sell';
}

export const TradingScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance } = useSelector((state: RootState) => state.wallet);
  const { language } = useSelector((state: RootState) => state.settings);
  
  const [refreshing, setRefreshing] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedPair, setSelectedPair] = useState<TradingPair | null>(null);
  const [tradingPairs, setTradingPairs] = useState<TradingPair[]>([]);
  const [orderBook, setOrderBook] = useState<{
    bids: OrderBookEntry[];
    asks: OrderBookEntry[];
  }>({ bids: [], asks: [] });
  const [recentTrades, setRecentTrades] = useState<RecentTrade[]>([]);
  const [chartData, setChartData] = useState<number[]>([]);
  const [viewMode, setViewMode] = useState('list');
  const [chartPeriod, setChartPeriod] = useState('1h');

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      // Fetch trading pairs
      const pairs = await TradingService.getTradingPairs();
      setTradingPairs(pairs);
      
      if (!selectedPair && pairs.length > 0) {
        setSelectedPair(pairs[0]);
      }
      
      if (selectedPair) {
        // Fetch order book
        const orderBookData = await TradingService.getOrderBook(selectedPair.symbol);
        setOrderBook(orderBookData);
        
        // Fetch recent trades
        const trades = await TradingService.getRecentTrades(selectedPair.symbol);
        setRecentTrades(trades);
        
        // Fetch chart data
        const chart = await PriceService.getChartData(selectedPair.symbol, chartPeriod);
        setChartData(chart);
      }
      
    } catch (error) {
      console.error('Error refreshing trading data:', error);
      Alert.alert('Error', 'Failed to refresh trading data. Please try again.');
    } finally {
      setRefreshing(false);
    }
  }, [selectedPair, chartPeriod]);

  useEffect(() => {
    onRefresh();
  }, [onRefresh]);

  const filteredPairs = tradingPairs.filter(pair =>
    pair.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
    pair.baseAsset.toLowerCase().includes(searchQuery.toLowerCase()) ||
    pair.quoteAsset.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const formatPrice = (price: number) => {
    return price.toLocaleString('en-IN', { 
      minimumFractionDigits: 2, 
      maximumFractionDigits: 6 
    });
  };

  const formatVolume = (volume: number) => {
    if (volume >= 1000000) {
      return `${(volume / 1000000).toFixed(1)}M`;
    } else if (volume >= 1000) {
      return `${(volume / 1000).toFixed(1)}K`;
    }
    return volume.toFixed(2);
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
      r: "4",
      strokeWidth: "2",
      stroke: colors.primary,
    },
  };

  const renderPairList = () => (
    <Card style={styles.pairListCard}>
      <Card.Content style={styles.pairListContent}>
        <Searchbar
          placeholder={language === 'hi' ? 'जोड़ी खोजें...' : 'Search pairs...'}
          onChangeText={setSearchQuery}
          value={searchQuery}
          style={styles.searchBar}
        />
        
        <ScrollView style={styles.pairList} showsVerticalScrollIndicator={false}>
          {filteredPairs.map((pair, index) => (
            <View key={pair.symbol}>
              <List.Item
                title={pair.symbol}
                description={`${pair.baseAsset}/${pair.quoteAsset}`}
                right={() => (
                  <View style={styles.pairInfo}>
                    <Text style={styles.pairPrice}>₹{formatPrice(pair.price)}</Text>
                    <Chip
                      textStyle={{ fontSize: 10 }}
                      style={[
                        styles.changeChip,
                        {
                          backgroundColor: pair.change24h >= 0 ? colors.success : colors.error
                        }
                      ]}
                    >
                      {pair.change24h >= 0 ? '+' : ''}{pair.change24h.toFixed(2)}%
                    </Chip>
                    <Text style={styles.pairVolume}>Vol: {formatVolume(pair.volume24h)}</Text>
                  </View>
                )}
                onPress={() => setSelectedPair(pair)}
                style={selectedPair?.symbol === pair.symbol ? styles.selectedPair : undefined}
              />
              {index < filteredPairs.length - 1 && <Divider />}
            </View>
          ))}
        </ScrollView>
      </Card.Content>
    </Card>
  );

  const renderTradingInterface = () => (
    <View>
      {/* Selected Pair Header */}
      {selectedPair && (
        <Animatable.View animation="fadeInUp" duration={800}>
          <Card style={styles.selectedPairCard}>
            <Card.Content>
              <View style={styles.selectedPairHeader}>
                <View>
                  <Text style={styles.selectedPairSymbol}>{selectedPair.symbol}</Text>
                  <Text style={styles.selectedPairName}>
                    {selectedPair.baseAsset}/{selectedPair.quoteAsset}
                  </Text>
                </View>
                <View style={styles.selectedPairPriceInfo}>
                  <Text style={styles.selectedPairPrice}>₹{formatPrice(selectedPair.price)}</Text>
                  <Chip
                    textStyle={{ fontSize: 12, color: 'white' }}
                    style={[
                      styles.selectedPairChange,
                      {
                        backgroundColor: selectedPair.change24h >= 0 ? colors.success : colors.error
                      }
                    ]}
                  >
                    {selectedPair.change24h >= 0 ? '+' : ''}{selectedPair.change24h.toFixed(2)}%
                  </Chip>
                </View>
              </View>
              
              <View style={styles.pairStats}>
                <View style={styles.statItem}>
                  <Text style={styles.statLabel}>24h High</Text>
                  <Text style={styles.statValue}>₹{formatPrice(selectedPair.high24h)}</Text>
                </View>
                <View style={styles.statItem}>
                  <Text style={styles.statLabel}>24h Low</Text>
                  <Text style={styles.statValue}>₹{formatPrice(selectedPair.low24h)}</Text>
                </View>
                <View style={styles.statItem}>
                  <Text style={styles.statLabel}>24h Volume</Text>
                  <Text style={styles.statValue}>{formatVolume(selectedPair.volume24h)}</Text>
                </View>
              </View>
            </Card.Content>
          </Card>
        </Animatable.View>
      )}

      {/* Chart */}
      {chartData.length > 0 && (
        <Animatable.View animation="fadeInUp" duration={800} delay={200}>
          <Card style={styles.chartCard}>
            <Card.Content>
              <View style={styles.chartHeader}>
                <Text style={styles.chartTitle}>
                  {language === 'hi' ? 'मूल्य चार्ट' : 'Price Chart'}
                </Text>
                <SegmentedButtons
                  value={chartPeriod}
                  onValueChange={setChartPeriod}
                  buttons={[
                    { value: '1h', label: '1H' },
                    { value: '4h', label: '4H' },
                    { value: '1d', label: '1D' },
                    { value: '1w', label: '1W' },
                  ]}
                  style={styles.chartPeriodButtons}
                />
              </View>
              
              <LineChart
                data={{
                  labels: ['', '', '', '', '', '', ''],
                  datasets: [{ data: chartData }],
                }}
                width={width - 80}
                height={200}
                chartConfig={chartConfig}
                bezier
                style={styles.chart}
              />
            </Card.Content>
          </Card>
        </Animatable.View>
      )}

      {/* Trading Actions */}
      <Animatable.View animation="fadeInUp" duration={800} delay={400}>
        <View style={styles.tradingActions}>
          <Surface style={[styles.tradingActionButton, { backgroundColor: colors.success }]}>
            <Button
              mode="text"
              onPress={() => navigation.navigate('CreateOrder' as never, { orderType: 'buy' } as never)}
              contentStyle={styles.tradingActionContent}
            >
              <View style={styles.tradingActionInner}>
                <Icon name="trending-up" size={24} color="white" />
                <Text style={styles.tradingActionText}>
                  {language === 'hi' ? 'खरीदें' : 'BUY'}
                </Text>
              </View>
            </Button>
          </Surface>
          
          <Surface style={[styles.tradingActionButton, { backgroundColor: colors.error }]}>
            <Button
              mode="text"
              onPress={() => navigation.navigate('CreateOrder' as never, { orderType: 'sell' } as never)}
              contentStyle={styles.tradingActionContent}
            >
              <View style={styles.tradingActionInner}>
                <Icon name="trending-down" size={24} color="white" />
                <Text style={styles.tradingActionText}>
                  {language === 'hi' ? 'बेचें' : 'SELL'}
                </Text>
              </View>
            </Button>
          </Surface>
        </View>
      </Animatable.View>

      {/* Order Book & Recent Trades */}
      <Animatable.View animation="fadeInUp" duration={800} delay={600}>
        <View style={styles.marketDataContainer}>
          {/* Order Book */}
          <Card style={styles.orderBookCard}>
            <Card.Content>
              <Text style={styles.orderBookTitle}>
                {language === 'hi' ? 'ऑर्डर बुक' : 'Order Book'}
              </Text>
              
              {/* Asks */}
              <View style={styles.orderBookSection}>
                <Text style={styles.orderBookSectionTitle}>
                  {language === 'hi' ? 'बिक्री ऑर्डर' : 'SELL ORDERS'}
                </Text>
                {orderBook.asks.slice(0, 5).map((ask, index) => (
                  <View key={index} style={styles.orderBookEntry}>
                    <Text style={[styles.orderBookPrice, { color: colors.error }]}>
                      {formatPrice(ask.price)}
                    </Text>
                    <Text style={styles.orderBookAmount}>{ask.amount.toFixed(4)}</Text>
                  </View>
                ))}
              </View>
              
              <Divider style={styles.orderBookDivider} />
              
              {/* Bids */}
              <View style={styles.orderBookSection}>
                <Text style={styles.orderBookSectionTitle}>
                  {language === 'hi' ? 'खरीद ऑर्डर' : 'BUY ORDERS'}
                </Text>
                {orderBook.bids.slice(0, 5).map((bid, index) => (
                  <View key={index} style={styles.orderBookEntry}>
                    <Text style={[styles.orderBookPrice, { color: colors.success }]}>
                      {formatPrice(bid.price)}
                    </Text>
                    <Text style={styles.orderBookAmount}>{bid.amount.toFixed(4)}</Text>
                  </View>
                ))}
              </View>
            </Card.Content>
          </Card>

          {/* Recent Trades */}
          <Card style={styles.recentTradesCard}>
            <Card.Content>
              <Text style={styles.recentTradesTitle}>
                {language === 'hi' ? 'हाल के ट्रेड' : 'Recent Trades'}
              </Text>
              
              {recentTrades.slice(0, 10).map((trade) => (
                <View key={trade.id} style={styles.recentTradeEntry}>
                  <Text style={[
                    styles.recentTradePrice,
                    { color: trade.side === 'buy' ? colors.success : colors.error }
                  ]}>
                    {formatPrice(trade.price)}
                  </Text>
                  <Text style={styles.recentTradeAmount}>{trade.amount.toFixed(4)}</Text>
                  <Text style={styles.recentTradeTime}>
                    {trade.timestamp.toLocaleTimeString()}
                  </Text>
                </View>
              ))}
            </Card.Content>
          </Card>
        </View>
      </Animatable.View>
    </View>
  );

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
              {language === 'hi' ? 'DEX ट्रेडिंग' : 'DEX Trading'}
            </Text>
            <View style={styles.headerActions}>
              <SegmentedButtons
                value={viewMode}
                onValueChange={setViewMode}
                buttons={[
                  { value: 'list', label: 'List', icon: 'list' },
                  { value: 'trade', label: 'Trade', icon: 'trending-up' },
                ]}
                style={styles.viewModeButtons}
              />
            </View>
          </View>
        </Animatable.View>

        {/* Content based on view mode */}
        {viewMode === 'list' ? renderPairList() : renderTradingInterface()}

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
  headerActions: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  viewModeButtons: {
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
  },
  // Pair List Styles
  pairListCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
    flex: 1,
  },
  pairListContent: {
    padding: 0,
  },
  searchBar: {
    margin: spacing.md,
    backgroundColor: colors.surface,
  },
  pairList: {
    maxHeight: 600,
  },
  pairInfo: {
    alignItems: 'flex-end',
  },
  pairPrice: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
  },
  changeChip: {
    height: 20,
    marginTop: 4,
    marginBottom: 4,
  },
  pairVolume: {
    fontSize: 10,
    color: colors.textSecondary,
  },
  selectedPair: {
    backgroundColor: colors.surface,
  },
  // Trading Interface Styles
  selectedPairCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  selectedPairHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  selectedPairSymbol: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.text,
  },
  selectedPairName: {
    fontSize: 14,
    color: colors.textSecondary,
  },
  selectedPairPriceInfo: {
    alignItems: 'flex-end',
  },
  selectedPairPrice: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.primary,
  },
  selectedPairChange: {
    marginTop: 4,
  },
  pairStats: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  statItem: {
    alignItems: 'center',
  },
  statLabel: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  statValue: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
    marginTop: 2,
  },
  // Chart Styles
  chartCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  chartHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  chartTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
  },
  chartPeriodButtons: {
    backgroundColor: colors.surface,
  },
  chart: {
    marginVertical: spacing.sm,
    borderRadius: 16,
  },
  // Trading Actions
  tradingActions: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingHorizontal: spacing.md,
    marginBottom: spacing.md,
  },
  tradingActionButton: {
    width: (width - spacing.md * 3) / 2,
    borderRadius: 12,
    elevation: 4,
  },
  tradingActionContent: {
    height: 60,
    justifyContent: 'center',
  },
  tradingActionInner: {
    alignItems: 'center',
  },
  tradingActionText: {
    color: 'white',
    fontSize: 14,
    fontWeight: 'bold',
    marginTop: spacing.xs,
  },
  // Market Data
  marketDataContainer: {
    flexDirection: 'row',
    paddingHorizontal: spacing.md,
    gap: spacing.md,
  },
  orderBookCard: {
    flex: 1,
    backgroundColor: colors.card,
  },
  orderBookTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
    marginBottom: spacing.md,
    textAlign: 'center',
  },
  orderBookSection: {
    marginVertical: spacing.sm,
  },
  orderBookSectionTitle: {
    fontSize: 12,
    fontWeight: 'bold',
    color: colors.textSecondary,
    marginBottom: spacing.sm,
  },
  orderBookEntry: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.xs,
  },
  orderBookPrice: {
    fontSize: 12,
    fontWeight: '600',
  },
  orderBookAmount: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  orderBookDivider: {
    marginVertical: spacing.sm,
  },
  recentTradesCard: {
    flex: 1,
    backgroundColor: colors.card,
  },
  recentTradesTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
    marginBottom: spacing.md,
    textAlign: 'center',
  },
  recentTradeEntry: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.xs,
  },
  recentTradePrice: {
    fontSize: 12,
    fontWeight: '600',
    flex: 1,
  },
  recentTradeAmount: {
    fontSize: 12,
    color: colors.textSecondary,
    flex: 1,
    textAlign: 'center',
  },
  recentTradeTime: {
    fontSize: 10,
    color: colors.textSecondary,
    flex: 1,
    textAlign: 'right',
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