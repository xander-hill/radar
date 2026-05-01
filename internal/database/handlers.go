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

	// 3. Sort by Momentum
	sort.Slice(plans, func(i, j int) bool {
		return logic.CalculateScore(plans[i]) > logic.CalculateScore(plans[j])
	})

	// 4. Return as JSON
	c.JSON(http.StatusOK, plans)
}
