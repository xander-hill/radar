# Radar 🛰️

Radar is a real-time social discovery engine built to digitize "word of mouth." Unlike static event directories, Radar uses a momentum-based scoring algorithm to surface what is happening _right now_ and what is gaining traction in the real world.

## Engineering Highlights

- **Momentum Scoring Engine:** Developed a time-decay algorithm in **Go** that ranks activities based on social signals (saves, check-ins) vs. temporal age.
- **Geospatial Intelligence:** Leverages **PostGIS** for high-performance coordinate indexing and radius-based discovery queries.
- **Real-Time Data Pipeline:** Architected to handle high-concurrency updates using a lightweight Go backend and connection pooling.
- **Containerized Infrastructure:** Fully orchestrated via **Docker** to ensure environment parity across development and production.

## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL with PostGIS extension
- **Infrastructure:** Docker & Docker Compose
- **Geospatial Logic:** Spherical coordinate geometry ($O(\log N)$ search)

## System Architecture

The system follows a **Clean Architecture** pattern, separating the core momentum logic from the infrastructure and delivery layers.

## 🚦 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### Installation

git clone [https://github.com/yourusername/radar.git](https://github.com/yourusername/radar.git) 2. Start the infrastructure:
docker-compose up -d 3. Run the engine:
go run main.go
