package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrMissingURI    = errors.New("missing MONGO_DB_URI environment variable")
	ErrMissingDBName = errors.New("missing MONGO_DB_NAME environment variable")
)

// Store holds the MongoDB client and database reference.
type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

// New creates a new Store by connecting to MongoDB using environment variables
// MONGO_DB_URI and MONGO_DB_NAME.
// It returns an error if the environment variables are missing or if connection
// or ping fails.
func New(ctx context.Context) (*Store, error) {
	uri, present := os.LookupEnv("MONGO_DB_URI")
	if !present {
		return nil, ErrMissingURI
	}

	dbName, present := os.LookupEnv("MONGO_DB_NAME")
	if !present {
		return nil, ErrMissingDBName
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verifying MongoDB connection: %w", err)
	}

	return &Store{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// Close disconnects the MongoDB client.
func (store *Store) Close(ctx context.Context) error {
	return store.client.Disconnect(ctx)
}
