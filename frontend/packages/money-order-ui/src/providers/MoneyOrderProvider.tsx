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

import React, { createContext, ReactNode } from 'react';
import { MoneyOrderConfig } from '../types';
import { DEFAULT_THEME_CONFIG } from '../constants';

interface MoneyOrderContextType {
  config: MoneyOrderConfig;
  client: any; // DeshChain SDK client
}

export const MoneyOrderContext = createContext<MoneyOrderContextType | null>(null);

interface MoneyOrderProviderProps {
  children: ReactNode;
  config: Partial<MoneyOrderConfig>;
  client?: any;
}

const defaultConfig: MoneyOrderConfig = {
  apiUrl: 'https://api.deshchain.org/v1/moneyorder',
  chainId: 'deshchain-1',
  defaultLanguage: 'en',
  enableCulturalFeatures: true,
  enableFestivalThemes: true,
  enablePatriotismRewards: true,
  maxSlippage: 0.05,
  defaultPoolType: 'fixed_rate',
  autoSelectPool: true,
  theme: DEFAULT_THEME_CONFIG
};

export const MoneyOrderProvider: React.FC<MoneyOrderProviderProps> = ({
  children,
  config: userConfig,
  client
}) => {
  const config = { ...defaultConfig, ...userConfig };

  const contextValue: MoneyOrderContextType = {
    config,
    client
  };

  return (
    <MoneyOrderContext.Provider value={contextValue}>
      {children}
    </MoneyOrderContext.Provider>
  );
};