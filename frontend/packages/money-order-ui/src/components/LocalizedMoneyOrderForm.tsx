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

import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Button,
  TextField,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  InputAdornment,
  Alert,
  Divider,
  Chip,
  useTheme
} from '@mui/material';
import {
  Send as SendIcon,
  AccountBalance as AccountIcon,
  Phone as PhoneIcon,
  LocationOn as LocationIcon,
  Translate as TranslateIcon
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { motion } from 'framer-motion';

import { useLanguage, useLocalizedNumbers, useRTL } from '../hooks/useLanguage';
import { LanguageSelector } from './LanguageSelector';
import { MoneyOrderFormData } from '../types';
import { CulturalQuote } from './cultural/CulturalQuote';

interface LocalizedMoneyOrderFormProps {
  onSubmit: (data: MoneyOrderFormData) => void;
  loading?: boolean;
}

export const LocalizedMoneyOrderForm: React.FC<LocalizedMoneyOrderFormProps> = ({
  onSubmit,
  loading = false
}) => {
  const theme = useTheme();
  const { t } = useTranslation();
  const { 
    currentLanguage, 
    formatCurrency, 
    formatIndianNumber,
    validatePhoneNumber,
    getGreeting,
    getCulturalPhrase
  } = useLanguage();
  const { localizeInput, formatForDisplay } = useLocalizedNumbers();
  const { isRTL, textAlign, direction } = useRTL();

  const [formData, setFormData] = useState<Partial<MoneyOrderFormData>>({
    amount: '',
    purpose: 'family_support'
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const normalizedValue = localizeInput(value);
    
    // Validate amount
    const numValue = parseFloat(normalizedValue);
    if (isNaN(numValue) || numValue < 10) {
      setErrors({ ...errors, amount: t('moneyOrder.validation.minAmount') });
    } else if (numValue > 100000) {
      setErrors({ ...errors, amount: t('moneyOrder.validation.maxAmount') });
    } else {
      const { amount, ...restErrors } = errors;
      setErrors(restErrors);
    }
    
    setFormData({ ...formData, amount: normalizedValue });
  };

  const handlePhoneChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    
    if (value && !validatePhoneNumber(value)) {
      setErrors({ ...errors, phone: t('moneyOrder.validation.invalidPhone') });
    } else {
      const { phone, ...restErrors } = errors;
      setErrors(restErrors);
    }
    
    setFormData({ ...formData, receiverPhone: value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (Object.keys(errors).length === 0) {
      onSubmit(formData as MoneyOrderFormData);
    }
  };

  return (
    <Box dir={direction}>
      {/* Language Selector in Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="bold">
          {t('moneyOrder.title')}
        </Typography>
        <LanguageSelector />
      </Box>

      {/* Greeting */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Alert 
          severity="info" 
          icon={<TranslateIcon />}
          sx={{ mb: 3, textAlign }}
        >
          <Typography variant="body1">
            {getGreeting()}, {t('app.tagline')}
          </Typography>
        </Alert>
      </motion.div>

      <Card>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              {/* Sender Section */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom sx={{ textAlign }}>
                  {t('moneyOrder.form.sender')}
                </Typography>
                <Divider sx={{ mb: 2 }} />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.senderAddress')}
                  value={formData.senderAddress || ''}
                  onChange={(e) => setFormData({ ...formData, senderAddress: e.target.value })}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position={isRTL ? 'end' : 'start'}>
                        <AccountIcon />
                      </InputAdornment>
                    )
                  }}
                  dir={direction}
                />
              </Grid>

              {/* Receiver Section */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom sx={{ textAlign, mt: 2 }}>
                  {t('moneyOrder.form.receiver')}
                </Typography>
                <Divider sx={{ mb: 2 }} />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.receiverName')}
                  value={formData.receiverName || ''}
                  onChange={(e) => setFormData({ ...formData, receiverName: e.target.value })}
                  required
                  dir={direction}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.receiverPhone')}
                  value={formData.receiverPhone || ''}
                  onChange={handlePhoneChange}
                  error={!!errors.phone}
                  helperText={errors.phone}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position={isRTL ? 'end' : 'start'}>
                        <PhoneIcon />
                      </InputAdornment>
                    )
                  }}
                  required
                  dir={direction}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.receiverAddress')}
                  value={formData.receiverAddress || ''}
                  onChange={(e) => setFormData({ ...formData, receiverAddress: e.target.value })}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position={isRTL ? 'end' : 'start'}>
                        <AccountIcon />
                      </InputAdornment>
                    )
                  }}
                  required
                  dir={direction}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.postalCode')}
                  value={formData.postalCode || ''}
                  onChange={(e) => setFormData({ ...formData, postalCode: e.target.value })}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position={isRTL ? 'end' : 'start'}>
                        <LocationIcon />
                      </InputAdornment>
                    )
                  }}
                  dir={direction}
                />
              </Grid>

              {/* Amount Section */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom sx={{ textAlign, mt: 2 }}>
                  {t('moneyOrder.form.amount')}
                </Typography>
                <Divider sx={{ mb: 2 }} />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={t('moneyOrder.form.amountToSend')}
                  value={formData.amount ? formatForDisplay(formData.amount) : ''}
                  onChange={handleAmountChange}
                  error={!!errors.amount}
                  helperText={errors.amount}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position={isRTL ? 'end' : 'start'}>
                        ₹
                      </InputAdornment>
                    )
                  }}
                  required
                  dir={direction}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <InputLabel>{t('moneyOrder.form.selectPurpose')}</InputLabel>
                  <Select
                    value={formData.purpose || 'family_support'}
                    onChange={(e) => setFormData({ ...formData, purpose: e.target.value })}
                    label={t('moneyOrder.form.selectPurpose')}
                  >
                    <MenuItem value="family_support">
                      {t('moneyOrder.purposes.family_support')}
                    </MenuItem>
                    <MenuItem value="education">
                      {t('moneyOrder.purposes.education')}
                    </MenuItem>
                    <MenuItem value="medical">
                      {t('moneyOrder.purposes.medical')}
                    </MenuItem>
                    <MenuItem value="business">
                      {t('moneyOrder.purposes.business')}
                    </MenuItem>
                    <MenuItem value="personal">
                      {t('moneyOrder.purposes.personal')}
                    </MenuItem>
                    <MenuItem value="other">
                      {t('moneyOrder.purposes.other')}
                    </MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12}>
                <TextField
                  fullWidth
                  multiline
                  rows={3}
                  label={t('moneyOrder.form.message')}
                  value={formData.message || ''}
                  onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                  dir={direction}
                />
              </Grid>

              {/* Cultural Quote */}
              <Grid item xs={12}>
                <CulturalQuote
                  quote={{
                    text: getCulturalPhrase('blessing'),
                    author: 'DeshChain',
                    language: currentLanguage
                  }}
                  variant="inline"
                />
              </Grid>

              {/* Action Buttons */}
              <Grid item xs={12}>
                <Box display="flex" gap={2} justifyContent={isRTL ? 'flex-start' : 'flex-end'}>
                  <Button
                    variant="outlined"
                    onClick={() => setFormData({})}
                    disabled={loading}
                  >
                    {t('moneyOrder.buttons.reset')}
                  </Button>
                  <Button
                    variant="contained"
                    type="submit"
                    startIcon={<SendIcon />}
                    disabled={loading || Object.keys(errors).length > 0}
                  >
                    {t('moneyOrder.buttons.send')}
                  </Button>
                </Box>
              </Grid>
            </Grid>
          </form>
        </CardContent>
      </Card>

      {/* Cultural Footer */}
      <Box mt={3} textAlign="center">
        <Typography variant="body2" color="text.secondary">
          {getCulturalPhrase('thanks')}
        </Typography>
        <Typography variant="caption" color="text.secondary">
          {t('common.poweredBy')}
        </Typography>
      </Box>
    </Box>
  );
};