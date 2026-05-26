package main

import "github.com/gin-gonic/gin"

func getHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Working fine as always",
		"status":  200,
	})
}

func main() {
	// create a router instance
	router := gin.Default()

	// setup routes
	router.GET("/", getHome)

	// run the server router
	router.Run(":3000")
}
