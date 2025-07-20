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
  Container,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Button,
  Grid,
  Card,
  CardContent,
  Fab,
  useTheme
} from '@mui/material';
import {
  Menu as MenuIcon,
  AccountBalanceWallet as WalletIcon,
  Send as SendIcon,
  Settings as SettingsIcon,
  Celebration as CelebrationIcon
} from '@mui/icons-material';

import { FestivalThemeProvider } from '../themes/FestivalThemeProvider';
import { FestivalBanner } from '../components/festival/FestivalBanner';
import { FestivalMoneyOrderBonus } from '../components/festival/FestivalMoneyOrderBonus';
import { LanguageSelector } from '../components/LanguageSelector';
import { LocalizedMoneyOrderForm } from '../components/LocalizedMoneyOrderForm';
import { useFestival } from '../hooks/useFestival';
import { useLanguage } from '../hooks/useLanguage';

// Main app with festival theming
const FestivalThemedAppContent: React.FC = () => {
  const theme = useTheme();
  const { currentFestival, isFestivalActive } = useFestival();
  const { formatCurrency, getGreeting } = useLanguage();

  return (
    <Box sx={{ flexGrow: 1, minHeight: '100vh' }}>
      {/* App Bar */}
      <AppBar 
        position="sticky" 
        className={isFestivalActive ? 'festival-glow' : ''}
        sx={{
          background: isFestivalActive && currentFestival
            ? currentFestival.theme.backgroundGradient
            : undefined
        }}
      >
        <Toolbar>
          <IconButton edge="start" color="inherit" sx={{ mr: 2 }}>
            <MenuIcon />
          </IconButton>
          
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            DeshChain Money Order
          </Typography>
          
          <Box display="flex" alignItems="center" gap={2}>
            <Button
              color="inherit"
              startIcon={<WalletIcon />}
              sx={{ display: { xs: 'none', sm: 'flex' } }}
            >
              ₹25,000
            </Button>
            
            <LanguageSelector variant="compact" />
            
            <IconButton color="inherit">
              <SettingsIcon />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      {/* Festival Banner */}
      <Container maxWidth="lg" sx={{ mt: 2 }}>
        <FestivalBanner variant="compact" />
      </Container>

      {/* Main Content */}
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={3}>
          {/* Welcome Card */}
          <Grid item xs={12}>
            <Card className={isFestivalActive ? 'festival-card' : ''}>
              <CardContent>
                <Typography variant="h4" gutterBottom>
                  {getGreeting()}! 👋
                </Typography>
                <Typography variant="body1" color="textSecondary">
                  Send money across India instantly with blockchain security.
                  {isFestivalActive && currentFestival && (
                    <Typography
                      component="span"
                      color="primary"
                      fontWeight="bold"
                      sx={{ ml: 1 }}
                    >
                      Special {currentFestival.name} offers active! 🎉
                    </Typography>
                  )}
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Quick Actions */}
          <Grid item xs={12} md={4}>
            <Card
              className={isFestivalActive ? 'festival-sparkle' : ''}
              sx={{ height: '100%', cursor: 'pointer' }}
            >
              <CardContent>
                <Box display="flex" alignItems="center" gap={2} mb={2}>
                  <SendIcon sx={{ fontSize: 40, color: 'primary.main' }} />
                  <Typography variant="h6">Send Money</Typography>
                </Box>
                <Typography variant="body2" color="textSecondary">
                  Transfer money instantly to any wallet or bank account
                </Typography>
                {isFestivalActive && (
                  <Box mt={2}>
                    <FestivalBonusChip amount={1000} />
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', cursor: 'pointer' }}>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2} mb={2}>
                  <WalletIcon sx={{ fontSize: 40, color: 'secondary.main' }} />
                  <Typography variant="h6">My Wallet</Typography>
                </Box>
                <Typography variant="body2" color="textSecondary">
                  Manage your funds and view transaction history
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', cursor: 'pointer' }}>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2} mb={2}>
                  <CelebrationIcon sx={{ fontSize: 40, color: 'error.main' }} />
                  <Typography variant="h6">Festival Offers</Typography>
                </Box>
                <Typography variant="body2" color="textSecondary">
                  Explore special festival discounts and bonuses
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Money Order Form */}
          <Grid item xs={12} lg={8}>
            <LocalizedMoneyOrderForm 
              onSubmit={(data) => console.log('Money order submitted:', data)}
            />
          </Grid>

          {/* Festival Bonus Card */}
          <Grid item xs={12} lg={4}>
            {isFestivalActive && (
              <FestivalMoneyOrderBonus 
                amount={5000}
                showDetails={true}
              />
            )}
          </Grid>
        </Grid>
      </Container>

      {/* Festival FAB */}
      {isFestivalActive && (
        <Fab
          color="primary"
          sx={{
            position: 'fixed',
            bottom: 24,
            right: 24,
            animation: 'dance 2s ease-in-out infinite'
          }}
        >
          <CelebrationIcon />
        </Fab>
      )}
    </Box>
  );
};

// Festival bonus chip component
const FestivalBonusChip: React.FC<{ amount: number }> = ({ amount }) => {
  const { calculateFestivalBonus } = useFestival();
  const bonus = calculateFestivalBonus(amount);
  
  if (bonus.totalBenefit === 0) {
    return null;
  }
  
  return (
    <Box
      sx={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 1,
        px: 2,
        py: 1,
        borderRadius: 2,
        bgcolor: 'success.light',
        color: 'success.contrastText'
      }}
    >
      <CelebrationIcon fontSize="small" />
      <Typography variant="body2" fontWeight="bold">
        Save ₹{bonus.totalBenefit.toFixed(0)}!
      </Typography>
    </Box>
  );
};

// Export wrapped with provider
export const FestivalThemedApp: React.FC = () => {
  return (
    <FestivalThemeProvider>
      <FestivalThemedAppContent />
    </FestivalThemeProvider>
  );
};