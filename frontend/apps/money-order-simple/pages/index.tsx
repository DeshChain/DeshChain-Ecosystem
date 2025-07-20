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
import { useRouter } from 'next/router';
import {
  Container,
  Grid,
  Typography,
  Box,
  Card,
  CardContent,
  Button,
  Chip,
  Alert,
  useTheme,
  useMediaQuery
} from '@mui/material';
import {
  Send as SendIcon,
  AccountBalance as BankIcon,
  Speed as SpeedIcon,
  Security as SecurityIcon,
  Language as LanguageIcon,
  Celebration as CelebrationIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import Confetti from 'react-confetti';

import { SimpleMoneyOrderForm } from '../components/SimpleMoneyOrderForm';
import { RecentTransactions } from '../components/RecentTransactions';
import { QuickActions } from '../components/QuickActions';
import { PoolStatusBar } from '../components/PoolStatusBar';
import { FestivalWidget } from '../components/FestivalWidget';
import { useWindowSize } from '../hooks/useWindowSize';
import { useCulturalContext } from '@deshchain/money-order-ui';

const Home: React.FC = () => {
  const router = useRouter();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const { width, height } = useWindowSize();
  const { currentFestival, currentLanguage } = useCulturalContext();

  const [showConfetti, setShowConfetti] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const handleSuccess = (receiptId: string) => {
    setSuccessMessage(`Money Order sent successfully! Receipt: ${receiptId}`);
    setShowConfetti(true);
    
    // Navigate to success page after a short delay
    setTimeout(() => {
      router.push(`/receipt/${receiptId}`);
    }, 3000);
  };

  const features = [
    {
      icon: <SpeedIcon />,
      title: 'Instant Transfers',
      description: 'Send money in seconds with our optimized network',
      color: theme.palette.primary.main
    },
    {
      icon: <SecurityIcon />,
      title: 'Secure & Transparent',
      description: 'Blockchain-secured with complete transaction transparency',
      color: theme.palette.secondary.main
    },
    {
      icon: <LanguageIcon />,
      title: '22 Languages',
      description: 'Use your preferred Indian language with native script support',
      color: theme.palette.info.main
    },
    {
      icon: <CelebrationIcon />,
      title: 'Festival Bonuses',
      description: 'Get special rewards during Indian festivals',
      color: theme.palette.warning.main
    }
  ];

  return (
    <>
      {showConfetti && (
        <Confetti
          width={width}
          height={height}
          recycle={false}
          numberOfPieces={200}
          gravity={0.2}
        />
      )}

      <Container maxWidth="lg" sx={{ py: 3 }}>
        {/* Hero Section */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <Box textAlign="center" mb={5}>
            <Typography variant="h1" gutterBottom>
              Send Money Orders
              <Chip
                label="Simple & Cultural"
                color="primary"
                sx={{ ml: 2, verticalAlign: 'middle' }}
              />
            </Typography>
            <Typography variant="h5" color="text.secondary" paragraph>
              Transfer money instantly with cultural quotes and festival bonuses
            </Typography>
          </Box>
        </motion.div>

        {/* Festival Widget */}
        {currentFestival && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.3 }}
          >
            <Box mb={4}>
              <FestivalWidget />
            </Box>
          </motion.div>
        )}

        {/* Pool Status Bar */}
        <Box mb={4}>
          <PoolStatusBar />
        </Box>

        {/* Main Content Grid */}
        <Grid container spacing={4}>
          {/* Left Column - Money Order Form */}
          <Grid item xs={12} md={8}>
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
            >
              <SimpleMoneyOrderForm onSuccess={handleSuccess} />
            </motion.div>

            {/* Features Grid */}
            <Grid container spacing={2} sx={{ mt: 3 }}>
              {features.map((feature, index) => (
                <Grid item xs={12} sm={6} key={index}>
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5, delay: 0.2 + index * 0.1 }}
                  >
                    <Card
                      sx={{
                        height: '100%',
                        transition: 'transform 0.2s',
                        '&:hover': {
                          transform: 'translateY(-4px)'
                        }
                      }}
                    >
                      <CardContent>
                        <Box display="flex" alignItems="center" mb={1}>
                          <Box
                            sx={{
                              p: 1,
                              borderRadius: 2,
                              bgcolor: `${feature.color}20`,
                              color: feature.color,
                              mr: 2
                            }}
                          >
                            {feature.icon}
                          </Box>
                          <Typography variant="h6">
                            {feature.title}
                          </Typography>
                        </Box>
                        <Typography variant="body2" color="text.secondary">
                          {feature.description}
                        </Typography>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              ))}
            </Grid>
          </Grid>

          {/* Right Column - Quick Actions & Recent */}
          <Grid item xs={12} md={4}>
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
            >
              <Box sx={{ position: 'sticky', top: 20 }}>
                <QuickActions />
                
                <Box mt={3}>
                  <RecentTransactions limit={5} />
                </Box>

                {/* Help Card */}
                <Card sx={{ mt: 3 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Need Help?
                    </Typography>
                    <Typography variant="body2" color="text.secondary" paragraph>
                      Our Money Order system is designed to be simple and intuitive.
                    </Typography>
                    <Button
                      variant="outlined"
                      fullWidth
                      onClick={() => router.push('/help')}
                    >
                      View Help Guide
                    </Button>
                  </CardContent>
                </Card>
              </Box>
            </motion.div>
          </Grid>
        </Grid>

        {/* Success Message */}
        {successMessage && (
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            style={{
              position: 'fixed',
              bottom: 20,
              left: '50%',
              transform: 'translateX(-50%)',
              zIndex: 1000
            }}
          >
            <Alert
              severity="success"
              sx={{
                boxShadow: theme.shadows[8],
                minWidth: 300
              }}
            >
              {successMessage}
            </Alert>
          </motion.div>
        )}

        {/* Cultural Quote at Bottom */}
        <Box mt={6} textAlign="center">
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{ fontStyle: 'italic' }}
          >
            "सर्वे भवन्तु सुखिनः - May all beings be happy"
          </Typography>
        </Box>
      </Container>
    </>
  );
};