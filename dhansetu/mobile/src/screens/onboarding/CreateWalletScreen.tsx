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
  Alert,
  ActivityIndicator,
  Clipboard,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import CheckBox from '@react-native-community/checkbox';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { useAppDispatch } from '@store/index';
import { createWallet } from '@store/slices/walletSlice';
import { HDWallet } from '@services/wallet/hdWallet';

export const CreateWalletScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  
  const [step, setStep] = useState<'info' | 'generate' | 'verify'>('info');
  const [mnemonic, setMnemonic] = useState<string[]>([]);
  const [isGenerating, setIsGenerating] = useState(false);
  const [showMnemonic, setShowMnemonic] = useState(false);
  const [hasBackedUp, setHasBackedUp] = useState(false);
  const [verificationWords, setVerificationWords] = useState<{ index: number; word: string }[]>([]);
  const [userInputs, setUserInputs] = useState<{ [key: number]: string }>({});
  const [isCreating, setIsCreating] = useState(false);

  const generateWallet = async () => {
    setIsGenerating(true);
    try {
      const hdWallet = await HDWallet.generate();
      const mnemonicPhrase = hdWallet.getMnemonic();
      setMnemonic(mnemonicPhrase.split(' '));
      setStep('generate');
      
      // Set up verification words (3 random words)
      const indices = [];
      while (indices.length < 3) {
        const randomIndex = Math.floor(Math.random() * 12);
        if (!indices.includes(randomIndex)) {
          indices.push(randomIndex);
        }
      }
      indices.sort((a, b) => a - b);
      setVerificationWords(indices.map(i => ({ index: i, word: mnemonicPhrase.split(' ')[i] })));
    } catch (error) {
      Alert.alert('Error', 'Failed to generate wallet. Please try again.');
    } finally {
      setIsGenerating(false);
    }
  };

  const copyMnemonic = () => {
    Clipboard.setString(mnemonic.join(' '));
    Alert.alert('Copied!', 'Recovery phrase copied to clipboard. Remember to delete it after saving!');
  };

  const verifyAndCreate = async () => {
    // Check if all verification words match
    const allCorrect = verificationWords.every(
      ({ index, word }) => userInputs[index]?.toLowerCase().trim() === word.toLowerCase()
    );

    if (!allCorrect) {
      Alert.alert('Incorrect', 'Please enter the correct words to verify you saved your recovery phrase.');
      return;
    }

    setIsCreating(true);
    try {
      await dispatch(createWallet(mnemonic.join(' '))).unwrap();
      navigation.navigate('PinSetup' as any, { isNewWallet: true });
    } catch (error) {
      Alert.alert('Error', 'Failed to create wallet. Please try again.');
    } finally {
      setIsCreating(false);
    }
  };

  const renderInfoStep = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.content}>
        <LinearGradient
          colors={theme.gradients.indianFlag}
          style={styles.iconContainer}
        >
          <Icon name="wallet-plus" size={64} color={COLORS.white} />
        </LinearGradient>

        <GradientText style={styles.title}>Create Your Wallet</GradientText>
        
        <Text style={styles.subtitle}>
          Your gateway to the DeshChain ecosystem
        </Text>

        <View style={styles.infoSection}>
          <Text style={styles.sectionTitle}>What is a wallet?</Text>
          <Text style={styles.infoText}>
            A crypto wallet is like your digital bank account. It stores your NAMO tokens and other digital assets securely.
          </Text>
        </View>

        <View style={styles.featuresContainer}>
          <View style={styles.featureItem}>
            <Icon name="shield-check" size={32} color={COLORS.green} />
            <Text style={styles.featureTitle}>Secure</Text>
            <Text style={styles.featureText}>Your keys, your coins</Text>
          </View>
          
          <View style={styles.featureItem}>
            <Icon name="key" size={32} color={COLORS.saffron} />
            <Text style={styles.featureTitle}>Private</Text>
            <Text style={styles.featureText}>Non-custodial wallet</Text>
          </View>
          
          <View style={styles.featureItem}>
            <Icon name="cellphone-lock" size={32} color={COLORS.navy} />
            <Text style={styles.featureTitle}>Protected</Text>
            <Text style={styles.featureText}>Biometric security</Text>
          </View>
        </View>

        <View style={styles.warningBox}>
          <Icon name="alert-circle" size={24} color={COLORS.warning} />
          <Text style={styles.warningText}>
            We will generate a 12-word recovery phrase. This is the ONLY way to recover your wallet if you lose access. Keep it safe!
          </Text>
        </View>

        <QuoteCard
          quote="सुरक्षा सर्वोपरि"
          translation="Security above all"
          language="hi"
          variant="minimal"
        />

        <CulturalButton
          title="Generate Wallet"
          onPress={generateWallet}
          size="large"
          style={styles.button}
          loading={isGenerating}
        />
      </View>
    </ScrollView>
  );

  const renderGenerateStep = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.content}>
        <Icon name="key-variant" size={48} color={COLORS.saffron} />
        
        <Text style={styles.title}>Your Recovery Phrase</Text>
        <Text style={styles.subtitle}>
          Write down these 12 words in order. This is your wallet backup.
        </Text>

        <TouchableOpacity
          style={styles.mnemonicContainer}
          onPress={() => setShowMnemonic(!showMnemonic)}
          activeOpacity={0.8}
        >
          {showMnemonic ? (
            <View style={styles.mnemonicGrid}>
              {mnemonic.map((word, index) => (
                <View key={index} style={styles.mnemonicItem}>
                  <Text style={styles.mnemonicNumber}>{index + 1}.</Text>
                  <Text style={styles.mnemonicWord}>{word}</Text>
                </View>
              ))}
            </View>
          ) : (
            <View style={styles.hiddenMnemonic}>
              <Icon name="eye-off" size={32} color={COLORS.gray600} />
              <Text style={styles.hiddenText}>Tap to reveal recovery phrase</Text>
            </View>
          )}
        </TouchableOpacity>

        {showMnemonic && (
          <TouchableOpacity style={styles.copyButton} onPress={copyMnemonic}>
            <Icon name="content-copy" size={20} color={COLORS.saffron} />
            <Text style={styles.copyText}>Copy to clipboard</Text>
          </TouchableOpacity>
        )}

        <View style={styles.warningBox}>
          <Icon name="alert" size={24} color={COLORS.error} />
          <Text style={styles.warningText}>
            Never share your recovery phrase with anyone! DeshChain support will NEVER ask for it.
          </Text>
        </View>

        <View style={styles.checkboxContainer}>
          <CheckBox
            value={hasBackedUp}
            onValueChange={setHasBackedUp}
            tintColors={{ true: COLORS.saffron, false: COLORS.gray400 }}
          />
          <Text style={styles.checkboxText}>
            I have written down my recovery phrase in a safe place
          </Text>
        </View>

        <CulturalButton
          title="Continue to Verify"
          onPress={() => setStep('verify')}
          size="large"
          style={styles.button}
          disabled={!hasBackedUp || !showMnemonic}
        />
      </View>
    </ScrollView>
  );

  const renderVerifyStep = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.content}>
        <Icon name="shield-check" size={48} color={COLORS.green} />
        
        <Text style={styles.title}>Verify Your Backup</Text>
        <Text style={styles.subtitle}>
          Enter the requested words to confirm you saved your recovery phrase
        </Text>

        <View style={styles.verificationContainer}>
          {verificationWords.map(({ index, word }) => (
            <View key={index} style={styles.verificationItem}>
              <Text style={styles.verificationLabel}>Word #{index + 1}</Text>
              <TextInput
                style={styles.verificationInput}
                placeholder="Enter word"
                value={userInputs[index] || ''}
                onChangeText={(text) => setUserInputs({ ...userInputs, [index]: text })}
                autoCapitalize="none"
                autoCorrect={false}
              />
            </View>
          ))}
        </View>

        <CulturalButton
          title="Create Wallet"
          onPress={verifyAndCreate}
          size="large"
          style={styles.button}
          loading={isCreating}
          disabled={Object.keys(userInputs).length < 3}
        />

        <TouchableOpacity
          style={styles.backButton}
          onPress={() => setStep('generate')}
        >
          <Text style={styles.backText}>Back to recovery phrase</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Icon name="arrow-left" size={24} color={COLORS.gray900} />
        </TouchableOpacity>
        <View style={styles.progressBar}>
          <View 
            style={[
              styles.progressFill,
              { width: step === 'info' ? '33%' : step === 'generate' ? '66%' : '100%' }
            ]}
          />
        </View>
        <View style={{ width: 24 }} />
      </View>

      {/* Content */}
      {step === 'info' && renderInfoStep()}
      {step === 'generate' && renderGenerateStep()}
      {step === 'verify' && renderVerifyStep()}
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
    gap: theme.spacing.md,
  },
  progressBar: {
    flex: 1,
    height: 4,
    backgroundColor: COLORS.gray200,
    borderRadius: 2,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.saffron,
    borderRadius: 2,
  },
  content: {
    padding: theme.spacing.lg,
    alignItems: 'center',
  },
  iconContainer: {
    width: 120,
    height: 120,
    borderRadius: 60,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.xl,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.gray900,
    textAlign: 'center',
    marginBottom: theme.spacing.sm,
  },
  subtitle: {
    fontSize: 16,
    color: COLORS.gray600,
    textAlign: 'center',
    marginBottom: theme.spacing.xl,
    paddingHorizontal: theme.spacing.lg,
  },
  infoSection: {
    width: '100%',
    marginBottom: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.sm,
  },
  infoText: {
    fontSize: 14,
    color: COLORS.gray700,
    lineHeight: 22,
  },
  featuresContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    width: '100%',
    marginBottom: theme.spacing.xl,
  },
  featureItem: {
    alignItems: 'center',
    flex: 1,
  },
  featureTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginTop: theme.spacing.sm,
  },
  featureText: {
    fontSize: 12,
    color: COLORS.gray600,
    marginTop: 4,
  },
  warningBox: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    backgroundColor: COLORS.warning + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.lg,
    width: '100%',
  },
  warningText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.warning,
    marginLeft: theme.spacing.sm,
    lineHeight: 20,
  },
  button: {
    marginTop: theme.spacing.lg,
    width: '100%',
  },
  mnemonicContainer: {
    width: '100%',
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    marginBottom: theme.spacing.md,
  },
  mnemonicGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
  },
  mnemonicItem: {
    flexDirection: 'row',
    alignItems: 'center',
    width: '48%',
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.sm,
    padding: theme.spacing.sm,
    marginBottom: theme.spacing.sm,
  },
  mnemonicNumber: {
    fontSize: 12,
    color: COLORS.gray500,
    marginRight: theme.spacing.sm,
    width: 20,
  },
  mnemonicWord: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  hiddenMnemonic: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
  },
  hiddenText: {
    fontSize: 16,
    color: COLORS.gray600,
    marginTop: theme.spacing.md,
  },
  copyButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  copyText: {
    fontSize: 14,
    color: COLORS.saffron,
    marginLeft: theme.spacing.sm,
    fontWeight: '600',
  },
  checkboxContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    width: '100%',
    marginBottom: theme.spacing.lg,
  },
  checkboxText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.gray700,
    marginLeft: theme.spacing.sm,
  },
  verificationContainer: {
    width: '100%',
    marginBottom: theme.spacing.xl,
  },
  verificationItem: {
    marginBottom: theme.spacing.lg,
  },
  verificationLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.gray700,
    marginBottom: theme.spacing.sm,
  },
  verificationInput: {
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    fontSize: 16,
    color: COLORS.gray900,
  },
  backButton: {
    marginTop: theme.spacing.md,
    padding: theme.spacing.md,
  },
  backText: {
    fontSize: 14,
    color: COLORS.gray600,
    textDecorationLine: 'underline',
  },
});