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
  Share,
  Alert,
  Clipboard,
  ActivityIndicator,
  Linking,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation, useRoute } from '@react-navigation/native';
import ViewShot from 'react-native-view-shot';
import { CameraRoll } from '@react-native-camera-roll/camera-roll';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { useAppSelector } from '@store/index';
import { DeshChainClient } from '@services/blockchain/deshchainClient';

interface RouteParams {
  txHash: string;
}

interface TransactionDetails {
  hash: string;
  type: 'send' | 'receive' | 'stake' | 'unstake' | 'swap' | 'reward' | 'moneyorder' | 'sikkebaaz' | 'suraksha';
  status: 'success' | 'pending' | 'failed';
  from: string;
  fromDhanPata?: string;
  to: string;
  toDhanPata?: string;
  amount: string;
  symbol: string;
  fee: string;
  timestamp: string;
  blockHeight: number;
  memo?: string;
  culturalQuote?: string;
  festivalBonus?: {
    festival: string;
    discount: string;
  };
  additionalInfo?: {
    orderId?: string;
    launchId?: string;
    accountId?: string;
  };
}

export const TransactionDetailsScreen: React.FC = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const viewShotRef = React.useRef<ViewShot>(null);
  
  const params = route.params as RouteParams;
  const txHash = params?.txHash;
  
  const [transaction, setTransaction] = useState<TransactionDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchTransactionDetails();
  }, [txHash]);

  const fetchTransactionDetails = async () => {
    try {
      // Simulate API call - in real app, fetch from blockchain
      setTimeout(() => {
        setTransaction(getMockTransaction());
        setLoading(false);
      }, 1000);
    } catch (err) {
      setError('Failed to fetch transaction details');
      setLoading(false);
    }
  };

  const getMockTransaction = (): TransactionDetails => {
    return {
      hash: txHash,
      type: 'send',
      status: 'success',
      from: currentAddress || 'desh1abc...',
      fromDhanPata: dhanPataAddress,
      to: 'desh1xyz...',
      toDhanPata: 'ramesh@dhan',
      amount: '1000',
      symbol: 'NAMO',
      fee: '0.025',
      timestamp: new Date().toISOString(),
      blockHeight: 12345678,
      memo: 'Thanks for lunch!',
      culturalQuote: 'दान देने वाला हाथ लेने वाले हाथ से ऊपर होता है - The giving hand is above the receiving hand',
      festivalBonus: {
        festival: 'Diwali',
        discount: '20%',
      },
    };
  };

  const copyToClipboard = (text: string, label: string) => {
    Clipboard.setString(text);
    Alert.alert('Copied!', `${label} copied to clipboard`);
  };

  const shareTransaction = async () => {
    if (!transaction) return;
    
    const message = `DeshChain Transaction

Amount: ${transaction.amount} ${transaction.symbol}
From: ${transaction.fromDhanPata || transaction.from}
To: ${transaction.toDhanPata || transaction.to}
Status: ${transaction.status}
Hash: ${transaction.hash}

View on DeshChain Explorer:
https://explorer.deshchain.com/tx/${transaction.hash}`;
    
    try {
      await Share.share({
        message,
        title: 'Transaction Details',
      });
    } catch (error) {
      console.error('Share error:', error);
    }
  };

  const saveReceipt = async () => {
    try {
      const uri = await viewShotRef.current?.capture?.();
      if (uri) {
        await CameraRoll.save(uri, { type: 'photo' });
        Alert.alert('Saved!', 'Transaction receipt saved to gallery');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to save receipt');
    }
  };

  const openExplorer = () => {
    const url = `https://explorer.deshchain.com/tx/${txHash}`;
    Linking.openURL(url).catch(() => {
      Alert.alert('Error', 'Failed to open explorer');
    });
  };

  const getTransactionIcon = () => {
    if (!transaction) return 'help-circle';
    
    switch (transaction.type) {
      case 'send': return 'arrow-up-circle';
      case 'receive': return 'arrow-down-circle';
      case 'stake': return 'lock';
      case 'unstake': return 'lock-open';
      case 'swap': return 'swap-horizontal';
      case 'reward': return 'gift';
      case 'moneyorder': return 'cash-fast';
      case 'sikkebaaz': return 'rocket-launch';
      case 'suraksha': return 'shield-check';
      default: return 'circle';
    }
  };

  const getTransactionColor = () => {
    if (!transaction) return COLORS.gray600;
    
    switch (transaction.type) {
      case 'send': return COLORS.error;
      case 'receive': return COLORS.success;
      case 'stake': return COLORS.saffron;
      case 'unstake': return COLORS.warning;
      case 'reward': return COLORS.festivalPrimary;
      default: return COLORS.gray700;
    }
  };

  const getStatusColor = () => {
    if (!transaction) return COLORS.gray600;
    
    switch (transaction.status) {
      case 'success': return COLORS.success;
      case 'pending': return COLORS.warning;
      case 'failed': return COLORS.error;
      default: return COLORS.gray600;
    }
  };

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={COLORS.saffron} />
          <Text style={styles.loadingText}>Fetching transaction details...</Text>
        </View>
      </SafeAreaView>
    );
  }

  if (error || !transaction) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.errorContainer}>
          <Icon name="alert-circle" size={64} color={COLORS.error} />
          <Text style={styles.errorText}>{error || 'Transaction not found'}</Text>
          <CulturalButton
            title="Go Back"
            onPress={() => navigation.goBack()}
            variant="outline"
            size="medium"
            style={styles.errorButton}
          />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Icon name="arrow-left" size={24} color={COLORS.gray900} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Transaction Details</Text>
        <TouchableOpacity onPress={shareTransaction}>
          <Icon name="share-variant" size={24} color={COLORS.gray900} />
        </TouchableOpacity>
      </View>

      <ScrollView showsVerticalScrollIndicator={false}>
        <ViewShot ref={viewShotRef} style={styles.receiptContainer}>
          {/* Transaction Status */}
          <LinearGradient
            colors={transaction.status === 'success' 
              ? [COLORS.success, COLORS.darkGreen]
              : transaction.status === 'pending'
              ? [COLORS.warning, COLORS.orange]
              : [COLORS.error, COLORS.darkRed]
            }
            style={styles.statusCard}
          >
            <Icon 
              name={getTransactionIcon()} 
              size={64} 
              color={COLORS.white} 
            />
            <Text style={styles.statusText}>
              {transaction.type.charAt(0).toUpperCase() + transaction.type.slice(1)} {transaction.status}
            </Text>
            <GradientText style={styles.amountText}>
              {transaction.type === 'send' ? '-' : '+'}{transaction.amount} {transaction.symbol}
            </GradientText>
          </LinearGradient>

          {/* Transaction Info */}
          <View style={styles.infoSection}>
            <View style={styles.infoItem}>
              <Text style={styles.infoLabel}>From</Text>
              <TouchableOpacity
                onPress={() => copyToClipboard(
                  transaction.from,
                  transaction.fromDhanPata ? 'DhanPata ID' : 'Address'
                )}
              >
                <Text style={styles.infoValue}>
                  {transaction.fromDhanPata || `${transaction.from.slice(0, 10)}...${transaction.from.slice(-8)}`}
                </Text>
                {transaction.fromDhanPata && (
                  <Text style={styles.infoSubvalue}>{transaction.from}</Text>
                )}
              </TouchableOpacity>
            </View>

            <Icon name="arrow-down" size={24} color={COLORS.gray400} style={styles.arrowIcon} />

            <View style={styles.infoItem}>
              <Text style={styles.infoLabel}>To</Text>
              <TouchableOpacity
                onPress={() => copyToClipboard(
                  transaction.to,
                  transaction.toDhanPata ? 'DhanPata ID' : 'Address'
                )}
              >
                <Text style={styles.infoValue}>
                  {transaction.toDhanPata || `${transaction.to.slice(0, 10)}...${transaction.to.slice(-8)}`}
                </Text>
                {transaction.toDhanPata && (
                  <Text style={styles.infoSubvalue}>{transaction.to}</Text>
                )}
              </TouchableOpacity>
            </View>

            <View style={styles.divider} />

            <View style={styles.infoRow}>
              <View style={[styles.infoItem, styles.infoItemHalf]}>
                <Text style={styles.infoLabel}>Fee</Text>
                <Text style={styles.infoValue}>{transaction.fee} {transaction.symbol}</Text>
              </View>
              <View style={[styles.infoItem, styles.infoItemHalf]}>
                <Text style={styles.infoLabel}>Block Height</Text>
                <Text style={styles.infoValue}>#{transaction.blockHeight.toLocaleString()}</Text>
              </View>
            </View>

            <View style={styles.infoItem}>
              <Text style={styles.infoLabel}>Date & Time</Text>
              <Text style={styles.infoValue}>
                {new Date(transaction.timestamp).toLocaleString()}
              </Text>
            </View>

            {transaction.memo && (
              <View style={styles.infoItem}>
                <Text style={styles.infoLabel}>Memo</Text>
                <Text style={styles.infoValue}>{transaction.memo}</Text>
              </View>
            )}

            {/* Festival Bonus */}
            {transaction.festivalBonus && (
              <View style={styles.festivalCard}>
                <Icon name="party-popper" size={24} color={COLORS.festivalPrimary} />
                <View style={styles.festivalContent}>
                  <Text style={styles.festivalTitle}>
                    {transaction.festivalBonus.festival} Bonus Applied!
                  </Text>
                  <Text style={styles.festivalText}>
                    You saved {transaction.festivalBonus.discount} on fees
                  </Text>
                </View>
              </View>
            )}

            {/* Transaction Hash */}
            <View style={styles.hashContainer}>
              <Text style={styles.hashLabel}>Transaction Hash</Text>
              <TouchableOpacity
                style={styles.hashButton}
                onPress={() => copyToClipboard(transaction.hash, 'Transaction hash')}
              >
                <Text style={styles.hashText} numberOfLines={1}>
                  {transaction.hash}
                </Text>
                <Icon name="content-copy" size={16} color={COLORS.gray600} />
              </TouchableOpacity>
            </View>
          </View>

          {/* Cultural Quote */}
          {transaction.culturalQuote && (
            <QuoteCard
              quote={transaction.culturalQuote.split(' - ')[0]}
              translation={transaction.culturalQuote.split(' - ')[1]}
              language="hi"
              variant="minimal"
              style={styles.quoteCard}
            />
          )}
        </ViewShot>

        {/* Action Buttons */}
        <View style={styles.actionButtons}>
          <CulturalButton
            title="View on Explorer"
            onPress={openExplorer}
            variant="outline"
            size="medium"
            style={styles.actionButton}
            icon="open-in-new"
          />
          <CulturalButton
            title="Save Receipt"
            onPress={saveReceipt}
            variant="secondary"
            size="medium"
            style={styles.actionButton}
            icon="download"
          />
        </View>

        {/* Related Actions */}
        {transaction.type === 'receive' && (
          <TouchableOpacity
            style={styles.relatedAction}
            onPress={() => navigation.navigate('Send' as any, {
              recipient: transaction.fromDhanPata || transaction.from,
            })}
          >
            <Icon name="reply" size={24} color={COLORS.saffron} />
            <Text style={styles.relatedActionText}>Send back to sender</Text>
            <Icon name="chevron-right" size={20} color={COLORS.gray400} />
          </TouchableOpacity>
        )}
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    fontSize: 16,
    color: COLORS.gray600,
    marginTop: theme.spacing.md,
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.xl,
  },
  errorText: {
    fontSize: 18,
    color: COLORS.error,
    textAlign: 'center',
    marginVertical: theme.spacing.lg,
  },
  errorButton: {
    paddingHorizontal: theme.spacing.xl,
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
  receiptContainer: {
    backgroundColor: COLORS.white,
  },
  statusCard: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
    marginHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
  },
  statusText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.white,
    marginTop: theme.spacing.md,
    textTransform: 'capitalize',
  },
  amountText: {
    fontSize: 32,
    fontWeight: 'bold',
    marginTop: theme.spacing.sm,
  },
  infoSection: {
    margin: theme.spacing.lg,
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
  },
  infoItem: {
    marginBottom: theme.spacing.lg,
  },
  infoItemHalf: {
    flex: 1,
  },
  infoRow: {
    flexDirection: 'row',
    gap: theme.spacing.lg,
  },
  infoLabel: {
    fontSize: 12,
    color: COLORS.gray600,
    marginBottom: 4,
    textTransform: 'uppercase',
    fontWeight: '600',
  },
  infoValue: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  infoSubvalue: {
    fontSize: 12,
    color: COLORS.gray500,
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
  festivalCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.festivalPrimary + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginTop: theme.spacing.md,
  },
  festivalContent: {
    flex: 1,
    marginLeft: theme.spacing.md,
  },
  festivalTitle: {
    fontSize: 14,
    fontWeight: 'bold',
    color: COLORS.festivalPrimary,
  },
  festivalText: {
    fontSize: 12,
    color: COLORS.festivalSecondary,
    marginTop: 2,
  },
  hashContainer: {
    marginTop: theme.spacing.lg,
  },
  hashLabel: {
    fontSize: 12,
    color: COLORS.gray600,
    marginBottom: 4,
    textTransform: 'uppercase',
    fontWeight: '600',
  },
  hashButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.sm,
    padding: theme.spacing.sm,
  },
  hashText: {
    flex: 1,
    fontSize: 12,
    color: COLORS.gray700,
    fontFamily: Platform.OS === 'ios' ? 'Courier New' : 'monospace',
  },
  quoteCard: {
    marginHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
  },
  actionButtons: {
    flexDirection: 'row',
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
    gap: theme.spacing.md,
  },
  actionButton: {
    flex: 1,
  },
  relatedAction: {
    flexDirection: 'row',
    alignItems: 'center',
    marginHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.xl,
    padding: theme.spacing.md,
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
  },
  relatedActionText: {
    flex: 1,
    fontSize: 16,
    color: COLORS.gray900,
    marginLeft: theme.spacing.md,
  },
});