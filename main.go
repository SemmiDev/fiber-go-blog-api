package main

import (
	"fmt"
	"github.com/SemmiDev/go-blog/app"
	"log"
	"os"
)

var server = app.Server{}

func main() {
	cfg := app.LoadConfig()

	redisService, err := app.NewRedisDB(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	token := app.NewToken()
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

	app.Load(server.DB)
	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	server.Run(apiPort)
}
