import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { theme } from '@constants/theme';

const AboutScreen = () => {
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <Text style={styles.title}>About DhanSetu</Text>
        <Text style={styles.version}>Version 1.0.0</Text>
        <Text style={styles.description}>
          India's Cultural DeFi Revolution
        </Text>
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: 10,
  },
  version: {
    fontSize: 18,
    color: theme.colors.saffron,
    marginBottom: 20,
  },
  description: {
    fontSize: 16,
    color: theme.colors.gray600,
    textAlign: 'center',
  },
});

export default AboutScreen;