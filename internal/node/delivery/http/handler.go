package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"seeder/internal/domain"
	"seeder/pkg/errors"
)

type nodeHandler struct {
	nodeUsecase domain.NodeUseCase
}

func NewNodeHandler(nUsecase domain.NodeUseCase) domain.NodeHandlers {
	return &nodeHandler{
		nodeUsecase: nUsecase,
	}
}

func (n *nodeHandler) GetNodesList() echo.HandlerFunc {
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

		list, err := n.nodeUsecase.GetNodesList(context.Background(), opts)
		if err != nil {
			if err == errors.ErrNotFound {
				return c.JSON(http.StatusOK, map[string]interface{}{"error": err.Error()})
			}
			return err
		}
		return c.JSON(http.StatusOK, list)
	}
}

func (n *nodeHandler) AuthNode() echo.HandlerFunc {
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
		err = n.nodeUsecase.AddNode(context.Background(), &node)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{"success": "true"})
	}
}
