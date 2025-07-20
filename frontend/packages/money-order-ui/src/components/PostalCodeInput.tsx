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

import React, { useState, useEffect, useCallback } from 'react';
import {
  TextField,
  InputAdornment,
  CircularProgress,
  Chip,
  Box,
  Typography,
  Collapse,
  Card,
  CardContent,
  IconButton,
  Tooltip,
  Alert
} from '@mui/material';
import {
  LocationOn as LocationIcon,
  Check as CheckIcon,
  Error as ErrorIcon,
  Info as InfoIcon,
  MyLocation as MyLocationIcon
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';

import { PostalCodeService, PostalCodeInfo } from '../services/postalCodeService';
import { useLanguage } from '../hooks/useLanguage';

interface PostalCodeInputProps {
  value: string;
  onChange: (value: string, info?: PostalCodeInfo) => void;
  label?: string;
  error?: boolean;
  helperText?: string;
  disabled?: boolean;
  required?: boolean;
  fullWidth?: boolean;
  showDetails?: boolean;
  autoDetect?: boolean;
  onPostalInfoLoad?: (info: PostalCodeInfo) => void;
}

export const PostalCodeInput: React.FC<PostalCodeInputProps> = ({
  value,
  onChange,
  label = 'Postal Code',
  error,
  helperText,
  disabled = false,
  required = false,
  fullWidth = true,
  showDetails = true,
  autoDetect = false,
  onPostalInfoLoad
}) => {
  const { currentLanguage } = useLanguage();
  const [loading, setLoading] = useState(false);
  const [postalInfo, setPostalInfo] = useState<PostalCodeInfo | null>(null);
  const [validationError, setValidationError] = useState<string>('');
  const [showInfo, setShowInfo] = useState(false);

  // Validate and fetch postal info
  const fetchPostalInfo = useCallback(async (pincode: string) => {
    if (!PostalCodeService.isValidPincode(pincode)) {
      setPostalInfo(null);
      setValidationError('');
      return;
    }

    setLoading(true);
    setValidationError('');

    try {
      const info = await PostalCodeService.getPostalCodeInfo(pincode);
      
      if (info) {
        setPostalInfo(info);
        setShowInfo(true);
        onPostalInfoLoad?.(info);
      } else {
        setValidationError('Invalid postal code');
        setPostalInfo(null);
      }
    } catch (error) {
      setValidationError('Failed to verify postal code');
      setPostalInfo(null);
    } finally {
      setLoading(false);
    }
  }, [onPostalInfoLoad]);

  // Debounced fetch
  useEffect(() => {
    if (value.length === 6) {
      const timer = setTimeout(() => {
        fetchPostalInfo(value);
      }, 500);
      
      return () => clearTimeout(timer);
    } else {
      setPostalInfo(null);
      setValidationError('');
    }
  }, [value, fetchPostalInfo]);

  // Handle input change
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = event.target.value.replace(/\D/g, '').slice(0, 6);
    onChange(newValue, postalInfo || undefined);
  };

  // Auto-detect location (mock)
  const handleAutoDetect = async () => {
    setLoading(true);
    
    // In real implementation, would use geolocation API
    setTimeout(() => {
      const mockPincode = '110001'; // Delhi example
      onChange(mockPincode);
      setLoading(false);
    }, 1000);
  };

  // Get validation state
  const getValidationState = () => {
    if (loading) return 'loading';
    if (validationError) return 'error';
    if (postalInfo) return 'success';
    return 'idle';
  };

  const validationState = getValidationState();

  return (
    <Box>
      <TextField
        value={value}
        onChange={handleChange}
        label={label}
        error={error || validationState === 'error'}
        helperText={helperText || validationError}
        disabled={disabled}
        required={required}
        fullWidth={fullWidth}
        placeholder="000000"
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <LocationIcon color={validationState === 'success' ? 'success' : 'action'} />
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position="end">
              <AnimatePresence mode="wait">
                {loading && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0 }}
                    animate={{ opacity: 1, scale: 1 }}
                    exit={{ opacity: 0, scale: 0 }}
                  >
                    <CircularProgress size={20} />
                  </motion.div>
                )}
                
                {!loading && validationState === 'success' && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0 }}
                    animate={{ opacity: 1, scale: 1 }}
                    exit={{ opacity: 0, scale: 0 }}
                  >
                    <CheckIcon color="success" />
                  </motion.div>
                )}
                
                {!loading && validationState === 'error' && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0 }}
                    animate={{ opacity: 1, scale: 1 }}
                    exit={{ opacity: 0, scale: 0 }}
                  >
                    <ErrorIcon color="error" />
                  </motion.div>
                )}
                
                {autoDetect && !loading && !value && (
                  <Tooltip title="Detect my location">
                    <IconButton size="small" onClick={handleAutoDetect}>
                      <MyLocationIcon />
                    </IconButton>
                  </Tooltip>
                )}
              </AnimatePresence>
            </InputAdornment>
          )
        }}
      />

      {/* Postal Info Display */}
      {showDetails && postalInfo && (
        <Collapse in={showInfo}>
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
          >
            <Card
              variant="outlined"
              sx={{
                mt: 1,
                borderColor: 'success.main',
                bgcolor: 'success.lighter'
              }}
            >
              <CardContent sx={{ py: 1.5 }}>
                <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
                  <Typography variant="subtitle2" color="success.dark">
                    Location Verified ✓
                  </Typography>
                  <IconButton
                    size="small"
                    onClick={() => setShowInfo(!showInfo)}
                  >
                    <InfoIcon fontSize="small" />
                  </IconButton>
                </Box>

                <Box display="flex" flexDirection="column" gap={0.5}>
                  <Typography variant="body2">
                    <strong>{postalInfo.officeName}</strong>
                  </Typography>
                  
                  <Box display="flex" gap={0.5} flexWrap="wrap">
                    <Chip
                      label={postalInfo.districtName}
                      size="small"
                      variant="outlined"
                    />
                    <Chip
                      label={postalInfo.stateName}
                      size="small"
                      variant="outlined"
                    />
                    <Chip
                      label={postalInfo.officeType}
                      size="small"
                      color={postalInfo.officeType === 'HO' ? 'primary' : 'default'}
                    />
                  </Box>

                  <Typography variant="caption" color="text.secondary">
                    {postalInfo.taluk} • {postalInfo.divisionName} • {postalInfo.circleName}
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </motion.div>
        </Collapse>
      )}

      {/* Auto-detect hint */}
      {autoDetect && !value && !loading && (
        <Alert
          severity="info"
          sx={{ mt: 1 }}
          action={
            <Button size="small" onClick={handleAutoDetect}>
              Detect
            </Button>
          }
        >
          Click detect to automatically find your postal code
        </Alert>
      )}
    </Box>
  );
};

// Postal code display component
export const PostalCodeDisplay: React.FC<{
  pincode: string;
  size?: 'small' | 'medium' | 'large';
  showIcon?: boolean;
}> = ({ pincode, size = 'medium', showIcon = true }) => {
  const [postalInfo, setPostalInfo] = useState<PostalCodeInfo | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (PostalCodeService.isValidPincode(pincode)) {
      setLoading(true);
      PostalCodeService.getPostalCodeInfo(pincode)
        .then(info => {
          setPostalInfo(info);
          setLoading(false);
        })
        .catch(() => setLoading(false));
    }
  }, [pincode]);

  const fontSize = size === 'small' ? '0.875rem' : size === 'large' ? '1.25rem' : '1rem';

  return (
    <Box display="inline-flex" alignItems="center" gap={0.5}>
      {showIcon && <LocationIcon sx={{ fontSize }} />}
      
      <Typography variant="body1" sx={{ fontSize, fontFamily: 'monospace' }}>
        {pincode}
      </Typography>
      
      {loading && <CircularProgress size={fontSize} />}
      
      {!loading && postalInfo && (
        <Tooltip title={`${postalInfo.officeName}, ${postalInfo.districtName}, ${postalInfo.stateName}`}>
          <CheckIcon color="success" sx={{ fontSize }} />
        </Tooltip>
      )}
    </Box>
  );
};