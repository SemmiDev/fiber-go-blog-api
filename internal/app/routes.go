package app

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	v1.Get("/", s.Check)

	// Auth Routes
	v1.Post("/auth/login", s.Login)
	v1.Post("/auth/logout", s.Logout)
	v1.Post("/auth/refresh", s.Refresh)

	//Users routes
	v1.Post("/users", s.CreateUser)
	v1.Get("/users", s.GetUsers)
	v1.Get("/users/:id", s.GetUser)
	v1.Put("/users/:id", s.UpdateUser)    // Require Authorization Header
	v1.Delete("/users/:id", s.DeleteUser) // Require Authorization Header

	//Posts routes
	v1.Post("/posts", s.CreatePost) // Require Authorization Header
	v1.Get("/posts", s.GetPosts)
	v1.Get("/posts/:id", s.GetPost)
	v1.Put("/posts/:id", s.UpdatePost)    // Require Authorization Header
	v1.Delete("/posts/:id", s.DeletePost) // Require Authorization Header
	v1.Get("/user_posts/:id", s.GetUserPosts)

	//Like route
	v1.Get("/likes/:id", s.GetLikes)
	v1.Post("/likes/:id", s.LikePost)     // Require Authorization Header
	v1.Delete("/likes/:id", s.UnLikePost) // Require Authorization Header

	//Comment routes
	v1.Post("/comments/:id", s.CreateComment) // Require Authorization Header
	v1.Get("/comments/:id", s.GetComments)
	v1.Put("/comments/:id", s.UpdateComment)    // Require Authorization Header
	v1.Delete("/comments/:id", s.DeleteComment) // Require Authorization Header
}

func (s *Server) Check(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("hallo 👋")
}
