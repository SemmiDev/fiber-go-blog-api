package server

import (
	"github.com/SemmiDev/fiber-go-blog/internal/app/handler"
	"github.com/SemmiDev/fiber-go-blog/internal/app/repository"
	"github.com/SemmiDev/fiber-go-blog/internal/app/service"
	"github.com/SemmiDev/fiber-go-blog/internal/db/mysql"
	"github.com/SemmiDev/fiber-go-blog/internal/db/redis"
	"github.com/SemmiDev/fiber-go-blog/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, mysqlClient mysql.Client, redisClient redis.Client) {
	// repo
	accountRepository := repository.NewAccountRepository(mysqlClient, redisClient)
	postRepository := repository.NewPostRepository(mysqlClient, redisClient)

	// service
	authService := service.NewAuthService(accountRepository)
	accountService := service.NewAccountService(accountRepository)
	postService := service.NewPostService(postRepository)

	// handler
	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	postHandler := handler.NewPostHandler(postService)

	// account apis.
	account := app.Group("api/v1/accounts")

	account.Get("/", accountHandler.List)
	account.Post("/", accountHandler.Create)
	account.Get("/:account_id", accountHandler.Get)
	account.Post("/auth", authHandler.Login)
	account.Put("/:account_id", middleware.JWTProtected(), accountHandler.Update)
	account.Put("/:account_id/password", middleware.JWTProtected(), accountHandler.UpdatePassword)
	account.Delete("/:account_id", middleware.JWTProtected(), accountHandler.Delete)

	// post apis.
	post := app.Group("api/v1/posts")

	post.Get("/", postHandler.List)
	post.Post("/", middleware.JWTProtected(), postHandler.Create)
	post.Get("/:post_id", postHandler.Get)
	post.Put("/:post_id", middleware.JWTProtected(), postHandler.Update)
	post.Delete("/:post_id", middleware.JWTProtected(), postHandler.Delete)
}