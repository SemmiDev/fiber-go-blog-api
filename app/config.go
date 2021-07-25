package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBDriver      string
	DBUser        string
	DBPassword    string
	DBPort        string
	DBHost        string
	DBName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func LoadConfig() *Config {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("getting the values")
	}

	return &Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBPort:        os.Getenv("DB_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBName:        os.Getenv("DB_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
