package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"seeder/internal/domain"
)

var (
	port          string
	serverAddr    string
	clientName    string
	clientVersion string
	clientType    string
)

func main() {
	flag.StringVar(&port, "port", "7887", "port for client")
	flag.StringVar(&serverAddr, "server.address", "http://127.0.0.1:5000", "seeder server address")
	flag.StringVar(&clientVersion, "client.version", "v1.0.0", "client version")
	flag.StringVar(&clientName, "client.name", "testClientName", "client name")
	flag.StringVar(&clientType, "client.type", "testClient", "client type")
	flag.Parse()

	go func() {
		time.Sleep(3 * time.Second)
		sendHelloRequest()
	}()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		node := domain.Node{
			IP:      getLocalIP() + ":" + port,
			Name:    clientName,
			Version: clientVersion,
			Client:  clientType,
		}
		log.Printf("Node %s %s %s received ping from seeder", node.IP, node.Name, node.Version)
		fmt.Fprintf(w, "{\"alive\":true}")

	})

	log.Println("Listening port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func sendHelloRequest() {
	node := domain.Node{
		IP:      getLocalIP() + ":" + port,
		Name:    clientName,
		Version: clientVersion,
		Client:  clientType,
	}

	postBody, err := json.Marshal(&node)
	if err != nil {
		log.Fatalf("Unable to marshall json: %v", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(serverAddr+"/v1/nodes", "application/json", responseBody)
	if err != nil {
		log.Fatalf("Unable to POST nodes request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Node %s %s %s introduced to a seeder", node.IP, node.Name, node.Version)
	log.Printf("Seeder response: %v", string(body))
}
