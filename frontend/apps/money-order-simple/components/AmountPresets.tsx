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
  Chip,
  Typography
} from '@mui/material';
import { motion } from 'framer-motion';

interface AmountPresetsProps {
  onSelect: (amount: string) => void;
  customPresets?: Array<{ label: string; value: string }>;
}

export const AmountPresets: React.FC<AmountPresetsProps> = ({ 
  onSelect, 
  customPresets 
}) => {
  const defaultPresets = [
    { label: '₹100', value: '1333' }, // ~100 INR at 0.075 rate
    { label: '₹500', value: '6667' },
    { label: '₹1000', value: '13333' },
    { label: '₹5000', value: '66667' },
    { label: '₹10000', value: '133333' }
  ];

  const presets = customPresets || defaultPresets;

  return (
    <Box>
      <Typography variant="caption" color="text.secondary" gutterBottom>
        Quick amounts:
      </Typography>
      <Box display="flex" gap={1} flexWrap="wrap" mt={0.5}>
        {presets.map((preset, index) => (
          <motion.div
            key={preset.value}
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.2, delay: index * 0.05 }}
          >
            <Chip
              label={preset.label}
              onClick={() => onSelect(preset.value)}
              variant="outlined"
              size="small"
              sx={{
                cursor: 'pointer',
                '&:hover': {
                  bgcolor: 'primary.main',
                  color: 'white',
                  borderColor: 'primary.main'
                }
              }}
            />
          </motion.div>
        ))}
      </Box>
    </Box>
  );
};