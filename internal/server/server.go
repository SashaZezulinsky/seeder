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
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo    *echo.Echo
	mongoDB *mongo.Client

	mongoCollection    string
	mongoURI           string
	mongoDatabase      string
	port               string
	nodesCheckInterval string
}

// NewServer New Server constructor
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
		Addr:           ":" + s.port,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("Server is listening on PORT: :`%s\n", s.port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalln("Error starting Server: ", err)
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

	log.Println("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
