package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ayushWeb07/Notes-Rest-API/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNote(pool *pgxpool.Pool, title string, description string, userId string) (*models.Note, error) {

	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// insert query
	newNote := models.Note{}

	query := "INSERT INTO notes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id, title, description, user_id, created_at, updated_at"

	err := pool.QueryRow(ctx, query, title, description, userId).Scan(
		&newNote.ID,
		&newNote.Title,
		&newNote.Description,
		&newNote.UserID,
		&newNote.CreatedAt,
		&newNote.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while inserting note into database:", err)
		return nil, err
	}

	return &newNote, nil
}

func GetAllNotes(pool *pgxpool.Pool, userId string) ([]models.Note, error) {

	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// get all query
	query := "SELECT id, title, description, user_id, created_at, updated_at FROM notes WHERE user_id=$1 ORDER BY updated_at DESC;"

	rows, err := pool.Query(ctx, query, userId)
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
			&note.UserID,
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

func GetNoteById(pool *pgxpool.Pool, id int, userId string) (*models.Note, error) {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// get single row query
	query := "SELECT id, title, description, user_id, created_at, updated_at FROM notes WHERE id=$1 AND user_id=$2"
	note := models.Note{}

	err := pool.QueryRow(ctx, query, id, userId).Scan(
		&note.ID,
		&note.Title,
		&note.Description,
		&note.UserID,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while fetching a note from database:", err)
		return nil, err
	}

	return &note, nil
}

func UpdateNote(pool *pgxpool.Pool, id int, userId string, title string, description string) (*models.Note, error) {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// update single row query
	query := "UPDATE notes SET title=$1, description=$2, updated_at= now() WHERE id=$3 AND user_id=$4 RETURNING id, title, description, user_id, created_at, updated_at"
	note := models.Note{}

	err := pool.QueryRow(ctx, query, title, description, id, userId).Scan(
		&note.ID,
		&note.Title,
		&note.Description,
		&note.UserID,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Something went wrong while updating the note in the database:", err)
		return nil, err
	}

	return &note, nil
}

func DeleteNote(pool *pgxpool.Pool, id int, userId string) error {
	// create a timed context to free up resources associated with it & avoid slow connections
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// release resource once all the operations are done
	defer cancelFunc()

	// delete row query
	query := "DELETE FROM notes WHERE id=$1 AND user_id=$2"

	cmdTag, err := pool.Exec(ctx, query, id, userId)

	if err != nil {
		fmt.Println("Something went wrong while deleting the note in the database:", err)
		return err
	}

	// no rows got deleted
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("Such note does not exist")
	}

	return nil
}
