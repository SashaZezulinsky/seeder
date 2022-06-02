package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"seeder/internal/domain"
)

var db *mongo.Client

const MONGO_INITDB_ROOT_USERNAME = "root"
const MONGO_INITDB_ROOT_PASSWORD = "password"

var node = &domain.Node{
	IP:      "127.0.0.1",
	Name:    "name",
	Version: "version",
	Client:  "client",
	Alive:   true,
	Date:    time.Now(),
}

func TestMain(m *testing.M) {
	// Setup
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=" + MONGO_INITDB_ROOT_USERNAME,
		"MONGO_INITDB_ROOT_PASSWORD=" + MONGO_INITDB_ROOT_PASSWORD,
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env:        environmentVariables,
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = pool.Retry(func() error {
		var err error
		db, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s@localhost:%s", MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD, resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		return db.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run tests
	exitCode := m.Run()

	// Teardown
	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// Exit
	os.Exit(exitCode)
}

func TestMongoDBRepo_AddFindNode(t *testing.T) {
	repo, err := NewMongoDBNodeRepository(db, "test", "test")
	assert.Nil(t, err)

	err = repo.AddNode(context.Background(), node)
	assert.Nil(t, err)

	err = repo.FindNode(context.Background(), node)
	assert.Nil(t, err)
}

func TestMongoDBRepo_GetNodesList(t *testing.T) {
	repo, err := NewMongoDBNodeRepository(db, "test2", "test2")
	assert.Nil(t, err)

	nodes := []*domain.Node{
		{
			IP:      "127.0.0.1",
			Name:    "name",
			Version: "version",
			Client:  "client",
			Alive:   true,
			Date:    time.Now(),
		},
		{
			IP:      "127.0.0.2",
			Name:    "name",
			Version: "version",
			Client:  "client",
			Alive:   true,
			Date:    time.Now(),
		},
	}
	err = repo.AddNode(context.Background(), nodes[0])
	assert.Nil(t, err)

	err = repo.AddNode(context.Background(), nodes[1])
	assert.Nil(t, err)

	nodesRes, err := repo.GetNodesList(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, nodesRes, nodes)
}

func TestMongoDBRepo_GetNodesListFilter(t *testing.T) {
	repo, err := NewMongoDBNodeRepository(db, "test3", "test3")
	assert.Nil(t, err)

	nodes := []*domain.Node{
		{
			IP:      "127.0.0.1",
			Name:    "name",
			Version: "version",
			Client:  "client",
			Alive:   true,
			Date:    time.Now(),
		},
		{
			IP:      "127.0.0.2",
			Name:    "name",
			Version: "version",
			Client:  "client",
			Alive:   true,
			Date:    time.Now(),
		},
	}
	err = repo.AddNode(context.Background(), nodes[0])
	assert.Nil(t, err)

	err = repo.AddNode(context.Background(), nodes[1])
	assert.Nil(t, err)

	nodesRes, err := repo.GetNodesList(context.Background(), domain.NodeListOptions{Ip: "127.0.0.1"})
	assert.Nil(t, err)
	assert.Len(t, nodesRes, 1)
}

func TestMongoDBRepo_UpdateNodeAliveStatus(t *testing.T) {
	repo, err := NewMongoDBNodeRepository(db, "test4", "test4")
	assert.Nil(t, err)

	err = repo.AddNode(context.Background(), node)
	assert.Nil(t, err)

	nodesRes, err := repo.GetNodesList(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, nodesRes[0].Alive, true)

	err = repo.UpdateNodeAliveStatus(context.Background(), node, false)
	assert.Nil(t, err)

	nodesRes, err = repo.GetNodesList(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, nodesRes[0].Alive, false)
}
