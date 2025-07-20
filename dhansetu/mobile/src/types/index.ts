export * from './navigation';

export interface User {
  id: string;
  name: string;
  email?: string;
  phone?: string;
  dhanPataId?: string;
  avatar?: string;
  createdAt: Date;
}

export interface Wallet {
  address: string;
  publicKey: string;
  isHD: boolean;
  derivationPath: string;
  coinType: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN';
}

export interface Balance {
  denom: string;
  amount: string;
  displayAmount: string;
  usdValue?: string;
}

export interface Transaction {
  hash: string;
  type: 'send' | 'receive' | 'swap' | 'delegate' | 'undelegate' | 'claim';
  status: 'pending' | 'success' | 'failed';
  amount: string;
  denom: string;
  from: string;
  to: string;
  fee: string;
  timestamp: number;
  memo?: string;
  blockHeight?: number;
  culturalQuote?: string;
}

export interface Token {
  denom: string;
  symbol: string;
  name: string;
  decimals: number;
  icon?: string;
  price?: number;
  priceChange24h?: number;
}

export interface MoneyOrder {
  id: string;
  creator: string;
  recipient: string;
  amount: string;
  denom: string;
  status: 'pending' | 'claimed' | 'expired' | 'cancelled';
  createdAt: Date;
  expiresAt: Date;
  code: string;
  culturalQuote: string;
}

export interface Launch {
  id: string;
  name: string;
  symbol: string;
  description: string;
  totalSupply: string;
  initialPrice: string;
  creator: string;
  culturalTheme: string;
  launchDate: Date;
  status: 'upcoming' | 'live' | 'ended';
  raised: string;
  participants: number;
}