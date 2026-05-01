package database

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xanderhill/radar/internal/logic"
)

// GetRadarHandler handles the GET /radar request
func (s *Store) GetRadarHandler(c *gin.Context) {
	// 1. Get query params (e.g., /radar?lat=34.05&lng=-118.24)
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius := 10000.0 // 10km default

	// 2. Fetch from DB
	plans, err := s.GetNearbyPlans(c.Request.Context(), lat, lng, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate scores for each plan before returning
	for i := range plans {
		plans[i].MomentumScore = logic.CalculateScore(plans[i])
		plans[i].Status = logic.GetState(plans[i].MomentumScore)
	}

	// 3. Sort by Momentum
	sort.Slice(plans, func(i, j int) bool {
		return logic.CalculateScore(plans[i]) > logic.CalculateScore(plans[j])
	})

	// 4. Return as JSON
	c.JSON(http.StatusOK, plans)
}

// PostCheckInHandler handles the "I am here" signal
func (s *Store) PostCheckInHandler(c *gin.Context) {
	// 1. Get the ID from the URL (e.g., /plans/1/checkin)
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// 2. Update the DB
	err := s.IncrementCheckIn(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check-in successful! Momentum increased."})
}
