package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectWithDatabase(databaseConnectionUri string) (*pgxpool.Pool, error) {
	// parse config and connect to database
	config, err := pgxpool.ParseConfig(databaseConnectionUri)

	if err != nil {
		fmt.Println("Something went wrong while parsing database config:", err)
		return nil, err
	}

	// create a connection pool
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		fmt.Println("Something went wrong while creating a connection pool:", err)
		return nil, err
	}

	// ping the database
	err = pool.Ping(ctx)

	if err != nil {
		fmt.Println("Failed to ping the database:", err)
		pool.Close()
		return nil, err
	}

	fmt.Println("Successfully connected to database")
	return pool, nil
}
