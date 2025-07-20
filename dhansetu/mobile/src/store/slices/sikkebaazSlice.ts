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

import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { DeshChainClient } from '@services/blockchain/deshchainClient';

export interface TokenLaunch {
  id: string;
  tokenName: string;
  tokenSymbol: string;
  description: string;
  totalSupply: string;
  targetAmount: string;
  raisedAmount: string;
  minimumInvestment: string;
  maximumInvestment: string;
  creator: string;
  creatorPincode: string;
  culturalCategory: string;
  imageUrl?: string;
  websiteUrl?: string;
  whitepaperUrl?: string;
  status: 'pending' | 'active' | 'successful' | 'failed' | 'vetoed';
  launchDate: string;
  endDate: string;
  tokenAddress?: string;
  participants: number;
  antiPumpConfig: {
    maxWalletPercent: number;
    tradingDelayHours: number;
    liquidityLockMonths: number;
    vetoThreshold: number;
  };
  festivalBonus?: {
    festivalId: string;
    bonusPercent: number;
  };
  culturalQuote?: string;
  vetoStatus?: {
    active: boolean;
    initiator: string;
    votes: number;
    threshold: number;
    endTime: string;
  };
}

export interface LaunchParticipation {
  launchId: string;
  participant: string;
  amount: string;
  timestamp: string;
  txHash: string;
}

export interface TradingMetrics {
  tokenAddress: string;
  totalVolume: string;
  dailyVolume: string;
  holders: number;
  price: string;
  priceChange24h: number;
  marketCap: string;
  liquidity: string;
}

interface CreateLaunchParams {
  tokenName: string;
  tokenSymbol: string;
  description: string;
  totalSupply: string;
  targetAmount: string;
  minimumInvestment: string;
  maximumInvestment: string;
  culturalCategory: string;
  creatorPincode: string;
  imageUrl?: string;
  websiteUrl?: string;
  whitepaperUrl?: string;
  antiPumpConfig: {
    maxWalletPercent: number;
    tradingDelayHours: number;
    liquidityLockMonths: number;
  };
  launchDurationDays: number;
}

interface SikkebaazState {
  launches: TokenLaunch[];
  myLaunches: TokenLaunch[];
  myParticipations: LaunchParticipation[];
  selectedLaunch: TokenLaunch | null;
  tradingMetrics: Record<string, TradingMetrics>;
  isLoading: boolean;
  error: string | null;
  stats: {
    totalLaunches: number;
    successfulLaunches: number;
    totalRaised: string;
    totalParticipants: number;
  };
}

const initialState: SikkebaazState = {
  launches: [],
  myLaunches: [],
  myParticipations: [],
  selectedLaunch: null,
  tradingMetrics: {},
  isLoading: false,
  error: null,
  stats: {
    totalLaunches: 0,
    successfulLaunches: 0,
    totalRaised: '0',
    totalParticipants: 0,
  },
};

// Async thunks
export const createTokenLaunch = createAsyncThunk(
  'sikkebaaz/createLaunch',
  async (params: CreateLaunchParams, { getState }) => {
    const state = getState() as any;
    const { currentAddress } = state.wallet;
    const client = new DeshChainClient({
      rpcEndpoint: process.env.EXPO_PUBLIC_RPC_URL!,
      apiEndpoint: process.env.EXPO_PUBLIC_API_URL!,
      chainId: 'deshchain-1',
      addressPrefix: 'desh',
      gasPrice: '0.025',
      gasDenom: 'namo',
    });

    // Calculate launch fee based on target amount
    const launchFee = calculateLaunchFee(params.targetAmount);
    
    // Apply festival bonus if active
    const festivalBonus = state.cultural.currentFestival
      ? {
          festivalId: state.cultural.currentFestival.id,
          bonusPercent: parseFloat(state.cultural.currentFestival.bonusRate),
        }
      : undefined;

    // Create launch message
    const msg = {
      typeUrl: '/deshchain.sikkebaaz.MsgCreateTokenLaunch',
      value: {
        creator: currentAddress,
        tokenName: params.tokenName,
        tokenSymbol: params.tokenSymbol,
        description: params.description,
        totalSupply: params.totalSupply,
        targetAmount: params.targetAmount,
        minimumInvestment: params.minimumInvestment,
        maximumInvestment: params.maximumInvestment,
        culturalCategory: params.culturalCategory,
        creatorPincode: params.creatorPincode,
        imageUrl: params.imageUrl,
        websiteUrl: params.websiteUrl,
        whitepaperUrl: params.whitepaperUrl,
        antiPumpConfig: {
          maxWalletPercent: params.antiPumpConfig.maxWalletPercent,
          tradingDelayHours: params.antiPumpConfig.tradingDelayHours,
          liquidityLockMonths: params.antiPumpConfig.liquidityLockMonths,
          vetoThreshold: 67, // 67% community veto threshold
        },
        launchDurationDays: params.launchDurationDays,
      },
    };

    // Sign and broadcast transaction
    const result = await client.sendNAMO(
      currentAddress,
      'desh1sikkebaaz...', // Sikkebaaz module address
      launchFee.toString(),
      `Creating ${params.tokenName} token launch`,
      getRandomCulturalQuote()
    );

    // Create launch object
    const launch: TokenLaunch = {
      id: result.transactionHash,
      ...params,
      creator: currentAddress,
      raisedAmount: '0',
      status: 'pending',
      launchDate: new Date().toISOString(),
      endDate: new Date(Date.now() + params.launchDurationDays * 24 * 60 * 60 * 1000).toISOString(),
      participants: 0,
      antiPumpConfig: {
        ...params.antiPumpConfig,
        vetoThreshold: 67,
      },
      festivalBonus,
      culturalQuote: getRandomCulturalQuote(),
    };

    return launch;
  }
);

export const participateInLaunch = createAsyncThunk(
  'sikkebaaz/participate',
  async ({ launchId, amount }: { launchId: string; amount: string }, { getState }) => {
    const state = getState() as any;
    const { currentAddress } = state.wallet;
    const launch = state.sikkebaaz.launches.find((l: TokenLaunch) => l.id === launchId);
    
    if (!launch) throw new Error('Launch not found');
    if (launch.status !== 'active') throw new Error('Launch is not active');
    
    // Validate investment amount
    const investmentAmount = parseFloat(amount);
    const minInvestment = parseFloat(launch.minimumInvestment);
    const maxInvestment = parseFloat(launch.maximumInvestment);
    
    if (investmentAmount < minInvestment) {
      throw new Error(`Minimum investment is ₹${minInvestment}`);
    }
    if (investmentAmount > maxInvestment) {
      throw new Error(`Maximum investment is ₹${maxInvestment}`);
    }

    // Create participation transaction
    const client = new DeshChainClient({
      rpcEndpoint: process.env.EXPO_PUBLIC_RPC_URL!,
      apiEndpoint: process.env.EXPO_PUBLIC_API_URL!,
      chainId: 'deshchain-1',
      addressPrefix: 'desh',
      gasPrice: '0.025',
      gasDenom: 'namo',
    });

    const result = await client.sendNAMO(
      currentAddress,
      launch.creator,
      amount,
      `Participating in ${launch.tokenName} launch`,
      launch.culturalQuote
    );

    const participation: LaunchParticipation = {
      launchId,
      participant: currentAddress,
      amount,
      timestamp: new Date().toISOString(),
      txHash: result.transactionHash,
    };

    return { participation, launchId };
  }
);

export const initiateVeto = createAsyncThunk(
  'sikkebaaz/initiateVeto',
  async ({ launchId, reason }: { launchId: string; reason: string }, { getState }) => {
    const state = getState() as any;
    const { currentAddress } = state.wallet;
    
    // Create veto initiation message
    const msg = {
      typeUrl: '/deshchain.sikkebaaz.MsgInitiateCommunityVeto',
      value: {
        initiator: currentAddress,
        launchId,
        reason,
      },
    };

    // Sign and broadcast
    // ... blockchain interaction

    return {
      launchId,
      vetoStatus: {
        active: true,
        initiator: currentAddress,
        votes: 1,
        threshold: 100, // Number of votes needed
        endTime: new Date(Date.now() + 72 * 60 * 60 * 1000).toISOString(), // 72 hours
      },
    };
  }
);

export const fetchLaunches = createAsyncThunk(
  'sikkebaaz/fetchLaunches',
  async ({ filter }: { filter?: 'trending' | 'new' | 'myTokens' } = {}) => {
    // Fetch from blockchain
    // For now, return mock data
    return getMockLaunches(filter);
  }
);

export const fetchTradingMetrics = createAsyncThunk(
  'sikkebaaz/fetchMetrics',
  async (tokenAddress: string) => {
    // Fetch trading metrics from blockchain
    const metrics: TradingMetrics = {
      tokenAddress,
      totalVolume: '1234567',
      dailyVolume: '123456',
      holders: 456,
      price: '0.123',
      priceChange24h: 15.67,
      marketCap: '12345678',
      liquidity: '1234567',
    };
    
    return metrics;
  }
);

const sikkebaazSlice = createSlice({
  name: 'sikkebaaz',
  initialState,
  reducers: {
    setSelectedLaunch: (state, action: PayloadAction<TokenLaunch | null>) => {
      state.selectedLaunch = action.payload;
    },
    updateLaunchStatus: (state, action: PayloadAction<{ launchId: string; status: TokenLaunch['status'] }>) => {
      const launch = state.launches.find(l => l.id === action.payload.launchId);
      if (launch) {
        launch.status = action.payload.status;
      }
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Create launch
      .addCase(createTokenLaunch.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createTokenLaunch.fulfilled, (state, action) => {
        state.isLoading = false;
        state.launches.push(action.payload);
        state.myLaunches.push(action.payload);
      })
      .addCase(createTokenLaunch.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to create launch';
      })
      
      // Participate
      .addCase(participateInLaunch.fulfilled, (state, action) => {
        state.myParticipations.push(action.payload.participation);
        
        // Update launch raised amount
        const launch = state.launches.find(l => l.id === action.payload.launchId);
        if (launch) {
          launch.raisedAmount = (
            parseFloat(launch.raisedAmount) + parseFloat(action.payload.participation.amount)
          ).toString();
          launch.participants += 1;
        }
      })
      
      // Initiate veto
      .addCase(initiateVeto.fulfilled, (state, action) => {
        const launch = state.launches.find(l => l.id === action.payload.launchId);
        if (launch) {
          launch.vetoStatus = action.payload.vetoStatus;
        }
      })
      
      // Fetch launches
      .addCase(fetchLaunches.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchLaunches.fulfilled, (state, action) => {
        state.isLoading = false;
        state.launches = action.payload;
      })
      
      // Fetch metrics
      .addCase(fetchTradingMetrics.fulfilled, (state, action) => {
        state.tradingMetrics[action.payload.tokenAddress] = action.payload;
      });
  },
});

export const { setSelectedLaunch, updateLaunchStatus, clearError } = sikkebaazSlice.actions;

export default sikkebaazSlice.reducer;

// Helper functions
function calculateLaunchFee(targetAmount: string): number {
  const target = parseFloat(targetAmount);
  
  // Tiered fee structure
  if (target <= 100000) return 1000; // ₹1,000 for up to ₹1 lakh
  if (target <= 1000000) return 5000; // ₹5,000 for up to ₹10 lakh
  if (target <= 10000000) return 25000; // ₹25,000 for up to ₹1 crore
  return 100000; // ₹1 lakh for above ₹1 crore
}

function getRandomCulturalQuote(): string {
  const quotes = [
    'सबका साथ, सबका विकास - Together we prosper',
    'वसुधैव कुटुम्बकम् - The world is one family',
    'Unity in diversity is India\'s strength',
    'From villages to cities, we grow together',
    'आत्मनिर्भर भारत - Self-reliant India',
  ];
  return quotes[Math.floor(Math.random() * quotes.length)];
}

function getMockLaunches(filter?: string): TokenLaunch[] {
  const allLaunches: TokenLaunch[] = [
    {
      id: '1',
      tokenName: 'Bollywood Coin',
      tokenSymbol: 'BOLLY',
      description: 'The official memecoin for Bollywood fans worldwide!',
      totalSupply: '1000000000',
      targetAmount: '1000000',
      raisedAmount: '750000',
      minimumInvestment: '100',
      maximumInvestment: '20000',
      creator: 'creator@dhan',
      creatorPincode: '400001',
      culturalCategory: 'Bollywood',
      status: 'active',
      launchDate: new Date().toISOString(),
      endDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
      participants: 234,
      antiPumpConfig: {
        maxWalletPercent: 2,
        tradingDelayHours: 24,
        liquidityLockMonths: 6,
        vetoThreshold: 67,
      },
    },
  ];
  
  return allLaunches;
}