import React from 'react';
import { Tabs } from 'expo-router';
import { SymbolView } from 'expo-symbols'; // If you're on iOS

export default function TabLayout() {
  return (
    <Tabs screenOptions={{ tabBarActiveTintColor: '#007AFF' }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Radar',
          tabBarIcon: ({ color }) => <SymbolView name="antenna.radiowaves.left.and.right" size={24} tintColor={color} />,
        }}
      />
      <Tabs.Screen
        name="explore"
        options={{
          title: 'Trending',
          tabBarIcon: ({ color }) => <SymbolView name="flame.fill" size={24} tintColor={color} />,
        }}
      />
    </Tabs>
  );
}
