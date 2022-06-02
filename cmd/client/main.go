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
	"seeder/internal/domain"
	"time"
)

func main() {
	var cstZone = time.FixedZone("GMT", 3*3600)
	time.Local = cstZone

	//Get Data
	{
		//resp, err := http.Get("http://localhost:5000/v1/nodes?ip=192.168.31.138&age=1m&alive=true")
		//if err != nil {
		//	log.Fatalln(err)
		//}

		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	log.Fatalln(err)
		//}

		//sb := string(body)
		//log.Printf("List of nodes: %v", sb)
	}

	// Send Data
	var port string
	flag.StringVar(&port, "port", "7887", "port for client")
	flag.Parse()

	node := domain.Node{
		IP:      getLocalIP() + ":" + port,
		Name:    "testName2",
		Version: "v1.0.0",
		Client:  "testClient2",
	}
	postBody, err := json.Marshal(&node)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:6000/v1/nodes", "application/json", responseBody)
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
