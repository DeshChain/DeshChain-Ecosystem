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
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Avatar,
  Button,
  Rating,
  TextField,
  InputAdornment,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  Badge,
  Tooltip,
  LinearProgress,
  useTheme,
  alpha
} from '@mui/material';
import {
  Store as StoreIcon,
  LocationOn as LocationIcon,
  Search as SearchIcon,
  Phone as PhoneIcon,
  Schedule as ScheduleIcon,
  Star as StarIcon,
  AccountBalance as CashInIcon,
  Payment as CashOutIcon,
  Language as LanguageIcon,
  Verified as VerifiedIcon,
  Map as MapIcon,
  DirectionsWalk as WalkIcon,
  DirectionsCar as CarIcon,
  FilterList as FilterIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { PostalCodeInput } from './PostalCodeInput';
import { useLanguage } from '../hooks/useLanguage';

interface Agent {
  agentId: string;
  businessName: string;
  address: string;
  postalCode: string;
  district: string;
  phone: string;
  languages: string[];
  services: ('CASH_IN' | 'CASH_OUT' | 'REMITTANCE' | 'BILL_PAYMENT')[];
  operatingHours: {
    day: string;
    open: string;
    close: string;
    isClosed: boolean;
  }[];
  rating: number;
  totalRatings: number;
  totalTransactions: number;
  dailyLimit: { amount: string; denom: string };
  currentUtilization: number; // percentage
  commissionRate: string;
  kycVerified: boolean;
  distance?: number; // km
  estimatedTime?: string;
}

interface AgentLocatorProps {
  userPostalCode?: string;
  serviceType?: 'CASH_IN' | 'CASH_OUT';
  onSelectAgent?: (agent: Agent) => void;
}

export const AgentLocator: React.FC<AgentLocatorProps> = ({
  userPostalCode = '',
  serviceType,
  onSelectAgent
}) => {
  const theme = useTheme();
  const { formatCurrency, formatNumber, currentLanguage } = useLanguage();
  const [postalCode, setPostalCode] = useState(userPostalCode);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [filteredAgents, setFilteredAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);
  const [showAgentDialog, setShowAgentDialog] = useState(false);
  const [searchRadius, setSearchRadius] = useState(5); // km
  const [selectedService, setSelectedService] = useState<string>(serviceType || 'ALL');
  const [loading, setLoading] = useState(false);

  // Mock agent data - in production would fetch from blockchain
  useEffect(() => {
    const mockAgents: Agent[] = [
      {
        agentId: 'AGT001',
        businessName: 'Sharma Digital Services',
        address: '123, Main Market, Connaught Place',
        postalCode: '110001',
        district: 'Central Delhi',
        phone: '+91-9876543210',
        languages: ['hi', 'en', 'pa'],
        services: ['CASH_IN', 'CASH_OUT', 'BILL_PAYMENT'],
        operatingHours: [
          { day: 'Monday', open: '09:00', close: '19:00', isClosed: false },
          { day: 'Tuesday', open: '09:00', close: '19:00', isClosed: false },
          { day: 'Wednesday', open: '09:00', close: '19:00', isClosed: false },
          { day: 'Thursday', open: '09:00', close: '19:00', isClosed: false },
          { day: 'Friday', open: '09:00', close: '19:00', isClosed: false },
          { day: 'Saturday', open: '10:00', close: '17:00', isClosed: false },
          { day: 'Sunday', open: '', close: '', isClosed: true }
        ],
        rating: 4.7,
        totalRatings: 312,
        totalTransactions: 1567,
        dailyLimit: { amount: '500000000000', denom: 'namo' },
        currentUtilization: 35,
        commissionRate: '2.0',
        kycVerified: true,
        distance: 0.5,
        estimatedTime: '5 min walk'
      },
      {
        agentId: 'AGT002',
        businessName: 'QuickPay Kendra',
        address: '45, Karol Bagh Market',
        postalCode: '110005',
        district: 'Central Delhi',
        phone: '+91-9988776655',
        languages: ['hi', 'en'],
        services: ['CASH_IN', 'CASH_OUT', 'REMITTANCE'],
        operatingHours: [
          { day: 'Monday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Tuesday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Wednesday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Thursday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Friday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Saturday', open: '08:00', close: '20:00', isClosed: false },
          { day: 'Sunday', open: '10:00', close: '14:00', isClosed: false }
        ],
        rating: 4.5,
        totalRatings: 189,
        totalTransactions: 892,
        dailyLimit: { amount: '300000000000', denom: 'namo' },
        currentUtilization: 67,
        commissionRate: '2.5',
        kycVerified: true,
        distance: 2.3,
        estimatedTime: '10 min car'
      },
      {
        agentId: 'AGT003',
        businessName: 'Digital India Point',
        address: '78, Lajpat Nagar',
        postalCode: '110024',
        district: 'South Delhi',
        phone: '+91-9111222333',
        languages: ['hi', 'en', 'bn'],
        services: ['CASH_IN', 'BILL_PAYMENT'],
        operatingHours: [
          { day: 'Monday', open: '10:00', close: '18:00', isClosed: false },
          { day: 'Tuesday', open: '10:00', close: '18:00', isClosed: false },
          { day: 'Wednesday', open: '10:00', close: '18:00', isClosed: false },
          { day: 'Thursday', open: '10:00', close: '18:00', isClosed: false },
          { day: 'Friday', open: '10:00', close: '18:00', isClosed: false },
          { day: 'Saturday', open: '10:00', close: '15:00', isClosed: false },
          { day: 'Sunday', open: '', close: '', isClosed: true }
        ],
        rating: 4.2,
        totalRatings: 98,
        totalTransactions: 445,
        dailyLimit: { amount: '200000000000', denom: 'namo' },
        currentUtilization: 82,
        commissionRate: '3.0',
        kycVerified: true,
        distance: 5.7,
        estimatedTime: '20 min car'
      }
    ];

    setAgents(mockAgents);
  }, []);

  // Filter agents based on criteria
  useEffect(() => {
    let filtered = [...agents];

    // Filter by postal code proximity
    if (postalCode) {
      filtered = filtered.filter(agent => {
        // Simple proximity check - in production would use actual distance
        const userPrefix = postalCode.substring(0, 3);
        const agentPrefix = agent.postalCode.substring(0, 3);
        return userPrefix === agentPrefix || (agent.distance && agent.distance <= searchRadius);
      });
    }

    // Filter by service type
    if (selectedService !== 'ALL') {
      filtered = filtered.filter(agent =>
        agent.services.includes(selectedService as any)
      );
    }

    // Sort by distance and rating
    filtered.sort((a, b) => {
      if (a.distance && b.distance) {
        return a.distance - b.distance;
      }
      return b.rating - a.rating;
    });

    setFilteredAgents(filtered);
  }, [agents, postalCode, searchRadius, selectedService]);

  const handleSelectAgent = (agent: Agent) => {
    setSelectedAgent(agent);
    setShowAgentDialog(true);
  };

  const handleConfirmAgent = () => {
    if (selectedAgent && onSelectAgent) {
      onSelectAgent(selectedAgent);
      setShowAgentDialog(false);
    }
  };

  const getServiceIcon = (service: string) => {
    switch (service) {
      case 'CASH_IN': return <CashInIcon />;
      case 'CASH_OUT': return <CashOutIcon />;
      default: return <StoreIcon />;
    }
  };

  const getServiceLabel = (service: string) => {
    switch (service) {
      case 'CASH_IN': return 'Cash to NAMO';
      case 'CASH_OUT': return 'NAMO to Cash';
      case 'REMITTANCE': return 'Send Money';
      case 'BILL_PAYMENT': return 'Bill Pay';
      default: return service;
    }
  };

  const isAgentOpen = (agent: Agent) => {
    const now = new Date();
    const day = now.toLocaleDateString('en-US', { weekday: 'long' });
    const currentTime = now.toTimeString().substring(0, 5);
    
    const todayHours = agent.operatingHours.find(h => h.day === day);
    if (!todayHours || todayHours.isClosed) return false;
    
    return currentTime >= todayHours.open && currentTime <= todayHours.close;
  };

  return (
    <Box>
      {/* Search Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h5" gutterBottom>
            Find DeshChain Agents
          </Typography>
          
          <Grid container spacing={2} alignItems="flex-end">
            <Grid item xs={12} md={6}>
              <PostalCodeInput
                value={postalCode}
                onChange={(value) => setPostalCode(value)}
                label="Enter Postal Code"
                autoDetect
              />
            </Grid>
            
            <Grid item xs={12} md={3}>
              <TextField
                select
                fullWidth
                label="Service Type"
                value={selectedService}
                onChange={(e) => setSelectedService(e.target.value)}
                SelectProps={{
                  native: true
                }}
              >
                <option value="ALL">All Services</option>
                <option value="CASH_IN">Cash to NAMO</option>
                <option value="CASH_OUT">NAMO to Cash</option>
                <option value="REMITTANCE">Send Money</option>
                <option value="BILL_PAYMENT">Bill Payment</option>
              </TextField>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <TextField
                type="number"
                fullWidth
                label="Search Radius (km)"
                value={searchRadius}
                onChange={(e) => setSearchRadius(parseInt(e.target.value))}
                InputProps={{
                  endAdornment: <InputAdornment position="end">km</InputAdornment>
                }}
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Results Section */}
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h6">
              {filteredAgents.length} Agents Found
            </Typography>
            <Button
              startIcon={<MapIcon />}
              variant="outlined"
              size="small"
            >
              View Map
            </Button>
          </Box>

          <List>
            {filteredAgents.map((agent, index) => (
              <React.Fragment key={agent.agentId}>
                {index > 0 && <Divider />}
                <ListItem
                  component={motion.div}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  sx={{
                    py: 2,
                    cursor: 'pointer',
                    '&:hover': {
                      bgcolor: alpha(theme.palette.primary.main, 0.05)
                    }
                  }}
                  onClick={() => handleSelectAgent(agent)}
                >
                  <ListItemAvatar>
                    <Badge
                      overlap="circular"
                      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                      badgeContent={
                        agent.kycVerified && (
                          <Avatar
                            sx={{
                              width: 22,
                              height: 22,
                              bgcolor: 'success.main'
                            }}
                          >
                            <VerifiedIcon sx={{ fontSize: 14 }} />
                          </Avatar>
                        )
                      }
                    >
                      <Avatar sx={{ bgcolor: 'primary.main', width: 56, height: 56 }}>
                        <StoreIcon />
                      </Avatar>
                    </Badge>
                  </ListItemAvatar>
                  
                  <ListItemText
                    primary={
                      <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="subtitle1" fontWeight="bold">
                          {agent.businessName}
                        </Typography>
                        <Chip
                          label={isAgentOpen(agent) ? 'Open' : 'Closed'}
                          size="small"
                          color={isAgentOpen(agent) ? 'success' : 'default'}
                        />
                      </Box>
                    }
                    secondary={
                      <Box>
                        <Box display="flex" alignItems="center" gap={1} mt={0.5}>
                          <LocationIcon fontSize="small" color="action" />
                          <Typography variant="body2">
                            {agent.address}
                          </Typography>
                        </Box>
                        
                        <Box display="flex" alignItems="center" gap={2} mt={1}>
                          <Box display="flex" alignItems="center" gap={0.5}>
                            <Rating value={agent.rating} readOnly size="small" />
                            <Typography variant="caption">
                              {agent.rating} ({agent.totalRatings})
                            </Typography>
                          </Box>
                          
                          {agent.distance && (
                            <Chip
                              label={`${agent.distance} km • ${agent.estimatedTime}`}
                              size="small"
                              icon={agent.distance < 2 ? <WalkIcon /> : <CarIcon />}
                            />
                          )}
                        </Box>
                        
                        <Box display="flex" gap={0.5} mt={1} flexWrap="wrap">
                          {agent.services.map(service => (
                            <Chip
                              key={service}
                              label={getServiceLabel(service)}
                              size="small"
                              variant="outlined"
                              icon={getServiceIcon(service)}
                            />
                          ))}
                        </Box>
                      </Box>
                    }
                  />
                  
                  <ListItemSecondaryAction>
                    <Box textAlign="right">
                      <Typography variant="body2" color="text.secondary">
                        Commission
                      </Typography>
                      <Typography variant="h6" color="primary">
                        {agent.commissionRate}%
                      </Typography>
                      
                      <Box mt={1}>
                        <Typography variant="caption" color="text.secondary">
                          Daily Limit
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={agent.currentUtilization}
                          sx={{ mt: 0.5, height: 6, borderRadius: 3 }}
                          color={agent.currentUtilization > 80 ? 'error' : 'primary'}
                        />
                        <Typography variant="caption">
                          {agent.currentUtilization}% used
                        </Typography>
                      </Box>
                    </Box>
                  </ListItemSecondaryAction>
                </ListItem>
              </React.Fragment>
            ))}
          </List>

          {filteredAgents.length === 0 && (
            <Box py={4} textAlign="center">
              <Typography color="text.secondary">
                No agents found in your area
              </Typography>
              <Button
                variant="outlined"
                sx={{ mt: 2 }}
                onClick={() => setSearchRadius(searchRadius + 5)}
              >
                Expand Search Radius
              </Button>
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Agent Details Dialog */}
      <Dialog
        open={showAgentDialog}
        onClose={() => setShowAgentDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        {selectedAgent && (
          <>
            <DialogTitle>
              <Box display="flex" alignItems="center" gap={2}>
                <Avatar sx={{ bgcolor: 'primary.main' }}>
                  <StoreIcon />
                </Avatar>
                <Box>
                  <Typography variant="h6">
                    {selectedAgent.businessName}
                  </Typography>
                  <Rating value={selectedAgent.rating} readOnly size="small" />
                </Box>
              </Box>
            </DialogTitle>
            
            <DialogContent>
              <Box display="flex" flexDirection="column" gap={2}>
                <Alert severity="info">
                  Commission Rate: {selectedAgent.commissionRate}% per transaction
                </Alert>
                
                <Box>
                  <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <LocationIcon color="action" />
                    <Typography variant="body2">
                      {selectedAgent.address}, {selectedAgent.district}
                    </Typography>
                  </Box>
                  
                  <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <PhoneIcon color="action" />
                    <Typography variant="body2">
                      {selectedAgent.phone}
                    </Typography>
                  </Box>
                  
                  <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <LanguageIcon color="action" />
                    <Box display="flex" gap={0.5}>
                      {selectedAgent.languages.map(lang => (
                        <Chip key={lang} label={lang.toUpperCase()} size="small" />
                      ))}
                    </Box>
                  </Box>
                </Box>
                
                <Divider />
                
                <Box>
                  <Typography variant="subtitle2" gutterBottom>
                    Operating Hours
                  </Typography>
                  <Grid container spacing={1}>
                    {selectedAgent.operatingHours.map(hours => (
                      <Grid item xs={6} key={hours.day}>
                        <Box display="flex" justifyContent="space-between">
                          <Typography variant="body2">
                            {hours.day}:
                          </Typography>
                          <Typography variant="body2" color={hours.isClosed ? 'text.secondary' : 'text.primary'}>
                            {hours.isClosed ? 'Closed' : `${hours.open} - ${hours.close}`}
                          </Typography>
                        </Box>
                      </Grid>
                    ))}
                  </Grid>
                </Box>
                
                <Divider />
                
                <Box>
                  <Typography variant="subtitle2" gutterBottom>
                    Performance
                  </Typography>
                  <Box display="flex" justifyContent="space-around" textAlign="center">
                    <Box>
                      <Typography variant="h6" color="primary">
                        {formatNumber(selectedAgent.totalTransactions)}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Total Transactions
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="h6" color="primary">
                        {selectedAgent.rating}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Average Rating
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="h6" color="primary">
                        {formatCurrency(parseFloat(selectedAgent.dailyLimit.amount) / 1000000)}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Daily Limit
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              </Box>
            </DialogContent>
            
            <DialogActions>
              <Button onClick={() => setShowAgentDialog(false)}>
                Cancel
              </Button>
              <Button
                variant="contained"
                onClick={handleConfirmAgent}
                startIcon={<LocationIcon />}
              >
                Select This Agent
              </Button>
            </DialogActions>
          </>
        )}
      </Dialog>
    </Box>
  );
};