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
  TouchableOpacity,
  Vibration,
  Alert,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation, useRoute } from '@react-navigation/native';
import TouchID from 'react-native-touch-id';
import Animated, {
  useAnimatedStyle,
  withSequence,
  withTiming,
  withSpring,
} from 'react-native-reanimated';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppDispatch, useAppSelector } from '@store/index';
import { setPin, authenticate } from '@store/slices/authSlice';
import { SecureStorage } from '@services/storage/secureStorage';

interface RouteParams {
  isNewWallet?: boolean;
  isChangingPin?: boolean;
}

export const PinSetupScreen: React.FC = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const dispatch = useAppDispatch();
  const { hasWallet } = useAppSelector((state) => state.wallet);
  
  const params = route.params as RouteParams;
  const isNewWallet = params?.isNewWallet || false;
  const isChangingPin = params?.isChangingPin || false;
  const isUnlocking = !isNewWallet && !isChangingPin;
  
  const [pin, setPin] = useState('');
  const [confirmPin, setConfirmPin] = useState('');
  const [isConfirming, setIsConfirming] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const [isLocked, setIsLocked] = useState(false);
  const [lockTime, setLockTime] = useState(0);
  const [showBiometric, setShowBiometric] = useState(false);
  const [shake, setShake] = useState(false);

  useEffect(() => {
    if (isUnlocking) {
      checkBiometricAvailability();
    }
  }, []);

  useEffect(() => {
    if (lockTime > 0) {
      const timer = setTimeout(() => {
        setLockTime(lockTime - 1);
        if (lockTime === 1) {
          setIsLocked(false);
          setAttempts(0);
        }
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [lockTime]);

  const checkBiometricAvailability = async () => {
    try {
      const biometryType = await TouchID.isSupported();
      if (biometryType) {
        setShowBiometric(true);
        // Auto-prompt for biometric on unlock
        setTimeout(() => authenticateWithBiometric(), 500);
      }
    } catch (error) {
      console.log('Biometric not available');
    }
  };

  const authenticateWithBiometric = async () => {
    try {
      const success = await TouchID.authenticate('Unlock your DhanSetu wallet', {
        title: 'Biometric Authentication',
        cancelText: 'Use PIN',
        fallbackLabel: 'Use PIN',
      });
      
      if (success) {
        dispatch(authenticate());
        navigation.reset({
          index: 0,
          routes: [{ name: 'MainTabs' as any }],
        });
      }
    } catch (error) {
      console.log('Biometric authentication failed');
    }
  };

  const handleNumberPress = (num: string) => {
    if (isLocked) return;
    
    if (!isConfirming) {
      if (pin.length < 6) {
        setPin(pin + num);
        
        if (isUnlocking && pin.length === 5) {
          // Auto-submit on 6th digit for unlocking
          setTimeout(() => verifyPin(pin + num), 100);
        } else if (!isUnlocking && pin.length === 5) {
          // Move to confirm step for new PIN
          setTimeout(() => {
            setIsConfirming(true);
          }, 100);
        }
      }
    } else {
      if (confirmPin.length < 6) {
        setConfirmPin(confirmPin + num);
        
        if (confirmPin.length === 5) {
          // Auto-submit on 6th digit
          setTimeout(() => handlePinSetup(confirmPin + num), 100);
        }
      }
    }
  };

  const handleDelete = () => {
    if (isLocked) return;
    
    if (!isConfirming) {
      setPin(pin.slice(0, -1));
    } else {
      setConfirmPin(confirmPin.slice(0, -1));
    }
  };

  const verifyPin = async (enteredPin: string) => {
    try {
      const secureStorage = new SecureStorage();
      const storedPin = await secureStorage.getPin();
      
      if (storedPin === enteredPin) {
        dispatch(authenticate());
        navigation.reset({
          index: 0,
          routes: [{ name: 'MainTabs' as any }],
        });
      } else {
        handleWrongPin();
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to verify PIN');
    }
  };

  const handleWrongPin = () => {
    Vibration.vibrate(200);
    setShake(true);
    setTimeout(() => setShake(false), 500);
    
    const newAttempts = attempts + 1;
    setAttempts(newAttempts);
    setPin('');
    
    if (newAttempts >= 3) {
      setIsLocked(true);
      setLockTime(30); // Lock for 30 seconds
      Alert.alert(
        'Too Many Attempts',
        'Please wait 30 seconds before trying again',
        [{ text: 'OK' }]
      );
    }
  };

  const handlePinSetup = async (confirmedPin: string) => {
    if (pin !== confirmedPin) {
      Vibration.vibrate(200);
      Alert.alert('PIN Mismatch', 'PINs do not match. Please try again.');
      setPin('');
      setConfirmPin('');
      setIsConfirming(false);
      return;
    }
    
    try {
      await dispatch(setPin(pin)).unwrap();
      
      if (isNewWallet) {
        Alert.alert(
          'Success!',
          'Your wallet is ready. Welcome to DhanSetu!',
          [{
            text: 'Let\'s Go!',
            onPress: () => {
              dispatch(authenticate());
              navigation.reset({
                index: 0,
                routes: [{ name: 'MainTabs' as any }],
              });
            },
          }]
        );
      } else if (isChangingPin) {
        Alert.alert('Success', 'PIN changed successfully');
        navigation.goBack();
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to set PIN');
    }
  };

  const renderPinDots = () => {
    const pinToShow = isConfirming ? confirmPin : pin;
    return (
      <Animated.View
        style={[
          styles.pinDotsContainer,
          useAnimatedStyle(() => ({
            transform: [
              {
                translateX: shake
                  ? withSequence(
                      withTiming(-10, { duration: 50 }),
                      withTiming(10, { duration: 50 }),
                      withTiming(-10, { duration: 50 }),
                      withTiming(0, { duration: 50 })
                    )
                  : 0,
              },
            ],
          })),
        ]}
      >
        {[...Array(6)].map((_, index) => (
          <View
            key={index}
            style={[
              styles.pinDot,
              pinToShow.length > index && styles.pinDotFilled,
            ]}
          />
        ))}
      </Animated.View>
    );
  };

  const renderNumberPad = () => {
    const numbers = [
      ['1', '2', '3'],
      ['4', '5', '6'],
      ['7', '8', '9'],
      ['biometric', '0', 'delete'],
    ];

    return (
      <View style={styles.numberPad}>
        {numbers.map((row, rowIndex) => (
          <View key={rowIndex} style={styles.numberRow}>
            {row.map((item) => {
              if (item === 'biometric') {
                return showBiometric && isUnlocking ? (
                  <TouchableOpacity
                    key={item}
                    style={styles.numberButton}
                    onPress={authenticateWithBiometric}
                    disabled={isLocked}
                  >
                    <Icon name="fingerprint" size={28} color={COLORS.saffron} />
                  </TouchableOpacity>
                ) : (
                  <View key={item} style={styles.numberButton} />
                );
              } else if (item === 'delete') {
                return (
                  <TouchableOpacity
                    key={item}
                    style={styles.numberButton}
                    onPress={handleDelete}
                    disabled={isLocked}
                  >
                    <Icon name="backspace" size={24} color={COLORS.gray700} />
                  </TouchableOpacity>
                );
              } else {
                return (
                  <TouchableOpacity
                    key={item}
                    style={[styles.numberButton, isLocked && styles.numberButtonDisabled]}
                    onPress={() => handleNumberPress(item)}
                    disabled={isLocked}
                  >
                    <Text style={[styles.numberText, isLocked && styles.numberTextDisabled]}>
                      {item}
                    </Text>
                  </TouchableOpacity>
                );
              }
            })}
          </View>
        ))}
      </View>
    );
  };

  const getTitle = () => {
    if (isUnlocking) return 'Enter PIN';
    if (isChangingPin) return isConfirming ? 'Confirm New PIN' : 'Enter New PIN';
    return isConfirming ? 'Confirm Your PIN' : 'Create a PIN';
  };

  const getSubtitle = () => {
    if (isUnlocking) return 'Enter your 6-digit PIN to unlock';
    if (isConfirming) return 'Re-enter your PIN to confirm';
    return 'Choose a 6-digit PIN to secure your wallet';
  };

  return (
    <SafeAreaView style={styles.container}>
      <LinearGradient
        colors={theme.gradients.subtle}
        style={styles.gradient}
      >
        {/* Header */}
        {!isUnlocking && (
          <View style={styles.header}>
            <TouchableOpacity onPress={() => navigation.goBack()}>
              <Icon name="arrow-left" size={24} color={COLORS.gray900} />
            </TouchableOpacity>
          </View>
        )}

        {/* Content */}
        <View style={styles.content}>
          <View style={styles.iconContainer}>
            <Icon name="lock" size={48} color={COLORS.saffron} />
          </View>

          <GradientText style={styles.title}>{getTitle()}</GradientText>
          <Text style={styles.subtitle}>{getSubtitle()}</Text>

          {isLocked && (
            <View style={styles.lockMessage}>
              <Icon name="timer" size={20} color={COLORS.error} />
              <Text style={styles.lockText}>
                Too many attempts. Try again in {lockTime} seconds
              </Text>
            </View>
          )}

          {attempts > 0 && !isLocked && (
            <Text style={styles.attemptsText}>
              {3 - attempts} attempts remaining
            </Text>
          )}

          {renderPinDots()}
        </View>

        {/* Number Pad */}
        {renderNumberPad()}

        {/* Skip Option for New Wallet (Not Recommended) */}
        {isNewWallet && !isConfirming && (
          <TouchableOpacity
            style={styles.skipButton}
            onPress={() => {
              Alert.alert(
                'Skip PIN Setup?',
                'Your wallet will not be protected. This is not recommended.',
                [
                  { text: 'Cancel', style: 'cancel' },
                  {
                    text: 'Skip',
                    style: 'destructive',
                    onPress: () => {
                      dispatch(authenticate());
                      navigation.reset({
                        index: 0,
                        routes: [{ name: 'MainTabs' as any }],
                      });
                    },
                  },
                ]
              );
            }}
          >
            <Text style={styles.skipText}>Skip for now</Text>
          </TouchableOpacity>
        )}
      </LinearGradient>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  gradient: {
    flex: 1,
  },
  header: {
    paddingHorizontal: theme.spacing.lg,
    paddingTop: theme.spacing.md,
  },
  content: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: theme.spacing.xl,
  },
  iconContainer: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: COLORS.saffron + '20',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.xl,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    marginBottom: theme.spacing.sm,
  },
  subtitle: {
    fontSize: 16,
    color: COLORS.gray600,
    textAlign: 'center',
    marginBottom: theme.spacing.xl,
  },
  lockMessage: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.error + '10',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.lg,
  },
  lockText: {
    fontSize: 14,
    color: COLORS.error,
    marginLeft: theme.spacing.sm,
  },
  attemptsText: {
    fontSize: 14,
    color: COLORS.warning,
    marginBottom: theme.spacing.md,
  },
  pinDotsContainer: {
    flexDirection: 'row',
    gap: theme.spacing.md,
    marginBottom: theme.spacing.xl,
  },
  pinDot: {
    width: 16,
    height: 16,
    borderRadius: 8,
    borderWidth: 2,
    borderColor: COLORS.gray400,
    backgroundColor: 'transparent',
  },
  pinDotFilled: {
    backgroundColor: COLORS.saffron,
    borderColor: COLORS.saffron,
  },
  numberPad: {
    paddingHorizontal: theme.spacing.xl,
    paddingBottom: theme.spacing.xl,
  },
  numberRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.md,
  },
  numberButton: {
    width: 75,
    height: 75,
    borderRadius: 37.5,
    backgroundColor: COLORS.white,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.small,
  },
  numberButtonDisabled: {
    backgroundColor: COLORS.gray100,
  },
  numberText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  numberTextDisabled: {
    color: COLORS.gray400,
  },
  skipButton: {
    alignItems: 'center',
    paddingBottom: theme.spacing.xl,
  },
  skipText: {
    fontSize: 14,
    color: COLORS.gray600,
    textDecorationLine: 'underline',
  },
});