package handlers

import (
	"net/http"
	"strconv"

	"github.com/ayushWeb07/Notes-Rest-API/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNoteInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateNoteInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func CreateNote(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from auth middleware
		userIdInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to create a new note as you are not authenticated",
			})

			return
		}

		userId := userIdInterface.(string)

		data := CreateNoteInput{}

		// bind the req json body with the note struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create a new note as some fields were missing",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		newNote, err := repository.CreateNote(pool, data.Title, data.Description, userId)

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
		// get user id from auth middleware
		userIdInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to fetch all the notes as you are not authenticated",
			})

			return
		}

		userId := userIdInterface.(string)

		// call the repository endpoint
		allNotes, err := repository.GetAllNotes(pool, userId)

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

func GetNoteById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from auth middleware
		userIdInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to fetch the note as you are not authenticated",
			})

			return
		}

		userId := userIdInterface.(string)

		// fetch the id from params
		idParam := c.Param("id")
		idNumParam, err := strconv.Atoi(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid id passed in the params",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		note, err := repository.GetNoteById(pool, idNumParam, userId)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Note either deleted or does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while fetching the note",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched the note",
			"note":    note,
		})
	}
}

func UpdateNote(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from auth middleware
		userIdInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to update the note as you are not authenticated",
			})

			return
		}

		userId := userIdInterface.(string)

		data := UpdateNoteInput{}

		// bind the req json body with the note struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to update the note",
				"error":   err.Error(),
			})

			return
		}

		// fetch the id from params
		idParam := c.Param("id")
		idNumParam, err := strconv.Atoi(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid id passed in the params",
				"error":   err.Error(),
			})

			return
		}

		// fetch the existing note
		existingNote, err := repository.GetNoteById(pool, idNumParam, userId)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Note either deleted or does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while updating the note",
				"error":   err.Error(),
			})

			return
		}

		title := existingNote.Title
		description := existingNote.Description

		// check if title has been modified and it's not empty
		if data.Title != nil && *data.Title != "" {
			title = *data.Title
		}

		// check if description has been modified and it's not empty
		if data.Description != nil && *data.Description != "" {
			description = *data.Description
		}

		// check if none of the fields have been modified
		if data.Title == nil && data.Description == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "At least specify one field to modify the note",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		updatedNote, err := repository.UpdateNote(pool, idNumParam, userId, title, description)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while updating the note",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully updated the note",
			"note":    updatedNote,
		})
	}
}

func DeleteNote(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from auth middleware
		userIdInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to delete the note as you are not authenticated",
			})

			return
		}

		userId := userIdInterface.(string)

		// fetch the id from params
		idParam := c.Param("id")
		idNumParam, err := strconv.Atoi(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid id passed in the params",
				"error":   err.Error(),
			})

			return
		}

		// fetch the existing note
		_, err = repository.GetNoteById(pool, idNumParam, userId)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Such note does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while deleting the note",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		err = repository.DeleteNote(pool, idNumParam, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while deleting the note",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted the note",
		})
	}
}
