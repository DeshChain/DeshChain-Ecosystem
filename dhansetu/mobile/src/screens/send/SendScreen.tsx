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
  TextInput,
  TouchableOpacity,
  Alert,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation, useRoute } from '@react-navigation/native';
import QRCodeScanner from 'react-native-qrcode-scanner';
import { RNCamera } from 'react-native-camera';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { useAppSelector, useAppDispatch } from '@store/index';
import { sendTransaction } from '@store/slices/walletSlice';
import { DeshChainClient } from '@services/blockchain/deshchainClient';
import { useFestival } from '@contexts/FestivalContext';

interface RouteParams {
  coinType?: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN';
  recipient?: string;
  amount?: string;
}

interface TransactionFee {
  base: string;
  festival?: string;
  total: string;
}

export const SendScreen: React.FC = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const dispatch = useAppDispatch();
  const { currentAddress, currentCoin, balances, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival } = useFestival();
  
  const params = route.params as RouteParams;
  
  const [recipient, setRecipient] = useState(params?.recipient || '');
  const [amount, setAmount] = useState(params?.amount || '');
  const [memo, setMemo] = useState('');
  const [culturalQuote, setCulturalQuote] = useState('');
  const [isScanning, setIsScanning] = useState(false);
  const [isValidating, setIsValidating] = useState(false);
  const [recipientValid, setRecipientValid] = useState<boolean | null>(null);
  const [recipientInfo, setRecipientInfo] = useState<{
    address: string;
    name?: string;
    pincode?: string;
    isDhanPata: boolean;
  } | null>(null);
  const [fee, setFee] = useState<TransactionFee>({
    base: '0.025',
    festival: currentFestival ? '0.005' : undefined,
    total: currentFestival ? '0.02' : '0.025',
  });
  const [sending, setSending] = useState(false);
  const [reviewMode, setReviewMode] = useState(false);

  // Quick amount buttons
  const quickAmounts = ['100', '500', '1000', '5000'];
  
  // Cultural quotes for transactions
  const transactionQuotes = [
    'दान देने वाला हाथ लेने वाले हाथ से ऊपर होता है - The giving hand is above the receiving hand',
    'जो देता है वो देवता है - One who gives is divine',
    'पैसा हाथ का मैल है - Money is meant to flow',
    'सबका साथ, सबका विकास - Together we prosper',
  ];

  useEffect(() => {
    // Set random cultural quote
    setCulturalQuote(transactionQuotes[Math.floor(Math.random() * transactionQuotes.length)]);
  }, []);

  useEffect(() => {
    if (recipient.length > 0) {
      validateRecipient();
    }
  }, [recipient]);

  const validateRecipient = async () => {
    setIsValidating(true);
    setRecipientValid(null);
    
    try {
      // Check if it's a DhanPata address
      const isDhanPata = recipient.includes('@dhan');
      
      if (isDhanPata) {
        // Validate DhanPata format
        const [username, domain] = recipient.split('@');
        if (domain !== 'dhan' || username.length < 3) {
          setRecipientValid(false);
          setIsValidating(false);
          return;
        }
        
        // Simulate DhanPata lookup
        setTimeout(() => {
          setRecipientValid(true);
          setRecipientInfo({
            address: 'desh1abc...' + username.slice(0, 4),
            name: username.charAt(0).toUpperCase() + username.slice(1),
            isDhanPata: true,
            pincode: '110001',
          });
          setIsValidating(false);
        }, 500);
      } else {
        // Validate blockchain address
        const isValidAddress = recipient.startsWith('desh1') && recipient.length === 44;
        setRecipientValid(isValidAddress);
        if (isValidAddress) {
          setRecipientInfo({
            address: recipient,
            isDhanPata: false,
          });
        }
        setIsValidating(false);
      }
    } catch (error) {
      setRecipientValid(false);
      setIsValidating(false);
    }
  };

  const handleQRScan = (e: any) => {
    setIsScanning(false);
    const data = e.data;
    
    // Parse QR data (could be address or payment request)
    if (data.startsWith('deshchain:')) {
      const [address, params] = data.replace('deshchain:', '').split('?');
      setRecipient(address);
      
      if (params) {
        const urlParams = new URLSearchParams(params);
        if (urlParams.get('amount')) {
          setAmount(urlParams.get('amount')!);
        }
        if (urlParams.get('memo')) {
          setMemo(urlParams.get('memo')!);
        }
      }
    } else {
      setRecipient(data);
    }
  };

  const calculateTotal = () => {
    const sendAmount = parseFloat(amount) || 0;
    const feeAmount = parseFloat(fee.total);
    return (sendAmount + feeAmount).toFixed(4);
  };

  const validateTransaction = () => {
    if (!recipientValid) {
      Alert.alert('Invalid Recipient', 'Please enter a valid address or DhanPata ID');
      return false;
    }
    
    const sendAmount = parseFloat(amount);
    if (!sendAmount || sendAmount <= 0) {
      Alert.alert('Invalid Amount', 'Please enter a valid amount');
      return false;
    }
    
    const balance = parseFloat(balances[currentCoin] || '0');
    const total = parseFloat(calculateTotal());
    
    if (total > balance) {
      Alert.alert('Insufficient Balance', `You need ${total} ${currentCoin} but only have ${balance}`);
      return false;
    }
    
    return true;
  };

  const handleReview = () => {
    if (validateTransaction()) {
      setReviewMode(true);
    }
  };

  const handleSend = async () => {
    setSending(true);
    
    try {
      const result = await dispatch(sendTransaction({
        recipient: recipientInfo!.address,
        amount,
        memo,
        culturalQuote,
        coinType: currentCoin,
      })).unwrap();
      
      Alert.alert(
        'Success!',
        `Sent ${amount} ${currentCoin} to ${recipientInfo?.name || recipient}`,
        [
          {
            text: 'View Transaction',
            onPress: () => navigation.navigate('TransactionDetails', { txHash: result.transactionHash }),
          },
          {
            text: 'Done',
            onPress: () => navigation.goBack(),
          },
        ]
      );
    } catch (error: any) {
      Alert.alert('Transaction Failed', error.message || 'Please try again');
    } finally {
      setSending(false);
    }
  };

  const renderScannerView = () => (
    <View style={styles.scannerContainer}>
      <QRCodeScanner
        onRead={handleQRScan}
        flashMode={RNCamera.Constants.FlashMode.auto}
        topContent={
          <Text style={styles.scannerText}>
            Scan recipient's QR code
          </Text>
        }
        bottomContent={
          <TouchableOpacity
            style={styles.scannerButton}
            onPress={() => setIsScanning(false)}
          >
            <Text style={styles.scannerButtonText}>Cancel</Text>
          </TouchableOpacity>
        }
      />
    </View>
  );

  const renderTransactionForm = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.formContainer}>
        {/* Recipient Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Send To</Text>
          <View style={styles.recipientInputContainer}>
            <TextInput
              style={styles.recipientInput}
              placeholder="Enter address or username@dhan"
              value={recipient}
              onChangeText={setRecipient}
              autoCapitalize="none"
              autoCorrect={false}
            />
            <TouchableOpacity
              style={styles.scanButton}
              onPress={() => setIsScanning(true)}
            >
              <Icon name="qrcode-scan" size={24} color={COLORS.saffron} />
            </TouchableOpacity>
          </View>
          
          {isValidating && (
            <View style={styles.validationRow}>
              <ActivityIndicator size="small" color={COLORS.saffron} />
              <Text style={styles.validatingText}>Validating...</Text>
            </View>
          )}
          
          {recipientValid !== null && !isValidating && (
            <View style={styles.validationRow}>
              <Icon
                name={recipientValid ? 'check-circle' : 'close-circle'}
                size={20}
                color={recipientValid ? COLORS.success : COLORS.error}
              />
              <Text style={[
                styles.validationText,
                { color: recipientValid ? COLORS.success : COLORS.error }
              ]}>
                {recipientValid
                  ? recipientInfo?.isDhanPata
                    ? `DhanPata: ${recipientInfo.name}`
                    : 'Valid address'
                  : 'Invalid recipient'
                }
              </Text>
            </View>
          )}
        </View>

        {/* Amount Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Amount</Text>
          <View style={styles.amountInputContainer}>
            <Text style={styles.currencySymbol}>{currentCoin === 'DESHCHAIN' ? '₹' : ''}</Text>
            <TextInput
              style={styles.amountInput}
              placeholder="0.00"
              value={amount}
              onChangeText={setAmount}
              keyboardType="decimal-pad"
            />
            <Text style={styles.currencyCode}>{currentCoin}</Text>
          </View>
          
          <View style={styles.balanceRow}>
            <Text style={styles.balanceText}>
              Available: {balances[currentCoin] || '0'} {currentCoin}
            </Text>
            <TouchableOpacity onPress={() => setAmount(balances[currentCoin] || '0')}>
              <Text style={styles.maxButton}>MAX</Text>
            </TouchableOpacity>
          </View>
          
          {/* Quick Amount Buttons */}
          <View style={styles.quickAmounts}>
            {quickAmounts.map((quickAmount) => (
              <TouchableOpacity
                key={quickAmount}
                style={styles.quickAmountButton}
                onPress={() => setAmount(quickAmount)}
              >
                <Text style={styles.quickAmountText}>₹{quickAmount}</Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        {/* Memo Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Memo (Optional)</Text>
          <TextInput
            style={styles.memoInput}
            placeholder="Add a note"
            value={memo}
            onChangeText={setMemo}
            multiline
            numberOfLines={3}
            textAlignVertical="top"
          />
        </View>

        {/* Festival Bonus */}
        {currentFestival && (
          <LinearGradient
            colors={theme.gradients.festivalGradient}
            style={styles.festivalCard}
          >
            <Icon name="party-popper" size={24} color={COLORS.white} />
            <View style={styles.festivalContent}>
              <Text style={styles.festivalTitle}>
                {currentFestival.name} Special!
              </Text>
              <Text style={styles.festivalText}>
                Reduced fees during festival period
              </Text>
            </View>
            <Text style={styles.festivalDiscount}>-20%</Text>
          </LinearGradient>
        )}

        {/* Cultural Quote */}
        <QuoteCard
          quote={culturalQuote.split(' - ')[0]}
          translation={culturalQuote.split(' - ')[1]}
          language="hi"
          variant="minimal"
        />

        {/* Fee Summary */}
        <View style={styles.feeSummary}>
          <View style={styles.feeRow}>
            <Text style={styles.feeLabel}>Transaction Fee</Text>
            <Text style={styles.feeValue}>{fee.base} {currentCoin}</Text>
          </View>
          {fee.festival && (
            <View style={styles.feeRow}>
              <Text style={styles.feeLabel}>Festival Discount</Text>
              <Text style={[styles.feeValue, styles.discountValue]}>-{fee.festival} {currentCoin}</Text>
            </View>
          )}
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Total</Text>
            <Text style={styles.totalValue}>{calculateTotal()} {currentCoin}</Text>
          </View>
        </View>

        {/* Send Button */}
        <CulturalButton
          title="Review Transaction"
          onPress={handleReview}
          size="large"
          style={styles.sendButton}
          disabled={!recipientValid || !amount}
        />
      </View>
    </ScrollView>
  );

  const renderReviewScreen = () => (
    <View style={styles.reviewContainer}>
      <LinearGradient
        colors={[COLORS.saffron, COLORS.orange]}
        style={styles.reviewHeader}
      >
        <Icon name="shield-check" size={48} color={COLORS.white} />
        <Text style={styles.reviewTitle}>Review Transaction</Text>
      </LinearGradient>

      <View style={styles.reviewDetails}>
        <View style={styles.reviewItem}>
          <Text style={styles.reviewLabel}>From</Text>
          <View style={styles.reviewValue}>
            <Text style={styles.reviewAddress}>
              {dhanPataAddress || `${currentAddress?.slice(0, 10)}...${currentAddress?.slice(-8)}`}
            </Text>
            <Text style={styles.reviewBalance}>
              Balance: {balances[currentCoin]} {currentCoin}
            </Text>
          </View>
        </View>

        <Icon name="arrow-down" size={24} color={COLORS.gray400} style={styles.arrowIcon} />

        <View style={styles.reviewItem}>
          <Text style={styles.reviewLabel}>To</Text>
          <View style={styles.reviewValue}>
            <Text style={styles.reviewAddress}>{recipient}</Text>
            {recipientInfo?.name && (
              <Text style={styles.reviewName}>{recipientInfo.name}</Text>
            )}
          </View>
        </View>

        <View style={styles.divider} />

        <View style={styles.reviewItem}>
          <Text style={styles.reviewLabel}>Amount</Text>
          <GradientText style={styles.reviewAmount}>
            {amount} {currentCoin}
          </GradientText>
        </View>

        <View style={styles.reviewItem}>
          <Text style={styles.reviewLabel}>Fee</Text>
          <Text style={styles.reviewFee}>{fee.total} {currentCoin}</Text>
        </View>

        {memo && (
          <View style={styles.reviewItem}>
            <Text style={styles.reviewLabel}>Memo</Text>
            <Text style={styles.reviewMemo}>{memo}</Text>
          </View>
        )}
      </View>

      <View style={styles.reviewActions}>
        <CulturalButton
          title="Edit"
          onPress={() => setReviewMode(false)}
          variant="outline"
          size="large"
          style={styles.editButton}
        />
        <CulturalButton
          title="Confirm & Send"
          onPress={handleSend}
          size="large"
          style={styles.confirmButton}
          loading={sending}
        />
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardAvoid}
      >
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity onPress={() => navigation.goBack()}>
            <Icon name="arrow-left" size={24} color={COLORS.gray900} />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Send {currentCoin}</Text>
          <View style={{ width: 24 }} />
        </View>

        {/* Content */}
        {isScanning ? renderScannerView() : 
         reviewMode ? renderReviewScreen() : renderTransactionForm()}
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  keyboardAvoid: {
    flex: 1,
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
  
  // Form styles
  formContainer: {
    padding: theme.spacing.lg,
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
  recipientInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    paddingHorizontal: theme.spacing.md,
    height: 48,
  },
  recipientInput: {
    flex: 1,
    fontSize: 16,
    color: COLORS.gray900,
  },
  scanButton: {
    padding: theme.spacing.sm,
  },
  validationRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: theme.spacing.sm,
  },
  validatingText: {
    fontSize: 12,
    color: COLORS.gray600,
    marginLeft: theme.spacing.sm,
  },
  validationText: {
    fontSize: 12,
    marginLeft: theme.spacing.xs,
  },
  amountInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    paddingHorizontal: theme.spacing.md,
    height: 56,
  },
  currencySymbol: {
    fontSize: 24,
    color: COLORS.gray700,
    fontWeight: 'bold',
  },
  amountInput: {
    flex: 1,
    fontSize: 24,
    color: COLORS.gray900,
    fontWeight: 'bold',
    marginHorizontal: theme.spacing.sm,
  },
  currencyCode: {
    fontSize: 16,
    color: COLORS.gray600,
    fontWeight: '600',
  },
  balanceRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: theme.spacing.sm,
  },
  balanceText: {
    fontSize: 12,
    color: COLORS.gray600,
  },
  maxButton: {
    fontSize: 12,
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  quickAmounts: {
    flexDirection: 'row',
    gap: theme.spacing.sm,
    marginTop: theme.spacing.md,
  },
  quickAmountButton: {
    flex: 1,
    paddingVertical: theme.spacing.sm,
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.sm,
    alignItems: 'center',
  },
  quickAmountText: {
    fontSize: 14,
    color: COLORS.gray700,
    fontWeight: '500',
  },
  memoInput: {
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    fontSize: 16,
    color: COLORS.gray900,
    minHeight: 80,
  },
  festivalCard: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.lg,
  },
  festivalContent: {
    flex: 1,
    marginLeft: theme.spacing.md,
  },
  festivalTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  festivalText: {
    fontSize: 12,
    color: COLORS.white,
    opacity: 0.9,
    marginTop: 2,
  },
  festivalDiscount: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  feeSummary: {
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
  feeLabel: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  feeValue: {
    fontSize: 14,
    color: COLORS.gray700,
    fontWeight: '500',
  },
  discountValue: {
    color: COLORS.success,
  },
  totalRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingTop: theme.spacing.sm,
    borderTopWidth: 1,
    borderTopColor: COLORS.gray200,
  },
  totalLabel: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  totalValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  sendButton: {
    marginTop: theme.spacing.md,
  },
  
  // Scanner styles
  scannerContainer: {
    flex: 1,
  },
  scannerText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
    textAlign: 'center',
    paddingHorizontal: theme.spacing.lg,
  },
  scannerButton: {
    paddingHorizontal: theme.spacing.xl,
    paddingVertical: theme.spacing.md,
  },
  scannerButtonText: {
    fontSize: 16,
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  
  // Review styles
  reviewContainer: {
    flex: 1,
  },
  reviewHeader: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
  },
  reviewTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
    marginTop: theme.spacing.md,
  },
  reviewDetails: {
    flex: 1,
    padding: theme.spacing.lg,
  },
  reviewItem: {
    marginBottom: theme.spacing.lg,
  },
  reviewLabel: {
    fontSize: 14,
    color: COLORS.gray600,
    marginBottom: theme.spacing.xs,
  },
  reviewValue: {},
  reviewAddress: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  reviewBalance: {
    fontSize: 12,
    color: COLORS.gray600,
    marginTop: 2,
  },
  reviewName: {
    fontSize: 14,
    color: COLORS.gray700,
    marginTop: 2,
  },
  arrowIcon: {
    alignSelf: 'center',
    marginVertical: theme.spacing.sm,
  },
  divider: {
    height: 1,
    backgroundColor: COLORS.gray200,
    marginVertical: theme.spacing.lg,
  },
  reviewAmount: {
    fontSize: 28,
    fontWeight: 'bold',
  },
  reviewFee: {
    fontSize: 16,
    color: COLORS.gray700,
  },
  reviewMemo: {
    fontSize: 14,
    color: COLORS.gray700,
    fontStyle: 'italic',
  },
  reviewActions: {
    flexDirection: 'row',
    padding: theme.spacing.lg,
    gap: theme.spacing.md,
  },
  editButton: {
    flex: 1,
  },
  confirmButton: {
    flex: 2,
  },
});