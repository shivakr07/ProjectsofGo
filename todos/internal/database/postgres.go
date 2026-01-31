package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error
	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Printf("Unable to parse the DATABASE_URL: %v", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create the connection pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Successfully connected to PostgresSQL database")
	return pool, nil
}
