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
	var cstZone = time.FixedZone("GMT", 3*3600)
	time.Local = cstZone

	flag.StringVar(&port, "port", "7887", "port for client")
	flag.StringVar(&serverAddr, "server_address", "http://127.0.0.1:5000", "seeder server address")
	flag.StringVar(&clientVersion, "client.version", "v1.0.0", "client version")
	flag.StringVar(&clientName, "client.name", "testClientName", "client name")
	flag.StringVar(&clientType, "client.type", "testClient", "client type")
	flag.Parse()

	sendHelloRequest()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	fmt.Printf("Starting server at port %v\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
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
		log.Fatalf("An Error Occured %v", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(serverAddr+"/v1/nodes", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	fmt.Printf("Response: %v\n", sb)
}
