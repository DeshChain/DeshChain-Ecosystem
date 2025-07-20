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
  Paper,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  TextField,
  InputAdornment,
  Slider,
  Chip,
  Avatar,
  AvatarGroup,
  LinearProgress,
  ToggleButton,
  ToggleButtonGroup,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  IconButton,
  Tooltip,
  Alert
} from '@mui/material';
import {
  Add as AddIcon,
  Remove as RemoveIcon,
  Refresh as RefreshIcon,
  Info as InfoIcon,
  TrendingUp as TrendingUpIcon,
  AccountBalance as PoolIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useMoneyOrder, usePoolData } from '@deshchain/money-order-ui';

interface LiquidityManagerProps {
  onPoolSelect?: (poolId: string) => void;
}

interface Position {
  id: string;
  poolId: string;
  poolName: string;
  shares: string;
  value: string;
  rewards: string;
  apy: number;
}

export const LiquidityManager: React.FC<LiquidityManagerProps> = ({ onPoolSelect }) => {
  const [activeView, setActiveView] = useState<'add' | 'positions'>('positions');
  const [selectedPool, setSelectedPool] = useState<string>('');
  const [tokenAAmount, setTokenAAmount] = useState('');
  const [tokenBAmount, setTokenBAmount] = useState('');
  const [sliderValue, setSliderValue] = useState(50);

  const { addLiquidity, removeLiquidity, getLiquidityPositions } = useMoneyOrder();
  const { pools, getPoolDetails } = usePoolData();

  // Mock positions data
  const positions: Position[] = [
    {
      id: '1',
      poolId: '1',
      poolName: 'NAMO/INR Fixed',
      shares: '1000000',
      value: '₹75,000',
      rewards: '₹1,250',
      apy: 18.5
    },
    {
      id: '2',
      poolId: '2',
      poolName: 'Village Pool Delhi',
      shares: '500000',
      value: '₹37,500',
      rewards: '₹625',
      apy: 22.3
    }
  ];

  const handleAddLiquidity = async () => {
    try {
      await addLiquidity({
        poolId: selectedPool,
        tokenA: {
          denom: 'unamo',
          amount: (parseFloat(tokenAAmount) * 1000000).toString()
        },
        tokenB: {
          denom: 'inr',
          amount: (parseFloat(tokenBAmount) * 1000000).toString()
        }
      });
    } catch (error) {
      console.error('Failed to add liquidity:', error);
    }
  };

  const handleRemoveLiquidity = async (positionId: string, percentage: number) => {
    try {
      const position = positions.find(p => p.id === positionId);
      if (!position) return;

      const sharesToRemove = (parseFloat(position.shares) * percentage / 100).toString();
      await removeLiquidity({
        positionId,
        sharesToRemove
      });
    } catch (error) {
      console.error('Failed to remove liquidity:', error);
    }
  };

  const renderAddLiquidity = () => (
    <Grid container spacing={3}>
      {/* Pool Selection */}
      <Grid item xs={12}>
        <Typography variant="h6" gutterBottom>
          Select Pool
        </Typography>
        <Grid container spacing={2}>
          {[
            { id: '1', name: 'NAMO/INR Fixed', apy: 18.5, tvl: '₹10M', type: 'fixed' },
            { id: '2', name: 'NAMO/INR AMM', apy: 24.3, tvl: '₹5M', type: 'amm' },
            { id: '3', name: 'Village Pool Delhi', apy: 22.1, tvl: '₹2M', type: 'village' }
          ].map(pool => (
            <Grid item xs={12} sm={6} md={4} key={pool.id}>
              <Card
                sx={{
                  cursor: 'pointer',
                  border: selectedPool === pool.id ? 2 : 1,
                  borderColor: selectedPool === pool.id ? 'primary.main' : 'divider',
                  '&:hover': { borderColor: 'primary.main' }
                }}
                onClick={() => {
                  setSelectedPool(pool.id);
                  onPoolSelect?.(pool.id);
                }}
              >
                <CardContent>
                  <Box display="flex" justifyContent="space-between" alignItems="start" mb={2}>
                    <Typography variant="h6">{pool.name}</Typography>
                    <Chip 
                      label={pool.type.toUpperCase()} 
                      size="small" 
                      color={pool.type === 'fixed' ? 'primary' : pool.type === 'village' ? 'secondary' : 'default'}
                    />
                  </Box>
                  <Box display="flex" justifyContent="space-between">
                    <Box>
                      <Typography variant="body2" color="text.secondary">TVL</Typography>
                      <Typography variant="h6">{pool.tvl}</Typography>
                    </Box>
                    <Box textAlign="right">
                      <Typography variant="body2" color="text.secondary">APY</Typography>
                      <Typography variant="h6" color="success.main">{pool.apy}%</Typography>
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Grid>

      {/* Add Liquidity Form */}
      {selectedPool && (
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Add Liquidity
            </Typography>

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="NAMO Amount"
                  type="number"
                  value={tokenAAmount}
                  onChange={(e) => setTokenAAmount(e.target.value)}
                  InputProps={{
                    endAdornment: <InputAdornment position="end">NAMO</InputAdornment>
                  }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="INR Amount"
                  type="number"
                  value={tokenBAmount}
                  onChange={(e) => setTokenBAmount(e.target.value)}
                  InputProps={{
                    startAdornment: <InputAdornment position="start">₹</InputAdornment>
                  }}
                />
              </Grid>
            </Grid>

            <Box mt={3}>
              <Typography variant="body2" gutterBottom>
                Pool Share: ~2.5%
              </Typography>
              <LinearProgress variant="determinate" value={2.5} sx={{ height: 8, borderRadius: 4 }} />
            </Box>

            <Box mt={3} p={2} bgcolor="background.default" borderRadius={1}>
              <Typography variant="subtitle2" gutterBottom>
                Expected Returns
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Daily</Typography>
                  <Typography variant="body1">₹{(75000 * 0.185 / 365).toFixed(2)}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Monthly</Typography>
                  <Typography variant="body1">₹{(75000 * 0.185 / 12).toFixed(2)}</Typography>
                </Grid>
              </Grid>
            </Box>

            <Button
              variant="contained"
              fullWidth
              size="large"
              startIcon={<AddIcon />}
              onClick={handleAddLiquidity}
              sx={{ mt: 3 }}
            >
              Add Liquidity
            </Button>
          </Paper>
        </Grid>
      )}
    </Grid>
  );

  const renderPositions = () => (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h6">Your Positions</Typography>
        <Button startIcon={<RefreshIcon />} size="small">
          Refresh
        </Button>
      </Box>

      <Grid container spacing={3}>
        {/* Summary Cards */}
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary">Total Value</Typography>
              <Typography variant="h5">₹112,500</Typography>
              <Chip label="+12.5%" color="success" size="small" sx={{ mt: 1 }} />
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary">Total Rewards</Typography>
              <Typography variant="h5">₹1,875</Typography>
              <Button size="small" sx={{ mt: 1 }}>Claim All</Button>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary">Avg APY</Typography>
              <Typography variant="h5" color="success.main">20.4%</Typography>
              <Typography variant="caption" sx={{ mt: 1 }}>Across all positions</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary">Active Pools</Typography>
              <AvatarGroup max={4} sx={{ mt: 1 }}>
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>F</Avatar>
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'secondary.main' }}>V</Avatar>
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'info.main' }}>A</Avatar>
              </AvatarGroup>
            </CardContent>
          </Card>
        </Grid>

        {/* Positions Table */}
        <Grid item xs={12}>
          <Paper>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Pool</TableCell>
                  <TableCell align="right">Value</TableCell>
                  <TableCell align="right">Rewards</TableCell>
                  <TableCell align="right">APY</TableCell>
                  <TableCell align="center">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {positions.map(position => (
                  <TableRow key={position.id}>
                    <TableCell>
                      <Box display="flex" alignItems="center" gap={1}>
                        <PoolIcon color="primary" />
                        <Box>
                          <Typography variant="body2">{position.poolName}</Typography>
                          <Typography variant="caption" color="text.secondary">
                            {position.shares} shares
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell align="right">
                      <Typography variant="body2" fontWeight="medium">{position.value}</Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Box>
                        <Typography variant="body2" color="success.main">{position.rewards}</Typography>
                        <Button size="small" sx={{ mt: 0.5 }}>Claim</Button>
                      </Box>
                    </TableCell>
                    <TableCell align="right">
                      <Box display="flex" alignItems="center" justifyContent="flex-end" gap={0.5}>
                        <TrendingUpIcon sx={{ fontSize: 16, color: 'success.main' }} />
                        <Typography variant="body2" color="success.main">{position.apy}%</Typography>
                      </Box>
                    </TableCell>
                    <TableCell align="center">
                      <Tooltip title="Add more liquidity">
                        <IconButton size="small">
                          <AddIcon />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Remove liquidity">
                        <IconButton size="small">
                          <RemoveIcon />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="View details">
                        <IconButton size="small">
                          <InfoIcon />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );

  return (
    <Box sx={{ p: 3 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h5">Liquidity Manager</Typography>
        <ToggleButtonGroup
          value={activeView}
          exclusive
          onChange={(_, value) => value && setActiveView(value)}
        >
          <ToggleButton value="positions">My Positions</ToggleButton>
          <ToggleButton value="add">Add Liquidity</ToggleButton>
        </ToggleButtonGroup>
      </Box>

      {activeView === 'positions' ? renderPositions() : renderAddLiquidity()}

      {/* Info Alert */}
      <Alert severity="info" sx={{ mt: 3 }}>
        <Typography variant="body2">
          Providing liquidity earns you trading fees and cultural rewards. Your share of the pool determines your portion of fees earned.
        </Typography>
      </Alert>
    </Box>
  );
};