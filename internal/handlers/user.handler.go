package handlers

import (
	"net/http"
	"time"

	"github.com/ayushWeb07/Notes-Rest-API/internal/config"
	"github.com/ayushWeb07/Notes-Rest-API/internal/repository"
	"github.com/ayushWeb07/Notes-Rest-API/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Email    string `json:"email" binding:"required,min=6,max=50"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type GetUserByEmailInput struct {
	Email string `json:"email" binding:"required,min=6,max=50"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,min=6,max=50"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func RegisterUser(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := RegisterUserInput{}

		// bind the req json body with the user struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to register as a new user",
				"error":   err.Error(),
			})

			return
		}

		// check if user already exist
		existingUser, _ := repository.GetUserByEmail(pool, data.Email)

		if existingUser != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "User with such email already exists",
			})

			return
		}

		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while hashing the password",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		newUser, err := repository.RegisterUser(pool, data.Email, string(hashedPassword))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while registering as a new user",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully registered as a new user",
			"user":    newUser,
		})
	}
}

func GetUserByEmail(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := GetUserByEmailInput{}

		// bind the req json body with the user struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to get user by email",
				"error":   err.Error(),
			})

			return
		}

		// call the repository endpoint
		user, err := repository.GetUserByEmail(pool, data.Email)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User either deleted or does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while fetching the user by email",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched the user",
			"user":    user,
		})
	}
}

func GetUserById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// fetch the id from params
		idParam := c.Param("id")

		if isValidUUID := utils.IsValidUUID(idParam); isValidUUID == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid id passed in the params",
			})

			return
		}

		// call the repository endpoint
		user, err := repository.GetUserById(pool, idParam)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User either deleted or does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while fetching the user by id",
				"error":   err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched the user",
			"user":    user,
		})
	}
}

func LoginUser(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := LoginUserInput{}

		// bind the req json body with the user struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to login as some fields are missing",
				"error":   err.Error(),
			})

			return
		}

		// fetch the user by email
		existingUser, err := repository.GetUserByEmail(pool, data.Email)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Such user does not exist",
					"error":   err.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong while logging",
				"error":   err.Error(),
			})

			return
		}

		// check if their passwords are the same
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(data.Password)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to login as invalid credentials are provided",
				"error":   err.Error(),
			})

			return
		}

		// generate the jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":    existingUser.ID,
			"user_email": existingUser.Email,
			"exp":        time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to login as something went wrong while generating tokens",
				"error":   err.Error(),
			})

			return
		}

		// return the token
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully logged in as " + existingUser.Email,
			"token":   tokenString,
		})
	}
}
