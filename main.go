package main

import (
	"github.com/SemmiDev/fiber-go-blog/internal/app/server"
	"github.com/SemmiDev/fiber-go-blog/internal/config"
	"github.com/SemmiDev/fiber-go-blog/internal/db/mysql"
	"github.com/SemmiDev/fiber-go-blog/internal/db/redis"
	"github.com/SemmiDev/fiber-go-blog/internal/middleware"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"log"
)

func main() {
	// Define mysql client.
	mysqlClient, err := mysql.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer mysqlClient.Close()

	// Define redis client.
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer redisClient.Close()

	// Define Fiber config.
	config := config.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Register Fiber's middleware for app.
	middleware.FiberMiddleware(app)

	// Define Routes.
	server.NewRouter(app, mysqlClient, redisClient)

	// Start server with graceful shutdown.
	server.StartServer(app)
}
