export type RootStackParamList = {
  Onboarding: undefined;
  CreateWallet: undefined;
  ImportWallet: undefined;
  PinSetup: { isNewWallet: boolean };
  Main: undefined;
};

export type MainTabParamList = {
  Home: undefined;
  Wallet: undefined;
  DEX: undefined;
  Sikkebaaz: undefined;
  Suraksha: undefined;
};

export type HomeStackParamList = {
  HomeScreen: undefined;
  Send: { coinType?: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN' };
  Receive: { coinType?: 'DESHCHAIN' | 'ETHEREUM' | 'BITCOIN' };
  TransactionDetails: { txHash: string };
  Profile: undefined;
  Settings: undefined;
};

export type DexStackParamList = {
  DexScreen: undefined;
  CreateMoneyOrder: undefined;
  MoneyOrderDetails: { orderId: string };
};

export type SikkebaazStackParamList = {
  SikkebaazScreen: undefined;
  CreateLaunch: undefined;
  LaunchDetails: { launchId: string };
};

export type SurakshaStackParamList = {
  SurakshaScreen: undefined;
  EnrollSuraksha: undefined;
  SurakshaDetails: { poolId: string };
};

export type SettingsStackParamList = {
  SettingsScreen: undefined;
  SecurityScreen: undefined;
  LanguageScreen: undefined;
  AboutScreen: undefined;
  DhanPataSetupScreen: undefined;
};