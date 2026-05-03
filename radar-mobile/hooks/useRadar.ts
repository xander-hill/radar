// src/hooks/useRadar.ts
import { useState } from 'react';
import { API_BASE_URL, API_KEY } from '../constants/Config';

export const useRadar = () => {
  const [plans, setPlans] = useState([]);

  const refreshRadar = async (lat: number, lng: number) => {
    const url = `${API_BASE_URL}/radar?lat=${lat}&lng=${lng}`;
    try {
      console.log(`📡 Requesting: ${url}`);
      const response = await fetch(url, {
        headers: { 'X-API-KEY': API_KEY }
      });
      
      if (!response.ok) {
        throw new Error(`Server responded with ${response.status}`);
      }

      const data = await response.json();
      setPlans(data);
    } catch (err) {
      // Log the specific error so we know if it's a timeout, DNS, or security block
      console.error("❌ Radar Fetch Error:", err);
    }
  };

  return { plans, refreshRadar };
};