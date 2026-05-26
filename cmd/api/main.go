package main

import (
	"log"

	"github.com/ayushWeb07/Notes-Rest-API/internal/config"
	"github.com/ayushWeb07/Notes-Rest-API/internal/database"
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

	// run the server router
	router.Run(":" + config.Port)
}
