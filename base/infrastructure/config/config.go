package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

type Config struct {
	Driver        string
	Host          string
	Password      string
	User          string
	DBName        string
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func NewConfig() Config {
	return Config{
		Driver:        os.Getenv("DB_DRIVER"),
		Host:          os.Getenv("DB_HOST"),
		Password:      os.Getenv("DB_PASSWORD"),
		User:          os.Getenv("DB_USER"),
		DBName:        os.Getenv("DB_NAME"),
		Port:          os.Getenv("DB_PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
