package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConnectionUri string
	Port                  string
}

func LoadConfig() (*Config, error) {
	// load env and check if any errors
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Something went wrong while loading the env variables:", err)
	}

	// load the actual envs
	c := &Config{
		DatabaseConnectionUri: os.Getenv("DATABASE_CONNECTION_URI"),
		Port:                  os.Getenv("PORT"),
	}

	return c, err
}
