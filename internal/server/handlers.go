package server

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"

	nodeHttp "seeder/internal/node/delivery/http"
	nodeRepo "seeder/internal/node/repository/mongodb"
	nodeUsecase "seeder/internal/node/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	iRepo, err := nodeRepo.NewMongoDBNodeRepository(s.mongoDB, s.mongoDatabase, s.mongoCollection)
	if err != nil {
		return err
	}
	nodeUC := nodeUsecase.NewNodeUseCase(iRepo)
	nodeHandlers := nodeHttp.NewNodeHandler(nodeUC)

	v1 := e.Group("/v1")
	nodeGroup := v1.Group("/nodes")
	nodeHttp.MapNodeRoutes(nodeGroup, nodeHandlers)

	interval, err := time.ParseDuration(s.nodesCheckInterval)
	if err != nil {
		return err
	}
	go nodeUC.CheckNodes(context.Background(), interval)
	return nil
}
