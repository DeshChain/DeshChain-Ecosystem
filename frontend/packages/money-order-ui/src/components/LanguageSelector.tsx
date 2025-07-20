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

import React, { useState, useRef, useEffect } from 'react';
import {
  Box,
  Button,
  Menu,
  MenuItem,
  Typography,
  Avatar,
  Chip,
  TextField,
  InputAdornment,
  Divider,
  ListItemIcon,
  ListItemText,
  Badge,
  useTheme,
  alpha,
  Tooltip
} from '@mui/material';
import {
  Language as LanguageIcon,
  Search as SearchIcon,
  Check as CheckIcon,
  Translate as TranslateIcon,
  Public as PublicIcon
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { motion, AnimatePresence } from 'framer-motion';

import { SUPPORTED_LANGUAGES, changeLanguage, getCurrentLanguage } from '../i18n';
import { LANGUAGE_METADATA, LanguageService } from '../i18n/LanguageService';

interface LanguageSelectorProps {
  variant?: 'button' | 'menu' | 'compact';
  showNativeName?: boolean;
  showSpeakers?: boolean;
  showRegion?: boolean;
  groupByScript?: boolean;
}

export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
  variant = 'button',
  showNativeName = true,
  showSpeakers = false,
  showRegion = false,
  groupByScript = true
}) => {
  const theme = useTheme();
  const { t } = useTranslation();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedLanguage, setSelectedLanguage] = useState(getCurrentLanguage());
  const searchInputRef = useRef<HTMLInputElement>(null);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
    setSearchQuery('');
  };

  const handleLanguageSelect = async (languageCode: string) => {
    await changeLanguage(languageCode);
    setSelectedLanguage(languageCode);
    handleClose();
    
    // Show success notification
    const selectedLang = SUPPORTED_LANGUAGES.find(lang => lang.code === languageCode);
    if (selectedLang) {
      // You can integrate with your notification system here
      console.log(`Language changed to ${selectedLang.nativeName}`);
    }
  };

  // Filter languages based on search
  const filteredLanguages = SUPPORTED_LANGUAGES.filter(lang => {
    const query = searchQuery.toLowerCase();
    return (
      lang.name.toLowerCase().includes(query) ||
      lang.nativeName.toLowerCase().includes(query) ||
      lang.code.toLowerCase().includes(query) ||
      LANGUAGE_METADATA[lang.code]?.region.some(r => r.toLowerCase().includes(query))
    );
  });

  // Group languages by script if enabled
  const groupedLanguages = groupByScript
    ? filteredLanguages.reduce((acc, lang) => {
        const script = LANGUAGE_METADATA[lang.code]?.script || 'Other';
        if (!acc[script]) acc[script] = [];
        acc[script].push(lang);
        return acc;
      }, {} as Record<string, typeof SUPPORTED_LANGUAGES>)
    : { All: filteredLanguages };

  // Get current language info
  const currentLang = SUPPORTED_LANGUAGES.find(lang => lang.code === selectedLanguage);
  const currentLangMeta = LANGUAGE_METADATA[selectedLanguage];

  // Focus search input when menu opens
  useEffect(() => {
    if (anchorEl && searchInputRef.current) {
      setTimeout(() => searchInputRef.current?.focus(), 100);
    }
  }, [anchorEl]);

  const renderLanguageButton = () => {
    switch (variant) {
      case 'compact':
        return (
          <Tooltip title={t('settings.language')}>
            <IconButton onClick={handleClick} size="small">
              <Badge
                badgeContent={currentLang?.code.toUpperCase()}
                color="primary"
                sx={{
                  '& .MuiBadge-badge': {
                    fontSize: '0.6rem',
                    height: 16,
                    minWidth: 16
                  }
                }}
              >
                <LanguageIcon />
              </Badge>
            </IconButton>
          </Tooltip>
        );

      case 'menu':
        return (
          <MenuItem onClick={handleClick}>
            <ListItemIcon>
              <TranslateIcon />
            </ListItemIcon>
            <ListItemText
              primary={t('settings.language')}
              secondary={currentLang?.nativeName}
            />
          </MenuItem>
        );

      default:
        return (
          <Button
            variant="outlined"
            startIcon={<LanguageIcon />}
            onClick={handleClick}
            sx={{
              borderRadius: 2,
              textTransform: 'none',
              minWidth: 120
            }}
          >
            {showNativeName ? currentLang?.nativeName : currentLang?.name}
          </Button>
        );
    }
  };

  return (
    <>
      {renderLanguageButton()}

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        PaperProps={{
          sx: {
            width: 360,
            maxHeight: 500,
            overflow: 'hidden',
            display: 'flex',
            flexDirection: 'column'
          }
        }}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        {/* Header */}
        <Box px={2} py={1.5} borderBottom={1} borderColor="divider">
          <Typography variant="h6" gutterBottom>
            {t('settings.language')}
          </Typography>
          <TextField
            ref={searchInputRef}
            fullWidth
            size="small"
            placeholder="Search languages..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon fontSize="small" />
                </InputAdornment>
              )
            }}
          />
        </Box>

        {/* Language List */}
        <Box sx={{ flex: 1, overflowY: 'auto' }}>
          <AnimatePresence>
            {Object.entries(groupedLanguages).map(([script, languages]) => (
              <motion.div
                key={script}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
              >
                {groupByScript && Object.keys(groupedLanguages).length > 1 && (
                  <Box px={2} py={0.5} bgcolor="grey.100">
                    <Typography variant="caption" color="text.secondary">
                      {script} Script
                    </Typography>
                  </Box>
                )}

                {languages.map((language) => {
                  const langMeta = LANGUAGE_METADATA[language.code];
                  const isSelected = language.code === selectedLanguage;
                  const isRTL = language.direction === 'rtl';

                  return (
                    <MenuItem
                      key={language.code}
                      onClick={() => handleLanguageSelect(language.code)}
                      selected={isSelected}
                      sx={{
                        py: 1.5,
                        direction: isRTL ? 'rtl' : 'ltr',
                        '&:hover': {
                          bgcolor: alpha(theme.palette.primary.main, 0.08)
                        }
                      }}
                    >
                      <ListItemIcon>
                        {isSelected ? (
                          <Avatar
                            sx={{
                              width: 32,
                              height: 32,
                              bgcolor: 'primary.main',
                              fontSize: '0.875rem'
                            }}
                          >
                            <CheckIcon fontSize="small" />
                          </Avatar>
                        ) : (
                          <Avatar
                            sx={{
                              width: 32,
                              height: 32,
                              bgcolor: 'grey.200',
                              color: 'text.primary',
                              fontSize: '0.875rem'
                            }}
                          >
                            {language.code.toUpperCase()}
                          </Avatar>
                        )}
                      </ListItemIcon>

                      <ListItemText
                        primary={
                          <Box display="flex" alignItems="center" gap={1}>
                            <Typography variant="body1">
                              {language.nativeName}
                            </Typography>
                            {langMeta?.official && (
                              <Chip
                                label="Official"
                                size="small"
                                color="primary"
                                variant="outlined"
                                sx={{ height: 20, fontSize: '0.7rem' }}
                              />
                            )}
                          </Box>
                        }
                        secondary={
                          <Box>
                            <Typography variant="caption" component="div">
                              {language.name}
                            </Typography>
                            {showRegion && langMeta?.region && (
                              <Typography variant="caption" color="text.disabled">
                                {langMeta.region.join(', ')}
                              </Typography>
                            )}
                            {showSpeakers && langMeta?.speakers && (
                              <Typography variant="caption" color="text.disabled">
                                {LanguageService.formatIndianNumber(langMeta.speakers, 'en')} speakers
                              </Typography>
                            )}
                          </Box>
                        }
                      />

                      {isRTL && (
                        <Chip
                          label="RTL"
                          size="small"
                          sx={{
                            height: 18,
                            fontSize: '0.65rem',
                            ml: 1
                          }}
                        />
                      )}
                    </MenuItem>
                  );
                })}
              </motion.div>
            ))}
          </AnimatePresence>

          {filteredLanguages.length === 0 && (
            <Box p={3} textAlign="center">
              <PublicIcon sx={{ fontSize: 48, color: 'text.disabled', mb: 1 }} />
              <Typography color="text.secondary">
                No languages found
              </Typography>
            </Box>
          )}
        </Box>

        {/* Footer */}
        <Box px={2} py={1} borderTop={1} borderColor="divider" bgcolor="grey.50">
          <Typography variant="caption" color="text.secondary" align="center">
            22 Official Languages of India 🇮🇳
          </Typography>
        </Box>
      </Menu>
    </>
  );
};

// Mini language selector for inline use
export const MiniLanguageSelector: React.FC = () => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const currentLang = SUPPORTED_LANGUAGES.find(lang => lang.code === getCurrentLanguage());

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLanguageSelect = async (languageCode: string) => {
    await changeLanguage(languageCode);
    handleClose();
  };

  return (
    <>
      <Chip
        label={currentLang?.code.toUpperCase()}
        onClick={handleClick}
        size="small"
        icon={<LanguageIcon />}
        sx={{ cursor: 'pointer' }}
      />

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        PaperProps={{ sx: { maxHeight: 300 } }}
      >
        {SUPPORTED_LANGUAGES.map((language) => (
          <MenuItem
            key={language.code}
            onClick={() => handleLanguageSelect(language.code)}
            selected={language.code === getCurrentLanguage()}
            dense
          >
            <Typography variant="body2">
              {language.nativeName} ({language.name})
            </Typography>
          </MenuItem>
        ))}
      </Menu>
    </>
  );
};

// Language indicator for status bar
export const LanguageIndicator: React.FC = () => {
  const currentLang = SUPPORTED_LANGUAGES.find(lang => lang.code === getCurrentLanguage());
  const theme = useTheme();

  return (
    <Box
      display="flex"
      alignItems="center"
      gap={0.5}
      px={1}
      py={0.5}
      borderRadius={1}
      bgcolor={alpha(theme.palette.primary.main, 0.1)}
    >
      <LanguageIcon sx={{ fontSize: 16, color: 'primary.main' }} />
      <Typography variant="caption" color="primary.main" fontWeight="medium">
        {currentLang?.code.toUpperCase()}
      </Typography>
    </Box>
  );
};