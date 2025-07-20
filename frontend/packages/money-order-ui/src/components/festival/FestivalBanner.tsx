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
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  IconButton,
  Collapse,
  Button,
  Alert,
  useTheme,
  alpha
} from '@mui/material';
import {
  Close as CloseIcon,
  Celebration as CelebrationIcon,
  CardGiftcard as GiftIcon,
  LocalOffer as OfferIcon,
  Timer as TimerIcon
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';

import { Festival, getFestivalGreeting, getDaysUntilNextFestival } from '../../themes/festivals';
import { useFestivalTheme } from '../../themes/FestivalThemeProvider';
import { useLanguage } from '../../hooks/useLanguage';

interface FestivalBannerProps {
  festival?: Festival;
  variant?: 'full' | 'compact' | 'minimal';
  showOffers?: boolean;
  dismissible?: boolean;
  onDismiss?: () => void;
}

export const FestivalBanner: React.FC<FestivalBannerProps> = ({
  festival: propFestival,
  variant = 'full',
  showOffers = true,
  dismissible = true,
  onDismiss
}) => {
  const theme = useTheme();
  const { currentFestival } = useFestivalTheme();
  const { currentLanguage, formatCurrency } = useLanguage();
  const [isVisible, setIsVisible] = useState(true);
  const [showFullBanner, setShowFullBanner] = useState(variant === 'full');

  const festival = propFestival || currentFestival;

  // Auto-hide banner after dismissal
  useEffect(() => {
    if (festival && dismissible) {
      const dismissedKey = `festival-banner-dismissed-${festival.id}`;
      const dismissed = localStorage.getItem(dismissedKey);
      if (dismissed === 'true') {
        setIsVisible(false);
      }
    }
  }, [festival, dismissible]);

  const handleDismiss = () => {
    if (festival) {
      const dismissedKey = `festival-banner-dismissed-${festival.id}`;
      localStorage.setItem(dismissedKey, 'true');
    }
    setIsVisible(false);
    onDismiss?.();
  };

  if (!festival || !isVisible) {
    return <NextFestivalReminder />;
  }

  const greeting = getFestivalGreeting(festival.id, currentLanguage);
  const localName = festival.localNames[currentLanguage] || festival.name;

  const renderMinimal = () => (
    <motion.div
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
    >
      <Alert
        severity="info"
        icon={<CelebrationIcon />}
        action={
          dismissible && (
            <IconButton size="small" onClick={handleDismiss}>
              <CloseIcon fontSize="small" />
            </IconButton>
          )
        }
        sx={{
          background: festival.theme.backgroundGradient,
          color: festival.theme.textPrimary,
          '& .MuiAlert-icon': {
            color: festival.theme.primary
          }
        }}
      >
        <Typography variant="body2" fontWeight="medium">
          {greeting}
        </Typography>
      </Alert>
    </motion.div>
  );

  const renderCompact = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.95 }}
    >
      <Card
        className="festival-card festival-glow"
        sx={{
          background: festival.theme.backgroundGradient,
          color: festival.theme.textPrimary,
          position: 'relative',
          overflow: 'hidden'
        }}
      >
        <CardContent sx={{ py: 2 }}>
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Box display="flex" alignItems="center" gap={2}>
              <CelebrationIcon sx={{ fontSize: 40, color: festival.theme.primary }} />
              <Box>
                <Typography variant="h6" className="festival-gradient-text">
                  {localName}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {greeting}
                </Typography>
              </Box>
            </Box>
            
            <Box display="flex" alignItems="center" gap={1}>
              {festival.specialOffers && showOffers && (
                <Chip
                  icon={<OfferIcon />}
                  label={`${festival.specialOffers[0].value}% OFF`}
                  color="error"
                  size="small"
                  className="festival-sparkle"
                />
              )}
              
              <Button
                size="small"
                onClick={() => setShowFullBanner(!showFullBanner)}
              >
                {showFullBanner ? 'Less' : 'More'}
              </Button>
              
              {dismissible && (
                <IconButton size="small" onClick={handleDismiss}>
                  <CloseIcon />
                </IconButton>
              )}
            </Box>
          </Box>
        </CardContent>
      </Card>
    </motion.div>
  );

  const renderFull = () => (
    <Collapse in={showFullBanner}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2 }}
      >
        <Box mt={2}>
          <Card
            className="festival-card"
            sx={{
              background: alpha(theme.palette.background.paper, 0.9),
              backdropFilter: 'blur(10px)'
            }}
          >
            <CardContent>
              <Typography variant="h5" gutterBottom className="festival-gradient-text">
                Celebrating {localName}!
              </Typography>
              
              {/* Festival Information */}
              <Box mb={2}>
                <Typography variant="body1" paragraph>
                  {greeting}
                </Typography>
                
                <Box display="flex" gap={1} flexWrap="wrap">
                  <Chip
                    size="small"
                    label={festival.type}
                    color="primary"
                    variant="outlined"
                  />
                  {festival.regions.map(region => (
                    <Chip
                      key={region}
                      size="small"
                      label={region}
                      variant="outlined"
                    />
                  ))}
                </Box>
              </Box>
              
              {/* Special Offers */}
              {festival.specialOffers && showOffers && (
                <Box mb={2}>
                  <Typography variant="h6" gutterBottom>
                    Festival Offers 🎁
                  </Typography>
                  {festival.specialOffers.map((offer, index) => (
                    <Alert
                      key={index}
                      severity="success"
                      icon={<GiftIcon />}
                      sx={{ mb: 1 }}
                    >
                      <Typography variant="body2">
                        {offer.message[currentLanguage] || offer.message.en}
                      </Typography>
                    </Alert>
                  ))}
                </Box>
              )}
              
              {/* Festival Duration */}
              <Box display="flex" alignItems="center" gap={1}>
                <TimerIcon fontSize="small" />
                <Typography variant="caption" color="textSecondary">
                  Festival period: {new Date(festival.startDate).toLocaleDateString()} - {new Date(festival.endDate).toLocaleDateString()}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Box>
      </motion.div>
    </Collapse>
  );

  return (
    <AnimatePresence>
      {variant === 'minimal' && renderMinimal()}
      {variant === 'compact' && (
        <>
          {renderCompact()}
          {renderFull()}
        </>
      )}
      {variant === 'full' && (
        <>
          {renderCompact()}
          {renderFull()}
        </>
      )}
    </AnimatePresence>
  );
};

// Component to show reminder for next festival
const NextFestivalReminder: React.FC = () => {
  const { formatNumber } = useLanguage();
  const nextFestival = getDaysUntilNextFestival();
  
  if (!nextFestival || nextFestival.days > 7) {
    return null;
  }
  
  return (
    <motion.div
      initial={{ opacity: 0, x: 100 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: 1 }}
    >
      <Alert
        severity="info"
        sx={{
          background: 'linear-gradient(135deg, #E3F2FD 0%, #BBDEFB 100%)',
          '& .MuiAlert-icon': {
            color: '#1976D2'
          }
        }}
      >
        <Typography variant="body2">
          {nextFestival.festival.name} in {formatNumber(nextFestival.days)} days! 🎉
        </Typography>
      </Alert>
    </motion.div>
  );
};