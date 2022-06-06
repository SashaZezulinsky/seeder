package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"seeder/pkg/errors"
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

	repo := mock.NewMockNodeRepository(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/nodes", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := getNodesList(repo)

	repo.EXPECT().GetNodesList(context.Background(), domain.NodeListOptions{}).Return([]*domain.Node{}, nil)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_GetNodesListFilters(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockNodeRepository(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/nodes?alive=true&age=1m&ip=127.0.0.1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := getNodesList(repo)
	alive := true
	age, _ := time.ParseDuration("1m")
	repo.EXPECT().GetNodesList(context.Background(), domain.NodeListOptions{Ip: "127.0.0.1", Alive: &alive, Age: age}).Return([]*domain.Node{}, nil)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_AuthNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"alive\":true}")

	})

	log.Println("Listening port", 8888)
	go http.ListenAndServe(":"+"8888", nil)
	time.Sleep(time.Second * 2)
	repo := mock.NewMockNodeRepository(ctrl)

	node := &domain.Node{
		IP:      "127.0.0.1:8888",
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
	handlerFunc := helloNode(repo)
	repo.EXPECT().AddNode(context.Background(), gomock.Any()).Return(nil)
	repo.EXPECT().FindNode(context.Background(), gomock.Any()).Return(errors.ErrNotFound)

	err := handlerFunc(c)
	require.Nil(t, err)
	require.NotNil(t, rec.Body.String())
}

func TestNodeHandler_AuthNodeBadData(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockNodeRepository(ctrl)

	node := &domain.Node{
		IP: "127.0.0.1",
	}
	nodeJson, _ := json.Marshal(&node)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewReader(nodeJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := helloNode(repo)

	err := handlerFunc(c)
	require.Nil(t, err)
}
