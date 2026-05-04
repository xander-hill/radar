# RADAR

A real-time social discovery engine that tracks the "momentum" of local spots. Built with a Go/PostGIS backend and a React Native frontend, it helps users find what's trending nearby.

## Core Features

- **Momentum Scoring:** Real-time ranking using a time-decay algorithm ($Energy / (Hours + 2)^{1.8}$). Spots gain "hype" from check-ins, which naturally cools off over time.
- **Reactive Mobile Map:** An interactive Expo-based map that updates marker colors (Gold/Orange/Blue) dynamically based on live momentum scores.
- **Smart Bottom Sheet:** A modern slide-up interaction panel for viewing spot details and performing instant check-ins.
- **Geospatial Discovery:** Advanced PostGIS queries to find spots within a specific radius using `ST_DWithin` and `ST_Distance`.
- **Soft Deletes:** Professional data management using `deleted_at` timestamps for instant recovery.

## Tech Stack

- **Backend:** Go (Golang) + Gin Gonic
- **Database:** PostgreSQL + PostGIS (Dockerized)
- **Frontend:** React Native + Expo + TypeScript
- **Maps:** React Native Maps + Expo Location
- **Styling:** Pulse/Glow marker logic with category-based color mapping.

## API Reference

### Discovery (Public)

- `GET /radar`: Get nearby plans sorted by momentum.
- `GET /radar/trending`: Get the top 3 trending spots in the region.

### Actions (Requires X-API-KEY)

- `POST /plans`: Create a new spot.
- `POST /plans/:id/checkin`: Register a visit (immediately spikes momentum).
- `DELETE /plans/:id`: Hide a spot from the radar.

## 📱 Mobile Setup

1. **Install Dependencies:**
   npm install
2. **Configure IP:**
   Update src/constants/Config.ts with your machine's local IP address to allow the physical device/emulator to communicate with the Dockerized backend.
3. **Start Expo:**
   npx expo start

## Backend Setup

1. **Environment:**
   Create a .env file with DB_URL, PORT, and INTERNAL_API_KEY
2. **Database:**
   docker-compose up -d
3. **Run Server:**
   make dev
