package main

import (
	"log"

	"github.com/ayushWeb07/Notes-Rest-API/internal/config"
	"github.com/ayushWeb07/Notes-Rest-API/internal/database"
	"github.com/ayushWeb07/Notes-Rest-API/internal/handlers"
	"github.com/gin-gonic/gin"
)

func getHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":  "Working fine as always",
		"database": "Connected",
		"status":   200,
	})
}

func main() {

	// load the config
	config, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	// connect to database and get the pool
	pool, err := database.ConnectWithDatabase(config.DatabaseConnectionUri)

	if err != nil {
		log.Fatal(err)
	}

	// close the pool finally at the end
	defer pool.Close()

	// create a router instance
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// setup routes
	router.GET("/", getHome)

	// notes routes
	router.POST("/notes", handlers.CreateNote(pool))
	router.GET("/notes", handlers.GetAllNotes(pool))
	router.GET("/notes/:id", handlers.GetNoteById(pool))
	router.PUT("/notes/:id", handlers.UpdateNote(pool))
	router.DELETE("/notes/:id", handlers.DeleteNote(pool))

	// user routes
	router.POST("/auth/register", handlers.RegisterUser(pool))
	router.GET("/users", handlers.GetUserByEmail(pool))
	router.GET("/users/:id", handlers.GetUserById(pool))

	// run the server router
	router.Run(":" + config.Port)
}
