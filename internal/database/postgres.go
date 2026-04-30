package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}
	return conn, nil
}
