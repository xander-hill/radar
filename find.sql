SELECT title, lat, lng, 
       (6371 * acos(cos(radians(44.974)) * cos(radians(lat)) * cos(radians(lng) - radians(-93.235)) + sin(radians(44.974)) * sin(radians(lat)))) AS distance 
FROM plans 
ORDER BY distance ASC;