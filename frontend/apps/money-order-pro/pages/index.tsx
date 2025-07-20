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
  Grid,
  Paper,
  Tab,
  Tabs,
  useTheme,
  useMediaQuery
} from '@mui/material';
import {
  ShowChart as ChartIcon,
  SwapHoriz as SwapIcon,
  AccountBalance as LiquidityIcon,
  Analytics as AnalyticsIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';

import { TradingView } from '../components/TradingView';
import { OrderBook } from '../components/OrderBook';
import { TradeForm } from '../components/TradeForm';
import { PositionsPanel } from '../components/PositionsPanel';
import { MarketDepth } from '../components/MarketDepth';
import { PoolAnalytics } from '../components/PoolAnalytics';
import { LiquidityManager } from '../components/LiquidityManager';
import { AdvancedSettings } from '../components/AdvancedSettings';
import { MarketOverview } from '../components/MarketOverview';
import { useCulturalContext } from '@deshchain/money-order-ui';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`trading-tabpanel-${index}`}
      aria-labelledby={`trading-tab-${index}`}
      {...other}
    >
      {value === index && <Box>{children}</Box>}
    </div>
  );
}

const ProTrading: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [activeTab, setActiveTab] = useState(0);
  const [selectedPool, setSelectedPool] = useState<string>('1');
  const [selectedPair, setSelectedPair] = useState('NAMO/INR');

  const { currentFestival } = useCulturalContext();

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const tabs = [
    { label: 'Trading', icon: <SwapIcon /> },
    { label: 'Liquidity', icon: <LiquidityIcon /> },
    { label: 'Analytics', icon: <AnalyticsIcon /> },
    { label: 'Settings', icon: <SettingsIcon /> }
  ];

  return (
    <Box sx={{ flexGrow: 1, height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* Market Overview Bar */}
      <Paper elevation={0} sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <MarketOverview />
      </Paper>

      {/* Main Content */}
      <Box sx={{ flexGrow: 1, overflow: 'hidden' }}>
        <Grid container sx={{ height: '100%' }}>
          {/* Left Panel - Charts & Analytics */}
          <Grid item xs={12} md={8} sx={{ height: '100%' }}>
            <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              {/* Tabs */}
              <Paper sx={{ borderRadius: 0 }}>
                <Tabs
                  value={activeTab}
                  onChange={handleTabChange}
                  variant={isMobile ? 'scrollable' : 'standard'}
                  scrollButtons={isMobile ? 'auto' : false}
                >
                  {tabs.map((tab, index) => (
                    <Tab
                      key={index}
                      icon={tab.icon}
                      label={!isMobile ? tab.label : undefined}
                      iconPosition="start"
                    />
                  ))}
                </Tabs>
              </Paper>

              {/* Tab Content */}
              <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
                <TabPanel value={activeTab} index={0}>
                  <Grid container sx={{ height: '100%' }}>
                    <Grid item xs={12} lg={8}>
                      <TradingView
                        pair={selectedPair}
                        poolId={selectedPool}
                      />
                    </Grid>
                    <Grid item xs={12} lg={4}>
                      <OrderBook
                        poolId={selectedPool}
                        pair={selectedPair}
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <MarketDepth
                        poolId={selectedPool}
                        pair={selectedPair}
                      />
                    </Grid>
                  </Grid>
                </TabPanel>

                <TabPanel value={activeTab} index={1}>
                  <LiquidityManager
                    onPoolSelect={setSelectedPool}
                  />
                </TabPanel>

                <TabPanel value={activeTab} index={2}>
                  <PoolAnalytics
                    poolId={selectedPool}
                  />
                </TabPanel>

                <TabPanel value={activeTab} index={3}>
                  <AdvancedSettings />
                </TabPanel>
              </Box>
            </Box>
          </Grid>

          {/* Right Panel - Trade Form & Positions */}
          <Grid item xs={12} md={4} sx={{ height: '100%', borderLeft: 1, borderColor: 'divider' }}>
            <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              {/* Trade Form */}
              <Box sx={{ p: 2 }}>
                <TradeForm
                  poolId={selectedPool}
                  pair={selectedPair}
                  festivalBonus={currentFestival?.bonusRate}
                />
              </Box>

              {/* Positions */}
              <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
                <PositionsPanel />
              </Box>
            </Box>
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
};