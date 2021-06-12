package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

// StartServer function for starting server with a graceful shutdown.
func StartServer(app *fiber.App) {

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
	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConsClosed
}