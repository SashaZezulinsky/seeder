package usecase

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"seeder/internal/domain"
	"seeder/pkg/errors"
	"strings"
	"time"
)

type nodeUsecase struct {
	nodeRepo domain.NodeRepository
}

func NewNodeUseCase(nRepo domain.NodeRepository) domain.NodeUseCase {
	return &nodeUsecase{
		nodeRepo: nRepo,
	}
}

func (n *nodeUsecase) GetNodesList(ctx context.Context, filter ...domain.NodeListOptions) ([]*domain.Node, error) {
	return n.nodeRepo.GetNodesList(ctx, filter...)
}

func (n *nodeUsecase) AddNode(ctx context.Context, node *domain.Node) error {
	node.Date = time.Now()
	node.Alive = true
	if err := n.nodeRepo.FindNode(ctx, node); err != nil {
		if err == errors.ErrNotFound {
			return n.nodeRepo.AddNode(ctx, node)
		}
		return err
	}
	return nil
}

func (n *nodeUsecase) CheckNodes(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			list, err := n.nodeRepo.GetNodesList(ctx)
			if err != nil {
				if err != errors.ErrNotFound {
					fmt.Printf("Error on getting nodes from list: %v\n", err)
				}
			}
			for _, node := range list {
				status, err := makePingRequest(node.IP)
				if err != nil && !strings.Contains(err.Error(), "connect: connection refused") {
					fmt.Printf("Error on pinging node %v:%v - %v\n", node.Client, node.IP, err)
				}
				if !status {
					fmt.Printf("Node %v:%v is not alive\n", node.Client, node.IP)
				}
				if node.Alive != status {
					err = n.nodeRepo.UpdateNodeAliveStatus(ctx, node, status)
					if err != nil {
						fmt.Printf("Error on getting nodes from list: %v\n", err)
					}
				}
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func makePingRequest(ip string) (ok bool, err error) {
	resp, err := http.Get("http://" + ip + "/ping")
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) == "pong" {
		return true, nil
	}
	return false, nil
}
