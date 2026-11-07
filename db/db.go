package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectToDb() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	Pool = dbpool
	err = Pool.Ping(ctx)
	if err := Pool.Ping(ctx); err != nil {
		Pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL using pgx!")
	return nil

}
