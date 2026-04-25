package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Connection string (Matches our docker-compose.yml)
	url := "postgres://user:password@localhost:5432/radar_db"

	// Try to connect
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println("Successfully connected to the Radar Database! 🚀")
}
