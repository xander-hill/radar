package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	connURL := os.Getenv("DB_URL")
	if connURL == "" {
		log.Fatal("DB_URL is not set in environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback default
	}

	// 1. Connect to DB
	conn, err := database.ConnectDB(connURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer conn.Close(context.Background())

	store := database.NewStore(conn)

	// 2. Setup Web Router
	r := gin.New() // Use New() instead of Default() for total control
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// This custom format shows the Method, Path, Status, and Query Params
		return fmt.Sprintf("[RADAR] %s | %d | %s | %s | %s\n",
			param.Method,
			param.StatusCode,
			param.ClientIP,
			param.Path,
			param.Request.URL.RawQuery, // This shows the lat/lng being searched!
		)
	}))
	r.Use(gin.Recovery()) // Prevents crashes from taking down the whole server

	// 3. Documentation (Public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 4. Read-Only Discovery (Public - No rate limit)
	r.GET("/radar", store.GetRadarHandler)
	r.GET("/radar/geojson", store.GetRadarGeoJSONHandler)
	r.GET("/radar/trending", store.GetTrendingHandler)

	// 5. Interaction Routes (Protected - Rate Limited)
	pulse := r.Group("/plans")
	pulse.Use(api.APIKeyAuth())
	pulse.Use(api.RateLimiter())
	{
		pulse.POST("", store.CreatePlanHandler)
		pulse.POST("/:id/checkin", store.PostCheckInHandler)
		pulse.POST("/:id/save", store.PostSaveHandler)
		pulse.DELETE("/:id", store.DeletePlanHandler)
		pulse.POST("/:id/restore", store.RestorePlanHandler)
	}

	// 6. Start the Engine
	log.Println("Radar API is live on http://localhost:8080")
	r.Run(":8080")
}
