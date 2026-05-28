package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ayushWeb07/Notes-Rest-API/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the headers
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Missing authorization token in header",
			})

			c.Abort()
			return
		}

		// trim bearer and get token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid authorization token in header",
			})

			c.Abort()
			return
		}

		// verify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

			// invalid signing method had been used for token generating
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Invalid signing method had been used")
			}

			// else return the jwt secret key
			return []byte(config.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token has been provided",
			})

			c.Abort()
			return
		}

		// parse token to decode the payload
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to decode payload from an invalid token",
			})

			c.Abort()
			return
		}

		// access the payload
		userId := claims["user_id"].(string)
		expiryTime := claims["exp"].(float64)

		// check if token has expired
		if time.Now().Unix() > int64(expiryTime) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token has expired. Please login again",
			})

			c.Abort()
			return
		}

		// attach the user id with the context and call next handler
		c.Set("user_id", userId)
		c.Next()
	}
}
