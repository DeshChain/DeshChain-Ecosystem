import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  FormGroup,
  FormControlLabel,
  Checkbox,
  Slider,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  Switch,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Button,
  Divider,
  IconButton,
  Tooltip,
  Badge,
} from '@mui/material';
import {
  ExpandMore,
  LocationOn,
  Star,
  Schedule,
  Language,
  Payment,
  Refresh,
  MyLocation,
  FilterList,
  Clear,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

// Types
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

interface FilterPanelProps {
  filters: FilterOptions;
  onFiltersChange: (filters: FilterOptions) => void;
  onLocationRequest: () => void;
  userLocation?: { lat: number; lng: number };
  activeFiltersCount?: number;
}

const StyledCard = styled(Card)(({ theme }) => ({
  borderRadius: theme.spacing(2),
  boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
}));

const FilterSection = styled(Box)(({ theme }) => ({
  marginBottom: theme.spacing(2),
  '&:last-child': {
    marginBottom: 0,
  },
}));

const StyledSlider = styled(Slider)(({ theme }) => ({
  '& .MuiSlider-thumb': {
    height: 20,
    width: 20,
  },
  '& .MuiSlider-valueLabel': {
    fontSize: 12,
    fontWeight: 'normal',
    top: -6,
    backgroundColor: 'unset',
    color: theme.palette.text.primary,
    '&:before': {
      display: 'none',
    },
    '& *': {
      background: 'transparent',
      color: theme.palette.mode === 'dark' ? '#fff' : '#000',
    },
  },
}));

const LanguageChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.25),
  fontSize: '0.75rem',
  height: '28px',
}));

const serviceOptions = [
  { value: 'CASH_IN', label: 'Cash In', icon: '💰' },
  { value: 'CASH_OUT', label: 'Cash Out', icon: '💸' },
  { value: 'REMITTANCE', label: 'Money Transfer', icon: '📤' },
  { value: 'BILL_PAYMENT', label: 'Bill Payment', icon: '🧾' },
];

const languageOptions = [
  'Hindi', 'English', 'Bengali', 'Telugu', 'Marathi', 'Tamil', 
  'Gujarati', 'Urdu', 'Kannada', 'Odia', 'Malayalam', 'Punjabi',
  'Assamese', 'Maithili', 'Santali', 'Kashmiri', 'Nepali', 'Konkani',
  'Sindhi', 'Dogri', 'Manipuri', 'Bodo'
];

const paymentMethodOptions = [
  { value: 'UPI', label: 'UPI', icon: '📱' },
  { value: 'IMPS', label: 'IMPS', icon: '🏦' },
  { value: 'NEFT', label: 'NEFT', icon: '🏛️' },
  { value: 'RTGS', label: 'RTGS', icon: '💼' },
  { value: 'CASH', label: 'Cash', icon: '💵' },
];

const trustScoreMarks = [
  { value: 0, label: '0' },
  { value: 25, label: '25' },
  { value: 50, label: '50' },
  { value: 75, label: '75' },
  { value: 100, label: '100' },
];

const distanceMarks = [
  { value: 1, label: '1km' },
  { value: 5, label: '5km' },
  { value: 10, label: '10km' },
  { value: 25, label: '25km' },
  { value: 50, label: '50km' },
];

export const FilterPanel: React.FC<FilterPanelProps> = ({
  filters,
  onFiltersChange,
  onLocationRequest,
  userLocation,
  activeFiltersCount = 0,
}) => {
  const [expandedSections, setExpandedSections] = useState<string[]>([
    'services', 'location', 'trust'
  ]);

  const handleSectionToggle = (section: string) => {
    setExpandedSections(prev =>
      prev.includes(section)
        ? prev.filter(s => s !== section)
        : [...prev, section]
    );
  };

  const handleServiceChange = (service: string, checked: boolean) => {
    const newServices = checked
      ? [...filters.services, service]
      : filters.services.filter(s => s !== service);
    
    onFiltersChange({ ...filters, services: newServices });
  };

  const handleLanguageChange = (language: string, checked: boolean) => {
    const newLanguages = checked
      ? [...filters.languages, language]
      : filters.languages.filter(l => l !== language);
    
    onFiltersChange({ ...filters, languages: newLanguages });
  };

  const handlePaymentMethodChange = (method: string, checked: boolean) => {
    const newMethods = checked
      ? [...filters.paymentMethods, method]
      : filters.paymentMethods.filter(m => m !== method);
    
    onFiltersChange({ ...filters, paymentMethods: newMethods });
  };

  const handleDistanceChange = (event: Event, newValue: number | number[]) => {
    onFiltersChange({ ...filters, maxDistance: newValue as number });
  };

  const handleTrustScoreChange = (event: Event, newValue: number | number[]) => {
    onFiltersChange({ ...filters, minTrustScore: newValue as number });
  };

  const handleResetFilters = () => {
    onFiltersChange({
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
  };

  const getTrustScoreBadge = (score: number) => {
    if (score >= 90) return { label: 'Diamond', color: '#9c27b0' };
    if (score >= 80) return { label: 'Platinum', color: '#607d8b' };
    if (score >= 70) return { label: 'Gold', color: '#ff9800' };
    if (score >= 60) return { label: 'Silver', color: '#9e9e9e' };
    if (score >= 50) return { label: 'Bronze', color: '#795548' };
    return { label: 'New', color: '#f44336' };
  };

  const currentTrustBadge = getTrustScoreBadge(filters.minTrustScore);

  return (
    <StyledCard>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <FilterList sx={{ mr: 1 }} />
            <Typography variant="h6" component="h2">
              Filters
            </Typography>
            {activeFiltersCount > 0 && (
              <Badge badgeContent={activeFiltersCount} color="primary" sx={{ ml: 1 }}>
                <Chip label="Active" size="small" color="primary" />
              </Badge>
            )}
          </Box>
          <Tooltip title="Reset all filters">
            <IconButton onClick={handleResetFilters} size="small">
              <Clear />
            </IconButton>
          </Tooltip>
        </Box>

        {/* Services Filter */}
        <Accordion 
          expanded={expandedSections.includes('services')}
          onChange={() => handleSectionToggle('services')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider', mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <Payment sx={{ mr: 1, fontSize: 20 }} />
              Services ({filters.services.length})
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <FormGroup>
              {serviceOptions.map((service) => (
                <FormControlLabel
                  key={service.value}
                  control={
                    <Checkbox
                      checked={filters.services.includes(service.value)}
                      onChange={(e) => handleServiceChange(service.value, e.target.checked)}
                      size="small"
                    />
                  }
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      <span style={{ marginRight: 8 }}>{service.icon}</span>
                      {service.label}
                    </Box>
                  }
                />
              ))}
            </FormGroup>
          </AccordionDetails>
        </Accordion>

        {/* Location Filter */}
        <Accordion 
          expanded={expandedSections.includes('location')}
          onChange={() => handleSectionToggle('location')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider', mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <LocationOn sx={{ mr: 1, fontSize: 20 }} />
              Location & Distance
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <FilterSection>
              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <TextField
                  label="Postal Code"
                  value={filters.postalCode}
                  onChange={(e) => onFiltersChange({ ...filters, postalCode: e.target.value })}
                  size="small"
                  placeholder="e.g., 110001"
                  sx={{ flex: 1 }}
                />
                <Tooltip title="Use my current location">
                  <IconButton 
                    onClick={onLocationRequest}
                    color={userLocation ? 'primary' : 'default'}
                  >
                    <MyLocation />
                  </IconButton>
                </Tooltip>
              </Box>
              
              <Typography variant="body2" gutterBottom>
                Maximum Distance: {filters.maxDistance} km
              </Typography>
              <StyledSlider
                value={filters.maxDistance}
                onChange={handleDistanceChange}
                min={1}
                max={50}
                marks={distanceMarks}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}km`}
              />
            </FilterSection>
          </AccordionDetails>
        </Accordion>

        {/* Trust Score Filter */}
        <Accordion 
          expanded={expandedSections.includes('trust')}
          onChange={() => handleSectionToggle('trust')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider', mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <Star sx={{ mr: 1, fontSize: 20 }} />
              Trust & Ratings
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <FilterSection>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <Typography variant="body2">
                  Minimum Trust Score: {filters.minTrustScore}
                </Typography>
                <Chip
                  label={currentTrustBadge.label}
                  size="small"
                  sx={{ 
                    ml: 1,
                    backgroundColor: currentTrustBadge.color,
                    color: 'white',
                    fontSize: '0.7rem'
                  }}
                />
              </Box>
              <StyledSlider
                value={filters.minTrustScore}
                onChange={handleTrustScoreChange}
                min={0}
                max={100}
                marks={trustScoreMarks}
                valueLabelDisplay="auto"
              />
            </FilterSection>

            <FormControlLabel
              control={
                <Switch
                  checked={filters.isKYCRequired}
                  onChange={(e) => onFiltersChange({ ...filters, isKYCRequired: e.target.checked })}
                  size="small"
                />
              }
              label="KYC Verified Only"
            />
          </AccordionDetails>
        </Accordion>

        {/* Availability Filter */}
        <Accordion 
          expanded={expandedSections.includes('availability')}
          onChange={() => handleSectionToggle('availability')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider', mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <Schedule sx={{ mr: 1, fontSize: 20 }} />
              Availability
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <FilterSection>
              <FormControlLabel
                control={
                  <Switch
                    checked={filters.isCurrentlyOpen}
                    onChange={(e) => onFiltersChange({ ...filters, isCurrentlyOpen: e.target.checked })}
                    size="small"
                  />
                }
                label="Currently Open Only"
              />
            </FilterSection>

            <FilterSection>
              <FormControl fullWidth size="small">
                <InputLabel>Sort By</InputLabel>
                <Select
                  value={filters.sortBy}
                  onChange={(e) => onFiltersChange({ ...filters, sortBy: e.target.value as any })}
                  label="Sort By"
                >
                  <MenuItem value="distance">Distance</MenuItem>
                  <MenuItem value="trustScore">Trust Score</MenuItem>
                  <MenuItem value="rating">Rating</MenuItem>
                  <MenuItem value="responseTime">Response Time</MenuItem>
                </Select>
              </FormControl>
            </FilterSection>
          </AccordionDetails>
        </Accordion>

        {/* Language Filter */}
        <Accordion 
          expanded={expandedSections.includes('language')}
          onChange={() => handleSectionToggle('language')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider', mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <Language sx={{ mr: 1, fontSize: 20 }} />
              Languages ({filters.languages.length})
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
              {languageOptions.map((language) => (
                <LanguageChip
                  key={language}
                  label={language}
                  variant={filters.languages.includes(language) ? 'filled' : 'outlined'}
                  color={filters.languages.includes(language) ? 'primary' : 'default'}
                  clickable
                  size="small"
                  onClick={() => handleLanguageChange(language, !filters.languages.includes(language))}
                />
              ))}
            </Box>
          </AccordionDetails>
        </Accordion>

        {/* Payment Methods Filter */}
        <Accordion 
          expanded={expandedSections.includes('payment')}
          onChange={() => handleSectionToggle('payment')}
          elevation={0}
          sx={{ border: '1px solid', borderColor: 'divider' }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center' }}>
              <Payment sx={{ mr: 1, fontSize: 20 }} />
              Payment Methods ({filters.paymentMethods.length})
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <FormGroup>
              {paymentMethodOptions.map((method) => (
                <FormControlLabel
                  key={method.value}
                  control={
                    <Checkbox
                      checked={filters.paymentMethods.includes(method.value)}
                      onChange={(e) => handlePaymentMethodChange(method.value, e.target.checked)}
                      size="small"
                    />
                  }
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      <span style={{ marginRight: 8 }}>{method.icon}</span>
                      {method.label}
                    </Box>
                  }
                />
              ))}
            </FormGroup>
          </AccordionDetails>
        </Accordion>

        <Divider sx={{ my: 2 }} />

        <Button
          variant="outlined"
          fullWidth
          startIcon={<Refresh />}
          onClick={handleResetFilters}
          size="small"
        >
          Reset All Filters
        </Button>
      </CardContent>
    </StyledCard>
  );
};

export default FilterPanel;