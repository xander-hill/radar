package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xanderhill/radar/internal/database"
)

func main() {
	connURL := "postgres://user:password@127.0.0.1:5432/radar_db"

	// 1. Connect to DB
	conn, err := database.ConnectDB(connURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer conn.Close(context.Background())

	store := database.NewStore(conn)

	// 2. Setup Web Router
	r := gin.Default()

	// 3. Define the Route
	r.GET("/radar", store.GetRadarHandler)

	// 4. Start the Engine
	log.Println("Radar API is live on http://localhost:8080")
	r.Run(":8080")
}
