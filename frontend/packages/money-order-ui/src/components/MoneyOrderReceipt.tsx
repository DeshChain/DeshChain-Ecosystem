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

import React, { useRef } from 'react';
import {
  Box,
  Paper,
  Typography,
  Divider,
  Grid,
  Chip,
  Button,
  Avatar,
  useTheme
} from '@mui/material';
import {
  CheckCircle as CheckIcon,
  Schedule as TimeIcon,
  Tag as TagIcon,
  AccountBalanceWallet as WalletIcon,
  Print as PrintIcon,
  Share as ShareIcon,
  GetApp as DownloadIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useReactToPrint } from 'react-to-print';

import { ReceiptData } from '../types';
import { QRCodeGenerator } from './utils/QRCodeGenerator';
import { CulturalQuote } from './cultural/CulturalQuote';
import { PatriotismBadge } from './cultural/PatriotismBadge';
import { useCulturalContext } from '../hooks/useCulturalContext';

interface MoneyOrderReceiptProps {
  receipt: ReceiptData;
  showActions?: boolean;
  showQRCode?: boolean;
  showCulturalElements?: boolean;
  variant?: 'full' | 'compact' | 'print';
}

export const MoneyOrderReceipt: React.FC<MoneyOrderReceiptProps> = ({
  receipt,
  showActions = true,
  showQRCode = true,
  showCulturalElements = true,
  variant = 'full'
}) => {
  const theme = useTheme();
  const printRef = useRef<HTMLDivElement>(null);
  const { currentFestival } = useCulturalContext();

  const handlePrint = useReactToPrint({
    content: () => printRef.current,
    documentTitle: `Money Order Receipt - ${receipt.receiptId}`
  });

  const handleDownload = () => {
    // Generate PDF or image download
    const element = document.createElement('a');
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(JSON.stringify(receipt, null, 2)));
    element.setAttribute('download', `receipt-${receipt.receiptId}.json`);
    element.style.display = 'none';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: `Money Order Receipt #${receipt.receiptId}`,
          text: `Money Order of ${receipt.amount.amount} ${receipt.amount.denom} sent successfully`,
          url: `https://deshchain.org/receipt/${receipt.receiptId}`
        });
      } catch (error) {
        console.error('Error sharing:', error);
      }
    }
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 12)}...${address.slice(-10)}`;
  };

  const getStatusColor = () => {
    switch (receipt.status) {
      case 'completed': return 'success';
      case 'processing': return 'warning';
      case 'failed': return 'error';
      default: return 'default';
    }
  };

  const ReceiptHeader = () => (
    <Box textAlign="center" mb={3}>
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ duration: 0.5 }}
      >
        <Avatar
          sx={{
            width: 80,
            height: 80,
            bgcolor: 'success.main',
            mx: 'auto',
            mb: 2
          }}
        >
          <CheckIcon sx={{ fontSize: 48 }} />
        </Avatar>
      </motion.div>
      
      <Typography variant="h4" fontWeight="bold" gutterBottom>
        Money Order Receipt
      </Typography>
      
      <Box display="flex" justifyContent="center" gap={1} alignItems="center">
        <Chip
          label={receipt.status.toUpperCase()}
          color={getStatusColor()}
          size="small"
        />
        <Typography variant="body2" color="text.secondary">
          #{receipt.receiptId}
        </Typography>
      </Box>
    </Box>
  );

  const TransactionDetails = () => (
    <Grid container spacing={3}>
      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            From
          </Typography>
          <Box display="flex" alignItems="center" gap={1} mt={0.5}>
            <WalletIcon color="action" fontSize="small" />
            <Typography variant="body1" fontFamily="monospace">
              {formatAddress(receipt.sender)}
            </Typography>
          </Box>
        </Box>
      </Grid>
      
      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            To
          </Typography>
          <Box display="flex" alignItems="center" gap={1} mt={0.5}>
            <WalletIcon color="action" fontSize="small" />
            <Typography variant="body1" fontFamily="monospace">
              {formatAddress(receipt.receiver)}
            </Typography>
          </Box>
        </Box>
      </Grid>

      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            Amount Sent
          </Typography>
          <Typography variant="h5" fontWeight="bold">
            {parseFloat(receipt.amount.amount) / 1000000} {receipt.amount.denom.toUpperCase()}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            ≈ ₹{(parseFloat(receipt.amount.amount) / 1000000 * parseFloat(receipt.exchangeRate)).toFixed(2)}
          </Typography>
        </Box>
      </Grid>

      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            Network Fee
          </Typography>
          <Typography variant="body1">
            {parseFloat(receipt.fee.amount) / 1000000} {receipt.fee.denom.toUpperCase()}
          </Typography>
          {receipt.patriotismBonus && (
            <Chip
              label={`Saved ${parseFloat(receipt.patriotismBonus.amount) / 1000000} ${receipt.patriotismBonus.denom.toUpperCase()}`}
              color="success"
              size="small"
              sx={{ mt: 0.5 }}
            />
          )}
        </Box>
      </Grid>

      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            Transaction Time
          </Typography>
          <Box display="flex" alignItems="center" gap={1} mt={0.5}>
            <TimeIcon color="action" fontSize="small" />
            <Typography variant="body1">
              {new Date(receipt.timestamp).toLocaleString()}
            </Typography>
          </Box>
        </Box>
      </Grid>

      <Grid item xs={12} sm={6}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            Verification Code
          </Typography>
          <Box display="flex" alignItems="center" gap={1} mt={0.5}>
            <TagIcon color="action" fontSize="small" />
            <Typography variant="body1" fontWeight="medium">
              {receipt.verificationCode}
            </Typography>
          </Box>
        </Box>
      </Grid>
    </Grid>
  );

  const BlockchainConfirmation = () => (
    receipt.blockchainConfirmation && (
      <Box mt={3} p={2} bgcolor="background.default" borderRadius={1}>
        <Typography variant="subtitle2" gutterBottom>
          Blockchain Confirmation
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Typography variant="caption" color="text.secondary">
              Transaction Hash
            </Typography>
            <Typography variant="body2" fontFamily="monospace" sx={{ wordBreak: 'break-all' }}>
              {receipt.blockchainConfirmation.transactionHash}
            </Typography>
          </Grid>
          <Grid item xs={6} md={3}>
            <Typography variant="caption" color="text.secondary">
              Block Height
            </Typography>
            <Typography variant="body2">
              {receipt.blockchainConfirmation.blockHeight.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={6} md={3}>
            <Typography variant="caption" color="text.secondary">
              Confirmations
            </Typography>
            <Typography variant="body2">
              {receipt.blockchainConfirmation.confirmations}
            </Typography>
          </Grid>
        </Grid>
      </Box>
    )
  );

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <Paper
        ref={printRef}
        sx={{
          p: variant === 'compact' ? 3 : 4,
          maxWidth: variant === 'print' ? 'none' : 800,
          mx: 'auto',
          background: currentFestival 
            ? `linear-gradient(135deg, ${theme.palette.background.paper}, rgba(255,215,0,0.05))`
            : theme.palette.background.paper
        }}
      >
        {/* Patriotism Badge */}
        {showCulturalElements && variant !== 'print' && (
          <Box position="absolute" top={16} right={16}>
            <PatriotismBadge score={85} size="small" />
          </Box>
        )}

        {/* Header */}
        <ReceiptHeader />

        {/* Main Content */}
        <Box>
          {/* Transaction Details */}
          <TransactionDetails />

          <Divider sx={{ my: 3 }} />

          {/* QR Code */}
          {showQRCode && (
            <Box mb={3}>
              <Typography variant="subtitle2" gutterBottom textAlign="center">
                Scan to Verify
              </Typography>
              <QRCodeGenerator
                data={{
                  receiptId: receipt.receiptId,
                  orderId: receipt.orderId,
                  amount: receipt.amount.amount,
                  currency: receipt.amount.denom,
                  sender: receipt.sender,
                  receiver: receipt.receiver,
                  timestamp: receipt.timestamp,
                  verificationCode: receipt.verificationCode,
                  culturalQuote: receipt.culturalQuote?.text,
                  festivalBonus: receipt.patriotismBonus?.amount
                }}
                size={200}
                showActions={false}
                culturalDesign={showCulturalElements}
              />
            </Box>
          )}

          {/* Cultural Quote */}
          {showCulturalElements && receipt.culturalQuote && (
            <Box mb={3}>
              <CulturalQuote
                quote={receipt.culturalQuote}
                variant="card"
                showActions={false}
              />
            </Box>
          )}

          {/* Blockchain Confirmation */}
          {variant === 'full' && <BlockchainConfirmation />}

          {/* Action Buttons */}
          {showActions && variant !== 'print' && (
            <Box display="flex" gap={2} justifyContent="center" mt={4}>
              <Button
                variant="outlined"
                startIcon={<PrintIcon />}
                onClick={handlePrint}
              >
                Print
              </Button>
              <Button
                variant="outlined"
                startIcon={<DownloadIcon />}
                onClick={handleDownload}
              >
                Download
              </Button>
              <Button
                variant="outlined"
                startIcon={<ShareIcon />}
                onClick={handleShare}
              >
                Share
              </Button>
            </Box>
          )}
        </Box>

        {/* Footer */}
        <Box mt={4} pt={2} borderTop={1} borderColor="divider" textAlign="center">
          <Typography variant="caption" color="text.secondary">
            This receipt is digitally signed and secured by DeshChain blockchain
          </Typography>
          {variant === 'print' && (
            <Typography variant="caption" display="block" mt={1}>
              Generated on {new Date().toLocaleString()}
            </Typography>
          )}
        </Box>
      </Paper>
    </motion.div>
  );
};