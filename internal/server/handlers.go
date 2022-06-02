package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"

	nodeHttp "seeder/internal/node/delivery/http"
	nodeRepo "seeder/internal/node/repository/mongodb"
	nodeUsecase "seeder/internal/node/usecase"
	"seeder/pkg/utils"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	iRepo, err := nodeRepo.NewMongoDBNodeRepository(s.mongoDB, s.cfg.MongoDB.Database, s.cfg.MongoDB.Collection)
	if err != nil {
		return err
	}
	// Init useCases
	nodeUC := nodeUsecase.NewNodeUseCase(iRepo)

	// Init handlers
	nodeHandlers := nodeHttp.NewNodeHandler(nodeUC)

	v1 := e.Group("/v1")

	health := v1.Group("/health")
	nodeGroup := v1.Group("/nodes")

	nodeHttp.MapNodeRoutes(nodeGroup, nodeHandlers)

	health.GET("", func(c echo.Context) error {
		log.Printf("Health check RequestID: %s\n", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	go nodeUC.CheckNodes(context.Background(), s.cfg.Server.NodesCheckInterval)
	return nil
}