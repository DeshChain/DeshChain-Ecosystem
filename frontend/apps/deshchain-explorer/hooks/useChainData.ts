'use client'

import { useState, useEffect } from 'react'
import { useQuery } from 'react-query'
import { useExplorer } from '@/providers/ExplorerProvider'

interface ChainStats {
  blockHeight: number
  totalTransactions: number
  activeValidators: number
  namoSupply: number
  blockTrend?: number
  txTrend?: number
  supplyTrend?: number
  avgBlockTime: number
  networkHealth: 'healthy' | 'degraded' | 'unhealthy'
  tps: number // transactions per second
  totalAccounts: number
}

interface Transaction {
  hash: string
  height: number
  timestamp: Date
  fee: string
  gas: string
  memo: string
  success: boolean
  type: string
  amount?: string
  from?: string
  to?: string
}

interface LendingStats {
  krishiMitra: {
    totalLoans: number
    totalDisbursed: string
    averageRate: number
    activeLoans: number
    defaultRate: number
  }
  vyavasayaMitra: {
    totalLoans: number
    totalDisbursed: string
    averageRate: number
    activeLoans: number
    defaultRate: number
  }
  shikshamitra: {
    totalLoans: number
    totalDisbursed: string
    averageRate: number
    activeLoans: number
    defaultRate: number
  }
  combined: {
    totalDisbursed: string
    avgInterestRate: number
    totalBorrowers: number
    totalLenders: number
  }
}

export function useChainData() {
  const { chainInfo, isConnected, client } = useExplorer()
  const [prevStats, setPrevStats] = useState<ChainStats | null>(null)

  // Fetch chain statistics
  const { 
    data: chainStats, 
    isLoading: statsLoading, 
    error: statsError,
    refetch: refetchStats
  } = useQuery(
    ['chainStats', chainInfo?.blockHeight],
    async () => {
      if (!chainInfo || !isConnected) return null

      // Mock data - in production this would come from actual chain queries
      const stats: ChainStats = {
        blockHeight: chainInfo.blockHeight,
        totalTransactions: 2547892,
        activeValidators: chainInfo.activeValidators,
        namoSupply: chainInfo.namoSupply,
        avgBlockTime: chainInfo.avgBlockTime,
        networkHealth: chainInfo.networkHealth,
        tps: Math.floor(Math.random() * 50) + 20, // 20-70 TPS
        totalAccounts: 156789,
        blockTrend: prevStats ? chainInfo.blockHeight - prevStats.blockHeight : 0,
        txTrend: prevStats ? Math.floor(Math.random() * 1000) : 0,
        supplyTrend: prevStats ? 0.5 : 0, // 0.5% increase
      }

      setPrevStats(stats)
      return stats
    },
    {
      enabled: !!chainInfo && isConnected,
      refetchInterval: 30000, // Refetch every 30 seconds
      staleTime: 15000, // Consider stale after 15 seconds
    }
  )

  // Fetch recent transactions
  const { 
    data: recentTransactions, 
    isLoading: transactionsLoading,
    error: transactionsError 
  } = useQuery(
    ['recentTransactions', chainInfo?.blockHeight],
    async () => {
      if (!chainInfo || !isConnected) return []

      // Mock data - in production this would fetch from actual transactions
      const transactions: Transaction[] = Array.from({ length: 10 }, (_, i) => ({
        hash: `0x${Math.random().toString(16).substr(2, 64)}`,
        height: chainInfo.blockHeight - i,
        timestamp: new Date(Date.now() - i * 6000), // 6 seconds apart
        fee: `${(Math.random() * 0.1).toFixed(6)} NAMO`,
        gas: `${Math.floor(Math.random() * 200000 + 100000)}`,
        memo: i % 3 === 0 ? 'Festival bonus transaction' : '',
        success: Math.random() > 0.05, // 95% success rate
        type: ['send', 'delegate', 'vote', 'loan_application', 'sikkebaaz_launch'][Math.floor(Math.random() * 5)],
        amount: `${(Math.random() * 1000).toFixed(2)} NAMO`,
        from: `deshchain1${Math.random().toString(36).substr(2, 39)}`,
        to: `deshchain1${Math.random().toString(36).substr(2, 39)}`,
      }))

      return transactions
    },
    {
      enabled: !!chainInfo && isConnected,
      refetchInterval: 15000, // Refetch every 15 seconds
    }
  )

  // Fetch lending module statistics
  const { 
    data: lendingStats, 
    isLoading: lendingLoading,
    error: lendingError 
  } = useQuery(
    'lendingStats',
    async () => {
      if (!isConnected) return null

      // Mock data - in production this would come from the lending modules
      const stats: LendingStats = {
        krishiMitra: {
          totalLoans: 2547,
          totalDisbursed: '₹12.4 Cr',
          averageRate: 7.2,
          activeLoans: 1892,
          defaultRate: 2.1,
        },
        vyavasayaMitra: {
          totalLoans: 1823,
          totalDisbursed: '₹45.7 Cr',
          averageRate: 9.8,
          activeLoans: 1456,
          defaultRate: 3.4,
        },
        shikshamitra: {
          totalLoans: 3421,
          totalDisbursed: '₹67.2 Cr',
          averageRate: 5.6,
          activeLoans: 2987,
          defaultRate: 1.2,
        },
        combined: {
          totalDisbursed: '₹125.3 Cr',
          avgInterestRate: 7.2,
          totalBorrowers: 6891,
          totalLenders: 892,
        }
      }

      return stats
    },
    {
      enabled: isConnected,
      refetchInterval: 60000, // Refetch every minute
      staleTime: 30000,
    }
  )

  const loading = statsLoading || transactionsLoading || lendingLoading
  const error = statsError || transactionsError || lendingError

  return {
    chainStats,
    recentTransactions: recentTransactions || [],
    lendingStats,
    loading,
    error,
    refetch: () => {
      refetchStats()
    }
  }
}

// Hook for real-time blockchain data
export function useRealtimeData() {
  const [realtimeStats, setRealtimeStats] = useState({
    currentTPS: 0,
    memPoolSize: 0,
    lastBlockTime: 0,
    networkLatency: 0,
  })

  useEffect(() => {
    const interval = setInterval(() => {
      setRealtimeStats({
        currentTPS: Math.floor(Math.random() * 50) + 20,
        memPoolSize: Math.floor(Math.random() * 1000) + 100,
        lastBlockTime: 6000 + Math.floor(Math.random() * 2000), // 6-8 seconds
        networkLatency: Math.floor(Math.random() * 100) + 50, // 50-150ms
      })
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  return realtimeStats
}

// Hook for historical data
export function useHistoricalData(timeframe: '1h' | '24h' | '7d' | '30d' = '24h') {
  return useQuery(
    ['historicalData', timeframe],
    async () => {
      // Mock historical data
      const now = Date.now()
      const points = timeframe === '1h' ? 12 : timeframe === '24h' ? 24 : timeframe === '7d' ? 7 : 30
      const interval = timeframe === '1h' ? 5 * 60 * 1000 : 
                     timeframe === '24h' ? 60 * 60 * 1000 :
                     timeframe === '7d' ? 24 * 60 * 60 * 1000 :
                     24 * 60 * 60 * 1000

      const data = Array.from({ length: points }, (_, i) => {
        const timestamp = now - (points - 1 - i) * interval
        return {
          timestamp,
          blockHeight: 1000000 + i * 100,
          transactions: Math.floor(Math.random() * 1000) + 500,
          tps: Math.floor(Math.random() * 30) + 20,
          validators: 100 + Math.floor(Math.random() * 20),
          avgBlockTime: 6000 + Math.floor(Math.random() * 2000),
        }
      })

      return data
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
      refetchInterval: timeframe === '1h' ? 30000 : 60000,
    }
  )
}