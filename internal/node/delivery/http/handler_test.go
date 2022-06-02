package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"seeder/internal/domain"
	"seeder/internal/node/mock"
)

func TestNodeHandler_GetNodesList(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock.NewMockNodeUseCase(ctrl)

	authHandlers := NewNodeHandler(mockUC)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/nodes", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := authHandlers.GetNodesList()

	mockUC.EXPECT().GetNodesList(context.Background(), domain.NodeListOptions{}).Return([]*domain.Node{}, nil)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_GetNodesListFilters(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock.NewMockNodeUseCase(ctrl)

	authHandlers := NewNodeHandler(mockUC)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/nodes?alive=true&age=1m&ip=127.0.0.1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := authHandlers.GetNodesList()
	alive := true
	age, _ := time.ParseDuration("1m")
	mockUC.EXPECT().GetNodesList(context.Background(), domain.NodeListOptions{Ip: "127.0.0.1", Alive: &alive, Age: age}).Return([]*domain.Node{}, nil)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_AuthNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock.NewMockNodeUseCase(ctrl)

	authHandlers := NewNodeHandler(mockUC)

	node := &domain.Node{
		IP:      "127.0.0.1",
		Name:    "name",
		Version: "version",
		Client:  "client",
	}
	nodeJson, _ := json.Marshal(&node)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(nodeJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := authHandlers.AuthNode()
	mockUC.EXPECT().AddNode(context.Background(), node).Return(nil)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_AuthNodeBadData(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock.NewMockNodeUseCase(ctrl)

	authHandlers := NewNodeHandler(mockUC)

	node := &domain.Node{
		IP: "127.0.0.1",
	}
	nodeJson, _ := json.Marshal(&node)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(nodeJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := authHandlers.AuthNode()

	err := handlerFunc(c)
	require.Nil(t, err)
}
