package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrMissingEnvVar = errors.New("missing either MONGO_DB_URI, MONGO_DB_NAME, or both environment variables")

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
	uri, uriPresent := os.LookupEnv("MONGO_DB_URI")
	dbName, dbNamePresent := os.LookupEnv("MONGO_DB_NAME")

	if !uriPresent || !dbNamePresent {
		return nil, ErrMissingEnvVar
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
