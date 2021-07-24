package main

import (
	"fmt"
	"github.com/SemmiDev/go-blog/internal/app/controllers"
	"github.com/SemmiDev/go-blog/internal/auth"
	"github.com/SemmiDev/go-blog/internal/config"
	"github.com/SemmiDev/go-blog/internal/seeder"
	"log"
	"os"
)

var server = controllers.Server{}

func main() {
	cfg := config.LoadConfig()

	redisService, err := auth.NewRedisDB(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	token := auth.NewToken()
	server.Initialize(
		cfg.DBDriver,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBPort,
		cfg.DBHost,
		cfg.DBName,
		redisService.Auth,
		token,
	)

	seeder.Load(server.DB)
	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	server.Run(apiPort)
}
