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
import { Box, Typography, CircularProgress } from '@mui/material';
import { motion } from 'framer-motion';

interface CountdownTimerProps {
  expiresAt: Date;
  onExpire?: () => void;
  size?: 'small' | 'medium' | 'large';
  showProgress?: boolean;
}

export const CountdownTimer: React.FC<CountdownTimerProps> = ({
  expiresAt,
  onExpire,
  size = 'medium',
  showProgress = true
}) => {
  const [timeLeft, setTimeLeft] = useState({
    hours: 0,
    minutes: 0,
    seconds: 0,
    totalSeconds: 0,
    progress: 100
  });
  
  useEffect(() => {
    const calculateTimeLeft = () => {
      const now = new Date().getTime();
      const expiry = expiresAt.getTime();
      const difference = expiry - now;
      
      if (difference <= 0) {
        setTimeLeft({
          hours: 0,
          minutes: 0,
          seconds: 0,
          totalSeconds: 0,
          progress: 0
        });
        if (onExpire) onExpire();
        return;
      }
      
      const totalSeconds = Math.floor(difference / 1000);
      const hours = Math.floor(totalSeconds / 3600);
      const minutes = Math.floor((totalSeconds % 3600) / 60);
      const seconds = totalSeconds % 60;
      
      // Calculate progress (assuming 24 hour total)
      const totalDuration = 24 * 60 * 60; // 24 hours in seconds
      const progress = (totalSeconds / totalDuration) * 100;
      
      setTimeLeft({
        hours,
        minutes,
        seconds,
        totalSeconds,
        progress: Math.min(progress, 100)
      });
    };
    
    calculateTimeLeft();
    const timer = setInterval(calculateTimeLeft, 1000);
    
    return () => clearInterval(timer);
  }, [expiresAt, onExpire]);
  
  const formatNumber = (num: number) => num.toString().padStart(2, '0');
  
  const getFontSize = () => {
    switch (size) {
      case 'small': return '1rem';
      case 'large': return '2rem';
      default: return '1.5rem';
    }
  };
  
  const getProgressSize = () => {
    switch (size) {
      case 'small': return 60;
      case 'large': return 120;
      default: return 80;
    }
  };
  
  const getColor = () => {
    if (timeLeft.progress > 50) return 'primary';
    if (timeLeft.progress > 25) return 'warning';
    return 'error';
  };
  
  return (
    <Box display="flex" alignItems="center" gap={2}>
      {showProgress && (
        <Box position="relative" display="inline-flex">
          <CircularProgress
            variant="determinate"
            value={timeLeft.progress}
            size={getProgressSize()}
            color={getColor() as any}
            thickness={3}
          />
          <Box
            sx={{
              top: 0,
              left: 0,
              bottom: 0,
              right: 0,
              position: 'absolute',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Typography 
              variant="caption" 
              component="div" 
              color="text.secondary"
              fontSize={size === 'small' ? '0.7rem' : '0.875rem'}
            >
              {timeLeft.hours}h
            </Typography>
          </Box>
        </Box>
      )}
      
      <Box display="flex" alignItems="baseline" gap={0.5}>
        <motion.div
          key={timeLeft.hours}
          initial={{ y: -10, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.3 }}
        >
          <Typography
            variant="h6"
            fontSize={getFontSize()}
            fontWeight="bold"
            color={getColor() as any}
          >
            {formatNumber(timeLeft.hours)}
          </Typography>
        </motion.div>
        
        <Typography variant="h6" fontSize={getFontSize()}>:</Typography>
        
        <motion.div
          key={timeLeft.minutes}
          initial={{ y: -10, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.3 }}
        >
          <Typography
            variant="h6"
            fontSize={getFontSize()}
            fontWeight="bold"
            color={getColor() as any}
          >
            {formatNumber(timeLeft.minutes)}
          </Typography>
        </motion.div>
        
        <Typography variant="h6" fontSize={getFontSize()}>:</Typography>
        
        <motion.div
          key={timeLeft.seconds}
          initial={{ y: -10, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.3 }}
        >
          <Typography
            variant="h6"
            fontSize={getFontSize()}
            fontWeight="bold"
            color={getColor() as any}
          >
            {formatNumber(timeLeft.seconds)}
          </Typography>
        </motion.div>
      </Box>
      
      {size !== 'small' && (
        <Typography variant="body2" color="text.secondary">
          {timeLeft.hours > 0 ? 'hours' : 'remaining'}
        </Typography>
      )}
    </Box>
  );
};