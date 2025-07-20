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

import React, { useEffect, useRef, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  ToggleButton,
  ToggleButtonGroup,
  IconButton,
  Menu,
  MenuItem,
  Chip
} from '@mui/material';
import {
  Settings as SettingsIcon,
  Fullscreen as FullscreenIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon
} from '@mui/icons-material';
import { createChart, IChartApi, ISeriesApi } from 'lightweight-charts';
import useWebSocket from 'react-use-websocket';

interface TradingViewProps {
  pair: string;
  poolId: string;
}

type TimeFrame = '1m' | '5m' | '15m' | '1h' | '4h' | '1d' | '1w';
type ChartType = 'candlestick' | 'line' | 'area';

export const TradingView: React.FC<TradingViewProps> = ({ pair, poolId }) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<any> | null>(null);

  const [timeframe, setTimeframe] = useState<TimeFrame>('15m');
  const [chartType, setChartType] = useState<ChartType>('candlestick');
  const [currentPrice, setCurrentPrice] = useState(0.075);
  const [priceChange, setPriceChange] = useState(0);
  const [volume24h, setVolume24h] = useState('0');
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  // WebSocket for real-time price updates
  const { lastMessage } = useWebSocket(
    `wss://api.deshchain.org/v1/moneyorder/ws/pools/${poolId}/trades`,
    {
      shouldReconnect: () => true,
      reconnectInterval: 3000
    }
  );

  // Initialize chart
  useEffect(() => {
    if (!chartContainerRef.current) return;

    const chart = createChart(chartContainerRef.current, {
      width: chartContainerRef.current.clientWidth,
      height: 500,
      layout: {
        background: { color: 'transparent' },
        textColor: '#333',
      },
      grid: {
        vertLines: { color: '#f0f0f0' },
        horzLines: { color: '#f0f0f0' },
      },
      crosshair: {
        mode: 1,
      },
      rightPriceScale: {
        borderColor: '#ccc',
      },
      timeScale: {
        borderColor: '#ccc',
        timeVisible: true,
        secondsVisible: false,
      },
    });

    chartRef.current = chart;

    // Add series based on chart type
    if (chartType === 'candlestick') {
      seriesRef.current = chart.addCandlestickSeries({
        upColor: '#26a69a',
        downColor: '#ef5350',
        borderVisible: false,
        wickUpColor: '#26a69a',
        wickDownColor: '#ef5350',
      });
    } else if (chartType === 'line') {
      seriesRef.current = chart.addLineSeries({
        color: '#2962FF',
        lineWidth: 2,
      });
    } else {
      seriesRef.current = chart.addAreaSeries({
        topColor: 'rgba(38, 166, 154, 0.56)',
        bottomColor: 'rgba(38, 166, 154, 0.04)',
        lineColor: 'rgba(38, 166, 154, 1)',
        lineWidth: 2,
      });
    }

    // Load initial data
    loadChartData();

    // Handle resize
    const handleResize = () => {
      if (chartContainerRef.current && chartRef.current) {
        chartRef.current.applyOptions({
          width: chartContainerRef.current.clientWidth,
        });
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      if (chartRef.current) {
        chartRef.current.remove();
      }
    };
  }, [chartType, timeframe, poolId]);

  // Update chart with WebSocket data
  useEffect(() => {
    if (lastMessage && seriesRef.current) {
      try {
        const data = JSON.parse(lastMessage.data);
        if (data.price) {
          setCurrentPrice(data.price);
          setPriceChange(data.priceChange24h || 0);
          setVolume24h(data.volume24h || '0');

          // Update chart
          if (chartType === 'candlestick') {
            seriesRef.current.update({
              time: Math.floor(Date.now() / 1000) as any,
              open: data.open || data.price,
              high: data.high || data.price,
              low: data.low || data.price,
              close: data.price,
            });
          } else {
            seriesRef.current.update({
              time: Math.floor(Date.now() / 1000) as any,
              value: data.price,
            });
          }
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    }
  }, [lastMessage, chartType]);

  const loadChartData = async () => {
    try {
      // This would fetch historical data from the API
      const mockData = generateMockData(timeframe);
      if (seriesRef.current) {
        seriesRef.current.setData(mockData as any);
      }
    } catch (error) {
      console.error('Failed to load chart data:', error);
    }
  };

  const generateMockData = (tf: TimeFrame) => {
    const now = Math.floor(Date.now() / 1000);
    const interval = getIntervalSeconds(tf);
    const data = [];

    for (let i = 100; i >= 0; i--) {
      const time = now - i * interval;
      const basePrice = 0.075;
      const variation = Math.random() * 0.005 - 0.0025;
      const open = basePrice + variation;
      const close = open + (Math.random() * 0.002 - 0.001);
      const high = Math.max(open, close) + Math.random() * 0.001;
      const low = Math.min(open, close) - Math.random() * 0.001;

      if (chartType === 'candlestick') {
        data.push({ time, open, high, low, close });
      } else {
        data.push({ time, value: close });
      }
    }

    return data;
  };

  const getIntervalSeconds = (tf: TimeFrame): number => {
    const intervals = {
      '1m': 60,
      '5m': 300,
      '15m': 900,
      '1h': 3600,
      '4h': 14400,
      '1d': 86400,
      '1w': 604800,
    };
    return intervals[tf];
  };

  const handleTimeframeChange = (event: React.MouseEvent<HTMLElement>, newTimeframe: TimeFrame | null) => {
    if (newTimeframe !== null) {
      setTimeframe(newTimeframe);
    }
  };

  const handleSettingsClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleChartTypeChange = (newType: ChartType) => {
    setChartType(newType);
    setAnchorEl(null);
  };

  return (
    <Paper sx={{ p: 2, height: '100%' }}>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Box>
          <Typography variant="h6" fontWeight="bold">
            {pair}
          </Typography>
          <Box display="flex" alignItems="center" gap={2}>
            <Typography variant="h4" fontWeight="bold">
              ₹{currentPrice.toFixed(4)}
            </Typography>
            <Chip
              label={`${priceChange >= 0 ? '+' : ''}${priceChange.toFixed(2)}%`}
              color={priceChange >= 0 ? 'success' : 'error'}
              icon={priceChange >= 0 ? <TrendingUpIcon /> : <TrendingDownIcon />}
              size="small"
            />
            <Typography variant="body2" color="text.secondary">
              Vol: ₹{parseFloat(volume24h).toLocaleString()}
            </Typography>
          </Box>
        </Box>

        <Box display="flex" gap={1}>
          <ToggleButtonGroup
            value={timeframe}
            exclusive
            onChange={handleTimeframeChange}
            size="small"
          >
            <ToggleButton value="1m">1m</ToggleButton>
            <ToggleButton value="5m">5m</ToggleButton>
            <ToggleButton value="15m">15m</ToggleButton>
            <ToggleButton value="1h">1h</ToggleButton>
            <ToggleButton value="4h">4h</ToggleButton>
            <ToggleButton value="1d">1D</ToggleButton>
            <ToggleButton value="1w">1W</ToggleButton>
          </ToggleButtonGroup>

          <IconButton size="small" onClick={handleSettingsClick}>
            <SettingsIcon />
          </IconButton>
          <IconButton size="small">
            <FullscreenIcon />
          </IconButton>
        </Box>
      </Box>

      {/* Chart Container */}
      <Box ref={chartContainerRef} sx={{ width: '100%', height: 'calc(100% - 100px)' }} />

      {/* Settings Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
      >
        <MenuItem onClick={() => handleChartTypeChange('candlestick')}>
          Candlestick Chart
        </MenuItem>
        <MenuItem onClick={() => handleChartTypeChange('line')}>
          Line Chart
        </MenuItem>
        <MenuItem onClick={() => handleChartTypeChange('area')}>
          Area Chart
        </MenuItem>
      </Menu>
    </Paper>
  );
};