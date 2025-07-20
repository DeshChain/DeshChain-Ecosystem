import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { GoogleMap, useJsApiLoader, Marker, InfoWindow, Circle } from '@react-google-maps/api';
import { Box, Card, CardContent, Typography, Chip, Avatar, Button, IconButton, Tooltip } from '@mui/material';
import { LocationOn, Star, Verified, Phone, Schedule, Navigation, FilterList } from '@mui/icons-material';
import { styled } from '@mui/material/styles';

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
}

interface OperatingHours {
  day: string;
  openTime: string;
  closeTime: string;
  isClosed: boolean;
}

interface MapViewProps {
  userLocation?: { lat: number; lng: number };
  selectedServices?: string[];
  maxDistance?: number;
  minTrustScore?: number;
  onSevaMitraSelect?: (mitra: SevaMitra) => void;
  showCurrentLocation?: boolean;
}

const mapContainerStyle = {
  width: '100%',
  height: '600px',
  borderRadius: '12px',
};

const defaultCenter = {
  lat: 28.6139, // Delhi
  lng: 77.2090,
};

const StyledCard = styled(Card)(({ theme }) => ({
  minWidth: 300,
  maxWidth: 350,
  '& .MuiCardContent-root': {
    padding: theme.spacing(2),
  },
}));

const TrustScoreChip = styled(Chip)<{ trustscore: number }>(({ theme, trustscore }) => ({
  fontWeight: 'bold',
  color: theme.palette.getContrastText(
    trustscore >= 90 ? '#9c27b0' : // Diamond - Purple
    trustscore >= 80 ? '#607d8b' : // Platinum - Blue Grey  
    trustscore >= 70 ? '#ff9800' : // Gold - Orange
    trustscore >= 60 ? '#9e9e9e' : // Silver - Grey
    trustscore >= 50 ? '#795548' : // Bronze - Brown
    '#f44336' // New User - Red
  ),
  backgroundColor:
    trustscore >= 90 ? '#9c27b0' :
    trustscore >= 80 ? '#607d8b' :
    trustscore >= 70 ? '#ff9800' :
    trustscore >= 60 ? '#9e9e9e' :
    trustscore >= 50 ? '#795548' :
    '#f44336',
}));

const ServiceChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.25),
  fontSize: '0.75rem',
  height: '24px',
}));

const libraries: ("places" | "geometry")[] = ["places", "geometry"];

export const SevaMitraMapView: React.FC<MapViewProps> = ({
  userLocation,
  selectedServices = [],
  maxDistance = 25,
  minTrustScore = 50,
  onSevaMitraSelect,
  showCurrentLocation = true,
}) => {
  const [sevaMitras, setSevaMitras] = useState<SevaMitra[]>([]);
  const [selectedMitra, setSelectedMitra] = useState<SevaMitra | null>(null);
  const [mapCenter, setMapCenter] = useState(userLocation || defaultCenter);
  const [mapZoom, setMapZoom] = useState(12);
  const [loading, setLoading] = useState(true);

  const { isLoaded, loadError } = useJsApiLoader({
    id: 'google-map-script',
    googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || '',
    libraries,
  });

  // Fetch Seva Mitras based on filters
  const fetchSevaMitras = useCallback(async () => {
    setLoading(true);
    try {
      // Mock data - replace with actual API call
      const mockData: SevaMitra[] = [
        {
          mitraId: 'SM-001',
          businessName: 'राम फाइनेंशियल सर्विसेज',
          address: 'Shop 15, Main Market, Connaught Place',
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
            { day: 'Sunday', openTime: '10:00', closeTime: '16:00', isClosed: false },
          ],
          distance: 2.3,
          responseTime: '< 5 min',
        },
        {
          mitraId: 'SM-002',
          businessName: 'Sharma Money Transfer',
          address: '42, Karol Bagh Metro Station',
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
            { day: 'Sunday', openTime: '11:00', closeTime: '17:00', isClosed: false },
          ],
          distance: 5.7,
          responseTime: '< 10 min',
        },
        {
          mitraId: 'SM-003',
          businessName: 'पटेल कैश पॉइंट',
          address: 'Near Bus Stand, Lajpat Nagar',
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
            { day: 'Sunday', openTime: '09:00', closeTime: '18:00', isClosed: false },
          ],
          distance: 8.2,
          responseTime: '< 15 min',
        },
      ];

      // Filter based on criteria
      const filtered = mockData.filter(mitra => {
        const serviceMatch = selectedServices.length === 0 || 
          selectedServices.some(service => mitra.services.includes(service));
        const trustMatch = mitra.trustScore >= minTrustScore;
        const distanceMatch = !mitra.distance || mitra.distance <= maxDistance;
        
        return serviceMatch && trustMatch && distanceMatch && mitra.isActive;
      });

      setSevaMitras(filtered);
    } catch (error) {
      console.error('Error fetching Seva Mitras:', error);
    } finally {
      setLoading(false);
    }
  }, [selectedServices, maxDistance, minTrustScore]);

  useEffect(() => {
    fetchSevaMitras();
  }, [fetchSevaMitras]);

  useEffect(() => {
    if (userLocation) {
      setMapCenter(userLocation);
    }
  }, [userLocation]);

  const getMarkerIcon = useCallback((mitra: SevaMitra) => {
    const color = mitra.trustScore >= 90 ? '#9c27b0' :
                 mitra.trustScore >= 80 ? '#607d8b' :
                 mitra.trustScore >= 70 ? '#ff9800' :
                 mitra.trustScore >= 60 ? '#9e9e9e' :
                 mitra.trustScore >= 50 ? '#795548' : '#f44336';
    
    return {
      path: google.maps.SymbolPath.CIRCLE,
      fillColor: color,
      fillOpacity: 0.8,
      strokeColor: '#ffffff',
      strokeWeight: 2,
      scale: 8,
    };
  }, []);

  const getTrustBadge = (trustScore: number) => {
    if (trustScore >= 90) return 'Diamond';
    if (trustScore >= 80) return 'Platinum';
    if (trustScore >= 70) return 'Gold';
    if (trustScore >= 60) return 'Silver';
    if (trustScore >= 50) return 'Bronze';
    return 'New';
  };

  const getServiceDisplayName = (service: string) => {
    const serviceNames: { [key: string]: string } = {
      'CASH_IN': 'Cash In',
      'CASH_OUT': 'Cash Out',
      'REMITTANCE': 'Money Transfer',
      'BILL_PAYMENT': 'Bill Payment',
    };
    return serviceNames[service] || service;
  };

  const getCurrentDayHours = (operatingHours: OperatingHours[]) => {
    const today = new Date().toLocaleLString('en-US', { weekday: 'long' });
    const todayHours = operatingHours.find(h => h.day === today);
    
    if (!todayHours || todayHours.isClosed) {
      return 'Closed today';
    }
    
    return `${todayHours.openTime} - ${todayHours.closeTime}`;
  };

  const isCurrentlyOpen = (operatingHours: OperatingHours[]) => {
    const now = new Date();
    const today = now.toLocaleDateString('en-US', { weekday: 'long' });
    const currentTime = now.toTimeString().slice(0, 5);
    
    const todayHours = operatingHours.find(h => h.day === today);
    if (!todayHours || todayHours.isClosed) return false;
    
    return currentTime >= todayHours.openTime && currentTime <= todayHours.closeTime;
  };

  const handleMarkerClick = (mitra: SevaMitra) => {
    setSelectedMitra(mitra);
    onSevaMitraSelect?.(mitra);
  };

  const handleDirections = (mitra: SevaMitra) => {
    const url = `https://www.google.com/maps/dir/?api=1&destination=${mitra.latitude},${mitra.longitude}`;
    window.open(url, '_blank');
  };

  const handleCall = (phone: string) => {
    window.open(`tel:${phone}`, '_self');
  };

  if (loadError) {
    return (
      <Box sx={{ p: 2, textAlign: 'center' }}>
        <Typography variant="h6" color="error">
          Error loading maps. Please check your internet connection.
        </Typography>
      </Box>
    );
  }

  if (!isLoaded) {
    return (
      <Box sx={{ p: 2, textAlign: 'center' }}>
        <Typography variant="h6">Loading map...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ position: 'relative', height: '600px' }}>
      <GoogleMap
        mapContainerStyle={mapContainerStyle}
        center={mapCenter}
        zoom={mapZoom}
        options={{
          zoomControl: true,
          streetViewControl: false,
          mapTypeControl: false,
          fullscreenControl: true,
        }}
      >
        {/* User's current location */}
        {showCurrentLocation && userLocation && (
          <>
            <Marker
              position={userLocation}
              icon={{
                path: google.maps.SymbolPath.CIRCLE,
                fillColor: '#4285f4',
                fillOpacity: 1,
                strokeColor: '#ffffff',
                strokeWeight: 3,
                scale: 10,
              }}
              title="Your Location"
            />
            <Circle
              center={userLocation}
              radius={maxDistance * 1000} // Convert km to meters
              options={{
                fillColor: '#4285f4',
                fillOpacity: 0.1,
                strokeColor: '#4285f4',
                strokeOpacity: 0.3,
                strokeWeight: 1,
              }}
            />
          </>
        )}

        {/* Seva Mitra markers */}
        {sevaMitras.map((mitra) => (
          <Marker
            key={mitra.mitraId}
            position={{ lat: mitra.latitude, lng: mitra.longitude }}
            icon={getMarkerIcon(mitra)}
            title={mitra.businessName}
            onClick={() => handleMarkerClick(mitra)}
          />
        ))}

        {/* Info Window for selected Seva Mitra */}
        {selectedMitra && (
          <InfoWindow
            position={{ lat: selectedMitra.latitude, lng: selectedMitra.longitude }}
            onCloseClick={() => setSelectedMitra(null)}
          >
            <StyledCard elevation={0}>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                  <Avatar sx={{ width: 40, height: 40, mr: 2, bgcolor: 'primary.main' }}>
                    {selectedMitra.businessName.charAt(0)}
                  </Avatar>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="h6" sx={{ fontSize: '1rem', fontWeight: 'bold' }}>
                      {selectedMitra.businessName}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <TrustScoreChip
                        label={`${getTrustBadge(selectedMitra.trustScore)} ${selectedMitra.trustScore}`}
                        size="small"
                        trustscore={selectedMitra.trustScore}
                      />
                      {selectedMitra.isKYCVerified && (
                        <Tooltip title="KYC Verified">
                          <Verified color="success" sx={{ fontSize: 16 }} />
                        </Tooltip>
                      )}
                    </Box>
                  </Box>
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                  <Star sx={{ color: '#ff9800', fontSize: 16, mr: 0.5 }} />
                  <Typography variant="body2" sx={{ mr: 1 }}>
                    {selectedMitra.averageRating} ({selectedMitra.totalRatings} reviews)
                  </Typography>
                  <Chip
                    label={isCurrentlyOpen(selectedMitra.operatingHours) ? 'Open' : 'Closed'}
                    size="small"
                    color={isCurrentlyOpen(selectedMitra.operatingHours) ? 'success' : 'error'}
                    sx={{ fontSize: '0.7rem', height: '20px' }}
                  />
                </Box>

                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  <LocationOn sx={{ fontSize: 14, mr: 0.5, verticalAlign: 'middle' }} />
                  {selectedMitra.address}
                </Typography>

                {selectedMitra.distance && (
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                    📍 {selectedMitra.distance.toFixed(1)} km away • Response: {selectedMitra.responseTime}
                  </Typography>
                )}

                <Typography variant="body2" sx={{ mb: 1 }}>
                  <Schedule sx={{ fontSize: 14, mr: 0.5, verticalAlign: 'middle' }} />
                  Today: {getCurrentDayHours(selectedMitra.operatingHours)}
                </Typography>

                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                    Services:
                  </Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                    {selectedMitra.services.map((service) => (
                      <ServiceChip
                        key={service}
                        label={getServiceDisplayName(service)}
                        size="small"
                        variant="outlined"
                        color="primary"
                      />
                    ))}
                  </Box>
                </Box>

                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                    Languages:
                  </Typography>
                  <Typography variant="body2">
                    {selectedMitra.languages.join(', ')}
                  </Typography>
                </Box>

                <Box sx={{ display: 'flex', gap: 1, justifyContent: 'space-between' }}>
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<Phone />}
                    onClick={() => handleCall(selectedMitra.phone)}
                    sx={{ flex: 1 }}
                  >
                    Call
                  </Button>
                  <Button
                    variant="contained"
                    size="small"
                    startIcon={<Navigation />}
                    onClick={() => handleDirections(selectedMitra)}
                    sx={{ flex: 1 }}
                  >
                    Directions
                  </Button>
                </Box>
              </CardContent>
            </StyledCard>
          </InfoWindow>
        )}
      </GoogleMap>

      {/* Loading overlay */}
      {loading && (
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'rgba(255, 255, 255, 0.8)',
            borderRadius: '12px',
          }}
        >
          <Typography variant="h6">Loading Seva Mitras...</Typography>
        </Box>
      )}
    </Box>
  );
};

export default SevaMitraMapView;