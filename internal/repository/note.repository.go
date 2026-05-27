package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ayushWeb07/Notes-Rest-API/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNote(pool *pgxpool.Pool, title string, description string) (*models.Note, error) {

	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// insert query
	newNote := models.Note{}

	query := "INSERT INTO notes (title, description) VALUES ($1, $2) RETURNING id, title, description, created_at, updated_at"

	err := pool.QueryRow(ctx, query, title, description).Scan(
		&newNote.ID,
		&newNote.Title,
		&newNote.Description,
		&newNote.CreatedAt,
		&newNote.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while inserting note into database:", err)
		return nil, err
	}

	return &newNote, nil
}
