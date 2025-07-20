import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useAppDispatch, useAppSelector } from '@store/index';
import { HDWallet } from '@services/wallet/hdWallet';
import { DeshChainClient } from '@services/blockchain/deshchainClient';
import { setWalletCreated, setAddress, setMnemonic } from '@store/slices/walletSlice';

interface WalletContextType {
  wallet: HDWallet | null;
  client: DeshChainClient | null;
  isWalletCreated: boolean;
  createWallet: () => Promise<string>;
  importWallet: (mnemonic: string) => Promise<void>;
  getBalance: () => Promise<any>;
}

const WalletContext = createContext<WalletContextType | undefined>(undefined);

export const useWallet = () => {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error('useWallet must be used within a WalletProvider');
  }
  return context;
};

export const WalletProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const dispatch = useAppDispatch();
  const { isWalletCreated } = useAppSelector(state => state.wallet);
  const [wallet, setWallet] = useState<HDWallet | null>(null);
  const [client, setClient] = useState<DeshChainClient | null>(null);

  useEffect(() => {
    initializeWallet();
  }, []);

  const initializeWallet = async () => {
    try {
      const hdWallet = new HDWallet();
      const hasSeed = await hdWallet.hasSeed();
      
      if (hasSeed) {
        await hdWallet.unlock(''); // In production, this would require the PIN
        const address = await hdWallet.getAddress();
        const deshClient = new DeshChainClient();
        await deshClient.connect();
        
        setWallet(hdWallet);
        setClient(deshClient);
        dispatch(setWalletCreated(true));
        dispatch(setAddress(address));
      }
    } catch (error) {
      console.error('Failed to initialize wallet:', error);
    }
  };

  const createWallet = async (): Promise<string> => {
    try {
      const hdWallet = new HDWallet();
      const mnemonic = HDWallet.generateMnemonic();
      await hdWallet.createWallet(mnemonic);
      
      const address = await hdWallet.getAddress();
      const deshClient = new DeshChainClient();
      await deshClient.connect();
      
      setWallet(hdWallet);
      setClient(deshClient);
      dispatch(setWalletCreated(true));
      dispatch(setAddress(address));
      dispatch(setMnemonic(mnemonic));
      
      return mnemonic;
    } catch (error) {
      console.error('Failed to create wallet:', error);
      throw error;
    }
  };

  const importWallet = async (mnemonic: string): Promise<void> => {
    try {
      const hdWallet = new HDWallet();
      await hdWallet.createWallet(mnemonic);
      
      const address = await hdWallet.getAddress();
      const deshClient = new DeshChainClient();
      await deshClient.connect();
      
      setWallet(hdWallet);
      setClient(deshClient);
      dispatch(setWalletCreated(true));
      dispatch(setAddress(address));
    } catch (error) {
      console.error('Failed to import wallet:', error);
      throw error;
    }
  };

  const getBalance = async () => {
    if (!wallet || !client) {
      throw new Error('Wallet not initialized');
    }
    
    const address = await wallet.getAddress();
    return await client.getBalance(address);
  };

  return (
    <WalletContext.Provider
      value={{
        wallet,
        client,
        isWalletCreated,
        createWallet,
        importWallet,
        getBalance,
      }}
    >
      {children}
    </WalletContext.Provider>
  );
};