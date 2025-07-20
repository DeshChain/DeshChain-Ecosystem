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
  Chip,
  Timeline,
  TimelineItem,
  TimelineSeparator,
  TimelineConnector,
  TimelineContent,
  TimelineDot,
  TimelineOppositeContent,
  LinearProgress,
  Avatar,
  Divider,
  Tooltip,
  Button,
  Collapse,
  useTheme,
  alpha
} from '@mui/material';
import {
  LocationOn as LocationIcon,
  FlightTakeoff as FlightIcon,
  Train as TrainIcon,
  DirectionsCar as CarIcon,
  LocalShipping as TruckIcon,
  Schedule as TimeIcon,
  CheckCircle as CheckIcon,
  RadioButtonUnchecked as PendingIcon,
  Speed as SpeedIcon,
  Route as RouteIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { PostalCodeService, PostalRoute } from '../services/postalCodeService';
import { useLanguage } from '../hooks/useLanguage';

interface DeliveryEstimateProps {
  fromPincode: string;
  toPincode: string;
  priority?: 'express' | 'standard' | 'economy';
  amount?: number;
  showDetails?: boolean;
  onRouteCalculated?: (route: PostalRoute) => void;
}

export const DeliveryEstimate: React.FC<DeliveryEstimateProps> = ({
  fromPincode,
  toPincode,
  priority = 'standard',
  amount = 0,
  showDetails = true,
  onRouteCalculated
}) => {
  const theme = useTheme();
  const { formatNumber, formatCurrency } = useLanguage();
  const [loading, setLoading] = useState(false);
  const [route, setRoute] = useState<PostalRoute | null>(null);
  const [estimate, setEstimate] = useState<{ minDays: number; maxDays: number; confidence: number } | null>(null);
  const [showFullDetails, setShowFullDetails] = useState(false);

  // Calculate route when inputs change
  useEffect(() => {
    if (fromPincode && toPincode && 
        PostalCodeService.isValidPincode(fromPincode) && 
        PostalCodeService.isValidPincode(toPincode)) {
      
      setLoading(true);
      PostalCodeService.calculateRoute(fromPincode, toPincode, priority)
        .then(calculatedRoute => {
          if (calculatedRoute) {
            setRoute(calculatedRoute);
            setEstimate(PostalCodeService.getDeliveryEstimate(calculatedRoute));
            onRouteCalculated?.(calculatedRoute);
          }
          setLoading(false);
        })
        .catch(() => setLoading(false));
    }
  }, [fromPincode, toPincode, priority, onRouteCalculated]);

  // Get transport icon
  const getTransportIcon = (mode: string) => {
    switch (mode) {
      case 'air': return <FlightIcon />;
      case 'rail': return <TrainIcon />;
      case 'road': return <CarIcon />;
      case 'combined': return <TruckIcon />;
      default: return <TruckIcon />;
    }
  };

  // Get priority color
  const getPriorityColor = () => {
    switch (priority) {
      case 'express': return 'error';
      case 'economy': return 'success';
      default: return 'primary';
    }
  };

  // Calculate delivery fee
  const calculateDeliveryFee = () => {
    if (!route || amount === 0) return 0;
    
    const baseFee = amount * 0.01; // 1% base
    const distanceFee = route.distance * 0.01; // ₹0.01 per km
    
    let multiplier = 1;
    if (priority === 'express') multiplier = 2;
    if (priority === 'economy') multiplier = 0.7;
    
    return (baseFee + distanceFee) * multiplier;
  };

  if (!route || !estimate) {
    return loading ? (
      <Card>
        <CardContent>
          <Box display="flex" alignItems="center" gap={2}>
            <CircularProgress size={24} />
            <Typography>Calculating delivery route...</Typography>
          </Box>
        </CardContent>
      </Card>
    ) : null;
  }

  const deliveryFee = calculateDeliveryFee();

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
    >
      <Card>
        <CardContent>
          {/* Header */}
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={1}>
              <Avatar sx={{ bgcolor: getPriorityColor() + '.main', width: 40, height: 40 }}>
                {getTransportIcon(route.transportMode)}
              </Avatar>
              <Box>
                <Typography variant="h6">
                  Delivery Estimate
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {route.from.districtName} → {route.to.districtName}
                </Typography>
              </Box>
            </Box>
            
            <Chip
              label={priority.toUpperCase()}
              color={getPriorityColor()}
              size="small"
              icon={<SpeedIcon />}
            />
          </Box>

          {/* Main Estimate */}
          <Box
            sx={{
              p: 2,
              borderRadius: 2,
              bgcolor: alpha(theme.palette.primary.main, 0.05),
              border: `1px solid ${alpha(theme.palette.primary.main, 0.2)}`
            }}
          >
            <Box display="flex" justifyContent="space-between" alignItems="center">
              <Box>
                <Typography variant="h4" color="primary" fontWeight="bold">
                  {estimate.minDays === estimate.maxDays 
                    ? `${formatNumber(estimate.minDays)} Day${estimate.minDays > 1 ? 's' : ''}`
                    : `${formatNumber(estimate.minDays)}-${formatNumber(estimate.maxDays)} Days`
                  }
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Estimated delivery time
                </Typography>
              </Box>
              
              <Box textAlign="right">
                <Typography variant="h6">
                  {formatNumber(route.distance)} km
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Total distance
                </Typography>
              </Box>
            </Box>

            {/* Confidence Indicator */}
            <Box mt={2}>
              <Box display="flex" justifyContent="space-between" mb={0.5}>
                <Typography variant="caption">Delivery confidence</Typography>
                <Typography variant="caption" fontWeight="bold">
                  {(estimate.confidence * 100).toFixed(0)}%
                </Typography>
              </Box>
              <LinearProgress
                variant="determinate"
                value={estimate.confidence * 100}
                sx={{
                  height: 6,
                  borderRadius: 3,
                  bgcolor: 'grey.200',
                  '& .MuiLinearProgress-bar': {
                    borderRadius: 3,
                    bgcolor: estimate.confidence > 0.8 ? 'success.main' : 'warning.main'
                  }
                }}
              />
            </Box>
          </Box>

          {/* Delivery Fee */}
          {amount > 0 && (
            <Box mt={2} display="flex" justifyContent="space-between" alignItems="center">
              <Typography variant="body2" color="text.secondary">
                Delivery Fee
              </Typography>
              <Typography variant="h6" color={getPriorityColor() + '.main'}>
                {formatCurrency(deliveryFee)}
              </Typography>
            </Box>
          )}

          {/* Show Details Button */}
          {showDetails && (
            <>
              <Divider sx={{ my: 2 }} />
              <Button
                fullWidth
                variant="text"
                onClick={() => setShowFullDetails(!showFullDetails)}
                endIcon={<RouteIcon />}
              >
                {showFullDetails ? 'Hide' : 'Show'} Route Details
              </Button>
            </>
          )}
        </CardContent>

        {/* Detailed Route Timeline */}
        <Collapse in={showFullDetails}>
          <Divider />
          <CardContent>
            <Typography variant="subtitle2" gutterBottom>
              Delivery Route
            </Typography>
            
            <Timeline position="alternate" sx={{ mt: 2 }}>
              {/* Origin */}
              <TimelineItem>
                <TimelineOppositeContent color="text.secondary">
                  <Typography variant="caption">Day 0</Typography>
                  <Typography variant="body2">Pickup</Typography>
                </TimelineOppositeContent>
                <TimelineSeparator>
                  <TimelineDot color="primary">
                    <LocationIcon />
                  </TimelineDot>
                  <TimelineConnector />
                </TimelineSeparator>
                <TimelineContent>
                  <Typography variant="body2" fontWeight="bold">
                    {route.from.officeName}
                  </Typography>
                  <Typography variant="caption">
                    {route.from.districtName}, {route.from.stateName}
                  </Typography>
                </TimelineContent>
              </TimelineItem>

              {/* Hubs */}
              {route.hubs?.map((hub, index) => (
                <TimelineItem key={index}>
                  <TimelineOppositeContent color="text.secondary">
                    <Typography variant="caption">
                      Day {Math.ceil((index + 1) * route.estimatedTime / (route.hubs!.length + 2) / 24)}
                    </Typography>
                    <Typography variant="body2">Transit Hub</Typography>
                  </TimelineOppositeContent>
                  <TimelineSeparator>
                    <TimelineDot color="secondary">
                      {getTransportIcon(route.transportMode)}
                    </TimelineDot>
                    <TimelineConnector />
                  </TimelineSeparator>
                  <TimelineContent>
                    <Typography variant="body2" fontWeight="bold">
                      {hub.officeName}
                    </Typography>
                    <Typography variant="caption">
                      {hub.districtName}, {hub.stateName}
                    </Typography>
                  </TimelineContent>
                </TimelineItem>
              ))}

              {/* Destination */}
              <TimelineItem>
                <TimelineOppositeContent color="text.secondary">
                  <Typography variant="caption">
                    Day {estimate.minDays}-{estimate.maxDays}
                  </Typography>
                  <Typography variant="body2">Delivery</Typography>
                </TimelineOppositeContent>
                <TimelineSeparator>
                  <TimelineDot color="success">
                    <CheckIcon />
                  </TimelineDot>
                </TimelineSeparator>
                <TimelineContent>
                  <Typography variant="body2" fontWeight="bold">
                    {route.to.officeName}
                  </Typography>
                  <Typography variant="caption">
                    {route.to.districtName}, {route.to.stateName}
                  </Typography>
                </TimelineContent>
              </TimelineItem>
            </Timeline>

            {/* Route Info */}
            <Box mt={3} display="flex" gap={2} flexWrap="wrap">
              <Chip
                icon={getTransportIcon(route.transportMode)}
                label={`${route.transportMode.charAt(0).toUpperCase() + route.transportMode.slice(1)} Transport`}
                variant="outlined"
                size="small"
              />
              <Chip
                icon={<RouteIcon />}
                label={`${route.routeType.replace('-', ' ').toUpperCase()} Route`}
                variant="outlined"
                size="small"
              />
              <Chip
                icon={<TimeIcon />}
                label={`~${route.estimatedTime} hours`}
                variant="outlined"
                size="small"
              />
            </Box>
          </CardContent>
        </Collapse>
      </Card>
    </motion.div>
  );
};

// Mini delivery estimate for inline display
export const DeliveryEstimateMini: React.FC<{
  fromPincode: string;
  toPincode: string;
  priority?: 'express' | 'standard' | 'economy';
}> = ({ fromPincode, toPincode, priority = 'standard' }) => {
  const [estimate, setEstimate] = useState<{ minDays: number; maxDays: number } | null>(null);
  const { formatNumber } = useLanguage();

  useEffect(() => {
    if (fromPincode && toPincode) {
      PostalCodeService.calculateRoute(fromPincode, toPincode, priority)
        .then(route => {
          if (route) {
            const est = PostalCodeService.getDeliveryEstimate(route);
            setEstimate(est);
          }
        });
    }
  }, [fromPincode, toPincode, priority]);

  if (!estimate) return null;

  return (
    <Chip
      icon={<TimeIcon />}
      label={`${formatNumber(estimate.minDays)}-${formatNumber(estimate.maxDays)} days`}
      size="small"
      color={priority === 'express' ? 'error' : priority === 'economy' ? 'success' : 'primary'}
    />
  );
};