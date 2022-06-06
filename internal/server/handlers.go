package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"seeder/internal/domain"
	"seeder/pkg/errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	nodeRepo "seeder/internal/node/repository/mongodb"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	repo, err := nodeRepo.NewMongoDBNodeRepository(s.mongoDB, s.mongoDatabase, s.mongoCollection)
	if err != nil {
		return err
	}

	v1 := e.Group("/v1")
	nodeGroup := v1.Group("/nodes")

	nodeGroup.GET("", getNodesList(repo))
	nodeGroup.POST("", helloNode(repo))

	interval, err := time.ParseDuration(s.nodesCheckInterval)
	if err != nil {
		return err
	}
	go checkNodes(context.Background(), repo, interval)
	return nil
}

func getNodesList(repo domain.NodeRepository) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var age time.Duration

		opts := domain.NodeListOptions{
			Ip:      c.QueryParam("ip"),
			Client:  c.QueryParam("client"),
			Version: c.QueryParam("version"),
		}

		if c.QueryParam("age") != "" {
			age, err = time.ParseDuration(c.QueryParam("age"))
			if err != nil {
				return err
			}
			opts.Age = age
		}
		if c.QueryParam("alive") != "" {
			aliveBool, err := strconv.ParseBool(c.QueryParam("alive"))
			if err != nil {
				return err
			}
			opts.Alive = &aliveBool
		}

		list, err := repo.GetNodesList(context.Background(), opts)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, list)
	}
}

func helloNode(repo domain.NodeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var node domain.Node

		nodeBytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(nodeBytes, &node)
		if err != nil {
			return err
		}
		if node.Name == "" || node.Client == "" || node.IP == "" || node.Version == "" {
			return c.JSON(http.StatusCreated, map[string]interface{}{"success": "false", "error": "please check node data"})
		}

		node.Date = time.Now()
		node.Alive = true
		if ok, err := makePingRequest(node.IP); !ok || err != nil {
			return c.JSON(http.StatusCreated, map[string]interface{}{"success": "false", "error": "cannot ping node"})
		}

		if err = repo.FindNode(context.Background(), &node); err != nil {
			if err == errors.ErrNotFound {
				err = repo.AddNode(context.Background(), &node)
				if err != nil {
					return err
				}
				return c.JSON(http.StatusCreated, map[string]interface{}{"success": "true"})
			}
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"success": "true"})
	}
}
