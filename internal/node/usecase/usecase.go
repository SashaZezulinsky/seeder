package usecase

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"

	"seeder/internal/domain"
	"seeder/pkg/errors"
)

type nodeUsecase struct {
	nodeRepo domain.NodeRepository
}

type PingResponse struct {
	Alive bool `json:"alive"`
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
					log.Printf("Error on getting nodes from list: %v\n", err)
				}
			}
			for _, node := range list {
				status, err := makePingRequest(node.IP)
				if err != nil && !strings.Contains(err.Error(), "connect: connection refused") {
					log.Printf("Error on pinging node %v %v, error: %v\n", node.IP, node.Client, err)
				}
				if !status {
					log.Printf("Node %v %v is not alive\n", node.IP, node.Client)
				}
				if node.Alive != status {
					err = n.nodeRepo.UpdateNodeAliveStatus(ctx, node, status)
					if err != nil {
						log.Printf("Error on getting nodes from list: %v\n", err)
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

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var pingrp PingResponse
	if err := json.Unmarshal(body, &pingrp); err != nil {
		return false, err
	}

	if pingrp.Alive {
		return true, nil
	} else {
		return false, nil
	}
}
