package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/xanderhill/radar/docs"
	"github.com/xanderhill/radar/internal/api"
	"github.com/xanderhill/radar/internal/database"
)

// @title Radar Momentum API
// @version 1.0
// @description This is a real-time social discovery engine.
// @host localhost:8080
// @BasePath /
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

	// 3. Documentation (Public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 4. Read-Only Discovery (Public - No rate limit)
	r.GET("/radar", store.GetRadarHandler)
	r.GET("/radar/geojson", store.GetRadarGeoJSONHandler)

	// 5. Interaction Routes (Protected - Rate Limited)
	// Using a group ensures the middleware only runs on these specific paths
	pulse := r.Group("/plans")
	pulse.Use(api.RateLimiter()) // Your new middleware from internal/api/middleware.go
	{
		pulse.POST("", store.CreatePlanHandler)
		pulse.POST("/:id/checkin", store.PostCheckInHandler)
		pulse.POST("/:id/save", store.PostSaveHandler)
	}

	// 6. Start the Engine
	log.Println("Radar API is live on http://localhost:8080")
	r.Run(":8080")
}
