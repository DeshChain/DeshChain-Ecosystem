import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  Avatar,
  Button,
  IconButton,
  Grid,
  List,
  ListItem,
  Divider,
  Rating,
  Tooltip,
  Badge,
  Dialog,
  DialogContent,
  DialogTitle,
  DialogActions,
  Collapse,
} from '@mui/material';
import {
  LocationOn,
  Phone,
  Star,
  Verified,
  Schedule,
  Navigation,
  Language,
  Payment,
  ExpandMore,
  ExpandLess,
  Info,
  ChatBubbleOutline,
  BookmarkBorder,
  Share,
} from '@mui/icons-material';
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
  commissionRate?: string;
}

interface OperatingHours {
  day: string;
  openTime: string;
  closeTime: string;
  isClosed: boolean;
}

interface ListViewProps {
  sevaMitras: SevaMitra[];
  onSevaMitraSelect?: (mitra: SevaMitra) => void;
  selectedMitraId?: string;
  loading?: boolean;
}

const StyledCard = styled(Card)(({ theme }) => ({
  marginBottom: theme.spacing(2),
  borderRadius: theme.spacing(2),
  transition: 'all 0.2s ease-in-out',
  cursor: 'pointer',
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: '0 8px 25px rgba(0,0,0,0.12)',
  },
  '&.selected': {
    borderColor: theme.palette.primary.main,
    borderWidth: 2,
    boxShadow: `0 0 0 2px ${theme.palette.primary.main}20`,
  },
}));

const TrustScoreChip = styled(Chip)<{ trustscore: number }>(({ theme, trustscore }) => ({
  fontWeight: 'bold',
  color: theme.palette.getContrastText(
    trustscore >= 90 ? '#9c27b0' : // Diamond
    trustscore >= 80 ? '#607d8b' : // Platinum
    trustscore >= 70 ? '#ff9800' : // Gold
    trustscore >= 60 ? '#9e9e9e' : // Silver
    trustscore >= 50 ? '#795548' : // Bronze
    '#f44336' // New User
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

const StatusBadge = styled(Badge)(({ theme }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: '#44b700',
    color: '#44b700',
    boxShadow: `0 0 0 2px ${theme.palette.background.paper}`,
    '&::after': {
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      borderRadius: '50%',
      animation: 'ripple 1.2s infinite ease-in-out',
      border: '1px solid currentColor',
      content: '""',
    },
  },
  '@keyframes ripple': {
    '0%': {
      transform: 'scale(.8)',
      opacity: 1,
    },
    '100%': {
      transform: 'scale(2.4)',
      opacity: 0,
    },
  },
}));

export const SevaMitraListView: React.FC<ListViewProps> = ({
  sevaMitras,
  onSevaMitraSelect,
  selectedMitraId,
  loading = false,
}) => {
  const [expandedDetails, setExpandedDetails] = useState<string[]>([]);
  const [detailsDialog, setDetailsDialog] = useState<SevaMitra | null>(null);

  const toggleExpanded = (mitraId: string, event: React.MouseEvent) => {
    event.stopPropagation();
    setExpandedDetails(prev =>
      prev.includes(mitraId)
        ? prev.filter(id => id !== mitraId)
        : [...prev, mitraId]
    );
  };

  const handleSevaMitraClick = (mitra: SevaMitra) => {
    onSevaMitraSelect?.(mitra);
  };

  const handleCall = (phone: string, event: React.MouseEvent) => {
    event.stopPropagation();
    window.open(`tel:${phone}`, '_self');
  };

  const handleDirections = (mitra: SevaMitra, event: React.MouseEvent) => {
    event.stopPropagation();
    const url = `https://www.google.com/maps/dir/?api=1&destination=${mitra.latitude},${mitra.longitude}`;
    window.open(url, '_blank');
  };

  const handleDetailsDialog = (mitra: SevaMitra, event: React.MouseEvent) => {
    event.stopPropagation();
    setDetailsDialog(mitra);
  };

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

  const getServiceIcon = (service: string) => {
    const serviceIcons: { [key: string]: string } = {
      'CASH_IN': '💰',
      'CASH_OUT': '💸',
      'REMITTANCE': '📤',
      'BILL_PAYMENT': '🧾',
    };
    return serviceIcons[service] || '⚙️';
  };

  const getCurrentDayHours = (operatingHours: OperatingHours[]) => {
    const today = new Date().toLocaleDateString('en-US', { weekday: 'long' });
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

  const getWeekSchedule = (operatingHours: OperatingHours[]) => {
    const daysOrder = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
    return daysOrder.map(day => {
      const hours = operatingHours.find(h => h.day === day);
      return {
        day,
        hours: hours ? (hours.isClosed ? 'Closed' : `${hours.openTime} - ${hours.closeTime}`) : 'Not specified',
        isClosed: hours?.isClosed || false,
      };
    });
  };

  if (loading) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="h6" color="text.secondary">
          Loading Seva Mitras...
        </Typography>
      </Box>
    );
  }

  if (sevaMitras.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          No Seva Mitras found
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Try adjusting your filters or expanding the search radius
        </Typography>
      </Box>
    );
  }

  return (
    <>
      <List sx={{ p: 0 }}>
        {sevaMitras.map((mitra) => {
          const isExpanded = expandedDetails.includes(mitra.mitraId);
          const isSelected = selectedMitraId === mitra.mitraId;
          const isOpen = isCurrentlyOpen(mitra.operatingHours);

          return (
            <ListItem key={mitra.mitraId} sx={{ p: 0, mb: 2 }}>
              <StyledCard
                className={isSelected ? 'selected' : ''}
                onClick={() => handleSevaMitraClick(mitra)}
                sx={{ width: '100%' }}
              >
                <CardContent>
                  {/* Header */}
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                    <StatusBadge
                      overlap="circular"
                      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                      variant="dot"
                      invisible={!isOpen}
                    >
                      <Avatar
                        sx={{
                          width: 56,
                          height: 56,
                          mr: 2,
                          bgcolor: 'primary.main',
                          fontSize: '1.5rem',
                          fontWeight: 'bold',
                        }}
                      >
                        {mitra.businessName.charAt(0)}
                      </Avatar>
                    </StatusBadge>

                    <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                      <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 0.5 }}>
                        {mitra.businessName}
                      </Typography>

                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1, flexWrap: 'wrap' }}>
                        <TrustScoreChip
                          label={`${getTrustBadge(mitra.trustScore)} ${mitra.trustScore}`}
                          size="small"
                          trustscore={mitra.trustScore}
                        />
                        
                        <Box sx={{ display: 'flex', alignItems: 'center' }}>
                          <Rating value={mitra.averageRating} readOnly size="small" precision={0.1} />
                          <Typography variant="body2" sx={{ ml: 0.5 }}>
                            ({mitra.totalRatings})
                          </Typography>
                        </Box>

                        {mitra.isKYCVerified && (
                          <Tooltip title="KYC Verified">
                            <Verified color="success" sx={{ fontSize: 18 }} />
                          </Tooltip>
                        )}

                        <Chip
                          label={isOpen ? 'Open' : 'Closed'}
                          size="small"
                          color={isOpen ? 'success' : 'error'}
                          sx={{ fontSize: '0.7rem', height: '20px' }}
                        />
                      </Box>

                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        <LocationOn sx={{ fontSize: 14, mr: 0.5, verticalAlign: 'middle' }} />
                        {mitra.address}
                        {mitra.distance && (
                          <span style={{ marginLeft: 8 }}>
                            • {mitra.distance.toFixed(1)} km away
                          </span>
                        )}
                      </Typography>

                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        <Schedule sx={{ fontSize: 14, mr: 0.5, verticalAlign: 'middle' }} />
                        Today: {getCurrentDayHours(mitra.operatingHours)}
                        {mitra.responseTime && (
                          <span style={{ marginLeft: 8 }}>
                            • Response: {mitra.responseTime}
                          </span>
                        )}
                      </Typography>

                      {/* Services */}
                      <Box sx={{ mb: 1 }}>
                        {mitra.services.slice(0, 3).map((service) => (
                          <ServiceChip
                            key={service}
                            label={`${getServiceIcon(service)} ${getServiceDisplayName(service)}`}
                            size="small"
                            variant="outlined"
                            color="primary"
                          />
                        ))}
                        {mitra.services.length > 3 && (
                          <ServiceChip
                            label={`+${mitra.services.length - 3} more`}
                            size="small"
                            variant="outlined"
                          />
                        )}
                      </Box>
                    </Box>

                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                      <Tooltip title="View details">
                        <IconButton
                          size="small"
                          onClick={(e) => handleDetailsDialog(mitra, e)}
                        >
                          <Info />
                        </IconButton>
                      </Tooltip>
                      
                      <Tooltip title={isExpanded ? 'Collapse' : 'Expand'}>
                        <IconButton
                          size="small"
                          onClick={(e) => toggleExpanded(mitra.mitraId, e)}
                        >
                          {isExpanded ? <ExpandLess /> : <ExpandMore />}
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </Box>

                  {/* Action Buttons */}
                  <Box sx={{ display: 'flex', gap: 1, mb: isExpanded ? 2 : 0 }}>
                    <Button
                      variant="outlined"
                      size="small"
                      startIcon={<Phone />}
                      onClick={(e) => handleCall(mitra.phone, e)}
                      sx={{ flex: 1 }}
                    >
                      Call
                    </Button>
                    <Button
                      variant="contained"
                      size="small"
                      startIcon={<Navigation />}
                      onClick={(e) => handleDirections(mitra, e)}
                      sx={{ flex: 1 }}
                    >
                      Directions
                    </Button>
                    <IconButton size="small" onClick={(e) => e.stopPropagation()}>
                      <BookmarkBorder />
                    </IconButton>
                    <IconButton size="small" onClick={(e) => e.stopPropagation()}>
                      <Share />
                    </IconButton>
                  </Box>

                  {/* Expanded Details */}
                  <Collapse in={isExpanded}>
                    <Divider sx={{ mb: 2 }} />
                    
                    <Grid container spacing={2}>
                      <Grid item xs={12} sm={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          <Language sx={{ fontSize: 14, mr: 0.5, verticalAlign: 'middle' }} />
                          Languages
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {mitra.languages.join(', ')}
                        </Typography>
                      </Grid>

                      <Grid item xs={12} sm={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Contact
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          📞 {mitra.phone}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          ✉️ {mitra.email}
                        </Typography>
                      </Grid>

                      <Grid item xs={12}>
                        <Typography variant="subtitle2" gutterBottom>
                          All Services
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {mitra.services.map((service) => (
                            <ServiceChip
                              key={service}
                              label={`${getServiceIcon(service)} ${getServiceDisplayName(service)}`}
                              size="small"
                              variant="outlined"
                              color="primary"
                            />
                          ))}
                        </Box>
                      </Grid>

                      {mitra.commissionRate && (
                        <Grid item xs={12}>
                          <Typography variant="subtitle2" gutterBottom>
                            Commission Rate
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            {mitra.commissionRate}% per transaction
                          </Typography>
                        </Grid>
                      )}
                    </Grid>
                  </Collapse>
                </CardContent>
              </StyledCard>
            </ListItem>
          );
        })}
      </List>

      {/* Details Dialog */}
      <Dialog
        open={!!detailsDialog}
        onClose={() => setDetailsDialog(null)}
        maxWidth="sm"
        fullWidth
      >
        {detailsDialog && (
          <>
            <DialogTitle>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Avatar sx={{ mr: 2, bgcolor: 'primary.main' }}>
                  {detailsDialog.businessName.charAt(0)}
                </Avatar>
                <Box>
                  <Typography variant="h6">{detailsDialog.businessName}</Typography>
                  <Typography variant="body2" color="text.secondary">
                    {detailsDialog.address}
                  </Typography>
                </Box>
              </Box>
            </DialogTitle>
            
            <DialogContent>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <Typography variant="subtitle2" gutterBottom>
                    Weekly Schedule
                  </Typography>
                  {getWeekSchedule(detailsDialog.operatingHours).map((schedule) => (
                    <Box key={schedule.day} sx={{ display: 'flex', justifyContent: 'space-between', py: 0.5 }}>
                      <Typography variant="body2" sx={{ fontWeight: schedule.day === new Date().toLocaleDateString('en-US', { weekday: 'long' }) ? 'bold' : 'normal' }}>
                        {schedule.day}
                      </Typography>
                      <Typography variant="body2" color={schedule.isClosed ? 'error' : 'text.secondary'}>
                        {schedule.hours}
                      </Typography>
                    </Box>
                  ))}
                </Grid>
              </Grid>
            </DialogContent>
            
            <DialogActions>
              <Button onClick={() => setDetailsDialog(null)}>Close</Button>
              <Button
                variant="contained"
                startIcon={<ChatBubbleOutline />}
                onClick={() => {
                  // Handle chat/contact
                  setDetailsDialog(null);
                }}
              >
                Start Chat
              </Button>
            </DialogActions>
          </>
        )}
      </Dialog>
    </>
  );
};

export default SevaMitraListView;