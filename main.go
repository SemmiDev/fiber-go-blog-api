package main

import (
	"github.com/SemmiDev/fiber-go-blog/infrastructure/auth"
	"github.com/SemmiDev/fiber-go-blog/infrastructure/persistence"
	"github.com/SemmiDev/fiber-go-blog/interfaces"
	"github.com/SemmiDev/fiber-go-blog/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbdriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()

	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	users := interfaces.NewUsers(services.User, redisService.Auth, tk)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	middleware.CORSMiddleware(app)

	//user routes
	usersRoute := app.Group("/users")
	usersRoute.Post("/", users.SaveUser)
	usersRoute.Get("/", users.FindAllUsers)
	usersRoute.Get("/users/:user_id", users.FindUserByID)

	//Starting the application
	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}

	StartServer(app, app_port)
}

// StartServer function for starting server with a graceful shutdown.
func StartServer(app *fiber.App, port string) {

	// Create channel for idle connections.
	idleConsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConsClosed)
	}()

	// Run server.
	if err := app.Listen(":" + port); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConsClosed
}
