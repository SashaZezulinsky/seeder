package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"seeder/internal/server"
)

var (
	port               string
	mongoURI           string
	mongoCollection    string
	mongoDatabase      string
	nodesCheckInterval string
)

func main() {
	var cstZone = time.FixedZone("GMT", 3*3600)
	time.Local = cstZone
	rand.Seed(time.Now().Unix())

	log.Println("Starting api server")

	flag.StringVar(&port, "port", "5000", "port for server")
	flag.StringVar(&mongoURI, "mongo.uri", "", "mongodb URI")
	flag.StringVar(&mongoDatabase, "mongo.database", "nodes", "mongodb database")
	flag.StringVar(&mongoCollection, "mongo.collection", "nodes", "mongodb collection")
	flag.StringVar(&nodesCheckInterval, "check_interval", "30s", "interval to check if node is alive")

	flag.Parse()

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Cannot connect to MongoDB: %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	s := server.NewServer(mongoCollection, mongoURI, mongoDatabase, port, nodesCheckInterval, mongoClient)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
