// src/app/index.tsx
import React, { useEffect } from 'react';
import { StyleSheet, View, Text, TouchableOpacity } from 'react-native';
import MapView, { Marker, Callout } from 'react-native-maps';
import * as Location from 'expo-location';
import { useRadar } from '../../hooks/useRadar';

export default function RadarScreen() {
  const { plans, refreshRadar } = useRadar();
  console.log("📍 Current Plans in State:", JSON.stringify(plans, null, 2));

  useEffect(() => {
    (async () => {
      let { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') return;

      let loc = await Location.getCurrentPositionAsync({});
      refreshRadar(loc.coords.latitude, loc.coords.longitude);
    })();
  }, []);

  return (
    <View style={styles.container}>
      <MapView 
        style={styles.map}
        showsUserLocation={true}
        zoomEnabled={true}            // <--- Make sure these are true
        scrollEnabled={true}          // <--- Make sure these are true
        rotateEnabled={true}
        initialRegion={{
          latitude: 44.974,
          longitude: -93.235,
          latitudeDelta: 0.02, // Smaller number = tighter zoom
          longitudeDelta: 0.02,
        }}
      >
        {(plans || []).map((plan: any) => (
          <Marker
            key={plan.id}
            coordinate={{ 
              latitude: Number(plan.lat), 
              longitude: Number(plan.lng) 
            }}
            pinColor={plan.status === 'trending' ? 'orange' : 'red'}
            onPress={() => console.log("Marker Pressed ID:", plan.id)}
          >
            {/* Remove title/description props from Marker and use this Callout */}
            <Callout tooltip={false}>
              <View style={{ width: 200, minHeight: 60, padding: 10, backgroundColor: 'white' }}>
                <Text style={{ fontWeight: 'bold', fontSize: 16 }}>{plan.title}</Text>
                <Text style={{ color: '#666' }}>{plan.description}</Text>
              </View>
            </Callout>
          </Marker>
        ))}
      </MapView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  map: { width: '100%', height: '100%' },
});