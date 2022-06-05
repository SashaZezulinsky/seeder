package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ctxTimeout = 5
)

type Server struct {
	echo    *echo.Echo
	mongoDB *mongo.Client
	mongoCollection    string
	mongoURI           string
	mongoDatabase      string
	port               string
	nodesCheckInterval string
}

func NewServer(mongoCollection, mongoURI, mongoDatabase, port, nodesCheckInterval string, mongoDB *mongo.Client) *Server {
	return &Server{
		echo:               echo.New(),
		mongoCollection:    mongoCollection,
		mongoDB:            mongoDB,
		mongoDatabase:      mongoDatabase,
		mongoURI:           mongoURI,
		port:               port,
		nodesCheckInterval: nodesCheckInterval,
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr: ":" + s.port,
	}

	go func() {
		log.Println("Listening port", s.port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Unable to start seeder: %v", err)
		}
	}()

	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Println("Error on request", "Path", c.Path(), "Params", c.QueryParams(), "Err", err)
		s.echo.DefaultHTTPErrorHandler(err, c)
	}

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	log.Println("Server exited properly")
	return s.echo.Server.Shutdown(ctx)
}
