package logic

import (
	"math"
	"time"

	"github.com/xanderhill/radar/internal/models"
)

// CalculateScore uses a time-decay formula to determine momentum
func CalculateScore(p models.Plan) float64 {
	// 1. Calculate raw "Energy" (Volume)
	// We weigh Saves slightly higher than Check-ins (intent vs presence)
	energy := p.BaseScore + (float64(p.CheckIns) * 2.0) + (float64(p.Saves) * 3.0)

	// 2. Calculate "Freshness"
	// time.Since gives us the duration between then and now
	hoursSinceActivity := time.Since(p.LastCheckinAt).Hours()

	// 3. Apply Gravity
	// Formula: Energy / (Hours + 2)^1.8
	// The "+2" ensures that brand new activity doesn't create an infinite score.
	gravity := 1.8
	momentum := energy / math.Pow(hoursSinceActivity+2, gravity)

	return momentum
}

// GetState translates a numerical score into a human-readable "vibe"
func GetStatus(score float64) string {
	if score > 15.0 {
		return "trending"
	}
	if score > 5.0 {
		return "steady"
	}
	if score > 0.5 {
		return "quiet"
	}
	return "cold"
}
