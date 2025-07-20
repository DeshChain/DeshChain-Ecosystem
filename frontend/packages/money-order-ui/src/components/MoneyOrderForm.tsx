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

import React, { useState, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Button,
  Alert,
  Stepper,
  Step,
  StepLabel,
  Collapse,
  Chip
} from '@mui/material';
import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { motion, AnimatePresence } from 'framer-motion';
import { Send as SendIcon, Receipt as ReceiptIcon, Star as StarIcon } from '@mui/icons-material';

import { MoneyOrderFormData, PoolInfo, SwapQuote } from '../types';
import { AmountInput } from './forms/AmountInput';
import { AddressInput } from './forms/AddressInput';
import { PoolSelect } from './forms/PoolSelect';
import { CulturalPreferences } from './forms/CulturalPreferences';
import { CulturalQuote } from './cultural/CulturalQuote';
import { FestivalBanner } from './cultural/FestivalBanner';
import { LoadingSpinner } from './utils/LoadingSpinner';
import { useMoneyOrder } from '../hooks/useMoneyOrder';
import { useCulturalContext } from '../hooks/useCulturalContext';
import { TRANSACTION_PRIORITIES } from '../constants';

const formSchema = z.object({
  sender: z.string().min(1, 'Sender address is required'),
  receiver: z.string().min(1, 'Receiver address is required'),
  amount: z.object({
    value: z.string().min(1, 'Amount is required'),
    denom: z.string().min(1, 'Denomination is required')
  }),
  poolId: z.string().min(1, 'Pool selection is required'),
  memo: z.string().max(200, 'Memo must be less than 200 characters').optional(),
  priority: z.enum(['standard', 'fast', 'instant']),
  culturalPreferences: z.object({
    language: z.string(),
    theme: z.string(),
    includeQuote: z.boolean()
  })
});

interface MoneyOrderFormProps {
  onSubmit: (data: MoneyOrderFormData) => Promise<void>;
  onQuoteUpdate?: (quote: SwapQuote | null) => void;
  initialData?: Partial<MoneyOrderFormData>;
  pools?: PoolInfo[];
  mode?: 'simple' | 'advanced';
  showCulturalFeatures?: boolean;
}

const steps = ['Details', 'Pool & Preferences', 'Review & Submit'];

export const MoneyOrderForm: React.FC<MoneyOrderFormProps> = ({
  onSubmit,
  onQuoteUpdate,
  initialData,
  pools = [],
  mode = 'simple',
  showCulturalFeatures = true
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [quote, setQuote] = useState<SwapQuote | null>(null);
  const [showAdvanced, setShowAdvanced] = useState(mode === 'advanced');

  const { createMoneyOrder, getSwapQuote } = useMoneyOrder();
  const { currentQuote, currentFestival, isPatriotismEnabled } = useCulturalContext();

  const methods = useForm<MoneyOrderFormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      sender: '',
      receiver: '',
      amount: { value: '', denom: 'unamo' },
      poolId: '',
      memo: '',
      priority: 'standard',
      culturalPreferences: {
        language: 'en',
        theme: 'prosperity',
        includeQuote: true
      },
      ...initialData
    }
  });

  const { handleSubmit, watch, formState: { errors, isValid } } = methods;

  const watchedValues = watch();

  // Get quote when relevant fields change
  React.useEffect(() => {
    const getQuote = async () => {
      if (watchedValues.amount.value && watchedValues.poolId) {
        try {
          const newQuote = await getSwapQuote({
            poolId: watchedValues.poolId,
            tokenIn: {
              denom: watchedValues.amount.denom,
              amount: watchedValues.amount.value
            },
            tokenOutDenom: watchedValues.amount.denom === 'unamo' ? 'inr' : 'unamo'
          });
          setQuote(newQuote);
          onQuoteUpdate?.(newQuote);
        } catch (error) {
          console.error('Failed to get quote:', error);
          setQuote(null);
          onQuoteUpdate?.(null);
        }
      }
    };

    const timeoutId = setTimeout(getQuote, 500);
    return () => clearTimeout(timeoutId);
  }, [watchedValues.amount, watchedValues.poolId, getSwapQuote, onQuoteUpdate]);

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const onFormSubmit = useCallback(async (data: MoneyOrderFormData) => {
    setIsSubmitting(true);
    try {
      await onSubmit(data);
    } finally {
      setIsSubmitting(false);
    }
  }, [onSubmit]);

  const getStepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <AddressInput
                name="sender"
                label="Sender Address"
                placeholder="desh1..."
                required
              />
            </Grid>
            <Grid item xs={12}>
              <AddressInput
                name="receiver"
                label="Receiver Address"
                placeholder="desh1..."
                required
              />
            </Grid>
            <Grid item xs={12}>
              <AmountInput
                name="amount"
                label="Amount"
                required
                showQuote={!!quote}
                quote={quote}
              />
            </Grid>
            <Grid item xs={12}>
              <AddressInput
                name="memo"
                label="Memo (Optional)"
                placeholder="Payment description..."
                multiline
                rows={2}
              />
            </Grid>
          </Grid>
        );

      case 1:
        return (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <PoolSelect
                name="poolId"
                label="Select Pool"
                pools={pools}
                required
                showPoolDetails
              />
            </Grid>
            
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Transaction Priority
              </Typography>
              <Grid container spacing={2}>
                {TRANSACTION_PRIORITIES.map((priority) => (
                  <Grid item xs={12} sm={4} key={priority.value}>
                    <Card
                      sx={{
                        cursor: 'pointer',
                        border: watchedValues.priority === priority.value ? 2 : 1,
                        borderColor: watchedValues.priority === priority.value ? 'primary.main' : 'divider',
                        '&:hover': { borderColor: 'primary.main' }
                      }}
                      onClick={() => methods.setValue('priority', priority.value as any)}
                    >
                      <CardContent>
                        <Typography variant="subtitle1" fontWeight="bold">
                          {priority.label}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {priority.description}
                        </Typography>
                        <Typography variant="caption" display="block" mt={1}>
                          Fee: {(priority.feeMultiplier * 100).toFixed(0)}% | Time: {priority.estimatedTime}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            </Grid>

            {showCulturalFeatures && (
              <Grid item xs={12}>
                <CulturalPreferences name="culturalPreferences" />
              </Grid>
            )}
          </Grid>
        );

      case 2:
        return (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Transaction Summary
                  </Typography>
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">From:</Typography>
                      <Typography variant="body1" fontFamily="monospace">
                        {watchedValues.sender.slice(0, 20)}...
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">To:</Typography>
                      <Typography variant="body1" fontFamily="monospace">
                        {watchedValues.receiver.slice(0, 20)}...
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Amount:</Typography>
                      <Typography variant="h6">
                        {watchedValues.amount.value} {watchedValues.amount.denom.toUpperCase()}
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Priority:</Typography>
                      <Chip 
                        label={watchedValues.priority.toUpperCase()} 
                        size="small" 
                        color="primary"
                      />
                    </Grid>
                  </Grid>

                  {quote && (
                    <Box mt={2}>
                      <Typography variant="subtitle2" gutterBottom>
                        Exchange Details
                      </Typography>
                      <Typography variant="body2">
                        Exchange Rate: {quote.exchangeRate} INR per NAMO
                      </Typography>
                      <Typography variant="body2">
                        Network Fee: {quote.fee.amount} {quote.fee.denom.toUpperCase()}
                      </Typography>
                      <Typography variant="body2">
                        Price Impact: {(quote.priceImpact * 100).toFixed(2)}%
                      </Typography>
                      {quote.culturalBonus && (
                        <Typography variant="body2" color="success.main">
                          Cultural Bonus: +{quote.culturalBonus.amount} {quote.culturalBonus.denom.toUpperCase()}
                        </Typography>
                      )}
                    </Box>
                  )}
                </CardContent>
              </Card>
            </Grid>

            {showCulturalFeatures && currentQuote && watchedValues.culturalPreferences.includeQuote && (
              <Grid item xs={12}>
                <CulturalQuote quote={currentQuote} showFullCard />
              </Grid>
            )}

            {isPatriotismEnabled && (
              <Grid item xs={12}>
                <Alert severity="info" icon={<StarIcon />}>
                  This transaction will contribute to your patriotism score and may earn additional rewards!
                </Alert>
              </Grid>
            )}
          </Grid>
        );

      default:
        return null;
    }
  };

  return (
    <FormProvider {...methods}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Card>
          <CardContent>
            {showCulturalFeatures && currentFestival && (
              <Box mb={3}>
                <FestivalBanner festival={currentFestival} />
              </Box>
            )}

            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h5" component="h1">
                {mode === 'simple' ? 'Send Money Order' : 'Advanced Trading'}
              </Typography>
              
              {mode === 'simple' && (
                <Button
                  variant="outlined"
                  size="small"
                  onClick={() => setShowAdvanced(!showAdvanced)}
                >
                  {showAdvanced ? 'Simple Mode' : 'Advanced Mode'}
                </Button>
              )}
            </Box>

            <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 4 }}>
              {steps.map((label) => (
                <Step key={label}>
                  <StepLabel>{label}</StepLabel>
                </Step>
              ))}
            </Stepper>

            <form onSubmit={handleSubmit(onFormSubmit)}>
              <AnimatePresence mode="wait">
                <motion.div
                  key={activeStep}
                  initial={{ opacity: 0, x: 50 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -50 }}
                  transition={{ duration: 0.2 }}
                >
                  {getStepContent(activeStep)}
                </motion.div>
              </AnimatePresence>

              <Box display="flex" justifyContent="space-between" mt={4}>
                <Button
                  disabled={activeStep === 0}
                  onClick={handleBack}
                  variant="outlined"
                >
                  Back
                </Button>

                <Box>
                  {activeStep === steps.length - 1 ? (
                    <Button
                      type="submit"
                      variant="contained"
                      disabled={!isValid || isSubmitting}
                      startIcon={isSubmitting ? <LoadingSpinner size={20} /> : <SendIcon />}
                      size="large"
                    >
                      {isSubmitting ? 'Processing...' : 'Send Money Order'}
                    </Button>
                  ) : (
                    <Button
                      variant="contained"
                      onClick={handleNext}
                      disabled={activeStep === 0 && (!watchedValues.sender || !watchedValues.receiver || !watchedValues.amount.value)}
                      startIcon={<ReceiptIcon />}
                    >
                      Next
                    </Button>
                  )}
                </Box>
              </Box>
            </form>

            <Collapse in={showAdvanced && mode === 'simple'}>
              <Box mt={3} p={2} bgcolor="background.paper" borderRadius={1}>
                <Typography variant="subtitle2" gutterBottom>
                  Advanced Options
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={6}>
                    <Typography variant="body2">
                      Max Slippage: 5%
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <Typography variant="body2">
                      Gas Price: Auto
                    </Typography>
                  </Grid>
                </Grid>
              </Box>
            </Collapse>

            {Object.keys(errors).length > 0 && (
              <Alert severity="error" sx={{ mt: 2 }}>
                Please fix the following errors:
                <ul>
                  {Object.entries(errors).map(([field, error]) => (
                    <li key={field}>{error.message}</li>
                  ))}
                </ul>
              </Alert>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </FormProvider>
  );
};