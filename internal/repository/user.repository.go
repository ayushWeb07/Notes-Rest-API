package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ayushWeb07/Notes-Rest-API/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUser(pool *pgxpool.Pool, email string, password string) (*models.User, error) {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// insert query
	newUser := models.User{}

	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at, updated_at"

	err := pool.QueryRow(ctx, query, email, password).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while registering user into database:", err)
		return nil, err
	}

	return &newUser, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// get single row query
	query := "SELECT id, email, created_at, updated_at FROM users WHERE email=$1"
	user := models.User{}

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while fetching a user from database:", err)
		return nil, err
	}

	return &user, nil
}

func GetUserById(pool *pgxpool.Pool, id int) (*models.User, error) {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// get single row query
	query := "SELECT id, email, created_at, updated_at FROM users WHERE id=$1"
	user := models.User{}

	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while fetching a user from database:", err)
		return nil, err
	}

	return &user, nil
}
