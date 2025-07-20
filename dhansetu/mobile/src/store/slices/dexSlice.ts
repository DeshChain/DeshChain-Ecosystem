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

export interface MoneyOrder {
  id: string;
  orderId: string;
  creator: string;
  creatorDhanPata?: string;
  receiver: string;
  receiverDhanPata?: string;
  amount: string;
  fee: string;
  totalAmount: string;
  status: 'pending' | 'accepted' | 'completed' | 'cancelled' | 'expired';
  paymentMode: 'namo' | 'fiat';
  createdAt: string;
  expiresAt: string;
  acceptedAt?: string;
  completedAt?: string;
  cancelledAt?: string;
  memo?: string;
  culturalQuote?: string;
  pincode?: string;
  escrowAddress?: string;
  txHash?: string;
  festivalBonus?: number;
}

interface CreateOrderParams {
  receiver: string;
  amount: string;
  paymentMode: 'namo' | 'fiat';
  memo?: string;
  pincode?: string;
  expiryHours?: number;
}

interface DexState {
  activeOrders: MoneyOrder[];
  completedOrders: MoneyOrder[];
  myOrders: MoneyOrder[];
  ordersForMe: MoneyOrder[];
  selectedOrder: MoneyOrder | null;
  isLoading: boolean;
  error: string | null;
  stats: {
    totalVolume: string;
    totalOrders: number;
    successRate: number;
    averageFee: string;
  };
}

const initialState: DexState = {
  activeOrders: [],
  completedOrders: [],
  myOrders: [],
  ordersForMe: [],
  selectedOrder: null,
  isLoading: false,
  error: null,
  stats: {
    totalVolume: '0',
    totalOrders: 0,
    successRate: 0,
    averageFee: '0',
  },
};

// Async thunks
export const createMoneyOrder = createAsyncThunk(
  'dex/createOrder',
  async (params: CreateOrderParams, { getState }) => {
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

    // Calculate fee based on amount and festival bonus
    const baseFee = parseFloat(params.amount) * 0.01; // 1% base fee
    const festivalBonus = state.cultural.festivalBonusRate || 0;
    const fee = baseFee * (1 - festivalBonus);
    const totalAmount = parseFloat(params.amount) + fee;

    // Create order message
    const msg = {
      typeUrl: '/deshchain.moneyorder.MsgCreateOrder',
      value: {
        creator: currentAddress,
        receiver: params.receiver,
        amount: {
          denom: 'namo',
          amount: Math.floor(parseFloat(params.amount) * 1e6).toString(), // Convert to unamo
        },
        paymentMode: params.paymentMode,
        memo: params.memo,
        pincode: params.pincode,
        expiryHours: params.expiryHours || 24,
      },
    };

    // Sign and broadcast transaction
    const result = await client.sendNAMO(
      currentAddress,
      'desh1moneyorder...', // Money order module address
      totalAmount.toString(),
      params.memo,
      getRandomCulturalQuote()
    );

    // Return created order
    const order: MoneyOrder = {
      id: result.transactionHash,
      orderId: `MO${Date.now()}`,
      creator: currentAddress,
      receiver: params.receiver,
      amount: params.amount,
      fee: fee.toFixed(2),
      totalAmount: totalAmount.toFixed(2),
      status: 'pending',
      paymentMode: params.paymentMode,
      createdAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + (params.expiryHours || 24) * 60 * 60 * 1000).toISOString(),
      memo: params.memo,
      culturalQuote: getRandomCulturalQuote(),
      pincode: params.pincode,
      txHash: result.transactionHash,
      festivalBonus: festivalBonus * 100,
    };

    return order;
  }
);

export const fetchActiveOrders = createAsyncThunk(
  'dex/fetchActive',
  async (_, { getState }) => {
    const state = getState() as any;
    const { currentAddress } = state.wallet;
    
    // Fetch from blockchain
    // For now, return mock data
    return getMockActiveOrders(currentAddress);
  }
);

export const acceptMoneyOrder = createAsyncThunk(
  'dex/acceptOrder',
  async (orderId: string, { getState }) => {
    const state = getState() as any;
    const order = state.dex.activeOrders.find((o: MoneyOrder) => o.id === orderId);
    
    if (!order) throw new Error('Order not found');
    
    // Accept order on blockchain
    const updatedOrder: MoneyOrder = {
      ...order,
      status: 'accepted',
      acceptedAt: new Date().toISOString(),
    };
    
    return updatedOrder;
  }
);

export const completeMoneyOrder = createAsyncThunk(
  'dex/completeOrder',
  async (orderId: string, { getState }) => {
    const state = getState() as any;
    const order = state.dex.activeOrders.find((o: MoneyOrder) => o.id === orderId);
    
    if (!order) throw new Error('Order not found');
    if (order.status !== 'accepted') throw new Error('Order must be accepted first');
    
    // Complete order on blockchain
    const updatedOrder: MoneyOrder = {
      ...order,
      status: 'completed',
      completedAt: new Date().toISOString(),
    };
    
    return updatedOrder;
  }
);

export const cancelMoneyOrder = createAsyncThunk(
  'dex/cancelOrder',
  async (orderId: string, { getState }) => {
    const state = getState() as any;
    const order = state.dex.activeOrders.find((o: MoneyOrder) => o.id === orderId);
    
    if (!order) throw new Error('Order not found');
    if (order.status !== 'pending') throw new Error('Only pending orders can be cancelled');
    
    // Cancel order on blockchain
    const updatedOrder: MoneyOrder = {
      ...order,
      status: 'cancelled',
      cancelledAt: new Date().toISOString(),
    };
    
    return updatedOrder;
  }
);

const dexSlice = createSlice({
  name: 'dex',
  initialState,
  reducers: {
    setSelectedOrder: (state, action: PayloadAction<MoneyOrder | null>) => {
      state.selectedOrder = action.payload;
    },
    updateOrderStatus: (state, action: PayloadAction<{ orderId: string; status: MoneyOrder['status'] }>) => {
      const order = state.activeOrders.find(o => o.id === action.payload.orderId);
      if (order) {
        order.status = action.payload.status;
        
        // Move to completed if necessary
        if (['completed', 'cancelled', 'expired'].includes(action.payload.status)) {
          state.activeOrders = state.activeOrders.filter(o => o.id !== action.payload.orderId);
          state.completedOrders.push(order);
        }
      }
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Create order
      .addCase(createMoneyOrder.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createMoneyOrder.fulfilled, (state, action) => {
        state.isLoading = false;
        state.activeOrders.push(action.payload);
        state.myOrders.push(action.payload);
      })
      .addCase(createMoneyOrder.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to create order';
      })
      
      // Fetch active orders
      .addCase(fetchActiveOrders.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchActiveOrders.fulfilled, (state, action) => {
        state.isLoading = false;
        state.activeOrders = action.payload;
      })
      .addCase(fetchActiveOrders.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch orders';
      })
      
      // Accept order
      .addCase(acceptMoneyOrder.fulfilled, (state, action) => {
        const index = state.activeOrders.findIndex(o => o.id === action.payload.id);
        if (index !== -1) {
          state.activeOrders[index] = action.payload;
        }
      })
      
      // Complete order
      .addCase(completeMoneyOrder.fulfilled, (state, action) => {
        state.activeOrders = state.activeOrders.filter(o => o.id !== action.payload.id);
        state.completedOrders.push(action.payload);
      })
      
      // Cancel order
      .addCase(cancelMoneyOrder.fulfilled, (state, action) => {
        state.activeOrders = state.activeOrders.filter(o => o.id !== action.payload.id);
        state.completedOrders.push(action.payload);
      });
  },
});

export const { setSelectedOrder, updateOrderStatus, clearError } = dexSlice.actions;

export default dexSlice.reducer;

// Helper functions
function getRandomCulturalQuote(): string {
  const quotes = [
    'वसुधैव कुटुम्बकम् - The world is one family',
    'सत्यमेव जयते - Truth alone triumphs',
    'अहिंसा परमो धर्मः - Non-violence is the supreme duty',
    'विद्या ददाति विनयं - Knowledge gives humility',
    'Unity in diversity is India\'s strength',
  ];
  return quotes[Math.floor(Math.random() * quotes.length)];
}

function getMockActiveOrders(userAddress: string): MoneyOrder[] {
  return [
    {
      id: '1',
      orderId: 'MO2024001',
      creator: userAddress,
      receiver: 'ramesh@dhan',
      amount: '5000',
      fee: '45',
      totalAmount: '5045',
      status: 'pending',
      paymentMode: 'namo',
      createdAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
      culturalQuote: 'Unity in diversity',
      pincode: '110001',
      festivalBonus: 10,
    },
  ];
}