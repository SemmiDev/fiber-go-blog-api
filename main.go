package main

import (
	auth2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/auth"
	config2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/config"
	persistence2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/persistence"
	interfaces2 "github.com/SemmiDev/fiber-go-blog/base/interfaces"
	middleware2 "github.com/SemmiDev/fiber-go-blog/base/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func panicIfNeeded(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config2.NewConfig()
	services, err := persistence2.NewRepositories(cfg.Driver, cfg.User, cfg.Password, cfg.Port, cfg.Host, cfg.DBName)
	panicIfNeeded(err)
	defer services.Close()

	// warning: dev only
	services.DropTables()
	// migrate
	services.Automigrate()

	redisService, err := auth2.NewRedisDB(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	panicIfNeeded(err)
	tk := auth2.NewToken()

	authenticate := interfaces2.NewAuthenticate(services.User, redisService.Auth, tk)
	users := interfaces2.NewUsers(services.User, redisService.Auth, tk)
	posts := interfaces2.NewPosts(services.User, services.Post, redisService.Auth, tk)

	app := fiber.New()
	middleware2.CORSMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("hii 👋!")
	})

	v1 := app.Group("/api/v1")

	//authentication routes
	v1.Post("/login", authenticate.Login)
	v1.Post("/logout", authenticate.Logout)
	v1.Post("/refresh", authenticate.Refresh)

	//user routes
	v1.Post("/users", users.SaveUser)
	v1.Get("/users", users.FindAllUsers)
	v1.Get("/users/:user_id", users.FindUserByID)
	v1.Delete("/users/:user_id", users.DeleteAUser).Use(middleware2.AuthMiddleware())

	//post routes
	v1.Post("/posts", posts.CreatePost).Use(middleware2.AuthMiddleware())
	v1.Get("/posts", posts.GetPosts)
	v1.Get("/posts/:post_id", posts.GetPost)
	//v1.Put("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdatePost)
	//v1.Delete("/posts/:id", middlewares.TokenAuthMiddleware(), s.DeletePost)
	v1.Get("/user_posts/:user_id", posts.GetUserPosts)

	//Starting the application
	appPort := os.Getenv("API_PORT") //using heroku host
	if appPort == "" {
		appPort = "9090" //localhost
	}

	log.Fatal(app.Listen(":" + appPort))

	//StartServer(app, appPort)
}

// StartServer function for starting server with a graceful shutdown.
//func StartServer(app *fiber.App, port string) {
//
//	// Create channel for idle connections.
//	idleConsClosed := make(chan struct{})
//
//	go func() {
//		sigint := make(chan os.Signal, 1)
//		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
//		<-sigint
//
//		// Received an interrupt signal, shutdown.
//		if err := app.Shutdown(); err != nil {
//			// Error from closing listeners, or context timeout:
//			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
//		}
//
//		close(idleConsClosed)
//	}()
//
//	// Run server.
//	if err := app.Listen(":" + port); err != nil {
//		log.Printf("Oops... Server is not running! Reason: %v", err)
//	}
//
//	<-idleConsClosed
//}
