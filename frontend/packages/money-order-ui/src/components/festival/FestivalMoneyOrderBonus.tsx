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
  LinearProgress,
  Tooltip,
  Avatar,
  Badge,
  useTheme,
  alpha
} from '@mui/material';
import {
  Celebration as CelebrationIcon,
  LocalOffer as OfferIcon,
  TrendingUp as BonusIcon,
  Star as StarIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { Festival } from '../../themes/festivals';
import { useFestivalTheme } from '../../themes/FestivalThemeProvider';
import { useLanguage } from '../../hooks/useLanguage';

interface FestivalMoneyOrderBonusProps {
  amount: number;
  festival?: Festival;
  showDetails?: boolean;
}

export const FestivalMoneyOrderBonus: React.FC<FestivalMoneyOrderBonusProps> = ({
  amount,
  festival: propFestival,
  showDetails = true
}) => {
  const theme = useTheme();
  const { currentFestival } = useFestivalTheme();
  const { formatCurrency, currentLanguage } = useLanguage();

  const festival = propFestival || currentFestival;

  if (!festival || !festival.specialOffers || amount <= 0) {
    return null;
  }

  // Calculate festival bonus
  const feeDiscount = festival.specialOffers.find(offer => offer.type === 'fee_discount');
  const bonusAmount = festival.specialOffers.find(offer => offer.type === 'bonus_amount');
  
  const baseFee = amount * 0.01; // 1% base fee
  const discountedFee = feeDiscount 
    ? baseFee * (1 - (feeDiscount.value as number) / 100)
    : baseFee;
  const savings = baseFee - discountedFee;
  const bonus = bonusAmount ? (bonusAmount.value as number) : 0;
  
  const totalBenefit = savings + bonus;
  const benefitPercentage = (totalBenefit / amount) * 100;

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Card
        className="festival-card festival-sparkle"
        sx={{
          background: `linear-gradient(135deg, ${alpha(festival.theme.primary, 0.1)} 0%, ${alpha(festival.theme.secondary, 0.1)} 100%)`,
          border: `2px solid ${festival.theme.primary}`,
          position: 'relative',
          overflow: 'hidden'
        }}
      >
        <CardContent>
          {/* Header */}
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={1}>
              <Badge
                badgeContent={<StarIcon sx={{ fontSize: 12 }} />}
                color="error"
                overlap="circular"
              >
                <Avatar
                  sx={{
                    bgcolor: festival.theme.primary,
                    width: 48,
                    height: 48
                  }}
                >
                  <CelebrationIcon />
                </Avatar>
              </Badge>
              <Box>
                <Typography variant="h6" className="festival-gradient-text">
                  Festival Bonus Active!
                </Typography>
                <Typography variant="caption" color="textSecondary">
                  {festival.localNames[currentLanguage] || festival.name}
                </Typography>
              </Box>
            </Box>
            
            <Chip
              icon={<BonusIcon />}
              label={`+${benefitPercentage.toFixed(1)}%`}
              color="success"
              size="small"
              className="festival-glow"
            />
          </Box>

          {showDetails && (
            <>
              {/* Savings Breakdown */}
              <Box mb={2}>
                <Box display="flex" justifyContent="space-between" mb={1}>
                  <Typography variant="body2" color="textSecondary">
                    Regular Fee
                  </Typography>
                  <Typography variant="body2" sx={{ textDecoration: 'line-through' }}>
                    {formatCurrency(baseFee)}
                  </Typography>
                </Box>
                
                {feeDiscount && (
                  <Box display="flex" justifyContent="space-between" mb={1}>
                    <Box display="flex" alignItems="center" gap={0.5}>
                      <OfferIcon fontSize="small" color="success" />
                      <Typography variant="body2" color="success.main">
                        Festival Discount ({feeDiscount.value}%)
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="success.main" fontWeight="bold">
                      -{formatCurrency(savings)}
                    </Typography>
                  </Box>
                )}
                
                {bonus > 0 && (
                  <Box display="flex" justifyContent="space-between" mb={1}>
                    <Box display="flex" alignItems="center" gap={0.5}>
                      <BonusIcon fontSize="small" color="primary" />
                      <Typography variant="body2" color="primary.main">
                        Bonus Amount
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="primary.main" fontWeight="bold">
                      +{formatCurrency(bonus)}
                    </Typography>
                  </Box>
                )}
                
                <Divider sx={{ my: 1 }} />
                
                <Box display="flex" justifyContent="space-between">
                  <Typography variant="body1" fontWeight="medium">
                    Your Festival Fee
                  </Typography>
                  <Typography variant="body1" fontWeight="bold" color="success.main">
                    {formatCurrency(discountedFee)}
                  </Typography>
                </Box>
              </Box>

              {/* Savings Progress */}
              <Box>
                <Box display="flex" justifyContent="space-between" mb={0.5}>
                  <Typography variant="caption" color="textSecondary">
                    Total Savings
                  </Typography>
                  <Typography variant="caption" color="success.main" fontWeight="bold">
                    {formatCurrency(totalBenefit)}
                  </Typography>
                </Box>
                <Tooltip title={`You save ${benefitPercentage.toFixed(1)}% on this transaction!`}>
                  <LinearProgress
                    variant="determinate"
                    value={Math.min(benefitPercentage * 10, 100)}
                    sx={{
                      height: 8,
                      borderRadius: 4,
                      backgroundColor: alpha(festival.theme.secondary, 0.2),
                      '& .MuiLinearProgress-bar': {
                        borderRadius: 4,
                        background: `linear-gradient(90deg, ${festival.theme.primary} 0%, ${festival.theme.secondary} 100%)`
                      }
                    }}
                  />
                </Tooltip>
              </Box>
            </>
          )}

          {/* Festival Message */}
          <Box mt={2} p={1} bgcolor={alpha(festival.theme.primary, 0.1)} borderRadius={1}>
            <Typography variant="caption" textAlign="center" display="block">
              {festival.greetings[currentLanguage] || festival.greetings.en}
            </Typography>
          </Box>
        </CardContent>

        {/* Decorative Elements */}
        <Box
          sx={{
            position: 'absolute',
            top: -20,
            right: -20,
            width: 80,
            height: 80,
            borderRadius: '50%',
            background: alpha(festival.theme.accent, 0.2),
            filter: 'blur(40px)'
          }}
        />
        <Box
          sx={{
            position: 'absolute',
            bottom: -30,
            left: -30,
            width: 100,
            height: 100,
            borderRadius: '50%',
            background: alpha(festival.theme.primary, 0.2),
            filter: 'blur(50px)'
          }}
        />
      </Card>
    </motion.div>
  );
};

// Mini version for inline display
export const FestivalBonusChip: React.FC<{ amount: number }> = ({ amount }) => {
  const { currentFestival } = useFestivalTheme();
  
  if (!currentFestival || !currentFestival.specialOffers) {
    return null;
  }

  const feeDiscount = currentFestival.specialOffers.find(offer => offer.type === 'fee_discount');
  
  if (!feeDiscount) {
    return null;
  }

  return (
    <motion.div
      initial={{ scale: 0 }}
      animate={{ scale: 1 }}
      transition={{ type: 'spring', stiffness: 500 }}
    >
      <Chip
        icon={<CelebrationIcon />}
        label={`${feeDiscount.value}% Festival Discount!`}
        color="success"
        size="small"
        className="festival-glow"
        sx={{
          background: currentFestival.theme.backgroundGradient,
          color: '#FFFFFF',
          fontWeight: 'bold'
        }}
      />
    </motion.div>
  );
};