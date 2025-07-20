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
  Paper,
  Typography,
  Grid,
  Divider,
  Avatar,
  Chip
} from '@mui/material';
import {
  ArrowForward as ArrowIcon,
  AccountBalanceWallet as WalletIcon,
  SwapHoriz as SwapIcon,
  Receipt as ReceiptIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

interface TransactionPreviewProps {
  sender: string;
  receiver: string;
  amount: string;
  memo?: string;
  exchangeRate: number;
  fee: string;
}

export const TransactionPreview: React.FC<TransactionPreviewProps> = ({
  sender,
  receiver,
  amount,
  memo,
  exchangeRate,
  fee
}) => {
  const formatAddress = (address: string) => {
    if (!address) return '';
    return `${address.slice(0, 12)}...${address.slice(-8)}`;
  };

  const calculateINR = () => {
    return (parseFloat(amount || '0') * exchangeRate).toFixed(2);
  };

  const finalAmount = () => {
    const amountNum = parseFloat(amount || '0');
    const feeNum = parseFloat(fee) / 1000000;
    return (amountNum - feeNum).toFixed(2);
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Paper
        elevation={1}
        sx={{
          p: 3,
          borderRadius: 2,
          background: 'linear-gradient(135deg, rgba(255,255,255,0.9), rgba(255,255,255,0.95))'
        }}
      >
        <Box display="flex" alignItems="center" gap={1} mb={2}>
          <ReceiptIcon color="primary" />
          <Typography variant="h6" fontWeight="medium">
            Transaction Preview
          </Typography>
        </Box>

        <Grid container spacing={3}>
          {/* From */}
          <Grid item xs={12} sm={5}>
            <Box textAlign="center">
              <Avatar sx={{ bgcolor: 'primary.light', mx: 'auto', mb: 1 }}>
                <WalletIcon />
              </Avatar>
              <Typography variant="body2" color="text.secondary">
                From
              </Typography>
              <Typography variant="body2" fontFamily="monospace" fontWeight="medium">
                {formatAddress(sender)}
              </Typography>
            </Box>
          </Grid>

          {/* Arrow */}
          <Grid item xs={12} sm={2}>
            <Box
              display="flex"
              alignItems="center"
              justifyContent="center"
              height="100%"
            >
              <motion.div
                animate={{ x: [0, 10, 0] }}
                transition={{ duration: 2, repeat: Infinity }}
              >
                <ArrowIcon sx={{ fontSize: 32, color: 'primary.main' }} />
              </motion.div>
            </Box>
          </Grid>

          {/* To */}
          <Grid item xs={12} sm={5}>
            <Box textAlign="center">
              <Avatar sx={{ bgcolor: 'secondary.light', mx: 'auto', mb: 1 }}>
                <WalletIcon />
              </Avatar>
              <Typography variant="body2" color="text.secondary">
                To
              </Typography>
              <Typography variant="body2" fontFamily="monospace" fontWeight="medium">
                {formatAddress(receiver)}
              </Typography>
            </Box>
          </Grid>
        </Grid>

        <Divider sx={{ my: 3 }} />

        {/* Amount Details */}
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="body2" color="text.secondary">
              Amount:
            </Typography>
            <Box display="flex" alignItems="baseline" gap={1}>
              <Typography variant="subtitle1" fontWeight="bold">
                {amount} NAMO
              </Typography>
              <Typography variant="caption" color="text.secondary">
                (≈ ₹{calculateINR()})
              </Typography>
            </Box>
          </Box>

          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="body2" color="text.secondary">
              Network Fee:
            </Typography>
            <Typography variant="body2">
              {(parseFloat(fee) / 1000000).toFixed(2)} NAMO
            </Typography>
          </Box>

          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="body2" color="text.secondary">
              Exchange Rate:
            </Typography>
            <Box display="flex" alignItems="center" gap={0.5}>
              <Typography variant="body2">
                1 NAMO = ₹{exchangeRate}
              </Typography>
              <SwapIcon sx={{ fontSize: 16, color: 'action.active' }} />
            </Box>
          </Box>

          {memo && (
            <Box>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Memo:
              </Typography>
              <Paper variant="outlined" sx={{ p: 1.5, bgcolor: 'background.default' }}>
                <Typography variant="body2">
                  {memo}
                </Typography>
              </Paper>
            </Box>
          )}

          <Divider />

          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="subtitle1" fontWeight="bold">
              Total to Receive:
            </Typography>
            <Box display="flex" alignItems="baseline" gap={1}>
              <Typography variant="h6" fontWeight="bold" color="primary">
                {finalAmount()} NAMO
              </Typography>
              <Typography variant="body2" color="text.secondary">
                (≈ ₹{(parseFloat(finalAmount()) * exchangeRate).toFixed(2)})
              </Typography>
            </Box>
          </Box>

          <Box display="flex" justifyContent="center" gap={1} mt={1}>
            <Chip label="Secure" size="small" color="success" variant="outlined" />
            <Chip label="Instant" size="small" color="info" variant="outlined" />
            <Chip label="Transparent" size="small" color="primary" variant="outlined" />
          </Box>
        </Box>
      </Paper>
    </motion.div>
  );
};