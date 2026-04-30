package logic

import (
	"math"
	"time"

	"github.com/xanderhill/radar/internal/models"
)

// CalculateScore uses a time-decay formula to determine momentum
func CalculateScore(plan models.Plan) float64 {
	// 1. Calculate gravity (how fast interest fades)
	const gravity = 1.8

	// 2. Calculate time elapsed in hours since the plan was created
	hoursSinceCreation := time.Since(plan.CreatedAt).Hours()

	// 3. Aggregate signals
	points := (float64(plan.Saves) * 2.0) + (float64(plan.CheckIns) * 3.0) + plan.BaseScore

	// 4. The Newton's Law of Cooling inspired Formula:
	score := points / math.Pow(hoursSinceCreation+2, gravity)

	return score
}

// GetState translates a numerical score into a human-readable "vibe"
func GetState(score float64) string {
	if score > 10 {
		return "HOT 🔥"
	} else if score > 5 {
		return "BUILDING 📈"
	} else if score > 2 {
		return "EARLY 🌱"
	}
	return "COLD ❄️"
}
