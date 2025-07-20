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
  FlatList,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector } from '@store/index';
import { useFestival } from '@contexts/FestivalContext';

interface MoneyOrder {
  id: string;
  orderId: string;
  creator: string;
  receiver: string;
  amount: string;
  fee: string;
  status: 'pending' | 'accepted' | 'completed' | 'cancelled' | 'expired';
  createdAt: string;
  expiresAt: string;
  paymentMode: 'namo' | 'fiat';
  memo?: string;
  culturalQuote?: string;
  pincode?: string;
}

export const DexScreen: React.FC = () => {
  const navigation = useNavigation();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival, festivalBonusRate } = useFestival();
  
  const [activeTab, setActiveTab] = useState<'create' | 'active' | 'history'>('create');
  const [orders, setOrders] = useState<MoneyOrder[]>([]);
  const [loading, setLoading] = useState(false);
  
  // Form state
  const [receiver, setReceiver] = useState('');
  const [amount, setAmount] = useState('');
  const [memo, setMemo] = useState('');
  const [paymentMode, setPaymentMode] = useState<'namo' | 'fiat'>('namo');
  const [pincode, setPincode] = useState('');

  useEffect(() => {
    fetchOrders();
  }, [activeTab]);

  const fetchOrders = async () => {
    setLoading(true);
    // Simulate API call
    setTimeout(() => {
      setOrders(getMockOrders());
      setLoading(false);
    }, 1000);
  };

  const calculateFee = () => {
    if (!amount) return '0';
    const baseAmount = parseFloat(amount);
    const baseFee = baseAmount * 0.01; // 1% base fee
    const discount = festivalBonusRate * baseFee; // Festival discount
    return (baseFee - discount).toFixed(2);
  };

  const createMoneyOrder = async () => {
    if (!receiver || !amount) {
      Alert.alert('Error', 'Please fill all required fields');
      return;
    }

    // Validate receiver (DhanPata or blockchain address)
    const isDhanPata = receiver.includes('@dhan');
    const isValidAddress = receiver.startsWith('desh1') && receiver.length === 44;

    if (!isDhanPata && !isValidAddress) {
      Alert.alert('Error', 'Invalid receiver address or DhanPata ID');
      return;
    }

    setLoading(true);
    
    try {
      // Create money order transaction
      Alert.alert(
        'Success!',
        `Money order created for ₹${amount} to ${receiver}`,
        [{ text: 'OK', onPress: () => resetForm() }]
      );
    } catch (error) {
      Alert.alert('Error', 'Failed to create money order');
    } finally {
      setLoading(false);
    }
  };

  const resetForm = () => {
    setReceiver('');
    setAmount('');
    setMemo('');
    setPincode('');
  };

  const renderCreateForm = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.formContainer}>
        <Text style={styles.formTitle}>Create Money Order</Text>
        
        {/* Receiver Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Receiver</Text>
          <View style={styles.inputWrapper}>
            <Icon name="account" size={20} color={COLORS.gray600} />
            <TextInput
              style={styles.input}
              placeholder="username@dhan or desh1..."
              value={receiver}
              onChangeText={setReceiver}
              autoCapitalize="none"
              autoCorrect={false}
            />
          </View>
          <Text style={styles.inputHint}>
            Enter DhanPata ID or blockchain address
          </Text>
        </View>

        {/* Amount Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Amount</Text>
          <View style={styles.inputWrapper}>
            <Text style={styles.currencySymbol}>₹</Text>
            <TextInput
              style={styles.input}
              placeholder="0.00"
              value={amount}
              onChangeText={setAmount}
              keyboardType="decimal-pad"
            />
          </View>
        </View>

        {/* Payment Mode */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Payment Mode</Text>
          <View style={styles.paymentModes}>
            <TouchableOpacity
              style={[
                styles.paymentMode,
                paymentMode === 'namo' && styles.paymentModeActive,
              ]}
              onPress={() => setPaymentMode('namo')}
            >
              <Icon 
                name="currency-inr" 
                size={20} 
                color={paymentMode === 'namo' ? COLORS.white : COLORS.gray600} 
              />
              <Text style={[
                styles.paymentModeText,
                paymentMode === 'namo' && styles.paymentModeTextActive,
              ]}>
                NAMO Token
              </Text>
            </TouchableOpacity>
            
            <TouchableOpacity
              style={[
                styles.paymentMode,
                paymentMode === 'fiat' && styles.paymentModeActive,
              ]}
              onPress={() => setPaymentMode('fiat')}
            >
              <Icon 
                name="cash" 
                size={20} 
                color={paymentMode === 'fiat' ? COLORS.white : COLORS.gray600} 
              />
              <Text style={[
                styles.paymentModeText,
                paymentMode === 'fiat' && styles.paymentModeTextActive,
              ]}>
                Cash/UPI
              </Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Pincode (Optional) */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>PIN Code (Optional)</Text>
          <View style={styles.inputWrapper}>
            <Icon name="map-marker" size={20} color={COLORS.gray600} />
            <TextInput
              style={styles.input}
              placeholder="6-digit PIN code"
              value={pincode}
              onChangeText={setPincode}
              keyboardType="number-pad"
              maxLength={6}
            />
          </View>
          <Text style={styles.inputHint}>
            For local Kshetra Coin rewards
          </Text>
        </View>

        {/* Memo */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Message (Optional)</Text>
          <TextInput
            style={styles.memoInput}
            placeholder="Add a message..."
            value={memo}
            onChangeText={setMemo}
            multiline
            numberOfLines={3}
          />
        </View>

        {/* Fee Display */}
        <View style={styles.feeContainer}>
          <View style={styles.feeRow}>
            <Text style={styles.feeLabel}>Service Fee</Text>
            <Text style={styles.feeAmount}>₹{calculateFee()}</Text>
          </View>
          {currentFestival && (
            <View style={styles.feeRow}>
              <Text style={styles.festivalLabel}>
                {currentFestival.name} Discount
              </Text>
              <Text style={styles.festivalDiscount}>
                -{currentFestival.bonusRate}%
              </Text>
            </View>
          )}
          <View style={[styles.feeRow, styles.totalRow]}>
            <Text style={styles.totalLabel}>Total Amount</Text>
            <Text style={styles.totalAmount}>
              ₹{(parseFloat(amount || '0') + parseFloat(calculateFee())).toFixed(2)}
            </Text>
          </View>
        </View>

        {/* Cultural Quote */}
        <View style={styles.quoteContainer}>
          <Icon name="format-quote-open" size={20} color={COLORS.saffron} />
          <Text style={styles.quoteText}>
            सत्यमेव जयते - Truth alone triumphs
          </Text>
        </View>

        {/* Create Button */}
        <CulturalButton
          title="Create Money Order"
          onPress={createMoneyOrder}
          loading={loading}
          size="large"
          style={styles.createButton}
        />
      </View>
    </ScrollView>
  );

  const renderOrderItem = ({ item }: { item: MoneyOrder }) => (
    <TouchableOpacity
      style={styles.orderCard}
      onPress={() => navigation.navigate('MoneyOrderDetails', { orderId: item.id })}
      activeOpacity={0.8}
    >
      <View style={styles.orderHeader}>
        <View style={styles.orderIdContainer}>
          <Text style={styles.orderId}>#{item.orderId}</Text>
          <View style={[styles.statusBadge, styles[`status${item.status}`]]}>
            <Text style={styles.statusText}>{item.status.toUpperCase()}</Text>
          </View>
        </View>
        <Text style={styles.orderAmount}>₹{item.amount}</Text>
      </View>
      
      <View style={styles.orderDetails}>
        <View style={styles.orderRow}>
          <Icon name="account-arrow-right" size={16} color={COLORS.gray600} />
          <Text style={styles.orderText} numberOfLines={1}>
            To: {item.receiver}
          </Text>
        </View>
        <View style={styles.orderRow}>
          <Icon name="clock-outline" size={16} color={COLORS.gray600} />
          <Text style={styles.orderText}>
            Created: {new Date(item.createdAt).toLocaleDateString()}
          </Text>
        </View>
      </View>
      
      {item.culturalQuote && (
        <Text style={styles.orderQuote}>"{item.culturalQuote}"</Text>
      )}
    </TouchableOpacity>
  );

  const renderActiveOrders = () => (
    <FlatList
      data={orders.filter(o => ['pending', 'accepted'].includes(o.status))}
      renderItem={renderOrderItem}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.ordersList}
      ListEmptyComponent={
        <View style={styles.emptyState}>
          <Icon name="inbox-outline" size={64} color={COLORS.gray400} />
          <Text style={styles.emptyText}>No active orders</Text>
          <Text style={styles.emptySubtext}>
            Create a money order to get started
          </Text>
        </View>
      }
    />
  );

  const renderHistory = () => (
    <FlatList
      data={orders.filter(o => ['completed', 'cancelled', 'expired'].includes(o.status))}
      renderItem={renderOrderItem}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.ordersList}
      ListEmptyComponent={
        <View style={styles.emptyState}>
          <Icon name="history" size={64} color={COLORS.gray400} />
          <Text style={styles.emptyText}>No order history</Text>
        </View>
      }
    />
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <LinearGradient
        colors={theme.gradients.indianFlag}
        style={styles.header}
      >
        <Text style={styles.headerTitle}>Money Order DEX</Text>
        <Text style={styles.headerSubtitle}>
          Digital money transfer for the modern age
        </Text>
      </LinearGradient>

      {/* Tabs */}
      <View style={styles.tabs}>
        {(['create', 'active', 'history'] as const).map((tab) => (
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

      {/* Content */}
      <View style={styles.content}>
        {activeTab === 'create' && renderCreateForm()}
        {activeTab === 'active' && renderActiveOrders()}
        {activeTab === 'history' && renderHistory()}
      </View>
    </SafeAreaView>
  );
};

// Mock data generator
const getMockOrders = (): MoneyOrder[] => [
  {
    id: '1',
    orderId: 'MO2024001',
    creator: 'desh1abc...xyz',
    receiver: 'ramesh@dhan',
    amount: '5000',
    fee: '50',
    status: 'pending',
    createdAt: new Date().toISOString(),
    expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
    paymentMode: 'namo',
    culturalQuote: 'Unity in diversity',
    pincode: '110001',
  },
  {
    id: '2',
    orderId: 'MO2024002',
    creator: 'desh1def...uvw',
    receiver: 'suresh@dhan',
    amount: '10000',
    fee: '100',
    status: 'completed',
    createdAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    expiresAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    paymentMode: 'fiat',
    memo: 'Thanks for the help!',
  },
];

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  header: {
    paddingVertical: theme.spacing.lg,
    paddingHorizontal: theme.spacing.lg,
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  headerSubtitle: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
    marginTop: 4,
  },
  tabs: {
    flexDirection: 'row',
    backgroundColor: COLORS.gray100,
    marginHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    padding: 4,
  },
  tab: {
    flex: 1,
    paddingVertical: theme.spacing.sm,
    alignItems: 'center',
    borderRadius: theme.borderRadius.sm,
  },
  activeTab: {
    backgroundColor: COLORS.white,
    ...theme.shadows.small,
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
  content: {
    flex: 1,
    paddingTop: theme.spacing.md,
  },
  formContainer: {
    paddingHorizontal: theme.spacing.lg,
    paddingBottom: theme.spacing.xl,
  },
  formTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.lg,
  },
  inputGroup: {
    marginBottom: theme.spacing.lg,
  },
  inputLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.gray700,
    marginBottom: theme.spacing.sm,
  },
  inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    paddingHorizontal: theme.spacing.md,
    height: 48,
  },
  input: {
    flex: 1,
    fontSize: 16,
    color: COLORS.gray900,
    marginLeft: theme.spacing.sm,
  },
  inputHint: {
    fontSize: 12,
    color: COLORS.gray500,
    marginTop: 4,
  },
  currencySymbol: {
    fontSize: 18,
    color: COLORS.gray700,
    fontWeight: 'bold',
  },
  paymentModes: {
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  paymentMode: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    backgroundColor: COLORS.gray100,
    gap: theme.spacing.sm,
  },
  paymentModeActive: {
    backgroundColor: COLORS.saffron,
  },
  paymentModeText: {
    fontSize: 14,
    color: COLORS.gray700,
    fontWeight: '500',
  },
  paymentModeTextActive: {
    color: COLORS.white,
  },
  memoInput: {
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    fontSize: 16,
    color: COLORS.gray900,
    minHeight: 80,
    textAlignVertical: 'top',
  },
  feeContainer: {
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  feeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.sm,
  },
  totalRow: {
    borderTopWidth: 1,
    borderTopColor: COLORS.gray200,
    paddingTop: theme.spacing.sm,
    marginBottom: 0,
  },
  feeLabel: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  feeAmount: {
    fontSize: 14,
    color: COLORS.gray800,
    fontWeight: '500',
  },
  festivalLabel: {
    fontSize: 14,
    color: COLORS.festivalPrimary,
  },
  festivalDiscount: {
    fontSize: 14,
    color: COLORS.festivalPrimary,
    fontWeight: '600',
  },
  totalLabel: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '600',
  },
  totalAmount: {
    fontSize: 18,
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  quoteContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.saffron + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  quoteText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.saffron,
    marginLeft: theme.spacing.sm,
    fontStyle: 'italic',
  },
  createButton: {
    marginTop: theme.spacing.md,
  },
  ordersList: {
    padding: theme.spacing.lg,
  },
  orderCard: {
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.md,
    ...theme.shadows.medium,
  },
  orderHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.sm,
  },
  orderIdContainer: {
    flex: 1,
  },
  orderId: {
    fontSize: 14,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  statusBadge: {
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: 2,
    borderRadius: theme.borderRadius.sm,
    marginTop: 4,
    alignSelf: 'flex-start',
  },
  statuspending: {
    backgroundColor: COLORS.warning + '20',
  },
  statusaccepted: {
    backgroundColor: COLORS.info + '20',
  },
  statuscompleted: {
    backgroundColor: COLORS.success + '20',
  },
  statuscancelled: {
    backgroundColor: COLORS.error + '20',
  },
  statusexpired: {
    backgroundColor: COLORS.gray300,
  },
  statusText: {
    fontSize: 10,
    fontWeight: 'bold',
    color: COLORS.gray800,
  },
  orderAmount: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.saffron,
  },
  orderDetails: {
    marginBottom: theme.spacing.sm,
  },
  orderRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 4,
  },
  orderText: {
    fontSize: 12,
    color: COLORS.gray600,
    marginLeft: theme.spacing.sm,
    flex: 1,
  },
  orderQuote: {
    fontSize: 12,
    color: COLORS.gray500,
    fontStyle: 'italic',
    marginTop: theme.spacing.sm,
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
  },
  emptySubtext: {
    fontSize: 14,
    color: COLORS.gray500,
    marginTop: theme.spacing.sm,
  },
});