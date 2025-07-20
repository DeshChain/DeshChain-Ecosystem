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

import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  IconButton,
  Tooltip,
  useTheme,
  keyframes
} from '@mui/material';
import {
  Celebration as CelebrationIcon,
  Star as StarIcon,
  Close as CloseIcon,
  LocalOffer as OfferIcon,
  Schedule as ScheduleIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { FestivalInfo } from '../../types';
import { useFestivalTheme } from '../../hooks/useFestivalTheme';

interface FestivalBannerProps {
  festival: FestivalInfo;
  showBonusInfo?: boolean;
  showDaysRemaining?: boolean;
  onDismiss?: () => void;
  variant?: 'full' | 'compact' | 'minimal';
  showGreeting?: boolean;
}

const sparkleAnimation = keyframes`
  0% { transform: scale(0) rotate(0deg); opacity: 0; }
  50% { transform: scale(1) rotate(180deg); opacity: 1; }
  100% { transform: scale(0) rotate(360deg); opacity: 0; }
`;

const pulseAnimation = keyframes`
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
`;

export const FestivalBanner: React.FC<FestivalBannerProps> = ({
  festival,
  showBonusInfo = true,
  showDaysRemaining = true,
  onDismiss,
  variant = 'full',
  showGreeting = true
}) => {
  const theme = useTheme();
  const { getFestivalColors, getFestivalIcon } = useFestivalTheme();

  const festivalColors = getFestivalColors(festival.festivalId);
  const FestivalIcon = getFestivalIcon(festival.festivalId);

  const getGradientBackground = () => {
    return `linear-gradient(135deg, 
      ${festivalColors.primary}20, 
      ${festivalColors.secondary}20, 
      ${festivalColors.accent}20
    )`;
  };

  const getBorderGradient = () => {
    return `linear-gradient(90deg, 
      ${festivalColors.primary}, 
      ${festivalColors.secondary}, 
      ${festivalColors.accent}
    )`;
  };

  const Sparkles = () => (
    <Box
      sx={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        pointerEvents: 'none',
        overflow: 'hidden'
      }}
    >
      {[...Array(6)].map((_, i) => (
        <Box
          key={i}
          sx={{
            position: 'absolute',
            top: `${20 + i * 15}%`,
            left: `${10 + i * 15}%`,
            animation: `${sparkleAnimation} 2s infinite`,
            animationDelay: `${i * 0.3}s`,
            color: festivalColors.accent
          }}
        >
          <StarIcon sx={{ fontSize: 12 }} />
        </Box>
      ))}
    </Box>
  );

  if (variant === 'minimal') {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1,
            p: 1,
            borderRadius: 1,
            background: getGradientBackground(),
            border: 1,
            borderColor: festivalColors.primary
          }}
        >
          <FestivalIcon sx={{ color: festivalColors.primary }} />
          <Typography variant="caption" fontWeight="medium">
            {festival.name}
          </Typography>
          {showBonusInfo && (
            <Chip
              label={`${(festival.bonusRate * 100).toFixed(0)}% Bonus`}
              size="small"
              sx={{
                bgcolor: festivalColors.primary,
                color: 'white',
                fontWeight: 'bold'
              }}
            />
          )}
        </Box>
      </motion.div>
    );
  }

  if (variant === 'compact') {
    return (
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Card
          sx={{
            background: getGradientBackground(),
            border: 2,
            borderColor: festivalColors.primary,
            position: 'relative',
            overflow: 'hidden'
          }}
        >
          <CardContent sx={{ p: 2 }}>
            <Box display="flex" alignItems="center" justifyContent="space-between">
              <Box display="flex" alignItems="center" gap={1}>
                <FestivalIcon sx={{ color: festivalColors.primary, fontSize: 24 }} />
                <Box>
                  <Typography variant="subtitle1" fontWeight="bold">
                    {festival.name}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {festival.description}
                  </Typography>
                </Box>
              </Box>

              <Box display="flex" alignItems="center" gap={1}>
                {showBonusInfo && (
                  <Chip
                    icon={<OfferIcon />}
                    label={`${(festival.bonusRate * 100).toFixed(0)}% Bonus`}
                    color="primary"
                    sx={{
                      bgcolor: festivalColors.primary,
                      color: 'white',
                      fontWeight: 'bold',
                      animation: `${pulseAnimation} 2s infinite`
                    }}
                  />
                )}

                {onDismiss && (
                  <IconButton size="small" onClick={onDismiss}>
                    <CloseIcon />
                  </IconButton>
                )}
              </Box>
            </Box>
          </CardContent>
        </Card>
      </motion.div>
    );
  }

  // Full variant
  return (
    <motion.div
      initial={{ opacity: 0, y: -30 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, type: 'spring' }}
    >
      <Card
        sx={{
          position: 'relative',
          overflow: 'hidden',
          background: getGradientBackground(),
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: 4,
            background: getBorderGradient()
          }
        }}
      >
        <Sparkles />
        
        <CardContent sx={{ p: 3 }}>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={2}>
              <Box
                sx={{
                  p: 1,
                  borderRadius: '50%',
                  bgcolor: festivalColors.primary,
                  color: 'white',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}
              >
                <FestivalIcon sx={{ fontSize: 32 }} />
              </Box>
              
              <Box>
                <Typography variant="h5" fontWeight="bold" color={festivalColors.primary}>
                  {festival.name}
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  {festival.description}
                </Typography>
                {festival.significance && (
                  <Typography variant="caption" display="block" sx={{ mt: 0.5 }}>
                    {festival.significance}
                  </Typography>
                )}
              </Box>
            </Box>

            {onDismiss && (
              <IconButton onClick={onDismiss}>
                <CloseIcon />
              </IconButton>
            )}
          </Box>

          <Box display="flex" flexWrap="wrap" gap={1} alignItems="center">
            {showBonusInfo && (
              <Chip
                icon={<OfferIcon />}
                label={`${(festival.bonusRate * 100).toFixed(0)}% Transaction Bonus`}
                sx={{
                  bgcolor: festivalColors.primary,
                  color: 'white',
                  fontWeight: 'bold',
                  animation: `${pulseAnimation} 2s infinite`
                }}
              />
            )}

            {showDaysRemaining && festival.daysRemaining !== undefined && (
              <Chip
                icon={<ScheduleIcon />}
                label={
                  festival.daysRemaining === 0 
                    ? 'Active Today!' 
                    : `${festival.daysRemaining} days remaining`
                }
                color={festival.daysRemaining <= 7 ? 'error' : 'default'}
                variant="outlined"
              />
            )}

            <Chip
              label={festival.region.replace('_', ' ').toUpperCase()}
              size="small"
              variant="outlined"
            />

            {festival.culturalTheme && (
              <Chip
                label={festival.culturalTheme.toUpperCase()}
                size="small"
                variant="outlined"
                sx={{ color: festivalColors.accent }}
              />
            )}
          </Box>

          {showGreeting && festival.traditionalGreeting && (
            <Box
              sx={{
                mt: 2,
                p: 2,
                borderRadius: 1,
                bgcolor: 'rgba(255,255,255,0.1)',
                border: 1,
                borderColor: festivalColors.secondary
              }}
            >
              <Typography
                variant="body1"
                textAlign="center"
                sx={{
                  fontStyle: 'italic',
                  color: festivalColors.primary,
                  fontWeight: 'medium'
                }}
              >
                {festival.traditionalGreeting}
              </Typography>
            </Box>
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
};