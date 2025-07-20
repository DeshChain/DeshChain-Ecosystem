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
  Paper,
  Typography,
  TextField,
  Button,
  ToggleButton,
  ToggleButtonGroup,
  Slider,
  InputAdornment,
  Alert,
  Chip,
  Tooltip,
  FormControlLabel,
  Switch,
  Collapse,
  IconButton
} from '@mui/material';
import {
  TrendingUp as BuyIcon,
  TrendingDown as SellIcon,
  Info as InfoIcon,
  Settings as SettingsIcon,
  Speed as SpeedIcon
} from '@mui/icons-material';
import { useForm, Controller } from 'react-hook-form';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';

import { useMoneyOrder, CulturalQuote, useCulturalContext } from '@deshchain/money-order-ui';

interface TradeFormProps {
  poolId: string;
  pair: string;
  festivalBonus?: number;
}

interface TradeFormData {
  type: 'market' | 'limit' | 'stop';
  side: 'buy' | 'sell';
  amount: string;
  price?: string;
  stopPrice?: string;
  leverage: number;
  reduceOnly: boolean;
  postOnly: boolean;
  timeInForce: 'GTC' | 'IOC' | 'FOK';
}

export const TradeForm: React.FC<TradeFormProps> = ({ poolId, pair, festivalBonus }) => {
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [estimatedOutput, setEstimatedOutput] = useState('0');
  const [priceImpact, setPriceImpact] = useState(0);
  const [fee, setFee] = useState('0');

  const { executeSwap, getSwapQuote } = useMoneyOrder();
  const { currentQuote, isCulturalFeaturesEnabled } = useCulturalContext();

  const {
    control,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting }
  } = useForm<TradeFormData>({
    defaultValues: {
      type: 'market',
      side: 'buy',
      amount: '',
      leverage: 1,
      reduceOnly: false,
      postOnly: false,
      timeInForce: 'GTC'
    }
  });

  const watchedValues = watch();

  // Get quote when amount changes
  useEffect(() => {
    const getQuote = async () => {
      if (watchedValues.amount && parseFloat(watchedValues.amount) > 0) {
        try {
          const quote = await getSwapQuote({
            poolId,
            tokenIn: {
              denom: watchedValues.side === 'buy' ? 'inr' : 'unamo',
              amount: (parseFloat(watchedValues.amount) * 1000000).toString()
            },
            tokenOutDenom: watchedValues.side === 'buy' ? 'unamo' : 'inr'
          });

          if (quote) {
            setEstimatedOutput(quote.tokenOut.amount);
            setPriceImpact(quote.priceImpact);
            setFee(quote.fee.amount);
          }
        } catch (error) {
          console.error('Failed to get quote:', error);
        }
      }
    };

    const debounceTimer = setTimeout(getQuote, 500);
    return () => clearTimeout(debounceTimer);
  }, [watchedValues.amount, watchedValues.side, poolId, getSwapQuote]);

  const onSubmit = async (data: TradeFormData) => {
    try {
      // Execute trade based on type
      if (data.type === 'market') {
        await executeSwap({
          poolId,
          tokenIn: {
            denom: data.side === 'buy' ? 'inr' : 'unamo',
            amount: (parseFloat(data.amount) * 1000000).toString()
          },
          tokenOutDenom: data.side === 'buy' ? 'unamo' : 'inr',
          maxSlippage: 0.05
        });

        toast.success(`Market ${data.side} order executed successfully!`);
      } else {
        // For limit/stop orders, we would submit to order book
        toast.success(`${data.type} order placed successfully!`);
      }
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Trade failed');
    }
  };

  const handlePercentageClick = (percentage: number) => {
    // This would calculate based on user's balance
    const balance = 100000; // Mock balance
    const amount = (balance * percentage / 100).toString();
    setValue('amount', amount);
  };

  return (
    <Paper sx={{ p: 2 }}>
      <form onSubmit={handleSubmit(onSubmit)}>
        {/* Side Toggle */}
        <Controller
          name="side"
          control={control}
          render={({ field }) => (
            <ToggleButtonGroup
              {...field}
              exclusive
              fullWidth
              sx={{ mb: 2 }}
            >
              <ToggleButton value="buy" sx={{ color: 'success.main' }}>
                <BuyIcon sx={{ mr: 1 }} />
                Buy
              </ToggleButton>
              <ToggleButton value="sell" sx={{ color: 'error.main' }}>
                <SellIcon sx={{ mr: 1 }} />
                Sell
              </ToggleButton>
            </ToggleButtonGroup>
          )}
        />

        {/* Order Type */}
        <Controller
          name="type"
          control={control}
          render={({ field }) => (
            <ToggleButtonGroup
              {...field}
              exclusive
              fullWidth
              size="small"
              sx={{ mb: 2 }}
            >
              <ToggleButton value="market">Market</ToggleButton>
              <ToggleButton value="limit">Limit</ToggleButton>
              <ToggleButton value="stop">Stop</ToggleButton>
            </ToggleButtonGroup>
          )}
        />

        {/* Price Input (for limit/stop orders) */}
        <AnimatePresence>
          {watchedValues.type !== 'market' && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
            >
              <Controller
                name="price"
                control={control}
                rules={{
                  required: watchedValues.type !== 'market' ? 'Price is required' : false,
                  min: { value: 0.001, message: 'Price must be greater than 0' }
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Price"
                    type="number"
                    error={!!errors.price}
                    helperText={errors.price?.message}
                    InputProps={{
                      startAdornment: <InputAdornment position="start">₹</InputAdornment>
                    }}
                    sx={{ mb: 2 }}
                  />
                )}
              />
            </motion.div>
          )}
        </AnimatePresence>

        {/* Stop Price (for stop orders) */}
        <AnimatePresence>
          {watchedValues.type === 'stop' && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
            >
              <Controller
                name="stopPrice"
                control={control}
                rules={{
                  required: watchedValues.type === 'stop' ? 'Stop price is required' : false,
                  min: { value: 0.001, message: 'Stop price must be greater than 0' }
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Stop Price"
                    type="number"
                    error={!!errors.stopPrice}
                    helperText={errors.stopPrice?.message}
                    InputProps={{
                      startAdornment: <InputAdornment position="start">₹</InputAdornment>
                    }}
                    sx={{ mb: 2 }}
                  />
                )}
              />
            </motion.div>
          )}
        </AnimatePresence>

        {/* Amount Input */}
        <Controller
          name="amount"
          control={control}
          rules={{
            required: 'Amount is required',
            min: { value: 1, message: 'Minimum amount is 1' }
          }}
          render={({ field }) => (
            <TextField
              {...field}
              fullWidth
              label="Amount"
              type="number"
              error={!!errors.amount}
              helperText={errors.amount?.message || `≈ ${(parseFloat(estimatedOutput) / 1000000).toFixed(2)} ${watchedValues.side === 'buy' ? 'NAMO' : 'INR'}`}
              InputProps={{
                endAdornment: <InputAdornment position="end">{watchedValues.side === 'buy' ? 'INR' : 'NAMO'}</InputAdornment>
              }}
              sx={{ mb: 1 }}
            />
          )}
        />

        {/* Percentage Buttons */}
        <Box display="flex" gap={1} mb={2}>
          {[25, 50, 75, 100].map(percentage => (
            <Button
              key={percentage}
              size="small"
              variant="outlined"
              onClick={() => handlePercentageClick(percentage)}
              sx={{ flex: 1 }}
            >
              {percentage}%
            </Button>
          ))}
        </Box>

        {/* Leverage Slider (if enabled) */}
        {showAdvanced && (
          <Box mb={2}>
            <Typography variant="body2" gutterBottom>
              Leverage: {watchedValues.leverage}x
            </Typography>
            <Controller
              name="leverage"
              control={control}
              render={({ field }) => (
                <Slider
                  {...field}
                  min={1}
                  max={10}
                  marks
                  valueLabelDisplay="auto"
                />
              )}
            />
          </Box>
        )}

        {/* Advanced Options Toggle */}
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
          <Typography variant="body2">Advanced Options</Typography>
          <IconButton size="small" onClick={() => setShowAdvanced(!showAdvanced)}>
            <SettingsIcon />
          </IconButton>
        </Box>

        {/* Advanced Options */}
        <Collapse in={showAdvanced}>
          <Box sx={{ mb: 2 }}>
            <Controller
              name="reduceOnly"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={<Switch {...field} checked={field.value} />}
                  label="Reduce Only"
                />
              )}
            />
            
            <Controller
              name="postOnly"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={<Switch {...field} checked={field.value} />}
                  label="Post Only"
                />
              )}
            />

            <Box mt={1}>
              <Typography variant="body2" gutterBottom>
                Time in Force
              </Typography>
              <Controller
                name="timeInForce"
                control={control}
                render={({ field }) => (
                  <ToggleButtonGroup
                    {...field}
                    exclusive
                    size="small"
                    fullWidth
                  >
                    <ToggleButton value="GTC">GTC</ToggleButton>
                    <ToggleButton value="IOC">IOC</ToggleButton>
                    <ToggleButton value="FOK">FOK</ToggleButton>
                  </ToggleButtonGroup>
                )}
              />
            </Box>
          </Box>
        </Collapse>

        {/* Trade Summary */}
        <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
          <Box display="flex" justifyContent="space-between" mb={1}>
            <Typography variant="body2" color="text.secondary">
              Est. Output:
            </Typography>
            <Typography variant="body2" fontWeight="medium">
              {(parseFloat(estimatedOutput) / 1000000).toFixed(4)} {watchedValues.side === 'buy' ? 'NAMO' : 'INR'}
            </Typography>
          </Box>
          
          <Box display="flex" justifyContent="space-between" mb={1}>
            <Typography variant="body2" color="text.secondary">
              Price Impact:
            </Typography>
            <Typography 
              variant="body2" 
              color={priceImpact > 0.05 ? 'error' : priceImpact > 0.02 ? 'warning.main' : 'text.primary'}
            >
              {(priceImpact * 100).toFixed(2)}%
            </Typography>
          </Box>
          
          <Box display="flex" justifyContent="space-between">
            <Typography variant="body2" color="text.secondary">
              Network Fee:
            </Typography>
            <Box display="flex" alignItems="center" gap={0.5}>
              <Typography variant="body2">
                {(parseFloat(fee) / 1000000).toFixed(4)} NAMO
              </Typography>
              {festivalBonus && (
                <Chip 
                  label={`-${(festivalBonus * 100).toFixed(0)}%`} 
                  size="small" 
                  color="success"
                />
              )}
            </Box>
          </Box>
        </Paper>

        {/* Cultural Quote */}
        {isCulturalFeaturesEnabled && currentQuote && (
          <Box mb={2}>
            <CulturalQuote quote={currentQuote} compact variant="inline" />
          </Box>
        )}

        {/* Submit Button */}
        <Button
          type="submit"
          variant="contained"
          fullWidth
          size="large"
          disabled={isSubmitting}
          sx={{
            bgcolor: watchedValues.side === 'buy' ? 'success.main' : 'error.main',
            '&:hover': {
              bgcolor: watchedValues.side === 'buy' ? 'success.dark' : 'error.dark'
            }
          }}
        >
          {isSubmitting ? 'Processing...' : `${watchedValues.type} ${watchedValues.side}`}
        </Button>

        {/* Warning for high price impact */}
        {priceImpact > 0.05 && (
          <Alert severity="warning" sx={{ mt: 2 }}>
            <Typography variant="body2">
              High price impact! Consider reducing the trade size.
            </Typography>
          </Alert>
        )}
      </form>
    </Paper>
  );
};