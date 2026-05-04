// src/app/index.tsx
import React, { useEffect } from 'react';
import { StyleSheet, View, Text, TouchableOpacity } from 'react-native';
import MapView, { Marker, Callout } from 'react-native-maps';
import * as Location from 'expo-location';
import { useRadar } from '../../hooks/useRadar';

const CATEGORY_COLORS: any = {
  coffee: '#6F4E37',    // Brown
  food: '#FF5733',      // Orange-Red
  nightlife: '#8E44AD', // Purple
  art: '#2ECC71',       // Green
  general: '#3498DB',   // Blue
};

export default function RadarScreen() {
  const { plans, refreshRadar } = useRadar();
  const [selectedPlan, setSelectedPlan] = React.useState<any>(null);
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
        zoomEnabled={true}           
        scrollEnabled={true}         
        rotateEnabled={true}
        initialRegion={{
          latitude: 44.974,
          longitude: -93.235,
          latitudeDelta: 0.02, 
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
            pinColor={
              plan.status === 'trending' ? '#FFD700' : 
              plan.status === 'steady'   ? '#FFA500' :
              (CATEGORY_COLORS[plan.category] || 'red')
            }
            onPress={() => setSelectedPlan(plan)}
          />
        ))}
      </MapView>
      {selectedPlan && (
        <View style={styles.bottomSheet}>
          <View style={styles.sheetHeader}>
            <View>
              <Text style={styles.categoryBadge}>{selectedPlan.category?.toUpperCase()}</Text>
              <Text style={styles.sheetTitle}>{selectedPlan.title}</Text>
            </View>
            <TouchableOpacity onPress={() => setSelectedPlan(null)}>
              <Text style={{ color: '#999', fontSize: 18 }}>✕</Text>
            </TouchableOpacity>
          </View>
          
          <Text style={styles.sheetDescription}>{selectedPlan.description}</Text>
          
          <TouchableOpacity style={styles.checkInButton}>
            <Text style={styles.buttonText}>Check In Here</Text>
          </TouchableOpacity>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  map: { width: '100%', height: '100%' },
  bottomSheet: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'white',
    padding: 20,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: -2 },
    shadowOpacity: 0.1,
    shadowRadius: 10,
    elevation: 5,
  },
  sheetHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 10,
  },
  sheetTitle: { fontSize: 22, fontWeight: 'bold' },
  sheetDescription: { fontSize: 16, color: '#666', marginBottom: 20 },
  checkInButton: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 12,
    alignItems: 'center',
  },
  buttonText: { color: 'white', fontWeight: '600', fontSize: 16 },
  categoryBadge: {
    fontSize: 10,
    fontWeight: 'bold',
    color: '#007AFF',
    letterSpacing: 1,
    marginBottom: 2,
  },
});