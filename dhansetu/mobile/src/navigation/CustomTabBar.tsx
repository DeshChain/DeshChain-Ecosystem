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

import React from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Platform,
} from 'react-native';
import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import LinearGradient from 'react-native-linear-gradient';
import Animated, {
  useAnimatedStyle,
  withSpring,
  useSharedValue,
} from 'react-native-reanimated';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { COLORS, theme } from '@constants/theme';

interface TabIconProps {
  name: string;
  icon: string;
  color: string;
  focused: boolean;
}

const TabIcon: React.FC<TabIconProps> = ({ name, icon, color, focused }) => {
  const scale = useSharedValue(focused ? 1 : 0.9);
  
  React.useEffect(() => {
    scale.value = withSpring(focused ? 1 : 0.9);
  }, [focused]);

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
  }));

  return (
    <Animated.View style={[styles.tabIconContainer, animatedStyle]}>
      <Icon name={icon} size={24} color={color} />
      {focused && (
        <Text style={[styles.tabLabel, { color }]}>
          {name}
        </Text>
      )}
    </Animated.View>
  );
};

export const CustomTabBar: React.FC<BottomTabBarProps> = ({
  state,
  descriptors,
  navigation,
}) => {
  const insets = useSafeAreaInsets();

  const getTabIcon = (routeName: string): string => {
    switch (routeName) {
      case 'Home':
        return 'home';
      case 'Wallet':
        return 'wallet';
      case 'DEX':
        return 'swap-horizontal';
      case 'Sikkebaaz':
        return 'rocket-launch';
      case 'Suraksha':
        return 'shield-check';
      default:
        return 'circle';
    }
  };

  const getTabLabel = (routeName: string): string => {
    switch (routeName) {
      case 'DEX':
        return 'Money Order';
      case 'Sikkebaaz':
        return 'Launchpad';
      case 'Suraksha':
        return 'Pension';
      default:
        return routeName;
    }
  };

  return (
    <View style={[styles.container, { paddingBottom: insets.bottom }]}>
      <LinearGradient
        colors={[COLORS.white, COLORS.gray100]}
        style={styles.gradient}
        start={{ x: 0, y: 0 }}
        end={{ x: 0, y: 1 }}
      >
        <View style={styles.tabBar}>
          {state.routes.map((route, index) => {
            const { options } = descriptors[route.key];
            const isFocused = state.index === index;

            const onPress = () => {
              const event = navigation.emit({
                type: 'tabPress',
                target: route.key,
                canPreventDefault: true,
              });

              if (!isFocused && !event.defaultPrevented) {
                navigation.navigate(route.name);
              }
            };

            const onLongPress = () => {
              navigation.emit({
                type: 'tabLongPress',
                target: route.key,
              });
            };

            return (
              <TouchableOpacity
                key={route.key}
                accessibilityRole="button"
                accessibilityState={isFocused ? { selected: true } : {}}
                accessibilityLabel={options.tabBarAccessibilityLabel}
                testID={options.tabBarTestID}
                onPress={onPress}
                onLongPress={onLongPress}
                style={styles.tab}
                activeOpacity={0.8}
              >
                <TabIcon
                  name={getTabLabel(route.name)}
                  icon={getTabIcon(route.name)}
                  color={isFocused ? COLORS.saffron : COLORS.gray600}
                  focused={isFocused}
                />
                {isFocused && (
                  <View style={styles.indicator}>
                    <LinearGradient
                      colors={theme.gradients.indianFlag}
                      style={styles.indicatorGradient}
                      start={{ x: 0, y: 0 }}
                      end={{ x: 1, y: 0 }}
                    />
                  </View>
                )}
              </TouchableOpacity>
            );
          })}
        </View>
      </LinearGradient>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: 'transparent',
  },
  gradient: {
    ...theme.shadows.medium,
  },
  tabBar: {
    flexDirection: 'row',
    height: Platform.OS === 'ios' ? 60 : 56,
    paddingHorizontal: theme.spacing.sm,
  },
  tab: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    position: 'relative',
  },
  tabIconContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  tabLabel: {
    fontSize: 12,
    fontFamily: theme.fonts.medium.fontFamily,
    marginTop: 2,
  },
  indicator: {
    position: 'absolute',
    bottom: 0,
    left: '20%',
    right: '20%',
    height: 3,
    borderRadius: 1.5,
    overflow: 'hidden',
  },
  indicatorGradient: {
    flex: 1,
  },
});