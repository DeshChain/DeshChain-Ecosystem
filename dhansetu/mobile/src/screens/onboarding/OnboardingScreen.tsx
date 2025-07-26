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

import React, { useState, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Dimensions,
  TouchableOpacity,
  Image,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import Swiper from 'react-native-swiper';
import Animated, {
  useAnimatedStyle,
  withSpring,
  interpolate,
} from 'react-native-reanimated';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';

const { width, height } = Dimensions.get('window');

interface OnboardingSlide {
  id: string;
  title: string;
  subtitle: string;
  description: string;
  icon: string;
  gradient: string[];
  culturalQuote?: {
    text: string;
    translation: string;
  };
}

export const OnboardingScreen: React.FC = () => {
  const navigation = useNavigation();
  const swiperRef = useRef<Swiper>(null);
  const [currentIndex, setCurrentIndex] = useState(0);

  const slides: OnboardingSlide[] = [
    {
      id: 'welcome',
      title: 'Welcome to DhanSetu',
      subtitle: 'Digital money transfer for the modern age',
      description: 'Your gateway to India\'s blockchain revolution. Send money, invest in memecoins, and secure your future with guaranteed pension.',
      icon: 'hand-wave',
      gradient: theme.gradients.indianFlag,
      culturalQuote: {
        text: 'वसुधैव कुटुम्बकम्',
        translation: 'The world is one family',
      },
    },
    {
      id: 'dhanpata',
      title: 'DhanPata Virtual IDs',
      subtitle: 'Simple as username@dhan',
      description: 'No more copying long blockchain addresses! Share your DhanPata ID like ramesh@dhan and receive money instantly.',
      icon: 'at',
      gradient: [COLORS.saffron, COLORS.orange],
      culturalQuote: {
        text: 'सरलता में ही सुंदरता है',
        translation: 'Beauty lies in simplicity',
      },
    },
    {
      id: 'moneyorder',
      title: 'Money Order DEX',
      subtitle: 'P2P transfers reimagined',
      description: 'Send money to anyone, anywhere in India. PIN code based routing ensures your money reaches the right place.',
      icon: 'cash-fast',
      gradient: [COLORS.green, COLORS.darkGreen],
      culturalQuote: {
        text: 'दूरियां मिट जाती हैं',
        translation: 'Distances disappear',
      },
    },
    {
      id: 'sikkebaaz',
      title: 'Sikkebaaz Launchpad',
      subtitle: 'Desi memecoins for everyone',
      description: 'Launch your own community token with anti-pump protection. From Bollywood to Cricket, celebrate what you love!',
      icon: 'rocket-launch',
      gradient: [COLORS.festivalPrimary, COLORS.festivalSecondary],
      culturalQuote: {
        text: 'हर सिक्के की अपनी कहानी',
        translation: 'Every coin has its own story',
      },
    },
    {
      id: 'suraksha',
      title: 'Gram Suraksha Pool',
      subtitle: 'Minimum 8% guaranteed, up to 50% returns',
      description: 'India\'s first blockchain pension with guaranteed returns. Start with just ₹100/month and secure your future.',
      icon: 'shield-check',
      gradient: [COLORS.navy, COLORS.sky],
      culturalQuote: {
        text: 'बूंद बूंद से सागर बनता है',
        translation: 'Drop by drop, an ocean is formed',
      },
    },
    {
      id: 'security',
      title: 'Your Keys, Your Coins',
      subtitle: 'Bank-grade security',
      description: 'Non-custodial wallet with biometric protection. Your private keys never leave your device.',
      icon: 'lock',
      gradient: [COLORS.gray700, COLORS.gray900],
      culturalQuote: {
        text: 'सुरक्षा सर्वोपरि',
        translation: 'Security above all',
      },
    },
  ];

  const handleNext = () => {
    if (currentIndex < slides.length - 1) {
      swiperRef.current?.scrollBy(1);
    } else {
      navigation.navigate('CreateWallet' as any);
    }
  };

  const handleSkip = () => {
    navigation.navigate('CreateWallet' as any);
  };

  const renderSlide = (slide: OnboardingSlide, index: number) => {
    const isActive = index === currentIndex;
    
    return (
      <View key={slide.id} style={styles.slide}>
        <LinearGradient
          colors={slide.gradient}
          style={styles.slideGradient}
        >
          <Animated.View
            style={[
              styles.iconContainer,
              useAnimatedStyle(() => ({
                transform: [
                  {
                    scale: withSpring(isActive ? 1 : 0.8),
                  },
                ],
              })),
            ]}
          >
            <Icon name={slide.icon} size={80} color={COLORS.white} />
          </Animated.View>
        </LinearGradient>

        <View style={styles.contentContainer}>
          <GradientText style={styles.title}>{slide.title}</GradientText>
          <Text style={styles.subtitle}>{slide.subtitle}</Text>
          <Text style={styles.description}>{slide.description}</Text>
          
          {slide.culturalQuote && (
            <View style={styles.quoteContainer}>
              <Text style={styles.quoteText}>"{slide.culturalQuote.text}"</Text>
              <Text style={styles.quoteTranslation}>{slide.culturalQuote.translation}</Text>
            </View>
          )}
        </View>
      </View>
    );
  };

  const renderPagination = () => (
    <View style={styles.pagination}>
      {slides.map((_, index) => (
        <View
          key={index}
          style={[
            styles.paginationDot,
            index === currentIndex && styles.paginationDotActive,
          ]}
        />
      ))}
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* Skip Button */}
      {currentIndex < slides.length - 1 && (
        <TouchableOpacity style={styles.skipButton} onPress={handleSkip}>
          <Text style={styles.skipText}>Skip</Text>
        </TouchableOpacity>
      )}

      {/* Swiper */}
      <Swiper
        ref={swiperRef}
        loop={false}
        showsPagination={false}
        onIndexChanged={setCurrentIndex}
        style={styles.swiper}
      >
        {slides.map((slide, index) => renderSlide(slide, index))}
      </Swiper>

      {/* Bottom Section */}
      <View style={styles.bottomSection}>
        {renderPagination()}
        
        <CulturalButton
          title={currentIndex === slides.length - 1 ? 'Get Started' : 'Next'}
          onPress={handleNext}
          size="large"
          style={styles.nextButton}
        />
        
        {currentIndex === slides.length - 1 && (
          <TouchableOpacity
            style={styles.importButton}
            onPress={() => navigation.navigate('ImportWallet' as any)}
          >
            <Text style={styles.importText}>Already have a wallet? Import</Text>
          </TouchableOpacity>
        )}
      </View>

      {/* Background Pattern */}
      <View style={styles.backgroundPattern}>
        {[...Array(20)].map((_, i) => (
          <View
            key={i}
            style={[
              styles.patternDot,
              {
                left: Math.random() * width,
                top: Math.random() * height,
                opacity: Math.random() * 0.1,
              },
            ]}
          />
        ))}
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  skipButton: {
    position: 'absolute',
    top: theme.spacing.xl,
    right: theme.spacing.lg,
    zIndex: 10,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  skipText: {
    fontSize: 16,
    color: COLORS.gray600,
    fontWeight: '600',
  },
  swiper: {
    flex: 1,
  },
  slide: {
    flex: 1,
  },
  slideGradient: {
    height: height * 0.4,
    justifyContent: 'center',
    alignItems: 'center',
    borderBottomLeftRadius: 50,
    borderBottomRightRadius: 50,
  },
  iconContainer: {
    width: 120,
    height: 120,
    borderRadius: 60,
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  contentContainer: {
    flex: 1,
    paddingHorizontal: theme.spacing.xl,
    paddingTop: theme.spacing.xl,
    alignItems: 'center',
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: theme.spacing.md,
  },
  subtitle: {
    fontSize: 18,
    color: COLORS.gray700,
    textAlign: 'center',
    marginBottom: theme.spacing.lg,
  },
  description: {
    fontSize: 16,
    color: COLORS.gray600,
    textAlign: 'center',
    lineHeight: 24,
    paddingHorizontal: theme.spacing.md,
  },
  quoteContainer: {
    marginTop: theme.spacing.xl,
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    width: '100%',
  },
  quoteText: {
    fontSize: 18,
    color: COLORS.gray800,
    fontStyle: 'italic',
    textAlign: 'center',
    marginBottom: theme.spacing.sm,
  },
  quoteTranslation: {
    fontSize: 14,
    color: COLORS.gray600,
    textAlign: 'center',
  },
  bottomSection: {
    paddingHorizontal: theme.spacing.xl,
    paddingBottom: theme.spacing.xl,
  },
  pagination: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginBottom: theme.spacing.lg,
  },
  paginationDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: COLORS.gray300,
    marginHorizontal: 4,
  },
  paginationDotActive: {
    width: 24,
    backgroundColor: COLORS.saffron,
  },
  nextButton: {
    marginBottom: theme.spacing.md,
  },
  importButton: {
    alignItems: 'center',
    paddingVertical: theme.spacing.md,
  },
  importText: {
    fontSize: 16,
    color: COLORS.saffron,
    fontWeight: '600',
  },
  backgroundPattern: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    pointerEvents: 'none',
  },
  patternDot: {
    position: 'absolute',
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.saffron,
  },
});