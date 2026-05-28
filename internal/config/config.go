package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConnectionUri string
	Port                  string
	JwtSecretKey          string
}

func LoadConfig() (*Config, error) {
	// load env and check if any errors
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Something went wrong while loading the env variables:", err)
		return nil, err
	}

	// load the actual envs
	c := &Config{
		DatabaseConnectionUri: os.Getenv("DATABASE_CONNECTION_URI"),
		Port:                  os.Getenv("PORT"),
		JwtSecretKey:          os.Getenv("JWT_SECRET_KEY"),
	}

	return c, nil
}
