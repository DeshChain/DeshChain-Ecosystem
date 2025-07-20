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

export const config = {
  // API Configuration
  api: {
    baseUrl: process.env.NEXT_PUBLIC_API_URL || 'https://api.deshchain.org/v1',
    timeout: 30000,
    retries: 3
  },

  // Blockchain Configuration
  blockchain: {
    chainId: process.env.NEXT_PUBLIC_CHAIN_ID || 'deshchain-1',
    chainName: 'DeshChain',
    rpcUrl: process.env.NEXT_PUBLIC_RPC_URL || 'https://rpc.deshchain.org',
    restUrl: process.env.NEXT_PUBLIC_REST_URL || 'https://api.deshchain.org',
    bech32Prefix: 'desh',
    coinDenom: 'NAMO',
    coinMinimalDenom: 'unamo',
    coinDecimals: 6,
    gasPrice: '0.025unamo'
  },

  // Money Order Configuration
  moneyOrder: {
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'https://api.deshchain.org/v1/moneyorder',
    chainId: process.env.NEXT_PUBLIC_CHAIN_ID || 'deshchain-1',
    defaultLanguage: 'en',
    enableCulturalFeatures: true,
    enableFestivalThemes: true,
    enablePatriotismRewards: true,
    maxSlippage: 0.05,
    defaultPoolType: 'fixed_rate',
    autoSelectPool: true,
    theme: {
      primary: '#FF6B35',
      secondary: '#138808',
      accent: '#000080',
      background: '#FAFAFA',
      surface: '#FFFFFF',
      text: '#212121',
      culturalElements: {
        borderStyle: 'traditional',
        pattern: 'mandala',
        iconSet: 'indian_classical'
      }
    }
  },

  // Cultural Configuration
  defaultLanguage: 'en',
  supportedLanguages: [
    'en', 'hi', 'bn', 'te', 'ta', 'mr', 'ur', 'gu', 
    'kn', 'ml', 'pa', 'or', 'as', 'ks', 'ne', 'si',
    'bho', 'raj', 'mai', 'sa', 'sd', 'kok'
  ],

  // Feature Flags
  features: {
    cultural: true,
    festivals: true,
    patriotism: true,
    quotes: true,
    themes: true,
    animations: true,
    addressBook: true,
    qrCode: true,
    notifications: true,
    analytics: true
  },

  // Transaction Limits
  limits: {
    minAmount: '1', // 1 NAMO
    maxAmount: '10000000', // 10M NAMO
    maxMemoLength: 200,
    maxSlippage: 0.1 // 10%
  },

  // UI Configuration
  ui: {
    animationDuration: 300,
    toastDuration: 4000,
    poolRefreshInterval: 60000, // 1 minute
    quoteRefreshInterval: 300000, // 5 minutes
    festivalCheckInterval: 3600000 // 1 hour
  },

  // Pool Configuration
  pools: {
    recommendedTypes: ['fixed_rate', 'village'],
    minLiquidity: '1000000', // 1M minimum liquidity
    maxPriceImpact: 0.05 // 5% max price impact
  },

  // Festival Bonuses
  festivalBonuses: {
    diwali: 0.15,
    holi: 0.10,
    independence_day: 0.20,
    republic_day: 0.20,
    eid: 0.12,
    durga_puja: 0.12,
    ganesh_chaturthi: 0.10,
    navratri: 0.08
  },

  // Priority Fees
  priorityFees: {
    standard: 1.0,
    fast: 1.5,
    instant: 2.0
  },

  // Analytics
  analytics: {
    trackingEnabled: process.env.NEXT_PUBLIC_ANALYTICS_ENABLED === 'true',
    trackingId: process.env.NEXT_PUBLIC_GA_TRACKING_ID || '',
    sendAnalytics: true
  },

  // External Services
  services: {
    ipfsGateway: 'https://ipfs.deshchain.org',
    explorerUrl: 'https://explorer.deshchain.org',
    docsUrl: 'https://docs.deshchain.org',
    supportUrl: 'https://support.deshchain.org'
  }
};