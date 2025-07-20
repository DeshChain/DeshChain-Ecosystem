import React, { useEffect, useState, useCallback } from 'react';
import {
  View,
  ScrollView,
  StyleSheet,
  RefreshControl,
  Dimensions,
  Alert,
  Linking,
} from 'react-native';
import {
  Text,
  Card,
  Button,
  Chip,
  Surface,
  IconButton,
  FAB,
  List,
  Divider,
  Searchbar,
  SegmentedButtons,
  Avatar,
  Badge,
  ProgressBar,
} from 'react-native-paper';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';
import * as Animatable from 'react-native-animatable';

import { RootState } from '../../store';
import { theme, spacing, colors } from '../../theme';
import { SevaMitraService } from '../../services/SevaMitraService';
import { LocationService } from '../../services/LocationService';

const { width } = Dimensions.get('window');

interface SevaMitra {
  id: string;
  name: string;
  trustScore: number;
  completedServices: number;
  responseTime: string;
  isOnline: boolean;
  location: {
    address: string;
    city: string;
    state: string;
    pincode: string;
    coordinates: {
      latitude: number;
      longitude: number;
    };
  };
  services: string[];
  languages: string[];
  workingHours: {
    start: string;
    end: string;
  };
  commission: {
    cashIn: number;
    cashOut: number;
  };
  maxAmount: number;
  minAmount: number;
  avatar?: string;
  verificationLevel: 'basic' | 'verified' | 'premium';
  lastActive: Date;
  distance?: number;
}

interface ServiceRequest {
  id: string;
  type: 'cash_in' | 'cash_out' | 'money_transfer';
  amount: number;
  status: 'pending' | 'accepted' | 'in_progress' | 'completed' | 'cancelled';
  sevaMitra: {
    id: string;
    name: string;
    trustScore: number;
  };
  location: {
    address: string;
    pincode: string;
  };
  scheduledTime?: Date;
  completedAt?: Date;
  fee: number;
}

export const SevaMitraScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance } = useSelector((state: RootState) => state.wallet);
  const { language } = useSelector((state: RootState) => state.settings);
  
  const [refreshing, setRefreshing] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [sevaMitras, setSevaMitras] = useState<SevaMitra[]>([]);
  const [serviceRequests, setServiceRequests] = useState<ServiceRequest[]>([]);
  const [selectedTab, setSelectedTab] = useState<'discover' | 'requests'>('discover');
  const [filterService, setFilterService] = useState<'all' | 'cash_in' | 'cash_out'>('all');
  const [sortBy, setSortBy] = useState<'distance' | 'trust' | 'commission'>('distance');
  const [userLocation, setUserLocation] = useState<{ latitude: number; longitude: number } | null>(null);
  const [isLocationLoading, setIsLocationLoading] = useState(false);

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      // Get user location first
      if (!userLocation) {
        setIsLocationLoading(true);
        const location = await LocationService.getCurrentLocation();
        setUserLocation(location);
        setIsLocationLoading(false);
      }
      
      // Fetch nearby Seva Mitras
      const mitrasData = await SevaMitraService.getNearbyMitras({
        location: userLocation,
        search: searchQuery,
        serviceType: filterService === 'all' ? undefined : filterService,
        sortBy,
        radius: 10, // 10km radius
      });
      setSevaMitras(mitrasData);
      
      // Fetch user's service requests
      const requestsData = await SevaMitraService.getUserServiceRequests();
      setServiceRequests(requestsData);
      
    } catch (error) {
      console.error('Error refreshing Seva Mitra data:', error);
      Alert.alert('Error', 'Failed to refresh data. Please try again.');
    } finally {
      setRefreshing(false);
      setIsLocationLoading(false);
    }
  }, [userLocation, searchQuery, filterService, sortBy]);

  useEffect(() => {
    onRefresh();
  }, [onRefresh]);

  const getTrustScoreBadge = (trustScore: number) => {
    if (trustScore >= 90) return { label: 'Diamond', color: colors.diamond };
    if (trustScore >= 80) return { label: 'Platinum', color: colors.platinum };
    if (trustScore >= 70) return { label: 'Gold', color: colors.goldBadge };
    if (trustScore >= 60) return { label: 'Silver', color: colors.silver };
    if (trustScore >= 50) return { label: 'Bronze', color: colors.bronze };
    return { label: 'New', color: colors.error };
  };

  const getVerificationIcon = (level: string) => {
    switch (level) {
      case 'premium': return { icon: 'verified', color: colors.gold };
      case 'verified': return { icon: 'verified-user', color: colors.success };
      case 'basic': return { icon: 'person', color: colors.textSecondary };
      default: return { icon: 'person', color: colors.textSecondary };
    }
  };

  const formatCurrency = (amount: number) => {
    return `₹${amount.toLocaleString('en-IN')}`;
  };

  const formatDistance = (distance: number) => {
    if (distance < 1) {
      return `${Math.round(distance * 1000)}m`;
    }
    return `${distance.toFixed(1)}km`;
  };

  const getServiceIcon = (service: string) => {
    switch (service.toLowerCase()) {
      case 'cash_in': return 'input';
      case 'cash_out': return 'output';
      case 'money_transfer': return 'swap-horiz';
      case 'bill_payment': return 'receipt';
      default: return 'payment';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return colors.warning;
      case 'accepted': return colors.info;
      case 'in_progress': return colors.primary;
      case 'completed': return colors.success;
      case 'cancelled': return colors.error;
      default: return colors.textSecondary;
    }
  };

  const callSevaMitra = (phoneNumber: string) => {
    Linking.openURL(`tel:${phoneNumber}`);
  };

  const openWhatsApp = (phoneNumber: string, message: string) => {
    const url = `whatsapp://send?phone=${phoneNumber}&text=${encodeURIComponent(message)}`;
    Linking.openURL(url).catch(() => {
      Alert.alert('Error', 'WhatsApp is not installed on this device');
    });
  };

  const renderSevaMitraCard = (mitra: SevaMitra) => {
    const trustBadge = getTrustScoreBadge(mitra.trustScore);
    const verification = getVerificationIcon(mitra.verificationLevel);
    
    return (
      <Card key={mitra.id} style={styles.mitraCard}>
        <Card.Content>
          <View style={styles.mitraHeader}>
            <View style={styles.mitraAvatarContainer}>
              <Avatar.Text
                size={48}
                label={mitra.name.charAt(0)}
                style={{ backgroundColor: trustBadge.color }}
              />
              {mitra.isOnline && (
                <Badge
                  style={styles.onlineBadge}
                  size={16}
                />
              )}
              <Icon
                name={verification.icon}
                size={16}
                color={verification.color}
                style={styles.verificationIcon}
              />
            </View>
            
            <View style={styles.mitraInfo}>
              <Text style={styles.mitraName}>{mitra.name}</Text>
              <Chip
                textStyle={{ fontSize: 10, color: 'white' }}
                style={[styles.trustChip, { backgroundColor: trustBadge.color }]}
              >
                {trustBadge.label} {mitra.trustScore}
              </Chip>
              <Text style={styles.mitraStats}>
                {mitra.completedServices} services • {mitra.responseTime}
              </Text>
            </View>
            
            <View style={styles.mitraDistance}>
              {mitra.distance && (
                <Chip
                  icon="location-on"
                  compact
                  style={styles.distanceChip}
                >
                  {formatDistance(mitra.distance)}
                </Chip>
              )}
              <Text style={styles.lastActive}>
                {mitra.isOnline ? 'Online' : `Active ${mitra.lastActive.toLocaleDateString()}`}
              </Text>
            </View>
          </View>

          <Divider style={styles.mitraDivider} />

          <View style={styles.mitraDetails}>
            <View style={styles.locationInfo}>
              <Icon name="location-on" size={16} color={colors.textSecondary} />
              <Text style={styles.locationText}>
                {mitra.location.address}, {mitra.location.city}
              </Text>
            </View>
            
            <View style={styles.servicesInfo}>
              <Text style={styles.servicesLabel}>
                {language === 'hi' ? 'सेवाएं' : 'Services'}
              </Text>
              <View style={styles.servicesList}>
                {mitra.services.slice(0, 3).map((service, index) => (
                  <Chip
                    key={index}
                    icon={getServiceIcon(service)}
                    compact
                    style={styles.serviceChip}
                  >
                    {service.replace('_', ' ').toUpperCase()}
                  </Chip>
                ))}
                {mitra.services.length > 3 && (
                  <Chip compact style={styles.serviceChip}>
                    +{mitra.services.length - 3}
                  </Chip>
                )}
              </View>
            </View>

            <View style={styles.limitsInfo}>
              <View style={styles.limitItem}>
                <Text style={styles.limitLabel}>Amount Range</Text>
                <Text style={styles.limitValue}>
                  {formatCurrency(mitra.minAmount)} - {formatCurrency(mitra.maxAmount)}
                </Text>
              </View>
              <View style={styles.limitItem}>
                <Text style={styles.limitLabel}>Commission</Text>
                <Text style={styles.limitValue}>
                  {mitra.commission.cashIn}% / {mitra.commission.cashOut}%
                </Text>
              </View>
            </View>

            <View style={styles.workingHours}>
              <Icon name="schedule" size={16} color={colors.textSecondary} />
              <Text style={styles.workingHoursText}>
                {mitra.workingHours.start} - {mitra.workingHours.end}
              </Text>
              <View style={styles.languagesList}>
                {mitra.languages.slice(0, 2).map((lang, index) => (
                  <Chip
                    key={index}
                    compact
                    style={styles.languageChip}
                  >
                    {lang}
                  </Chip>
                ))}
              </View>
            </View>
          </View>

          <View style={styles.mitraActions}>
            <Button
              mode="outlined"
              icon="phone"
              style={styles.mitraActionButton}
              onPress={() => callSevaMitra('+91' + mitra.id)}
            >
              {language === 'hi' ? 'कॉल' : 'Call'}
            </Button>
            <Button
              mode="outlined"
              icon="chat"
              style={styles.mitraActionButton}
              onPress={() => openWhatsApp('+91' + mitra.id, 'Hi, I need help with money transfer')}
            >
              Chat
            </Button>
            <Button
              mode="contained"
              style={styles.mitraActionButton}
              onPress={() => navigation.navigate('SevaMitraDetails' as never, { 
                mitraId: mitra.id 
              } as never)}
            >
              {language === 'hi' ? 'बुक करें' : 'Book'}
            </Button>
          </View>
        </Card.Content>
      </Card>
    );
  };

  const renderServiceRequestCard = (request: ServiceRequest) => {
    return (
      <Card key={request.id} style={styles.requestCard}>
        <Card.Content>
          <View style={styles.requestHeader}>
            <View>
              <Text style={styles.requestTitle}>
                {request.type.replace('_', ' ').toUpperCase()} Service
              </Text>
              <Text style={styles.requestAmount}>
                {formatCurrency(request.amount)}
              </Text>
            </View>
            <Chip
              textStyle={{ fontSize: 10, color: 'white' }}
              style={[
                styles.requestStatusChip,
                { backgroundColor: getStatusColor(request.status) }
              ]}
            >
              {request.status.replace('_', ' ').toUpperCase()}
            </Chip>
          </View>

          <View style={styles.requestDetails}>
            <View style={styles.requestInfo}>
              <Text style={styles.requestInfoLabel}>Seva Mitra</Text>
              <Text style={styles.requestInfoValue}>
                {request.sevaMitra.name} (Trust: {request.sevaMitra.trustScore})
              </Text>
            </View>
            
            <View style={styles.requestInfo}>
              <Text style={styles.requestInfoLabel}>Location</Text>
              <Text style={styles.requestInfoValue}>
                {request.location.address}
              </Text>
            </View>
            
            <View style={styles.requestInfo}>
              <Text style={styles.requestInfoLabel}>Fee</Text>
              <Text style={styles.requestInfoValue}>
                {formatCurrency(request.fee)}
              </Text>
            </View>

            {request.scheduledTime && (
              <View style={styles.requestInfo}>
                <Text style={styles.requestInfoLabel}>Scheduled</Text>
                <Text style={styles.requestInfoValue}>
                  {request.scheduledTime.toLocaleDateString()} at {request.scheduledTime.toLocaleTimeString()}
                </Text>
              </View>
            )}
          </View>

          <View style={styles.requestActions}>
            {request.status === 'pending' && (
              <Button
                mode="outlined"
                style={styles.requestActionButton}
                onPress={() => {
                  Alert.alert(
                    'Cancel Request',
                    'Are you sure you want to cancel this service request?',
                    [
                      { text: 'No', style: 'cancel' },
                      { text: 'Yes', onPress: () => console.log('Cancel request:', request.id) },
                    ]
                  );
                }}
              >
                Cancel
              </Button>
            )}
            <Button
              mode="contained"
              style={styles.requestActionButton}
              onPress={() => {
                // Navigate to request details or chat
                console.log('View request details:', request.id);
              }}
            >
              Details
            </Button>
          </View>
        </Card.Content>
      </Card>
    );
  };

  return (
    <View style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {/* Header */}
        <Animatable.View animation="fadeInDown" duration={800} style={styles.header}>
          <View style={styles.headerContent}>
            <Text style={styles.headerTitle}>
              {language === 'hi' ? 'सेवा मित्र' : 'Seva Mitra'}
            </Text>
            <View style={styles.headerActions}>
              <IconButton
                icon="map"
                size={24}
                iconColor="white"
                onPress={() => navigation.navigate('SevaMitraMap' as never)}
              />
              <IconButton
                icon="person-add"
                size={24}
                iconColor="white"
                onPress={() => navigation.navigate('SevaMitraRegistration' as never)}
              />
            </View>
          </View>
        </Animatable.View>

        {/* Location Loading */}
        {isLocationLoading && (
          <Animatable.View animation="fadeIn" duration={500}>
            <Card style={styles.locationCard}>
              <Card.Content>
                <View style={styles.locationLoading}>
                  <Icon name="location-searching" size={24} color={colors.primary} />
                  <Text style={styles.locationLoadingText}>
                    {language === 'hi' ? 'स्थान खोजा जा रहा है...' : 'Finding your location...'}
                  </Text>
                </View>
                <ProgressBar indeterminate style={styles.locationProgress} />
              </Card.Content>
            </Card>
          </Animatable.View>
        )}

        {/* Tab Navigator */}
        <Animatable.View animation="fadeInUp" duration={800} delay={200}>
          <SegmentedButtons
            value={selectedTab}
            onValueChange={(value) => setSelectedTab(value as 'discover' | 'requests')}
            buttons={[
              { 
                value: 'discover', 
                label: language === 'hi' ? 'खोजें' : 'Discover',
                icon: 'search'
              },
              { 
                value: 'requests', 
                label: language === 'hi' ? 'रिक्वेस्ट्स' : 'My Requests',
                icon: 'list'
              },
            ]}
            style={styles.tabButtons}
          />
        </Animatable.View>

        {selectedTab === 'discover' && (
          <>
            {/* Search and Filters */}
            <Animatable.View animation="fadeInUp" duration={800} delay={400}>
              <View style={styles.searchContainer}>
                <Searchbar
                  placeholder={language === 'hi' ? 'सेवा मित्र खोजें...' : 'Search Seva Mitra...'}
                  onChangeText={setSearchQuery}
                  value={searchQuery}
                  style={styles.searchBar}
                />
                
                <View style={styles.filterChips}>
                  <Chip
                    selected={filterService === 'all'}
                    onPress={() => setFilterService('all')}
                    style={styles.filterChip}
                  >
                    All
                  </Chip>
                  <Chip
                    selected={filterService === 'cash_in'}
                    onPress={() => setFilterService('cash_in')}
                    style={styles.filterChip}
                  >
                    Cash In
                  </Chip>
                  <Chip
                    selected={filterService === 'cash_out'}
                    onPress={() => setFilterService('cash_out')}
                    style={styles.filterChip}
                  >
                    Cash Out
                  </Chip>
                </View>
              </View>
            </Animatable.View>

            {/* Seva Mitras List */}
            <Animatable.View animation="fadeInUp" duration={800} delay={600}>
              {sevaMitras.length > 0 ? (
                <View style={styles.mitrasList}>
                  {sevaMitras.map(renderSevaMitraCard)}
                </View>
              ) : (
                <Card style={styles.emptyCard}>
                  <Card.Content style={styles.emptyContent}>
                    <Icon name="location-off" size={64} color={colors.disabled} />
                    <Text style={styles.emptyText}>
                      {language === 'hi' ? 'कोई सेवा मित्र नहीं मिला' : 'No Seva Mitra found'}
                    </Text>
                    <Text style={styles.emptySubtext}>
                      {language === 'hi' 
                        ? 'अपना स्थान चालू करें या खोज की सीमा बढ़ाएं'
                        : 'Enable location or expand search radius'}
                    </Text>
                  </Card.Content>
                </Card>
              )}
            </Animatable.View>
          </>
        )}

        {selectedTab === 'requests' && (
          <Animatable.View animation="fadeInUp" duration={800} delay={400}>
            {serviceRequests.length > 0 ? (
              <View style={styles.requestsList}>
                {serviceRequests.map(renderServiceRequestCard)}
              </View>
            ) : (
              <Card style={styles.emptyCard}>
                <Card.Content style={styles.emptyContent}>
                  <Icon name="assignment" size={64} color={colors.disabled} />
                  <Text style={styles.emptyText}>
                    {language === 'hi' ? 'कोई सेवा रिक्वेस्ट नहीं' : 'No service requests'}
                  </Text>
                  <Text style={styles.emptySubtext}>
                    {language === 'hi' 
                      ? 'सेवा मित्र बुक करने के लिए खोजें टैब देखें'
                      : 'Check discover tab to book a Seva Mitra'}
                  </Text>
                </Card.Content>
              </Card>
            )}
          </Animatable.View>
        )}

        <View style={styles.bottomSpacer} />
      </ScrollView>

      {/* Floating Action Button */}
      <FAB
        icon="add"
        style={styles.fab}
        onPress={() => {
          Alert.alert(
            'Quick Service',
            'What do you need help with?',
            [
              { text: 'Cash In', onPress: () => console.log('Quick cash in') },
              { text: 'Cash Out', onPress: () => console.log('Quick cash out') },
              { text: 'Money Transfer', onPress: () => console.log('Quick transfer') },
              { text: 'Cancel', style: 'cancel' },
            ]
          );
        }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    backgroundColor: colors.primary,
    paddingTop: 20,
    paddingBottom: 20,
    paddingHorizontal: spacing.md,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: 'white',
  },
  headerActions: {
    flexDirection: 'row',
  },
  locationCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  locationLoading: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  locationLoadingText: {
    fontSize: 14,
    color: colors.text,
    marginLeft: spacing.sm,
  },
  locationProgress: {
    backgroundColor: colors.surface,
  },
  tabButtons: {
    margin: spacing.md,
    backgroundColor: colors.surface,
  },
  searchContainer: {
    paddingHorizontal: spacing.md,
  },
  searchBar: {
    backgroundColor: colors.surface,
    marginBottom: spacing.md,
  },
  filterChips: {
    flexDirection: 'row',
    gap: spacing.sm,
    marginBottom: spacing.md,
  },
  filterChip: {
    backgroundColor: colors.surface,
  },
  // Seva Mitras List
  mitrasList: {
    paddingHorizontal: spacing.md,
    gap: spacing.md,
  },
  mitraCard: {
    backgroundColor: colors.card,
    elevation: 2,
  },
  mitraHeader: {
    flexDirection: 'row',
    marginBottom: spacing.md,
  },
  mitraAvatarContainer: {
    position: 'relative',
    marginRight: spacing.md,
  },
  onlineBadge: {
    position: 'absolute',
    top: -2,
    right: -2,
    backgroundColor: colors.success,
  },
  verificationIcon: {
    position: 'absolute',
    bottom: -2,
    right: -2,
    backgroundColor: colors.background,
    borderRadius: 8,
  },
  mitraInfo: {
    flex: 1,
  },
  mitraName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
    marginBottom: spacing.xs,
  },
  trustChip: {
    alignSelf: 'flex-start',
    height: 20,
    marginBottom: spacing.xs,
  },
  mitraStats: {
    fontSize: 12,
    color: colors.textSecondary,
  },
  mitraDistance: {
    alignItems: 'flex-end',
  },
  distanceChip: {
    backgroundColor: colors.surface,
    height: 24,
    marginBottom: spacing.xs,
  },
  lastActive: {
    fontSize: 10,
    color: colors.textSecondary,
  },
  mitraDivider: {
    marginVertical: spacing.md,
  },
  mitraDetails: {
    marginBottom: spacing.md,
  },
  locationInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  locationText: {
    fontSize: 12,
    color: colors.textSecondary,
    marginLeft: spacing.xs,
    flex: 1,
  },
  servicesInfo: {
    marginBottom: spacing.sm,
  },
  servicesLabel: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  servicesList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.xs,
  },
  serviceChip: {
    backgroundColor: colors.surface,
    height: 24,
  },
  limitsInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.sm,
  },
  limitItem: {
    flex: 1,
  },
  limitLabel: {
    fontSize: 10,
    color: colors.textSecondary,
  },
  limitValue: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.text,
  },
  workingHours: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  workingHoursText: {
    fontSize: 12,
    color: colors.textSecondary,
    marginLeft: spacing.xs,
    flex: 1,
  },
  languagesList: {
    flexDirection: 'row',
    gap: spacing.xs,
  },
  languageChip: {
    backgroundColor: colors.surface,
    height: 20,
  },
  mitraActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: spacing.sm,
  },
  mitraActionButton: {
    flex: 1,
  },
  // Service Requests
  requestsList: {
    paddingHorizontal: spacing.md,
    gap: spacing.md,
  },
  requestCard: {
    backgroundColor: colors.card,
    elevation: 2,
  },
  requestHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  requestTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: colors.text,
  },
  requestAmount: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.primary,
    marginTop: spacing.xs,
  },
  requestStatusChip: {
    height: 24,
  },
  requestDetails: {
    marginBottom: spacing.md,
  },
  requestInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.sm,
  },
  requestInfoLabel: {
    fontSize: 12,
    color: colors.textSecondary,
    flex: 1,
  },
  requestInfoValue: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.text,
    flex: 2,
    textAlign: 'right',
  },
  requestActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: spacing.md,
  },
  requestActionButton: {
    flex: 0.4,
  },
  // Empty States
  emptyCard: {
    margin: spacing.md,
    backgroundColor: colors.card,
  },
  emptyContent: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
  },
  emptyText: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.textSecondary,
    marginTop: spacing.md,
  },
  emptySubtext: {
    fontSize: 14,
    color: colors.textSecondary,
    marginTop: spacing.sm,
    textAlign: 'center',
  },
  fab: {
    position: 'absolute',
    margin: 16,
    right: 0,
    bottom: 0,
    backgroundColor: colors.primary,
  },
  bottomSpacer: {
    height: 100,
  },
});