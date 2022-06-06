package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"seeder/internal/domain"
	"seeder/pkg/errors"
	"strings"
	"time"
)

func checkNodes(ctx context.Context, repo domain.NodeRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			list, err := repo.GetNodesList(ctx)
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
					err = repo.UpdateNodeAliveStatus(ctx, node, status)
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

type PingResponse struct {
	Alive bool `json:"alive"`
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
	return pingrp.Alive, nil
}
