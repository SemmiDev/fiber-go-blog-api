package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"log"
	"os"
	"os/signal"
)

type Server struct {
	DB     *gorm.DB
	Router *fiber.App
	Rd     AuthInterface
	Tk     TokenInterface
}

var errList = make(map[string]string)

func (s *Server) Initialize(DBDriver, DbUser, DbPassword, DbPort, DbHost, DbName string,
	rd AuthInterface,
	tk TokenInterface) {

	s.Rd = rd
	s.Tk = tk
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	log.Println(DBURL)
	s.DB, err = gorm.Open(DBDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", DBDriver)
		log.Fatal("This is the error connecting to postgres:", err)
	} else {
		fmt.Printf("We are connected to the %s database", DBDriver)
	}

	s.Router = fiber.New()
	s.Router.Use(cors.New())

	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	log.Fatal(s.Router.Listen(addr))
}

func (s *Server) RUnWithGracefulShutdown(addr string) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := s.Router.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := s.Router.Listen(addr); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
