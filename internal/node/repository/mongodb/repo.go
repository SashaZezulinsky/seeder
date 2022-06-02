package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"seeder/internal/domain"
	"time"

	"seeder/pkg/errors"
)

type mongoDBRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBNodeRepository(client *mongo.Client, database, collection string) (domain.NodeRepository, error) {
	return &mongoDBRepo{
		client:     client,
		collection: client.Database(database).Collection(collection),
	}, nil
}

func (m *mongoDBRepo) GetNodesList(ctx context.Context, filters ...domain.NodeListOptions) ([]*domain.Node, error) {
	var (
		nodes []*domain.Node
		cur   *mongo.Cursor
		err   error
	)
	query := bson.D{}

	switch {
	case len(filters) > 1:
		return nil, fmt.Errorf("filters length must be less than 2")
	case len(filters) > 0:
		if filters[0].Ip != "" {
			query = append(query, bson.E{Key: "ip", Value: filters[0].Ip})
		}
		if filters[0].Client != "" {
			query = append(query, bson.E{Key: "client", Value: filters[0].Client})
		}
		if filters[0].Age.Seconds() != 0 {
			filterTime := time.Now().Add(time.Duration(-filters[0].Age.Seconds()) * time.Second)
			query = append(query, bson.E{Key: "date", Value: bson.D{{"$lt", filterTime}}})
		}
		if filters[0].Version != "" {
			query = append(query, bson.E{Key: "version", Value: filters[0].Version})
		}
		if filters[0].Alive != nil {
			query = append(query, bson.E{Key: "alive", Value: *filters[0].Alive})
		}

		cur, err = m.collection.Aggregate(ctx, mongo.Pipeline{bson.D{{"$match", query}}})
	default:
		cur, err = m.collection.Find(ctx, bson.D{{}})
	}

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var n domain.Node
		err := cur.Decode(&n)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(nodes) == 0 {
		return nil, errors.ErrNotFound
	}

	return nodes, nil
}

func (m *mongoDBRepo) AddNode(ctx context.Context, node *domain.Node) error {
	_, err := m.collection.InsertOne(ctx, node)
	return err
}

func (m *mongoDBRepo) UpdateNodeAliveStatus(ctx context.Context, node *domain.Node, alive bool) error {
	filter := bson.M{
		"$and": []bson.M{
			{"ip": node.IP},
			{"version": node.Version},
			{"name": node.Name},
			{"client": node.Client},
		},
	}
	_, err := m.collection.UpdateOne(ctx,
		filter,
		bson.D{{"$set", bson.D{{"alive", alive}}}},
	)
	return err
}

func (m *mongoDBRepo) FindNode(ctx context.Context, node *domain.Node) error {
	filter := bson.M{
		"$and": []bson.M{
			{"ip": node.IP},
			{"version": node.Version},
			{"name": node.Name},
			{"client": node.Client},
		},
	}
	res := m.collection.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return errors.ErrNotFound
	}
	return res.Err()
}
