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
  Typography,
  Paper,
  Chip,
  Tooltip,
  LinearProgress
} from '@mui/material';
import {
  Info as InfoIcon,
  LocalOffer as OfferIcon,
  Speed as SpeedIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

interface FeeDisplayProps {
  amount: string;
  fee: string;
  festivalBonus?: number;
  priority?: 'standard' | 'fast' | 'instant';
  showDetails?: boolean;
}

export const FeeDisplay: React.FC<FeeDisplayProps> = ({
  amount,
  fee,
  festivalBonus,
  priority = 'standard',
  showDetails = true
}) => {
  const baseFee = parseFloat(fee) / 1000000; // Convert from micro units
  const priorityMultiplier = priority === 'fast' ? 1.5 : priority === 'instant' ? 2 : 1;
  const totalFee = baseFee * priorityMultiplier;
  const festivalDiscount = festivalBonus ? totalFee * festivalBonus : 0;
  const finalFee = totalFee - festivalDiscount;

  const getPriorityInfo = () => {
    switch (priority) {
      case 'fast':
        return { label: 'Fast', time: '2-5 min', color: 'warning' as const };
      case 'instant':
        return { label: 'Instant', time: '< 1 min', color: 'error' as const };
      default:
        return { label: 'Standard', time: '5-10 min', color: 'info' as const };
    }
  };

  const priorityInfo = getPriorityInfo();

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
    >
      <Paper 
        elevation={1} 
        sx={{ 
          p: 2, 
          borderRadius: 2,
          background: 'linear-gradient(135deg, rgba(255,107,53,0.05), rgba(19,136,8,0.05))'
        }}
      >
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
          <Typography variant="subtitle2" fontWeight="medium">
            Transaction Fee
          </Typography>
          <Chip
            icon={<SpeedIcon />}
            label={priorityInfo.label}
            size="small"
            color={priorityInfo.color}
          />
        </Box>

        {showDetails && (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center">
              <Typography variant="body2" color="text.secondary">
                Base Fee:
              </Typography>
              <Typography variant="body2">
                {baseFee.toFixed(2)} NAMO
              </Typography>
            </Box>

            {priority !== 'standard' && (
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Typography variant="body2" color="text.secondary">
                  Priority Multiplier:
                </Typography>
                <Typography variant="body2">
                  {priorityMultiplier}x
                </Typography>
              </Box>
            )}

            {festivalBonus && festivalBonus > 0 && (
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box display="flex" alignItems="center" gap={0.5}>
                  <OfferIcon sx={{ fontSize: 16, color: 'success.main' }} />
                  <Typography variant="body2" color="success.main">
                    Festival Discount:
                  </Typography>
                </Box>
                <Typography variant="body2" color="success.main" fontWeight="medium">
                  -{festivalDiscount.toFixed(2)} NAMO
                </Typography>
              </Box>
            )}

            <Box
              sx={{
                pt: 1,
                mt: 1,
                borderTop: 1,
                borderColor: 'divider'
              }}
              display="flex"
              justifyContent="space-between"
              alignItems="center"
            >
              <Typography variant="subtitle2" fontWeight="bold">
                Total Fee:
              </Typography>
              <Box display="flex" alignItems="baseline" gap={0.5}>
                <Typography variant="subtitle1" fontWeight="bold" color="primary">
                  {finalFee.toFixed(2)} NAMO
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  (≈ ₹{(finalFee * 0.075).toFixed(2)})
                </Typography>
              </Box>
            </Box>

            <Box display="flex" alignItems="center" gap={1} mt={1}>
              <InfoIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
              <Typography variant="caption" color="text.secondary">
                Estimated processing time: {priorityInfo.time}
              </Typography>
            </Box>
          </Box>
        )}

        {!showDetails && (
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="body2" color="text.secondary">
              Total:
            </Typography>
            <Typography variant="body2" fontWeight="medium">
              {finalFee.toFixed(2)} NAMO
            </Typography>
          </Box>
        )}

        {/* Visual Progress Indicator */}
        <Box mt={2}>
          <LinearProgress
            variant="determinate"
            value={priority === 'instant' ? 100 : priority === 'fast' ? 66 : 33}
            sx={{
              height: 4,
              borderRadius: 2,
              bgcolor: 'action.hover',
              '& .MuiLinearProgress-bar': {
                borderRadius: 2,
                bgcolor: priorityInfo.color === 'error' ? 'error.main' :
                        priorityInfo.color === 'warning' ? 'warning.main' : 'info.main'
              }
            }}
          />
        </Box>
      </Paper>
    </motion.div>
  );
};