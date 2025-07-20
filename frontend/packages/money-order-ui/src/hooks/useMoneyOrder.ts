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

import { useState, useCallback, useContext } from 'react';
import { MoneyOrderFormData, PoolInfo, SwapQuote, ReceiptData, LiquidityPosition } from '../types';
import { MoneyOrderContext } from '../providers/MoneyOrderProvider';

interface SwapParams {
  poolId: string;
  tokenIn: {
    denom: string;
    amount: string;
  };
  tokenOutDenom: string;
  maxSlippage?: number;
}

interface AddLiquidityParams {
  poolId: string;
  tokenA: {
    denom: string;
    amount: string;
  };
  tokenB: {
    denom: string;
    amount: string;
  };
  minShares?: string;
}

interface RemoveLiquidityParams {
  positionId: string;
  sharesToRemove: string;
  minTokenA?: string;
  minTokenB?: string;
}

export const useMoneyOrder = () => {
  const context = useContext(MoneyOrderContext);
  
  if (!context) {
    throw new Error('useMoneyOrder must be used within MoneyOrderProvider');
  }

  const { config, client } = context;
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createMoneyOrder = useCallback(async (data: MoneyOrderFormData): Promise<ReceiptData> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${config.apiUrl}/money-orders`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Chain-ID': config.chainId,
          'X-Culture-Language': data.culturalPreferences.language,
          'X-Festival-Context': data.culturalPreferences.theme
        },
        body: JSON.stringify({
          sender: data.sender,
          receiver: data.receiver,
          amount: {
            denom: data.amount.denom,
            amount: data.amount.value
          },
          pool_id: data.poolId,
          memo: data.memo,
          cultural_preferences: data.culturalPreferences,
          priority: data.priority
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error?.message || 'Failed to create money order');
      }

      const result = await response.json();
      return result.receipt;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [config, client]);

  const getSwapQuote = useCallback(async (params: SwapParams): Promise<SwapQuote> => {
    try {
      const queryParams = new URLSearchParams({
        pool_id: params.poolId,
        token_in_denom: params.tokenIn.denom,
        token_in_amount: params.tokenIn.amount,
        token_out_denom: params.tokenOutDenom,
        ...(params.maxSlippage && { max_slippage: params.maxSlippage.toString() })
      });

      const response = await fetch(`${config.apiUrl}/swaps/quote?${queryParams}`, {
        headers: {
          'X-Chain-ID': config.chainId
        }
      });

      if (!response.ok) {
        throw new Error('Failed to get swap quote');
      }

      const result = await response.json();
      return result.quote;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get quote';
      setError(errorMessage);
      throw err;
    }
  }, [config]);

  const executeSwap = useCallback(async (params: SwapParams): Promise<any> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${config.apiUrl}/swaps`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Chain-ID': config.chainId
        },
        body: JSON.stringify({
          pool_id: params.poolId,
          token_in: params.tokenIn,
          token_out_denom: params.tokenOutDenom,
          max_slippage: params.maxSlippage || config.maxSlippage
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error?.message || 'Failed to execute swap');
      }

      return await response.json();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Swap failed';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [config]);

  const addLiquidity = useCallback(async (params: AddLiquidityParams): Promise<LiquidityPosition> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${config.apiUrl}/liquidity/add`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Chain-ID': config.chainId
        },
        body: JSON.stringify({
          pool_id: params.poolId,
          token_a_amount: params.tokenA,
          token_b_amount: params.tokenB,
          min_shares: params.minShares
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error?.message || 'Failed to add liquidity');
      }

      const result = await response.json();
      return result.position;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to add liquidity';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [config]);

  const removeLiquidity = useCallback(async (params: RemoveLiquidityParams): Promise<any> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${config.apiUrl}/liquidity/remove`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Chain-ID': config.chainId
        },
        body: JSON.stringify({
          position_id: params.positionId,
          shares_to_remove: params.sharesToRemove,
          min_token_a: params.minTokenA,
          min_token_b: params.minTokenB
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error?.message || 'Failed to remove liquidity');
      }

      return await response.json();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to remove liquidity';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [config]);

  const getMoneyOrderStatus = useCallback(async (orderId: string): Promise<any> => {
    try {
      const response = await fetch(`${config.apiUrl}/money-orders/${orderId}`, {
        headers: {
          'X-Chain-ID': config.chainId
        }
      });

      if (!response.ok) {
        throw new Error('Failed to get money order status');
      }

      return await response.json();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get status';
      setError(errorMessage);
      throw err;
    }
  }, [config]);

  const getReceipt = useCallback(async (receiptId: string): Promise<ReceiptData> => {
    try {
      const response = await fetch(`${config.apiUrl}/receipts/${receiptId}`, {
        headers: {
          'X-Chain-ID': config.chainId
        }
      });

      if (!response.ok) {
        throw new Error('Failed to get receipt');
      }

      const result = await response.json();
      return result.receipt;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get receipt';
      setError(errorMessage);
      throw err;
    }
  }, [config]);

  const getLiquidityPositions = useCallback(async (owner: string): Promise<LiquidityPosition[]> => {
    try {
      const response = await fetch(`${config.apiUrl}/liquidity/positions?owner=${owner}`, {
        headers: {
          'X-Chain-ID': config.chainId
        }
      });

      if (!response.ok) {
        throw new Error('Failed to get liquidity positions');
      }

      const result = await response.json();
      return result.positions;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get positions';
      setError(errorMessage);
      throw err;
    }
  }, [config]);

  const getTransactionHistory = useCallback(async (address: string, limit = 50): Promise<any[]> => {
    try {
      const response = await fetch(`${config.apiUrl}/analytics/transactions?address=${address}&limit=${limit}`, {
        headers: {
          'X-Chain-ID': config.chainId
        }
      });

      if (!response.ok) {
        throw new Error('Failed to get transaction history');
      }

      const result = await response.json();
      return result.transactions;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get transaction history';
      setError(errorMessage);
      throw err;
    }
  }, [config]);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    // State
    isLoading,
    error,
    
    // Money Order operations
    createMoneyOrder,
    getMoneyOrderStatus,
    getReceipt,
    
    // Trading operations
    getSwapQuote,
    executeSwap,
    
    // Liquidity operations
    addLiquidity,
    removeLiquidity,
    getLiquidityPositions,
    
    // Analytics
    getTransactionHistory,
    
    // Utilities
    clearError
  };
};