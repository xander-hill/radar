package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/xanderhill/radar/internal/database"
	"github.com/xanderhill/radar/internal/logic"
)

func main() {
	ctx := context.Background()
	connURL := "postgres://user:password@127.0.0.1:5432/radar_db"

	// 1. Connect
	conn, err := database.ConnectDB(connURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	// 2. Init Store
	store := database.NewStore(conn)

	// 3. Perform a Radar Sweep (User at 34.0, -118.0)
	userLat, userLng := 34.0522, -118.2437
	radius := 10000.0 // 10km

	plans, err := store.GetNearbyPlans(ctx, userLat, userLng, radius)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// 4. Sort by Momentum Score
	sort.Slice(plans, func(i, j int) bool {
		return logic.CalculateScore(plans[i]) > logic.CalculateScore(plans[j])
	})

	fmt.Printf("Found %d plans in your area:\n", len(plans))
	for _, p := range plans {
		fmt.Printf("- %s (Score: %.2f)\n", p.Title, logic.CalculateScore(p))
	}
}
