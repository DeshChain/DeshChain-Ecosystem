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
  Card,
  CardContent,
  Typography,
  IconButton,
  Tooltip,
  Chip,
  Fade,
  useTheme
} from '@mui/material';
import {
  FormatQuote as QuoteIcon,
  Share as ShareIcon,
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { CulturalQuoteData } from '../../types';
import { useCulturalContext } from '../../hooks/useCulturalContext';

interface CulturalQuoteProps {
  quote: CulturalQuoteData;
  showFullCard?: boolean;
  showActions?: boolean;
  compact?: boolean;
  onFavorite?: (quoteId: string) => void;
  onShare?: (quote: CulturalQuoteData) => void;
  isFavorited?: boolean;
  variant?: 'card' | 'banner' | 'inline';
  animation?: 'fade' | 'slide' | 'scale' | 'none';
}

export const CulturalQuote: React.FC<CulturalQuoteProps> = ({
  quote,
  showFullCard = false,
  showActions = true,
  compact = false,
  onFavorite,
  onShare,
  isFavorited = false,
  variant = 'card',
  animation = 'fade'
}) => {
  const theme = useTheme();
  const { currentFestival, getLanguageNativeName } = useCulturalContext();

  const getQuoteGradient = () => {
    if (currentFestival && theme.palette.mode === 'light') {
      const festivalColors = theme.palette.primary;
      return `linear-gradient(135deg, ${festivalColors.light}20, ${festivalColors.main}20)`;
    }
    return theme.palette.mode === 'dark' 
      ? 'linear-gradient(135deg, rgba(255,255,255,0.05), rgba(255,255,255,0.1))'
      : 'linear-gradient(135deg, rgba(0,0,0,0.02), rgba(0,0,0,0.05))';
  };

  const getCategoryColor = (category: string) => {
    const colors = {
      wisdom: 'primary',
      motivation: 'success',
      patriotism: 'error',
      prosperity: 'warning',
      family: 'info',
      community: 'secondary',
      tradition: 'primary',
      spirituality: 'secondary'
    } as const;
    return colors[category as keyof typeof colors] || 'default';
  };

  const handleShare = () => {
    if (onShare) {
      onShare(quote);
    } else if (navigator.share) {
      navigator.share({
        title: `Quote by ${quote.author}`,
        text: `"${quote.text}" - ${quote.author}`,
        url: window.location.href
      });
    } else {
      // Fallback to clipboard
      navigator.clipboard.writeText(`"${quote.text}" - ${quote.author}`);
    }
  };

  const QuoteContent = () => (
    <Box>
      <Box display="flex" alignItems="flex-start" gap={1} mb={2}>
        <QuoteIcon 
          sx={{ 
            color: 'text.secondary', 
            fontSize: compact ? 16 : 20,
            mt: 0.5,
            transform: 'scaleX(-1)'
          }} 
        />
        <Box flex={1}>
          <Typography
            variant={compact ? 'body2' : 'body1'}
            sx={{
              fontStyle: 'italic',
              lineHeight: 1.6,
              fontFamily: quote.language === 'hindi' || quote.language === 'sanskrit' 
                ? 'Noto Sans Devanagari, serif' 
                : 'inherit'
            }}
          >
            {quote.text}
          </Typography>
          
          {quote.translation && quote.translation !== quote.text && (
            <Typography
              variant="caption"
              display="block"
              sx={{ 
                mt: 1, 
                color: 'text.secondary',
                fontStyle: 'normal'
              }}
            >
              {quote.translation}
            </Typography>
          )}
        </Box>
        <QuoteIcon 
          sx={{ 
            color: 'text.secondary', 
            fontSize: compact ? 16 : 20,
            mt: 0.5
          }} 
        />
      </Box>

      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Typography
            variant={compact ? 'caption' : 'subtitle2'}
            sx={{ fontWeight: 'medium' }}
          >
            — {quote.author}
          </Typography>
          
          {showFullCard && (
            <Box display="flex" gap={1} mt={1} flexWrap="wrap">
              <Chip
                label={quote.category}
                size="small"
                color={getCategoryColor(quote.category)}
                variant="outlined"
              />
              
              {quote.language !== 'en' && (
                <Chip
                  label={getLanguageNativeName(quote.language)}
                  size="small"
                  variant="outlined"
                />
              )}
              
              {quote.occasion !== 'general' && (
                <Chip
                  label={quote.occasion}
                  size="small"
                  color="secondary"
                  variant="outlined"
                />
              )}
            </Box>
          )}
        </Box>

        {showActions && (
          <Box display="flex" gap={0.5}>
            <Tooltip title={isFavorited ? 'Remove from favorites' : 'Add to favorites'}>
              <IconButton
                size="small"
                onClick={() => onFavorite?.(quote.quoteId)}
                color={isFavorited ? 'error' : 'default'}
              >
                {isFavorited ? <FavoriteIcon /> : <FavoriteBorderIcon />}
              </IconButton>
            </Tooltip>
            
            <Tooltip title="Share quote">
              <IconButton size="small" onClick={handleShare}>
                <ShareIcon />
              </IconButton>
            </Tooltip>
            
            {quote.context && (
              <Tooltip title={quote.context}>
                <IconButton size="small">
                  <InfoIcon />
                </IconButton>
              </Tooltip>
            )}
          </Box>
        )}
      </Box>
    </Box>
  );

  const getAnimationProps = () => {
    switch (animation) {
      case 'slide':
        return {
          initial: { opacity: 0, x: -20 },
          animate: { opacity: 1, x: 0 },
          transition: { duration: 0.3 }
        };
      case 'scale':
        return {
          initial: { opacity: 0, scale: 0.95 },
          animate: { opacity: 1, scale: 1 },
          transition: { duration: 0.3 }
        };
      case 'fade':
        return {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: 0.3 }
        };
      default:
        return {};
    }
  };

  if (variant === 'inline') {
    return (
      <motion.div {...(animation !== 'none' ? getAnimationProps() : {})}>
        <Box
          sx={{
            p: compact ? 1 : 2,
            borderLeft: 3,
            borderColor: 'primary.main',
            bgcolor: getQuoteGradient(),
            borderRadius: 1
          }}
        >
          <QuoteContent />
        </Box>
      </motion.div>
    );
  }

  if (variant === 'banner') {
    return (
      <motion.div {...(animation !== 'none' ? getAnimationProps() : {})}>
        <Box
          sx={{
            p: 2,
            background: getQuoteGradient(),
            borderRadius: 2,
            border: 1,
            borderColor: 'divider',
            textAlign: 'center'
          }}
        >
          <QuoteContent />
        </Box>
      </motion.div>
    );
  }

  return (
    <motion.div {...(animation !== 'none' ? getAnimationProps() : {})}>
      <Card
        sx={{
          background: getQuoteGradient(),
          border: 1,
          borderColor: 'divider',
          '&:hover': {
            boxShadow: theme.shadows[4],
            transform: 'translateY(-2px)'
          },
          transition: 'all 0.2s ease-in-out'
        }}
      >
        <CardContent sx={{ p: compact ? 2 : 3 }}>
          <QuoteContent />
        </CardContent>
      </Card>
    </motion.div>
  );
};