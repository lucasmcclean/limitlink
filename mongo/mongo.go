package mongo

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var indexes = []mongo.IndexModel{
	{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("slug_idx"),
	},
	{
		Keys:    bson.D{{Key: "admin_token", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("admin_token_idx"),
	},
}

func New(ctx context.Context) (*MongoDB, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, errors.New("MONGO_URI is required\n")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return nil, errors.New("MONGO_DB_NAME is required\n")
	}

	collName := os.Getenv("MONGO_COLLECTION")
	if collName == "" {
		return nil, errors.New("MONGO_COLLECTION is required\n")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collName)

	names, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}
	log.Printf("Created indexes: %v\n", strings.Join(names, ", "))

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	repo := &MongoDB{
		client:     client,
		collection: collection,
	}

	repo.scheduleLinksCleanup(ctx)

	return repo, nil
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
