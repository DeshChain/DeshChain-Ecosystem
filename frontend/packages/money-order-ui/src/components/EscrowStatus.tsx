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
  Stepper,
  Step,
  StepLabel,
  StepContent,
  Button,
  Alert,
  AlertTitle,
  Chip,
  LinearProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Tooltip,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
  useTheme,
  alpha
} from '@mui/material';
import {
  Lock as LockIcon,
  CheckCircle as CheckIcon,
  Cancel as CancelIcon,
  Warning as WarningIcon,
  Timer as TimerIcon,
  Chat as ChatIcon,
  Report as ReportIcon,
  Info as InfoIcon,
  Payment as PaymentIcon,
  AccountBalance as BankIcon,
  Security as SecurityIcon,
  Speed as SpeedIcon
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';

import { useLanguage } from '../hooks/useLanguage';
import { CountdownTimer } from './CountdownTimer';

interface EscrowInfo {
  escrowId: string;
  tradeId: string;
  amount: { amount: string; denom: string };
  platformFee: { amount: string; denom: string };
  status: 'ACTIVE' | 'RELEASED' | 'REFUNDED' | 'DISPUTED' | 'EXPIRED';
  createdAt: Date;
  expiresAt: Date;
  buyer: string;
  seller: string;
  currentUser: 'buyer' | 'seller';
  paymentMethod: {
    type: string;
    provider?: string;
  };
  tradeStatus: 'MATCHED' | 'PAYMENT_PENDING' | 'PAYMENT_CONFIRMED' | 'COMPLETED';
}

interface EscrowStatusProps {
  escrow: EscrowInfo;
  onConfirmPayment?: () => void;
  onDispute?: (reason: string) => void;
  onCancel?: () => void;
}

export const EscrowStatus: React.FC<EscrowStatusProps> = ({
  escrow,
  onConfirmPayment,
  onDispute,
  onCancel
}) => {
  const theme = useTheme();
  const { formatCurrency, t } = useLanguage();
  const [showDisputeDialog, setShowDisputeDialog] = useState(false);
  const [disputeReason, setDisputeReason] = useState('');
  const [showInfoDialog, setShowInfoDialog] = useState(false);
  
  const steps = [
    {
      label: 'Escrow Created',
      icon: <LockIcon />,
      description: 'Funds securely locked in escrow',
      completed: true
    },
    {
      label: 'Trade Matched',
      icon: <CheckIcon />,
      description: 'Buyer and seller matched',
      completed: true
    },
    {
      label: 'Payment Sent',
      icon: <PaymentIcon />,
      description: escrow.currentUser === 'buyer' ? 'Send payment via selected method' : 'Waiting for buyer payment',
      completed: escrow.tradeStatus !== 'MATCHED'
    },
    {
      label: 'Payment Confirmed',
      icon: <BankIcon />,
      description: escrow.currentUser === 'seller' ? 'Confirm payment received' : 'Waiting for seller confirmation',
      completed: escrow.tradeStatus === 'PAYMENT_CONFIRMED' || escrow.tradeStatus === 'COMPLETED'
    },
    {
      label: 'Trade Completed',
      icon: <CheckIcon />,
      description: 'Funds released from escrow',
      completed: escrow.tradeStatus === 'COMPLETED'
    }
  ];
  
  const getCurrentStep = () => {
    switch (escrow.tradeStatus) {
      case 'MATCHED': return 2;
      case 'PAYMENT_PENDING': return 3;
      case 'PAYMENT_CONFIRMED': return 4;
      case 'COMPLETED': return 5;
      default: return 2;
    }
  };
  
  const handleDispute = () => {
    if (disputeReason.trim() && onDispute) {
      onDispute(disputeReason);
      setShowDisputeDialog(false);
    }
  };
  
  const getStatusColor = () => {
    switch (escrow.status) {
      case 'ACTIVE': return 'primary';
      case 'RELEASED': return 'success';
      case 'REFUNDED': return 'info';
      case 'DISPUTED': return 'error';
      case 'EXPIRED': return 'default';
      default: return 'default';
    }
  };
  
  const getTimeRemaining = () => {
    const now = new Date();
    const remaining = escrow.expiresAt.getTime() - now.getTime();
    const hours = Math.floor(remaining / (1000 * 60 * 60));
    const minutes = Math.floor((remaining % (1000 * 60 * 60)) / (1000 * 60));
    return { hours, minutes, isExpired: remaining <= 0 };
  };
  
  const timeInfo = getTimeRemaining();
  
  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Box>
            <Typography variant="h5" gutterBottom>
              Escrow Status
            </Typography>
            <Box display="flex" gap={1} alignItems="center">
              <Chip
                label={escrow.status}
                color={getStatusColor()}
                size="small"
              />
              <Typography variant="body2" color="text.secondary">
                ID: {escrow.escrowId}
              </Typography>
            </Box>
          </Box>
          
          <IconButton onClick={() => setShowInfoDialog(true)}>
            <InfoIcon />
          </IconButton>
        </Box>
        
        {/* Timer */}
        {escrow.status === 'ACTIVE' && (
          <Alert 
            severity={timeInfo.hours < 1 ? 'warning' : 'info'}
            icon={<TimerIcon />}
            sx={{ mb: 3 }}
          >
            <AlertTitle>Time Remaining</AlertTitle>
            {timeInfo.isExpired ? (
              'Escrow has expired - funds will be refunded automatically'
            ) : (
              <CountdownTimer 
                expiresAt={escrow.expiresAt} 
                onExpire={() => {}}
              />
            )}
          </Alert>
        )}
        
        {/* Amount */}
        <Box sx={{ 
          p: 2, 
          mb: 3, 
          bgcolor: alpha(theme.palette.primary.main, 0.05),
          borderRadius: 2
        }}>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            Escrow Amount
          </Typography>
          <Typography variant="h4" color="primary">
            {formatCurrency(parseFloat(escrow.amount.amount) / 1000000)}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Platform fee: {formatCurrency(parseFloat(escrow.platformFee.amount) / 1000000)}
          </Typography>
        </Box>
        
        {/* Progress Steps */}
        <Stepper activeStep={getCurrentStep() - 1} orientation="vertical">
          {steps.map((step, index) => (
            <Step key={index} completed={step.completed}>
              <StepLabel
                StepIconComponent={() => (
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    transition={{ delay: index * 0.1 }}
                  >
                    <Box
                      sx={{
                        width: 40,
                        height: 40,
                        borderRadius: '50%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        bgcolor: step.completed ? 'success.main' : 'grey.300',
                        color: 'white'
                      }}
                    >
                      {step.icon}
                    </Box>
                  </motion.div>
                )}
              >
                {step.label}
              </StepLabel>
              <StepContent>
                <Typography variant="body2" color="text.secondary">
                  {step.description}
                </Typography>
                
                {/* Action buttons for current step */}
                {index === getCurrentStep() - 1 && escrow.status === 'ACTIVE' && (
                  <Box mt={2}>
                    {escrow.currentUser === 'buyer' && escrow.tradeStatus === 'MATCHED' && (
                      <Alert severity="info" sx={{ mb: 2 }}>
                        Send payment to seller using {escrow.paymentMethod.provider || escrow.paymentMethod.type}
                      </Alert>
                    )}
                    
                    {escrow.currentUser === 'seller' && escrow.tradeStatus === 'PAYMENT_PENDING' && (
                      <Button
                        variant="contained"
                        color="success"
                        onClick={onConfirmPayment}
                        startIcon={<CheckIcon />}
                      >
                        Confirm Payment Received
                      </Button>
                    )}
                  </Box>
                )}
              </StepContent>
            </Step>
          ))}
        </Stepper>
        
        {/* Actions */}
        <Box mt={3} display="flex" gap={2} justifyContent="flex-end">
          {escrow.status === 'ACTIVE' && escrow.tradeStatus !== 'COMPLETED' && (
            <>
              <Button
                variant="outlined"
                color="error"
                startIcon={<ReportIcon />}
                onClick={() => setShowDisputeDialog(true)}
              >
                Raise Dispute
              </Button>
              
              {escrow.tradeStatus === 'MATCHED' && (
                <Button
                  variant="outlined"
                  onClick={onCancel}
                  startIcon={<CancelIcon />}
                >
                  Cancel Trade
                </Button>
              )}
            </>
          )}
        </Box>
        
        {/* Security Info */}
        <Box mt={3} p={2} bgcolor={alpha(theme.palette.info.main, 0.05)} borderRadius={1}>
          <Box display="flex" alignItems="center" gap={1} mb={1}>
            <SecurityIcon color="info" />
            <Typography variant="subtitle2">
              Your funds are secure
            </Typography>
          </Box>
          <Typography variant="body2" color="text.secondary">
            All funds are held in a decentralized escrow smart contract until both parties confirm the transaction.
            In case of disputes, our resolution team will review and decide fairly.
          </Typography>
        </Box>
      </CardContent>
      
      {/* Dispute Dialog */}
      <Dialog
        open={showDisputeDialog}
        onClose={() => setShowDisputeDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Raise a Dispute</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ mb: 2 }}>
            Please try to resolve issues with your trading partner first. 
            Disputes should only be raised for serious issues.
          </Alert>
          
          <TextField
            fullWidth
            multiline
            rows={4}
            label="Reason for dispute"
            value={disputeReason}
            onChange={(e) => setDisputeReason(e.target.value)}
            helperText="Please provide detailed information about the issue"
          />
          
          <Typography variant="caption" color="text.secondary" sx={{ mt: 2, display: 'block' }}>
            Common reasons: Payment not received, incorrect amount, fraudulent behavior
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowDisputeDialog(false)}>
            Cancel
          </Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleDispute}
            disabled={!disputeReason.trim()}
          >
            Submit Dispute
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Info Dialog */}
      <Dialog
        open={showInfoDialog}
        onClose={() => setShowInfoDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Escrow Information</DialogTitle>
        <DialogContent>
          <List>
            <ListItem>
              <ListItemIcon>
                <LockIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary="Secure Escrow"
                secondary="Your funds are locked in a smart contract until trade completion"
              />
            </ListItem>
            
            <Divider />
            
            <ListItem>
              <ListItemIcon>
                <TimerIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary="24 Hour Protection"
                secondary="If no match found, full refund including fees"
              />
            </ListItem>
            
            <Divider />
            
            <ListItem>
              <ListItemIcon>
                <SpeedIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary="Fast Resolution"
                secondary="Most trades complete within 30 minutes"
              />
            </ListItem>
            
            <Divider />
            
            <ListItem>
              <ListItemIcon>
                <WarningIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary="Dispute Protection"
                secondary="Fair resolution process with evidence review"
              />
            </ListItem>
          </List>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowInfoDialog(false)}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
};