package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xanderhill/radar/internal/models"
)

type Store struct {
	Conn *pgx.Conn
}

func NewStore(conn *pgx.Conn) *Store {
	return &Store{Conn: conn}
}

// GetNearbyPlans fetches plans within a radius (in meters)
func (s *Store) GetNearbyPlans(ctx context.Context, userLat, userLng float64, radius float64) ([]models.Plan, error) {
	query := `
		SELECT id, title, ST_Y(location::geometry), ST_X(location::geometry), base_score, saves, check_ins, created_at
		FROM plans
		WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326), $3)
	`
	rows, err := s.Conn.Query(ctx, query, userLng, userLat, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var p models.Plan
		err := rows.Scan(&p.ID, &p.Title, &p.Lat, &p.Lng, &p.BaseScore, &p.Saves, &p.CheckIns, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}
