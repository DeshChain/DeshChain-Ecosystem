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
  Image,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector, useAppDispatch } from '@store/index';
import { createDhanPataAddress } from '@store/slices/walletSlice';

interface MitraProfile {
  name: string;
  dhanPataId: string;
  trustScore: number;
  totalTransactions: number;
  memberSince: string;
  verificationLevel: 'basic' | 'kyc' | 'enhanced';
  badges: string[];
  pincode?: string;
  preferredLanguage: string;
  culturalQuote?: string;
}

export const ProfileScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  const { currentAddress, dhanPataAddress } = useAppSelector((state) => state.wallet);
  
  const [isEditingProfile, setIsEditingProfile] = useState(false);
  const [showDhanPataSetup, setShowDhanPataSetup] = useState(false);
  
  // Profile state
  const [profile, setProfile] = useState<MitraProfile>({
    name: 'Satoshi Nakamoto',
    dhanPataId: dhanPataAddress || '',
    trustScore: 85,
    totalTransactions: 156,
    memberSince: 'January 2024',
    verificationLevel: 'kyc',
    badges: ['Early Adopter', 'Trusted Trader', 'Festival Participant'],
    pincode: '110001',
    preferredLanguage: 'Hindi',
    culturalQuote: 'वसुधैव कुटुम्बकम्',
  });

  // DhanPata setup state
  const [dhanPataUsername, setDhanPataUsername] = useState('');
  const [isCheckingAvailability, setIsCheckingAvailability] = useState(false);

  const checkDhanPataAvailability = async () => {
    if (!dhanPataUsername || dhanPataUsername.length < 3) {
      Alert.alert('Error', 'Username must be at least 3 characters');
      return;
    }

    setIsCheckingAvailability(true);
    
    // Simulate API call
    setTimeout(() => {
      setIsCheckingAvailability(false);
      const isAvailable = Math.random() > 0.3; // 70% chance of availability
      
      if (isAvailable) {
        Alert.alert(
          'Available!',
          `${dhanPataUsername}@dhan is available. Would you like to claim it?`,
          [
            { text: 'Cancel', style: 'cancel' },
            { text: 'Claim', onPress: claimDhanPata },
          ]
        );
      } else {
        Alert.alert('Not Available', `${dhanPataUsername}@dhan is already taken`);
      }
    }, 1500);
  };

  const claimDhanPata = async () => {
    try {
      await dispatch(createDhanPataAddress({
        username: dhanPataUsername,
        address: currentAddress!,
      })).unwrap();
      
      setProfile({ ...profile, dhanPataId: `${dhanPataUsername}@dhan` });
      setShowDhanPataSetup(false);
      Alert.alert('Success!', `You now own ${dhanPataUsername}@dhan`);
    } catch (error) {
      Alert.alert('Error', 'Failed to claim DhanPata address');
    }
  };

  const renderHeader = () => (
    <LinearGradient
      colors={theme.gradients.indianFlag}
      style={styles.header}
    >
      <View style={styles.avatarContainer}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>
            {profile.name.split(' ').map(n => n[0]).join('')}
          </Text>
        </View>
        <TouchableOpacity style={styles.editAvatarButton}>
          <Icon name="camera" size={20} color={COLORS.white} />
        </TouchableOpacity>
      </View>
      
      <Text style={styles.profileName}>{profile.name}</Text>
      
      {profile.dhanPataId ? (
        <View style={styles.dhanPataContainer}>
          <Icon name="at" size={16} color={COLORS.white} />
          <Text style={styles.dhanPataText}>{profile.dhanPataId}</Text>
        </View>
      ) : (
        <TouchableOpacity
          style={styles.createDhanPataButton}
          onPress={() => setShowDhanPataSetup(true)}
        >
          <Icon name="plus" size={16} color={COLORS.saffron} />
          <Text style={styles.createDhanPataText}>Create DhanPata ID</Text>
        </TouchableOpacity>
      )}
      
      <View style={styles.statsContainer}>
        <View style={styles.statItem}>
          <Text style={styles.statValue}>{profile.trustScore}</Text>
          <Text style={styles.statLabel}>Trust Score</Text>
        </View>
        <View style={styles.statDivider} />
        <View style={styles.statItem}>
          <Text style={styles.statValue}>{profile.totalTransactions}</Text>
          <Text style={styles.statLabel}>Transactions</Text>
        </View>
        <View style={styles.statDivider} />
        <View style={styles.statItem}>
          <Text style={styles.statValue}>{profile.badges.length}</Text>
          <Text style={styles.statLabel}>Badges</Text>
        </View>
      </View>
    </LinearGradient>
  );

  const renderDhanPataSetup = () => (
    <View style={styles.dhanPataSetupContainer}>
      <View style={styles.dhanPataHeader}>
        <Text style={styles.dhanPataTitle}>Create Your DhanPata ID</Text>
        <TouchableOpacity onPress={() => setShowDhanPataSetup(false)}>
          <Icon name="close" size={24} color={COLORS.gray700} />
        </TouchableOpacity>
      </View>
      
      <Text style={styles.dhanPataDescription}>
        Choose a unique username for easy payments. People can send you money using just your DhanPata ID!
      </Text>
      
      <View style={styles.dhanPataInputContainer}>
        <TextInput
          style={styles.dhanPataInput}
          placeholder="Choose username"
          value={dhanPataUsername}
          onChangeText={setDhanPataUsername}
          autoCapitalize="none"
          autoCorrect={false}
        />
        <Text style={styles.dhanPataSuffix}>@dhan</Text>
      </View>
      
      <Text style={styles.dhanPataHint}>
        Minimum 3 characters, only letters and numbers
      </Text>
      
      <CulturalButton
        title="Check Availability"
        onPress={checkDhanPataAvailability}
        loading={isCheckingAvailability}
        size="large"
        style={styles.checkButton}
      />
      
      <View style={styles.dhanPataFeatures}>
        <View style={styles.featureItem}>
          <Icon name="check-circle" size={20} color={COLORS.success} />
          <Text style={styles.featureText}>No need to share long addresses</Text>
        </View>
        <View style={styles.featureItem}>
          <Icon name="check-circle" size={20} color={COLORS.success} />
          <Text style={styles.featureText}>Works across all DeshChain apps</Text>
        </View>
        <View style={styles.featureItem}>
          <Icon name="check-circle" size={20} color={COLORS.success} />
          <Text style={styles.featureText}>Permanent and transferable</Text>
        </View>
      </View>
    </View>
  );

  const renderProfileInfo = () => (
    <View style={styles.profileSection}>
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionTitle}>Profile Information</Text>
        <TouchableOpacity onPress={() => setIsEditingProfile(!isEditingProfile)}>
          <Icon name={isEditingProfile ? 'check' : 'pencil'} size={20} color={COLORS.saffron} />
        </TouchableOpacity>
      </View>
      
      <View style={styles.infoItem}>
        <Icon name="shield-check" size={20} color={COLORS.green} />
        <Text style={styles.infoLabel}>Verification Level</Text>
        <View style={[styles.verificationBadge, styles[`verification_${profile.verificationLevel}`]]}>
          <Text style={styles.verificationText}>
            {profile.verificationLevel.toUpperCase()}
          </Text>
        </View>
      </View>
      
      <View style={styles.infoItem}>
        <Icon name="map-marker" size={20} color={COLORS.saffron} />
        <Text style={styles.infoLabel}>PIN Code</Text>
        <Text style={styles.infoValue}>{profile.pincode || 'Not set'}</Text>
      </View>
      
      <View style={styles.infoItem}>
        <Icon name="translate" size={20} color={COLORS.navy} />
        <Text style={styles.infoLabel}>Preferred Language</Text>
        <Text style={styles.infoValue}>{profile.preferredLanguage}</Text>
      </View>
      
      <View style={styles.infoItem}>
        <Icon name="calendar" size={20} color={COLORS.gray600} />
        <Text style={styles.infoLabel}>Member Since</Text>
        <Text style={styles.infoValue}>{profile.memberSince}</Text>
      </View>
    </View>
  );

  const renderBadges = () => (
    <View style={styles.profileSection}>
      <Text style={styles.sectionTitle}>Achievements & Badges</Text>
      
      <View style={styles.badgesContainer}>
        {profile.badges.map((badge, index) => (
          <View key={index} style={styles.badge}>
            <Icon 
              name={getBadgeIcon(badge)} 
              size={24} 
              color={getBadgeColor(badge)} 
            />
            <Text style={styles.badgeText}>{badge}</Text>
          </View>
        ))}
      </View>
    </View>
  );

  const renderQuickActions = () => (
    <View style={styles.profileSection}>
      <Text style={styles.sectionTitle}>Quick Actions</Text>
      
      <TouchableOpacity style={styles.actionItem}>
        <Icon name="qrcode" size={24} color={COLORS.saffron} />
        <Text style={styles.actionText}>Share Payment QR</Text>
        <Icon name="chevron-right" size={20} color={COLORS.gray400} />
      </TouchableOpacity>
      
      <TouchableOpacity style={styles.actionItem}>
        <Icon name="history" size={24} color={COLORS.green} />
        <Text style={styles.actionText}>Transaction History</Text>
        <Icon name="chevron-right" size={20} color={COLORS.gray400} />
      </TouchableOpacity>
      
      <TouchableOpacity style={styles.actionItem}>
        <Icon name="account-group" size={24} color={COLORS.navy} />
        <Text style={styles.actionText}>Refer Friends</Text>
        <Icon name="chevron-right" size={20} color={COLORS.gray400} />
      </TouchableOpacity>
      
      <TouchableOpacity 
        style={styles.actionItem}
        onPress={() => navigation.navigate('Settings')}
      >
        <Icon name="cog" size={24} color={COLORS.gray600} />
        <Text style={styles.actionText}>Settings</Text>
        <Icon name="chevron-right" size={20} color={COLORS.gray400} />
      </TouchableOpacity>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {renderHeader()}
        
        {showDhanPataSetup ? (
          renderDhanPataSetup()
        ) : (
          <>
            {renderProfileInfo()}
            {renderBadges()}
            {renderQuickActions()}
          </>
        )}
        
        {/* Cultural Quote */}
        {profile.culturalQuote && (
          <View style={styles.quoteContainer}>
            <Text style={styles.quoteText}>"{profile.culturalQuote}"</Text>
            <Text style={styles.quoteTranslation}>The world is one family</Text>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
};

// Helper functions
const getBadgeIcon = (badge: string): string => {
  switch (badge) {
    case 'Early Adopter': return 'star';
    case 'Trusted Trader': return 'shield-check';
    case 'Festival Participant': return 'party-popper';
    default: return 'medal';
  }
};

const getBadgeColor = (badge: string): string => {
  switch (badge) {
    case 'Early Adopter': return COLORS.saffron;
    case 'Trusted Trader': return COLORS.green;
    case 'Festival Participant': return COLORS.festivalPrimary;
    default: return COLORS.gray600;
  }
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.white,
  },
  header: {
    paddingTop: theme.spacing.xl,
    paddingBottom: theme.spacing.xl,
    alignItems: 'center',
  },
  avatarContainer: {
    position: 'relative',
    marginBottom: theme.spacing.md,
  },
  avatar: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.white,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.medium,
  },
  avatarText: {
    fontSize: 36,
    fontWeight: 'bold',
    color: COLORS.saffron,
  },
  editAvatarButton: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: COLORS.saffron,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.small,
  },
  profileName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
    marginBottom: theme.spacing.sm,
  },
  dhanPataContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.full,
    marginBottom: theme.spacing.lg,
  },
  dhanPataText: {
    fontSize: 14,
    color: COLORS.white,
    marginLeft: 4,
  },
  createDhanPataButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.white,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.full,
    marginBottom: theme.spacing.lg,
  },
  createDhanPataText: {
    fontSize: 14,
    color: COLORS.saffron,
    marginLeft: 4,
    fontWeight: '600',
  },
  statsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  statItem: {
    alignItems: 'center',
    paddingHorizontal: theme.spacing.lg,
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.white,
  },
  statLabel: {
    fontSize: 12,
    color: COLORS.white,
    opacity: 0.8,
    marginTop: 2,
  },
  statDivider: {
    width: 1,
    height: 30,
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
  },
  
  // DhanPata Setup
  dhanPataSetupContainer: {
    margin: theme.spacing.lg,
    padding: theme.spacing.lg,
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.lg,
    ...theme.shadows.medium,
  },
  dhanPataHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  dhanPataTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  dhanPataDescription: {
    fontSize: 14,
    color: COLORS.gray600,
    lineHeight: 20,
    marginBottom: theme.spacing.lg,
  },
  dhanPataInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    paddingRight: theme.spacing.md,
    marginBottom: theme.spacing.sm,
  },
  dhanPataInput: {
    flex: 1,
    fontSize: 18,
    color: COLORS.gray900,
    padding: theme.spacing.md,
  },
  dhanPataSuffix: {
    fontSize: 18,
    color: COLORS.gray600,
    fontWeight: '600',
  },
  dhanPataHint: {
    fontSize: 12,
    color: COLORS.gray500,
    marginBottom: theme.spacing.lg,
  },
  checkButton: {
    marginBottom: theme.spacing.lg,
  },
  dhanPataFeatures: {
    paddingTop: theme.spacing.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.gray200,
  },
  featureItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  featureText: {
    fontSize: 14,
    color: COLORS.gray700,
    marginLeft: theme.spacing.sm,
    flex: 1,
  },
  
  // Profile Info
  profileSection: {
    margin: theme.spacing.lg,
    marginBottom: 0,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.gray900,
  },
  infoItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray100,
  },
  infoLabel: {
    flex: 1,
    fontSize: 14,
    color: COLORS.gray700,
    marginLeft: theme.spacing.md,
  },
  infoValue: {
    fontSize: 14,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  verificationBadge: {
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: 2,
    borderRadius: theme.borderRadius.sm,
  },
  verification_basic: {
    backgroundColor: COLORS.gray300,
  },
  verification_kyc: {
    backgroundColor: COLORS.warning + '20',
  },
  verification_enhanced: {
    backgroundColor: COLORS.success + '20',
  },
  verificationText: {
    fontSize: 10,
    fontWeight: 'bold',
    color: COLORS.gray800,
  },
  
  // Badges
  badgesContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginTop: theme.spacing.md,
    gap: theme.spacing.sm,
  },
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.gray100,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.full,
  },
  badgeText: {
    fontSize: 12,
    color: COLORS.gray700,
    marginLeft: theme.spacing.xs,
    fontWeight: '500',
  },
  
  // Quick Actions
  actionItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.gray100,
  },
  actionText: {
    flex: 1,
    fontSize: 16,
    color: COLORS.gray900,
    marginLeft: theme.spacing.md,
  },
  
  // Quote
  quoteContainer: {
    alignItems: 'center',
    padding: theme.spacing.xl,
  },
  quoteText: {
    fontSize: 16,
    color: COLORS.gray700,
    fontStyle: 'italic',
    textAlign: 'center',
  },
  quoteTranslation: {
    fontSize: 14,
    color: COLORS.gray500,
    marginTop: theme.spacing.sm,
  },
});