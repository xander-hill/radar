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

func (s *Store) GetNearbyPlans(ctx context.Context, lat, lng float64, radius float64, category string) ([]models.Plan, error) {
	// 1. Match SELECT to your Model (added last_checkin_at and description)
	query := `
        SELECT id, title, category, description, 
               ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng, 
               base_score, check_ins, saves, 
               last_checkin_at, created_at, updated_at
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
		// 2. Scan must match the SELECT order EXACTLY
		err := rows.Scan(
			&p.ID, &p.Title, &p.Category, &p.Description,
			&p.Lat, &p.Lng,
			&p.BaseScore, &p.CheckIns, &p.Saves,
			&p.LastCheckinAt, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

// IncrementCheckIn bumps the count and updates the timestamp
func (s *Store) IncrementCheckIn(ctx context.Context, id int) error {
	query := `
        UPDATE plans 
        SET check_ins = check_ins + 1, 
            last_checkin_at = NOW(), 
            updated_at = NOW() 
        WHERE id = $1
    `
	_, err := s.Conn.Exec(ctx, query, id)
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
