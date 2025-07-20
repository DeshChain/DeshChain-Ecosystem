import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Container,
  Grid,
  Typography,
  Tabs,
  Tab,
  Paper,
  IconButton,
  Tooltip,
  Fab,
  Snackbar,
  Alert,
  Switch,
  FormControlLabel,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Map,
  ViewList,
  MyLocation,
  Refresh,
  Settings,
  FilterList,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

import SevaMitraMapView from './MapView';
import FilterPanel from './FilterPanel';
import SevaMitraListView from './ListView';

// Types
interface SevaMitra {
  mitraId: string;
  businessName: string;
  address: string;
  postalCode: string;
  district: string;
  state: string;
  latitude: number;
  longitude: number;
  phone: string;
  email: string;
  languages: string[];
  services: string[];
  trustScore: number;
  averageRating: number;
  totalRatings: number;
  isActive: boolean;
  isKYCVerified: boolean;
  operatingHours: OperatingHours[];
  distance?: number;
  responseTime?: string;
  commissionRate?: string;
}

interface OperatingHours {
  day: string;
  openTime: string;
  closeTime: string;
  isClosed: boolean;
}

interface FilterOptions {
  services: string[];
  maxDistance: number;
  minTrustScore: number;
  languages: string[];
  isKYCRequired: boolean;
  isCurrentlyOpen: boolean;
  paymentMethods: string[];
  sortBy: 'distance' | 'trustScore' | 'rating' | 'responseTime';
  postalCode: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const StyledPaper = styled(Paper)(({ theme }) => ({
  borderRadius: theme.spacing(2),
  boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
  overflow: 'hidden',
}));

const FloatingActionButton = styled(Fab)(({ theme }) => ({
  position: 'fixed',
  bottom: theme.spacing(3),
  right: theme.spacing(3),
  zIndex: 1000,
}));

const StyledTabs = styled(Tabs)(({ theme }) => ({
  borderBottom: `1px solid ${theme.palette.divider}`,
  '& .MuiTabs-indicator': {
    height: 3,
    borderRadius: '3px 3px 0 0',
  },
}));

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`seva-mitra-tabpanel-${index}`}
      aria-labelledby={`seva-mitra-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 2 }}>{children}</Box>}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `seva-mitra-tab-${index}`,
    'aria-controls': `seva-mitra-tabpanel-${index}`,
  };
}

export const SevaMitraDiscovery: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  const [activeTab, setActiveTab] = useState(0);
  const [selectedMitra, setSelectedMitra] = useState<SevaMitra | null>(null);
  const [userLocation, setUserLocation] = useState<{ lat: number; lng: number } | undefined>();
  const [sevaMitras, setSevaMitras] = useState<SevaMitra[]>([]);
  const [loading, setLoading] = useState(false);
  const [showFilters, setShowFilters] = useState(!isMobile);
  const [locationError, setLocationError] = useState<string>('');
  const [autoRefresh, setAutoRefresh] = useState(false);
  
  const [filters, setFilters] = useState<FilterOptions>({
    services: [],
    maxDistance: 25,
    minTrustScore: 50,
    languages: [],
    isKYCRequired: false,
    isCurrentlyOpen: false,
    paymentMethods: [],
    sortBy: 'distance',
    postalCode: '',
  });

  // Request user's current location
  const requestLocation = useCallback(() => {
    setLocationError('');
    
    if (!navigator.geolocation) {
      setLocationError('Geolocation is not supported by this browser');
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        setUserLocation({ lat: latitude, lng: longitude });
        
        // Auto-detect postal code based on location
        // In production, would use reverse geocoding API
        detectPostalCode(latitude, longitude);
      },
      (error) => {
        let errorMessage = 'Unable to get location';
        
        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = 'Location permission denied';
            break;
          case error.POSITION_UNAVAILABLE:
            errorMessage = 'Location information unavailable';
            break;
          case error.TIMEOUT:
            errorMessage = 'Location request timed out';
            break;
        }
        
        setLocationError(errorMessage);
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000, // 5 minutes
      }
    );
  }, []);

  // Detect postal code from coordinates (mock implementation)
  const detectPostalCode = async (lat: number, lng: number) => {
    try {
      // In production, would use Google Geocoding API or similar
      // For now, use mock postal code based on Delhi coordinates
      if (lat >= 28.4 && lat <= 28.9 && lng >= 76.8 && lng <= 77.5) {
        setFilters(prev => ({ ...prev, postalCode: '110001' }));
      }
    } catch (error) {
      console.error('Error detecting postal code:', error);
    }
  };

  // Fetch Seva Mitras based on current filters
  const fetchSevaMitras = useCallback(async () => {
    setLoading(true);
    
    try {
      // Mock API call - replace with actual service
      const response = await mockFetchSevaMitras(filters, userLocation);
      setSevaMitras(response);
    } catch (error) {
      console.error('Error fetching Seva Mitras:', error);
      setSevaMitras([]);
    } finally {
      setLoading(false);
    }
  }, [filters, userLocation]);

  // Mock API function - replace with actual implementation
  const mockFetchSevaMitras = async (
    filters: FilterOptions,
    location?: { lat: number; lng: number }
  ): Promise<SevaMitra[]> => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock data with more comprehensive examples
    const mockData: SevaMitra[] = [
      {
        mitraId: 'SM-001',
        businessName: 'राम फाइनेंशियल सर्विसेज',
        address: 'Shop 15, Main Market, Connaught Place, New Delhi',
        postalCode: '110001',
        district: 'New Delhi',
        state: 'Delhi',
        latitude: 28.6304,
        longitude: 77.2177,
        phone: '+91 98765 43210',
        email: 'ram.financial@example.com',
        languages: ['Hindi', 'English', 'Punjabi'],
        services: ['CASH_IN', 'CASH_OUT', 'REMITTANCE'],
        trustScore: 92,
        averageRating: 4.8,
        totalRatings: 156,
        isActive: true,
        isKYCVerified: true,
        operatingHours: [
          { day: 'Monday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          { day: 'Tuesday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          { day: 'Wednesday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          { day: 'Thursday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          { day: 'Friday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          { day: 'Saturday', openTime: '10:00', closeTime: '17:00', isClosed: false },
          { day: 'Sunday', openTime: '10:00', closeTime: '16:00', isClosed: false },
        ],
        distance: 2.3,
        responseTime: '< 5 min',
        commissionRate: '1.8',
      },
      {
        mitraId: 'SM-002',
        businessName: 'Sharma Money Transfer',
        address: '42, Karol Bagh Metro Station, Karol Bagh, New Delhi',
        postalCode: '110005',
        district: 'New Delhi',
        state: 'Delhi',
        latitude: 28.6519,
        longitude: 77.1909,
        phone: '+91 98765 43211',
        email: 'sharma.money@example.com',
        languages: ['Hindi', 'English'],
        services: ['CASH_IN', 'REMITTANCE', 'BILL_PAYMENT'],
        trustScore: 78,
        averageRating: 4.2,
        totalRatings: 89,
        isActive: true,
        isKYCVerified: true,
        operatingHours: [
          { day: 'Monday', openTime: '10:00', closeTime: '19:00', isClosed: false },
          { day: 'Tuesday', openTime: '10:00', closeTime: '19:00', isClosed: false },
          { day: 'Wednesday', openTime: '10:00', closeTime: '19:00', isClosed: false },
          { day: 'Thursday', openTime: '10:00', closeTime: '19:00', isClosed: false },
          { day: 'Friday', openTime: '10:00', closeTime: '19:00', isClosed: false },
          { day: 'Saturday', openTime: '11:00', closeTime: '18:00', isClosed: false },
          { day: 'Sunday', openTime: '11:00', closeTime: '17:00', isClosed: false },
        ],
        distance: 5.7,
        responseTime: '< 10 min',
        commissionRate: '2.2',
      },
      {
        mitraId: 'SM-003',
        businessName: 'पटेल कैश पॉइंट',
        address: 'Near Bus Stand, Lajpat Nagar, New Delhi',
        postalCode: '110024',
        district: 'New Delhi',
        state: 'Delhi',
        latitude: 28.5676,
        longitude: 77.2436,
        phone: '+91 98765 43212',
        email: 'patel.cash@example.com',
        languages: ['Hindi', 'Gujarati', 'English'],
        services: ['CASH_OUT', 'REMITTANCE'],
        trustScore: 85,
        averageRating: 4.5,
        totalRatings: 203,
        isActive: true,
        isKYCVerified: true,
        operatingHours: [
          { day: 'Monday', openTime: '08:30', closeTime: '20:00', isClosed: false },
          { day: 'Tuesday', openTime: '08:30', closeTime: '20:00', isClosed: false },
          { day: 'Wednesday', openTime: '08:30', closeTime: '20:00', isClosed: false },
          { day: 'Thursday', openTime: '08:30', closeTime: '20:00', isClosed: false },
          { day: 'Friday', openTime: '08:30', closeTime: '20:00', isClosed: false },
          { day: 'Saturday', openTime: '09:00', closeTime: '19:00', isClosed: false },
          { day: 'Sunday', openTime: '09:00', closeTime: '18:00', isClosed: false },
        ],
        distance: 8.2,
        responseTime: '< 15 min',
        commissionRate: '2.0',
      },
    ];

    // Apply filters
    return mockData.filter(mitra => {
      // Service filter
      if (filters.services.length > 0) {
        const hasService = filters.services.some(service => mitra.services.includes(service));
        if (!hasService) return false;
      }

      // Trust score filter
      if (mitra.trustScore < filters.minTrustScore) return false;

      // KYC filter
      if (filters.isKYCRequired && !mitra.isKYCVerified) return false;

      // Distance filter (if user location available)
      if (location && mitra.distance && mitra.distance > filters.maxDistance) return false;

      // Currently open filter
      if (filters.isCurrentlyOpen) {
        const isOpen = isCurrentlyOpenCheck(mitra.operatingHours);
        if (!isOpen) return false;
      }

      // Language filter
      if (filters.languages.length > 0) {
        const hasLanguage = filters.languages.some(lang => mitra.languages.includes(lang));
        if (!hasLanguage) return false;
      }

      // Postal code filter
      if (filters.postalCode && !mitra.postalCode.startsWith(filters.postalCode.slice(0, 3))) {
        return false;
      }

      return true;
    }).sort((a, b) => {
      // Sort by selected criteria
      switch (filters.sortBy) {
        case 'distance':
          return (a.distance || 999) - (b.distance || 999);
        case 'trustScore':
          return b.trustScore - a.trustScore;
        case 'rating':
          return b.averageRating - a.averageRating;
        case 'responseTime':
          // Simple sort by response time (would need proper parsing in real app)
          return (a.responseTime || 'z').localeCompare(b.responseTime || 'z');
        default:
          return 0;
      }
    });
  };

  const isCurrentlyOpenCheck = (operatingHours: OperatingHours[]) => {
    const now = new Date();
    const today = now.toLocaleDateString('en-US', { weekday: 'long' });
    const currentTime = now.toTimeString().slice(0, 5);
    
    const todayHours = operatingHours.find(h => h.day === today);
    if (!todayHours || todayHours.isClosed) return false;
    
    return currentTime >= todayHours.openTime && currentTime <= todayHours.closeTime;
  };

  // Effects
  useEffect(() => {
    fetchSevaMitras();
  }, [fetchSevaMitras]);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(fetchSevaMitras, 30000); // Refresh every 30 seconds
      return () => clearInterval(interval);
    }
  }, [autoRefresh, fetchSevaMitras]);

  // Handlers
  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const handleSevaMitraSelect = (mitra: SevaMitra) => {
    setSelectedMitra(mitra);
    // If on mobile and selecting from list, switch to map view
    if (isMobile && activeTab === 1) {
      setActiveTab(0);
    }
  };

  const handleFiltersChange = (newFilters: FilterOptions) => {
    setFilters(newFilters);
  };

  const getActiveFiltersCount = () => {
    let count = 0;
    if (filters.services.length > 0) count++;
    if (filters.minTrustScore > 50) count++;
    if (filters.maxDistance < 25) count++;
    if (filters.languages.length > 0) count++;
    if (filters.isKYCRequired) count++;
    if (filters.isCurrentlyOpen) count++;
    if (filters.paymentMethods.length > 0) count++;
    if (filters.postalCode) count++;
    return count;
  };

  return (
    <Container maxWidth="xl" sx={{ py: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
          Seva Mitra Discovery
        </Typography>
        <Typography variant="body1" color="text.secondary" gutterBottom>
          Find trusted community service friends near you for cash-in, cash-out, and money transfer services
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 2, flexWrap: 'wrap' }}>
          <FormControlLabel
            control={
              <Switch
                checked={showFilters}
                onChange={(e) => setShowFilters(e.target.checked)}
                size="small"
              />
            }
            label="Show Filters"
          />
          
          <FormControlLabel
            control={
              <Switch
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                size="small"
              />
            }
            label="Auto Refresh"
          />

          <Tooltip title="Get current location">
            <IconButton onClick={requestLocation} color="primary">
              <MyLocation />
            </IconButton>
          </Tooltip>

          <Tooltip title="Refresh data">
            <IconButton onClick={fetchSevaMitras} disabled={loading}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      <Grid container spacing={3}>
        {/* Filters Panel */}
        {showFilters && (
          <Grid item xs={12} md={3}>
            <FilterPanel
              filters={filters}
              onFiltersChange={handleFiltersChange}
              onLocationRequest={requestLocation}
              userLocation={userLocation}
              activeFiltersCount={getActiveFiltersCount()}
            />
          </Grid>
        )}

        {/* Main Content */}
        <Grid item xs={12} md={showFilters ? 9 : 12}>
          <StyledPaper>
            {/* Tabs */}
            <StyledTabs
              value={activeTab}
              onChange={handleTabChange}
              aria-label="seva mitra discovery tabs"
            >
              <Tab icon={<Map />} label="Map View" {...a11yProps(0)} />
              <Tab icon={<ViewList />} label={`List View (${sevaMitras.length})`} {...a11yProps(1)} />
            </StyledTabs>

            {/* Tab Panels */}
            <TabPanel value={activeTab} index={0}>
              <SevaMitraMapView
                userLocation={userLocation}
                selectedServices={filters.services}
                maxDistance={filters.maxDistance}
                minTrustScore={filters.minTrustScore}
                onSevaMitraSelect={handleSevaMitraSelect}
                showCurrentLocation={!!userLocation}
              />
            </TabPanel>

            <TabPanel value={activeTab} index={1}>
              <SevaMitraListView
                sevaMitras={sevaMitras}
                onSevaMitraSelect={handleSevaMitraSelect}
                selectedMitraId={selectedMitra?.mitraId}
                loading={loading}
              />
            </TabPanel>
          </StyledPaper>
        </Grid>
      </Grid>

      {/* Mobile Toggle Filters FAB */}
      {isMobile && (
        <FloatingActionButton
          color="primary"
          onClick={() => setShowFilters(!showFilters)}
          size="medium"
        >
          <FilterList />
        </FloatingActionButton>
      )}

      {/* Location Error Snackbar */}
      <Snackbar
        open={!!locationError}
        autoHideDuration={6000}
        onClose={() => setLocationError('')}
      >
        <Alert onClose={() => setLocationError('')} severity="warning">
          {locationError}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default SevaMitraDiscovery;