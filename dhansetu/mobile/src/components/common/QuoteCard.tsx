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

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import { COLORS, theme } from '@constants/theme';
import { GradientText } from './GradientText';

interface QuoteCardProps {
  quote: string;
  author?: string;
  translation?: string;
  language?: 'en' | 'hi' | 'sa';
  variant?: 'default' | 'festival' | 'minimal';
  showQuoteMarks?: boolean;
  style?: any;
}

export const QuoteCard: React.FC<QuoteCardProps> = ({
  quote,
  author,
  translation,
  language = 'en',
  variant = 'default',
  showQuoteMarks = true,
  style,
}) => {
  const getGradientColors = () => {
    switch (variant) {
      case 'festival':
        return theme.gradients.festivalGradient;
      case 'minimal':
        return [COLORS.white, COLORS.gray100];
      default:
        return [COLORS.white, COLORS.gray100, COLORS.white];
    }
  };

  const getQuoteStyle = () => {
    switch (language) {
      case 'hi':
      case 'sa':
        return styles.quoteTextHindi;
      default:
        return styles.quoteText;
    }
  };

  return (
    <LinearGradient
      colors={getGradientColors()}
      style={[styles.container, variant === 'minimal' && styles.minimalContainer]}
      start={{ x: 0, y: 0 }}
      end={{ x: 1, y: 1 }}
    >
      <View style={styles.content}>
        {showQuoteMarks && variant !== 'minimal' && (
          <GradientText style={styles.quoteMark}>"</GradientText>
        )}
        
        <Text style={[styles.quoteText, getQuoteStyle()]}>
          {quote}
        </Text>
        
        {showQuoteMarks && variant !== 'minimal' && (
          <GradientText style={[styles.quoteMark, styles.quoteMarkEnd]}>"</GradientText>
        )}
        
        {author && (
          <View style={styles.authorContainer}>
            <View style={styles.divider} />
            <Text style={styles.authorText}>— {author}</Text>
          </View>
        )}
      </View>
    </LinearGradient>
  );
};

const styles = StyleSheet.create({
  container: {
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    marginVertical: theme.spacing.sm,
    ...theme.shadows.medium,
  },
  minimalContainer: {
    padding: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    ...theme.shadows.small,
  },
  content: {
    position: 'relative',
  },
  quoteMark: {
    fontSize: 48,
    opacity: 0.3,
    position: 'absolute',
    top: -20,
    left: -10,
  },
  quoteMarkEnd: {
    top: 'auto',
    bottom: -20,
    right: -10,
    left: 'auto',
    transform: [{ rotate: '180deg' }],
  },
  quoteText: {
    fontSize: 16,
    lineHeight: 24,
    color: COLORS.gray800,
    fontStyle: 'italic',
    textAlign: 'center',
    paddingHorizontal: theme.spacing.md,
  },
  quoteTextHindi: {
    fontSize: 18,
    lineHeight: 28,
    color: COLORS.gray800,
    fontStyle: 'italic',
    textAlign: 'center',
    paddingHorizontal: theme.spacing.md,
    fontFamily: 'NotoSansDevanagari-Regular',
  },
  authorContainer: {
    marginTop: theme.spacing.md,
    alignItems: 'center',
  },
  divider: {
    width: 40,
    height: 2,
    backgroundColor: COLORS.saffron,
    marginBottom: theme.spacing.sm,
  },
  authorText: {
    fontSize: 14,
    color: COLORS.gray600,
    fontStyle: 'italic',
  },
});