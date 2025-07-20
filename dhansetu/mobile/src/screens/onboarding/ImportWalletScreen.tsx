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

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  Alert,
  KeyboardAvoidingView,
  Platform,
  Clipboard,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { QuoteCard } from '@components/common/QuoteCard';
import { useAppDispatch } from '@store/index';
import { importWallet } from '@store/slices/walletSlice';
import { HDWallet } from '@services/wallet/hdWallet';

export const ImportWalletScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  
  const [importMethod, setImportMethod] = useState<'phrase' | 'key'>('phrase');
  const [mnemonicWords, setMnemonicWords] = useState<string[]>(Array(12).fill(''));
  const [privateKey, setPrivateKey] = useState('');
  const [isImporting, setIsImporting] = useState(false);
  const [showPrivateKey, setShowPrivateKey] = useState(false);
  const [isValidating, setIsValidating] = useState(false);
  const [validationErrors, setValidationErrors] = useState<number[]>([]);

  const handleWordChange = (index: number, word: string) => {
    const newWords = [...mnemonicWords];
    newWords[index] = word.trim().toLowerCase();
    setMnemonicWords(newWords);
    
    // Clear validation error for this word
    setValidationErrors(errors => errors.filter(i => i !== index));
  };

  const pasteFromClipboard = async () => {
    try {
      const clipboardContent = await Clipboard.getString();
      const words = clipboardContent.trim().split(/\s+/);
      
      if (words.length === 12) {
        setMnemonicWords(words.map(w => w.toLowerCase()));
      } else if (words.length === 24) {
        Alert.alert(
          '24-word phrase detected',
          'DhanSetu supports 12-word phrases. Please use a compatible wallet.',
          [{ text: 'OK' }]
        );
      } else {
        Alert.alert('Invalid', 'Clipboard does not contain a valid recovery phrase');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to read from clipboard');
    }
  };

  const validateMnemonic = async () => {
    setIsValidating(true);
    const errors: number[] = [];
    
    // Check each word
    mnemonicWords.forEach((word, index) => {
      if (!word || word.length < 3) {
        errors.push(index);
      }
    });
    
    if (errors.length > 0) {
      setValidationErrors(errors);
      setIsValidating(false);
      return false;
    }
    
    // Validate mnemonic phrase
    try {
      const hdWallet = new HDWallet();
      const isValid = await hdWallet.validateMnemonic(mnemonicWords.join(' '));
      setIsValidating(false);
      return isValid;
    } catch (error) {
      setIsValidating(false);
      return false;
    }
  };

  const handleImport = async () => {
    if (importMethod === 'phrase') {
      const isValid = await validateMnemonic();
      if (!isValid) {
        Alert.alert('Invalid', 'Please check your recovery phrase and try again');
        return;
      }
    } else if (!privateKey || privateKey.length < 64) {
      Alert.alert('Invalid', 'Please enter a valid private key');
      return;
    }

    setIsImporting(true);
    try {
      if (importMethod === 'phrase') {
        await dispatch(importWallet({
          type: 'mnemonic',
          value: mnemonicWords.join(' '),
        })).unwrap();
      } else {
        await dispatch(importWallet({
          type: 'privateKey',
          value: privateKey,
        })).unwrap();
      }
      
      navigation.navigate('PinSetup' as any, { isNewWallet: false });
    } catch (error: any) {
      Alert.alert('Import Failed', error.message || 'Please check your input and try again');
    } finally {
      setIsImporting(false);
    }
  };

  const renderPhraseImport = () => (
    <View style={styles.importContainer}>
      <View style={styles.importHeader}>
        <Text style={styles.importTitle}>Enter Recovery Phrase</Text>
        <TouchableOpacity onPress={pasteFromClipboard}>
          <Icon name="content-paste" size={24} color={COLORS.saffron} />
        </TouchableOpacity>
      </View>
      
      <Text style={styles.importSubtitle}>
        Enter your 12-word recovery phrase in the correct order
      </Text>

      <View style={styles.wordsGrid}>
        {mnemonicWords.map((word, index) => (
          <View
            key={index}
            style={[
              styles.wordInputContainer,
              validationErrors.includes(index) && styles.wordInputError,
            ]}
          >
            <Text style={styles.wordNumber}>{index + 1}.</Text>
            <TextInput
              style={styles.wordInput}
              placeholder="word"
              value={word}
              onChangeText={(text) => handleWordChange(index, text)}
              autoCapitalize="none"
              autoCorrect={false}
              returnKeyType={index === 11 ? 'done' : 'next'}
            />
          </View>
        ))}
      </View>

      <View style={styles.infoBox}>
        <Icon name="information" size={20} color={COLORS.info} />
        <Text style={styles.infoText}>
          Recovery phrases are typically 12 words. Make sure you enter them in the exact order.
        </Text>
      </View>
    </View>
  );

  const renderKeyImport = () => (
    <View style={styles.importContainer}>
      <Text style={styles.importTitle}>Enter Private Key</Text>
      <Text style={styles.importSubtitle}>
        Enter your wallet's private key (64 characters)
      </Text>

      <View style={styles.keyInputContainer}>
        <TextInput
          style={styles.keyInput}
          placeholder="Enter private key"
          value={privateKey}
          onChangeText={setPrivateKey}
          autoCapitalize="none"
          autoCorrect={false}
          secureTextEntry={!showPrivateKey}
          multiline={showPrivateKey}
          numberOfLines={showPrivateKey ? 3 : 1}
        />
        <TouchableOpacity
          style={styles.eyeButton}
          onPress={() => setShowPrivateKey(!showPrivateKey)}
        >
          <Icon
            name={showPrivateKey ? 'eye-off' : 'eye'}
            size={24}
            color={COLORS.gray600}
          />
        </TouchableOpacity>
      </View>

      <View style={styles.warningBox}>
        <Icon name="alert" size={20} color={COLORS.error} />
        <Text style={styles.warningText}>
          Private keys provide full access to your wallet. Only import from trusted sources.
        </Text>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardAvoid}
      >
        <ScrollView showsVerticalScrollIndicator={false}>
          {/* Header */}
          <View style={styles.header}>
            <TouchableOpacity onPress={() => navigation.goBack()}>
              <Icon name="arrow-left" size={24} color={COLORS.gray900} />
            </TouchableOpacity>
            <Text style={styles.headerTitle}>Import Wallet</Text>
            <View style={{ width: 24 }} />
          </View>

          <View style={styles.content}>
            <LinearGradient
              colors={[COLORS.saffron, COLORS.orange]}
              style={styles.iconContainer}
            >
              <Icon name="import" size={64} color={COLORS.white} />
            </LinearGradient>

            <GradientText style={styles.title}>Import Existing Wallet</GradientText>
            <Text style={styles.subtitle}>
              Restore your wallet using recovery phrase or private key
            </Text>

            {/* Import Method Toggle */}
            <View style={styles.methodToggle}>
              <TouchableOpacity
                style={[
                  styles.methodButton,
                  importMethod === 'phrase' && styles.methodButtonActive,
                ]}
                onPress={() => setImportMethod('phrase')}
              >
                <Icon
                  name="format-list-numbered"
                  size={20}
                  color={importMethod === 'phrase' ? COLORS.white : COLORS.gray600}
                />
                <Text
                  style={[
                    styles.methodText,
                    importMethod === 'phrase' && styles.methodTextActive,
                  ]}
                >
                  Recovery Phrase
                </Text>
              </TouchableOpacity>
              
              <TouchableOpacity
                style={[
                  styles.methodButton,
                  importMethod === 'key' && styles.methodButtonActive,
                ]}
                onPress={() => setImportMethod('key')}
              >
                <Icon
                  name="key"
                  size={20}
                  color={importMethod === 'key' ? COLORS.white : COLORS.gray600}
                />
                <Text
                  style={[
                    styles.methodText,
                    importMethod === 'key' && styles.methodTextActive,
                  ]}
                >
                  Private Key
                </Text>
              </TouchableOpacity>
            </View>

            {/* Import Forms */}
            {importMethod === 'phrase' ? renderPhraseImport() : renderKeyImport()}

            {/* Cultural Quote */}
            <QuoteCard
              quote="पुरानी चाबी से नया दरवाज़ा"
              translation="Old key opens new doors"
              language="hi"
              variant="minimal"
            />

            {/* Import Button */}
            <CulturalButton
              title="Import Wallet"
              onPress={handleImport}
              size="large"
              style={styles.importButton}
              loading={isImporting || isValidating}
              disabled={
                importMethod === 'phrase'
                  ? mnemonicWords.some(w => !w)
                  : !privateKey
              }
            />
          </View>
        </ScrollView>
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
    textAlign: 'center',
    marginBottom: theme.spacing.sm,
  },
  subtitle: {
    fontSize: 16,
    color: COLORS.gray600,
    textAlign: 'center',
    marginBottom: theme.spacing.xl,
  },
  methodToggle: {
    flexDirection: 'row',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: 4,
    marginBottom: theme.spacing.xl,
    width: '100%',
  },
  methodButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: theme.spacing.md,
    borderRadius: theme.borderRadius.sm,
    gap: theme.spacing.sm,
  },
  methodButtonActive: {
    backgroundColor: COLORS.saffron,
  },
  methodText: {
    fontSize: 14,
    color: COLORS.gray600,
    fontWeight: '500',
  },
  methodTextActive: {
    color: COLORS.white,
    fontWeight: 'bold',
  },
  importContainer: {
    width: '100%',
    marginBottom: theme.spacing.lg,
  },
  importHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  importTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  importSubtitle: {
    fontSize: 14,
    color: COLORS.gray600,
    marginBottom: theme.spacing.lg,
  },
  wordsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.lg,
  },
  wordInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    width: '48%',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.sm,
    paddingHorizontal: theme.spacing.sm,
    marginBottom: theme.spacing.sm,
    height: 40,
  },
  wordInputError: {
    borderWidth: 1,
    borderColor: COLORS.error,
  },
  wordNumber: {
    fontSize: 12,
    color: COLORS.gray500,
    marginRight: theme.spacing.xs,
    width: 20,
  },
  wordInput: {
    flex: 1,
    fontSize: 14,
    color: COLORS.gray900,
  },
  keyInputContainer: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    paddingRight: theme.spacing.sm,
    marginBottom: theme.spacing.lg,
  },
  keyInput: {
    flex: 1,
    fontSize: 14,
    color: COLORS.gray900,
    padding: theme.spacing.md,
    minHeight: 80,
  },
  eyeButton: {
    padding: theme.spacing.md,
  },
  infoBox: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    backgroundColor: COLORS.info + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
  },
  infoText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.info,
    marginLeft: theme.spacing.sm,
    lineHeight: 20,
  },
  warningBox: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    backgroundColor: COLORS.error + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
  },
  warningText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.error,
    marginLeft: theme.spacing.sm,
    lineHeight: 20,
  },
  importButton: {
    marginTop: theme.spacing.lg,
    width: '100%',
  },
});