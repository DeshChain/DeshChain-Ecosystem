'use client'

import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react'
import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'

interface ChainInfo {
  chainId: string
  rpcEndpoint: string
  restEndpoint: string
  blockHeight: number
  blockTime: number
  totalTransactions: number
  activeValidators: number
  namoSupply: number
  avgBlockTime: number
  networkHealth: 'healthy' | 'degraded' | 'unhealthy'
}

interface ExplorerState {
  chainInfo: ChainInfo | null
  client: StargateClient | null
  tmClient: Tendermint34Client | null
  isConnected: boolean
  isLoading: boolean
  error: string | null
  lastUpdated: Date | null
}

type ExplorerAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_CHAIN_INFO'; payload: ChainInfo }
  | { type: 'SET_CLIENTS'; payload: { client: StargateClient; tmClient: Tendermint34Client } }
  | { type: 'SET_CONNECTED'; payload: boolean }
  | { type: 'UPDATE_LAST_UPDATED' }

const initialState: ExplorerState = {
  chainInfo: null,
  client: null,
  tmClient: null,
  isConnected: false,
  isLoading: false,
  error: null,
  lastUpdated: null,
}

function explorerReducer(state: ExplorerState, action: ExplorerAction): ExplorerState {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload }
    case 'SET_ERROR':
      return { ...state, error: action.payload, isLoading: false }
    case 'SET_CHAIN_INFO':
      return { ...state, chainInfo: action.payload, error: null }
    case 'SET_CLIENTS':
      return { 
        ...state, 
        client: action.payload.client, 
        tmClient: action.payload.tmClient,
        isConnected: true,
        error: null 
      }
    case 'SET_CONNECTED':
      return { ...state, isConnected: action.payload }
    case 'UPDATE_LAST_UPDATED':
      return { ...state, lastUpdated: new Date() }
    default:
      return state
  }
}

interface ExplorerContextType extends ExplorerState {
  connect: () => Promise<void>
  disconnect: () => void
  refreshChainInfo: () => Promise<void>
  getTransaction: (hash: string) => Promise<any>
  getBlock: (height: number) => Promise<any>
  searchAddress: (address: string) => Promise<any>
}

const ExplorerContext = createContext<ExplorerContextType | undefined>(undefined)

interface ExplorerProviderProps {
  children: ReactNode
}

export function ExplorerProvider({ children }: ExplorerProviderProps) {
  const [state, dispatch] = useReducer(explorerReducer, initialState)

  const rpcEndpoint = process.env.NEXT_PUBLIC_RPC_ENDPOINT || 'http://localhost:26657'
  const restEndpoint = process.env.NEXT_PUBLIC_REST_ENDPOINT || 'http://localhost:1317'

  const connect = async () => {
    dispatch({ type: 'SET_LOADING', payload: true })
    
    try {
      // Connect to Tendermint RPC
      const tmClient = await Tendermint34Client.connect(rpcEndpoint)
      
      // Connect to Stargate client
      const client = await StargateClient.connectWithSigner(rpcEndpoint, {} as any)
      
      dispatch({ 
        type: 'SET_CLIENTS', 
        payload: { client, tmClient } 
      })
      
      // Fetch initial chain info
      await refreshChainInfo()
      
    } catch (error) {
      console.error('Failed to connect to chain:', error)
      dispatch({ 
        type: 'SET_ERROR', 
        payload: error instanceof Error ? error.message : 'Connection failed'
      })
    }
  }

  const disconnect = () => {
    if (state.client) {
      state.client.disconnect()
    }
    if (state.tmClient) {
      state.tmClient.disconnect()
    }
    dispatch({ type: 'SET_CONNECTED', payload: false })
  }

  const refreshChainInfo = async () => {
    if (!state.tmClient) return

    try {
      const status = await state.tmClient.status()
      const validators = await state.tmClient.validatorsAll()
      
      // Get additional chain data
      const chainInfo: ChainInfo = {
        chainId: status.nodeInfo.network,
        rpcEndpoint,
        restEndpoint,
        blockHeight: status.syncInfo.latestBlockHeight,
        blockTime: new Date(status.syncInfo.latestBlockTime).getTime(),
        totalTransactions: 0, // This would need to be fetched from a custom API
        activeValidators: validators.validators.length,
        namoSupply: 1000000000, // This would need to be fetched from bank module
        avgBlockTime: 6000, // 6 seconds average
        networkHealth: status.syncInfo.catchingUp ? 'degraded' : 'healthy',
      }

      dispatch({ type: 'SET_CHAIN_INFO', payload: chainInfo })
      dispatch({ type: 'UPDATE_LAST_UPDATED' })
      
    } catch (error) {
      console.error('Failed to refresh chain info:', error)
      dispatch({ 
        type: 'SET_ERROR', 
        payload: error instanceof Error ? error.message : 'Failed to refresh chain info'
      })
    }
  }

  const getTransaction = async (hash: string) => {
    if (!state.tmClient) throw new Error('Not connected to chain')
    
    try {
      const tx = await state.tmClient.tx({ hash: Buffer.from(hash, 'hex') })
      return tx
    } catch (error) {
      throw error
    }
  }

  const getBlock = async (height: number) => {
    if (!state.tmClient) throw new Error('Not connected to chain')
    
    try {
      const block = await state.tmClient.block(height)
      return block
    } catch (error) {
      throw error
    }
  }

  const searchAddress = async (address: string) => {
    if (!state.client) throw new Error('Not connected to chain')
    
    try {
      const account = await state.client.getAccount(address)
      const balances = await state.client.getAllBalances(address)
      
      return {
        account,
        balances,
        address,
      }
    } catch (error) {
      throw error
    }
  }

  // Auto-connect on mount
  useEffect(() => {
    connect()
    
    // Set up periodic refresh
    const interval = setInterval(() => {
      if (state.isConnected) {
        refreshChainInfo()
      }
    }, 30000) // Refresh every 30 seconds

    return () => {
      clearInterval(interval)
      disconnect()
    }
  }, [])

  const value: ExplorerContextType = {
    ...state,
    connect,
    disconnect,
    refreshChainInfo,
    getTransaction,
    getBlock,
    searchAddress,
  }

  return (
    <ExplorerContext.Provider value={value}>
      {children}
    </ExplorerContext.Provider>
  )
}

export function useExplorer() {
  const context = useContext(ExplorerContext)
  if (context === undefined) {
    throw new Error('useExplorer must be used within an ExplorerProvider')
  }
  return context
}