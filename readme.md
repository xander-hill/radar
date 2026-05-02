# RADAR

A real-time social discovery engine that tracks the "momentum" of local spots. Built with Go and PostGIS, it helps users find what's trending nearby in categories like coffee, nightlife, art, and food.

## Core Features

- **Momentum Scoring:** Real-time ranking using a time-decay algorithm. Spots gain "hype" from check-ins and saves, which naturally cools off over time.
- **Geospatial Discovery:** Advanced PostGIS queries to find spots within a specific radius (default 10km) or city-wide trending spots (50km).
- **Soft Deletes:** Professional data management using `deleted_at` timestamps, allowing for instant data recovery and audit trails.
- **Security & Safety:**
  - API Key authentication (`X-API-KEY`) for all write operations.
  - Category validation to ensure data integrity.
  - Rate limiting to prevent spamming check-ins or saves.
- **GeoJSON Support:** Native GeoJSON output for seamless integration with Mapbox, Leaflet, or Google Maps.

## Tech Stack

- **Backend:** Go (Golang)
- **Database:** PostgreSQL + PostGIS (Dockerized)
- **Framework:** Gin Gonic
- **Documentation:** Swagger (swag)

## API Reference

### Discovery (Public)

- `GET /radar`: Get nearby plans sorted by momentum.
- `GET /radar/trending`: Get the top 3 trending spots in the region.
- `GET /radar/geojson`: Get data in GeoJSON format for map rendering.

### Actions (Requires X-API-KEY)

- `POST /plans`: Create a new spot.
- `POST /plans/:id/checkin`: Register a visit.
- `POST /plans/:id/save`: Mark interest.
- `DELETE /plans/:id`: Hide a spot from the radar.
- `POST /plans/:id/restore`: Bring a hidden spot back.

## 🛠 Setup

1. **Environment:** Create a `.env` file with `DB_URL`, `PORT`, and `INTERNAL_API_KEY`.
2. **Database:**
   docker-compose up -d
3. **Run Server:**
   make dev
