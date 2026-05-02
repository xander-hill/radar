package database

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xanderhill/radar/internal/logic"
	"github.com/xanderhill/radar/internal/models"
)

// GetRadarHandler godoc
// @Summary Get nearby trending plans
// @Description Fetches plans within a radius and sorts by momentum score
// @Tags plans
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Param category query string false "Filter by category (e.g. coffee, art)"
// @Success 200 {array} models.Plan "Includes distance_meters and momentum_score"
// @Router /radar [get]
func (s *Store) GetRadarHandler(c *gin.Context) {
	// 1. Get query params (e.g., /radar?lat=34.05&lng=-118.24)
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	category := c.Query("category")
	radius := 10000.0 // 10km default

	var plans []models.Plan
	var err error

	// 2. Fetch from DB
	plans, err = s.GetNearbyPlans(c.Request.Context(), lat, lng, radius, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate scores for each plan before returning
	for i := range plans {
		plans[i].MomentumScore = logic.CalculateScore(plans[i])
		plans[i].Status = logic.GetStatus(plans[i].MomentumScore) // Dynamic assignment
	}

	// 3. Sort by Momentum
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].MomentumScore == plans[j].MomentumScore {
			return plans[i].Distance < plans[j].Distance
		}
		return plans[i].MomentumScore > plans[j].MomentumScore
	})

	// 4. Return as JSON
	c.JSON(http.StatusOK, plans)
}

// PostCheckInHandler godoc
// @Summary      Check in to a plan
// @Description  Increments the check-in count for a specific plan
// @Tags         plans
// @Param        id   path      int  true  "Plan ID"
// @Success      200  {object}  map[string]string
// @Router       /plans/{id}/checkin [post]
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

// PostSaveHandler godoc
// @Summary      Save a plan
// @Description  Increments the save count for a specific plan
// @Tags         plans
// @Param        id   path      int  true  "Plan ID"
// @Success      200  {object}  map[string]string
// @Router       /plans/{id}/save [post]
func (s *Store) PostSaveHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.IncrementSave(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan saved! Interest signal tracked."})
}

// CreatePlanHandler adds a new activity to the radar
// @Summary Create a new plan
// @ Description Adds a new location to the databse with geospatial coordinates
// @Tags plans
// @Accept json
// @Produce json
// @Param plan body models.Plan true "Plan data"
// @Success 201 {object} models.Plan
// @Router /plans [post]
func (s *Store) CreatePlanHandler(c *gin.Context) {
	var p models.Plan
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	p.BaseScore = 1.0

	if p.Category != "" && !models.IsValidCategory(p.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category. Allowed: coffee, nightlife, art, food",
		})
		return
	}

	// 1. Coordinate Validation
	if p.Lat < -90 || p.Lat > 90 || p.Lng < -180 || p.Lng > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coordinates out of bounds"})
		return
	}

	// 2. Sanitize Inputs (Avoid XSS or weird strings)
	if len(p.Title) < 3 || len(p.Title) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be between 3 and 100 chars"})
		return
	}

	// 3. Security: Cap the BaseScore
	if p.BaseScore > 50 {
		p.BaseScore = 50 // Hard cap for public users
	}

	id, err := s.CreatePlan(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}

	p.ID = id
	c.JSON(http.StatusCreated, p)
}

// GetRadarGeoJSONHandler returns nearby plans in GeoJSON format
// @Summary Get nearby plans as GeoJSON
// @Tags plans
// @Produce json
// @Router /radar/geojson [get]
func (s *Store) GetRadarGeoJSONHandler(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius := 10000.0

	plans, err := s.GetNearbyPlans(c.Request.Context(), lat, lng, radius, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mapbox/Leaflet compatible structure
	features := make([]map[string]interface{}, len(plans))
	for i, p := range plans {
		// Recalculate momentum for the properties
		momentum := logic.CalculateScore(p)
		weight := momentum / 20.0
		if weight > 1.0 {
			weight = 1.0
		}

		features[i] = map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{p.Lng, p.Lat}, // GeoJSON is [Lng, Lat]
			},
			"properties": map[string]interface{}{
				"id":       p.ID,
				"title":    p.Title,
				"momentum": momentum,
				"category": p.Category,
				"weight":   weight,
				"status":   logic.GetStatus(momentum),
				"distance": p.Distance,
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"type":     "FeatureCollection",
		"features": features,
	})
}

func (s *Store) DeletePlanHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := s.SoftDeletePlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Store) GetTrendingHandler(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	// Search a much wider area for "Trending" (50km)
	plans, err := s.GetNearbyPlans(c.Request.Context(), lat, lng, 50000.0, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var trending []models.Plan
	for i := range plans {
		plans[i].MomentumScore = logic.CalculateScore(plans[i])
		plans[i].Status = logic.GetStatus(plans[i].MomentumScore)

		// Only include "Trending" items
		if plans[i].Status == "trending" {
			trending = append(trending, plans[i])
		}
	}

	// Sort by Momentum Score descending
	sort.Slice(trending, func(i, j int) bool {
		return trending[i].MomentumScore > trending[j].MomentumScore
	})

	// Return top 3
	if len(trending) > 3 {
		trending = trending[:3]
	}

	c.JSON(http.StatusOK, trending)
}
