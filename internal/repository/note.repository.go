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

func GetAllNotes(pool *pgxpool.Pool) ([]models.Note, error) {

	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// get all query
	query := "SELECT id, title, description, created_at, updated_at FROM notes ORDER BY updated_at DESC;"

	rows, err := pool.Query(ctx, query)
	defer rows.Close()

	if err != nil {
		fmt.Println("Something went wrong while getting all notes from database:", err)
		return nil, err
	}

	// scan each row and insert into slice
	var allNotes []models.Note

	// scanning row
	for rows.Next() {
		note := models.Note{}

		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Description,
			&note.CreatedAt,
			&note.UpdatedAt,
		)

		// error while scanning an individual row
		if err != nil {
			return nil, err
		}

		allNotes = append(allNotes, note)
	}

	// error which scanning all the rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allNotes, nil
}
