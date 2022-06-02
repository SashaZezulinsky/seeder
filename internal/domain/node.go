//go:generate mockgen -source node.go -destination ../node/mock/mock.go -package mock

package domain

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type Node struct {
	IP      string    `json:"ip" bson:"ip"`
	Alive   bool      `json:"alive" bson:"alive"`
	Name    string    `json:"name" bson:"name"`
	Version string    `json:"version" bson:"version"`
	Client  string    `json:"client" bson:"client"`
	Date    time.Time `json:"date" bson:"date"`
}

type NodeListOptions struct {
	Alive   *bool
	Ip      string
	Client  string
	Version string
	Age     time.Duration
}

type NodeUseCase interface {
	GetNodesList(ctx context.Context, filters ...NodeListOptions) ([]*Node, error)
	AddNode(ctx context.Context, node *Node) error
	CheckNodes(ctx context.Context, interval time.Duration)
}

type NodeRepository interface {
	GetNodesList(ctx context.Context, filters ...NodeListOptions) ([]*Node, error)
	AddNode(ctx context.Context, node *Node) error
	UpdateNodeAliveStatus(ctx context.Context, node *Node, status bool) error
	FindNode(ctx context.Context, node *Node) error
}

type NodeHandlers interface {
	GetNodesList() echo.HandlerFunc
	AuthNode() echo.HandlerFunc
}
