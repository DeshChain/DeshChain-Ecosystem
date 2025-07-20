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
  Card,
  CardContent,
  Typography,
  Box,
  TextField,
  Button,
  InputAdornment,
  Divider,
  Alert,
  Collapse,
  LinearProgress,
  Chip,
  Tooltip,
  IconButton
} from '@mui/material';
import {
  Send as SendIcon,
  AccountBalanceWallet as WalletIcon,
  SwapHoriz as SwapIcon,
  Info as InfoIcon,
  CheckCircle as CheckIcon,
  ContactMail as ContactIcon,
  CurrencyRupee as RupeeIcon
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';
import { useForm, Controller } from 'react-hook-form';
import toast from 'react-hot-toast';

import {
  useMoneyOrder,
  useCulturalContext,
  usePoolData,
  CulturalQuote,
  PatriotismBadge
} from '@deshchain/money-order-ui';
import { AddressBook } from './AddressBook';
import { AmountPresets } from './AmountPresets';
import { FeeDisplay } from './FeeDisplay';
import { TransactionPreview } from './TransactionPreview';

interface SimpleFormData {
  sender: string;
  receiver: string;
  amount: string;
  memo?: string;
}

interface SimpleMoneyOrderFormProps {
  onSuccess: (receiptId: string) => void;
}

export const SimpleMoneyOrderForm: React.FC<SimpleMoneyOrderFormProps> = ({ onSuccess }) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showAddressBook, setShowAddressBook] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [exchangeRate, setExchangeRate] = useState<number>(0.075); // 1 NAMO = 0.075 INR
  const [estimatedFee, setEstimatedFee] = useState<string>('0');

  const { createMoneyOrder, getSwapQuote, isLoading } = useMoneyOrder();
  const { currentQuote, isCulturalFeaturesEnabled, currentFestival } = useCulturalContext();
  const { getRecommendedPool } = usePoolData();

  const {
    control,
    handleSubmit,
    setValue,
    watch,
    formState: { errors, isValid }
  } = useForm<SimpleFormData>({
    mode: 'onChange',
    defaultValues: {
      sender: '',
      receiver: '',
      amount: '',
      memo: ''
    }
  });

  const watchedAmount = watch('amount');
  const watchedReceiver = watch('receiver');

  // Auto-populate sender from wallet if connected
  useEffect(() => {
    const getWalletAddress = async () => {
      try {
        // This would connect to Keplr or other Cosmos wallet
        const address = localStorage.getItem('walletAddress');
        if (address) {
          setValue('sender', address);
        }
      } catch (error) {
        console.error('Failed to get wallet address:', error);
      }
    };
    getWalletAddress();
  }, [setValue]);

  // Get quote when amount changes
  useEffect(() => {
    const getQuote = async () => {
      if (watchedAmount && parseFloat(watchedAmount) > 0) {
        try {
          const pool = await getRecommendedPool();
          if (pool) {
            const quote = await getSwapQuote({
              poolId: pool.poolId,
              tokenIn: {
                denom: 'unamo',
                amount: (parseFloat(watchedAmount) * 1000000).toString()
              },
              tokenOutDenom: 'inr'
            });
            
            if (quote) {
              setExchangeRate(parseFloat(quote.exchangeRate));
              setEstimatedFee(quote.fee.amount);
            }
          }
        } catch (error) {
          console.error('Failed to get quote:', error);
        }
      }
    };

    const debounceTimer = setTimeout(getQuote, 500);
    return () => clearTimeout(debounceTimer);
  }, [watchedAmount, getSwapQuote, getRecommendedPool]);

  const onSubmit = async (data: SimpleFormData) => {
    setIsSubmitting(true);

    try {
      // Get recommended pool
      const pool = await getRecommendedPool();
      if (!pool) {
        throw new Error('No suitable pool found');
      }

      // Create money order
      const receipt = await createMoneyOrder({
        sender: data.sender,
        receiver: data.receiver,
        amount: {
          value: (parseFloat(data.amount) * 1000000).toString(),
          denom: 'unamo'
        },
        poolId: pool.poolId,
        memo: data.memo,
        priority: 'standard',
        culturalPreferences: {
          language: 'en',
          theme: currentFestival?.culturalTheme || 'prosperity',
          includeQuote: isCulturalFeaturesEnabled
        }
      });

      // Show success
      toast.success('Money Order sent successfully!');
      onSuccess(receipt.receiptId);

    } catch (error) {
      console.error('Failed to send money order:', error);
      toast.error(error instanceof Error ? error.message : 'Failed to send money order');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleAddressSelect = (address: string, type: 'sender' | 'receiver') => {
    setValue(type, address);
    setShowAddressBook(false);
  };

  const handleAmountPreset = (amount: string) => {
    setValue('amount', amount);
  };

  const calculateINRAmount = () => {
    if (watchedAmount && parseFloat(watchedAmount) > 0) {
      return (parseFloat(watchedAmount) * exchangeRate).toFixed(2);
    }
    return '0.00';
  };

  return (
    <Card>
      <CardContent sx={{ p: 3 }}>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
          <Typography variant="h5" fontWeight="bold">
            Simple Money Order
          </Typography>
          <PatriotismBadge score={85} size="small" />
        </Box>

        <form onSubmit={handleSubmit(onSubmit)}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            {/* Sender Address */}
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                From (Your Address)
              </Typography>
              <Controller
                name="sender"
                control={control}
                rules={{ 
                  required: 'Sender address is required',
                  pattern: {
                    value: /^desh1[a-z0-9]{38}$/,
                    message: 'Invalid DeshChain address'
                  }
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    placeholder="desh1..."
                    error={!!errors.sender}
                    helperText={errors.sender?.message}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <WalletIcon color="action" />
                        </InputAdornment>
                      ),
                      endAdornment: (
                        <InputAdornment position="end">
                          <Tooltip title="Connect Wallet">
                            <IconButton size="small">
                              <ContactIcon />
                            </IconButton>
                          </Tooltip>
                        </InputAdornment>
                      )
                    }}
                  />
                )}
              />
            </Box>

            {/* Receiver Address */}
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                To (Receiver Address)
              </Typography>
              <Controller
                name="receiver"
                control={control}
                rules={{ 
                  required: 'Receiver address is required',
                  pattern: {
                    value: /^desh1[a-z0-9]{38}$/,
                    message: 'Invalid DeshChain address'
                  }
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    placeholder="desh1..."
                    error={!!errors.receiver}
                    helperText={errors.receiver?.message}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <AccountBalanceWallet color="action" />
                        </InputAdornment>
                      ),
                      endAdornment: (
                        <InputAdornment position="end">
                          <Tooltip title="Address Book">
                            <IconButton 
                              size="small"
                              onClick={() => setShowAddressBook(!showAddressBook)}
                            >
                              <ContactIcon />
                            </IconButton>
                          </Tooltip>
                        </InputAdornment>
                      )
                    }}
                  />
                )}
              />
            </Box>

            {/* Address Book Collapse */}
            <Collapse in={showAddressBook}>
              <AddressBook onSelect={(address) => handleAddressSelect(address, 'receiver')} />
            </Collapse>

            {/* Amount */}
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Amount
              </Typography>
              <Controller
                name="amount"
                control={control}
                rules={{ 
                  required: 'Amount is required',
                  min: { value: 1, message: 'Minimum amount is 1 NAMO' },
                  max: { value: 10000000, message: 'Maximum amount is 10M NAMO' }
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    type="number"
                    placeholder="0"
                    error={!!errors.amount}
                    helperText={errors.amount?.message || `≈ ₹${calculateINRAmount()}`}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <Typography variant="subtitle1" fontWeight="bold">
                            NAMO
                          </Typography>
                        </InputAdornment>
                      ),
                      endAdornment: (
                        <InputAdornment position="end">
                          <SwapIcon color="action" />
                        </InputAdornment>
                      )
                    }}
                  />
                )}
              />
              
              {/* Amount Presets */}
              <Box mt={1}>
                <AmountPresets onSelect={handleAmountPreset} />
              </Box>
            </Box>

            {/* Memo (Optional) */}
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Memo (Optional)
              </Typography>
              <Controller
                name="memo"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    multiline
                    rows={2}
                    placeholder="Add a note..."
                    inputProps={{ maxLength: 200 }}
                  />
                )}
              />
            </Box>

            {/* Fee Display */}
            {watchedAmount && parseFloat(watchedAmount) > 0 && (
              <FeeDisplay
                amount={watchedAmount}
                fee={estimatedFee}
                festivalBonus={currentFestival?.bonusRate}
              />
            )}

            {/* Cultural Quote */}
            {isCulturalFeaturesEnabled && currentQuote && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3 }}
              >
                <CulturalQuote quote={currentQuote} compact variant="inline" />
              </motion.div>
            )}

            <Divider />

            {/* Submit Button */}
            <Box>
              <Button
                type="submit"
                variant="contained"
                size="large"
                fullWidth
                disabled={!isValid || isSubmitting || isLoading}
                startIcon={isSubmitting ? null : <SendIcon />}
              >
                {isSubmitting ? (
                  <>
                    <LinearProgress 
                      sx={{ 
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        right: 0
                      }} 
                    />
                    Processing...
                  </>
                ) : (
                  'Send Money Order'
                )}
              </Button>

              {/* Show Preview Option */}
              <Button
                variant="text"
                size="small"
                fullWidth
                sx={{ mt: 1 }}
                onClick={() => setShowPreview(!showPreview)}
                disabled={!isValid}
              >
                Preview Transaction
              </Button>
            </Box>

            {/* Transaction Preview */}
            <Collapse in={showPreview && isValid}>
              <TransactionPreview
                sender={watch('sender')}
                receiver={watch('receiver')}
                amount={watch('amount')}
                memo={watch('memo')}
                exchangeRate={exchangeRate}
                fee={estimatedFee}
              />
            </Collapse>

            {/* Security Note */}
            <Alert severity="info" icon={<InfoIcon />}>
              <Typography variant="caption">
                Your transaction is secured by blockchain technology and will include a cultural quote on the receipt.
                {currentFestival && ` Enjoy ${(currentFestival.bonusRate * 100).toFixed(0)}% festival bonus!`}
              </Typography>
            </Alert>
          </Box>
        </form>
      </CardContent>
    </Card>
  );
};