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
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab,
  LinearProgress,
  useTheme
} from '@mui/material';
import { FixedSizeList as List } from 'react-window';
import useWebSocket from 'react-use-websocket';

interface OrderBookProps {
  poolId: string;
  pair: string;
}

interface Order {
  price: number;
  amount: number;
  total: number;
  percentage: number;
}

interface OrderBookData {
  bids: Order[];
  asks: Order[];
  spread: number;
  spreadPercentage: number;
}

export const OrderBook: React.FC<OrderBookProps> = ({ poolId, pair }) => {
  const theme = useTheme();
  const [activeTab, setActiveTab] = useState(0);
  const [orderBook, setOrderBook] = useState<OrderBookData>({
    bids: [],
    asks: [],
    spread: 0,
    spreadPercentage: 0
  });

  // WebSocket for real-time order book updates
  const { lastMessage } = useWebSocket(
    `wss://api.deshchain.org/v1/moneyorder/ws/pools/${poolId}/orderbook`,
    {
      shouldReconnect: () => true,
      reconnectInterval: 3000
    }
  );

  useEffect(() => {
    // Generate mock order book data
    generateMockOrderBook();
  }, [poolId, pair]);

  useEffect(() => {
    if (lastMessage) {
      try {
        const data = JSON.parse(lastMessage.data);
        if (data.bids && data.asks) {
          processOrderBookData(data);
        }
      } catch (error) {
        console.error('Failed to parse order book data:', error);
      }
    }
  }, [lastMessage]);

  const generateMockOrderBook = () => {
    const midPrice = 0.075;
    const bids: Order[] = [];
    const asks: Order[] = [];

    // Generate bids
    for (let i = 0; i < 20; i++) {
      const price = midPrice - (i + 1) * 0.0001;
      const amount = Math.random() * 100000 + 10000;
      bids.push({
        price,
        amount,
        total: price * amount,
        percentage: 0
      });
    }

    // Generate asks
    for (let i = 0; i < 20; i++) {
      const price = midPrice + (i + 1) * 0.0001;
      const amount = Math.random() * 100000 + 10000;
      asks.push({
        price,
        amount,
        total: price * amount,
        percentage: 0
      });
    }

    // Calculate cumulative percentages
    const maxBidTotal = bids.reduce((sum, bid) => sum + bid.total, 0);
    const maxAskTotal = asks.reduce((sum, ask) => sum + ask.total, 0);
    
    let cumBid = 0;
    bids.forEach(bid => {
      cumBid += bid.total;
      bid.percentage = (cumBid / maxBidTotal) * 100;
    });

    let cumAsk = 0;
    asks.forEach(ask => {
      cumAsk += ask.total;
      ask.percentage = (cumAsk / maxAskTotal) * 100;
    });

    const spread = asks[0].price - bids[0].price;
    const spreadPercentage = (spread / midPrice) * 100;

    setOrderBook({ bids, asks, spread, spreadPercentage });
  };

  const processOrderBookData = (data: any) => {
    // Process real WebSocket data
    setOrderBook({
      bids: data.bids.map((bid: any) => ({
        price: bid.price,
        amount: bid.amount,
        total: bid.price * bid.amount,
        percentage: bid.percentage || 0
      })),
      asks: data.asks.map((ask: any) => ({
        price: ask.price,
        amount: ask.amount,
        total: ask.price * ask.amount,
        percentage: ask.percentage || 0
      })),
      spread: data.spread || 0,
      spreadPercentage: data.spreadPercentage || 0
    });
  };

  const OrderRow = ({ order, type }: { order: Order; type: 'bid' | 'ask' }) => (
    <TableRow
      sx={{
        position: 'relative',
        '&:hover': { bgcolor: 'action.hover' },
        cursor: 'pointer'
      }}
    >
      <TableCell 
        sx={{ 
          color: type === 'bid' ? 'success.main' : 'error.main',
          fontWeight: 'medium',
          fontSize: '0.875rem'
        }}
      >
        {order.price.toFixed(4)}
      </TableCell>
      <TableCell align="right" sx={{ fontSize: '0.875rem' }}>
        {order.amount.toLocaleString()}
      </TableCell>
      <TableCell align="right" sx={{ fontSize: '0.875rem' }}>
        ₹{order.total.toFixed(2)}
      </TableCell>
      
      {/* Background depth visualization */}
      <Box
        sx={{
          position: 'absolute',
          top: 0,
          right: 0,
          bottom: 0,
          width: `${order.percentage}%`,
          bgcolor: type === 'bid' ? 'success.main' : 'error.main',
          opacity: 0.1,
          zIndex: 0
        }}
      />
    </TableRow>
  );

  const renderOrderList = (orders: Order[], type: 'bid' | 'ask') => (
    <TableContainer sx={{ maxHeight: 300 }}>
      <Table size="small" stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Price (INR)</TableCell>
            <TableCell align="right">Amount</TableCell>
            <TableCell align="right">Total</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {orders.map((order, index) => (
            <OrderRow key={index} order={order} type={type} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );

  return (
    <Paper sx={{ p: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Typography variant="h6" gutterBottom>
        Order Book
      </Typography>

      {/* Spread Display */}
      <Box
        sx={{
          p: 1,
          mb: 2,
          bgcolor: 'background.default',
          borderRadius: 1,
          textAlign: 'center'
        }}
      >
        <Typography variant="body2" color="text.secondary">
          Spread: ₹{orderBook.spread.toFixed(4)} ({orderBook.spreadPercentage.toFixed(2)}%)
        </Typography>
      </Box>

      {/* Tabs for different views */}
      <Tabs value={activeTab} onChange={(_, v) => setActiveTab(v)} sx={{ mb: 2 }}>
        <Tab label="All" />
        <Tab label="Bids" />
        <Tab label="Asks" />
      </Tabs>

      {/* Order Book Content */}
      <Box sx={{ flexGrow: 1, overflow: 'hidden' }}>
        {activeTab === 0 && (
          <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Asks (reversed for display) */}
            <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
              {renderOrderList([...orderBook.asks].reverse(), 'ask')}
            </Box>
            
            {/* Current Price Divider */}
            <Box
              sx={{
                p: 1,
                bgcolor: 'primary.main',
                color: 'primary.contrastText',
                textAlign: 'center'
              }}
            >
              <Typography variant="body2" fontWeight="bold">
                ₹0.0750
              </Typography>
            </Box>
            
            {/* Bids */}
            <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
              {renderOrderList(orderBook.bids, 'bid')}
            </Box>
          </Box>
        )}

        {activeTab === 1 && renderOrderList(orderBook.bids, 'bid')}
        {activeTab === 2 && renderOrderList(orderBook.asks, 'ask')}
      </Box>
    </Paper>
  );
};