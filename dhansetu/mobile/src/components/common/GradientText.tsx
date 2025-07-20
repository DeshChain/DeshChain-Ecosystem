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
import { Text, TextProps, View, StyleSheet } from 'react-native';
import MaskedView from '@react-native-masked-view/masked-view';
import LinearGradient from 'react-native-linear-gradient';
import { COLORS } from '@constants/theme';

interface GradientTextProps extends TextProps {
  colors?: string[];
  locations?: number[];
  start?: { x: number; y: number };
  end?: { x: number; y: number };
  children: React.ReactNode;
}

export const GradientText: React.FC<GradientTextProps> = ({
  colors = [COLORS.saffron, COLORS.white, COLORS.green],
  locations = [0, 0.5, 1],
  start = { x: 0, y: 0 },
  end = { x: 1, y: 0 },
  style,
  children,
  ...props
}) => {
  return (
    <MaskedView
      maskElement={
        <Text style={[styles.maskText, style]} {...props}>
          {children}
        </Text>
      }
    >
      <LinearGradient
        colors={colors}
        locations={locations}
        start={start}
        end={end}
        style={StyleSheet.absoluteFillObject}
      />
    </MaskedView>
  );
};

const styles = StyleSheet.create({
  maskText: {
    backgroundColor: 'transparent',
    fontSize: 24,
    fontWeight: 'bold',
  },
});