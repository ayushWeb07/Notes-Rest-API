package main

import (
	"log"

	"github.com/ayushWeb07/Notes-Rest-API/internal/config"
	"github.com/ayushWeb07/Notes-Rest-API/internal/database"
	"github.com/ayushWeb07/Notes-Rest-API/internal/handlers"
	"github.com/ayushWeb07/Notes-Rest-API/internal/middlewares"
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

	// auth routes
	authRouter := router.Group("/auth")
	authRouter.POST("/register", handlers.RegisterUser(pool))
	authRouter.POST("/login", handlers.LoginUser(pool, config))

	// user routes
	userRouter := router.Group("/users")
	userRouter.GET("", handlers.GetUserByEmail(pool))
	userRouter.GET("/:id", handlers.GetUserById(pool))

	// notes routes
	notesRouter := router.Group("/notes")
	notesRouter.Use(middlewares.AuthMiddleware(config)) // protect the routes by jwt

	notesRouter.POST("", handlers.CreateNote(pool))
	notesRouter.GET("", handlers.GetAllNotes(pool))
	notesRouter.GET("/:id", handlers.GetNoteById(pool))
	notesRouter.PUT("/:id", handlers.UpdateNote(pool))
	notesRouter.DELETE("/:id", handlers.DeleteNote(pool))

	// run the server router
	router.Run(":" + config.Port)
}
