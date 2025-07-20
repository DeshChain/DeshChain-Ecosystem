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
  withRepeat,
  withTiming,
  Easing,
} from 'react-native-reanimated';
import { LineChart } from 'react-native-chart-kit';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { useAppSelector } from '@store/index';
import { useFestival } from '@contexts/FestivalContext';

const { width } = Dimensions.get('window');

interface SurakshaAccount {
  id: string;
  accountHolder: string;
  principalAmount: string;
  maturityAmount: string;
  monthlyContribution: string;
  startDate: string;
  maturityDate: string;
  tenureYears: number;
  guaranteedReturn: number;
  currentValue: string;
  status: 'active' | 'matured' | 'withdrawn';
  nomineeInfo: {
    name: string;
    relationship: string;
    percentage: number;
  }[];
}

export const SurakshaScreen: React.FC = () => {
  const navigation = useNavigation();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival } = useFestival();
  
  const [activeTab, setActiveTab] = useState<'overview' | 'calculator' | 'myAccounts'>('overview');
  const [accounts, setAccounts] = useState<SurakshaAccount[]>([]);
  const [loading, setLoading] = useState(false);
  
  // Calculator state
  const [monthlyAmount, setMonthlyAmount] = useState('');
  const [tenureYears, setTenureYears] = useState('10');
  const [calculatedMaturity, setCalculatedMaturity] = useState<{
    totalInvestment: string;
    maturityAmount: string;
    totalReturns: string;
    monthlyPension: string;
  } | null>(null);

  useEffect(() => {
    fetchAccounts();
  }, []);

  const fetchAccounts = async () => {
    setLoading(true);
    // Simulate API call
    setTimeout(() => {
      setAccounts(getMockAccounts());
      setLoading(false);
    }, 1000);
  };

  const calculateMaturity = () => {
    if (!monthlyAmount || parseFloat(monthlyAmount) <= 0) {
      Alert.alert('Error', 'Please enter a valid monthly amount');
      return;
    }

    const monthly = parseFloat(monthlyAmount);
    const years = parseInt(tenureYears);
    const totalMonths = years * 12;
    const guaranteedReturn = 0.5; // 50% guaranteed return
    
    const totalInvestment = monthly * totalMonths;
    const returns = totalInvestment * guaranteedReturn;
    const maturityAmount = totalInvestment + returns;
    
    // Calculate monthly pension (assuming 20 years payout)
    const monthlyPension = maturityAmount / (20 * 12);

    setCalculatedMaturity({
      totalInvestment: totalInvestment.toFixed(0),
      maturityAmount: maturityAmount.toFixed(0),
      totalReturns: returns.toFixed(0),
      monthlyPension: monthlyPension.toFixed(0),
    });
  };

  const renderOverview = () => {
    // Sample data for growth chart
    const chartData = {
      labels: ['Year 1', 'Year 5', 'Year 10', 'Year 15', 'Year 20'],
      datasets: [{
        data: [12000, 72000, 180000, 360000, 600000],
      }],
    };

    return (
      <ScrollView showsVerticalScrollIndicator={false}>
        <View style={styles.overviewContainer}>
          {/* Hero Section */}
          <LinearGradient
            colors={[COLORS.green, COLORS.darkGreen]}
            style={styles.heroSection}
          >
            <Icon name="shield-check" size={64} color={COLORS.white} />
            <Text style={styles.heroTitle}>Gram Suraksha Pool</Text>
            <Text style={styles.heroSubtitle}>
              India's first blockchain-powered pension scheme
            </Text>
            <View style={styles.guaranteeBox}>
              <GradientText style={styles.guaranteeText}>50% GUARANTEED RETURNS</GradientText>
            </View>
          </LinearGradient>

          {/* Key Features */}
          <View style={styles.featuresSection}>
            <Text style={styles.sectionTitle}>Why Choose Gram Suraksha?</Text>
            
            <View style={styles.featureCard}>
              <View style={styles.featureIcon}>
                <Icon name="percent" size={24} color={COLORS.green} />
              </View>
              <View style={styles.featureContent}>
                <Text style={styles.featureTitle}>50% Guaranteed Returns</Text>
                <Text style={styles.featureDescription}>
                  Fixed returns backed by DeshChain protocol
                </Text>
              </View>
            </View>

            <View style={styles.featureCard}>
              <View style={styles.featureIcon}>
                <Icon name="lock" size={24} color={COLORS.saffron} />
              </View>
              <View style={styles.featureContent}>
                <Text style={styles.featureTitle}>100% Secure</Text>
                <Text style={styles.featureDescription}>
                  Smart contract protected, no middlemen
                </Text>
              </View>
            </View>

            <View style={styles.featureCard}>
              <View style={styles.featureIcon}>
                <Icon name="account-group" size={24} color={COLORS.navy} />
              </View>
              <View style={styles.featureContent}>
                <Text style={styles.featureTitle}>Community Powered</Text>
                <Text style={styles.featureDescription}>
                  Pooled investments for better returns
                </Text>
              </View>
            </View>

            <View style={styles.featureCard}>
              <View style={styles.featureIcon}>
                <Icon name="cash-multiple" size={24} color={COLORS.festivalPrimary} />
              </View>
              <View style={styles.featureContent}>
                <Text style={styles.featureTitle}>Flexible Contributions</Text>
                <Text style={styles.featureDescription}>
                  Start with as low as ₹100 per month
                </Text>
              </View>
            </View>
          </View>

          {/* Growth Chart */}
          <View style={styles.chartSection}>
            <Text style={styles.sectionTitle}>Your Money Grows Steadily</Text>
            <LineChart
              data={chartData}
              width={width - 40}
              height={220}
              chartConfig={{
                backgroundColor: COLORS.white,
                backgroundGradientFrom: COLORS.white,
                backgroundGradientTo: COLORS.white,
                decimalPlaces: 0,
                color: (opacity = 1) => `rgba(19, 136, 8, ${opacity})`,
                labelColor: (opacity = 1) => `rgba(0, 0, 0, ${opacity})`,
                style: {
                  borderRadius: 16,
                },
                propsForDots: {
                  r: '6',
                  strokeWidth: '2',
                  stroke: COLORS.green,
                },
              }}
              bezier
              style={styles.chart}
            />
          </View>

          {/* Cultural Quote */}
          <QuoteCard
            quote="बूंद बूंद से सागर बनता है"
            translation="Drop by drop, an ocean is formed"
            author="Indian Proverb"
            language="hi"
            variant="minimal"
          />

          {/* CTA Button */}
          <CulturalButton
            title="Start Your Suraksha Journey"
            onPress={() => navigation.navigate('EnrollSuraksha')}
            size="large"
            style={styles.ctaButton}
          />
        </View>
      </ScrollView>
    );
  };

  const renderCalculator = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.calculatorContainer}>
        <Text style={styles.calculatorTitle}>Pension Calculator</Text>
        <Text style={styles.calculatorSubtitle}>
          See how your savings grow with 50% guaranteed returns
        </Text>

        {/* Monthly Contribution */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Monthly Contribution</Text>
          <View style={styles.inputWrapper}>
            <Text style={styles.currencySymbol}>₹</Text>
            <TextInput
              style={styles.input}
              placeholder="Enter amount"
              value={monthlyAmount}
              onChangeText={setMonthlyAmount}
              keyboardType="decimal-pad"
            />
          </View>
          <Text style={styles.inputHint}>Minimum ₹100 per month</Text>
        </View>

        {/* Tenure Selection */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Investment Period</Text>
          <View style={styles.tenureOptions}>
            {['5', '10', '15', '20'].map((years) => (
              <TouchableOpacity
                key={years}
                style={[
                  styles.tenureButton,
                  tenureYears === years && styles.tenureButtonActive,
                ]}
                onPress={() => setTenureYears(years)}
              >
                <Text style={[
                  styles.tenureText,
                  tenureYears === years && styles.tenureTextActive,
                ]}>
                  {years} Years
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        {/* Calculate Button */}
        <CulturalButton
          title="Calculate Returns"
          onPress={calculateMaturity}
          size="large"
          style={styles.calculateButton}
        />

        {/* Results */}
        {calculatedMaturity && (
          <Animated.View style={styles.resultsContainer}>
            <Text style={styles.resultsTitle}>Your Pension Plan</Text>
            
            <View style={styles.resultCard}>
              <Icon name="piggy-bank" size={32} color={COLORS.green} />
              <View style={styles.resultContent}>
                <Text style={styles.resultLabel}>Total Investment</Text>
                <Text style={styles.resultValue}>₹{calculatedMaturity.totalInvestment}</Text>
              </View>
            </View>

            <View style={styles.resultCard}>
              <Icon name="trending-up" size={32} color={COLORS.saffron} />
              <View style={styles.resultContent}>
                <Text style={styles.resultLabel}>Maturity Amount</Text>
                <GradientText style={styles.resultValue}>
                  ₹{calculatedMaturity.maturityAmount}
                </GradientText>
              </View>
            </View>

            <View style={styles.resultCard}>
              <Icon name="cash" size={32} color={COLORS.navy} />
              <View style={styles.resultContent}>
                <Text style={styles.resultLabel}>Total Returns</Text>
                <Text style={[styles.resultValue, styles.returnsValue]}>
                  +₹{calculatedMaturity.totalReturns}
                </Text>
              </View>
            </View>

            <View style={styles.resultCard}>
              <Icon name="calendar-month" size={32} color={COLORS.festivalPrimary} />
              <View style={styles.resultContent}>
                <Text style={styles.resultLabel}>Monthly Pension</Text>
                <Text style={styles.resultValue}>₹{calculatedMaturity.monthlyPension}</Text>
                <Text style={styles.resultHint}>For 20 years after maturity</Text>
              </View>
            </View>

            <CulturalButton
              title="Enroll Now"
              onPress={() => navigation.navigate('EnrollSuraksha')}
              variant="secondary"
              size="large"
              style={styles.enrollButton}
            />
          </Animated.View>
        )}
      </View>
    </ScrollView>
  );

  const renderMyAccounts = () => {
    const totalValue = accounts.reduce((sum, acc) => 
      sum + parseFloat(acc.currentValue), 0
    );

    return (
      <ScrollView showsVerticalScrollIndicator={false}>
        <View style={styles.accountsContainer}>
          {/* Summary Card */}
          <LinearGradient
            colors={[COLORS.green, COLORS.darkGreen]}
            style={styles.summaryCard}
          >
            <Text style={styles.summaryLabel}>Total Portfolio Value</Text>
            <GradientText style={styles.summaryValue}>
              ₹{totalValue.toFixed(0)}
            </GradientText>
            <Text style={styles.summaryAccounts}>
              {accounts.length} Active {accounts.length === 1 ? 'Account' : 'Accounts'}
            </Text>
          </LinearGradient>

          {/* Account List */}
          {accounts.length === 0 ? (
            <View style={styles.emptyState}>
              <Icon name="piggy-bank-outline" size={64} color={COLORS.gray400} />
              <Text style={styles.emptyText}>No Suraksha accounts yet</Text>
              <Text style={styles.emptySubtext}>
                Start your pension journey today
              </Text>
              <CulturalButton
                title="Create First Account"
                onPress={() => navigation.navigate('EnrollSuraksha')}
                variant="outline"
                size="medium"
                style={styles.emptyButton}
              />
            </View>
          ) : (
            accounts.map((account) => (
              <TouchableOpacity
                key={account.id}
                style={styles.accountCard}
                onPress={() => navigation.navigate('SurakshaDetails', { accountId: account.id })}
                activeOpacity={0.8}
              >
                <View style={styles.accountHeader}>
                  <View>
                    <Text style={styles.accountId}>Account #{account.id.slice(-6)}</Text>
                    <Text style={styles.accountTenure}>
                      {account.tenureYears} Year Plan
                    </Text>
                  </View>
                  <View style={[styles.statusBadge, styles[`status_${account.status}`]]}>
                    <Text style={styles.statusText}>{account.status.toUpperCase()}</Text>
                  </View>
                </View>

                <View style={styles.accountDetails}>
                  <View style={styles.detailRow}>
                    <Text style={styles.detailLabel}>Monthly Contribution</Text>
                    <Text style={styles.detailValue}>₹{account.monthlyContribution}</Text>
                  </View>
                  <View style={styles.detailRow}>
                    <Text style={styles.detailLabel}>Current Value</Text>
                    <Text style={styles.detailValue}>₹{account.currentValue}</Text>
                  </View>
                  <View style={styles.detailRow}>
                    <Text style={styles.detailLabel}>Maturity Amount</Text>
                    <Text style={[styles.detailValue, styles.maturityValue]}>
                      ₹{account.maturityAmount}
                    </Text>
                  </View>
                </View>

                <View style={styles.progressBar}>
                  <Animated.View 
                    style={[
                      styles.progressFill,
                      { 
                        width: `${(parseFloat(account.currentValue) / parseFloat(account.maturityAmount)) * 100}%` 
                      }
                    ]} 
                  />
                </View>

                <Text style={styles.maturityDate}>
                  Matures on {new Date(account.maturityDate).toLocaleDateString()}
                </Text>
              </TouchableOpacity>
            ))
          )}
        </View>
      </ScrollView>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Gram Suraksha</Text>
        <Text style={styles.headerSubtitle}>Your Blockchain Pension</Text>
      </View>

      {/* Tabs */}
      <View style={styles.tabs}>
        {(['overview', 'calculator', 'myAccounts'] as const).map((tab) => (
          <TouchableOpacity
            key={tab}
            style={[styles.tab, activeTab === tab && styles.activeTab]}
            onPress={() => setActiveTab(tab)}
          >
            <Text style={[
              styles.tabText,
              activeTab === tab && styles.activeTabText,
            ]}>
              {tab === 'overview' ? 'Overview' : 
               tab === 'calculator' ? 'Calculator' : 'My Accounts'}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {/* Content */}
      <View style={styles.content}>
        {activeTab === 'overview' && renderOverview()}
        {activeTab === 'calculator' && renderCalculator()}
        {activeTab === 'myAccounts' && renderMyAccounts()}
      </View>
    </SafeAreaView>
  );
};

// Mock data
const getMockAccounts = (): SurakshaAccount[] => [
  {
    id: 'GSP001',
    accountHolder: 'user@dhan',
    principalAmount: '120000',
    maturityAmount: '180000',
    monthlyContribution: '1000',
    startDate: new Date(Date.now() - 365 * 24 * 60 * 60 * 1000).toISOString(),
    maturityDate: new Date(Date.now() + 9 * 365 * 24 * 60 * 60 * 1000).toISOString(),
    tenureYears: 10,
    guaranteedReturn: 50,
    currentValue: '12000',
    status: 'active',
    nomineeInfo: [{
      name: 'Spouse',
      relationship: 'Wife',
      percentage: 100,
    }],
  },
];

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  header: {
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    backgroundColor: COLORS.white,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray200,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  headerSubtitle: {
    fontSize: 14,
    color: COLORS.gray600,
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
    color: COLORS.green,
    fontWeight: 'bold',
  },
  content: {
    flex: 1,
    paddingTop: theme.spacing.md,
  },
  
  // Overview styles
  overviewContainer: {
    paddingBottom: theme.spacing.xl,
  },
  heroSection: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
    marginHorizontal: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    marginBottom: theme.spacing.lg,
  },
  heroTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.white,
    marginTop: theme.spacing.md,
  },
  heroSubtitle: {
    fontSize: 16,
    color: COLORS.white,
    opacity: 0.9,
    marginTop: theme.spacing.sm,
    textAlign: 'center',
    paddingHorizontal: theme.spacing.lg,
  },
  guaranteeBox: {
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.full,
    marginTop: theme.spacing.lg,
  },
  guaranteeText: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  featuresSection: {
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.md,
  },
  featureCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.sm,
  },
  featureIcon: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: COLORS.white,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: theme.spacing.md,
    ...theme.shadows.small,
  },
  featureContent: {
    flex: 1,
  },
  featureTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.gray900,
  },
  featureDescription: {
    fontSize: 14,
    color: COLORS.gray600,
    marginTop: 2,
  },
  chartSection: {
    paddingHorizontal: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
  },
  chart: {
    marginVertical: theme.spacing.md,
    borderRadius: theme.borderRadius.lg,
  },
  ctaButton: {
    marginHorizontal: theme.spacing.lg,
    marginTop: theme.spacing.lg,
  },
  
  // Calculator styles
  calculatorContainer: {
    paddingHorizontal: theme.spacing.lg,
    paddingBottom: theme.spacing.xl,
  },
  calculatorTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.sm,
  },
  calculatorSubtitle: {
    fontSize: 14,
    color: COLORS.gray600,
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
  currencySymbol: {
    fontSize: 18,
    color: COLORS.gray700,
    fontWeight: 'bold',
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
  tenureOptions: {
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  tenureButton: {
    flex: 1,
    paddingVertical: theme.spacing.md,
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
  },
  tenureButtonActive: {
    backgroundColor: COLORS.green,
  },
  tenureText: {
    fontSize: 14,
    color: COLORS.gray700,
    fontWeight: '500',
  },
  tenureTextActive: {
    color: COLORS.white,
    fontWeight: 'bold',
  },
  calculateButton: {
    marginTop: theme.spacing.md,
  },
  resultsContainer: {
    marginTop: theme.spacing.xl,
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
  },
  resultsTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.lg,
    textAlign: 'center',
  },
  resultCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.white,
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.md,
    ...theme.shadows.small,
  },
  resultContent: {
    flex: 1,
    marginLeft: theme.spacing.md,
  },
  resultLabel: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  resultValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginTop: 2,
  },
  returnsValue: {
    color: COLORS.green,
  },
  resultHint: {
    fontSize: 12,
    color: COLORS.gray500,
    marginTop: 2,
  },
  enrollButton: {
    marginTop: theme.spacing.md,
  },
  
  // Accounts styles
  accountsContainer: {
    paddingHorizontal: theme.spacing.lg,
    paddingBottom: theme.spacing.xl,
  },
  summaryCard: {
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    alignItems: 'center',
    marginBottom: theme.spacing.lg,
  },
  summaryLabel: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  summaryValue: {
    fontSize: 36,
    fontWeight: 'bold',
    marginVertical: theme.spacing.sm,
  },
  summaryAccounts: {
    fontSize: 14,
    color: COLORS.white,
    opacity: 0.9,
  },
  accountCard: {
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    marginBottom: theme.spacing.md,
    ...theme.shadows.medium,
  },
  accountHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.md,
  },
  accountId: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  accountTenure: {
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
  status_matured: {
    backgroundColor: COLORS.info + '20',
  },
  status_withdrawn: {
    backgroundColor: COLORS.gray300,
  },
  statusText: {
    fontSize: 10,
    fontWeight: 'bold',
    color: COLORS.gray800,
  },
  accountDetails: {
    marginBottom: theme.spacing.md,
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.sm,
  },
  detailLabel: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  detailValue: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.gray900,
  },
  maturityValue: {
    color: COLORS.green,
  },
  progressBar: {
    height: 6,
    backgroundColor: COLORS.gray200,
    borderRadius: 3,
    overflow: 'hidden',
    marginBottom: theme.spacing.sm,
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.green,
    borderRadius: 3,
  },
  maturityDate: {
    fontSize: 12,
    color: COLORS.gray500,
    textAlign: 'center',
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
    marginBottom: theme.spacing.lg,
  },
  emptyButton: {
    paddingHorizontal: theme.spacing.xl,
  },
});