import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { CompositeNavigationProp } from '@react-navigation/native';
import { BottomTabNavigationProp } from '@react-navigation/bottom-tabs';

// Import navigation param lists
import type { 
  RootStackParamList, 
  MainTabParamList,
  HomeStackParamList,
  DexStackParamList,
  SikkebaazStackParamList,
  SurakshaStackParamList,
} from '../types/navigation';

// Type for root stack navigation
export type RootStackNavigationProp = StackNavigationProp<RootStackParamList>;

// Type for main tab navigation
export type MainTabNavigationProp = BottomTabNavigationProp<MainTabParamList>;

// Type for home stack navigation
export type HomeStackNavigationProp = CompositeNavigationProp<
  StackNavigationProp<HomeStackParamList>,
  CompositeNavigationProp<
    BottomTabNavigationProp<MainTabParamList>,
    StackNavigationProp<RootStackParamList>
  >
>;

// Type for dex stack navigation
export type DexStackNavigationProp = CompositeNavigationProp<
  StackNavigationProp<DexStackParamList>,
  CompositeNavigationProp<
    BottomTabNavigationProp<MainTabParamList>,
    StackNavigationProp<RootStackParamList>
  >
>;

// Type for sikkebaaz stack navigation
export type SikkebaazStackNavigationProp = CompositeNavigationProp<
  StackNavigationProp<SikkebaazStackParamList>,
  CompositeNavigationProp<
    BottomTabNavigationProp<MainTabParamList>,
    StackNavigationProp<RootStackParamList>
  >
>;

// Type for suraksha stack navigation
export type SurakshaStackNavigationProp = CompositeNavigationProp<
  StackNavigationProp<SurakshaStackParamList>,
  CompositeNavigationProp<
    BottomTabNavigationProp<MainTabParamList>,
    StackNavigationProp<RootStackParamList>
  >
>;

// Typed navigation hooks
export const useRootNavigation = () => useNavigation<RootStackNavigationProp>();
export const useHomeNavigation = () => useNavigation<HomeStackNavigationProp>();
export const useDexNavigation = () => useNavigation<DexStackNavigationProp>();
export const useSikkebaazNavigation = () => useNavigation<SikkebaazStackNavigationProp>();
export const useSurakshaNavigation = () => useNavigation<SurakshaStackNavigationProp>();