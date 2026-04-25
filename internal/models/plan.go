package models

import "time"

type Plan struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
