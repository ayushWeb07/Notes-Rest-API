package handlers

import (
	"net/http"

	"github.com/ayushWeb07/Notes-Rest-API/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Note struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func CreateNote(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := Note{}

		// bind the req json body with the note struct
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create a new note as some fields were missing",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		newNote, err := repository.CreateNote(pool, data.Title, data.Description)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while creating a new note",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully created a new note",
			"note":    newNote,
		})
	}
}

func GetAllNotes(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// call the repository endpoint
		allNotes, err := repository.GetAllNotes(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while reading all the notes",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched all the notes",
			"notes":   allNotes,
		})
	}
}
