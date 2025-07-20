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
  Share,
  Alert,
  Clipboard,
  Dimensions,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import QRCode from 'react-native-qrcode-svg';
import { useNavigation, useRoute } from '@react-navigation/native';
import ViewShot from 'react-native-view-shot';
import { CameraRoll } from '@react-native-camera-roll/camera-roll';

import { COLORS, theme } from '@constants/theme';
import { CulturalButton } from '@components/common/CulturalButton';
import { GradientText } from '@components/common/GradientText';
import { useAppSelector } from '@store/index';
import { useFestival } from '@contexts/FestivalContext';

const { width } = Dimensions.get('window');

interface RouteParams {
  coinType?: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN';
}

export const ReceiveScreen: React.FC = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const { currentAddress, currentCoin, dhanPataAddress } = useAppSelector((state) => state.wallet);
  const { currentFestival } = useFestival();
  
  const params = route.params as RouteParams;
  const selectedCoin = params?.coinType || currentCoin;
  
  const [activeTab, setActiveTab] = useState<'address' | 'request'>('address');
  const [requestAmount, setRequestAmount] = useState('');
  const [requestMemo, setRequestMemo] = useState('');
  const [showDhanPata, setShowDhanPata] = useState(true);
  const viewShotRef = React.useRef<ViewShot>(null);
  
  // Get display address
  const displayAddress = showDhanPata && dhanPataAddress ? dhanPataAddress : currentAddress || '';
  const shortAddress = displayAddress.includes('@dhan') 
    ? displayAddress 
    : `${displayAddress.slice(0, 10)}...${displayAddress.slice(-8)}`;
  
  // Generate QR data
  const generateQRData = () => {
    let baseData = displayAddress;
    
    if (activeTab === 'request' && (requestAmount || requestMemo)) {
      const params = new URLSearchParams();
      if (requestAmount) params.append('amount', requestAmount);
      if (requestMemo) params.append('memo', requestMemo);
      baseData = `deshchain:${displayAddress}?${params.toString()}`;
    }
    
    return baseData;
  };
  
  const copyAddress = () => {
    Clipboard.setString(displayAddress);
    Alert.alert('Copied!', `${showDhanPata && dhanPataAddress ? 'DhanPata ID' : 'Address'} copied to clipboard`);
  };
  
  const shareAddress = async () => {
    try {
      let message = `My ${selectedCoin} address:\n${displayAddress}`;
      
      if (dhanPataAddress && showDhanPata) {
        message = `Send me ${selectedCoin} using my DhanPata ID:\n${dhanPataAddress}\n\nNo need to remember long addresses!`;
      }
      
      if (activeTab === 'request' && requestAmount) {
        message += `\n\nAmount requested: ${requestAmount} ${selectedCoin}`;
        if (requestMemo) {
          message += `\nNote: ${requestMemo}`;
        }
      }
      
      await Share.share({
        message,
        title: `Receive ${selectedCoin}`,
      });
    } catch (error) {
      console.error('Share error:', error);
    }
  };
  
  const saveQRCode = async () => {
    try {
      const uri = await viewShotRef.current?.capture?.();
      if (uri) {
        await CameraRoll.save(uri, { type: 'photo' });
        Alert.alert('Saved!', 'QR code saved to gallery');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to save QR code');
    }
  };
  
  const renderAddressTab = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.content}>
        {/* Toggle between DhanPata and Address */}
        {dhanPataAddress && (
          <View style={styles.toggleContainer}>
            <TouchableOpacity
              style={[styles.toggleButton, showDhanPata && styles.toggleButtonActive]}
              onPress={() => setShowDhanPata(true)}
            >
              <Icon name="at" size={16} color={showDhanPata ? COLORS.white : COLORS.gray600} />
              <Text style={[styles.toggleText, showDhanPata && styles.toggleTextActive]}>
                DhanPata ID
              </Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={[styles.toggleButton, !showDhanPata && styles.toggleButtonActive]}
              onPress={() => setShowDhanPata(false)}
            >
              <Icon name="wallet" size={16} color={!showDhanPata ? COLORS.white : COLORS.gray600} />
              <Text style={[styles.toggleText, !showDhanPata && styles.toggleTextActive]}>
                Wallet Address
              </Text>
            </TouchableOpacity>
          </View>
        )}
        
        {/* QR Code */}
        <ViewShot ref={viewShotRef} style={styles.qrContainer}>
          <LinearGradient
            colors={currentFestival ? theme.gradients.festivalGradient : [COLORS.white, COLORS.white]}
            style={styles.qrCard}
          >
            {currentFestival && (
              <View style={styles.festivalBadge}>
                <Icon name="party-popper" size={16} color={COLORS.white} />
                <Text style={styles.festivalText}>{currentFestival.name}</Text>
              </View>
            )}
            
            <View style={styles.qrCodeWrapper}>
              <QRCode
                value={generateQRData()}
                size={width * 0.6}
                backgroundColor="white"
                color={COLORS.navy}
                logo={require('@assets/logo.png')}
                logoSize={50}
                logoBackgroundColor="white"
                logoBorderRadius={25}
              />
            </View>
            
            <Text style={styles.scanText}>Scan to send {selectedCoin}</Text>
          </LinearGradient>
        </ViewShot>
        
        {/* Address Display */}
        <TouchableOpacity
          style={styles.addressContainer}
          onPress={copyAddress}
          activeOpacity={0.8}
        >
          <View style={styles.addressHeader}>
            <Text style={styles.addressLabel}>
              {showDhanPata && dhanPataAddress ? 'DhanPata ID' : `${selectedCoin} Address`}
            </Text>
            <Icon name="content-copy" size={20} color={COLORS.saffron} />
          </View>
          <Text style={styles.address}>{shortAddress}</Text>
          {showDhanPata && dhanPataAddress && (
            <Text style={styles.addressHint}>Easy to remember, easy to share!</Text>
          )}
        </TouchableOpacity>
        
        {/* Action Buttons */}
        <View style={styles.actionButtons}>
          <TouchableOpacity style={styles.actionButton} onPress={copyAddress}>
            <Icon name="content-copy" size={24} color={COLORS.saffron} />
            <Text style={styles.actionText}>Copy</Text>
          </TouchableOpacity>
          
          <TouchableOpacity style={styles.actionButton} onPress={shareAddress}>
            <Icon name="share-variant" size={24} color={COLORS.green} />
            <Text style={styles.actionText}>Share</Text>
          </TouchableOpacity>
          
          <TouchableOpacity style={styles.actionButton} onPress={saveQRCode}>
            <Icon name="download" size={24} color={COLORS.navy} />
            <Text style={styles.actionText}>Save QR</Text>
          </TouchableOpacity>
        </View>
        
        {/* Instructions */}
        <View style={styles.instructions}>
          <Text style={styles.instructionTitle}>How to receive {selectedCoin}:</Text>
          <View style={styles.instructionItem}>
            <Text style={styles.instructionNumber}>1</Text>
            <Text style={styles.instructionText}>
              Share your {showDhanPata && dhanPataAddress ? 'DhanPata ID' : 'address'} with the sender
            </Text>
          </View>
          <View style={styles.instructionItem}>
            <Text style={styles.instructionNumber}>2</Text>
            <Text style={styles.instructionText}>
              Or let them scan your QR code
            </Text>
          </View>
          <View style={styles.instructionItem}>
            <Text style={styles.instructionNumber}>3</Text>
            <Text style={styles.instructionText}>
              Funds will appear in your wallet instantly
            </Text>
          </View>
        </View>
        
        {/* Cultural Note */}
        {showDhanPata && dhanPataAddress && (
          <View style={styles.culturalNote}>
            <Icon name="information" size={20} color={COLORS.info} />
            <Text style={styles.culturalNoteText}>
              DhanPata makes receiving money as easy as sharing your username. 
              No more copying long addresses!
            </Text>
          </View>
        )}
      </View>
    </ScrollView>
  );
  
  const renderRequestTab = () => (
    <ScrollView showsVerticalScrollIndicator={false}>
      <View style={styles.content}>
        <Text style={styles.requestTitle}>Create Payment Request</Text>
        <Text style={styles.requestSubtitle}>
          Generate a QR code with specific amount and memo
        </Text>
        
        {/* Amount Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Amount (Optional)</Text>
          <View style={styles.amountInputContainer}>
            <Text style={styles.currencySymbol}>{selectedCoin === 'DESHCHAIN' ? '₹' : ''}</Text>
            <TextInput
              style={styles.amountInput}
              placeholder="0.00"
              value={requestAmount}
              onChangeText={setRequestAmount}
              keyboardType="decimal-pad"
            />
            <Text style={styles.currencyCode}>{selectedCoin}</Text>
          </View>
        </View>
        
        {/* Memo Input */}
        <View style={styles.inputGroup}>
          <Text style={styles.inputLabel}>Memo (Optional)</Text>
          <TextInput
            style={styles.memoInput}
            placeholder="What's this payment for?"
            value={requestMemo}
            onChangeText={setRequestMemo}
            multiline
            numberOfLines={3}
            textAlignVertical="top"
          />
        </View>
        
        {/* Preview */}
        {(requestAmount || requestMemo) && (
          <View style={styles.requestPreview}>
            <Text style={styles.previewTitle}>Request Preview</Text>
            {requestAmount && (
              <View style={styles.previewItem}>
                <Text style={styles.previewLabel}>Amount:</Text>
                <Text style={styles.previewValue}>{requestAmount} {selectedCoin}</Text>
              </View>
            )}
            {requestMemo && (
              <View style={styles.previewItem}>
                <Text style={styles.previewLabel}>Memo:</Text>
                <Text style={styles.previewValue}>{requestMemo}</Text>
              </View>
            )}
          </View>
        )}
        
        {/* Generate Button */}
        <CulturalButton
          title="Generate Request QR"
          onPress={() => {
            // QR will update automatically with request data
            Alert.alert('Success!', 'Payment request QR generated');
          }}
          size="large"
          style={styles.generateButton}
          disabled={!requestAmount && !requestMemo}
        />
        
        {/* QR Code with Request */}
        {(requestAmount || requestMemo) && (
          <View style={styles.requestQRContainer}>
            <QRCode
              value={generateQRData()}
              size={width * 0.5}
              backgroundColor="white"
              color={COLORS.navy}
            />
            <Text style={styles.requestQRText}>Payment Request QR</Text>
          </View>
        )}
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
        <Text style={styles.headerTitle}>Receive {selectedCoin}</Text>
        <View style={{ width: 24 }} />
      </View>
      
      {/* Tabs */}
      <View style={styles.tabs}>
        <TouchableOpacity
          style={[styles.tab, activeTab === 'address' && styles.activeTab]}
          onPress={() => setActiveTab('address')}
        >
          <Icon 
            name="qrcode" 
            size={20} 
            color={activeTab === 'address' ? COLORS.saffron : COLORS.gray600}
          />
          <Text style={[
            styles.tabText,
            activeTab === 'address' && styles.activeTabText,
          ]}>
            My Address
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity
          style={[styles.tab, activeTab === 'request' && styles.activeTab]}
          onPress={() => setActiveTab('request')}
        >
          <Icon 
            name="cash-register" 
            size={20} 
            color={activeTab === 'request' ? COLORS.saffron : COLORS.gray600}
          />
          <Text style={[
            styles.tabText,
            activeTab === 'request' && styles.activeTabText,
          ]}>
            Request Payment
          </Text>
        </TouchableOpacity>
      </View>
      
      {/* Content */}
      {activeTab === 'address' ? renderAddressTab() : renderRequestTab()}
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
  tabs: {
    flexDirection: 'row',
    paddingHorizontal: theme.spacing.lg,
    paddingTop: theme.spacing.md,
  },
  tab: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
    gap: theme.spacing.sm,
  },
  activeTab: {
    borderBottomColor: COLORS.saffron,
  },
  tabText: {
    fontSize: 14,
    color: COLORS.gray600,
    fontWeight: '500',
  },
  activeTabText: {
    color: COLORS.saffron,
    fontWeight: 'bold',
  },
  content: {
    padding: theme.spacing.lg,
  },
  
  // Address Tab
  toggleContainer: {
    flexDirection: 'row',
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: 4,
    marginBottom: theme.spacing.lg,
  },
  toggleButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.sm,
    gap: theme.spacing.xs,
  },
  toggleButtonActive: {
    backgroundColor: COLORS.saffron,
  },
  toggleText: {
    fontSize: 14,
    color: COLORS.gray600,
    fontWeight: '500',
  },
  toggleTextActive: {
    color: COLORS.white,
    fontWeight: 'bold',
  },
  qrContainer: {
    alignItems: 'center',
    marginBottom: theme.spacing.lg,
  },
  qrCard: {
    padding: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    alignItems: 'center',
    ...theme.shadows.medium,
  },
  festivalBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.full,
    marginBottom: theme.spacing.md,
  },
  festivalText: {
    fontSize: 12,
    color: COLORS.white,
    fontWeight: 'bold',
    marginLeft: theme.spacing.xs,
  },
  qrCodeWrapper: {
    padding: theme.spacing.md,
    backgroundColor: COLORS.white,
    borderRadius: theme.borderRadius.md,
  },
  scanText: {
    fontSize: 14,
    color: currentFestival => currentFestival ? COLORS.white : COLORS.gray600,
    marginTop: theme.spacing.md,
  },
  addressContainer: {
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  addressHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  addressLabel: {
    fontSize: 12,
    color: COLORS.gray600,
    fontWeight: '600',
  },
  address: {
    fontSize: 16,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  addressHint: {
    fontSize: 12,
    color: COLORS.success,
    marginTop: theme.spacing.xs,
    fontStyle: 'italic',
  },
  actionButtons: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: theme.spacing.lg,
  },
  actionButton: {
    alignItems: 'center',
    padding: theme.spacing.md,
  },
  actionText: {
    fontSize: 12,
    color: COLORS.gray700,
    marginTop: theme.spacing.xs,
  },
  instructions: {
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.lg,
    marginBottom: theme.spacing.lg,
  },
  instructionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.md,
  },
  instructionItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.md,
  },
  instructionNumber: {
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: COLORS.saffron,
    color: COLORS.white,
    textAlign: 'center',
    lineHeight: 24,
    fontWeight: 'bold',
    marginRight: theme.spacing.md,
  },
  instructionText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.gray700,
    lineHeight: 20,
  },
  culturalNote: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    backgroundColor: COLORS.info + '10',
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
  },
  culturalNoteText: {
    flex: 1,
    fontSize: 14,
    color: COLORS.info,
    marginLeft: theme.spacing.sm,
    lineHeight: 20,
  },
  
  // Request Tab
  requestTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.gray900,
    marginBottom: theme.spacing.sm,
  },
  requestSubtitle: {
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
  memoInput: {
    backgroundColor: COLORS.gray100,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    fontSize: 16,
    color: COLORS.gray900,
    minHeight: 80,
  },
  requestPreview: {
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.md,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.lg,
  },
  previewTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.gray700,
    marginBottom: theme.spacing.sm,
  },
  previewItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.xs,
  },
  previewLabel: {
    fontSize: 14,
    color: COLORS.gray600,
  },
  previewValue: {
    fontSize: 14,
    color: COLORS.gray900,
    fontWeight: '500',
  },
  generateButton: {
    marginBottom: theme.spacing.lg,
  },
  requestQRContainer: {
    alignItems: 'center',
    backgroundColor: COLORS.gray50,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
  },
  requestQRText: {
    fontSize: 14,
    color: COLORS.gray600,
    marginTop: theme.spacing.md,
  },
});