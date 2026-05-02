package models

import "time"

// Plan represents a real-world activity with "momentum"
type Plan struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`

	// Analytics signals for the Momentum Engine
	BaseScore float64 `json:"base_score"` // Organizer's initial "weight"
	CheckIns  int     `json:"check_ins"`  // Real-time presence
	Saves     int     `json:"saves"`      // Intent signal

	// Enriched Fields
	MomentumScore float64 `json:"momentum_score"`
	Status        string  `json:"status"`

	LastCheckinAt time.Time `json:"last_checkin_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
