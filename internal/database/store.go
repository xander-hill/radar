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
func (s *Store) GetNearbyPlans(ctx context.Context, lat, lng float64, radius float64, category string) ([]models.Plan, error) {
	query := `
		SELECT id, title, category, ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng, 
		       base_score, check_ins, saves, created_at, updated_at
		FROM plans
		WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3)
	`

	args := []interface{}{lng, lat, radius}

	if category != "" {
		query += " AND category = $4"
		args = append(args, category)
	}
	rows, err := s.Conn.Query(ctx, query, args...)
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

// IncrementCheckIn bumps the count and updates the timestamp
func (s *Store) IncrementCheckIn(ctx context.Context, planID int) error {
	query := `
		UPDATE plans 
		SET check_ins = check_ins + 1, 
		    updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`
	_, err := s.Conn.Exec(ctx, query, planID)
	return err
}

// IncrementSave bumps the save count
func (s *Store) IncrementSave(ctx context.Context, planID int) error {
	query := `
		UPDATE plans 
		SET saves = saves + 1, 
		    updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`
	_, err := s.Conn.Exec(ctx, query, planID)
	return err
}
