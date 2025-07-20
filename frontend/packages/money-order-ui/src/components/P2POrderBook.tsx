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
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Button,
  Avatar,
  Tooltip,
  IconButton,
  Badge,
  TextField,
  InputAdornment,
  ToggleButton,
  ToggleButtonGroup,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  LinearProgress,
  useTheme,
  alpha
} from '@mui/material';
import {
  ShoppingCart as BuyIcon,
  Sell as SellIcon,
  LocationOn as LocationIcon,
  Star as StarIcon,
  Search as SearchIcon,
  FilterList as FilterIcon,
  AccountBalance as BankIcon,
  PhoneAndroid as UPIIcon,
  Money as CashIcon,
  Lock as EscrowIcon,
  Timer as TimerIcon,
  Person as PersonIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { PostalCodeService } from '../services/postalCodeService';
import { useLanguage } from '../hooks/useLanguage';

interface P2POrder {
  orderId: string;
  creator: string;
  orderType: 'BUY_NAMO' | 'SELL_NAMO';
  amount: { amount: string; denom: string };
  fiatAmount: { amount: string; currency: string };
  rate: number;
  postalCode: string;
  district: string;
  paymentMethods: PaymentMethod[];
  minAmount: { amount: string; denom: string };
  maxAmount: { amount: string; denom: string };
  trustScore: number;
  completedTrades: number;
  responseTime: string;
  isKYCVerified: boolean;
  languages: string[];
  status: 'ACTIVE' | 'MATCHED' | 'COMPLETED';
  expiresIn: number; // minutes
}

interface PaymentMethod {
  type: 'UPI' | 'IMPS' | 'NEFT' | 'CASH';
  provider?: string;
}

interface P2POrderBookProps {
  userPostalCode: string;
  onSelectOrder: (order: P2POrder) => void;
}

export const P2POrderBook: React.FC<P2POrderBookProps> = ({
  userPostalCode,
  onSelectOrder
}) => {
  const theme = useTheme();
  const { formatCurrency, formatNumber } = useLanguage();
  const [orderType, setOrderType] = useState<'BUY' | 'SELL'>('BUY');
  const [orders, setOrders] = useState<P2POrder[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<P2POrder[]>([]);
  const [searchAmount, setSearchAmount] = useState('');
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<string>('ALL');
  const [selectedOrder, setSelectedOrder] = useState<P2POrder | null>(null);
  const [showTradeDialog, setShowTradeDialog] = useState(false);
  const [loading, setLoading] = useState(false);

  // Mock data - in production would fetch from blockchain
  useEffect(() => {
    const mockOrders: P2POrder[] = [
      {
        orderId: 'P2P001',
        creator: 'desh1abc...xyz',
        orderType: 'SELL_NAMO',
        amount: { amount: '1000000000', denom: 'namo' },
        fiatAmount: { amount: '82500', currency: 'INR' },
        rate: 82.50,
        postalCode: '110001',
        district: 'Central Delhi',
        paymentMethods: [
          { type: 'UPI', provider: 'GPay' },
          { type: 'UPI', provider: 'PhonePe' },
          { type: 'IMPS' }
        ],
        minAmount: { amount: '100000000', denom: 'namo' },
        maxAmount: { amount: '1000000000', denom: 'namo' },
        trustScore: 95,
        completedTrades: 247,
        responseTime: '< 5 min',
        isKYCVerified: true,
        languages: ['en', 'hi'],
        status: 'ACTIVE',
        expiresIn: 120
      },
      {
        orderId: 'P2P002',
        creator: 'desh2def...uvw',
        orderType: 'SELL_NAMO',
        amount: { amount: '500000000', denom: 'namo' },
        fiatAmount: { amount: '41000', currency: 'INR' },
        rate: 82.00,
        postalCode: '110005',
        district: 'Central Delhi',
        paymentMethods: [
          { type: 'CASH' },
          { type: 'UPI', provider: 'PayTM' }
        ],
        minAmount: { amount: '50000000', denom: 'namo' },
        maxAmount: { amount: '500000000', denom: 'namo' },
        trustScore: 88,
        completedTrades: 156,
        responseTime: '< 15 min',
        isKYCVerified: true,
        languages: ['hi', 'pa'],
        status: 'ACTIVE',
        expiresIn: 90
      },
      {
        orderId: 'P2P003',
        creator: 'desh3ghi...rst',
        orderType: 'BUY_NAMO',
        amount: { amount: '2000000000', denom: 'namo' },
        fiatAmount: { amount: '166000', currency: 'INR' },
        rate: 83.00,
        postalCode: '400001',
        district: 'Mumbai',
        paymentMethods: [
          { type: 'NEFT' },
          { type: 'IMPS' }
        ],
        minAmount: { amount: '500000000', denom: 'namo' },
        maxAmount: { amount: '2000000000', denom: 'namo' },
        trustScore: 92,
        completedTrades: 512,
        responseTime: '< 10 min',
        isKYCVerified: true,
        languages: ['en', 'mr'],
        status: 'ACTIVE',
        expiresIn: 180
      }
    ];

    setOrders(mockOrders);
  }, []);

  // Filter orders based on type and user preferences
  useEffect(() => {
    let filtered = orders.filter(order => {
      // Filter by order type (opposite of what user wants)
      if (orderType === 'BUY' && order.orderType !== 'SELL_NAMO') return false;
      if (orderType === 'SELL' && order.orderType !== 'BUY_NAMO') return false;

      // Filter by amount if specified
      if (searchAmount) {
        const amount = parseFloat(searchAmount) * 1000000; // Convert to base units
        const minAmount = parseFloat(order.minAmount.amount);
        const maxAmount = parseFloat(order.maxAmount.amount);
        if (amount < minAmount || amount > maxAmount) return false;
      }

      // Filter by payment method
      if (selectedPaymentMethod !== 'ALL') {
        const hasMethod = order.paymentMethods.some(
          pm => pm.type === selectedPaymentMethod
        );
        if (!hasMethod) return false;
      }

      return true;
    });

    // Sort by distance and trust score
    filtered.sort((a, b) => {
      // Prioritize same postal code
      if (a.postalCode === userPostalCode && b.postalCode !== userPostalCode) return -1;
      if (b.postalCode === userPostalCode && a.postalCode !== userPostalCode) return 1;
      
      // Then by trust score
      return b.trustScore - a.trustScore;
    });

    setFilteredOrders(filtered);
  }, [orders, orderType, searchAmount, selectedPaymentMethod, userPostalCode]);

  const handleSelectOrder = (order: P2POrder) => {
    setSelectedOrder(order);
    setShowTradeDialog(true);
  };

  const handleConfirmTrade = async () => {
    if (!selectedOrder) return;
    
    setLoading(true);
    // In production, would create escrow and initiate trade
    setTimeout(() => {
      setLoading(false);
      setShowTradeDialog(false);
      onSelectOrder(selectedOrder);
    }, 2000);
  };

  const getPaymentIcon = (type: string) => {
    switch (type) {
      case 'UPI': return <UPIIcon fontSize="small" />;
      case 'CASH': return <CashIcon fontSize="small" />;
      default: return <BankIcon fontSize="small" />;
    }
  };

  const getDistanceText = (orderPostalCode: string) => {
    if (orderPostalCode === userPostalCode) return 'Same Area';
    if (orderPostalCode.substring(0, 3) === userPostalCode.substring(0, 3)) return 'Nearby';
    if (orderPostalCode.substring(0, 2) === userPostalCode.substring(0, 2)) return 'Same District';
    return 'Same City';
  };

  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box mb={3}>
          <Typography variant="h5" gutterBottom>
            P2P Trading
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Trade NAMO directly with other users in your area
          </Typography>
        </Box>

        {/* Controls */}
        <Box display="flex" gap={2} mb={3} flexWrap="wrap">
          <ToggleButtonGroup
            value={orderType}
            exclusive
            onChange={(e, value) => value && setOrderType(value)}
            size="small"
          >
            <ToggleButton value="BUY">
              <BuyIcon sx={{ mr: 1 }} />
              Buy NAMO
            </ToggleButton>
            <ToggleButton value="SELL">
              <SellIcon sx={{ mr: 1 }} />
              Sell NAMO
            </ToggleButton>
          </ToggleButtonGroup>

          <TextField
            placeholder="Amount"
            value={searchAmount}
            onChange={(e) => setSearchAmount(e.target.value)}
            size="small"
            sx={{ width: 200 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">₹</InputAdornment>
              )
            }}
          />

          <ToggleButtonGroup
            value={selectedPaymentMethod}
            exclusive
            onChange={(e, value) => value && setSelectedPaymentMethod(value)}
            size="small"
          >
            <ToggleButton value="ALL">All</ToggleButton>
            <ToggleButton value="UPI">UPI</ToggleButton>
            <ToggleButton value="CASH">Cash</ToggleButton>
            <ToggleButton value="IMPS">IMPS</ToggleButton>
          </ToggleButtonGroup>
        </Box>

        {/* Orders Table */}
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Trader</TableCell>
                <TableCell>Rate</TableCell>
                <TableCell>Available/Limit</TableCell>
                <TableCell>Payment</TableCell>
                <TableCell>Location</TableCell>
                <TableCell align="right">Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredOrders.map((order) => (
                <TableRow
                  key={order.orderId}
                  hover
                  sx={{ cursor: 'pointer' }}
                  onClick={() => handleSelectOrder(order)}
                >
                  <TableCell>
                    <Box display="flex" alignItems="center" gap={1}>
                      <Badge
                        overlap="circular"
                        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                        badgeContent={
                          order.isKYCVerified && (
                            <Tooltip title="KYC Verified">
                              <Avatar
                                sx={{
                                  width: 20,
                                  height: 20,
                                  bgcolor: 'success.main'
                                }}
                              >
                                ✓
                              </Avatar>
                            </Tooltip>
                          )
                        }
                      >
                        <Avatar sx={{ bgcolor: 'primary.main' }}>
                          <PersonIcon />
                        </Avatar>
                      </Badge>
                      <Box>
                        <Typography variant="body2">
                          {order.creator.substring(0, 8)}...
                        </Typography>
                        <Box display="flex" alignItems="center" gap={0.5}>
                          <StarIcon sx={{ fontSize: 16, color: 'warning.main' }} />
                          <Typography variant="caption">
                            {order.trustScore}% ({order.completedTrades})
                          </Typography>
                        </Box>
                      </Box>
                    </Box>
                  </TableCell>
                  
                  <TableCell>
                    <Typography variant="body1" fontWeight="bold">
                      ₹{order.rate}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      per NAMO
                    </Typography>
                  </TableCell>
                  
                  <TableCell>
                    <Typography variant="body2">
                      {formatNumber(parseFloat(order.amount.amount) / 1000000)} NAMO
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {formatCurrency(parseFloat(order.minAmount.amount) / 1000000 * order.rate)} - 
                      {formatCurrency(parseFloat(order.maxAmount.amount) / 1000000 * order.rate)}
                    </Typography>
                  </TableCell>
                  
                  <TableCell>
                    <Box display="flex" gap={0.5} flexWrap="wrap">
                      {order.paymentMethods.map((pm, index) => (
                        <Chip
                          key={index}
                          label={pm.provider || pm.type}
                          size="small"
                          icon={getPaymentIcon(pm.type)}
                          variant="outlined"
                        />
                      ))}
                    </Box>
                  </TableCell>
                  
                  <TableCell>
                    <Box display="flex" alignItems="center" gap={0.5}>
                      <LocationIcon fontSize="small" color="action" />
                      <Box>
                        <Typography variant="body2">
                          {getDistanceText(order.postalCode)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {order.district}
                        </Typography>
                      </Box>
                    </Box>
                  </TableCell>
                  
                  <TableCell align="right">
                    <Button
                      variant="contained"
                      size="small"
                      color={orderType === 'BUY' ? 'success' : 'error'}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleSelectOrder(order);
                      }}
                    >
                      {orderType === 'BUY' ? 'Buy' : 'Sell'}
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {filteredOrders.length === 0 && (
          <Box py={4} textAlign="center">
            <Typography color="text.secondary">
              No orders found matching your criteria
            </Typography>
            <Button sx={{ mt: 2 }}>Create Alert</Button>
          </Box>
        )}

        {/* Trade Dialog */}
        <Dialog
          open={showTradeDialog}
          onClose={() => setShowTradeDialog(false)}
          maxWidth="sm"
          fullWidth
        >
          {selectedOrder && (
            <>
              <DialogTitle>
                {orderType === 'BUY' ? 'Buy' : 'Sell'} NAMO
              </DialogTitle>
              <DialogContent>
                <Alert severity="info" sx={{ mb: 2 }}>
                  <Typography variant="body2">
                    Your funds will be held in escrow until the trade is completed.
                    If no match is found within 24 hours, you'll receive a full refund including fees.
                  </Typography>
                </Alert>

                <Box display="flex" flexDirection="column" gap={2}>
                  <Box>
                    <Typography variant="caption" color="text.secondary">
                      Rate
                    </Typography>
                    <Typography variant="h6">
                      ₹{selectedOrder.rate} per NAMO
                    </Typography>
                  </Box>

                  <TextField
                    label="Amount to trade"
                    type="number"
                    fullWidth
                    InputProps={{
                      endAdornment: <InputAdornment position="end">NAMO</InputAdornment>
                    }}
                    helperText={`Min: ${formatNumber(parseFloat(selectedOrder.minAmount.amount) / 1000000)} - Max: ${formatNumber(parseFloat(selectedOrder.maxAmount.amount) / 1000000)}`}
                  />

                  <Box>
                    <Typography variant="caption" color="text.secondary">
                      Payment Method
                    </Typography>
                    <Box display="flex" gap={1} mt={0.5}>
                      {selectedOrder.paymentMethods.map((pm, index) => (
                        <Chip
                          key={index}
                          label={pm.provider || pm.type}
                          icon={getPaymentIcon(pm.type)}
                        />
                      ))}
                    </Box>
                  </Box>

                  <Box display="flex" alignItems="center" gap={1}>
                    <TimerIcon color="action" />
                    <Typography variant="body2" color="text.secondary">
                      Expires in {selectedOrder.expiresIn} minutes
                    </Typography>
                  </Box>

                  <Box display="flex" alignItems="center" gap={1}>
                    <EscrowIcon color="primary" />
                    <Typography variant="body2">
                      Secure escrow protection
                    </Typography>
                  </Box>
                </Box>
              </DialogContent>
              <DialogActions>
                <Button onClick={() => setShowTradeDialog(false)}>
                  Cancel
                </Button>
                <Button
                  variant="contained"
                  onClick={handleConfirmTrade}
                  disabled={loading}
                  startIcon={loading && <CircularProgress size={16} />}
                >
                  {loading ? 'Creating Escrow...' : 'Confirm & Deposit'}
                </Button>
              </DialogActions>
            </>
          )}
        </Dialog>
      </CardContent>
    </Card>
  );
};